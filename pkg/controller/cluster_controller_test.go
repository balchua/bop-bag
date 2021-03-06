package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/balchua/bopbag/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClusterService struct {
	mock.Mock
}

func (m *MockClusterService) GetClusterInfo() ([]domain.ClusterInfo, error) {
	args := m.Called()
	return args.Get(0).([]domain.ClusterInfo), args.Error(1)
}

func (m *MockClusterService) RemoveNode(address string) (string, error) {
	args := m.Called(address)
	return args.String(0), args.Error(1)
}

func TestMustReturnClusterInfo(t *testing.T) {
	app := setupApp()

	var clusterInfo []domain.ClusterInfo
	clusterInfo = append(clusterInfo, domain.ClusterInfo{ID: uint64(1234), Address: "127.0.0.1:50000", Role: 0})
	// create an instance of our test object
	mockClusterService := new(MockClusterService)

	mockClusterService.On("GetClusterInfo").Return(clusterInfo, nil)

	controller := NewClusterController(mockClusterService)

	app.Get("/api/v1/clusterInfo", controller.ShowCluster)
	req := httptest.NewRequest("GET", "/api/v1/clusterInfo", nil)
	resp, _ := app.Test(req, 1)
	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "Show cluster info")

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "127.0.0.1:50000")
}

func TestMustFailWithBadResponse(t *testing.T) {
	app := setupApp()

	var clusterInfo []domain.ClusterInfo
	clusterInfo = append(clusterInfo, domain.ClusterInfo{ID: uint64(1234), Address: "127.0.0.1:50000", Role: 0})
	// create an instance of our test object
	mockClusterService := new(MockClusterService)

	mockClusterService.On("GetClusterInfo").Return(clusterInfo, fmt.Errorf("bad service call"))

	controller := NewClusterController(mockClusterService)

	app.Get("/api/v1/clusterInfo", controller.ShowCluster)
	req := httptest.NewRequest("GET", "/api/v1/clusterInfo", nil)
	resp, _ := app.Test(req, 1)
	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Show cluster info")

}

func TestMustTryToRemoveNode(t *testing.T) {
	app := setupApp()

	removedNode := "norse:9000"
	// create an instance of our test object
	mockClusterService := new(MockClusterService)

	mockClusterService.On("RemoveNode", removedNode).Return(removedNode, nil)

	controller := NewClusterController(mockClusterService)

	app.Delete("/api/v1/node/:nodeId", controller.RemoveNode)
	req := httptest.NewRequest("DELETE", "/api/v1/node/"+removedNode, nil)
	resp, _ := app.Test(req, 1)
	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "Snode removed")
}

func TestFailRemoveNode(t *testing.T) {
	app := setupApp()

	removedNode := "norse:9000"
	// create an instance of our test object
	mockClusterService := new(MockClusterService)

	mockClusterService.On("RemoveNode", removedNode).Return(removedNode, fmt.Errorf("unable to remove node %s", removedNode))

	controller := NewClusterController(mockClusterService)

	app.Delete("/api/v1/node/:nodeId", controller.RemoveNode)
	req := httptest.NewRequest("DELETE", "/api/v1/node/"+removedNode, nil)
	resp, _ := app.Test(req, 1)
	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "node removed")
}
