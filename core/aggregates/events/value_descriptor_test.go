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
	"gopkg.in/mgo.v2/bson"
)

func TestAddValueDescriptor(t *testing.T) {
	id := bson.NewObjectId().Hex()
	val, err := getDatabase().ValueDescriptorById(id) //using this as a factory
	if err != nil {
		t.Error(err.Error())
		return
	}

	if id != val.Id.Hex() {
		t.Errorf("id value mismatch")
	}
}

func TestDeleteValueDescriptorByIdFailure(t *testing.T) {
	id := bson.NewObjectId().Hex()
	err := DeleteValueDescriptorById(id)
	if err == nil {
		t.Errorf("error expected due to existing readings")
		return
	}
}

func TestDeleteValueDescriptorById(t *testing.T) {
	err := DeleteValueDescriptorById(mockParams.NoReadingsKey)
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetAllValueDescriptors(t *testing.T) {
	vals, err := GetAllValueDescriptors()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(vals) == 0 {
		t.Errorf("zero length list of value descriptors")
	}
}

func TestGetValueDescriptorsByDeviceId(t *testing.T) {
	_, err := GetValueDescriptorsByDeviceId(mockParams.DeviceId.Hex())
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetValueDescriptorsByInvalidDeviceIdFailure(t *testing.T) {
	_, err := GetValueDescriptorsByDeviceId("something random")
	if err == nil {
		t.Error("error should have been thrown")
		return
	}
}

func TestGetValueDescriptorsByDeviceName(t *testing.T) {
	_, err := GetValueDescriptorsByDeviceName(mockParams.DeviceName)
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetValueDescriptorsByInvalidDeviceNameFailure(t *testing.T) {
	_, err := GetValueDescriptorsByDeviceName("something random")
	if err == nil {
		t.Error("error should have been thrown")
		return
	}
}

func TestGetValueDescriptorById(t *testing.T) {
	id := bson.NewObjectId().Hex()
	val, err := GetValueDescriptorById(id)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if id != val.Id.Hex() {
		t.Errorf("expected bson ObjectId, received %s", id)
	}
}

func TestGetValueDescriptorByName(t *testing.T) {
	val, err := GetValueDescriptorByName(mockParams.ValueDescriptorName)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if val.Name != mockParams.ValueDescriptorName {
		t.Errorf("unexpected value descriptor name %s", val.Name)
	}
}

func TestGetValueDescriptorsByLabel(t *testing.T) {
	val, err := GetValueDescriptorsByLabel("temp")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(val) == 0 {
		t.Errorf("no values returned")
	}
}

func TestGetValueDescriptorsByUomLabel(t *testing.T) {
	val, err := GetValueDescriptorsByUomLabel("C")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(val) == 0 {
		t.Errorf("no values returned")
	}
}

func TestUpdateValueDescriptor(t *testing.T) {
	id := bson.NewObjectId().Hex()
	val, err := getDatabase().ValueDescriptorById(id) //using this as a factory
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = UpdateValueDescriptor(val)
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestUpdateValueDescriptorInvalidFormatFailure(t *testing.T) {
	id := bson.NewObjectId().Hex()
	val, err := getDatabase().ValueDescriptorById(id) //using this as a factory
	if err != nil {
		t.Error(err.Error())
		return
	}
	val.Formatting = "something random"
	err = UpdateValueDescriptor(val)
	if err == nil {
		t.Errorf("error should have been thrown, format check")
	}
}

func TestUpdateValueDescriptorExistingReadingFailure(t *testing.T) {
	id := bson.NewObjectId().Hex()
	val, err := getDatabase().ValueDescriptorById(id) //using this as a factory
	if err != nil {
		t.Error(err.Error())
		return
	}
	val.Name = "something random"
	err = UpdateValueDescriptor(val)
	if err == nil {
		t.Errorf("error should have been thrown, existing readings")
	}
}

func TestValidateFormatString(t *testing.T) {
	id := bson.NewObjectId().Hex()
	val, err := getDatabase().ValueDescriptorById(id) //using this as a factory
	if err != nil {
		t.Error(err.Error())
		return
	}
	bit, err := ValidateFormatString(val)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !bit{
		t.Errorf("format string should be valid %s", val.Formatting)
	}
}

func TestValidateFormatStringInvalidFailure(t *testing.T) {
	id := bson.NewObjectId().Hex()
	val, err := getDatabase().ValueDescriptorById(id) //using this as a factory
	if err != nil {
		t.Error(err.Error())
		return
	}
	val.Formatting = "something random"
	bit, err := ValidateFormatString(val)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if bit{
		t.Errorf("format string should not be valid %s", val.Formatting)
	}
}