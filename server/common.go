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

// sliceGwType - Type of the Slice Gateway
type sliceGwType string

// sliceGwType - Slice Gateway Host type  -  client/server
const (
	SLICE_GW_SERVER sliceGwType = "SLICE_GW_SERVER"
	SLICE_GW_CLIENT sliceGwType = "SLICE_GW_CLIENT"
)

// classType - Type of the Class
type classType string

type SliceGwInfo struct {
	// Slice GW ID
	sliceGwId    string
	gwType       sliceGwType
	localPorts    []string
	remotePorts   []string
	tcConfigured bool
}

// SliceInfo - the Slice information
type SliceInfo struct {
	// Name of the slice
	sliceName string
	// QoS profile of the slice
	qosProfile *SliceQosProfile
	// Slice Tc parent class ID
	tcParentClassId uint32
	// Slice Tc fully qualified parent class ID.
	// It will be in the form "x:y", where x is the root qdisc ID and y is
	// the parent class ID for the slice.
	tcParentClassFqId string
	// Slice Tc fully qualified leaf class ID.
	// It will be in the form "x:y", where x is the root qdisc ID and y is
	// the leaf class ID for the slice.
	tcLeafClassFqId string
	// Flag to check if the parent class has been configured for the slice.
	tcInited bool
	// Tc configuration received from the slice controller for the slice.
	tc *TcInfo
	// SliceGw info
	sliceGwInfo map[string]*SliceGwInfo
}

// tcInfo - the TC information
type TcInfo struct {
	// ClassType
	class classType
	// Bandwidth Ceiling in Kbps
	bwCeiling uint32
	// Bandwidth Guaranteed
	bwGuaranteed uint32
	// Priority
	priority uint32
}

// sliceQosProfile structure to store slice QoS Profile
type SliceQosProfile struct {
	// ClassType
	class classType
	// Bandwidth Ceiling in Kbps
	bwCeiling uint32
	// Bandwidth Guaranteed
	bwGuaranteed uint32
	// Priority
	priority uint32
}
