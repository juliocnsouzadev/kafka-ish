package settings

type StorageType string

const (
	MongoDB   StorageType = "mongodb"
	FileStore StorageType = "filestore"
)

type Settings struct {
	StorageType
	TcpPort string
}
