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
	"fmt"
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
		fmt.Errorf("zero readings returned, expected at least 1")
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