/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protos

import (
	context "context"
	reflect "reflect"
	sync "sync"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	v1 "k8s.io/api/core/v1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BestOptionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options []*Option `protobuf:"bytes,1,rep,name=options,proto3" json:"options,omitempty"`
	// key is node id from options
	NodeMap map[string]*v1.Node `protobuf:"bytes,2,rep,name=nodeMap,proto3" json:"nodeMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *BestOptionsRequest) Reset() {
	*x = BestOptionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BestOptionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BestOptionsRequest) ProtoMessage() {}

func (x *BestOptionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BestOptionsRequest.ProtoReflect.Descriptor instead.
func (*BestOptionsRequest) Descriptor() ([]byte, []int) {
	return file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescGZIP(), []int{0}
}

func (x *BestOptionsRequest) GetOptions() []*Option {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *BestOptionsRequest) GetNodeMap() map[string]*v1.Node {
	if x != nil {
		return x.NodeMap
	}
	return nil
}

type BestOptionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options []*Option `protobuf:"bytes,1,rep,name=options,proto3" json:"options,omitempty"`
}

func (x *BestOptionsResponse) Reset() {
	*x = BestOptionsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BestOptionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BestOptionsResponse) ProtoMessage() {}

func (x *BestOptionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BestOptionsResponse.ProtoReflect.Descriptor instead.
func (*BestOptionsResponse) Descriptor() ([]byte, []int) {
	return file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescGZIP(), []int{1}
}

func (x *BestOptionsResponse) GetOptions() []*Option {
	if x != nil {
		return x.Options
	}
	return nil
}

type Option struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// only need the ID of node to uniquely identify the nodeGroup, used in the nodeInfo map.
	NodeGroupId string    `protobuf:"bytes,1,opt,name=nodeGroupId,proto3" json:"nodeGroupId,omitempty"`
	NodeCount   int32     `protobuf:"varint,2,opt,name=nodeCount,proto3" json:"nodeCount,omitempty"`
	Debug       string    `protobuf:"bytes,3,opt,name=debug,proto3" json:"debug,omitempty"`
	Pod         []*v1.Pod `protobuf:"bytes,4,rep,name=pod,proto3" json:"pod,omitempty"`
}

func (x *Option) Reset() {
	*x = Option{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Option) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Option) ProtoMessage() {}

func (x *Option) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Option.ProtoReflect.Descriptor instead.
func (*Option) Descriptor() ([]byte, []int) {
	return file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescGZIP(), []int{2}
}

func (x *Option) GetNodeGroupId() string {
	if x != nil {
		return x.NodeGroupId
	}
	return ""
}

func (x *Option) GetNodeCount() int32 {
	if x != nil {
		return x.NodeCount
	}
	return 0
}

func (x *Option) GetDebug() string {
	if x != nil {
		return x.Debug
	}
	return ""
}

func (x *Option) GetPod() []*v1.Pod {
	if x != nil {
		return x.Pod
	}
	return nil
}

var File_cluster_autoscaler_expander_grpcplugin_protos_expander_proto protoreflect.FileDescriptor

var file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63,
	0x61, 0x6c, 0x65, 0x72, 0x2f, 0x65, 0x78, 0x70, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f,
	0x65, 0x78, 0x70, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x67, 0x72, 0x70, 0x63, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x1a, 0x22, 0x6b, 0x38, 0x73, 0x2e,
	0x69, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdf,
	0x01, 0x0a, 0x12, 0x42, 0x65, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x45, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x70, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x42, 0x65, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x70, 0x1a, 0x54, 0x0a, 0x0c, 0x4e, 0x6f,
	0x64, 0x65, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2e, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6b, 0x38,
	0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x43, 0x0a, 0x13, 0x42, 0x65, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x89, 0x01, 0x0a, 0x06, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x12, 0x29, 0x0a, 0x03, 0x70, 0x6f, 0x64, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x64, 0x52, 0x03, 0x70, 0x6f,
	0x64, 0x32, 0x5c, 0x0a, 0x08, 0x45, 0x78, 0x70, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x50, 0x0a,
	0x0b, 0x42, 0x65, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1e, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x42, 0x65, 0x73, 0x74, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x42, 0x65, 0x73, 0x74, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x2f, 0x5a, 0x2d, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x61, 0x75, 0x74, 0x6f, 0x73,
	0x63, 0x61, 0x6c, 0x65, 0x72, 0x2f, 0x65, 0x78, 0x70, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescOnce sync.Once
	file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescData = file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDesc
)

func file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescGZIP() []byte {
	file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescOnce.Do(func() {
		file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescData = protoimpl.X.CompressGZIP(file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescData)
	})
	return file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDescData
}

var file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_goTypes = []interface{}{
	(*BestOptionsRequest)(nil),  // 0: grpcplugin.BestOptionsRequest
	(*BestOptionsResponse)(nil), // 1: grpcplugin.BestOptionsResponse
	(*Option)(nil),              // 2: grpcplugin.Option
	nil,                         // 3: grpcplugin.BestOptionsRequest.NodeMapEntry
	(*v1.Pod)(nil),              // 4: k8s.io.api.core.v1.Pod
	(*v1.Node)(nil),             // 5: k8s.io.api.core.v1.Node
}
var file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_depIdxs = []int32{
	2, // 0: grpcplugin.BestOptionsRequest.options:type_name -> grpcplugin.Option
	3, // 1: grpcplugin.BestOptionsRequest.nodeMap:type_name -> grpcplugin.BestOptionsRequest.NodeMapEntry
	2, // 2: grpcplugin.BestOptionsResponse.options:type_name -> grpcplugin.Option
	4, // 3: grpcplugin.Option.pod:type_name -> k8s.io.api.core.v1.Pod
	5, // 4: grpcplugin.BestOptionsRequest.NodeMapEntry.value:type_name -> k8s.io.api.core.v1.Node
	0, // 5: grpcplugin.Expander.BestOptions:input_type -> grpcplugin.BestOptionsRequest
	1, // 6: grpcplugin.Expander.BestOptions:output_type -> grpcplugin.BestOptionsResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_init() }
func file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_init() {
	if File_cluster_autoscaler_expander_grpcplugin_protos_expander_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BestOptionsRequest); i {
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
		file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BestOptionsResponse); i {
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
		file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Option); i {
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
			RawDescriptor: file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_goTypes,
		DependencyIndexes: file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_depIdxs,
		MessageInfos:      file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_msgTypes,
	}.Build()
	File_cluster_autoscaler_expander_grpcplugin_protos_expander_proto = out.File
	file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_rawDesc = nil
	file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_goTypes = nil
	file_cluster_autoscaler_expander_grpcplugin_protos_expander_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ExpanderClient is the client API for Expander service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ExpanderClient interface {
	BestOptions(ctx context.Context, in *BestOptionsRequest, opts ...grpc.CallOption) (*BestOptionsResponse, error)
}

type expanderClient struct {
	cc grpc.ClientConnInterface
}

func NewExpanderClient(cc grpc.ClientConnInterface) ExpanderClient {
	return &expanderClient{cc}
}

func (c *expanderClient) BestOptions(ctx context.Context, in *BestOptionsRequest, opts ...grpc.CallOption) (*BestOptionsResponse, error) {
	out := new(BestOptionsResponse)
	err := c.cc.Invoke(ctx, "/grpcplugin.Expander/BestOptions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExpanderServer is the server API for Expander service.
type ExpanderServer interface {
	BestOptions(context.Context, *BestOptionsRequest) (*BestOptionsResponse, error)
}

// UnimplementedExpanderServer can be embedded to have forward compatible implementations.
type UnimplementedExpanderServer struct {
}

func (*UnimplementedExpanderServer) BestOptions(context.Context, *BestOptionsRequest) (*BestOptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BestOptions not implemented")
}

func RegisterExpanderServer(s *grpc.Server, srv ExpanderServer) {
	s.RegisterService(&_Expander_serviceDesc, srv)
}

func _Expander_BestOptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BestOptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpanderServer).BestOptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcplugin.Expander/BestOptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpanderServer).BestOptions(ctx, req.(*BestOptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Expander_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpcplugin.Expander",
	HandlerType: (*ExpanderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BestOptions",
			Handler:    _Expander_BestOptions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster-autoscaler/expander/grpcplugin/protos/expander.proto",
}
