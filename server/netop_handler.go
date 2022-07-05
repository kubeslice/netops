/*  Copyright (c) 2022 Avesha, Inc. All rights reserved.
 *
 *  SPDX-License-Identifier: Apache-2.0
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package server

import (
	"errors"
	"fmt"
	"os"

	netops "github.com/kubeslice/netops/pkg/proto"

	"github.com/kubeslice/netops/logger"
)

var (
	// NetOpHandle is a global handle to store the slice information in a map.
	NetOpHandle map[string]*SliceInfo
	// Map of tc class ID to slice name
	tcClassIdMap map[uint32]string
	netIface     string
	// Handle for htb root qdisc. Try to keep the handle ID obscure to avoid
	// interfering with the exisiting config on the intf.
	// TODO:
	// We should choose the handle by reading the existing tc config on the
	// intf. This can be done when we use a programmable golang package to
	// manage tc.
	htbRootHandleId uint32 = 17
	// This is used to derive parent class for the slice.
	// Every slice will have a parent class that is attached to the root htb.
	// We are choosing a multiple of 11 so that we can have 10 classes under the
	// parent slice class.
	// Slice 1 would have its parent class ID as 11. While slice 2 would have its
	// parent class ID as 22 and so on in multiples of 11.
	// Child classes under slice 1 would range from :12 to :21, while child classes
	// under slice 2 from :23 to :32.
	tcParentClassIdMultiple uint32 = 11
)

const MAX_NUM_OF_SLICE uint32 = 100

// BootstrapNetOpPod handles the bootstrap of the NetOp Pod.
func BootstrapNetOpPod() error {

	NetOpHandle = make(map[string]*SliceInfo)
	tcClassIdMap = make(map[uint32]string)
	netIface = "eth0"
	if os.Getenv("NETWORK_INTERFACE") != "" {
		netIface = os.Getenv("NETWORK_INTERFACE")
	}

	// Start with a clean slate:  delete TC root qdisc
	err := netOpDelTcRootQdisc()
	if err != nil {
		return err
	}

	logger.GlobalLogger.Infof("NetOp Pod is Bootstraped Successfully")
	return nil
}

func TcCmdError(tcCmd string, err error, cmdOut string) string {
	errStr := fmt.Sprintf("tc Command: %v execution failed with err: %v and stderr : %v", tcCmd, err, cmdOut)
	return errStr
}

func tcCmdOut(tcCmd string, cmdOut string) string {
	output := fmt.Sprintf("tc Command: %v output :%v", tcCmd, cmdOut)
	return output
}

func tcCmdShowNetInf(netIface string) string {
	output := fmt.Sprintf("tc qdisc show dev %s", netIface)
	return output
}

func sliceIdNotFound(sliceID string) string {
	output := fmt.Sprintf("SliceId %v is not found", sliceID)
	return output
}

func netOpAddTcRootQdisc() error {
	tcCmd := fmt.Sprintf("tc qdisc add dev %s root handle %d: htb default 30", netIface, htbRootHandleId)
	cmdOut, err := runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return errors.New(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))
	tcCmd = tcCmdShowNetInf(netIface)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))

	return nil
}

func netOpDelTcRootQdisc() error {
	tcCmd := fmt.Sprintf("tc qdisc delete dev %s root", netIface)
	cmdOut, err := runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return errors.New(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))
	tcCmd = tcCmdShowNetInf(tcCmd)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))

	return nil
}

func (s *NetOps) configureTcForSliceGwPort(gwType sliceGwType, localPort string, remotePort string, prio uint32, flowId string) error {
	//  This command adds a filter to the qdisc 1: of dev eth0, set the
	//  priority of the filter to 1, matches packets with a
	//  destination port 32100, and make the class 1:10 process the
	//  packets that match.
	tcCmd := ""
	if gwType == SLICE_GW_CLIENT {
		tcCmd = fmt.Sprintf("tc filter add dev %s protocol ip parent %d: prio %d u32 match ip dport %s 0xffff flowid %s",
			netIface, htbRootHandleId, prio, remotePort, flowId)
	} else if gwType == SLICE_GW_SERVER {
		tcCmd = fmt.Sprintf("tc filter add dev %s protocol ip parent %d: prio %d u32 match ip sport %s 0xffff flowid %s",
			netIface, htbRootHandleId, prio, localPort, flowId)
	} else {
		return nil
	}
	cmdOut, err := runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return errors.New(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))

	tcCmd = fmt.Sprintf("tc filter show dev %s", netIface)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))

	return nil
}

func (s *NetOps) configureTcForSliceGw(sliceID string, newTc *TcInfo) error {
	sliceInfo, found := NetOpHandle[sliceID]
	if !found {
		errVal := sliceIdNotFound(sliceID)
		return errors.New(errVal)
	}
	for k := range sliceInfo.sliceGwInfo {
		if !sliceInfo.sliceGwInfo[k].tcConfigured {
			err := s.configureTcForSliceGwPort(
				sliceInfo.sliceGwInfo[k].gwType,
				sliceInfo.sliceGwInfo[k].localPort,
				sliceInfo.sliceGwInfo[k].remotePort,
				sliceInfo.tc.priority, sliceInfo.tcLeafClassFqId)
			if err == nil {
				sliceInfo.sliceGwInfo[k].tcConfigured = true
			}
		}
	}

	return nil
}

func (s *NetOps) deleteTcForSliceGwAll() error {
	tcCmd := fmt.Sprintf("tc filter delete dev %s parent %d:", netIface, htbRootHandleId)
	cmdOut, err := runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return err
	}

	logger.GlobalLogger.Infof("Deleting tc config for all slice GWsi across all slices\n")

	return nil
}

func (s *NetOps) invalidateSliceGwTcConfig(sliceID string) {
	_, found := NetOpHandle[sliceID]
	if !found {
		return
	}
	for k := range NetOpHandle[sliceID].sliceGwInfo {
		NetOpHandle[sliceID].sliceGwInfo[k].tcConfigured = false
	}
	logger.GlobalLogger.Infof("Invalidated tc config for slice GWs. slice id: %s\n", sliceID)
}

func (s *NetOps) configureParentTcForSlice(sliceID string, newTc *TcInfo) error {
	// Create a tc class object for the slice under the root qdisc. We will have a parent
	// class under root qdisc for each slice.
	// tc class add dev eth0 parent 17: classid 17:1 htb rate 5mbit burst 64k
	classIdStr := fmt.Sprintf("%d:%d", htbRootHandleId, NetOpHandle[sliceID].tcParentClassId)
	tcCmd := fmt.Sprintf("tc class add dev %s parent %d: classid %s htb rate %dkbit burst 64k",
		netIface, htbRootHandleId, classIdStr, newTc.bwCeiling)
	cmdOut, err := runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return errors.New(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))
	tcCmd = tcCmdShowNetInf(netIface)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))

	NetOpHandle[sliceID].tcParentClassFqId = classIdStr
	NetOpHandle[sliceID].tcInited = true

	return nil
}

func (s *NetOps) deleteTcForSlice(sliceID string) error {
	sliceInfo, found := NetOpHandle[sliceID]
	if !found {
		logger.GlobalLogger.Infof("Delete slice tc: SliceId %v is not found", sliceID)
		return nil
	}

	err := s.deleteTcForSliceGwAll()
	if err != nil {
		logger.GlobalLogger.Errorf("Failed to delete TC settings for sliceGWs: %v, err: %v", sliceID, err)
		return err
	}

	// Delete the leaf class for the slice
	tcCmd := fmt.Sprintf("tc class delete dev %s parent %s classid %s",
		netIface, sliceInfo.tcParentClassFqId, sliceInfo.tcLeafClassFqId)
	cmdOut, err := runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		// Do not return error yet. Lets try deleting the parent class which would in turn cleanup the child classes.
	}

	// Delete the parent class for the slice
	tcCmd = fmt.Sprintf("tc class delete dev %s parent %d: classid %s",
		netIface, htbRootHandleId, sliceInfo.tcParentClassFqId)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return err
	}

	return nil
}

func (s *NetOps) configureTcForSlice(sliceID string, newTc *TcInfo) error {
	sliceInfo, found := NetOpHandle[sliceID]
	if !found {
		errVal := sliceIdNotFound(sliceID)
		return errors.New(errVal)
	}

	if sliceInfo.tc == nil {
		logger.GlobalLogger.Infof("Recieved nil value for tc configuration from slice controller")
		return nil
	}

	// Check if there are any changes in TC parameters.
	if (*sliceInfo.tc) == (*newTc) {
		logger.GlobalLogger.Infof("No change in Slice TC params, ignoring update")
		return nil
	}
	logger.GlobalLogger.Infof("Slice TC params updated. Old: %v, New: %v", sliceInfo.tc, newTc)
	// Modify parent class config
	tcCmd := fmt.Sprintf("tc class replace dev %s parent %d: classid %s htb rate %dkbit burst 64k",
		netIface, htbRootHandleId, sliceInfo.tcParentClassFqId, newTc.bwCeiling)
	cmdOut, err := runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return errors.New(errStr)
	}

	// Modify leaf class config
	tcCmd = fmt.Sprintf("tc class replace dev %s parent %s classid %s htb rate %dkbit ceil %dkbit burst 32k",
		netIface, sliceInfo.tcParentClassFqId, sliceInfo.tcLeafClassFqId, newTc.bwGuaranteed, newTc.bwCeiling)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return errors.New(errStr)
	}
	sliceInfo.tc = newTc

	if !sliceInfo.tcInited {
		err := s.configureParentTcForSlice(sliceID, newTc)
		if err != nil {
			return err
		}
		sliceInfo.tcInited = true
	}

	// Based on the numSlices create the child class id
	// # Class 1:10, which has a rate of 3mbit
	// %tc class add dev eth0 parent 1:1 classid 1:10 htb rate 3mbit ceil 5mbit burst 32k
	// TODO:
	// We only have one child class under the parent class right now. Hence, incrementing by 1 to form
	// the child class ID is ok for now. Needs to be modified if there is a use case in the future that
	// requires us to create multiple child classes under the parent class.
	classID := fmt.Sprintf("%d:%d", htbRootHandleId, sliceInfo.tcParentClassId+1)
	tcCmd = fmt.Sprintf("tc class add dev %s parent %s classid %s htb rate %dkbit ceil %dkbit burst 32k",
		netIface, sliceInfo.tcParentClassFqId, classID, newTc.bwGuaranteed, newTc.bwCeiling)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return errors.New(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))
	NetOpHandle[sliceID].tcLeafClassFqId = classID

	// Martin Devera, author of HTB, then recommends SFQ for beneath these classes:
	handleID := fmt.Sprintf("%d", sliceInfo.tcParentClassId)
	tcCmd = fmt.Sprintf("tc qdisc add dev %s parent %s handle %s: sfq perturb 10", netIface, classID, handleID)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
		return errors.New(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))

	tcCmd = tcCmdShowNetInf(netIface)
	cmdOut, err = runTcCommand(tcCmd)
	if err != nil {
		errStr := TcCmdError(tcCmd, err, cmdOut)
		logger.GlobalLogger.Errorf(errStr)
	}
	logger.GlobalLogger.Infof(tcCmdOut(tcCmd, cmdOut))

	sliceInfo.tc = newTc

	return nil
}

func (s *NetOps) enforceSliceTc(sliceID string, newTc *TcInfo) error {
	_, found := NetOpHandle[sliceID]
	if !found {
		errVal := sliceIdNotFound(sliceID)
		return errors.New(errVal)
	}

	err := s.configureTcForSlice(sliceID, newTc)
	if err != nil {
		return err
	}

	err = s.configureTcForSliceGw(sliceID, newTc)
	if err != nil {
		return err
	}

	return nil
}

func (s *NetOps) getClassIdForSlice(sliceName string) (uint32, error) {
	var i uint32 = 0
	for i = 1; i <= MAX_NUM_OF_SLICE; i++ {
		_, found := tcClassIdMap[i*tcParentClassIdMultiple]
		if !found {
			return i * tcParentClassIdMultiple, nil
		}
	}

	return 0, errors.New("could not find a free class ID")
}

func (s *NetOps) enforceSliceQosPolicy(sliceID string, sliceName string, qosProfile *SliceQosProfile) error {
	_, found := NetOpHandle[sliceID]
	if !found {
		if len(NetOpHandle) == 0 {
			// Add root qdisc
			// tc qdisc add dev eth0 root handle 1: htb default 30
			err := netOpAddTcRootQdisc()
			if err != nil {
				return err
			}
		}
		NetOpHandle[sliceID] = &SliceInfo{}
		NetOpHandle[sliceID].sliceName = sliceName
		NetOpHandle[sliceID].qosProfile = qosProfile
		parentClassId, err := s.getClassIdForSlice(sliceName)
		if err != nil {
			logger.GlobalLogger.Errorf("Failed to assign class ID for slice: %v, err: %v", sliceName, err)
			return err
		}
		logger.GlobalLogger.Infof("Assigning class ID: %v to slice: %v", parentClassId, sliceName)
		NetOpHandle[sliceID].tcParentClassId = parentClassId
		tcClassIdMap[parentClassId] = sliceName
		NetOpHandle[sliceID].sliceGwInfo = make(map[string]*SliceGwInfo)
	}

	sliceTc := &TcInfo{
		class:        qosProfile.class,
		bwCeiling:    qosProfile.bwCeiling,
		bwGuaranteed: qosProfile.bwGuaranteed,
		priority:     qosProfile.priority,
	}

	err := s.enforceSliceTc(sliceID, sliceTc)
	if err != nil {
		logger.GlobalLogger.Errorf("Failed to enforce TC settings for slice: %v, tc: %v, err: %v", sliceID, sliceTc, err)
	}

	return err
}

func (s *NetOps) handleSliceLifeCycleEvent(sliceName string, sliceEvent netops.EventType) error {
	logger.GlobalLogger.Infof("Received slice life cycle event %v for slice %v\n", sliceEvent, sliceName)

	if sliceEvent != netops.EventType_EV_DELETE {
		return nil
	}

	foundSliceToDel := false
	for k := range NetOpHandle {
		if NetOpHandle[k].sliceName == sliceName {
			err := s.deleteTcForSlice(k)
			if err != nil {
				logger.GlobalLogger.Errorf("Failed to delete TC settings for sliceGWs: %v, err: %v", sliceName, err)
				return err
			}
			delete(tcClassIdMap, NetOpHandle[k].tcParentClassId)
			delete(NetOpHandle, k)
			logger.GlobalLogger.Infof("Deleted tc config for slice: name: %v, id: %v\n", sliceName, k)
			foundSliceToDel = true
			break
		}
	}

	if foundSliceToDel {
		// When we delete the tc filter applied on sliceGW node port, we end up deleting all
		// filters under the root qdisc. To reapply the filter on slices that are still configured
		// on the system, we mark sliceGw tc configure as invalid so that next time when we
		// receive qosProfile from the slice controller we go ahead and create the filters for all
		// existing slices and their sliceGWs.
		for k := range NetOpHandle {
			s.invalidateSliceGwTcConfig(k)
		}

		// If there are no slices anymore, remove the root qdisc. This helps with cleanup of tc config
		// when the mesh is uninstalled from the cluster.
		if len(NetOpHandle) == 0 {
			logger.GlobalLogger.Infof("Deleting root tc config as no slices present on the node\n")
			err := netOpDelTcRootQdisc()
			if err != nil {
				logger.GlobalLogger.Errorf("Failed to delete root qdisc, err: %v\n", err)
			}
		}
	}
	return nil
}

func updateSliceGwInfo(sliceID string, gwInfo *SliceGwInfo) {
	_, found := NetOpHandle[sliceID]
	if !found {
		logger.GlobalLogger.Infof("Slice info not available yet: %v. Cannot update GW info", sliceID)
		return
	}
	_, found = NetOpHandle[sliceID].sliceGwInfo[gwInfo.sliceGwId]
	if !found {
		NetOpHandle[sliceID].sliceGwInfo[gwInfo.sliceGwId] = gwInfo
	} else {
		// Check if sliceGW info has changed.
		if NetOpHandle[sliceID].sliceGwInfo[gwInfo.sliceGwId].gwType != gwInfo.gwType ||
			NetOpHandle[sliceID].sliceGwInfo[gwInfo.sliceGwId].localPort != gwInfo.localPort ||
			NetOpHandle[sliceID].sliceGwInfo[gwInfo.sliceGwId].remotePort != gwInfo.remotePort {
			NetOpHandle[sliceID].sliceGwInfo[gwInfo.sliceGwId] = gwInfo
			NetOpHandle[sliceID].sliceGwInfo[gwInfo.sliceGwId].tcConfigured = false
		}
	}
}
