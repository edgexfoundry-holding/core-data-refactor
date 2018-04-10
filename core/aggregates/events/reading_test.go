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
	"testing"
	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"gopkg.in/mgo.v2/bson"
)

func TestGetAllReadings(t *testing.T) {
	readings, err := GetAllReadings()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(readings) == 0 {
		t.Errorf("zero readings returned, expected at least 1")
	}
}

func TestAddNewReading(t *testing.T) {
	getConfiguration().PersistData = true
	reading := models.Reading{Id:bson.NewObjectId(), Name:"Temperature", Device:mockParams.DeviceName}
	id, err := AddNewReading(reading)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if !bson.IsObjectIdHex(id) {
		t.Errorf("expected bson ObjectId, received %s", id)
		return
	}
}

func TestAddNewReadingWithoutPersistence(t *testing.T) {
	getConfiguration().PersistData = false
	reading := models.Reading{Id:bson.NewObjectId(), Name:"Temperature", Device:mockParams.DeviceName}
	id, err := AddNewReading(reading)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if id != "unsaved" {
		t.Errorf("unexpected readingId, received %s", id)
		return
	}
}

func TestAddNewReadingInvalidDeviceFailure(t *testing.T) {
	getConfiguration().PersistData = true
	reading := models.Reading{Id:bson.NewObjectId(), Name:"Temperature", Device:"something random"}
	_, err := AddNewReading(reading)
	if err == nil {
		t.Error("error should have been thrown")
		return
	}
}

func TestCountReadings(t *testing.T) {
	count, err := CountReadings()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count <= 0 {
		t.Errorf("unexpected value returned %v", count)
	}
}

func TestDeleteReadingById(t *testing.T) {
	id := bson.NewObjectId().Hex()
	err := DeleteEventById(id)
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetReadingsByCreateTime(t *testing.T){
	getConfiguration().ReadMaxLimit = 3
	readings, err := GetReadingsByCreateTime(mockParams.EventAgeInTicks, mockParams.EventAgeInTicks + 10, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(readings) != 2 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetReadingsByCreateTimeLimit1(t *testing.T){
	getConfiguration().ReadMaxLimit = 1
	readings, err := GetReadingsByCreateTime(mockParams.EventAgeInTicks, mockParams.EventAgeInTicks + 10, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(readings) != 1 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetEventsByDevice(t *testing.T) {
	getConfiguration().ReadMaxLimit = 3
	readings, err := GetReadingsByDevice(mockParams.DeviceName, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 2 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetEventsByDeviceLimit1(t *testing.T) {
	getConfiguration().ReadMaxLimit = 1
	readings, err := GetReadingsByDevice(mockParams.DeviceName, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 1 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetEventsByInvalidDeviceFailure(t *testing.T) {
	_, err := GetReadingsByDevice("something random", 2)
	if err == nil {
		t.Error("error should have been thrown")
		return
	}
}

func TestGetReadingById(t *testing.T) {
	id := bson.NewObjectId().Hex()
	reading, err := GetReadingById(id)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if id != reading.Id.Hex() {
		t.Errorf("mismatched reading returned, expected %s got %s", id, reading.Id.Hex())
	}
}

func TestGetReadingsByType(t *testing.T) {
	getConfiguration().ReadMaxLimit = 3
	readings, err := GetReadingsByType(mockParams.ValueDescriptorName, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 2 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetReadingsByTypeLimit1(t *testing.T) {
	getConfiguration().ReadMaxLimit = 1
	readings, err := GetReadingsByType(mockParams.ValueDescriptorName, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 1 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetReadingsByUomLabel(t *testing.T) {
	getConfiguration().ReadMaxLimit = 3
	readings, err := GetReadingsByUomLabel("C", 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 2 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetReadingsByUomLabelLimit1(t *testing.T) {
	getConfiguration().ReadMaxLimit = 1
	readings, err := GetReadingsByUomLabel("C", 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 1 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetReadingsByValueDescriptor(t *testing.T) {
	getConfiguration().ReadMaxLimit = 3
	readings, err := GetReadingsByValueDescriptor(mockParams.ValueDescriptorName, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 2 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetReadingsByValueDescriptorLimit1(t *testing.T) {
	getConfiguration().ReadMaxLimit = 1
	readings, err := GetReadingsByValueDescriptor(mockParams.ValueDescriptorName, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 1 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetReadingsByValueDescriptorLabel(t *testing.T) {
	getConfiguration().ReadMaxLimit = 3
	readings, err := GetReadingsByValueDescriptorLabel("temp", 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 2 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestGetReadingsByValueDescriptorLabelLimit1(t *testing.T) {
	getConfiguration().ReadMaxLimit = 1
	readings, err := GetReadingsByValueDescriptorLabel("temp", 2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(readings) != 1 {
		t.Errorf("unexpected reading count, received %v", len(readings))
	}
}

func TestUpdateReading(t *testing.T) {
	reading, err := getDatabase().ReadingById(mockParams.EventId.Hex()) //using this as a factory
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = UpdateReading(reading)
	if err != nil {
		t.Error(err.Error())
		return
	}
}