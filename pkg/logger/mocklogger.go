package logger

import (
	"fmt"

	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
	"github.com/stretchr/testify/mock"
)


var _ loggeriface.Service = &MockLogger {}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string) {
	fmt.Println("INFO: " + msg)
}

func (m *MockLogger) Error(msg string) {
	fmt.Println("ERROR: " + msg)
}

func (m *MockLogger) Fatal(msg string) {
	fmt.Println("FATAL: " + msg)
}

func (m *MockLogger) Warn(msg string) {
	fmt.Println("WARN: " + msg)
}

func (m *MockLogger) Stop() {}