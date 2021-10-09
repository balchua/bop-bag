package usecase

import (
	"fmt"
	"testing"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClusterRepository struct {
	mock.Mock
}

func (m *MockClusterRepository) ClusterInfo() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockClusterRepository) RemoveNode(address string) (string, error) {
	args := m.Called(address)
	return args.String(0), args.Error(1)
}

func (m *MockClusterRepository) FindLeader() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestMustSuccessfullyReturnClusterInfo(t *testing.T) {
	logger := applog.NewLogger()
	mockClusterRepo := new(MockClusterRepository)
	// setup expectations
	// var clusterInfo []domain.ClusterInfo
	// clusterInfo = append(clusterInfo, domain.ClusterInfo{ID: uint64(123456), Address: "127.0.0.1:50000", Role: 0})
	data := []byte(`[
		{
		  "ID": 3297041220608546300,
		  "Address": "norse:9000",
		  "Role": 0
		},
		{
		  "ID": 7997991560008497000,
		  "Address": "norse:9001",
		  "Role": 2
		}
	  ]`)
	mockClusterRepo.On("ClusterInfo").Return(data, nil)

	service := NewClusterService(mockClusterRepo, logger)

	response, _ := service.GetClusterInfo()
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response))
}

func TestMustFailIfUnSuccessfulCall(t *testing.T) {
	logger := applog.NewLogger()
	mockClusterRepo := new(MockClusterRepository)
	data := []byte(`[
		{
		  "ID": 3297041220608546300,
		  "Address": "norse:9000",
		  "Role": 0
		},
		{
		  "ID": 7997991560008497000,
		  "Address": "norse:9001",
		  "Role": 2
		}
	  ]`)
	mockClusterRepo.On("ClusterInfo").Return(data, fmt.Errorf("dqlite error"))

	service := NewClusterService(mockClusterRepo, logger)

	response, err := service.GetClusterInfo()
	assert.Nil(t, response)
	assert.NotNil(t, err)
}

func TestMustFailIfInvalidJSON(t *testing.T) {
	logger := applog.NewLogger()
	mockClusterRepo := new(MockClusterRepository)
	data := []byte(`[
		{
		  "ID": 3297041220608546300,
		  "Address": "norse:9000",
		  "Role": 0
		}
	  `)
	mockClusterRepo.On("ClusterInfo").Return(data, nil)

	service := NewClusterService(mockClusterRepo, logger)

	_, err := service.GetClusterInfo()
	assert.NotNil(t, err)
}

func TestReturnRemovedNode(t *testing.T) {
	logger := applog.NewLogger()
	mockClusterRepo := new(MockClusterRepository)
	data := "localhost:50000"
	mockClusterRepo.On("RemoveNode", data).Return(data, nil)

	service := NewClusterService(mockClusterRepo, logger)

	node, err := service.RemoveNode("localhost:50000")
	assert.Equal(t, "localhost:50000", node)

	assert.Nil(t, err)
}

func TestFailedRemovedNode(t *testing.T) {
	logger := applog.NewLogger()
	mockClusterRepo := new(MockClusterRepository)
	data := "localhost:50000"
	mockClusterRepo.On("RemoveNode", data).Return(data, fmt.Errorf("lost leader"))

	service := NewClusterService(mockClusterRepo, logger)

	_, err := service.RemoveNode("localhost:50000")
	assert.NotNil(t, err)
}
