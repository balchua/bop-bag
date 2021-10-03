package usecase

import (
	"encoding/json"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/domain"
)

type ClusterService struct {
	clusterRepo ClusterRepository
	logger      *applog.Logger
}

func NewClusterService(clusterRepo ClusterRepository, logger *applog.Logger) *ClusterService {
	return &ClusterService{
		clusterRepo: clusterRepo,
		logger:      logger,
	}
}

func (c *ClusterService) GetClusterInfo() ([]domain.ClusterInfo, error) {
	clusterInfoInBytes, err := c.clusterRepo.ClusterInfo()
	if err != nil {
		return nil, err
	}
	clusterInfo := make([]domain.ClusterInfo, 0)

	if err := json.Unmarshal(clusterInfoInBytes, &clusterInfo); err != nil {
		return nil, err
	}
	return clusterInfo, nil
}
