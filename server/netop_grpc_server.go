package server

import (
	"bitbucket.org/realtimeai/kubeslice-netops/logger"
	"bitbucket.org/realtimeai/kubeslice-netops/pkg/proto"
	"context"

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
