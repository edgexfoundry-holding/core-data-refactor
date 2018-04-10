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
 * @author: Trevor Conn, Dell
 * @version: 0.5.0
 *******************************************************************************/
package clients

import (
	"time"

	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"gopkg.in/mgo.v2/bson"
)

type MockParams struct {
	DeviceId bson.ObjectId
	EventId bson.ObjectId
	EventCount int
	DeviceName string
	EventAgeInTicks int64
	Origin int64
	ReadingName string
	ValueDescriptorName string
}

var mockParams *MockParams

func GetMockParams() *MockParams {
	return mockParams
}

func init() {
	mockParams = &MockParams{
		DeviceId:bson.NewObjectId(),
		EventId:bson.NewObjectId(),
		EventCount:1,
		DeviceName:"Test Device",
		EventAgeInTicks:1257894000,
		Origin:123456789,
	    ReadingName:"Temperature",
		ValueDescriptorName:"Temperature"}
}

type MockDb struct {

}

func (mc *MockDb) AddReading(r models.Reading) (bson.ObjectId, error) {
	return bson.NewObjectId(), nil
}

//DatabaseClient interface methods
func (mc *MockDb) Events() ([]models.Event, error) {
	ticks := time.Now().Unix()
	events := []models.Event{}

	evt1 := models.Event{ID:mockParams.EventId, Pushed:1, Device:mockParams.DeviceName, Created:ticks, Modified:ticks,
		Origin:mockParams.Origin, Schedule:"TestScheduleA", Event:"SampleEvent", Readings:[]models.Reading{}}

	events = append(events, evt1)
	return events, nil
}

func (mc *MockDb) AddEvent(e *models.Event) (bson.ObjectId, error){
	return e.ID, nil
}

func (mc *MockDb) UpdateEvent(e models.Event) error {
	return nil
}

func (mc *MockDb) EventById(id string) (models.Event, error){
	ticks := time.Now().Unix()

	if id == mockParams.EventId.Hex() {
		return models.Event{ID:mockParams.EventId, Pushed:1, Device:mockParams.DeviceName, Created:ticks, Modified:ticks,
			Origin:mockParams.Origin, Schedule:"TestScheduleA", Event:"SampleEvent", Readings:[]models.Reading{}}, nil
	}
	return models.Event{}, nil
}

func (mc *MockDb) EventCount() (int, error) {
	return mockParams.EventCount, nil
}

func (mc *MockDb) EventCountByDeviceId(id string) (int, error) {
	return mockParams.EventCount, nil
}

func (mc *MockDb) DeleteEventById(id string) error {
	return nil
}

func (mc *MockDb) EventsForDeviceLimit(id string, limit int) ([]models.Event, error) {
	ticks := time.Now().Unix()
	events := []models.Event{}

	evt1 := models.Event{ID:mockParams.EventId, Pushed:1, Device:mockParams.DeviceName, Created:ticks, Modified:ticks,
		Origin:mockParams.Origin, Schedule:"TestScheduleA", Event:"SampleEvent", Readings:[]models.Reading{}}

	events = append(events, evt1)
	return events, nil
}

func (mc *MockDb) EventsForDevice(id string) ([]models.Event, error){
	events := []models.Event{}
	ticks := time.Now().Unix()

	if id == mockParams.DeviceId.Hex() || id == mockParams.DeviceName {
		for len(events) < 4 {
			readings := buildListOfMockReadings()
			evt1 := models.Event{ID: mockParams.EventId, Pushed: 1, Device: mockParams.DeviceName, Created: ticks, Modified: ticks,
				Origin: mockParams.Origin, Schedule: "TestScheduleA", Event: "SampleEvent", Readings: readings}

			events = append(events, evt1)
		}
	}

	return events, nil
}

func (mc *MockDb) EventsByCreationTime(startTime, endTime int64, limit int) ([]models.Event, error) {
	ticks := time.Now().Unix()
	events := []models.Event{}

	if startTime == mockParams.EventAgeInTicks {
		if limit > 0 {
			evt1 := models.Event{ID: bson.NewObjectId(), Pushed: 1, Device: mockParams.DeviceName, Created: ticks, Modified: ticks,
				Origin: mockParams.Origin, Schedule: "TestScheduleA", Event: "SampleEvent1", Readings: []models.Reading{}}
			events = append(events, evt1)
		}
		if limit > 1 {
			evt2 := models.Event{ID: bson.NewObjectId(), Pushed: 1, Device: mockParams.DeviceName, Created: ticks, Modified: ticks,
				Origin: mockParams.Origin, Schedule: "TestScheduleA", Event: "SampleEvent2", Readings: []models.Reading{}}
			events = append(events, evt2)
		}
		if limit > 2 {
			evt3 := models.Event{ID: bson.NewObjectId(), Pushed: 1, Device: mockParams.DeviceName, Created: ticks, Modified: ticks,
				Origin: mockParams.Origin, Schedule: "TestScheduleA", Event: "SampleEvent3", Readings: []models.Reading{}}
			events = append(events, evt3)
		}
	}
	return events, nil
}

func (mc *MockDb) ReadingsByDeviceAndValueDescriptor(deviceId, valueDescriptor string, limit int) ([]models.Reading, error) {
	return []models.Reading{}, nil
}

func (mc *MockDb) EventsOlderThanAge(age int64) ([]models.Event, error) {
	events := []models.Event{}
	if age == mockParams.EventAgeInTicks {
		evt1 := models.Event{ID:bson.NewObjectId(), Pushed:1, Device:mockParams.DeviceName, Created:1257893999, Modified:1257893999,
			Origin:mockParams.Origin, Schedule:"TestScheduleA", Event:"SampleEvent", Readings:[]models.Reading{}}

		events = append(events, evt1)
	}
	return events, nil
}

func (mc *MockDb) EventsPushed() ([]models.Event, error) {
	events := []models.Event{}
	ticks := time.Now().Unix()

	for len(events) < 4 {
		readings := buildListOfMockReadings()
		evt1 := models.Event{ID: bson.NewObjectId(), Pushed: 1, Device: mockParams.DeviceName, Created: ticks, Modified: ticks,
			Origin: mockParams.Origin, Schedule: "TestScheduleA", Event: "SampleEvent", Readings: readings}

		events = append(events, evt1)
	}

	return events, nil
}

func (mc *MockDb) ScrubAllEvents() error {
	return nil
}

func (mc *MockDb) Readings() ([]models.Reading, error) {
	return buildListOfMockReadings(), nil
}

func (mc *MockDb) UpdateReading(r models.Reading) error {
	return nil
}

func (mc *MockDb) ReadingById(id string) (models.Reading, error) {
	if bson.IsObjectIdHex(id) {
		readings := buildListOfMockReadings()
		readings[0].Id = bson.ObjectIdHex(id)
		return readings[0], nil
	}
	return models.Reading{}, nil
}

func (mc *MockDb) ReadingCount() (int, error) {
	return 2, nil
}

func (mc *MockDb) DeleteReadingById(id string) error {
	return nil
}

func (mc *MockDb) ReadingsByDevice(id string, limit int) ([]models.Reading, error) {
	return buildListOfMockReadingsWithLimit(limit), nil
}

func (mc *MockDb) ReadingsByValueDescriptor(name string, limit int) ([]models.Reading, error) {
	return buildListOfMockReadingsWithLimit(limit), nil
}

func (mc *MockDb) ReadingsByValueDescriptorNames(names []string, limit int) ([]models.Reading, error) {
	return buildListOfMockReadingsWithLimit(limit), nil
}

func (mc *MockDb) ReadingsByCreationTime(start, end int64, limit int) ([]models.Reading, error) {
	return buildListOfMockReadingsWithLimit(limit), nil
}

func (mc *MockDb) AddValueDescriptor(v models.ValueDescriptor) (bson.ObjectId, error) {
	return bson.NewObjectId(), nil
}

func (mc *MockDb) ValueDescriptors() ([]models.ValueDescriptor, error) {
	return []models.ValueDescriptor{}, nil
}

func (mc *MockDb) UpdateValueDescriptor(v models.ValueDescriptor) error {
	return nil
}

func (mc *MockDb) DeleteValueDescriptorById(id string) error {
	return nil
}

func (mc *MockDb) ValueDescriptorByName(name string) (models.ValueDescriptor, error) {
	return buildValueDescriptor(), nil
}

func (mc *MockDb) ValueDescriptorsByName(names []string) ([]models.ValueDescriptor, error) {
	return buildListofMockValueDescritors(), nil
}

func (mc *MockDb) ValueDescriptorById(id string) (models.ValueDescriptor, error) {
	return buildValueDescriptor(), nil
}

func (mc *MockDb) ValueDescriptorsByUomLabel(uomLabel string) ([]models.ValueDescriptor, error) {
	return buildListofMockValueDescritors(), nil
}

func (mc *MockDb) ValueDescriptorsByLabel(label string) ([]models.ValueDescriptor, error) {
	return buildListofMockValueDescritors(), nil
}

func (mc *MockDb) ValueDescriptorsByType(t string) ([]models.ValueDescriptor, error) {
	return buildListofMockValueDescritors(), nil
}

func buildListOfMockReadings() []models.Reading {
	ticks := time.Now().Unix()
	r1 := models.Reading{Id:bson.NewObjectId(),
		Name:"Temperature",
		Value:"45",
		Origin:mockParams.Origin,
		Created:ticks,
		Modified:ticks,
		Pushed:ticks,
		Device:mockParams.DeviceName}

	r2 := models.Reading{Id:bson.NewObjectId(),
		Name:"Pressure",
		Value:"1.01325",
		Origin:mockParams.Origin,
		Created:ticks,
		Modified:ticks,
		Pushed:ticks,
		Device:mockParams.DeviceName}
	readings := []models.Reading{}
	readings = append(readings, r1, r2)
	return readings
}

func buildListOfMockReadingsWithLimit(limit int) []models.Reading {
	source := buildListOfMockReadings()
	readings := []models.Reading{}

	i := 0
	for _, r := range source {
		if i >= limit {
			break
		}
		readings  = append(readings, r)
		i++
	}
	return readings
}

func buildValueDescriptor() models.ValueDescriptor {
	ticks := time.Now().Unix()
	v := models.ValueDescriptor{Id:bson.NewObjectId(),
								Created:ticks,
								Description:"test description",
								Modified:ticks,
								Origin:mockParams.Origin,
								Name:mockParams.ValueDescriptorName,
								Min:-70,
								Max:140,
								DefaultValue:32,
								Type:"I",
								UomLabel:"C",
								Formatting:"%d",
								Labels:[]string{"temp", "room temp"}}

	return v

}

func buildListofMockValueDescritors() []models.ValueDescriptor {
	ticks := time.Now().Unix()
	vals := []models.ValueDescriptor{}

	v1 := models.ValueDescriptor{Id:bson.NewObjectId(),
		Created:ticks,
		Description:"test description",
		Modified:ticks,
		Origin:mockParams.Origin,
		Name:mockParams.ValueDescriptorName,
		Min:-70,
		Max:140,
		DefaultValue:32,
		Type:"I",
		UomLabel:"C",
		Formatting:"%d",
		Labels:[]string{"temp", "room temp"}}

	v2 := models.ValueDescriptor{Id:bson.NewObjectId(),
		Created:ticks,
		Description:"test description",
		Modified:ticks,
		Origin:mockParams.Origin,
		Name:mockParams.ValueDescriptorName,
		Min:-70,
		Max:140,
		DefaultValue:32,
		Type:"I",
		UomLabel:"C",
		Formatting:"%d",
		Labels:[]string{"temp", "room temp"}}

	v3 := models.ValueDescriptor{Id:bson.NewObjectId(),
		Created:ticks,
		Description:"test description",
		Modified:ticks,
		Origin:mockParams.Origin,
		Name:mockParams.ValueDescriptorName,
		Min:-70,
		Max:140,
		DefaultValue:32,
		Type:"I",
		UomLabel:"C",
		Formatting:"%d",
		Labels:[]string{"temp", "room temp"}}

	vals = append(vals, v1, v2, v3)
	return vals
}
