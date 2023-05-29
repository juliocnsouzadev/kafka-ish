package settings

type StorageType string

const (
	MongoDB StorageType = "MongoDB"
)

type Settings struct {
	StorageType
}
