package task

type Repository interface {
	Save(task Task) error
	FindByID(id string) (*Task, error)
	FindAll() ([]Task, error)
	UpdateStatus(id string, status Status) error
	Delete(id string) error
}