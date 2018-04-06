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
 * @microservice: core-command-go service
 * @author: Trevor Conn, Dell
 * @version: 0.5.0
 *******************************************************************************/
package config

// ConfigurationStruct : Struct used to pase the JSON configuration file
type ConfigurationStruct struct {
	ApplicationName           string
	ConsulProfilesActive      string
	ReadMaxLimit              int
	ServicePort               int
	HeartBeatTime             int
	ConsulPort                int
	ServiceTimeout            int
	CheckInterval             string
	ServiceAddress            string
	ServiceName               string
	DeviceServiceProtocol     string
	HeartBeatMsg              string
	AppOpenMsg                string
	URLProtocol               string
	URLDevicePath             string
	ConsulHost                string
	ConsulCheckAddress        string
	EnableRemoteLogging       bool
	LogFile                   string
	LoggingRemoteURL          string
	MetaAddressableURL        string
	MetaDeviceServiceURL      string
	MetaDeviceProfileURL      string
	MetaDeviceURL             string
	MetaDeviceReportURL       string
	MetaCommandURL            string
	MetaEventURL              string
	MetaScheduleURL           string
	MetaProvisionWatcherURL   string
}

// Configuration data for the metadata service
var Configuration *ConfigurationStruct

