package orders

type Order struct {
	ID           int     `json:"id"`
	CustomerName string  `json:"customer_name"`
	TotalPrice   float64 `json:"total_price"`
}
