package entity

type HealthCheck struct {
	Name      string `json:"name"`
	IsHealthy bool   `json:"is_healthy"`
}
type ReportTransaction struct {
	TotalRevenue      int64             `json:"total_revenue"`
	TotalTransactions int               `json:"total_transaksi"`
	MostSoldProduct   []MostSoldProduct `json:"produk_terlaris"`
}

type MostSoldProduct struct {
	Name      string `json:"nama"`
	ProductID int    `json:"id,omitempty"`
	QtySold   int    `json:"qty_terjual"`
}
