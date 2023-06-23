package mock_service

import (
	models "GuppyTech/modules/v1/utilities/device/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddDevice mocks base method.
func (m *MockService) AddDevice(arg0 models.DeviceInput, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDevice", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDevice indicates an expected call of AddDevice.
func (mr *MockServiceMockRecorder) AddDevice(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDevice", reflect.TypeOf((*MockService)(nil).AddDevice), arg0, arg1)
}

// Control mocks base method.
func (m *MockService) Control(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Control", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Control indicates an expected call of Control.
func (mr *MockServiceMockRecorder) Control(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Control", reflect.TypeOf((*MockService)(nil).Control), arg0, arg1, arg2)
}

// DeleteDevice mocks base method.
func (m *MockService) DeleteDevice(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDevice", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDevice indicates an expected call of DeleteDevice.
func (mr *MockServiceMockRecorder) DeleteDevice(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDevice", reflect.TypeOf((*MockService)(nil).DeleteDevice), arg0)
}

// GetAllDevices mocks base method.
func (m *MockService) GetAllDevices(arg0 string) ([]models.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDevices", arg0)
	ret0, _ := ret[0].([]models.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDevices indicates an expected call of GetAllDevices.
func (mr *MockServiceMockRecorder) GetAllDevices(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDevices", reflect.TypeOf((*MockService)(nil).GetAllDevices), arg0)
}

// GetDatafromWebhook mocks base method.
func (m *MockService) GetDatafromWebhook(arg0, arg1 string) (models.ConnectionDat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatafromWebhook", arg0, arg1)
	ret0, _ := ret[0].(models.ConnectionDat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDatafromWebhook indicates an expected call of GetDatafromWebhook.
func (mr *MockServiceMockRecorder) GetDatafromWebhook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatafromWebhook", reflect.TypeOf((*MockService)(nil).GetDatafromWebhook), arg0, arg1)
}

// GetDeviceBrands mocks base method.
func (m *MockService) GetDeviceBrands() ([]models.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceBrands")
	ret0, _ := ret[0].([]models.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceBrands indicates an expected call of GetDeviceBrands.
func (mr *MockServiceMockRecorder) GetDeviceBrands() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceBrands", reflect.TypeOf((*MockService)(nil).GetDeviceBrands))
}

// GetDeviceById mocks base method.
func (m *MockService) GetDeviceById(arg0, arg1 string) (models.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceById", arg0, arg1)
	ret0, _ := ret[0].(models.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceById indicates an expected call of GetDeviceById.
func (mr *MockServiceMockRecorder) GetDeviceById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceById", reflect.TypeOf((*MockService)(nil).GetDeviceById), arg0, arg1)
}

// GetDeviceHistory mocks base method.
func (m *MockService) GetDeviceHistory(arg0 string) ([]models.DeviceHistory, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceHistory", arg0)
	ret0, _ := ret[0].([]models.DeviceHistory)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceHistory indicates an expected call of GetDeviceHistory.
func (mr *MockServiceMockRecorder) GetDeviceHistory(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceHistory", reflect.TypeOf((*MockService)(nil).GetDeviceHistory), arg0)
}

// GetDeviceHistoryById mocks base method.
func (m *MockService) GetDeviceHistoryById(arg0, arg1 string) ([]models.DeviceHistory, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceHistoryById", arg0, arg1)
	ret0, _ := ret[0].([]models.DeviceHistory)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceHistoryById indicates an expected call of GetDeviceHistoryById.
func (mr *MockServiceMockRecorder) GetDeviceHistoryById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceHistoryById", reflect.TypeOf((*MockService)(nil).GetDeviceHistoryById), arg0, arg1)
}

// PostControlAntares mocks base method.
func (m *MockService) PostControlAntares(arg0, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostControlAntares", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostControlAntares indicates an expected call of PostControlAntares.
func (mr *MockServiceMockRecorder) PostControlAntares(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostControlAntares", reflect.TypeOf((*MockService)(nil).PostControlAntares), arg0, arg1, arg2, arg3)
}

// UpdateDeviceById mocks base method.
func (m *MockService) UpdateDeviceById(arg0 models.DeviceInput, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDeviceById", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDeviceById indicates an expected call of UpdateDeviceById.
func (mr *MockServiceMockRecorder) UpdateDeviceById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDeviceById", reflect.TypeOf((*MockService)(nil).UpdateDeviceById), arg0, arg1)
}
