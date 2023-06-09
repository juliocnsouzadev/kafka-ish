package repository

type Repository[T any] interface {
	Insert(t T) error
	FindByTopic(topic string) ([]T, error)
}
