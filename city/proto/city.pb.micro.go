// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: city.proto

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

// Api Endpoints for CityService service

func NewCityServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for CityService service

type CityService interface {
	//City Methods
	CreateCity(ctx context.Context, in *CreateCityRequest, opts ...client.CallOption) (*CreateCityResponse, error)
	UpdateCity(ctx context.Context, in *UpdateCityRequest, opts ...client.CallOption) (*UpdateResponse, error)
	GetCities(ctx context.Context, in *GetCityRequest, opts ...client.CallOption) (*GetCitiesResponse, error)
	GetCityById(ctx context.Context, in *GetCityRequest, opts ...client.CallOption) (*GetCityResponse, error)
	GetCityByGeoPoint(ctx context.Context, in *GetCityRequest, opts ...client.CallOption) (*GetCityResponse, error)
	//Fences Methods
	CreateFence(ctx context.Context, in *CreateFenceRequest, opts ...client.CallOption) (*CreateFenceResponse, error)
	UpdateFence(ctx context.Context, in *UpdateFenceRequest, opts ...client.CallOption) (*UpdateResponse, error)
	GetFences(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFencesResponse, error)
	GetFenceById(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFenceResponse, error)
	GetFenceByGeoPoint(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFenceResponse, error)
	GetFenceByGeoPointAndId(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFenceResponse, error)
	GetHubFences(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFencesResponse, error)
	//City Settings
	SetDeliveryDays(ctx context.Context, in *SetDeliveryDaysRequest, opts ...client.CallOption) (*UpdateResponse, error)
	GetDeliveryDays(ctx context.Context, in *GetCityRequest, opts ...client.CallOption) (*GetDeliveryDaysResponse, error)
}

type cityService struct {
	c    client.Client
	name string
}

func NewCityService(name string, c client.Client) CityService {
	return &cityService{
		c:    c,
		name: name,
	}
}

func (c *cityService) CreateCity(ctx context.Context, in *CreateCityRequest, opts ...client.CallOption) (*CreateCityResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.CreateCity", in)
	out := new(CreateCityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) UpdateCity(ctx context.Context, in *UpdateCityRequest, opts ...client.CallOption) (*UpdateResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.UpdateCity", in)
	out := new(UpdateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetCities(ctx context.Context, in *GetCityRequest, opts ...client.CallOption) (*GetCitiesResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetCities", in)
	out := new(GetCitiesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetCityById(ctx context.Context, in *GetCityRequest, opts ...client.CallOption) (*GetCityResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetCityById", in)
	out := new(GetCityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetCityByGeoPoint(ctx context.Context, in *GetCityRequest, opts ...client.CallOption) (*GetCityResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetCityByGeoPoint", in)
	out := new(GetCityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) CreateFence(ctx context.Context, in *CreateFenceRequest, opts ...client.CallOption) (*CreateFenceResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.CreateFence", in)
	out := new(CreateFenceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) UpdateFence(ctx context.Context, in *UpdateFenceRequest, opts ...client.CallOption) (*UpdateResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.UpdateFence", in)
	out := new(UpdateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetFences(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFencesResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetFences", in)
	out := new(GetFencesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetFenceById(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFenceResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetFenceById", in)
	out := new(GetFenceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetFenceByGeoPoint(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFenceResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetFenceByGeoPoint", in)
	out := new(GetFenceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetFenceByGeoPointAndId(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFenceResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetFenceByGeoPointAndId", in)
	out := new(GetFenceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetHubFences(ctx context.Context, in *GetFenceRequest, opts ...client.CallOption) (*GetFencesResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetHubFences", in)
	out := new(GetFencesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) SetDeliveryDays(ctx context.Context, in *SetDeliveryDaysRequest, opts ...client.CallOption) (*UpdateResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.SetDeliveryDays", in)
	out := new(UpdateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityService) GetDeliveryDays(ctx context.Context, in *GetCityRequest, opts ...client.CallOption) (*GetDeliveryDaysResponse, error) {
	req := c.c.NewRequest(c.name, "CityService.GetDeliveryDays", in)
	out := new(GetDeliveryDaysResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CityService service

type CityServiceHandler interface {
	//City Methods
	CreateCity(context.Context, *CreateCityRequest, *CreateCityResponse) error
	UpdateCity(context.Context, *UpdateCityRequest, *UpdateResponse) error
	GetCities(context.Context, *GetCityRequest, *GetCitiesResponse) error
	GetCityById(context.Context, *GetCityRequest, *GetCityResponse) error
	GetCityByGeoPoint(context.Context, *GetCityRequest, *GetCityResponse) error
	//Fences Methods
	CreateFence(context.Context, *CreateFenceRequest, *CreateFenceResponse) error
	UpdateFence(context.Context, *UpdateFenceRequest, *UpdateResponse) error
	GetFences(context.Context, *GetFenceRequest, *GetFencesResponse) error
	GetFenceById(context.Context, *GetFenceRequest, *GetFenceResponse) error
	GetFenceByGeoPoint(context.Context, *GetFenceRequest, *GetFenceResponse) error
	GetFenceByGeoPointAndId(context.Context, *GetFenceRequest, *GetFenceResponse) error
	GetHubFences(context.Context, *GetFenceRequest, *GetFencesResponse) error
	//City Settings
	SetDeliveryDays(context.Context, *SetDeliveryDaysRequest, *UpdateResponse) error
	GetDeliveryDays(context.Context, *GetCityRequest, *GetDeliveryDaysResponse) error
}

func RegisterCityServiceHandler(s server.Server, hdlr CityServiceHandler, opts ...server.HandlerOption) error {
	type cityService interface {
		CreateCity(ctx context.Context, in *CreateCityRequest, out *CreateCityResponse) error
		UpdateCity(ctx context.Context, in *UpdateCityRequest, out *UpdateResponse) error
		GetCities(ctx context.Context, in *GetCityRequest, out *GetCitiesResponse) error
		GetCityById(ctx context.Context, in *GetCityRequest, out *GetCityResponse) error
		GetCityByGeoPoint(ctx context.Context, in *GetCityRequest, out *GetCityResponse) error
		CreateFence(ctx context.Context, in *CreateFenceRequest, out *CreateFenceResponse) error
		UpdateFence(ctx context.Context, in *UpdateFenceRequest, out *UpdateResponse) error
		GetFences(ctx context.Context, in *GetFenceRequest, out *GetFencesResponse) error
		GetFenceById(ctx context.Context, in *GetFenceRequest, out *GetFenceResponse) error
		GetFenceByGeoPoint(ctx context.Context, in *GetFenceRequest, out *GetFenceResponse) error
		GetFenceByGeoPointAndId(ctx context.Context, in *GetFenceRequest, out *GetFenceResponse) error
		GetHubFences(ctx context.Context, in *GetFenceRequest, out *GetFencesResponse) error
		SetDeliveryDays(ctx context.Context, in *SetDeliveryDaysRequest, out *UpdateResponse) error
		GetDeliveryDays(ctx context.Context, in *GetCityRequest, out *GetDeliveryDaysResponse) error
	}
	type CityService struct {
		cityService
	}
	h := &cityServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&CityService{h}, opts...))
}

type cityServiceHandler struct {
	CityServiceHandler
}

func (h *cityServiceHandler) CreateCity(ctx context.Context, in *CreateCityRequest, out *CreateCityResponse) error {
	return h.CityServiceHandler.CreateCity(ctx, in, out)
}

func (h *cityServiceHandler) UpdateCity(ctx context.Context, in *UpdateCityRequest, out *UpdateResponse) error {
	return h.CityServiceHandler.UpdateCity(ctx, in, out)
}

func (h *cityServiceHandler) GetCities(ctx context.Context, in *GetCityRequest, out *GetCitiesResponse) error {
	return h.CityServiceHandler.GetCities(ctx, in, out)
}

func (h *cityServiceHandler) GetCityById(ctx context.Context, in *GetCityRequest, out *GetCityResponse) error {
	return h.CityServiceHandler.GetCityById(ctx, in, out)
}

func (h *cityServiceHandler) GetCityByGeoPoint(ctx context.Context, in *GetCityRequest, out *GetCityResponse) error {
	return h.CityServiceHandler.GetCityByGeoPoint(ctx, in, out)
}

func (h *cityServiceHandler) CreateFence(ctx context.Context, in *CreateFenceRequest, out *CreateFenceResponse) error {
	return h.CityServiceHandler.CreateFence(ctx, in, out)
}

func (h *cityServiceHandler) UpdateFence(ctx context.Context, in *UpdateFenceRequest, out *UpdateResponse) error {
	return h.CityServiceHandler.UpdateFence(ctx, in, out)
}

func (h *cityServiceHandler) GetFences(ctx context.Context, in *GetFenceRequest, out *GetFencesResponse) error {
	return h.CityServiceHandler.GetFences(ctx, in, out)
}

func (h *cityServiceHandler) GetFenceById(ctx context.Context, in *GetFenceRequest, out *GetFenceResponse) error {
	return h.CityServiceHandler.GetFenceById(ctx, in, out)
}

func (h *cityServiceHandler) GetFenceByGeoPoint(ctx context.Context, in *GetFenceRequest, out *GetFenceResponse) error {
	return h.CityServiceHandler.GetFenceByGeoPoint(ctx, in, out)
}

func (h *cityServiceHandler) GetFenceByGeoPointAndId(ctx context.Context, in *GetFenceRequest, out *GetFenceResponse) error {
	return h.CityServiceHandler.GetFenceByGeoPointAndId(ctx, in, out)
}

func (h *cityServiceHandler) GetHubFences(ctx context.Context, in *GetFenceRequest, out *GetFencesResponse) error {
	return h.CityServiceHandler.GetHubFences(ctx, in, out)
}

func (h *cityServiceHandler) SetDeliveryDays(ctx context.Context, in *SetDeliveryDaysRequest, out *UpdateResponse) error {
	return h.CityServiceHandler.SetDeliveryDays(ctx, in, out)
}

func (h *cityServiceHandler) GetDeliveryDays(ctx context.Context, in *GetCityRequest, out *GetDeliveryDaysResponse) error {
	return h.CityServiceHandler.GetDeliveryDays(ctx, in, out)
}
