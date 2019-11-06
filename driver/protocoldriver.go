// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example implementation of
// ProtocolDriver interface.
//
package driver

import (
	"fmt"
	"os/exec"
	"time"
	"strings"

	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

type ProtocolDriver struct {
	lc           logger.LoggingClient
	asyncCh      chan<- *dsModels.AsyncValues
}

// Define application
var  df1ErrCount int32 = 0

func restartDf1d(){
	if df1ErrCount >5 {
		stopCmd := []string{"pkill", "df1d"}
		restartCmd := []string{"df1d"}

		exec.Command(stopCmd[0], stopCmd[1]).Run()
		exec.Command(restartCmd[0], restartCmd[1:]...).Run()

		df1ErrCount = 0
	}
}

func cmdFunc(cmd... string) (string, error) {
//	fmt.Printf("cmd len: %d, value:%v\n", len(cmd),  cmd)
	result, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return strings.TrimSpace(string(result)), err 
	}

	return strings.TrimSpace(string(result)), nil
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *ProtocolDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	s.lc = lc
	s.asyncCh = asyncCh
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *ProtocolDriver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {

	if len(reqs) != 1 {
		err = fmt.Errorf("ProtocolDriver.HandleReadCommands; too many command requests; only one supported")
		return
	}

	addr, ok := reqs[0].Attributes["addr"]
	if ok==false {
		err = fmt.Errorf("ProtocolDriver.HandleReadCommands; DeviceResource without addr attribution")
		return
	} 
	rwCmd := []string{"df1c", "127.0.0.1", "N7:1"}
	res = make([]*dsModels.CommandValue, 1)
	now := time.Now().UnixNano() / int64(time.Millisecond)
	rwCmd[2] = addr

	s.lc.Debug(fmt.Sprintf("ProtocolDriver.HandleReadCommands: protocols: %v, resource: %v, param: %v", protocols, reqs[0].DeviceResourceName, rwCmd))

	str, err := cmdFunc(rwCmd...)
	if err != nil {
		df1ErrCount++
		restartDf1d()
		err = fmt.Errorf("ProtocolDriver.HandleReadCommands: [%v]:%v, failed:%d", err, string(str), df1ErrCount)
		s.lc.Debug(err.Error())
		return
	}

	cvi, err := dsModels.NewCommandValue(reqs[0].DeviceResourceName, now, str, reqs[0].Type)
	res[0] = cvi

	return
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *ProtocolDriver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {

	if len(reqs) != 1 {
		err := fmt.Errorf("ProtocolDriver.HandleWriteCommands; too many command requests; only one supported")
		return err
	}
	if len(params) != 1 {
		err := fmt.Errorf("ProtocolDriver.HandleWriteCommands; the number of parameter is not correct; only one supported")
		return err
	}

	var err error
	rwCmd := []string{"df1c", "127.0.0.1", "N7:1"}
	addr, ok := reqs[0].Attributes["addr"]
	if ok == false {
		err = fmt.Errorf("ProtocolDriver.HandleReadCommands; DeviceResource without addr attribution")
		return err
	} 
	rwCmd[2] = addr + "=" + params[0].ValueToString()

	s.lc.Debug(fmt.Sprintf("ProtocolDriver.HandleReadCommands: protocols: %v, resource: %v, param: %v, write param: %v", protocols, reqs[0].DeviceResourceName, params, rwCmd))

	str, err := cmdFunc(rwCmd...)
	if err != nil {
		df1ErrCount++
		restartDf1d()
		err = fmt.Errorf("ProtocolDriver.HandleReadCommands: [%v]:%v, failed:%d", err, string(str), df1ErrCount)
		s.lc.Debug(err.Error())
		return err
	}

	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *ProtocolDriver) Stop(force bool) error {
	s.lc.Debug(fmt.Sprintf("ProtocolDriver.Stop called: force=%v", force))
	return nil
}
