package repository

import "github.com/balchua/bopbag/pkg/applog"

type ClusterRepository struct {
	db  Db
	log *applog.Logger
}

func NewClusterRepository(db Db) *ClusterRepository {
	return &ClusterRepository{
		db: db,
	}
}

func (c *ClusterRepository) ClusterInfo() (string, error) {
	return c.db.GetClusterInfo()
}
