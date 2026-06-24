package expense

import (
	"context"
	"errors"
	"time"
)

type Expense struct {
	ID               string
	GroupID          string
	Description      string
	Amount           float64
	CategoryID       string
	CompetenceDate   time.Time
	DueDate          time.Time
	PaidBy           string
	SplitMode        string
	Installments     int32
	FirstDueDate     *time.Time
	IsFixed          bool
	ParentExpenseID  *string
	InstallmentIndex *int32
	InstallmentTotal *int32
	CreatedBy        *string
	FixedSourceID    *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateExpenseParams struct {
	GroupID          string
	Description      string
	Amount           float64
	CategoryID       string
	CompetenceDate   time.Time
	DueDate          time.Time
	PaidBy           string
	SplitMode        string
	Installments     int32
	FirstDueDate     *time.Time
	IsFixed          bool
	ParentExpenseID  *string
	InstallmentIndex *int32
	InstallmentTotal *int32
	CreatedBy        *string
	FixedSourceID    *string
}

type UpdateExpenseParams struct {
	Description     *string
	Amount          *float64
	CompetenceDate  *time.Time
	DueDate         *time.Time
	CategoryID      *string
	SplitMode       *string
	SelectedUserIDs []string
	Splits          []CreateExpenseSplitParams
}

type ExpenseCategory struct {
	ID        string
	Slug      string
	Label     string
	CreatedAt time.Time
}

type ExpenseSplit struct {
	ID              string
	ExpenseID       string
	UserID          string
	Amount          float64
	IsPaid          bool
	PaidAt          *time.Time
	CreatedAt       time.Time
	ReceiptData     *string
	ReceiptType     *string
	ReceiptFileName *string
}

type ExpenseSplitWithUser struct {
	ExpenseSplit
	UserName  string
	UserEmail string
}

type UserBalance struct {
	TotalOwed float64
	TotalPaid float64
}

type ExpenseRepository interface {
	GetByID(ctx context.Context, id string) (Expense, error)
	ListByGroup(ctx context.Context, groupID string, dateFrom, dateTo *time.Time, limit, offset int32) ([]Expense, error)
	ListByGroupFull(ctx context.Context, groupID string, dateFrom, dateTo *time.Time, search string, myUserID *string, limit, offset int32) ([]Expense, int64, error)
	CountByGroup(ctx context.Context, groupID string, dateFrom, dateTo *time.Time) (int, error)
	ListByUser(ctx context.Context, userID string) ([]Expense, error)
	ListByCategory(ctx context.Context, categoryID string) ([]Expense, error)
	ListByParent(ctx context.Context, parentID string) ([]Expense, error)
	ListFixedOrigins(ctx context.Context) ([]Expense, error)
	CountCloneByMonth(ctx context.Context, fixedSourceID string, dateFrom, dateTo time.Time) (int64, error)
	Create(ctx context.Context, params CreateExpenseParams) (Expense, error)
	Update(ctx context.Context, id string, params UpdateExpenseParams) (Expense, error)
	Delete(ctx context.Context, id string) error
	DeleteByParent(ctx context.Context, parentID string) error
	GetTotalByGroup(ctx context.Context, groupID string) (float64, error)
}

type ExpenseCategoryRepository interface {
	GetByID(ctx context.Context, id string) (ExpenseCategory, error)
	ListAll(ctx context.Context) ([]ExpenseCategory, error)
	Create(ctx context.Context, slug, label string) (ExpenseCategory, error)
	Update(ctx context.Context, id, slug, label string) (ExpenseCategory, error)
	Delete(ctx context.Context, id string) error
}

type ExpenseSplitRepository interface {
	GetByID(ctx context.Context, id string) (ExpenseSplit, error)
	ListByExpense(ctx context.Context, expenseID string) ([]ExpenseSplitWithUser, error)
	ListByExpenseIDs(ctx context.Context, expenseIDs []string) ([]ExpenseSplitWithUser, error)
	ListByUser(ctx context.Context, userID, groupID string) ([]ExpenseSplit, error)
	Create(ctx context.Context, expenseID, userID string, amount float64) (ExpenseSplit, error)
	CreateMany(ctx context.Context, expenseID string, splits []CreateExpenseSplitParams) (int64, error)
	MarkAsPaid(ctx context.Context, id string, receiptData, receiptType, receiptFileName *string) error
	HasPaidSplits(ctx context.Context, expenseID string) (bool, error)
	DeleteByExpense(ctx context.Context, expenseID string) error
	GetUserBalance(ctx context.Context, userID, groupID string) (UserBalance, error)
}

type CreateExpenseSplitParams struct {
	UserID string
	Amount float64
}

type CalculateSplitsParams struct {
	GroupID         string
	Amount          float64
	SplitMode       string
	SelectedUserIDs []string
}

type CreateExpenseWithSplitsParams struct {
	Expense    CreateExpenseParams
	Splits     []CreateExpenseSplitParams
	CalcParams CalculateSplitsParams
}

var (
	ErrExpenseNotFound             = errors.New("despesa não encontrada")
	ErrInvalidSplitTotal           = errors.New("os valores das divisões devem ser iguais ao valor total")
	ErrCategoryNotFound            = errors.New("categoria não encontrada")
	ErrForbidden                   = errors.New("acesso proibido")
	ErrCannotDeleteWithPaidSplits  = errors.New("não é possível excluir uma despesa com divisões pagas")
	ErrCannotEditWithPaidSplits    = errors.New("não é possível editar uma despesa com divisões pagas")
	ErrInvalidSplitMode            = errors.New("modo de divisão inválido")
	ErrNoSelectedMembers           = errors.New("selecione pelo menos 2 membros para alguma divisão")
	ErrSplitAlreadyPaid            = errors.New("divisão já paga")
	ErrCannotEditInstallmentChild  = errors.New("não é possível editar uma despesa filha de parcelamento diretamente")
	ErrCannotEditInstallmentParent = errors.New("não é possível editar uma despesa pai de parcelamento diretamente")
	ErrFixedWithInstallments       = errors.New("não é possível criar uma despesa fixa com parcelamentos")
)
