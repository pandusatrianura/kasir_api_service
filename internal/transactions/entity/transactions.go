package entity

type Transaction struct {
	ID          int    `json:"id"`
	TotalAmount int    `json:"total_amount"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type Checkout struct {
	Checkouts []CheckoutRequest `json:"checkout"`
}

type CheckoutRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutProductDetail struct {
	ID           int    `json:"product_id"`
	Name         string `json:"product_name"`
	Quantity     int    `json:"quantity"`
	Price        int    `json:"price"`
	Stock        int    `json:"stock"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type CheckoutProduct struct {
	ID           int    `json:"product_id"`
	Name         string `json:"product_name"`
	Quantity     int    `json:"quantity"`
	Subtotal     int    `json:"subtotal"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type CheckoutResponse struct {
	Transaction      Transaction       `json:"transaction"`
	CheckoutProducts []CheckoutProduct `json:"transaction_details"`
}
