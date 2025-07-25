// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: account.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AccountService_GetKakaoLoginURL_FullMethodName = "/go.escape.ship.proto.v1.AccountService/GetKakaoLoginURL"
	AccountService_GetKakaoCallBack_FullMethodName = "/go.escape.ship.proto.v1.AccountService/GetKakaoCallBack"
	AccountService_Login_FullMethodName            = "/go.escape.ship.proto.v1.AccountService/Login"
	AccountService_Register_FullMethodName         = "/go.escape.ship.proto.v1.AccountService/Register"
)

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	GetKakaoLoginURL(ctx context.Context, in *GetKakaoLoginURLRequest, opts ...grpc.CallOption) (*GetKakaoLoginURLResponse, error)
	GetKakaoCallBack(ctx context.Context, in *GetKakaoCallBackRequest, opts ...grpc.CallOption) (*GetKakaoCallBackResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) GetKakaoLoginURL(ctx context.Context, in *GetKakaoLoginURLRequest, opts ...grpc.CallOption) (*GetKakaoLoginURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetKakaoLoginURLResponse)
	err := c.cc.Invoke(ctx, AccountService_GetKakaoLoginURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetKakaoCallBack(ctx context.Context, in *GetKakaoCallBackRequest, opts ...grpc.CallOption) (*GetKakaoCallBackResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetKakaoCallBackResponse)
	err := c.cc.Invoke(ctx, AccountService_GetKakaoCallBack_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AccountService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, AccountService_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations must embed UnimplementedAccountServiceServer
// for forward compatibility.
type AccountServiceServer interface {
	GetKakaoLoginURL(context.Context, *GetKakaoLoginURLRequest) (*GetKakaoLoginURLResponse, error)
	GetKakaoCallBack(context.Context, *GetKakaoCallBackRequest) (*GetKakaoCallBackResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	mustEmbedUnimplementedAccountServiceServer()
}

// UnimplementedAccountServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAccountServiceServer struct{}

func (UnimplementedAccountServiceServer) GetKakaoLoginURL(context.Context, *GetKakaoLoginURLRequest) (*GetKakaoLoginURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKakaoLoginURL not implemented")
}
func (UnimplementedAccountServiceServer) GetKakaoCallBack(context.Context, *GetKakaoCallBackRequest) (*GetKakaoCallBackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKakaoCallBack not implemented")
}
func (UnimplementedAccountServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAccountServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {}
func (UnimplementedAccountServiceServer) testEmbeddedByValue()                        {}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	// If the following call pancis, it indicates UnimplementedAccountServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_GetKakaoLoginURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKakaoLoginURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetKakaoLoginURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetKakaoLoginURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetKakaoLoginURL(ctx, req.(*GetKakaoLoginURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetKakaoCallBack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKakaoCallBackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetKakaoCallBack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetKakaoCallBack_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetKakaoCallBack(ctx, req.(*GetKakaoCallBackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go.escape.ship.proto.v1.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetKakaoLoginURL",
			Handler:    _AccountService_GetKakaoLoginURL_Handler,
		},
		{
			MethodName: "GetKakaoCallBack",
			Handler:    _AccountService_GetKakaoCallBack_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AccountService_Login_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _AccountService_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account.proto",
}
