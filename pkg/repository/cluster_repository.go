package repository

import (
	"github.com/balchua/bopbag/pkg/applog"
)

type ClusterRepository struct {
	db  Db
	log *applog.Logger
}

func NewClusterRepository(db Db) *ClusterRepository {
	return &ClusterRepository{
		db: db,
	}
}

func (c *ClusterRepository) ClusterInfo() ([]byte, error) {
	clusterInfoInBytes, err := c.db.GetClusterInfo()
	if err != nil {
		return nil, err
	}
	return clusterInfoInBytes, nil
}

func (c *ClusterRepository) RemoveNode(address string) (string, error) {
	nodeRemoved, err := c.db.RemoveNode(address)
	if err != nil {
		return "", err
	}
	return nodeRemoved, nil
}

func (c *ClusterRepository) FindLeader() (string, error) {
	leadeNodeAddress, err := c.db.Leader()
	if err != nil {
		return "", err
	}
	return leadeNodeAddress, nil
}
