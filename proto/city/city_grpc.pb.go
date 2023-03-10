// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: city.proto

package justify

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CityServiceClient is the client API for CityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CityServiceClient interface {
	//City Methods
	CreateCity(ctx context.Context, in *CreateCityRequest, opts ...grpc.CallOption) (*CreateCityResponse, error)
	UpdateCity(ctx context.Context, in *UpdateCityRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	GetCities(ctx context.Context, in *GetCityRequest, opts ...grpc.CallOption) (*GetCitiesResponse, error)
	GetCityById(ctx context.Context, in *GetCityRequest, opts ...grpc.CallOption) (*GetCityResponse, error)
	GetCityByGeoPoint(ctx context.Context, in *GetCityRequest, opts ...grpc.CallOption) (*GetCityResponse, error)
	//Fences Methods
	CreateFence(ctx context.Context, in *CreateFenceRequest, opts ...grpc.CallOption) (*CreateFenceResponse, error)
	UpdateFence(ctx context.Context, in *UpdateFenceRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	GetFences(ctx context.Context, in *GetFenceRequest, opts ...grpc.CallOption) (*GetFencesResponse, error)
	GetHubFences(ctx context.Context, in *GetFenceRequest, opts ...grpc.CallOption) (*GetFencesResponse, error)
	GetFenceById(ctx context.Context, in *GetFenceRequest, opts ...grpc.CallOption) (*GetFenceResponse, error)
	GetFenceByGeoPoint(ctx context.Context, in *GetFenceRequest, opts ...grpc.CallOption) (*GetFenceResponse, error)
	//City Settings
	SetDeliveryDays(ctx context.Context, in *SetDeliveryDaysRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	GetDeliveryDays(ctx context.Context, in *GetCityRequest, opts ...grpc.CallOption) (*GetDeliveryDaysResponse, error)
}

type cityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCityServiceClient(cc grpc.ClientConnInterface) CityServiceClient {
	return &cityServiceClient{cc}
}

func (c *cityServiceClient) CreateCity(ctx context.Context, in *CreateCityRequest, opts ...grpc.CallOption) (*CreateCityResponse, error) {
	out := new(CreateCityResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/CreateCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) UpdateCity(ctx context.Context, in *UpdateCityRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/UpdateCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) GetCities(ctx context.Context, in *GetCityRequest, opts ...grpc.CallOption) (*GetCitiesResponse, error) {
	out := new(GetCitiesResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/GetCities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) GetCityById(ctx context.Context, in *GetCityRequest, opts ...grpc.CallOption) (*GetCityResponse, error) {
	out := new(GetCityResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/GetCityById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) GetCityByGeoPoint(ctx context.Context, in *GetCityRequest, opts ...grpc.CallOption) (*GetCityResponse, error) {
	out := new(GetCityResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/GetCityByGeoPoint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) CreateFence(ctx context.Context, in *CreateFenceRequest, opts ...grpc.CallOption) (*CreateFenceResponse, error) {
	out := new(CreateFenceResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/CreateFence", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) UpdateFence(ctx context.Context, in *UpdateFenceRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/UpdateFence", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) GetFences(ctx context.Context, in *GetFenceRequest, opts ...grpc.CallOption) (*GetFencesResponse, error) {
	out := new(GetFencesResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/GetFences", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) GetHubFences(ctx context.Context, in *GetFenceRequest, opts ...grpc.CallOption) (*GetFencesResponse, error) {
	out := new(GetFencesResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/GetHubFences", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) GetFenceById(ctx context.Context, in *GetFenceRequest, opts ...grpc.CallOption) (*GetFenceResponse, error) {
	out := new(GetFenceResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/GetFenceById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) GetFenceByGeoPoint(ctx context.Context, in *GetFenceRequest, opts ...grpc.CallOption) (*GetFenceResponse, error) {
	out := new(GetFenceResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/GetFenceByGeoPoint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) SetDeliveryDays(ctx context.Context, in *SetDeliveryDaysRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/SetDeliveryDays", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cityServiceClient) GetDeliveryDays(ctx context.Context, in *GetCityRequest, opts ...grpc.CallOption) (*GetDeliveryDaysResponse, error) {
	out := new(GetDeliveryDaysResponse)
	err := c.cc.Invoke(ctx, "/justify.CityService/GetDeliveryDays", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CityServiceServer is the server API for CityService service.
// All implementations must embed UnimplementedCityServiceServer
// for forward compatibility
type CityServiceServer interface {
	//City Methods
	CreateCity(context.Context, *CreateCityRequest) (*CreateCityResponse, error)
	UpdateCity(context.Context, *UpdateCityRequest) (*UpdateResponse, error)
	GetCities(context.Context, *GetCityRequest) (*GetCitiesResponse, error)
	GetCityById(context.Context, *GetCityRequest) (*GetCityResponse, error)
	GetCityByGeoPoint(context.Context, *GetCityRequest) (*GetCityResponse, error)
	//Fences Methods
	CreateFence(context.Context, *CreateFenceRequest) (*CreateFenceResponse, error)
	UpdateFence(context.Context, *UpdateFenceRequest) (*UpdateResponse, error)
	GetFences(context.Context, *GetFenceRequest) (*GetFencesResponse, error)
	GetHubFences(context.Context, *GetFenceRequest) (*GetFencesResponse, error)
	GetFenceById(context.Context, *GetFenceRequest) (*GetFenceResponse, error)
	GetFenceByGeoPoint(context.Context, *GetFenceRequest) (*GetFenceResponse, error)
	//City Settings
	SetDeliveryDays(context.Context, *SetDeliveryDaysRequest) (*UpdateResponse, error)
	GetDeliveryDays(context.Context, *GetCityRequest) (*GetDeliveryDaysResponse, error)
	mustEmbedUnimplementedCityServiceServer()
}

// UnimplementedCityServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCityServiceServer struct {
}

func (UnimplementedCityServiceServer) CreateCity(context.Context, *CreateCityRequest) (*CreateCityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCity not implemented")
}
func (UnimplementedCityServiceServer) UpdateCity(context.Context, *UpdateCityRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCity not implemented")
}
func (UnimplementedCityServiceServer) GetCities(context.Context, *GetCityRequest) (*GetCitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCities not implemented")
}
func (UnimplementedCityServiceServer) GetCityById(context.Context, *GetCityRequest) (*GetCityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCityById not implemented")
}
func (UnimplementedCityServiceServer) GetCityByGeoPoint(context.Context, *GetCityRequest) (*GetCityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCityByGeoPoint not implemented")
}
func (UnimplementedCityServiceServer) CreateFence(context.Context, *CreateFenceRequest) (*CreateFenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFence not implemented")
}
func (UnimplementedCityServiceServer) UpdateFence(context.Context, *UpdateFenceRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFence not implemented")
}
func (UnimplementedCityServiceServer) GetFences(context.Context, *GetFenceRequest) (*GetFencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFences not implemented")
}
func (UnimplementedCityServiceServer) GetHubFences(context.Context, *GetFenceRequest) (*GetFencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHubFences not implemented")
}
func (UnimplementedCityServiceServer) GetFenceById(context.Context, *GetFenceRequest) (*GetFenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFenceById not implemented")
}
func (UnimplementedCityServiceServer) GetFenceByGeoPoint(context.Context, *GetFenceRequest) (*GetFenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFenceByGeoPoint not implemented")
}
func (UnimplementedCityServiceServer) SetDeliveryDays(context.Context, *SetDeliveryDaysRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDeliveryDays not implemented")
}
func (UnimplementedCityServiceServer) GetDeliveryDays(context.Context, *GetCityRequest) (*GetDeliveryDaysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeliveryDays not implemented")
}
func (UnimplementedCityServiceServer) mustEmbedUnimplementedCityServiceServer() {}

// UnsafeCityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CityServiceServer will
// result in compilation errors.
type UnsafeCityServiceServer interface {
	mustEmbedUnimplementedCityServiceServer()
}

func RegisterCityServiceServer(s grpc.ServiceRegistrar, srv CityServiceServer) {
	s.RegisterService(&CityService_ServiceDesc, srv)
}

func _CityService_CreateCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).CreateCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/CreateCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).CreateCity(ctx, req.(*CreateCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_UpdateCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).UpdateCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/UpdateCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).UpdateCity(ctx, req.(*UpdateCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_GetCities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).GetCities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/GetCities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).GetCities(ctx, req.(*GetCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_GetCityById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).GetCityById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/GetCityById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).GetCityById(ctx, req.(*GetCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_GetCityByGeoPoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).GetCityByGeoPoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/GetCityByGeoPoint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).GetCityByGeoPoint(ctx, req.(*GetCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_CreateFence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).CreateFence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/CreateFence",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).CreateFence(ctx, req.(*CreateFenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_UpdateFence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).UpdateFence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/UpdateFence",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).UpdateFence(ctx, req.(*UpdateFenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_GetFences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).GetFences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/GetFences",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).GetFences(ctx, req.(*GetFenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_GetHubFences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).GetHubFences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/GetHubFences",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).GetHubFences(ctx, req.(*GetFenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_GetFenceById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).GetFenceById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/GetFenceById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).GetFenceById(ctx, req.(*GetFenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_GetFenceByGeoPoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).GetFenceByGeoPoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/GetFenceByGeoPoint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).GetFenceByGeoPoint(ctx, req.(*GetFenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_SetDeliveryDays_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDeliveryDaysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).SetDeliveryDays(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/SetDeliveryDays",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).SetDeliveryDays(ctx, req.(*SetDeliveryDaysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CityService_GetDeliveryDays_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CityServiceServer).GetDeliveryDays(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/justify.CityService/GetDeliveryDays",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CityServiceServer).GetDeliveryDays(ctx, req.(*GetCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CityService_ServiceDesc is the grpc.ServiceDesc for CityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "justify.CityService",
	HandlerType: (*CityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCity",
			Handler:    _CityService_CreateCity_Handler,
		},
		{
			MethodName: "UpdateCity",
			Handler:    _CityService_UpdateCity_Handler,
		},
		{
			MethodName: "GetCities",
			Handler:    _CityService_GetCities_Handler,
		},
		{
			MethodName: "GetCityById",
			Handler:    _CityService_GetCityById_Handler,
		},
		{
			MethodName: "GetCityByGeoPoint",
			Handler:    _CityService_GetCityByGeoPoint_Handler,
		},
		{
			MethodName: "CreateFence",
			Handler:    _CityService_CreateFence_Handler,
		},
		{
			MethodName: "UpdateFence",
			Handler:    _CityService_UpdateFence_Handler,
		},
		{
			MethodName: "GetFences",
			Handler:    _CityService_GetFences_Handler,
		},
		{
			MethodName: "GetHubFences",
			Handler:    _CityService_GetHubFences_Handler,
		},
		{
			MethodName: "GetFenceById",
			Handler:    _CityService_GetFenceById_Handler,
		},
		{
			MethodName: "GetFenceByGeoPoint",
			Handler:    _CityService_GetFenceByGeoPoint_Handler,
		},
		{
			MethodName: "SetDeliveryDays",
			Handler:    _CityService_SetDeliveryDays_Handler,
		},
		{
			MethodName: "GetDeliveryDays",
			Handler:    _CityService_GetDeliveryDays_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "city.proto",
}
