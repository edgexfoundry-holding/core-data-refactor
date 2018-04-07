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
package events

import (
	"os"
	"testing"

	"github.com/edgexfoundry/edgex-go/core/clients/metadataclients"
	"github.com/edgexfoundry/edgex-go/core/clients/metadataclients/mocks"
	"github.com/edgexfoundry/edgex-go/core/data/clients"
	"github.com/edgexfoundry/edgex-go/core/data/config"
	"github.com/edgexfoundry/edgex-go/core/data/log"
	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"github.com/edgexfoundry/edgex-go/support/logging-client"

	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

func TestMain(m *testing.M) {
	deviceClient = registerMockMethods()
	_, _ = clients.NewDBClient(clients.DBConfiguration{DbType: clients.MOCK})
	log.Logger = logger.NewMockClient()
	config.Configuration = &config.ConfigurationStruct{}

	os.Exit(m.Run())
}

func TestGetEventsByDevice(t *testing.T) {
	events, err := GetEventsByDevice(bson.NewObjectId().Hex(), 10)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(events) == 0 {
		t.Errorf("no events returned")
	}
}

func registerMockMethods() metadataclients.DeviceClient {
	client := &mocks.MockDeviceClient{}

	mockAddressable := models.Addressable{
		Address: "localhost",
		Name:    "Test Addressable",
		Port:    3000,
		Protocol: "http"}

	mockDeviceResultFn := func(id string) models.Device {
		device := models.Device{Id:bson.ObjectIdHex(id), Name:"Test Device", Addressable:mockAddressable}

		return device
	}
	client.On("Device", mock.MatchedBy(func(id string) bool {
		return id != "" && bson.IsObjectIdHex(id)
	})).Return(mockDeviceResultFn, nil)

	mockDeviceForNameResultFn := func(name string) models.Device {
		device := models.Device{Id:bson.NewObjectId(), Name:"Test Device", Addressable:mockAddressable}

		return device
	}
	client.On("DeviceForName", mock.MatchedBy(func(name string) bool {
		return name == "Test Device"
	})).Return(mockDeviceForNameResultFn, nil)

	return client
}