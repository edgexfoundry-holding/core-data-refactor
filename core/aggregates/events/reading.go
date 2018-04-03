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

	"github.com/edgexfoundry/edgex-go/core/domain/models" //for now
	"github.com/edgexfoundry/edgex-go/core/domain/errs"
)

func GetAllReadings() ([]models.Reading, error) {
	r, err := getDatabase().Readings()
	if err != nil {
		getLogger().Error(err.Error())
		return nil, fmt.Errorf(err.Error())
	}

	return r, err
}

func AddNewReading(reading models.Reading) (string, error) {
	// Check the value descriptor
	_, err := getDatabase().ValueDescriptorByName(reading.Name)
	if err != nil {
		var msg string
		if err == errs.ErrNotFound {
			msg = fmt.Sprintf("value descriptor not found: %s", reading.Name)
		} else {
			msg = err.Error()
		}
		getLogger().Error(msg)
		return "", err
	}

	// Check device
	if reading.Device != "" {
		// Try by name
		d, err := mdc.DeviceForName(reading.Device)
		// Try by ID
		if err != nil {
			d, err = mdc.Device(reading.Device)
			if err != nil {
				getLogger().Error(err.Error(), "")
				return "", fmt.Errorf(err.Error())
			}
		}
		reading.Device = d.Name
	}

	retVal := "unsaved"
	if getConfiguration().PersistData {
		id, err := getDatabase().AddReading(reading)
		if err != nil {
			getLogger().Error(err.Error())
			return "", fmt.Errorf(err.Error())
		}
		retVal = id.Hex()
	}
	return retVal, nil
}

func CountReadings() (int, error) {
	count, err := getDatabase().ReadingCount()
	if err != nil {
		getLogger().Error(err.Error())
		return -1, err
	}
	return count, nil
}

func DeleteReadingById(id string) error {
	reading, err := getDatabase().ReadingById(id)
	if err != nil {
		var msg string
		if err == errs.ErrNotFound { //One could argue this shouldn't throw an error
			msg = fmt.Sprintf("Reading not found %s", id)
		} else {
			msg = err.Error()
		}
		getLogger().Error(msg)
		return err
	}

	err = getDatabase().DeleteReadingById(reading.Id.Hex())
	if err != nil {
		getLogger().Error(err.Error())
		return err
	}
	return nil
}

func GetReadingsByCreateTime(startTime, endTime int64, limit int) ([]models.Reading, error) {
	if limit > getConfiguration().ReadMaxLimit {
		limit = getConfiguration().ReadMaxLimit
	}

	readings, err := getDatabase().ReadingsByCreationTime(startTime, endTime, limit)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	return readings, err
}

func GetReadingsByDevice(device string, limit int) ([]models.Reading, error) {
	if limit > getConfiguration().ReadMaxLimit {
		limit = getConfiguration().ReadMaxLimit
	}

	// Try to get device
	// First check by name
	var d models.Device
	d, err := mdc.DeviceForName(device)
	if err != nil {
		// Then check by ID
		d, err = mdc.Device(device)
		if err != nil {
			if getConfiguration().MetaDataCheck {
				getLogger().Error("Error getting readings for a non-existent device: " + device)
				return nil, errs.ErrNotFound
			}
		}
	}

	readings, err := getDatabase().ReadingsByDevice(d.Name, limit)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	return readings, err
}

func GetReadingById(id string) (models.Reading, error) {
	reading, err := getDatabase().ReadingById(id)
	if err != nil {
		var msg string
		if err == errs.ErrNotFound {
			msg = fmt.Sprintf("Reading not found %s", id)
		} else {
			msg = err.Error()
		}
		getLogger().Error(msg)
		return models.Reading{}, err
	}
	return reading, nil
}

func GetReadingsByType(valType string, limit int) ([]models.Reading, error) {
	// Limit exceeds max limit
	if limit > getConfiguration().ReadMaxLimit {
		limit = getConfiguration().ReadMaxLimit
	}

	// Get the value descriptors
	vdList, err := getDatabase().ValueDescriptorsByType(valType)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	var vdNames []string
	for _, vd := range vdList {
		vdNames = append(vdNames, vd.Name)
	}

	readings, err := getDatabase().ReadingsByValueDescriptorNames(vdNames, limit)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	return readings, err
}

func GetReadingsByUomLabel(label string, limit int) ([]models.Reading, error) {
	// Limit was exceeded
	if limit > getConfiguration().ReadMaxLimit {
		limit = getConfiguration().ReadMaxLimit
	}

	// Get the value descriptors
	vList, err := GetValueDescriptorsByUomLabel(label)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}

	var vNames []string
	for _, v := range vList {
		vNames = append(vNames, v.Name)
	}

	readings, err := getDatabase().ReadingsByValueDescriptorNames(vNames, limit)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	return readings, err
}

func GetReadingsByValueDescriptor(descriptor string, limit int) ([]models.Reading, error) {
	// Check for value descriptor
	_, err := GetValueDescriptorByName(descriptor)
	if err != nil {
		var msg string
		if err == errs.ErrNotFound {
			msg = fmt.Sprintf("Value Descriptor not found %s", descriptor)
		} else {
			msg = err.Error()
		}
		getLogger().Error(msg)
		return nil, err
	}

	if limit > getConfiguration().ReadMaxLimit {
		limit = getConfiguration().ReadMaxLimit
	}

	readings, err := getDatabase().ReadingsByValueDescriptor(descriptor, limit)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	return readings, err
}

func GetReadingsByValueDescriptorLabel(label string, limit int) ([]models.Reading, error) {
	// Limit is too large
	if limit > getConfiguration().ReadMaxLimit {
		limit = getConfiguration().ReadMaxLimit
	}

	// Get the value descriptors
	vdList, err := GetValueDescriptorsByLabel(label)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	var vdNames []string
	for _, vd := range vdList {
		vdNames = append(vdNames, vd.Name)
	}

	readings, err := getDatabase().ReadingsByValueDescriptorNames(vdNames, limit)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	return readings, err
}

func UpdateReading(from models.Reading) error {
	// Check if the reading exists
	to, err := getDatabase().ReadingById(from.Id.Hex())
	if err != nil {
		var msg string
		if err == errs.ErrNotFound {
			msg = fmt.Sprintf("Reading not found %s", from.Id.Hex())
		} else {
			msg = err.Error()
		}
		getLogger().Error(msg)
		return err
	}

	//Update the fields
	if from.Value != "" {
		to.Value = from.Value
	}
	if from.Name != "" {
		_, err := getDatabase().ValueDescriptorByName(from.Name)
		if err != nil {
			var msg string
			if err == errs.ErrNotFound {
				msg = fmt.Sprintf("no value descriptor for reading %s", from.Name)
			} else {
				msg = err.Error()
			}
			getLogger().Error(msg)
			return err
		}
		to.Name = from.Name
	}
	if from.Origin != 0 {
		to.Origin = from.Origin
	}

	err = getDatabase().UpdateReading(to)
	if err != nil {
		getLogger().Error(err.Error())
	}
	return err
}
