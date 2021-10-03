package controller

import (
	fiber "github.com/gofiber/fiber/v2"
)

type ClusterController struct {
	service ClusterService
}

func NewClusterController(clusterService ClusterService) *ClusterController {
	return &ClusterController{
		service: clusterService,
	}
}
func (cl *ClusterController) ShowCluster(c *fiber.Ctx) error {

	clusterInfo, err := cl.service.GetClusterInfo()

	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}
	return c.JSON(clusterInfo)
}
