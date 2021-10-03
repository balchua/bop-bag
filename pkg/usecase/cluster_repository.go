package usecase

type ClusterRepository interface {
	ClusterInfo() ([]byte, error)
}
