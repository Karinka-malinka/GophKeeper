// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
// source: gophkeeper.proto

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	UserService_Register_FullMethodName = "/gophkeeper.UserService/Register"
	UserService_Login_FullMethodName    = "/gophkeeper.UserService/Login"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	Register(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	Login(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Register(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, UserService_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Login(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, UserService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	Register(context.Context, *UserRequest) (*UserResponse, error)
	Login(context.Context, *UserRequest) (*UserResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) Register(context.Context, *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUserServiceServer) Login(context.Context, *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Register(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Login(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gophkeeper.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _UserService_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _UserService_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gophkeeper.proto",
}

const (
	SyncService_ListLoginData_FullMethodName = "/gophkeeper.SyncService/ListLoginData"
	SyncService_ListText_FullMethodName      = "/gophkeeper.SyncService/ListText"
	SyncService_ListFile_FullMethodName      = "/gophkeeper.SyncService/ListFile"
	SyncService_ListBankCard_FullMethodName  = "/gophkeeper.SyncService/ListBankCard"
)

// SyncServiceClient is the client API for SyncService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SyncServiceClient interface {
	ListLoginData(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*LoginDataResponse, error)
	ListText(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*TextResponse, error)
	ListFile(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*FileResponse, error)
	ListBankCard(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BankCardResponse, error)
}

type syncServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSyncServiceClient(cc grpc.ClientConnInterface) SyncServiceClient {
	return &syncServiceClient{cc}
}

func (c *syncServiceClient) ListLoginData(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*LoginDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginDataResponse)
	err := c.cc.Invoke(ctx, SyncService_ListLoginData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *syncServiceClient) ListText(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*TextResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TextResponse)
	err := c.cc.Invoke(ctx, SyncService_ListText_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *syncServiceClient) ListFile(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*FileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, SyncService_ListFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *syncServiceClient) ListBankCard(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BankCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BankCardResponse)
	err := c.cc.Invoke(ctx, SyncService_ListBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SyncServiceServer is the server API for SyncService service.
// All implementations must embed UnimplementedSyncServiceServer
// for forward compatibility
type SyncServiceServer interface {
	ListLoginData(context.Context, *empty.Empty) (*LoginDataResponse, error)
	ListText(context.Context, *empty.Empty) (*TextResponse, error)
	ListFile(context.Context, *empty.Empty) (*FileResponse, error)
	ListBankCard(context.Context, *empty.Empty) (*BankCardResponse, error)
	mustEmbedUnimplementedSyncServiceServer()
}

// UnimplementedSyncServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSyncServiceServer struct {
}

func (UnimplementedSyncServiceServer) ListLoginData(context.Context, *empty.Empty) (*LoginDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLoginData not implemented")
}
func (UnimplementedSyncServiceServer) ListText(context.Context, *empty.Empty) (*TextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListText not implemented")
}
func (UnimplementedSyncServiceServer) ListFile(context.Context, *empty.Empty) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFile not implemented")
}
func (UnimplementedSyncServiceServer) ListBankCard(context.Context, *empty.Empty) (*BankCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBankCard not implemented")
}
func (UnimplementedSyncServiceServer) mustEmbedUnimplementedSyncServiceServer() {}

// UnsafeSyncServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SyncServiceServer will
// result in compilation errors.
type UnsafeSyncServiceServer interface {
	mustEmbedUnimplementedSyncServiceServer()
}

func RegisterSyncServiceServer(s grpc.ServiceRegistrar, srv SyncServiceServer) {
	s.RegisterService(&SyncService_ServiceDesc, srv)
}

func _SyncService_ListLoginData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncServiceServer).ListLoginData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SyncService_ListLoginData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncServiceServer).ListLoginData(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SyncService_ListText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncServiceServer).ListText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SyncService_ListText_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncServiceServer).ListText(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SyncService_ListFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncServiceServer).ListFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SyncService_ListFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncServiceServer).ListFile(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SyncService_ListBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncServiceServer).ListBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SyncService_ListBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncServiceServer).ListBankCard(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// SyncService_ServiceDesc is the grpc.ServiceDesc for SyncService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SyncService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gophkeeper.SyncService",
	HandlerType: (*SyncServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListLoginData",
			Handler:    _SyncService_ListLoginData_Handler,
		},
		{
			MethodName: "ListText",
			Handler:    _SyncService_ListText_Handler,
		},
		{
			MethodName: "ListFile",
			Handler:    _SyncService_ListFile_Handler,
		},
		{
			MethodName: "ListBankCard",
			Handler:    _SyncService_ListBankCard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gophkeeper.proto",
}

const (
	ManagementService_AddLoginData_FullMethodName    = "/gophkeeper.ManagementService/AddLoginData"
	ManagementService_AddText_FullMethodName         = "/gophkeeper.ManagementService/AddText"
	ManagementService_AddFile_FullMethodName         = "/gophkeeper.ManagementService/AddFile"
	ManagementService_AddBankCard_FullMethodName     = "/gophkeeper.ManagementService/AddBankCard"
	ManagementService_EditLoginData_FullMethodName   = "/gophkeeper.ManagementService/EditLoginData"
	ManagementService_GetFile_FullMethodName         = "/gophkeeper.ManagementService/GetFile"
	ManagementService_DeleteLoginData_FullMethodName = "/gophkeeper.ManagementService/DeleteLoginData"
	ManagementService_DeleteText_FullMethodName      = "/gophkeeper.ManagementService/DeleteText"
	ManagementService_DeleteFile_FullMethodName      = "/gophkeeper.ManagementService/DeleteFile"
	ManagementService_DeleteBankCard_FullMethodName  = "/gophkeeper.ManagementService/DeleteBankCard"
)

// ManagementServiceClient is the client API for ManagementService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ManagementServiceClient interface {
	AddLoginData(ctx context.Context, in *LoginData, opts ...grpc.CallOption) (*LoginData, error)
	AddText(ctx context.Context, in *Text, opts ...grpc.CallOption) (*Text, error)
	AddFile(ctx context.Context, in *File, opts ...grpc.CallOption) (*File, error)
	AddBankCard(ctx context.Context, in *BankCard, opts ...grpc.CallOption) (*BankCard, error)
	EditLoginData(ctx context.Context, in *LoginData, opts ...grpc.CallOption) (*empty.Empty, error)
	GetFile(ctx context.Context, in *File, opts ...grpc.CallOption) (*File, error)
	DeleteLoginData(ctx context.Context, in *LoginData, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteText(ctx context.Context, in *Text, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteFile(ctx context.Context, in *File, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteBankCard(ctx context.Context, in *BankCard, opts ...grpc.CallOption) (*empty.Empty, error)
}

type managementServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewManagementServiceClient(cc grpc.ClientConnInterface) ManagementServiceClient {
	return &managementServiceClient{cc}
}

func (c *managementServiceClient) AddLoginData(ctx context.Context, in *LoginData, opts ...grpc.CallOption) (*LoginData, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginData)
	err := c.cc.Invoke(ctx, ManagementService_AddLoginData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) AddText(ctx context.Context, in *Text, opts ...grpc.CallOption) (*Text, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Text)
	err := c.cc.Invoke(ctx, ManagementService_AddText_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) AddFile(ctx context.Context, in *File, opts ...grpc.CallOption) (*File, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(File)
	err := c.cc.Invoke(ctx, ManagementService_AddFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) AddBankCard(ctx context.Context, in *BankCard, opts ...grpc.CallOption) (*BankCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BankCard)
	err := c.cc.Invoke(ctx, ManagementService_AddBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) EditLoginData(ctx context.Context, in *LoginData, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, ManagementService_EditLoginData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) GetFile(ctx context.Context, in *File, opts ...grpc.CallOption) (*File, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(File)
	err := c.cc.Invoke(ctx, ManagementService_GetFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) DeleteLoginData(ctx context.Context, in *LoginData, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, ManagementService_DeleteLoginData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) DeleteText(ctx context.Context, in *Text, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, ManagementService_DeleteText_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) DeleteFile(ctx context.Context, in *File, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, ManagementService_DeleteFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) DeleteBankCard(ctx context.Context, in *BankCard, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, ManagementService_DeleteBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ManagementServiceServer is the server API for ManagementService service.
// All implementations must embed UnimplementedManagementServiceServer
// for forward compatibility
type ManagementServiceServer interface {
	AddLoginData(context.Context, *LoginData) (*LoginData, error)
	AddText(context.Context, *Text) (*Text, error)
	AddFile(context.Context, *File) (*File, error)
	AddBankCard(context.Context, *BankCard) (*BankCard, error)
	EditLoginData(context.Context, *LoginData) (*empty.Empty, error)
	GetFile(context.Context, *File) (*File, error)
	DeleteLoginData(context.Context, *LoginData) (*empty.Empty, error)
	DeleteText(context.Context, *Text) (*empty.Empty, error)
	DeleteFile(context.Context, *File) (*empty.Empty, error)
	DeleteBankCard(context.Context, *BankCard) (*empty.Empty, error)
	mustEmbedUnimplementedManagementServiceServer()
}

// UnimplementedManagementServiceServer must be embedded to have forward compatible implementations.
type UnimplementedManagementServiceServer struct {
}

func (UnimplementedManagementServiceServer) AddLoginData(context.Context, *LoginData) (*LoginData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLoginData not implemented")
}
func (UnimplementedManagementServiceServer) AddText(context.Context, *Text) (*Text, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddText not implemented")
}
func (UnimplementedManagementServiceServer) AddFile(context.Context, *File) (*File, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddFile not implemented")
}
func (UnimplementedManagementServiceServer) AddBankCard(context.Context, *BankCard) (*BankCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBankCard not implemented")
}
func (UnimplementedManagementServiceServer) EditLoginData(context.Context, *LoginData) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditLoginData not implemented")
}
func (UnimplementedManagementServiceServer) GetFile(context.Context, *File) (*File, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (UnimplementedManagementServiceServer) DeleteLoginData(context.Context, *LoginData) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLoginData not implemented")
}
func (UnimplementedManagementServiceServer) DeleteText(context.Context, *Text) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteText not implemented")
}
func (UnimplementedManagementServiceServer) DeleteFile(context.Context, *File) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedManagementServiceServer) DeleteBankCard(context.Context, *BankCard) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBankCard not implemented")
}
func (UnimplementedManagementServiceServer) mustEmbedUnimplementedManagementServiceServer() {}

// UnsafeManagementServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ManagementServiceServer will
// result in compilation errors.
type UnsafeManagementServiceServer interface {
	mustEmbedUnimplementedManagementServiceServer()
}

func RegisterManagementServiceServer(s grpc.ServiceRegistrar, srv ManagementServiceServer) {
	s.RegisterService(&ManagementService_ServiceDesc, srv)
}

func _ManagementService_AddLoginData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).AddLoginData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_AddLoginData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).AddLoginData(ctx, req.(*LoginData))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_AddText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Text)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).AddText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_AddText_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).AddText(ctx, req.(*Text))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_AddFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(File)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).AddFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_AddFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).AddFile(ctx, req.(*File))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_AddBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BankCard)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).AddBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_AddBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).AddBankCard(ctx, req.(*BankCard))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_EditLoginData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).EditLoginData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_EditLoginData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).EditLoginData(ctx, req.(*LoginData))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_GetFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(File)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).GetFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_GetFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).GetFile(ctx, req.(*File))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_DeleteLoginData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).DeleteLoginData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_DeleteLoginData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).DeleteLoginData(ctx, req.(*LoginData))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_DeleteText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Text)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).DeleteText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_DeleteText_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).DeleteText(ctx, req.(*Text))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_DeleteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(File)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).DeleteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_DeleteFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).DeleteFile(ctx, req.(*File))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_DeleteBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BankCard)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).DeleteBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_DeleteBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).DeleteBankCard(ctx, req.(*BankCard))
	}
	return interceptor(ctx, in, info, handler)
}

// ManagementService_ServiceDesc is the grpc.ServiceDesc for ManagementService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ManagementService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gophkeeper.ManagementService",
	HandlerType: (*ManagementServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddLoginData",
			Handler:    _ManagementService_AddLoginData_Handler,
		},
		{
			MethodName: "AddText",
			Handler:    _ManagementService_AddText_Handler,
		},
		{
			MethodName: "AddFile",
			Handler:    _ManagementService_AddFile_Handler,
		},
		{
			MethodName: "AddBankCard",
			Handler:    _ManagementService_AddBankCard_Handler,
		},
		{
			MethodName: "EditLoginData",
			Handler:    _ManagementService_EditLoginData_Handler,
		},
		{
			MethodName: "GetFile",
			Handler:    _ManagementService_GetFile_Handler,
		},
		{
			MethodName: "DeleteLoginData",
			Handler:    _ManagementService_DeleteLoginData_Handler,
		},
		{
			MethodName: "DeleteText",
			Handler:    _ManagementService_DeleteText_Handler,
		},
		{
			MethodName: "DeleteFile",
			Handler:    _ManagementService_DeleteFile_Handler,
		},
		{
			MethodName: "DeleteBankCard",
			Handler:    _ManagementService_DeleteBankCard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gophkeeper.proto",
}
