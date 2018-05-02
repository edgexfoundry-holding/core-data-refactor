/*******************************************************************************
 * Copyright 2018 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * @microservice: core-data-go library
 * @author: Trevor Conn, Dell
 * @version: 0.5.0
 *******************************************************************************/
package devices

import (
	"fmt"
	"os"
	"testing"

	"github.com/edgexfoundry/edgex-go/core/aggregates/devices/mocks"
	"github.com/edgexfoundry/edgex-go/core/data/clients"
	"github.com/edgexfoundry/edgex-go/core/data/config"
	"github.com/edgexfoundry/edgex-go/core/data/log"
	"github.com/edgexfoundry/edgex-go/core/domain/models"
    "github.com/edgexfoundry/edgex-go/support/logging-client"
	"github.com/stretchr/testify/mock"

	"gopkg.in/mgo.v2/bson"
)

var mockParams *clients.MockParams

func TestMain(m *testing.M) {
	mockParams = clients.GetMockParams()
	dc = registerDeviceMethods()
	sc = registerServiceMethods()
	_, _ = clients.NewDBClient(clients.DBConfiguration{DbType: clients.MOCK})
	log.Logger = logger.NewMockClient()
	config.Configuration = &config.ConfigurationStruct{ MetaDataCheck:true, DeviceUpdateLastConnected:true}

	os.Exit(m.Run())
}

func TestSkipUpdateDeviceLastReportedConnected(t *testing.T) {
	getConfiguration().DeviceUpdateLastConnected = false
	err := updateDeviceLastReportedConnected(mockParams.DeviceName)
	if err == nil {
		t.Errorf("device update should have been skipped")
	}
	getConfiguration().DeviceUpdateLastConnected = true
}

func TestUpdateDeviceLastReportedConnectedFailNoDevice(t *testing.T) {
	err := updateDeviceLastReportedConnected("something random")
	if err == nil {
		t.Errorf("device should not have been found")
	}
}

func TestUpdateDeviceLastReportedConnectedPassById(t *testing.T) {
	err := updateDeviceLastReportedConnected(mockParams.DeviceId.Hex())
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUpdateDeviceLastReportedConnectedFailByUnexpectedId(t *testing.T) {
	err := updateDeviceLastReportedConnected(bson.NewObjectId().Hex())
	if err == nil {
		t.Errorf("device should not have been found")
	}
}

func TestUpdateDeviceLastReportedConnectedPassByName(t *testing.T) {
	err := updateDeviceLastReportedConnected(mockParams.DeviceName)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestSkipUpdateDeviceServiceLastReportedConnected(t *testing.T) {
	getConfiguration().ServiceUpdateLastConnected = false
	err := updateDeviceServiceLastReportedConnected(mockParams.DeviceName)
	if err == nil {
		t.Errorf("service update should have been skipped")
	}
	getConfiguration().ServiceUpdateLastConnected = true
}

func TestUpdateDeviceServiceLastReportedFailNoDevice(t *testing.T) {
	err := updateDeviceServiceLastReportedConnected("something random")
	if err == nil {
		t.Errorf("device should not have been found")
	}
}

func TestUpdateDeviceServiceLastReportedPassById(t *testing.T) {
	err := updateDeviceServiceLastReportedConnected(mockParams.DeviceId.Hex())
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUpdateDeviceServiceLastReportedPassByName(t *testing.T) {
	err := updateDeviceServiceLastReportedConnected(mockParams.DeviceName)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUpdateDeviceDeviceLastReportedFailByUnexpectedId(t *testing.T) {
	err := updateDeviceServiceLastReportedConnected(bson.NewObjectId().Hex())
	if err == nil {
		t.Errorf("device should not have been found")
	}
}

func registerDeviceMethods() deviceClient {
	client := &mocks.DeviceClient{}
	mockDevice := newMockDevice()

	mockDeviceResultFn := func(id string) models.Device {
		if bson.IsObjectIdHex(id) {
			mockDevice.Id = bson.ObjectIdHex(id)
			return mockDevice
		}
		return models.Device{}
	}
	client.On("Device", mockParams.DeviceId.Hex()).Return(mockDeviceResultFn, nil)
	client.On("Device", mock.AnythingOfType("string")).Return(mockDeviceResultFn, fmt.Errorf("device not found for id"))

	mockDeviceForNameResultFn := func(name string) models.Device {
		mockDevice.Id = bson.NewObjectId()

		return mockDevice
	}
	client.On("DeviceForName", mockParams.DeviceName).Return(mockDeviceForNameResultFn, nil)

	mockDeviceForNameNotFoundResultFn := func(name string) models.Device {
		return models.Device{}
	}
	client.On("DeviceForName", mock.AnythingOfType("string")).Return(mockDeviceForNameNotFoundResultFn, fmt.Errorf("no device found for name"))


	client.On("UpdateLastConnected", mockParams.DeviceId.Hex(), mock.AnythingOfType("int64")).Return(nil)
	client.On("UpdateLastConnected", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(fmt.Errorf("id is not bson ObjectIdHex"))

	client.On("UpdateLastConnectedByName", mockParams.DeviceName, mock.AnythingOfType("int64")).Return(nil)
	client.On("UpdateLastConnectedByName", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(fmt.Errorf("name is not %s", mockParams.DeviceName))

	client.On("UpdateLastReported", mockParams.DeviceId.Hex(), mock.AnythingOfType("int64")).Return(nil)
	client.On("UpdateLastReported", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(fmt.Errorf("id is not bson ObjectIdHex"))

	client.On("UpdateLastReportedByName", mockParams.DeviceName, mock.AnythingOfType("int64")).Return(nil)
	client.On("UpdateLastReportedByName", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(fmt.Errorf("name is not %s", mockParams.DeviceName))


	return client
}

func registerServiceMethods() serviceClient {
	client := &mocks.ServiceClient{}

	client.On("UpdateLastConnected", mockParams.ServiceId.Hex(), mock.AnythingOfType("int64")).Return(nil)
	client.On("UpdateLastConnected", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(fmt.Errorf("id is not bson ObjectIdHex"))

	client.On("UpdateLastReported", mockParams.ServiceId.Hex(), mock.AnythingOfType("int64")).Return(nil)
	client.On("UpdateLastReported", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(fmt.Errorf("id is not bson ObjectIdHex"))


	return client
}

func newMockDevice() models.Device {
	mockAddressable := models.Addressable{
		Address: "localhost",
		Name:    "Test Addressable",
		Port:    3000,
		Protocol: "http"}

	mockService := models.Service{
		Id:				mockParams.ServiceId,
		Name: 			"Test Service",
		OperatingState: models.Enabled,
		Labels: 		[]string{"MODBUS", "TEMP"},
		Addressable:	mockAddressable}

	mockDeviceService := models.DeviceService{
		Service:	mockService,
		AdminState: models.Unlocked}

	mockDevice := models.Device {
		Name:mockParams.DeviceName,
		Addressable:mockAddressable,
		Service:mockDeviceService}

	return mockDevice
}