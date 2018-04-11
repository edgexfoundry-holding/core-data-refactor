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
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/edgexfoundry/edgex-go/core/aggregates"
	"github.com/edgexfoundry/edgex-go/core/clients/metadataclients"
	"github.com/edgexfoundry/edgex-go/core/clients/metadataclients/mocks"
	"github.com/edgexfoundry/edgex-go/core/data/clients"
	"github.com/edgexfoundry/edgex-go/core/data/config"
	"github.com/edgexfoundry/edgex-go/core/data/log"
	"github.com/edgexfoundry/edgex-go/core/data/messaging"
	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"github.com/edgexfoundry/edgex-go/support/logging-client"

	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

var mockParams *clients.MockParams

func TestMain(m *testing.M) {
	mockParams = clients.GetMockParams()
	deviceClient = registerMockMethods()
	_, _ = clients.NewDBClient(clients.DBConfiguration{DbType: clients.MOCK})
	_ = messaging.NewMQPublisher("", messaging.MOCK)
	log.Logger = logger.NewMockClient()
	config.Configuration = &config.ConfigurationStruct{ MetaDataCheck:true}

	os.Exit(m.Run())
}

func TestGetEventsByDeviceId(t *testing.T) {
	events, err := GetEventsByDevice(mockParams.DeviceId.Hex(), 10)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(events) == 0 {
		t.Errorf("no events returned")
	}
}

func TestGetEventsByDeviceName(t *testing.T) {
	events, err := GetEventsByDevice(mockParams.DeviceName, 10)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(events) == 0 {
		t.Errorf("no events returned")
	}
}

func TestGetEventsByDeviceNameForFailure(t *testing.T) {
	_, err := GetEventsByDevice("something random", 10)

	if err == nil {
		t.Error("error should have been thrown")
		return
	}
}

func TestEventsCount(t *testing.T) {
	count, err := CountEvents()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count != mockParams.EventCount {
		t.Error(fmt.Errorf("event count %v does not match expected mock value %v", count, mockParams.EventCount))
		return
	}
}

func TestEventsCountByDeviceId(t *testing.T) {
	count, err := CountByDevice(mockParams.DeviceId.Hex())
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count != mockParams.EventCount {
		t.Error(fmt.Errorf("event count %v does not match expected mock value %v", count, mockParams.EventCount))
		return
	}
}

func TestEventsCountByDeviceName(t *testing.T) {
	count, err := CountByDevice(mockParams.DeviceName)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count != mockParams.EventCount {
		t.Error(fmt.Errorf("event count %v does not match expected mock value %v", count, mockParams.EventCount))
		return
	}
}

func TestEventsCountByDeviceNameForFailure(t *testing.T) {
	_, err := CountByDevice("something random")
	if err == nil {
		t.Error("error should have been thrown")
		return
	}
}

func TestDeleteByAge(t *testing.T) {
	del, err := DeleteByAge(mockParams.EventAgeInTicks)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if del == 0 {
		t.Errorf("no events were deleted, should be at least 1")
		return
	}
}

func TestDeleteByAgeForFailure(t *testing.T) {
	del, err := DeleteByAge(mockParams.EventAgeInTicks - 10)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if del != 0 {
		t.Errorf("events were deleted, should have been zero")
		return
	}
}

func TestDeleteByDeviceId(t *testing.T) {
	del, err := DeleteByDevice(mockParams.DeviceId.Hex())
	if err != nil {
		t.Error(err.Error())
		return
	}

	if del == 0 {
		t.Errorf("no events were deleted, should be at least 1")
		return
	}
}

func TestDeleteByDeviceName(t *testing.T) {
	del, err := DeleteByDevice(mockParams.DeviceName)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if del == 0 {
		t.Errorf("no events were deleted, should be at least 1")
		return
	}
}

func TestDeleteByDeviceNameForFailure(t *testing.T) {
	_, err := DeleteByDevice("something random")
	if err == nil {
		t.Error("error should have been thrown")
		return
	}
}

func TestDeleteEventById(t *testing.T) {
	err := DeleteEventById(mockParams.EventId.Hex())
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetAllEvents(t *testing.T) {
	events, err := GetAllEvents()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(events) == 0 {
		t.Errorf("no events returned, expected at least 1")
		return
	}
}

func TestGetEventsByCreateTime(t *testing.T) {
	getConfiguration().ReadMaxLimit = 2
	events, err := GetEventsByCreateTime(mockParams.EventAgeInTicks, mockParams.EventAgeInTicks + 10, 3)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(events) != 2 {
		t.Errorf("unexpected number of events %v, expected 2", len(events))
		return
	}
}

func TestGetEventById(t *testing.T) {
	event, err := GetEventById(mockParams.EventId.Hex())
	if err != nil {
		t.Error(err.Error())
		return
	}
	testEventWithoutReadings(event, t)
}

func TestGetReadingsByDeviceAndValueDescriptor(t *testing.T) {
	getConfiguration().ReadMaxLimit = 2
	readings, err := GetReadingsByDeviceAndValueDescriptor(mockParams.DeviceName, mockParams.ReadingName, 3)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) == 0 {
		t.Errorf("no readings returned, expected at least 1")
		return
	}
	if len(readings) > getConfiguration().ReadMaxLimit {
		t.Errorf("readings returned %v exceeded ReadMaxLimit %v", len(readings), getConfiguration().ReadMaxLimit)
		return
	}
	for _, r := range readings {
		if r.Name != mockParams.ReadingName {
			t.Errorf("unexpected reading name returned %s", r.Name)
			return
		}
	}
}

func TestGetReadingsByDeviceAndValueDescriptorDeviceNotFound(t *testing.T) {
	getConfiguration().ReadMaxLimit = 2
	_, err := GetReadingsByDeviceAndValueDescriptor("something random", mockParams.ReadingName, 3)
	if err == nil {
		t.Errorf("error was expected when supplying bogus device name")
		return
	}
}

func TestPurge(t *testing.T) {
	err := Purge()
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestPurgeIfPublished(t *testing.T) {
	count, err := PurgeIfPublished()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count == 0 {
		t.Errorf("%v events deleted, expected at least 1", count)
		return
	}
}

func TestTouchInvalidIdFailure(t *testing.T) {
	err := Touch("test")
	if err == nil {
		t.Errorf("error expected for non-bson object id eventId")
		return
	}
}

func TestTouch(t *testing.T) {
	err := Touch(mockParams.EventId.Hex())
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestUpdateEvent(t *testing.T) {
	event, _ := getDatabase().EventById(mockParams.EventId.Hex()) //using this as a factory
	err := UpdateEvent(event)
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestUpdateEventInvalidDeviceFailure(t *testing.T) {
	event, _ := getDatabase().EventById(mockParams.EventId.Hex()) //using this as a factory
	event.Device = "something random"
	err := UpdateEvent(event)
	if err == nil {
		t.Errorf("expected error for invalid device name")
		return
	}
}

func TestAddNewEvent(t *testing.T) {
	getConfiguration().PersistData = true
	event, _ := getDatabase().EventById(mockParams.EventId.Hex()) //using this as a factory
	//wire up handlers to listen for device events
	bitEvents := make([]bool, 2)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go handleDomainEvents(bitEvents, &wg, t)

	id, err := AddNewEvent(event)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if id != event.ID.Hex() {
		t.Errorf("mismatched ID after save %s, expected %s", id, event.ID.Hex())
		return
	}

	wg.Wait()
	for i, val := range bitEvents {
		if !val {
			t.Errorf("event not received in timely fashion, index %v", i)
			return
		}
	}
}

func TestAddNewEventWithoutPersistence(t *testing.T) {
	getConfiguration().PersistData = false
	event, _ := getDatabase().EventById(mockParams.EventId.Hex()) //using this as a factory
	//wire up handlers to listen for device events
	bitEvents := make([]bool, 2)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go handleDomainEvents(bitEvents, &wg, t)

	id, err := AddNewEvent(event)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if id != "unsaved" {
		t.Errorf("mismatched ID after save %s, expected %s", id, "unsaved")
		return
	}

	wg.Wait()
	for i, val := range bitEvents {
		if !val {
			t.Errorf("event not received in timely fashion, index %v", i)
			return
		}
	}
}

func handleDomainEvents(bitEvents []bool, wait *sync.WaitGroup, t *testing.T) {
		until := time.Now().Add(250 * time.Millisecond) //Kill this loop after quarter second.
		for time.Now().Before(until) {
			select {
			case evt := <- EventAggregateEvents:
				switch evt.(type) {
				case aggregates.DeviceLastReported:
					fmt.Println("TestAddNewEventWithoutPersistence aggregates.DeviceLastReported")
					e := evt.(aggregates.DeviceLastReported)
					if e.DeviceName != mockParams.DeviceName {
						t.Errorf("DeviceLastReported name mistmatch %s", e.DeviceName)
						return
					}
					setEventBit(0, true, bitEvents)
					break;
				case aggregates.DeviceServiceLastReported:
					fmt.Println("TestAddNewEventWithoutPersistence aggregates.DeviceServiceLastReported")
					e := evt.(aggregates.DeviceServiceLastReported)
					if e.DeviceName != mockParams.DeviceName {
						t.Errorf("DeviceLastReported name mistmatch %s", e.DeviceName)
						return
					}
					setEventBit(1, true, bitEvents)
					break;
				}
			default:
				//	Without a default case in here, the select block will hang.
			}
		}
		wait.Done()
}

func setEventBit(index int, value bool, source []bool) {
	for i, oldVal := range source {
		if i == index {
			source[i] = value
		} else {
			source[i] = oldVal
		}
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
		if bson.IsObjectIdHex(id) {
			return models.Device{Id:bson.ObjectIdHex(id), Name:mockParams.DeviceName, Addressable:mockAddressable}
		}
		return models.Device{}
	}
	client.On("Device", mock.MatchedBy(func(id string) bool {
		return bson.IsObjectIdHex(id)
	})).Return(mockDeviceResultFn, nil)
    client.On("Device", mock.MatchedBy(func(id string) bool {
    	return !bson.IsObjectIdHex(id)
	})).Return(mockDeviceResultFn, fmt.Errorf("id is not bson ObjectIdHex"))

	mockDeviceForNameResultFn := func(name string) models.Device {
		device := models.Device{Id:bson.NewObjectId(), Name:name, Addressable:mockAddressable}

		return device
	}
	client.On("DeviceForName", mock.MatchedBy(func(name string) bool {
		return name == mockParams.DeviceName
	})).Return(mockDeviceForNameResultFn, nil)
	client.On("DeviceForName", mock.MatchedBy(func(name string) bool {
		return name != mockParams.DeviceName
	})).Return(mockDeviceForNameResultFn, fmt.Errorf("no device found for name"))

	return client
}

func testEventWithoutReadings(event models.Event, t *testing.T) {
	if event.ID.Hex() != mockParams.EventId.Hex() {
		t.Error("eventId mismatch. expected " + mockParams.EventId.Hex() + " received " + event.ID.Hex())
	}

	if event.Device != mockParams.DeviceName {
		t.Error("device mismatch. expected " + mockParams.DeviceName + " received " + event.Device)
	}

	if event.Origin != mockParams.Origin {
		t.Error("origin mismatch. expected " + strconv.FormatInt(mockParams.Origin, 10) + " received " + strconv.FormatInt(event.Origin, 10))
	}
}