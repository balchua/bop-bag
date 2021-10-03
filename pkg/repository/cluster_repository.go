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
