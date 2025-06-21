package model_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// MockMesheryController is a mock implementation of controllers.IMesheryController
type MockMesheryController struct {
	Name    string
	Version string
	Status  controllers.MesheryControllerStatus
	Err     error // To simulate errors for methods that return error
}

func (m *MockMesheryController) GetName() string {
	return m.Name
}

func (m *MockMesheryController) GetVersion() (string, error) {
	return m.Version, m.Err
}

func (m *MockMesheryController) GetStatus() controllers.MesheryControllerStatus {
	return m.Status
}

func (m *MockMesheryController) GetEndpointForPort(portName string) (string, error) {
	return "", m.Err
}

func (m *MockMesheryController) GetPublicEndpoint() (string, error) {
	return "", m.Err
}

func (m *MockMesheryController) Deploy(force bool) error {
	return m.Err
}

func (m *MockMesheryController) Undeploy() error {
	return m.Err
}

func TestGetMeshSyncInfo_NilBroker(t *testing.T) {
	// Create a mock meshsync controller
	mockMeshsync := &MockMesheryController{
		Name:   "meshsync",
		Status: controllers.Connected, // Initial status for meshsync
	}

	// Create a logger
	// Assuming logger.SyslogLogFormat is a valid constant in the logger package.
	// If not, this might need adjustment (e.g., to logger.JSONLogFormat or leaving it default if possible)
	// For tests, sending output to io.Discard is common.
	log, err := logger.New("test", logger.Options{
		Format:   logger.SyslogLogFormat, // This might need to be logger.JSONLogFormat or another valid one
		LogLevel: int(logrus.ErrorLevel),
		Output:   io.Discard,
	})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Call GetMeshSyncInfo with a nil broker
	statusResult := model.GetMeshSyncInfo(mockMeshsync, nil, log)

	// Assert that the status is Unknown
	assert.Equal(t, model.Status(controllers.Unknown.String()), statusResult.Status)
	assert.Equal(t, "meshsync", statusResult.Name) // Also check if the name is correctly propagated
}

func TestGetMeshSyncInfo_BrokerEndpointError(t *testing.T) {
	// Create a mock meshsync controller
	mockMeshsync := &MockMesheryController{
		Name:   "meshsync",
		Status: controllers.Connected,
	}

	// Create a mock broker controller that will return an error for GetPublicEndpoint
	mockBroker := &MockMesheryController{
		Name:   "broker",
		Status: controllers.Connected,
		Err:    fmt.Errorf("mock public endpoint error"), // Configure GetPublicEndpoint to return this error
	}

	// Create a logger
	log, err := logger.New("test-broker-error", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: int(logrus.ErrorLevel),
		Output:   io.Discard,
	})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Call GetMeshSyncInfo with the mock broker
	statusResult := model.GetMeshSyncInfo(mockMeshsync, mockBroker, log)

	// Assert that the status is Connected (without any endpoint appended)
	assert.Equal(t, model.Status(controllers.Connected.String()), statusResult.Status)
	assert.Equal(t, "meshsync", statusResult.Name) // Name should still be meshsync's name
	// Optionally, verify that no error was logged to t.Log or similar if the logger was configured for that
}

func TestGetMeshSyncInfo_NilMeshsync(t *testing.T) {
	// Create a logger
	// Assuming logger.SyslogLogFormat is a valid constant in the logger package.
	// If not, this might need adjustment (e.g., to logger.JSONLogFormat or leaving it default if possible)
	// For tests, sending output to io.Discard is common.
	log, err := logger.New("test", logger.Options{
		Format:   logger.SyslogLogFormat, // This might need to be logger.JSONLogFormat or another valid one
		LogLevel: int(logrus.ErrorLevel),
		Output:   io.Discard,
	})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Call GetMeshSyncInfo with a nil broker
	statusResult := model.GetMeshSyncInfo(nil, nil, log)

	// Assert that the status is Unknown
	assert.Equal(t, model.StatusUnknown, statusResult.Status)
}
