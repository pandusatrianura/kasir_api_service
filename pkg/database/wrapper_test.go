package database

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

type fakeDriver struct{}

func (d *fakeDriver) Open(name string) (driver.Conn, error) {
	h, ok := harnesses.Load(name)
	if !ok {
		return nil, errors.New("dsn not found")
	}
	return &fakeConn{h: h.(*fakeHarness)}, nil
}

type fakeConn struct {
	h *fakeHarness
}

func (c *fakeConn) Prepare(query string) (driver.Stmt, error) {
	c.h.mu.Lock()
	err := c.h.prepareErr[query]
	c.h.mu.Unlock()
	if err != nil {
		return nil, err
	}
	return &fakeStmt{h: c.h, query: query}, nil
}

func (c *fakeConn) Close() error {
	return nil
}

func (c *fakeConn) Begin() (driver.Tx, error) {
	c.h.mu.Lock()
	err := c.h.beginErr
	c.h.mu.Unlock()
	if err != nil {
		return nil, err
	}
	return &fakeTx{h: c.h}, nil
}

type fakeStmt struct {
	h     *fakeHarness
	query string
}

func (s *fakeStmt) Close() error {
	return nil
}

func (s *fakeStmt) NumInput() int {
	return -1
}

func (s *fakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	return driver.RowsAffected(0), nil
}

func (s *fakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	s.h.mu.Lock()
	err := s.h.queryErr[s.query]
	data := s.h.queryData[s.query]
	s.h.mu.Unlock()
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = &fakeRowsData{}
	}
	return &fakeRows{data: data}, nil
}

type fakeTx struct {
	h *fakeHarness
}

func (tx *fakeTx) Commit() error {
	tx.h.mu.Lock()
	tx.h.commitCount++
	err := tx.h.commitErr
	tx.h.mu.Unlock()
	return err
}

func (tx *fakeTx) Rollback() error {
	tx.h.mu.Lock()
	tx.h.rollbackCount++
	err := tx.h.rollbackErr
	tx.h.mu.Unlock()
	return err
}

type fakeRowsData struct {
	columns    []string
	values     [][]driver.Value
	columnsErr error
	nextErr    error
	closeErr   error
}

type fakeRows struct {
	data *fakeRowsData
	idx  int
}

func (r *fakeRows) Columns() []string {
	if r.data.columnsErr != nil {
		return nil
	}
	return r.data.columns
}

func (r *fakeRows) Close() error {
	if r.data.closeErr != nil {
		return r.data.closeErr
	}
	return nil
}

func (r *fakeRows) Next(dest []driver.Value) error {
	if r.data.nextErr != nil {
		return r.data.nextErr
	}
	if r.idx >= len(r.data.values) {
		return io.EOF
	}
	row := r.data.values[r.idx]
	for i := range dest {
		if i < len(row) {
			dest[i] = row[i]
		} else {
			dest[i] = nil
		}
	}
	r.idx++
	return nil
}

type fakeHarness struct {
	dsn           string
	mu            sync.Mutex
	prepareErr    map[string]error
	queryErr      map[string]error
	queryData     map[string]*fakeRowsData
	beginErr      error
	commitErr     error
	rollbackErr   error
	commitCount   int
	rollbackCount int
}

var (
	registerOnce sync.Once
	harnesses    sync.Map
	dsnCounter   int64
)

func registerFakeDriver() {
	registerOnce.Do(func() {
		sql.Register("fakedb", &fakeDriver{})
	})
}

func newHarness() *fakeHarness {
	id := fmt.Sprintf("dsn-%d", atomic.AddInt64(&dsnCounter, 1))
	h := &fakeHarness{
		dsn:        id,
		prepareErr: map[string]error{},
		queryErr:   map[string]error{},
		queryData:  map[string]*fakeRowsData{},
	}
	harnesses.Store(id, h)
	return h
}

func openTestDB(t *testing.T) (*DB, *fakeHarness) {
	t.Helper()
	registerFakeDriver()
	h := newHarness()
	db, err := Open("fakedb", h.dsn)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	return db, h
}

func TestOpen(t *testing.T) {
	registerFakeDriver()
	h := newHarness()
	tests := []struct {
		name       string
		driverName string
		dsn        string
		wantErr    bool
	}{
		{name: "ok", driverName: "fakedb", dsn: h.dsn, wantErr: false},
		{name: "unknown", driverName: "nope", dsn: "x", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.driverName, tt.dsn)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if db == nil || db.DB == nil {
				t.Fatalf("db not initialized")
			}
			if !db.Logging {
				t.Fatalf("expected Logging true")
			}
		})
	}
}

func TestDBWithStmt(t *testing.T) {
	tests := []struct {
		name       string
		prepareErr error
		fnErr      error
		wantErr    bool
		wantLog    bool
		wantLogMsg string
	}{
		{name: "ok", wantErr: false, wantLog: true, wantLogMsg: "q"},
		{name: "fnerr", fnErr: errors.New("boom"), wantErr: true, wantLog: true, wantLogMsg: "boom"},
		{name: "prepare", prepareErr: errors.New("prep"), wantErr: true, wantLog: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, h := openTestDB(t)
			if tt.prepareErr != nil {
				h.prepareErr["q"] = tt.prepareErr
			}
			var logMsg string
			orig := LogFn
			LogFn = func(format string, args ...interface{}) {
				logMsg = fmt.Sprintf(format, args...)
			}
			t.Cleanup(func() { LogFn = orig })
			called := false
			err := db.WithStmt("q", func(stmt *Stmt) error {
				called = true
				return tt.fnErr
			})
			if tt.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.prepareErr != nil && called {
				t.Fatalf("fn should not be called")
			}
			if tt.wantLog {
				if logMsg == "" {
					t.Fatalf("expected log")
				}
				if tt.wantLogMsg != "" && !strings.Contains(logMsg, tt.wantLogMsg) {
					t.Fatalf("log missing message: %s", logMsg)
				}
			} else if logMsg != "" {
				t.Fatalf("unexpected log")
			}
		})
	}
}

func TestDBWithTx(t *testing.T) {
	tests := []struct {
		name         string
		beginErr     error
		fnErr        error
		commitErr    error
		wantErr      bool
		wantCommit   int
		wantRollback int
	}{
		{name: "begin", beginErr: errors.New("begin"), wantErr: true},
		{name: "fnerr", fnErr: errors.New("fn"), wantErr: true, wantRollback: 1},
		{name: "commit", commitErr: errors.New("commit"), wantErr: true, wantCommit: 1},
		{name: "ok", wantErr: false, wantCommit: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, h := openTestDB(t)
			h.beginErr = tt.beginErr
			h.commitErr = tt.commitErr
			err := db.WithTx(func(tx *Tx) error {
				return tt.fnErr
			})
			if tt.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if h.commitCount != tt.wantCommit {
				t.Fatalf("commit count: %d", h.commitCount)
			}
			if h.rollbackCount != tt.wantRollback {
				t.Fatalf("rollback count: %d", h.rollbackCount)
			}
		})
	}
}

func TestDBQueryRow(t *testing.T) {
	tests := []struct {
		name       string
		prepareErr error
		data       *fakeRowsData
		wantErr    bool
	}{
		{name: "queryerr", prepareErr: errors.New("prep"), wantErr: true},
		{name: "ok", data: &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(7)}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, h := openTestDB(t)
			if tt.prepareErr != nil {
				h.prepareErr["q"] = tt.prepareErr
			}
			if tt.data != nil {
				h.queryData["q"] = tt.data
			}
			row := db.QueryRow("q")
			var v int
			err := row.Scan(&v)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v != 7 {
				t.Fatalf("unexpected value: %d", v)
			}
		})
	}
}

func TestStmtQuery(t *testing.T) {
	tests := []struct {
		name     string
		queryErr error
		data     *fakeRowsData
		fnErr    error
		wantErr  bool
		wantRows int
	}{
		{name: "ok", data: &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(1)}, {int64(2)}}}, wantRows: 2},
		{name: "rowfnerr", data: &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(1)}}}, fnErr: errors.New("fn"), wantErr: true, wantRows: 1},
		{name: "queryerr", queryErr: errors.New("qerr"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, h := openTestDB(t)
			h.queryErr["q"] = tt.queryErr
			if tt.data != nil {
				h.queryData["q"] = tt.data
			}
			stmt, err := db.Prepare("q")
			if err != nil {
				t.Fatalf("Prepare: %v", err)
			}
			defer stmt.Close()
			wstmt := &Stmt{stmt}
			count := 0
			err = wstmt.Query(func(rows *Rows) error {
				count++
				if tt.fnErr != nil {
					return tt.fnErr
				}
				return nil
			})
			if tt.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if count != tt.wantRows {
				t.Fatalf("row count: %d", count)
			}
		})
	}
}

func TestStmtQueryRow(t *testing.T) {
	tests := []struct {
		name     string
		queryErr error
		data     *fakeRowsData
		wantErr  bool
	}{
		{name: "queryerr", queryErr: errors.New("qerr"), wantErr: true},
		{name: "ok", data: &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(3)}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, h := openTestDB(t)
			h.queryErr["q"] = tt.queryErr
			if tt.data != nil {
				h.queryData["q"] = tt.data
			}
			stmt, err := db.Prepare("q")
			if err != nil {
				t.Fatalf("Prepare: %v", err)
			}
			defer stmt.Close()
			wstmt := &Stmt{stmt}
			row := wstmt.QueryRow()
			var v int
			err = row.Scan(&v)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v != 3 {
				t.Fatalf("unexpected value: %d", v)
			}
		})
	}
}

func TestTxWithStmt(t *testing.T) {
	tests := []struct {
		name       string
		prepareErr error
		fnErr      error
		wantErr    bool
		wantLog    bool
		wantLogMsg string
	}{
		{name: "ok", wantErr: false, wantLog: true, wantLogMsg: "tx: q"},
		{name: "fnerr", fnErr: errors.New("boom"), wantErr: true, wantLog: true, wantLogMsg: "boom"},
		{name: "prepare", prepareErr: errors.New("prep"), wantErr: true, wantLog: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, h := openTestDB(t)
			if tt.prepareErr != nil {
				h.prepareErr["q"] = tt.prepareErr
			}
			sqlTx, err := db.Begin()
			if err != nil {
				t.Fatalf("Begin: %v", err)
			}
			defer sqlTx.Rollback()
			wtx := &Tx{Tx: sqlTx}
			var logMsg string
			orig := LogFn
			LogFn = func(format string, args ...interface{}) {
				logMsg = fmt.Sprintf(format, args...)
			}
			t.Cleanup(func() { LogFn = orig })
			called := false
			err = wtx.WithStmt("q", func(stmt *Stmt) error {
				called = true
				return tt.fnErr
			})
			if tt.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.prepareErr != nil && called {
				t.Fatalf("fn should not be called")
			}
			if tt.wantLog {
				if logMsg == "" {
					t.Fatalf("expected log")
				}
				if tt.wantLogMsg != "" && !strings.Contains(logMsg, tt.wantLogMsg) {
					t.Fatalf("log missing message: %s", logMsg)
				}
			} else if logMsg != "" {
				t.Fatalf("unexpected log")
			}
		})
	}
}

func TestRowError(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) *Row
		wantErr string
	}{
		{
			name: "errset",
			setup: func(t *testing.T) *Row {
				return &Row{err: errors.New("bad")}
			},
			wantErr: "bad",
		},
		{
			name: "rowserr",
			setup: func(t *testing.T) *Row {
				db, h := openTestDB(t)
				h.queryData["q"] = &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(1)}}, nextErr: errors.New("next")}
				row := db.QueryRow("q")
				row.rows.Next()
				return row
			},
			wantErr: "next",
		},
		{
			name: "none",
			setup: func(t *testing.T) *Row {
				db, h := openTestDB(t)
				h.queryData["q"] = &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(1)}}}
				return db.QueryRow("q")
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := tt.setup(t)
			if got := row.Error(); got != tt.wantErr {
				t.Fatalf("error: %q", got)
			}
		})
	}
}

func TestRowScan(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*Row, interface{})
		wantErr string
		wantVal interface{}
	}{
		{
			name: "rerr",
			setup: func(t *testing.T) (*Row, interface{}) {
				var v int
				return &Row{err: errors.New("bad")}, &v
			},
			wantErr: "bad",
		},
		{
			name: "rawbytes",
			setup: func(t *testing.T) (*Row, interface{}) {
				db, h := openTestDB(t)
				h.queryData["q"] = &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{"x"}}}
				var v sql.RawBytes
				return db.QueryRow("q"), &v
			},
			wantErr: "sql: RawBytes isn't allowed on Row.Scan",
		},
		{
			name: "norows",
			setup: func(t *testing.T) (*Row, interface{}) {
				db, h := openTestDB(t)
				h.queryData["q"] = &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{}}
				var v int
				return db.QueryRow("q"), &v
			},
			wantErr: sql.ErrNoRows.Error(),
		},
		{
			name: "scanerr",
			setup: func(t *testing.T) (*Row, interface{}) {
				db, h := openTestDB(t)
				h.queryData["q"] = &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{"x"}}}
				var v int
				return db.QueryRow("q"), &v
			},
			wantErr: "converting",
		},
		{
			name: "ok",
			setup: func(t *testing.T) (*Row, interface{}) {
				db, h := openTestDB(t)
				h.queryData["q"] = &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(9)}}}
				var v int
				return db.QueryRow("q"), &v
			},
			wantVal: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, dest := tt.setup(t)
			err := row.Scan(dest)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error mismatch: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantVal != nil {
				if got := *(dest.(*int)); got != tt.wantVal.(int) {
					t.Fatalf("value: %d", got)
				}
			}
		})
	}
}

func TestRowsScan(t *testing.T) {
	type item struct {
		A int `sql:"a"`
	}
	tests := []struct {
		name    string
		data    *fakeRowsData
		dest    interface{}
		wantErr string
		wantVal int
	}{
		{name: "maperr", data: &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(1)}}}, dest: &struct {
			B int `sql:"b"`
		}{}, wantErr: "Could not find column"},
		{name: "ok", data: &fakeRowsData{columns: []string{"a"}, values: [][]driver.Value{{int64(4)}}}, dest: &item{}, wantVal: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, h := openTestDB(t)
			h.queryData["q"] = tt.data
			rows, err := db.Query("q")
			if err != nil {
				t.Fatalf("Query: %v", err)
			}
			defer rows.Close()
			rows.Next()
			wrows := &Rows{rows}
			err = wrows.Scan(tt.dest)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error mismatch: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantVal != 0 {
				if got := tt.dest.(*item).A; got != tt.wantVal {
					t.Fatalf("value: %d", got)
				}
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		value  string
		want   int
	}{
		{name: "hit", values: []string{"a", "b"}, value: "b", want: 1},
		{name: "miss", values: []string{"a"}, value: "x", want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := find(tt.values, tt.value); got != tt.want {
				t.Fatalf("got %d", got)
			}
		})
	}
}

func TestMapColumns(t *testing.T) {
	type simple struct {
		A int    `sql:"a"`
		B string `sql:"b"`
	}
	type nested struct {
		Inner struct {
			C int `sql:"c"`
		}
	}
	type item struct {
		A int `sql:"a"`
	}
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "struct",
			run: func(t *testing.T) {
				var s simple
				dest := make([]interface{}, 2)
				j := 0
				if err := mapColumns(dest, &s, []string{"a", "b"}, "", &j); err != nil {
					t.Fatalf("mapColumns: %v", err)
				}
				*(dest[0].(*int)) = 1
				*(dest[1].(*string)) = "x"
				if s.A != 1 || s.B != "x" {
					t.Fatalf("unexpected values")
				}
			},
		},
		{
			name: "nested",
			run: func(t *testing.T) {
				var n nested
				dest := make([]interface{}, 1)
				j := 0
				if err := mapColumns(dest, &n, []string{"c"}, "", &j); err != nil {
					t.Fatalf("mapColumns: %v", err)
				}
				*(dest[0].(*int)) = 2
				if n.Inner.C != 2 {
					t.Fatalf("unexpected value")
				}
			},
		},
		{
			name: "default",
			run: func(t *testing.T) {
				v := 5
				dest := make([]interface{}, 1)
				j := 0
				if err := mapColumns(dest, &v, []string{"x"}, "", &j); err != nil {
					t.Fatalf("mapColumns: %v", err)
				}
				*(dest[0].(*int)) = 6
				if v != 6 {
					t.Fatalf("unexpected value")
				}
			},
		},
		{
			name: "missing",
			run: func(t *testing.T) {
				var s simple
				dest := make([]interface{}, 2)
				j := 0
				if err := mapColumns(dest, &s, []string{"a"}, "", &j); err == nil {
					t.Fatalf("expected error")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
