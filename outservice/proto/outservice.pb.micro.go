// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: outservice.proto

package justify

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for OutService service

func NewOutServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for OutService service

type OutService interface {
	SendSMS(ctx context.Context, in *SendSMSRequest, opts ...client.CallOption) (*SendSMSResponse, error)
	SendMobileNotification(ctx context.Context, in *SendNotificaitonRequest, opts ...client.CallOption) (*SendSMSResponse, error)
}

type outService struct {
	c    client.Client
	name string
}

func NewOutService(name string, c client.Client) OutService {
	return &outService{
		c:    c,
		name: name,
	}
}

func (c *outService) SendSMS(ctx context.Context, in *SendSMSRequest, opts ...client.CallOption) (*SendSMSResponse, error) {
	req := c.c.NewRequest(c.name, "OutService.SendSMS", in)
	out := new(SendSMSResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *outService) SendMobileNotification(ctx context.Context, in *SendNotificaitonRequest, opts ...client.CallOption) (*SendSMSResponse, error) {
	req := c.c.NewRequest(c.name, "OutService.SendMobileNotification", in)
	out := new(SendSMSResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OutService service

type OutServiceHandler interface {
	SendSMS(context.Context, *SendSMSRequest, *SendSMSResponse) error
	SendMobileNotification(context.Context, *SendNotificaitonRequest, *SendSMSResponse) error
}

func RegisterOutServiceHandler(s server.Server, hdlr OutServiceHandler, opts ...server.HandlerOption) error {
	type outService interface {
		SendSMS(ctx context.Context, in *SendSMSRequest, out *SendSMSResponse) error
		SendMobileNotification(ctx context.Context, in *SendNotificaitonRequest, out *SendSMSResponse) error
	}
	type OutService struct {
		outService
	}
	h := &outServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&OutService{h}, opts...))
}

type outServiceHandler struct {
	OutServiceHandler
}

func (h *outServiceHandler) SendSMS(ctx context.Context, in *SendSMSRequest, out *SendSMSResponse) error {
	return h.OutServiceHandler.SendSMS(ctx, in, out)
}

func (h *outServiceHandler) SendMobileNotification(ctx context.Context, in *SendNotificaitonRequest, out *SendSMSResponse) error {
	return h.OutServiceHandler.SendMobileNotification(ctx, in, out)
}