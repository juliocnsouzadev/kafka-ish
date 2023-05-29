package repository

type Repository[T any] interface {
	Insert(t T) error
	FindByTopic(topic string) ([]T, error)
	//FindById(id string) T
	//Update(t T) (T, error)
	//Delete(id int)
}
