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
	 "time"

	 "github.com/edgexfoundry/edgex-go/core/aggregates"
	 "github.com/edgexfoundry/edgex-go/core/domain/errs"
	 "github.com/edgexfoundry/edgex-go/core/domain/models" //for now
	 "gopkg.in/mgo.v2/bson"
 )

func CountEvents() (int, error) {
	count, err := getDatabase().EventCount()
	if err != nil {
		getLogger().Error(err.Error())
		return -1, fmt.Errorf(err.Error())
	}
	return count, nil
}

//TODO: Again with the double purposing of device as either ID or name. Need to address this.
func CountByDevice(device string) (int, error) {
	// Get the device
	// Try by ID
	d, err := getDeviceClient().Device(device)
	if err != nil {
		// Try by Name
		d, err = getDeviceClient().DeviceForName(device)
		if err != nil {
			getLogger().Error("error finding device " + device + ": " + err.Error(), "")
			return -1, fmt.Errorf("error finding device %s: %v", device, err)
		}
	}

	count, err := getDatabase().EventCountByDeviceId(d.Name)
	if err != nil {
		getLogger().Error(err.Error())
		return -1, fmt.Errorf("error obtaining count for device %s: %v", device, err)
	}
	return count, err
}

func DeleteByAge(age int64) (int, error) {
	events, err := getDatabase().EventsOlderThanAge(age)
	if err != nil {
		getLogger().Error(err.Error())
		return -1, fmt.Errorf(err.Error())
	}

	// Delete all the events
	count := len(events)
	for _, event := range events {
		if err = deleteEvent(event); err != nil {
			getLogger().Error(err.Error())
			return -1, fmt.Errorf(err.Error())
		}
	}
	return count, nil
}

func DeleteByDevice(deviceId string) (int, error) {
	// Get the device
	deviceFound := true
	d, err := getDeviceClient().Device(deviceId)
	if err != nil {
		d, err = getDeviceClient().DeviceForName(deviceId)
		if err != nil {
			deviceFound = false
		}
	}

	if deviceFound {
		deviceId = d.Name
	}

	// See if you need to check metadata
	if getConfiguration().MetaDataCheck && !deviceFound {
		getLogger().Error("Device not found for event: "+err.Error(), "")
		return -1, fmt.Errorf("device not found %s: %v ", deviceId, err)
	}

	// Get the events by the device name
	events, err := getDatabase().EventsForDevice(deviceId)
	if err != nil {
		getLogger().Error(err.Error())
		return -1, fmt.Errorf(err.Error())
	}

	getLogger().Info("Deleting the events for device: " + deviceId)

	// Delete the events
	count := len(events)
	for _, event := range events {
		if err = deleteEvent(event); err != nil {
			getLogger().Error(err.Error())
			return -1, fmt.Errorf(err.Error())
		}
	}
	return count, nil
}

func DeleteEventById(id string) error {
	e, err := getDatabase().EventById(id)
	if err != nil {
		getLogger().Error(err.Error())
		return fmt.Errorf(err.Error())
	}

	getLogger().Info("Deleting event: " + e.ID.Hex())

	err = deleteEvent(e)
	if err != nil {
		getLogger().Error(err.Error())
		return fmt.Errorf(err.Error())
	}
	return nil
}

func GetAllEvents() ([]models.Event, error) {
	events, err := getDatabase().Events()

	if err != nil {
		getLogger().Error(err.Error())
		return nil, fmt.Errorf(err.Error())
	}
	return events, err
}

func GetEventsByCreateTime(startTime, endTime int64, limit int) ([]models.Event, error) {
	if limit > getConfiguration().ReadMaxLimit {
		limit = getConfiguration().ReadMaxLimit
	}

	e, err := getDatabase().EventsByCreationTime(startTime, endTime, limit)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, fmt.Errorf(err.Error())
	}
	return e, nil
}

func GetEventsByDevice(deviceId string, limitNum int) ([]models.Event, error) {
	//TODO: Resolve the double purposing of deviceId, as elsewhere
	// Get the device
	deviceFound := true
	// Try by ID
	d, err := getDeviceClient().Device(deviceId)
	if err != nil {
		// Try by Name
		d, err = getDeviceClient().DeviceForName(deviceId)
		if err != nil {
			deviceFound = false
		}
	}

	if deviceFound {
		deviceId = d.Name
	}

	// See if you need to check metadata for the device
	if getConfiguration().MetaDataCheck && !deviceFound {
		msg := "error getting readings for non-existent device " + deviceId
		getLogger().Error(msg)
		return nil, fmt.Errorf(msg)
	}

	if limitNum > getConfiguration().ReadMaxLimit {
		limitNum = getConfiguration().ReadMaxLimit
	}

	eventList, err := getDatabase().EventsForDeviceLimit(deviceId, limitNum)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, fmt.Errorf(err.Error())
	}
	return eventList, nil
}

func GetEventById(id string) (models.Event, error) {
	evt, err := getDatabase().EventById(id)
	if err != nil {
		getLogger().Error(err.Error())
		return models.Event{}, fmt.Errorf(err.Error())
	}
	return evt, nil
}

func GetReadingsByDeviceAndValueDescriptor(deviceId string, descriptor string, limitNum int) ([]models.Reading, error) {
	if limitNum > getConfiguration().ReadMaxLimit {
		limitNum = getConfiguration().ReadMaxLimit
	}

	// Get the device
	deviceFound := true
	// Try by id
	d, err := getDeviceClient().Device(deviceId)
	if err != nil {
		// Try by name
		d, err = getDeviceClient().DeviceForName(deviceId)
		if err != nil {
			deviceFound = false
		}
	}

	if deviceFound {
		deviceId = d.Name
	}

	// See if you need to check metadata
	if getConfiguration().MetaDataCheck && !deviceFound {
		getLogger().Error("device " + deviceId + " not found for event: "+err.Error(), "")
		return nil, fmt.Errorf("device %s not found: %v", deviceId, err)
	}

	// Get all the events for the device
	e, err := getDatabase().EventsForDevice(deviceId)
	if err != nil {
		getLogger().Error(err.Error())
		return nil, fmt.Errorf(err.Error())
	}

	// Only pick the readings who match the value descriptor
	readings := []models.Reading{}
	count := 0 // Make sure we stay below the limit
	for _, event := range e {
		if count >= limitNum {
			break
		}
		for _, reading := range event.Readings {
			if count >= limitNum {
				break
			}
			if reading.Name == descriptor {
				readings = append(readings, reading)
				count += 1
			}
		}
	}
	return readings, nil
}

func AddNewEvent(evt models.Event) (string, error) {
	// Get device from metadata
	deviceFound := true
	// Try by ID
	d, err := getDeviceClient().Device(evt.Device) //TODO: Why is this property double-purposed?
	if err != nil {
		// Try by name
		d, err = getDeviceClient().DeviceForName(evt.Device)
		if err != nil {
			deviceFound = false
		}
	}
	// Make sure the identifier is the device name
	if deviceFound {
		evt.Device = d.Name
	}

	// See if metadata checking is enabled
	if getConfiguration().MetaDataCheck && !deviceFound {
		return "", fmt.Errorf("device %s not found for new event", evt.Device)
	}

	if getConfiguration().ValidateCheck {
		getLogger().Debug("Validation enabled, parsing events")
		for reading := range evt.Readings {
			valid, err := isValidValueDescriptor(evt.Readings[reading], evt)
			if !valid {
				return "", fmt.Errorf("validation failed: %s", err.Error())
			}
		}
	}

	// Add the readings to the database
	retVal := "unsaved"
	if getConfiguration().PersistData {
		for i, reading := range evt.Readings {
			// Check value descriptor
			_, err := getDatabase().ValueDescriptorByName(reading.Name)
			if err != nil {
				getLogger().Error(err.Error())
				if err == errs.ErrNotFound {
					return "", fmt.Errorf("value descriptor not found for a reading (%s): %v", reading.Name, err)
				} else {
					return "", err
				}
			}

			reading.Device = evt.Device // Update the device for the reading

			// Add the reading
			id, err := getDatabase().AddReading(reading)
			if err != nil {
				getLogger().Error(err.Error())
				return "", fmt.Errorf(err.Error())
			}

			evt.Readings[i].Id = id // Set the ID for referencing later
		}

		// Add the event to the database
		id, err := getDatabase().AddEvent(&evt)
		if err != nil {
			getLogger().Error(err.Error())
			return "", fmt.Errorf(err.Error())
		}
		retVal = id.Hex()
	}

	publishExternalEvent(evt)                                 // Push the aux struct to export service (It has the actual readings)
	EventAggregateEvents <- aggregates.DeviceLastReported{DeviceName:evt.Device} // update last reported connected (device)
	EventAggregateEvents <- aggregates.DeviceServiceLastReported{DeviceName:evt.Device} // update last reported connected (device service)

	return retVal, nil
}

func Purge() error {
	err := getDatabase().ScrubAllEvents()
	if err != nil {
		getLogger().Error("error purging all events/readings: " + err.Error())
		return fmt.Errorf(err.Error())
	}
	return nil
}

func PurgeIfPublished() (int, error) {
	//TODO: Would be more performant to do this in one shot, not as get followed by loop. There could be thousands.
	// Get the events
	events, err := getDatabase().EventsPushed()
	if err != nil {
		getLogger().Error(err.Error())
		return -1, fmt.Errorf(err.Error())
	}

	// Delete all the events
	count := len(events)
	for _, event := range events {
		if err = deleteEvent(event); err != nil {
			getLogger().Error(err.Error())
			return -1, fmt.Errorf(err.Error())
		}
	}
	return count, nil
}

func Touch(id string) error {
	if !bson.IsObjectIdHex(id) {
		msg := fmt.Sprintf("%s is not a valid bson objectId", id)
		getLogger().Error(msg)
		return fmt.Errorf(msg)
	}
	// Check if the event exists
	evt, err := getDatabase().EventById(id)
	if err != nil {
		getLogger().Error(err.Error())
		return fmt.Errorf(err.Error())
	}

	getLogger().Info("Updating event: " + evt.ID.Hex())

	evt.Pushed = time.Now().UnixNano() / int64(time.Millisecond)
	err = getDatabase().UpdateEvent(evt)
	if err != nil {
		getLogger().Error(err.Error())
		return fmt.Errorf(err.Error())
	}
	return nil
}

func UpdateEvent(from models.Event) error {
	// Check if the event exists
	to, err := getDatabase().EventById(from.ID.Hex())
	if err != nil {
		getLogger().Error(err.Error())
		return fmt.Errorf(err.Error())
	}

	getLogger().Info("Updating event: " + from.ID.Hex())

	// Update the fields
	if from.Device != "" {
		deviceFound := true
		d, err := getDeviceClient().Device(from.Device)
		if err != nil {
			d, err = getDeviceClient().DeviceForName(from.Device)
			if err != nil {
				deviceFound = false
			}
		}

		// See if we need to check metadata
		if getConfiguration().MetaDataCheck && !deviceFound {
			getLogger().Error("Error updating device, device " + from.Device + " doesn't exist")
			return fmt.Errorf("device not found: %v", err)
		}

		if deviceFound {
			to.Device = d.Name
		} else {
			to.Device = from.Device
		}
	}
	if from.Pushed != 0 {
		to.Pushed = from.Pushed
	}
	if from.Origin != 0 {
		to.Origin = from.Origin
	}

	// Update
	if err = getDatabase().UpdateEvent(to); err != nil {
		getLogger().Error(err.Error())
		return fmt.Errorf(err.Error())
	}
	return nil
}

// Delete the event and readings
func deleteEvent(e models.Event) error {
	for _, reading := range e.Readings {
		if err := getDatabase().DeleteReadingById(reading.Id.Hex()); err != nil {
			return err
		}
	}
	if err := getDatabase().DeleteEventById(e.ID.Hex()); err != nil {
		return err
	}

	return nil
}

// Put event on the message queue to be processed by the rules engine
func publishExternalEvent(e models.Event) {
	getLogger().Info("Putting event on message queue", "")
	//	Have multiple implementations (start with ZeroMQ)
	err := getMQPublisher().SendEventMessage(e)
	if err != nil {
		getLogger().Error("Unable to send message for event: " + e.String())
	}
}