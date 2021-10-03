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
