// protoc --go_out=plugins=grpc:. *.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.11.4
// source: time.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type NowReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NowReq) Reset() {
	*x = NowReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_time_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NowReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NowReq) ProtoMessage() {}

func (x *NowReq) ProtoReflect() protoreflect.Message {
	mi := &file_time_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NowReq.ProtoReflect.Descriptor instead.
func (*NowReq) Descriptor() ([]byte, []int) {
	return file_time_proto_rawDescGZIP(), []int{0}
}

type NowRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time string `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *NowRes) Reset() {
	*x = NowRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_time_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NowRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NowRes) ProtoMessage() {}

func (x *NowRes) ProtoReflect() protoreflect.Message {
	mi := &file_time_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NowRes.ProtoReflect.Descriptor instead.
func (*NowRes) Descriptor() ([]byte, []int) {
	return file_time_proto_rawDescGZIP(), []int{1}
}

func (x *NowRes) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

var File_time_proto protoreflect.FileDescriptor

var file_time_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x08, 0x0a, 0x06, 0x4e, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x22, 0x1c, 0x0a,
	0x06, 0x4e, 0x6f, 0x77, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x32, 0x2d, 0x0a, 0x04, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x03, 0x4e, 0x6f, 0x77, 0x12, 0x0d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4e, 0x6f, 0x77, 0x52, 0x65, 0x73, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_time_proto_rawDescOnce sync.Once
	file_time_proto_rawDescData = file_time_proto_rawDesc
)

func file_time_proto_rawDescGZIP() []byte {
	file_time_proto_rawDescOnce.Do(func() {
		file_time_proto_rawDescData = protoimpl.X.CompressGZIP(file_time_proto_rawDescData)
	})
	return file_time_proto_rawDescData
}

var file_time_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_time_proto_goTypes = []interface{}{
	(*NowReq)(nil), // 0: proto.NowReq
	(*NowRes)(nil), // 1: proto.NowRes
}
var file_time_proto_depIdxs = []int32{
	0, // 0: proto.Time.Now:input_type -> proto.NowReq
	1, // 1: proto.Time.Now:output_type -> proto.NowRes
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_time_proto_init() }
func file_time_proto_init() {
	if File_time_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_time_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NowReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_time_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NowRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_time_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_time_proto_goTypes,
		DependencyIndexes: file_time_proto_depIdxs,
		MessageInfos:      file_time_proto_msgTypes,
	}.Build()
	File_time_proto = out.File
	file_time_proto_rawDesc = nil
	file_time_proto_goTypes = nil
	file_time_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TimeClient is the client API for Time service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TimeClient interface {
	Now(ctx context.Context, in *NowReq, opts ...grpc.CallOption) (*NowRes, error)
}

type timeClient struct {
	cc grpc.ClientConnInterface
}

func NewTimeClient(cc grpc.ClientConnInterface) TimeClient {
	return &timeClient{cc}
}

func (c *timeClient) Now(ctx context.Context, in *NowReq, opts ...grpc.CallOption) (*NowRes, error) {
	out := new(NowRes)
	err := c.cc.Invoke(ctx, "/proto.Time/Now", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TimeServer is the server API for Time service.
type TimeServer interface {
	Now(context.Context, *NowReq) (*NowRes, error)
}

// UnimplementedTimeServer can be embedded to have forward compatible implementations.
type UnimplementedTimeServer struct {
}

func (*UnimplementedTimeServer) Now(context.Context, *NowReq) (*NowRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Now not implemented")
}

func RegisterTimeServer(s *grpc.Server, srv TimeServer) {
	s.RegisterService(&_Time_serviceDesc, srv)
}

func _Time_Now_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NowReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimeServer).Now(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Time/Now",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimeServer).Now(ctx, req.(*NowReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Time_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Time",
	HandlerType: (*TimeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Now",
			Handler:    _Time_Now_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "time.proto",
}