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
	"errors"
	"fmt"
	"regexp"

	"github.com/edgexfoundry/edgex-go/core/domain/errs"
	"github.com/edgexfoundry/edgex-go/core/domain/models" //for now
)

const (
	formatSpecifier          = "%(\\d+\\$)?([-#+ 0,(\\<]*)?(\\d+)?(\\.\\d+)?([tT])?([a-zA-Z%])"
)

func AddValueDescriptor(v models.ValueDescriptor) (string, error) {
	id, err := getDatabase().AddValueDescriptor(v)
	if err != nil {
		getLogger().Error(err.Error())
		return "", err
	}
	return id.Hex(), err
}

func DeleteValueDescriptorById(id string) error {
	vd, err := getDatabase().ValueDescriptorById(id)
	if err != nil {
		getLogger().Error(err.Error())
		return err
	}

	// Check if the value descriptor is still in use by readings
	readings, err := getDatabase().ReadingsByValueDescriptor(vd.Name, 10)
	if err != nil {
		getLogger().Error(err.Error())
		return err
	}
	if len(readings) > 0 {
		err = errors.New("Data integrity issue.  Value Descriptor is still referenced by existing readings.")
		getLogger().Error(err.Error())
		return err
	}

	// Delete the value descriptor
	if err = getDatabase().DeleteValueDescriptorById(vd.Id.Hex()); err != nil {
		getLogger().Error(err.Error())
		return err
	}

	return nil
}

func GetAllValueDescriptors() ([]models.ValueDescriptor, error) {
	vList, err := getDatabase().ValueDescriptors()
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}

	return vList, err
}

func GetValueDescriptorsByDeviceId(deviceId string) ([]models.ValueDescriptor, error) {
	// Get the device
	d, err := mdc.Device(deviceId)
	if err != nil {
		getLogger().Error("Device not found: " + err.Error())
		return nil, errs.ErrNotFound
	}

	// Get the names of the value descriptors
	vdNames := []string{}
	d.AllAssociatedValueDescriptors(&vdNames)

	// Get the value descriptors
	vdList := []models.ValueDescriptor{}
	for _, name := range vdNames {
		vd, err := getDatabase().ValueDescriptorByName(name)

		// Not an error if not found
		if err == errs.ErrNotFound {
			continue
		}

		if err != nil {
			getLogger().Error(err.Error())
			return vdList, err
		}

		vdList = append(vdList, vd)
	}

	return vdList, nil
}

func GetValueDescriptorsByDeviceName(device string) ([]models.ValueDescriptor, error) {
	// Get the device
	d, err := mdc.DeviceForName(device)
	if err != nil {
		getLogger().Error("Device not found: " + err.Error())
		return nil, errs.ErrNotFound
	}

	// Get the names of the value descriptors
	vdNames := []string{}
	d.AllAssociatedValueDescriptors(&vdNames)

	// Get the value descriptors
	vdList := []models.ValueDescriptor{}
	for _, name := range vdNames {
		vd, err := getDatabase().ValueDescriptorByName(name)

		// Not an error if not found
		if err == errs.ErrNotFound {
			continue
		}

		if err != nil {
			getLogger().Error(err.Error())
			return vdList, err
		}

		vdList = append(vdList, vd)
	}

	return vdList, nil
}

func GetValueDescriptorById(id string) (models.ValueDescriptor, error) {
	v, err := getDatabase().ValueDescriptorById(id)
	if err != nil {
		getLogger().Error(err.Error())
		return models.ValueDescriptor{}, err
	}
	return v, err
}

func GetValueDescriptorByName(name string) (models.ValueDescriptor, error) {
	v, err := getDatabase().ValueDescriptorByName(name)
	if err != nil {
		getLogger().Error(err.Error())
		return models.ValueDescriptor{}, err
	}
	return v, err
}

func GetValueDescriptorsByLabel(label string) ([]models.ValueDescriptor, error) {
	vdList, err := getDatabase().ValueDescriptorsByLabel(label)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	return vdList, err
}

func GetValueDescriptorsByUomLabel(label string) ([]models.ValueDescriptor, error) {
	vList, err := getDatabase().ValueDescriptorsByUomLabel(label)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, err
	}
	return vList, err
}

func UpdateValueDescriptor(from models.ValueDescriptor) error {
	// Find the value descriptor thats being updated
	// Try by ID
	to, err := getDatabase().ValueDescriptorById(from.Id.Hex())
	if err != nil {
		to, err = getDatabase().ValueDescriptorByName(from.Name)
		if err != nil {
			var msg string
			if err == errs.ErrNotFound {
				msg = fmt.Sprintf("value descriptor not found %s %s", from.Id, from.Name)
			} else {
				msg = err.Error()
			}
			getLogger().Error(msg)
			return err
		}
	}

	// Update the fields
	if from.DefaultValue != "" {
		to.DefaultValue = from.DefaultValue
	}
	if from.Formatting != "" {
		match, err := regexp.MatchString(formatSpecifier, from.Formatting)
		if err != nil {
			getLogger().Error(fmt.Sprintf("invalid format for updated value descriptor %s", from.Formatting))
			return err
		}
		if !match {
			return fmt.Errorf("Value descriptor's format string doesn't match the required pattern %s ", formatSpecifier)
		}
		to.Formatting = from.Formatting
	}
	if from.Labels != nil {
		to.Labels = from.Labels
	}

	if from.Max != "" {
		to.Max = from.Max
	}
	if from.Min != "" {
		to.Min = from.Min
	}
	if from.Name != "" {
		// Check if value descriptor is still in use by readings if the name changes
		if from.Name != to.Name {
			r, err := getDatabase().ReadingsByValueDescriptor(to.Name, 10) // Arbitrary limit, we're just checking if there are any readings
			if err != nil {
				getLogger().Error("Error checking the readings for the value descriptor: " + err.Error())
				return err
			}
			// Value descriptor is still in use
			if len(r) != 0 {
				msg := fmt.Sprintf("Data integrity issue. Value Descriptor %s is referenced by existing readings.", from.Name)
				getLogger().Error(msg)
				return fmt.Errorf(msg)
			}
		}
		to.Name = from.Name
	}
	if from.Origin != 0 {
		to.Origin = from.Origin
	}
	if from.Type != "" {
		to.Type = from.Type
	}
	if from.UomLabel != "" {
		to.UomLabel = from.UomLabel
	}

	// Push the updated valuedescriptor to the database
	err = getDatabase().UpdateValueDescriptor(to)
	if err != nil {
		getLogger().Error(err.Error())
	}
	return err
}

func ValidateFormatString(v models.ValueDescriptor) (bool, error) {
	// No formatting specified
	if v.Formatting == "" {
		return true, nil
	} else {
		return regexp.MatchString(formatSpecifier, v.Formatting)
	}
}
