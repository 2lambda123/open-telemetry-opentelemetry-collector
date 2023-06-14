// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: opentelemetry/proto/collector/profiles/v1/profiles_service.proto

package v1

import (
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	v1 "go.opentelemetry.io/collector/pdata/internal/data/protogen/profiles/v1"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ExportProfilesServiceRequest struct {
	// An array of ResourceProfiles.
	// For data coming from a single resource this array will typically contain one
	// element. Intermediary nodes (such as OpenTelemetry Collector) that receive
	// data from multiple origins typically batch the data before forwarding further and
	// in that case this array will contain multiple elements.
	ResourceProfiles []*v1.ResourceProfiles `protobuf:"bytes,1,rep,name=resource_profiles,json=resourceProfiles,proto3" json:"resource_profiles,omitempty"`
}

func (m *ExportProfilesServiceRequest) Reset()         { *m = ExportProfilesServiceRequest{} }
func (m *ExportProfilesServiceRequest) String() string { return proto.CompactTextString(m) }
func (*ExportProfilesServiceRequest) ProtoMessage()    {}
func (*ExportProfilesServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0139d7c125155ea1, []int{0}
}
func (m *ExportProfilesServiceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExportProfilesServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExportProfilesServiceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExportProfilesServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportProfilesServiceRequest.Merge(m, src)
}
func (m *ExportProfilesServiceRequest) XXX_Size() int {
	return m.Size()
}
func (m *ExportProfilesServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportProfilesServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExportProfilesServiceRequest proto.InternalMessageInfo

func (m *ExportProfilesServiceRequest) GetResourceProfiles() []*v1.ResourceProfiles {
	if m != nil {
		return m.ResourceProfiles
	}
	return nil
}

type ExportProfilesServiceResponse struct {
	// The details of a partially successful export request.
	//
	// If the request is only partially accepted
	// (i.e. when the server accepts only parts of the data and rejects the rest)
	// the server MUST initialize the `partial_success` field and MUST
	// set the `rejected_<signal>` with the number of items it rejected.
	//
	// Servers MAY also make use of the `partial_success` field to convey
	// warnings/suggestions to senders even when the request was fully accepted.
	// In such cases, the `rejected_<signal>` MUST have a value of `0` and
	// the `error_message` MUST be non-empty.
	//
	// A `partial_success` message with an empty value (rejected_<signal> = 0 and
	// `error_message` = "") is equivalent to it not being set/present. Senders
	// SHOULD interpret it the same way as in the full success case.
	PartialSuccess ExportProfilesPartialSuccess `protobuf:"bytes,1,opt,name=partial_success,json=partialSuccess,proto3" json:"partial_success"`
}

func (m *ExportProfilesServiceResponse) Reset()         { *m = ExportProfilesServiceResponse{} }
func (m *ExportProfilesServiceResponse) String() string { return proto.CompactTextString(m) }
func (*ExportProfilesServiceResponse) ProtoMessage()    {}
func (*ExportProfilesServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0139d7c125155ea1, []int{1}
}
func (m *ExportProfilesServiceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExportProfilesServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExportProfilesServiceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExportProfilesServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportProfilesServiceResponse.Merge(m, src)
}
func (m *ExportProfilesServiceResponse) XXX_Size() int {
	return m.Size()
}
func (m *ExportProfilesServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportProfilesServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExportProfilesServiceResponse proto.InternalMessageInfo

func (m *ExportProfilesServiceResponse) GetPartialSuccess() ExportProfilesPartialSuccess {
	if m != nil {
		return m.PartialSuccess
	}
	return ExportProfilesPartialSuccess{}
}

type ExportProfilesPartialSuccess struct {
	// The number of rejected profile records.
	//
	// A `rejected_<signal>` field holding a `0` value indicates that the
	// request was fully accepted.
	RejectedProfileRecords int64 `protobuf:"varint,1,opt,name=rejected_profile_records,json=rejectedProfileRecords,proto3" json:"rejected_profile_records,omitempty"`
	// A developer-facing human-readable message in English. It should be used
	// either to explain why the server rejected parts of the data during a partial
	// success or to convey warnings/suggestions during a full success. The message
	// should offer guidance on how users can address such issues.
	//
	// error_message is an optional field. An error_message with an empty value
	// is equivalent to it not being set.
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
}

func (m *ExportProfilesPartialSuccess) Reset()         { *m = ExportProfilesPartialSuccess{} }
func (m *ExportProfilesPartialSuccess) String() string { return proto.CompactTextString(m) }
func (*ExportProfilesPartialSuccess) ProtoMessage()    {}
func (*ExportProfilesPartialSuccess) Descriptor() ([]byte, []int) {
	return fileDescriptor_0139d7c125155ea1, []int{2}
}
func (m *ExportProfilesPartialSuccess) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExportProfilesPartialSuccess) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExportProfilesPartialSuccess.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExportProfilesPartialSuccess) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportProfilesPartialSuccess.Merge(m, src)
}
func (m *ExportProfilesPartialSuccess) XXX_Size() int {
	return m.Size()
}
func (m *ExportProfilesPartialSuccess) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportProfilesPartialSuccess.DiscardUnknown(m)
}

var xxx_messageInfo_ExportProfilesPartialSuccess proto.InternalMessageInfo

func (m *ExportProfilesPartialSuccess) GetRejectedProfileRecords() int64 {
	if m != nil {
		return m.RejectedProfileRecords
	}
	return 0
}

func (m *ExportProfilesPartialSuccess) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

func init() {
	proto.RegisterType((*ExportProfilesServiceRequest)(nil), "opentelemetry.proto.collector.profiles.v1.ExportProfilesServiceRequest")
	proto.RegisterType((*ExportProfilesServiceResponse)(nil), "opentelemetry.proto.collector.profiles.v1.ExportProfilesServiceResponse")
	proto.RegisterType((*ExportProfilesPartialSuccess)(nil), "opentelemetry.proto.collector.profiles.v1.ExportProfilesPartialSuccess")
}

func init() {
	proto.RegisterFile("opentelemetry/proto/collector/profiles/v1/profiles_service.proto", fileDescriptor_0139d7c125155ea1)
}

var fileDescriptor_0139d7c125155ea1 = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xc1, 0xca, 0xd3, 0x40,
	0x14, 0x85, 0x33, 0xff, 0x2f, 0x05, 0xa7, 0x6a, 0x35, 0x14, 0x29, 0x45, 0x63, 0x89, 0x9b, 0x16,
	0x64, 0x42, 0xea, 0xc6, 0xa5, 0x54, 0x44, 0x37, 0x62, 0x48, 0xc5, 0x85, 0x0b, 0x43, 0x9c, 0x5e,
	0x43, 0x24, 0xcd, 0xc4, 0x3b, 0xd3, 0xa0, 0x0b, 0x7d, 0x06, 0x77, 0xae, 0x7c, 0x01, 0x77, 0xbe,
	0x45, 0x97, 0x5d, 0xba, 0x10, 0x91, 0xf6, 0x45, 0x24, 0x99, 0x26, 0x98, 0x10, 0xa1, 0xe8, 0x6e,
	0xe6, 0xcc, 0xbd, 0xdf, 0x39, 0x73, 0x87, 0xa1, 0xf7, 0x45, 0x06, 0xa9, 0x82, 0x04, 0xd6, 0xa0,
	0xf0, 0xbd, 0x93, 0xa1, 0x50, 0xc2, 0xe1, 0x22, 0x49, 0x80, 0x2b, 0x81, 0xc5, 0xfe, 0x75, 0x9c,
	0x80, 0x74, 0x72, 0xb7, 0x5e, 0x07, 0x12, 0x30, 0x8f, 0x39, 0xb0, 0xb2, 0xd8, 0x9c, 0x35, 0x08,
	0x5a, 0x64, 0x35, 0x81, 0x55, 0x5d, 0x2c, 0x77, 0xc7, 0xc3, 0x48, 0x44, 0x42, 0x5b, 0x14, 0x2b,
	0x5d, 0x3b, 0x66, 0x5d, 0x11, 0xba, 0x8c, 0x75, 0xbd, 0xfd, 0x91, 0xde, 0x78, 0xf8, 0x2e, 0x13,
	0xa8, 0xbc, 0xa3, 0xbe, 0xd4, 0x79, 0x7c, 0x78, 0xbb, 0x01, 0xa9, 0xcc, 0x97, 0xf4, 0x1a, 0x82,
	0x14, 0x1b, 0xe4, 0x10, 0x54, 0xad, 0x23, 0x32, 0x39, 0x9f, 0xf6, 0xe7, 0x2e, 0xeb, 0x0a, 0xfb,
	0x47, 0x44, 0xe6, 0x1f, 0x3b, 0x2b, 0xb6, 0x7f, 0x15, 0x5b, 0x8a, 0xfd, 0x99, 0xd0, 0x9b, 0x7f,
	0x09, 0x20, 0x33, 0x91, 0x4a, 0x30, 0x73, 0x3a, 0xc8, 0x42, 0x54, 0x71, 0x98, 0x04, 0x72, 0xc3,
	0x39, 0xc8, 0xc2, 0x9f, 0x4c, 0xfb, 0xf3, 0x47, 0xec, 0xe4, 0x61, 0xb1, 0xa6, 0x85, 0xa7, 0x79,
	0x4b, 0x8d, 0x5b, 0x5c, 0xd8, 0xfe, 0xbc, 0x65, 0xf8, 0x57, 0xb2, 0x86, 0x6a, 0x7f, 0x68, 0x4f,
	0xa6, 0xd9, 0x65, 0xde, 0xa3, 0x23, 0x84, 0x37, 0xc0, 0x15, 0xac, 0xaa, 0xc9, 0x04, 0x08, 0x5c,
	0xe0, 0x4a, 0x07, 0x3c, 0xf7, 0xaf, 0x57, 0xe7, 0x47, 0x82, 0xaf, 0x4f, 0xcd, 0xdb, 0xf4, 0x32,
	0x20, 0x0a, 0x0c, 0xd6, 0x20, 0x65, 0x18, 0xc1, 0xe8, 0x6c, 0x42, 0xa6, 0x17, 0xfd, 0x4b, 0xa5,
	0xf8, 0x44, 0x6b, 0xf3, 0x6f, 0x84, 0x0e, 0x5a, 0x23, 0x31, 0xbf, 0x10, 0xda, 0xd3, 0x99, 0xcc,
	0x7f, 0xbf, 0x7c, 0xf3, 0x81, 0xc7, 0x8f, 0xff, 0x1f, 0xa4, 0x1f, 0xca, 0x36, 0x16, 0x3f, 0xc8,
	0x76, 0x6f, 0x91, 0xdd, 0xde, 0x22, 0xbf, 0xf6, 0x16, 0xf9, 0x74, 0xb0, 0x8c, 0xdd, 0xc1, 0x32,
	0xbe, 0x1f, 0x2c, 0x83, 0xde, 0x89, 0xc5, 0xe9, 0x46, 0x8b, 0x61, 0xcb, 0xc3, 0x2b, 0x6a, 0x3d,
	0xf2, 0xc2, 0x8b, 0xda, 0x94, 0xb8, 0xf1, 0xc1, 0x56, 0xa1, 0x0a, 0x9d, 0x38, 0x55, 0x80, 0x69,
	0x98, 0x38, 0xe5, 0xae, 0xb4, 0x89, 0x20, 0xed, 0xfe, 0x87, 0x5f, 0xcf, 0x66, 0x4f, 0x33, 0x48,
	0x9f, 0xd5, 0xbc, 0xd2, 0x89, 0x3d, 0xa8, 0x53, 0x55, 0x41, 0xd8, 0x73, 0xf7, 0x55, 0xaf, 0x64,
	0xdd, 0xfd, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x51, 0xba, 0x0e, 0xbc, 0xe7, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProfilesServiceClient is the client API for ProfilesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProfilesServiceClient interface {
	// For performance reasons, it is recommended to keep this RPC
	// alive for the entire life of the application.
	Export(ctx context.Context, in *ExportProfilesServiceRequest, opts ...grpc.CallOption) (*ExportProfilesServiceResponse, error)
}

type profilesServiceClient struct {
	cc *grpc.ClientConn
}

func NewProfilesServiceClient(cc *grpc.ClientConn) ProfilesServiceClient {
	return &profilesServiceClient{cc}
}

func (c *profilesServiceClient) Export(ctx context.Context, in *ExportProfilesServiceRequest, opts ...grpc.CallOption) (*ExportProfilesServiceResponse, error) {
	out := new(ExportProfilesServiceResponse)
	err := c.cc.Invoke(ctx, "/opentelemetry.proto.collector.profiles.v1.ProfilesService/Export", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfilesServiceServer is the server API for ProfilesService service.
type ProfilesServiceServer interface {
	// For performance reasons, it is recommended to keep this RPC
	// alive for the entire life of the application.
	Export(context.Context, *ExportProfilesServiceRequest) (*ExportProfilesServiceResponse, error)
}

// UnimplementedProfilesServiceServer can be embedded to have forward compatible implementations.
type UnimplementedProfilesServiceServer struct {
}

func (*UnimplementedProfilesServiceServer) Export(ctx context.Context, req *ExportProfilesServiceRequest) (*ExportProfilesServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Export not implemented")
}

func RegisterProfilesServiceServer(s *grpc.Server, srv ProfilesServiceServer) {
	s.RegisterService(&_ProfilesService_serviceDesc, srv)
}

func _ProfilesService_Export_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportProfilesServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfilesServiceServer).Export(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/opentelemetry.proto.collector.profiles.v1.ProfilesService/Export",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfilesServiceServer).Export(ctx, req.(*ExportProfilesServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProfilesService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "opentelemetry.proto.collector.profiles.v1.ProfilesService",
	HandlerType: (*ProfilesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Export",
			Handler:    _ProfilesService_Export_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "opentelemetry/proto/collector/profiles/v1/profiles_service.proto",
}

func (m *ExportProfilesServiceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExportProfilesServiceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExportProfilesServiceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ResourceProfiles) > 0 {
		for iNdEx := len(m.ResourceProfiles) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ResourceProfiles[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProfilesService(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ExportProfilesServiceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExportProfilesServiceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExportProfilesServiceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.PartialSuccess.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintProfilesService(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *ExportProfilesPartialSuccess) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExportProfilesPartialSuccess) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExportProfilesPartialSuccess) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ErrorMessage) > 0 {
		i -= len(m.ErrorMessage)
		copy(dAtA[i:], m.ErrorMessage)
		i = encodeVarintProfilesService(dAtA, i, uint64(len(m.ErrorMessage)))
		i--
		dAtA[i] = 0x12
	}
	if m.RejectedProfileRecords != 0 {
		i = encodeVarintProfilesService(dAtA, i, uint64(m.RejectedProfileRecords))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintProfilesService(dAtA []byte, offset int, v uint64) int {
	offset -= sovProfilesService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ExportProfilesServiceRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ResourceProfiles) > 0 {
		for _, e := range m.ResourceProfiles {
			l = e.Size()
			n += 1 + l + sovProfilesService(uint64(l))
		}
	}
	return n
}

func (m *ExportProfilesServiceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.PartialSuccess.Size()
	n += 1 + l + sovProfilesService(uint64(l))
	return n
}

func (m *ExportProfilesPartialSuccess) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RejectedProfileRecords != 0 {
		n += 1 + sovProfilesService(uint64(m.RejectedProfileRecords))
	}
	l = len(m.ErrorMessage)
	if l > 0 {
		n += 1 + l + sovProfilesService(uint64(l))
	}
	return n
}

func sovProfilesService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProfilesService(x uint64) (n int) {
	return sovProfilesService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ExportProfilesServiceRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProfilesService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExportProfilesServiceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExportProfilesServiceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResourceProfiles", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProfilesService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProfilesService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProfilesService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ResourceProfiles = append(m.ResourceProfiles, &v1.ResourceProfiles{})
			if err := m.ResourceProfiles[len(m.ResourceProfiles)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProfilesService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProfilesService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ExportProfilesServiceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProfilesService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExportProfilesServiceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExportProfilesServiceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PartialSuccess", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProfilesService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProfilesService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProfilesService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PartialSuccess.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProfilesService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProfilesService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ExportProfilesPartialSuccess) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProfilesService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExportProfilesPartialSuccess: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExportProfilesPartialSuccess: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectedProfileRecords", wireType)
			}
			m.RejectedProfileRecords = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProfilesService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RejectedProfileRecords |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorMessage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProfilesService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProfilesService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProfilesService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ErrorMessage = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProfilesService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProfilesService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProfilesService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProfilesService
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProfilesService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProfilesService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthProfilesService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProfilesService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProfilesService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProfilesService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProfilesService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProfilesService = fmt.Errorf("proto: unexpected end of group")
)
