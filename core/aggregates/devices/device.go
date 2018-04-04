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
	"github.com/edgexfoundry/edgex-go/core/clients/metadataclients"
	"github.com/edgexfoundry/edgex-go/core/data/config"
	"github.com/edgexfoundry/edgex-go/core/data/clients"
	"github.com/edgexfoundry/edgex-go/core/data/log"
	"github.com/edgexfoundry/edgex-go/core/data/messaging"
	"github.com/edgexfoundry/edgex-go/support/logging-client"
)

func getConfiguration() *config.ConfigurationStruct {
	return config.Configuration
}

func getDatabase() clients.DBClient {
	return clients.CurrentClient
}

func getDeviceClient() metadataclients.DeviceClient {
	return metadataclients.GetDeviceClient()
}

func getServiceClient() metadataclients.ServiceClient {
	return metadataclients.GetServiceClient()
}

func getLogger() logger.LoggingClient {
	return log.Logger
}

func getMQPublisher() messaging.EventPublisher {
	return messaging.CurrentPublisher
}

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

func updateDeviceLastReportedConnected(device string) {
	// Config set to skip update last reported
	if !getConfiguration().DeviceUpdateLastConnected {
		getLogger().Debug("Skipping update of device connected/reported times for:  " + device)
		return
	}

	t := time.Now().UnixNano() / int64(time.Millisecond)

	// Get the device by name
	d, err := getDeviceClient().DeviceForName(device)
	if err != nil {
		getLogger().Error("Error getting device " + device + ": " + err.Error())
		return
	}

	// Couldn't find by name
	if &d == nil {
		// Get the device by ID
		d, err = getDeviceClient().Device(device)
		if err != nil {
			getLogger().Error("Error getting device " + device + ": " + err.Error())
			return
		}

		// Couldn't find device
		if &d == nil {
			getLogger().Error("Error updating device connected/reported times.  Unknown device with identifier of:  " + device)
			return
		}

		// Got device by ID, now update lastReported/Connected by ID
		err = getDeviceClient().UpdateLastConnected(d.Id.Hex(), t)
		if err != nil {
			getLogger().Error("Problems updating last connected value for device: " + d.Id.Hex())
			return
		}
		err = getDeviceClient().UpdateLastReported(d.Id.Hex(), t)
		if err != nil {
			getLogger().Error("Problems updating last reported value for device: " + d.Id.Hex())
		}
		return
	}

	// Found by name, now update lastReported
	err = getDeviceClient().UpdateLastConnectedByName(d.Name, t)
	if err != nil {
		getLogger().Error("Problems updating last connected value for device: " + d.Name)
		return
	}
	err = getDeviceClient().UpdateLastReportedByName(d.Name, t)
	if err != nil {
		getLogger().Error("Problems updating last reported value for device: " + d.Name)
	}
	return
}

func updateDeviceServiceLastReportedConnected(device string) {
	if !getConfiguration().ServiceUpdateLastConnected {
		getLogger().Debug("Skipping update of device service connected/reported times for:  " + device)
		return
	}

	t := time.Now().UnixNano() / int64(time.Millisecond)

	// Get the device
	d, err := getDeviceClient().DeviceForName(device)
	if err != nil {
		getLogger().Error("Error getting device " + device + ": " + err.Error())
		return
	}

	// Couldn't find by name
	if &d == nil {
		d, err = getDeviceClient().Device(device)
		if err != nil {
			getLogger().Error("Error getting device " + device + ": " + err.Error())
			return
		}
		// Couldn't find device
		if &d == nil {
			getLogger().Error("Error updating device connected/reported times.  Unknown device with identifier of:  " + device)
			return
		}
	}

	// Get the device service
	s := d.Service
	if &s == nil {
		getLogger().Error("Error updating device service connected/reported times.  Unknown device service in device:  " + d.Id.Hex())
		return
	}

	getServiceClient().UpdateLastConnected(s.Service.Id.Hex(), t)
	getServiceClient().UpdateLastReported(s.Service.Id.Hex(), t)
}