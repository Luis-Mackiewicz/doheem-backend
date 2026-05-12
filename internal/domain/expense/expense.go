package expense

type Expense struct {
	ID          string
	Description string
	TotalAmount float64
	Category    string
	CreatedBy   string
	Splits      []Split
}

type Split struct {
	UserID string
	Amount float64
	Paid   bool
}