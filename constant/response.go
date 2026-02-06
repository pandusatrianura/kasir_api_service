package constants

const (
	SuccessCode = 1000
	ErrorCode   = 2000

	ErrInvalidCategoryID      = "invalid category id"
	ErrInvalidCategoryRequest = "invalid category request"
	ErrInvalidMethod          = "invalid http method"
	ErrInvalidProductID       = "invalid product id"
	ErrInvalidProductRequest  = "invalid product request"
	ErrInvalidCheckoutRequest = "invalid checkout request"
	ErrTransactionNotFound    = "transactions not found"
	ErrRequiredDate           = "start date and end date are required"
	ErrStarDate               = "start date cannot be greater than end date"
	ErrReportRequest          = "invalid report request"
	ErrProductNotFound        = "product not found"
	ErrStockNotEnough         = "stock not enough"
	ErrStockEmpty             = "stock is empty"
	ErrCategoryNotFound       = "category not found"
)
