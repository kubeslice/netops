//go:build !server
// +build !server

package server

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/kubeslice/netops/logger"
	netops "github.com/kubeslice/netops/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func MockBootstrapNetOpPod() error {
	mockSliceGwInfo := make(map[string]*SliceGwInfo)
	mockSliceGwInfo["test-slice"] = &SliceGwInfo{tcConfigured: false, localPort: "5000", remotePort: "5000", gwType: sliceGwType("SLICE_GW_SERVER")}
	NetOpHandle = make(map[string]*SliceInfo)
	NetOpHandle["randomid"] = &SliceInfo{sliceName: "test-slice", qosProfile: &SliceQosProfile{class: classType(netops.ClassType_HTB.String()), bwCeiling: 1, bwGuaranteed: htbRootHandleId, priority: 2}, sliceGwInfo: mockSliceGwInfo, tcLeafClassFqId: "randomLeafID", tc: &TcInfo{class: classType(netops.ClassType_HTB.String()), bwCeiling: 1, bwGuaranteed: htbRootHandleId, priority: 2}, tcInited: false}
	NetOpHandle["randomid2"] = &SliceInfo{sliceName: "test-slice", qosProfile: &SliceQosProfile{class: classType(netops.ClassType_HTB.String()), bwCeiling: 1, bwGuaranteed: htbRootHandleId, priority: 2}, sliceGwInfo: mockSliceGwInfo, tcLeafClassFqId: "randomLeafID"}
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

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	grpcServer := grpc.NewServer()

	netops.RegisterNetOpsServiceServer(grpcServer, &NetOps{})

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestUpdateSliceQosProfileWithMockNetOpHandle(t *testing.T) {
	tests := []struct {
		testCase string
		res      *netops.Response
		sqos     *netops.SliceQosProfile
		errMsg   string
		errCode  codes.Code
		ctxCanel bool
	}{
		{
			"Test for successfully enforcing slice QoS policy",
			&netops.Response{StatusMsg: "Slice QoS policy enforced successfully"},
			&netops.SliceQosProfile{
				SliceName:      "test-slice",
				SliceId:        "randomid",
				QosProfileName: "",
				TcType:         netops.TcType_BANDWIDTH_CONTROL,
				ClassType:      netops.ClassType_HTB,
				BwCeiling:      1,
				BwGuaranteed:   htbRootHandleId,
				Priority:       2,
				DscpClass:      "",
			},
			"",
			codes.OK,
			false,
		},
		{
			"Test to validate Empty Slice Qos Profile",
			&netops.Response{StatusMsg: ""},
			&netops.SliceQosProfile{SliceName: "random-name"},
			"Qos profile message is empty",
			codes.InvalidArgument,
			false,
		},
		{
			"Test for Cancelled context",
			&netops.Response{StatusMsg: ""},
			&netops.SliceQosProfile{SliceName: "test-slice", SliceId: "12"},
			"context canceled",
			codes.Canceled,
			true,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	logger.GlobalLogger = logger.NewLogger("ERROR")
	err := MockBootstrapNetOpPod()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// var emptyProfile *netops.SliceQosProfile
	client := netops.NewNetOpsServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.testCase, func(t *testing.T) {
			request := tt.sqos
			if tt.ctxCanel {
				cancel()
			}
			response, err := client.UpdateSliceQosProfile(ctx, request)
			if response != nil {
				t.Log(response)
				if response.StatusMsg != tt.res.StatusMsg {
					t.Error("response: expected", tt.res.StatusMsg, "received", response.StatusMsg)
				}
			}
			if err != nil {
				t.Log("in ERROR")
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}

}

func TestUpdateSliceQosProfile(t *testing.T) {
	tests := []struct {
		testCase string
		res      *netops.Response
		sqos     *netops.SliceQosProfile
		errMsg   string
		errCode  codes.Code
		ctxCanel bool
	}{
		{
			"Test for successfully enforcing slice QoS policy",
			&netops.Response{StatusMsg: "Slice QoS policy enforced successfully"},
			&netops.SliceQosProfile{
				SliceName:      "test-slice",
				SliceId:        "randomid",
				QosProfileName: "",
				TcType:         netops.TcType_BANDWIDTH_CONTROL,
				ClassType:      netops.ClassType_HTB,
				BwCeiling:      1,
				BwGuaranteed:   htbRootHandleId,
				Priority:       2,
				DscpClass:      "",
			},
			"",
			codes.OK,
			false,
		},
		{
			"Empty Slice Qos Profile",
			&netops.Response{StatusMsg: ""},
			&netops.SliceQosProfile{},
			"Qos profile message is empty",
			codes.InvalidArgument,
			false,
		},
		{
			"Test for Cancelled context",
			&netops.Response{StatusMsg: ""},
			&netops.SliceQosProfile{SliceName: "12"},
			"context canceled",
			codes.Canceled,
			true,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	logger.GlobalLogger = logger.NewLogger("ERROR")
	err := BootstrapNetOpPod()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// var emptyProfile *netops.SliceQosProfile
	client := netops.NewNetOpsServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.testCase, func(t *testing.T) {
			request := tt.sqos
			if tt.ctxCanel {
				cancel()
			}
			response, err := client.UpdateSliceQosProfile(ctx, request)
			if response != nil {
				t.Log(response)
				if response.StatusMsg != tt.res.StatusMsg {
					t.Error("response: expected", tt.res.StatusMsg, "received", response.StatusMsg)
				}
			}
			if err != nil {
				t.Log("in ERROR")
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}

}

func TestUpdateSliceLifeCycleEvent(t *testing.T) {
	tests := []struct {
		testCase string
		res      *netops.Response
		slce     *netops.SliceLifeCycleEvent
		errMsg   string
		errCode  codes.Code
		ctxCanel bool
	}{
		{
			"Test for Slice Create Event",
			&netops.Response{StatusMsg: "Slice life cycle event handled successfully"},
			&netops.SliceLifeCycleEvent{
				SliceName: "test-slice",
				Event:     netops.EventType_EV_CREATE,
			},
			"",
			codes.OK,
			false,
		},
		{
			"Test for Slice Delete Event",
			&netops.Response{StatusMsg: "Slice life cycle event handled successfully"},
			&netops.SliceLifeCycleEvent{
				SliceName: "test-slice",
				Event:     netops.EventType_EV_DELETE,
			},
			"",
			codes.OK,
			false,
		},
		{
			"Test for Slice Update Event",
			&netops.Response{StatusMsg: "Slice life cycle event handled successfully"},
			&netops.SliceLifeCycleEvent{
				SliceName: "test-slice",
				Event:     netops.EventType_EV_UPDATE,
			},
			"",
			codes.OK,
			false,
		},
		{
			"Empty Slice Event",
			&netops.Response{StatusMsg: ""},
			&netops.SliceLifeCycleEvent{SliceName: ""},
			"Slice lifecycle message is empty",
			codes.InvalidArgument,
			false,
		},
		{
			"Test for Cancelled context",
			&netops.Response{StatusMsg: ""},
			&netops.SliceLifeCycleEvent{
				SliceName: netIface,
				Event:     netops.EventType_EV_CREATE,
			},
			"context canceled",
			codes.Canceled,
			true,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	logger.GlobalLogger = logger.NewLogger("DEBUG")
	err := MockBootstrapNetOpPod()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := netops.NewNetOpsServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.testCase, func(t *testing.T) {
			request := tt.slce
			if tt.ctxCanel {
				cancel()
			}
			response, err := client.UpdateSliceLifeCycleEvent(ctx, request)
			if response != nil {
				if response.StatusMsg != tt.res.StatusMsg {
					t.Error("response: expected", tt.res.StatusMsg, "received", response.StatusMsg)
				}
			}
			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}

}
