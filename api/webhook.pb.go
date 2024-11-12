// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.3
// source: webhook.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type WebhookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uri     string           `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
	Method  string           `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Headers *structpb.Struct `protobuf:"bytes,3,opt,name=headers,proto3" json:"headers,omitempty"`
	Body    *structpb.Struct `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *WebhookRequest) Reset() {
	*x = WebhookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookRequest) ProtoMessage() {}

func (x *WebhookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookRequest.ProtoReflect.Descriptor instead.
func (*WebhookRequest) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{0}
}

func (x *WebhookRequest) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

func (x *WebhookRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *WebhookRequest) GetHeaders() *structpb.Struct {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *WebhookRequest) GetBody() *structpb.Struct {
	if x != nil {
		return x.Body
	}
	return nil
}

type WebhookResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  uint32           `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Headers *structpb.Struct `protobuf:"bytes,2,opt,name=headers,proto3" json:"headers,omitempty"`
	Body    *structpb.Struct `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *WebhookResponse) Reset() {
	*x = WebhookResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookResponse) ProtoMessage() {}

func (x *WebhookResponse) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookResponse.ProtoReflect.Descriptor instead.
func (*WebhookResponse) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{1}
}

func (x *WebhookResponse) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *WebhookResponse) GetHeaders() *structpb.Struct {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *WebhookResponse) GetBody() *structpb.Struct {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_webhook_proto protoreflect.FileDescriptor

var file_webhook_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x61, 0x70, 0x69, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x0e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x69, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x69, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12,
	0x31, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x12, 0x2b, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22,
	0x89, 0x01, 0x0a, 0x0f, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x31, 0x0a, 0x07, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x2b,
	0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x32, 0x46, 0x0a, 0x0e, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x34, 0x0a,
	0x07, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x13, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_webhook_proto_rawDescOnce sync.Once
	file_webhook_proto_rawDescData = file_webhook_proto_rawDesc
)

func file_webhook_proto_rawDescGZIP() []byte {
	file_webhook_proto_rawDescOnce.Do(func() {
		file_webhook_proto_rawDescData = protoimpl.X.CompressGZIP(file_webhook_proto_rawDescData)
	})
	return file_webhook_proto_rawDescData
}

var file_webhook_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_webhook_proto_goTypes = []interface{}{
	(*WebhookRequest)(nil),  // 0: api.WebhookRequest
	(*WebhookResponse)(nil), // 1: api.WebhookResponse
	(*structpb.Struct)(nil), // 2: google.protobuf.Struct
}
var file_webhook_proto_depIdxs = []int32{
	2, // 0: api.WebhookRequest.headers:type_name -> google.protobuf.Struct
	2, // 1: api.WebhookRequest.body:type_name -> google.protobuf.Struct
	2, // 2: api.WebhookResponse.headers:type_name -> google.protobuf.Struct
	2, // 3: api.WebhookResponse.body:type_name -> google.protobuf.Struct
	0, // 4: api.WebhookService.Webhook:input_type -> api.WebhookRequest
	1, // 5: api.WebhookService.Webhook:output_type -> api.WebhookResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_webhook_proto_init() }
func file_webhook_proto_init() {
	if File_webhook_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_webhook_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookRequest); i {
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
		file_webhook_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookResponse); i {
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
			RawDescriptor: file_webhook_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_webhook_proto_goTypes,
		DependencyIndexes: file_webhook_proto_depIdxs,
		MessageInfos:      file_webhook_proto_msgTypes,
	}.Build()
	File_webhook_proto = out.File
	file_webhook_proto_rawDesc = nil
	file_webhook_proto_goTypes = nil
	file_webhook_proto_depIdxs = nil
}
