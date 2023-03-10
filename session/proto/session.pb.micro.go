// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: session.proto

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

// Api Endpoints for SessionService service

func NewSessionServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for SessionService service

type SessionService interface {
	GetSessionById(ctx context.Context, in *GetSessionByIdRequest, opts ...client.CallOption) (*GetSessionResponse, error)
	GetSession(ctx context.Context, in *GetSessionRequest, opts ...client.CallOption) (*GetSessionResponse, error)
	GetUserSession(ctx context.Context, in *GetSessionRequest, opts ...client.CallOption) (*GetSessionsResponse, error)
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...client.CallOption) (*GetSessionResponse, error)
}

type sessionService struct {
	c    client.Client
	name string
}

func NewSessionService(name string, c client.Client) SessionService {
	return &sessionService{
		c:    c,
		name: name,
	}
}

func (c *sessionService) GetSessionById(ctx context.Context, in *GetSessionByIdRequest, opts ...client.CallOption) (*GetSessionResponse, error) {
	req := c.c.NewRequest(c.name, "SessionService.GetSessionById", in)
	out := new(GetSessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionService) GetSession(ctx context.Context, in *GetSessionRequest, opts ...client.CallOption) (*GetSessionResponse, error) {
	req := c.c.NewRequest(c.name, "SessionService.GetSession", in)
	out := new(GetSessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionService) GetUserSession(ctx context.Context, in *GetSessionRequest, opts ...client.CallOption) (*GetSessionsResponse, error) {
	req := c.c.NewRequest(c.name, "SessionService.GetUserSession", in)
	out := new(GetSessionsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionService) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...client.CallOption) (*GetSessionResponse, error) {
	req := c.c.NewRequest(c.name, "SessionService.CreateSession", in)
	out := new(GetSessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SessionService service

type SessionServiceHandler interface {
	GetSessionById(context.Context, *GetSessionByIdRequest, *GetSessionResponse) error
	GetSession(context.Context, *GetSessionRequest, *GetSessionResponse) error
	GetUserSession(context.Context, *GetSessionRequest, *GetSessionsResponse) error
	CreateSession(context.Context, *CreateSessionRequest, *GetSessionResponse) error
}

func RegisterSessionServiceHandler(s server.Server, hdlr SessionServiceHandler, opts ...server.HandlerOption) error {
	type sessionService interface {
		GetSessionById(ctx context.Context, in *GetSessionByIdRequest, out *GetSessionResponse) error
		GetSession(ctx context.Context, in *GetSessionRequest, out *GetSessionResponse) error
		GetUserSession(ctx context.Context, in *GetSessionRequest, out *GetSessionsResponse) error
		CreateSession(ctx context.Context, in *CreateSessionRequest, out *GetSessionResponse) error
	}
	type SessionService struct {
		sessionService
	}
	h := &sessionServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&SessionService{h}, opts...))
}

type sessionServiceHandler struct {
	SessionServiceHandler
}

func (h *sessionServiceHandler) GetSessionById(ctx context.Context, in *GetSessionByIdRequest, out *GetSessionResponse) error {
	return h.SessionServiceHandler.GetSessionById(ctx, in, out)
}

func (h *sessionServiceHandler) GetSession(ctx context.Context, in *GetSessionRequest, out *GetSessionResponse) error {
	return h.SessionServiceHandler.GetSession(ctx, in, out)
}

func (h *sessionServiceHandler) GetUserSession(ctx context.Context, in *GetSessionRequest, out *GetSessionsResponse) error {
	return h.SessionServiceHandler.GetUserSession(ctx, in, out)
}

func (h *sessionServiceHandler) CreateSession(ctx context.Context, in *CreateSessionRequest, out *GetSessionResponse) error {
	return h.SessionServiceHandler.CreateSession(ctx, in, out)
}
