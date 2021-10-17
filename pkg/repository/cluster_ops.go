package repository

import (
	"context"
)

type ClusterOps interface {
	GetClusterInfo() ([]byte, error)
	RemoveNode(address string) (string, error)
	Leader() (string, error)
	Shutdown(ctx context.Context)
}
