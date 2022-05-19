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
	"context"
	"log"
	"testing"

	"github.com/kubeslice/netops/logger"
	netops "github.com/kubeslice/netops/pkg/proto"
	"google.golang.org/grpc"
)

func TestConfigureTcForSliceGwPort(t *testing.T) {
	testCases := []struct {
		Case       string
		GwType     sliceGwType
		localPort  string
		RemotePort string
		Priority   uint32
		FlowId     string
		ErrStr     string
	}{
		{
			"Testing for Gateway type SLICE_GW_CLIENT",
			sliceGwType("SLICE_GW_CLIENT"),
			"5000",
			"5000",
			2,
			"randomid",
			"",
		},
		{
			"Testing for Gateway type SLICE_GW_SERVER",
			sliceGwType("SLICE_GW_SERVER"),
			"5000",
			"5000",
			2,
			"randomid",
			"",
		},
	}
	ctx := context.Background()
	logger.GlobalLogger = logger.NewLogger("ERROR")
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := &NetOps{}
	for _, tt := range testCases {
		err := client.configureTcForSliceGwPort(tt.GwType, tt.localPort, tt.RemotePort, tt.Priority, tt.FlowId)
		if err != nil {
			t.Log(err.Error())
			if err.Error() != tt.ErrStr {
				t.Error("Expected :", tt.ErrStr, " but got ", err)
			}
		}
	}
}

func TestConfigureTcForSliceGw(t *testing.T) {
	testCases := []struct {
		Case     string
		EmptyMap bool
		SliceID  string
		Tc       *TcInfo
		ErrStr   string
	}{
		{
			"Testing while the NetOpHandle map is empty",
			true,
			"randomid",
			&TcInfo{class: classType(netops.ClassType_HTB.String()), bwCeiling: 1, bwGuaranteed: htbRootHandleId, priority: 2},
			"SliceId randomid is not found",
		},
		{
			"Testing with NetOpHandle populated with a value",
			false,
			"randomid",
			&TcInfo{class: classType(netops.ClassType_HTB.String()), bwCeiling: 1, bwGuaranteed: htbRootHandleId, priority: 2},
			"",
		},
	}
	ctx := context.Background()
	logger.GlobalLogger = logger.NewLogger("ERROR")
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := &NetOps{}
	for _, tt := range testCases {
		if !tt.EmptyMap {
			err := MockBootstrapNetOpPod()
			if err != nil {
				t.Error(err)
			}
		}
		err := client.configureTcForSliceGw(tt.SliceID, tt.Tc)
		if err != nil {
			t.Log(err.Error())
			if err.Error() != tt.ErrStr {
				t.Error("Expected :", tt.ErrStr, " but got ", err)
			}
		}
	}
}
func TestHandleSliceLifeCycleEvent(t *testing.T) {
	testCases := []struct {
		Case       string
		EmptyMap   bool
		SliceName  string
		SliceEvent netops.EventType
		ErrStr     string
	}{
		{
			"Providing the wrong slice Event",
			true,
			"test-slice",
			netops.EventType_EV_CREATE,
			"",
		},
		{
			"Testing for delete slice event",
			false,
			"test-slice",
			netops.EventType_EV_DELETE,
			"",
		},
	}
	ctx := context.Background()
	logger.GlobalLogger = logger.NewLogger("ERROR")
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := &NetOps{}
	for _, tt := range testCases {
		if !tt.EmptyMap {
			err := MockBootstrapNetOpPod()
			if err != nil {
				t.Error(err)
			}
		}
		err := client.handleSliceLifeCycleEvent(tt.SliceName, tt.SliceEvent)
		if err != nil {
			t.Log(err.Error())
			if err.Error() != tt.ErrStr {
				t.Error("Expected :", tt.ErrStr, " but got ", err)
			}
		}
	}
}
func TestConfigureTcForSlice(t *testing.T) {
	testCases := []struct {
		Case     string
		EmptyMap bool
		SliceID  string
		Tc       *TcInfo
		ErrStr   string
	}{
		{
			"Testing with empty NetopHandle map",
			true,
			"randomid",
			&TcInfo{class: classType(netops.ClassType_HTB.String()), bwCeiling: 1, bwGuaranteed: htbRootHandleId, priority: 2},
			"SliceId randomid is not found",
		},
		{
			"Test for ignoring the update",
			false,
			"randomid",
			&TcInfo{class: classType(netops.ClassType_HTB.String()), bwCeiling: 1, bwGuaranteed: htbRootHandleId, priority: 2},
			"No change in Slice TC params, ignoring update",
		},
		{
			"Test for updating slice tc parameters",
			false,
			"randomid",
			&TcInfo{class: classType(netops.ClassType_HTB.String()), bwCeiling: 3, bwGuaranteed: htbRootHandleId, priority: 1},
			"",
		},
		{
			"Test for updating paremeters tcLeafClassFqId and tcInited",
			false,
			"randomid2",
			&TcInfo{class: classType(netops.ClassType_HTB.String()), bwCeiling: 3, bwGuaranteed: htbRootHandleId, priority: 1},
			"",
		},
	}
	ctx := context.Background()
	logger.GlobalLogger = logger.NewLogger("INFO")
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := &NetOps{}
	for _, tt := range testCases {
		if !tt.EmptyMap {
			err := MockBootstrapNetOpPod()
			if err != nil {
				t.Error(err)
			}
		}
		err := client.configureTcForSlice(tt.SliceID, tt.Tc)
		if err != nil {
			t.Log(err.Error())
			if err.Error() != tt.ErrStr {
				t.Error("Expected :", tt.ErrStr, " but got ", err)
			}
		}
	}
}

func TestEnforceSliceQosPolicy(t *testing.T) {
	testCases := []struct {
		Case                string
		SliceID             string
		SliceName           string
		QosProfile          *SliceQosProfile
		EmptyNetOpHandleMap bool
		ErrStr              string
	}{
		{
			"Test when the NetopHandle Map is Empty",
			"randomid",
			"test-slice",
			&SliceQosProfile{class: classType(netops.ClassType_HTB.String()), bwCeiling: 3, bwGuaranteed: htbRootHandleId, priority: 1},
			true,
			"",
		},
		{
			"Test when the NetopHandle Map is Empty",
			"randomid",
			"test-slice",
			&SliceQosProfile{class: classType(netops.ClassType_HTB.String()), bwCeiling: 3, bwGuaranteed: htbRootHandleId, priority: 1},
			false,
			"",
		},
	}
	ctx := context.Background()
	logger.GlobalLogger = logger.NewLogger("ERROR")
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := &NetOps{}
	for _, tt := range testCases {
		if !tt.EmptyNetOpHandleMap {
			err := MockBootstrapNetOpPod()
			if err != nil {
				t.Error(err)
			}
		}
		err := client.enforceSliceQosPolicy(tt.SliceID, tt.SliceName, tt.QosProfile)
		if err != nil {
			t.Log(err.Error())
			if err.Error() != tt.ErrStr {
				t.Error("Expected :", tt.ErrStr, " but got ", err)
			}
		}
	}
}
func TestDeleteTcForSliceGwAll(t *testing.T) {
	testCases := []struct {
		Case               string
		CallBootstrapNetOp bool
		ErrStr             string
	}{
		{
			"Test without calling BoostrapNetOpPod",
			true,
			"No command defined :",
		},
		{
			"Testing the function with calling the helper BoostrapNetOpPod",
			false,
			"",
		},
	}
	ctx := context.Background()
	logger.GlobalLogger = logger.NewLogger("ERROR")
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := &NetOps{}
	for _, tt := range testCases {
		if !tt.CallBootstrapNetOp {
			err := MockBootstrapNetOpPod()
			if err != nil {
				t.Error(err)
			}
		}
		err := client.deleteTcForSliceGwAll()
		if err != nil {
			t.Log(err.Error())
			if err.Error() != tt.ErrStr {
				t.Error("Expected :", tt.ErrStr, " but got ", err)
			}
		}
	}
}
