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

	"github.com/kubeslice/netops/logger"
	netops "github.com/kubeslice/netops/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NetOps represents the GRPC NetOps
type NetOps struct {
	netops.UnimplementedNetOpsServiceServer
}

// UpdateSliceQosProfile implements the QoS Policy for a slice
func (s *NetOps) UpdateSliceQosProfile(ctx context.Context, qosProfile *netops.SliceQosProfile) (*netops.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "Client canceled, ignoring qos update message.")
	}
	if qosProfile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Qos profile message is empty")
	}

	logger.GlobalLogger.Debugf("SliceQosProfile : %v", qosProfile)

	err := s.enforceSliceQosPolicy(
		qosProfile.GetSliceId(),
		qosProfile.GetSliceName(),
		&SliceQosProfile{
			class:        classType(qosProfile.GetClassType().String()),
			bwCeiling:    qosProfile.GetBwCeiling(),
			bwGuaranteed: qosProfile.GetBwGuaranteed(),
			priority:     qosProfile.GetPriority(),
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to enforce QoS policy: %v", err)
	}
	logger.GlobalLogger.Debugf("Slice QoS policy enforced successfully")

	return &netops.Response{StatusMsg: "Slice QoS policy enforced successfully"}, nil
}

// UpdateSliceLifeCycleEvent handles slice life cycle events
func (s *NetOps) UpdateSliceLifeCycleEvent(ctx context.Context, sliceEvent *netops.SliceLifeCycleEvent) (*netops.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "Client canceled, ignoring qos update message.")
	}
	if sliceEvent == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Slice lifecycle message is empty")
	}

	logger.GlobalLogger.Infof("SliceLifeCycleEvent : %v", sliceEvent)

	err := s.handleSliceLifeCycleEvent(sliceEvent.GetSliceName(), sliceEvent.GetEvent())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to handle slice lifecycle event: %v", err)
	}
	logger.GlobalLogger.Infof("Slice life cycle event handled successfully")

	return &netops.Response{StatusMsg: "Slice life cycle event handled successfully"}, nil
}

// UpdateConnectionContext updates the connection context and adds the route
func (s *NetOps) UpdateConnectionContext(ctx context.Context, conContext *netops.NetOpConnectionContext) (*netops.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
	}
	if conContext == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Connection Context is Empty")
	}
	if conContext.GetLocalSliceGwNodePort() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Slice Gateway Node Port")
	}
	if conContext.GetLocalSliceGwHostType().String() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Slice Gateway Host Type")
	}
	logger.GlobalLogger.Infof("conContext : %v", conContext)

	updateSliceGwInfo(
		conContext.GetSliceId(),
		&SliceGwInfo{
			sliceGwId:  conContext.GetLocalSliceGwId(),
			gwType:     sliceGwType(conContext.GetLocalSliceGwHostType().String()),
			localPort:  conContext.GetLocalSliceGwNodePort(),
			remotePort: conContext.GetRemoteSliceGwNodePort(),
		},
	)

	logger.GlobalLogger.Infof("Connection Context Updated Successfully")

	return &netops.Response{StatusMsg: "Connection Context Updated Successfully in netops pod"}, nil
}
