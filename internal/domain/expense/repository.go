package expense

type Repository interface {
	Save(expense Expense) error
	FindByID(id string) (*Expense, error)
	FindAll() ([]Expense, error)
	UpdateSplitStatus(expenseID, userID string, paid bool) error
}