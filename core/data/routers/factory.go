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

package routers

import (
	"github.com/edgexfoundry/edgex-go/routing"
	"github.com/gorilla/mux"
	"fmt"
)

type routerType int

const (
	Gorilla routerType = iota
	Opentrace
)

type RouterTyper interface {
	Type() routerType
}

func (t routerType) Type() routerType {
	return t
}

func NewRouter(t RouterTyper) (routing.RestRouter, error) {
	switch t {
	case Gorilla:
		g := &gorillaRouter{mux.NewRouter() }
		return g, nil
		break;
	case Opentrace:
		return nil, fmt.Errorf("opentrace not yet supported")
		break;
	}
	//effectively default case
	return nil, fmt.Errorf("unrecognized router type")
}