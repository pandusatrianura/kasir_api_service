package entity

type HealthCheck struct {
	Name      string `json:"name"`
	IsHealthy bool   `json:"is_healthy"`
}
type ReportTransaction struct {
	TotalRevenue      int64             `json:"total_revenue"`
	TotalTransactions int               `json:"total_transactions"`
	MostSoldProduct   []MostSoldProduct `json:"most_sold_product"`
}

type MostSoldProduct struct {
	Name      string `json:"name"`
	ProductID int    `json:"id,omitempty"`
	QtySold   int    `json:"quantity_sold"`
}
