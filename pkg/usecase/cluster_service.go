package usecase

import (
	"encoding/json"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/domain"
)

type ClusterService struct {
	clusterRepo domain.ClusterRepository
	logger      *applog.Logger
}

func NewClusterService(clusterRepo domain.ClusterRepository, logger *applog.Logger) *ClusterService {
	return &ClusterService{
		clusterRepo: clusterRepo,
		logger:      logger,
	}
}

func (c *ClusterService) GetClusterInfo() ([]domain.ClusterInfo, error) {
	clusterInfoInBytes, err := c.clusterRepo.ClusterInfo()
	c.logger.Log.Sugar().Infof("cluster info retrieved")
	leader, err := c.clusterRepo.FindLeader()
	c.logger.Log.Sugar().Infof("leader found")
	if err != nil {
		return nil, err
	}
	clusterInfo := make([]domain.ClusterInfo, 0)

	if err := json.Unmarshal(clusterInfoInBytes, &clusterInfo); err != nil {
		return nil, err
	}
	for i, info := range clusterInfo {
		c.logger.Log.Sugar().Infof("Nodes %s, leader %s", info.Address, leader)
		if info.Address == leader {
			clusterInfo[i].Leader = true
		}
	}
	return clusterInfo, nil
}

func (c *ClusterService) RemoveNode(address string) (string, error) {
	removedNode, err := c.clusterRepo.RemoveNode(address)
	if err != nil {
		return "", err
	}

	return removedNode, nil
}
