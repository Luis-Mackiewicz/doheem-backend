package expense

type ServicePort interface {
	CreateExpense(expense Expense) error
	PaySplit(expenseID, userID string) error
	GetSummary() ([]Split, error)
}