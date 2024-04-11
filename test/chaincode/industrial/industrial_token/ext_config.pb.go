// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: ext_config.proto

package industrialtoken

import (
	proto "github.com/anoideaopen/foundation/proto"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type ExtConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name             string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Decimals         uint32        `protobuf:"varint,2,opt,name=decimals,proto3" json:"decimals,omitempty"`
	UnderlyingAsset  string        `protobuf:"bytes,3,opt,name=underlying_asset,json=underlyingAsset,proto3" json:"underlying_asset,omitempty"`
	DeliveryForm     string        `protobuf:"bytes,4,opt,name=delivery_form,json=deliveryForm,proto3" json:"delivery_form,omitempty"`
	UnitOfMeasure    string        `protobuf:"bytes,5,opt,name=unit_of_measure,json=unitOfMeasure,proto3" json:"unit_of_measure,omitempty"`
	TokensForUnit    string        `protobuf:"bytes,6,opt,name=tokens_for_unit,json=tokensForUnit,proto3" json:"tokens_for_unit,omitempty"`
	PaymentTerms     string        `protobuf:"bytes,7,opt,name=payment_terms,json=paymentTerms,proto3" json:"payment_terms,omitempty"`
	Price            string        `protobuf:"bytes,8,opt,name=price,proto3" json:"price,omitempty"`
	Issuer           *proto.Wallet `protobuf:"bytes,9,opt,name=issuer,proto3" json:"issuer,omitempty"`
	FeeSetter        *proto.Wallet `protobuf:"bytes,10,opt,name=fee_setter,json=feeSetter,proto3" json:"fee_setter,omitempty"`
	FeeAddressSetter *proto.Wallet `protobuf:"bytes,11,opt,name=fee_address_setter,json=feeAddressSetter,proto3" json:"fee_address_setter,omitempty"`
}

func (x *ExtConfig) Reset() {
	*x = ExtConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ext_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExtConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtConfig) ProtoMessage() {}

func (x *ExtConfig) ProtoReflect() protoreflect.Message {
	mi := &file_ext_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtConfig.ProtoReflect.Descriptor instead.
func (*ExtConfig) Descriptor() ([]byte, []int) {
	return file_ext_config_proto_rawDescGZIP(), []int{0}
}

func (x *ExtConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ExtConfig) GetDecimals() uint32 {
	if x != nil {
		return x.Decimals
	}
	return 0
}

func (x *ExtConfig) GetUnderlyingAsset() string {
	if x != nil {
		return x.UnderlyingAsset
	}
	return ""
}

func (x *ExtConfig) GetDeliveryForm() string {
	if x != nil {
		return x.DeliveryForm
	}
	return ""
}

func (x *ExtConfig) GetUnitOfMeasure() string {
	if x != nil {
		return x.UnitOfMeasure
	}
	return ""
}

func (x *ExtConfig) GetTokensForUnit() string {
	if x != nil {
		return x.TokensForUnit
	}
	return ""
}

func (x *ExtConfig) GetPaymentTerms() string {
	if x != nil {
		return x.PaymentTerms
	}
	return ""
}

func (x *ExtConfig) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *ExtConfig) GetIssuer() *proto.Wallet {
	if x != nil {
		return x.Issuer
	}
	return nil
}

func (x *ExtConfig) GetFeeSetter() *proto.Wallet {
	if x != nil {
		return x.FeeSetter
	}
	return nil
}

func (x *ExtConfig) GetFeeAddressSetter() *proto.Wallet {
	if x != nil {
		return x.FeeAddressSetter
	}
	return nil
}

var File_ext_config_proto protoreflect.FileDescriptor

var file_ext_config_proto_rawDesc = []byte{
	0x0a, 0x10, 0x65, 0x78, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0f, 0x69, 0x6e, 0x64, 0x75, 0x73, 0x74, 0x72, 0x69, 0x61, 0x6c, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x1a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd8, 0x03, 0x0a, 0x09, 0x45,
	0x78, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x08, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x2a, 0x02, 0x20, 0x00,
	0x52, 0x08, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x75, 0x6e,
	0x64, 0x65, 0x72, 0x6c, 0x79, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x75, 0x6e, 0x64, 0x65, 0x72, 0x6c, 0x79, 0x69, 0x6e, 0x67,
	0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x79, 0x5f, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65,
	0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x46, 0x6f, 0x72, 0x6d, 0x12, 0x26, 0x0a, 0x0f, 0x75, 0x6e,
	0x69, 0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x75, 0x6e, 0x69, 0x74, 0x4f, 0x66, 0x4d, 0x65, 0x61, 0x73, 0x75,
	0x72, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x5f, 0x66, 0x6f, 0x72,
	0x5f, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x73, 0x46, 0x6f, 0x72, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x65, 0x72, 0x6d, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06,
	0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x12, 0x36, 0x0a, 0x0a, 0x66, 0x65, 0x65, 0x5f, 0x73, 0x65,
	0x74, 0x74, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01,
	0x02, 0x10, 0x01, 0x52, 0x09, 0x66, 0x65, 0x65, 0x53, 0x65, 0x74, 0x74, 0x65, 0x72, 0x12, 0x45,
	0x0a, 0x12, 0x66, 0x65, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x73, 0x65,
	0x74, 0x74, 0x65, 0x72, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01,
	0x02, 0x10, 0x01, 0x52, 0x10, 0x66, 0x65, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x53,
	0x65, 0x74, 0x74, 0x65, 0x72, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x6f, 0x69, 0x64, 0x65, 0x61, 0x6f, 0x70, 0x65, 0x6e, 0x2f,
	0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x69, 0x6e, 0x64, 0x75, 0x73, 0x74,
	0x72, 0x69, 0x61, 0x6c, 0x2f, 0x69, 0x6e, 0x64, 0x75, 0x73, 0x74, 0x72, 0x69, 0x61, 0x6c, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ext_config_proto_rawDescOnce sync.Once
	file_ext_config_proto_rawDescData = file_ext_config_proto_rawDesc
)

func file_ext_config_proto_rawDescGZIP() []byte {
	file_ext_config_proto_rawDescOnce.Do(func() {
		file_ext_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_ext_config_proto_rawDescData)
	})
	return file_ext_config_proto_rawDescData
}

var file_ext_config_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ext_config_proto_goTypes = []interface{}{
	(*ExtConfig)(nil),    // 0: industrialtoken.ExtConfig
	(*proto.Wallet)(nil), // 1: proto.Wallet
}
var file_ext_config_proto_depIdxs = []int32{
	1, // 0: industrialtoken.ExtConfig.issuer:type_name -> proto.Wallet
	1, // 1: industrialtoken.ExtConfig.fee_setter:type_name -> proto.Wallet
	1, // 2: industrialtoken.ExtConfig.fee_address_setter:type_name -> proto.Wallet
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_ext_config_proto_init() }
func file_ext_config_proto_init() {
	if File_ext_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ext_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExtConfig); i {
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
			RawDescriptor: file_ext_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ext_config_proto_goTypes,
		DependencyIndexes: file_ext_config_proto_depIdxs,
		MessageInfos:      file_ext_config_proto_msgTypes,
	}.Build()
	File_ext_config_proto = out.File
	file_ext_config_proto_rawDesc = nil
	file_ext_config_proto_goTypes = nil
	file_ext_config_proto_depIdxs = nil
}
