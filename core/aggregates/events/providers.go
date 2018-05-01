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
	"github.com/edgexfoundry/edgex-go/core/clients/metadataclients"
	"github.com/edgexfoundry/edgex-go/core/data/clients"
	"github.com/edgexfoundry/edgex-go/core/data/config"
	"github.com/edgexfoundry/edgex-go/core/data/log"
	"github.com/edgexfoundry/edgex-go/core/data/messaging"
	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"github.com/edgexfoundry/edgex-go/support/logging-client"
)

//Only set this manually when providing a mock for unit testing
var dc deviceClient

var EventAggregateEvents chan interface{}

func init() {
	if EventAggregateEvents == nil {
		EventAggregateEvents = make(chan interface{}, 10)
	}
}

func getDeviceClient() deviceClient {
	if dc == nil { //We check here in case a mock was supplied.
		dc = metadataclients.GetDeviceClient()
	}
	return dc
}

func getConfiguration() *config.ConfigurationStruct {
	return config.Configuration
}

func getDatabase() clients.DBClient {
	return clients.CurrentClient
}

func getLogger() logger.LoggingClient {
	return log.Logger
}

func getMQPublisher() messaging.EventPublisher {
	return messaging.CurrentPublisher
}

type deviceClient interface {
	Device(id string) (models.Device, error)
	DeviceForName(name string) (models.Device, error)
}