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
	"time"

	"github.com/edgexfoundry/edgex-go/core/aggregates"
	"github.com/edgexfoundry/edgex-go/core/aggregates/events"
	"fmt"
)

func init() {
	go func() {
		evt := <- events.EventAggregateEvents
		switch evt.(type) {
		case aggregates.DeviceLastReported:
			e := evt.(aggregates.DeviceLastReported)
			updateDeviceLastReportedConnected(e.DeviceName)
			break;
		case aggregates.DeviceServiceLastReported:
			e := evt.(aggregates.DeviceServiceLastReported)
			updateDeviceServiceLastReportedConnected(e.DeviceName)
			break;
		}
	}()
}

func updateDeviceLastReportedConnected(device string) error {
	// Config set to skip update last reported
	if !getConfiguration().DeviceUpdateLastConnected {
		err := fmt.Errorf("Skipping update of device connected/reported times for:  %s", device)
		getLogger().Debug(err.Error())
		return err
	}

	t := time.Now().UnixNano() / int64(time.Millisecond)

	// Get the device by name
	d, err := getDeviceClient().DeviceForName(device)
	if err != nil {
		msg := fmt.Sprintf("error getting device by name %s: %v", device, err)
		getLogger().Error(msg)
	}

	// Couldn't find by name
	if len(d.Name) == 0 {
		// Get the device by ID
		d, err = getDeviceClient().Device(device)
		if err != nil {
			msg := fmt.Sprintf("Error getting device %s: %v", device, err)
			getLogger().Error(msg)
			return err
		}

		// Couldn't find device
		if len(d.Name) == 0 {
			msg := fmt.Sprintf("Error updating device connected/reported times. Unknown device with identifier of: %s", device)
			getLogger().Error(msg)
			return err
		}

		// Got device by ID, now update lastReported/Connected by ID
		err = getDeviceClient().UpdateLastConnected(d.Id.Hex(), t)
		if err != nil {
			msg := fmt.Sprintf("Problems updating last connected value for device: %s", d.Id.Hex())
			getLogger().Error(msg)
			return err
		}
		err = getDeviceClient().UpdateLastReported(d.Id.Hex(), t)
		if err != nil {
			msg := fmt.Sprintf("Problems updating last reported value for device: %s", d.Id.Hex())
			getLogger().Error(msg)
			return err
		}
	}

	// Found by name, now update lastReported
	err = getDeviceClient().UpdateLastConnectedByName(d.Name, t)
	if err != nil {
		msg := fmt.Sprintf("Problems updating last connected value for device: %s", d.Name)
		getLogger().Error(msg)
		return err
	}
	err = getDeviceClient().UpdateLastReportedByName(d.Name, t)
	if err != nil {
		msg := fmt.Sprintf("Problems updating last reported value for device: %s", d.Name)
		getLogger().Error(msg)
		return err
	}
	return nil
}

func updateDeviceServiceLastReportedConnected(device string) error {
	if !getConfiguration().ServiceUpdateLastConnected {
		err := fmt.Errorf("Skipping update of device service connected/reported times for:  " + device)
		getLogger().Error(err.Error())
		return err
	}

	t := time.Now().UnixNano() / int64(time.Millisecond)

	// Get the device
	d, err := getDeviceClient().DeviceForName(device)
	if err != nil {
		msg := fmt.Sprintf("Error getting device %s: %v", device, err.Error())
		getLogger().Error(msg)
	}

	// Couldn't find by name
	if len(d.Name) == 0 {
		d, err = getDeviceClient().Device(device)
		if err != nil {
			msg := fmt.Sprintf("Error getting device %s: %v", device, err.Error())
			getLogger().Error(msg)
			return err
		}
		// Couldn't find device
		if len(d.Name) == 0 {
			msg := fmt.Sprintf("Error updating device connected/reported times.  Unknown device with identifier of: %s", device)
			getLogger().Error(msg)
			return err
		}
	}

	// Get the device service
	s := d.Service
	if &s == nil {
		msg := fmt.Sprintf("Error updating device service connected/reported times.  Unknown device service in device: %s", d.Id.Hex())
		getLogger().Error(msg)
		return err
	}

	err = getServiceClient().UpdateLastConnected(s.Service.Id.Hex(), t)
	if err != nil {
		msg := fmt.Sprintf("error updating service connection %s: %v", s.Service.Id.Hex(), err.Error())
		getLogger().Error(msg)
		return err
	}
	err = getServiceClient().UpdateLastReported(s.Service.Id.Hex(), t)
	if err != nil {
		msg := fmt.Sprintf("error updating service reported %s: %v", s.Service.Id.Hex(), err.Error())
		getLogger().Error(msg)
		return err
	}
	return nil
}