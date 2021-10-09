package usecase

type ClusterRepository interface {
	ClusterInfo() ([]byte, error)
	RemoveNode(address string) (string, error)
	FindLeader() (string, error)
}
