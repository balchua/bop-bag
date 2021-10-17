package repository

import (
	"github.com/balchua/bopbag/pkg/applog"
)

type ClusterRepository struct {
	clusterOps ClusterOps
	log        *applog.Logger
}

func NewClusterRepository(clusterOps ClusterOps) *ClusterRepository {
	return &ClusterRepository{
		clusterOps: clusterOps,
	}
}

func (c *ClusterRepository) ClusterInfo() ([]byte, error) {
	clusterInfoInBytes, err := c.clusterOps.GetClusterInfo()
	if err != nil {
		return nil, err
	}
	return clusterInfoInBytes, nil
}

func (c *ClusterRepository) RemoveNode(address string) (string, error) {
	nodeRemoved, err := c.clusterOps.RemoveNode(address)
	if err != nil {
		return "", err
	}
	return nodeRemoved, nil
}

func (c *ClusterRepository) FindLeader() (string, error) {
	leadeNodeAddress, err := c.clusterOps.Leader()
	if err != nil {
		return "", err
	}
	return leadeNodeAddress, nil
}
