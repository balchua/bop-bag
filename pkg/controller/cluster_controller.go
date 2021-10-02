package controller

import (
	"github.com/balchua/bopbag/pkg/repository"
	fiber "github.com/gofiber/fiber/v2"
)

type ClusterController struct {
	repo *repository.ClusterRepository
}

func NewClusterController(clusterRepo *repository.ClusterRepository) *ClusterController {
	return &ClusterController{
		repo: clusterRepo,
	}
}
func (cl *ClusterController) ShowCluster(c *fiber.Ctx) error {

	clusterInfo, err := cl.repo.ClusterInfo()

	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}
	return c.Send([]byte(clusterInfo))
}
