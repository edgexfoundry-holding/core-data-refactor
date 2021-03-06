/*******************************************************************************
 * Copyright 2017 Dell Inc.
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
 * @author: Ryan Comer, Dell
 * @version: 0.5.0
 *******************************************************************************/
package messaging

import (
	"encoding/json"
	"sync"

	"github.com/edgexfoundry/edgex-go/core/domain/models"
	zmq "github.com/pebbe/zmq4"
)

type MQClient interface {
	SendEventMessage(e models.Event) error
}

// Configuration struct for ZeroMQ
type zeroMQConfiguration struct {
	AddressPort string
}

// ZeroMQ implementation of the event publisher
type zeroMQClient struct {
	socket	*zmq.Socket
	mux     sync.Mutex
}

func newZeroMQEventPublisher(config zeroMQConfiguration) MQClient {
	newSocket, _ := zmq.NewSocket(zmq.PUB)
	newSocket.Bind(config.AddressPort)

	return &zeroMQClient{
		socket: newSocket,
	}
}

func (zep *zeroMQClient) SendEventMessage(e models.Event) error {
	s, err := json.Marshal(&e)
	if err != nil {
		return err
	}
	zep.mux.Lock()
	defer zep.mux.Unlock()
	_, err = zep.socket.SendBytes(s, 0)
	if err != nil {
		return err
	}

	return nil
}
