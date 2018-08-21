// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mnemosynerpc/session.proto

package mnemosynerpc // import "github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import wrappers "github.com/golang/protobuf/ptypes/wrappers"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Session struct {
	AccessToken          string               `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	SubjectId            string               `protobuf:"bytes,2,opt,name=subject_id,json=subjectId,proto3" json:"subject_id,omitempty"`
	SubjectClient        string               `protobuf:"bytes,3,opt,name=subject_client,json=subjectClient,proto3" json:"subject_client,omitempty"`
	Bag                  map[string]string    `protobuf:"bytes,4,rep,name=bag,proto3" json:"bag,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ExpireAt             *timestamp.Timestamp `protobuf:"bytes,5,opt,name=expire_at,json=expireAt,proto3" json:"expire_at,omitempty"`
	RefreshToken         string               `protobuf:"bytes,6,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Session) Reset()         { *m = Session{} }
func (m *Session) String() string { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()    {}
func (*Session) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{0}
}
func (m *Session) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Session.Unmarshal(m, b)
}
func (m *Session) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Session.Marshal(b, m, deterministic)
}
func (dst *Session) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Session.Merge(dst, src)
}
func (m *Session) XXX_Size() int {
	return xxx_messageInfo_Session.Size(m)
}
func (m *Session) XXX_DiscardUnknown() {
	xxx_messageInfo_Session.DiscardUnknown(m)
}

var xxx_messageInfo_Session proto.InternalMessageInfo

func (m *Session) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *Session) GetSubjectId() string {
	if m != nil {
		return m.SubjectId
	}
	return ""
}

func (m *Session) GetSubjectClient() string {
	if m != nil {
		return m.SubjectClient
	}
	return ""
}

func (m *Session) GetBag() map[string]string {
	if m != nil {
		return m.Bag
	}
	return nil
}

func (m *Session) GetExpireAt() *timestamp.Timestamp {
	if m != nil {
		return m.ExpireAt
	}
	return nil
}

func (m *Session) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

type GetRequest struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{1}
}
func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (dst *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(dst, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

type GetResponse struct {
	Session              *Session `protobuf:"bytes,1,opt,name=session,proto3" json:"session,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{2}
}
func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (dst *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(dst, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

func (m *GetResponse) GetSession() *Session {
	if m != nil {
		return m.Session
	}
	return nil
}

type ContextResponse struct {
	Session              *Session `protobuf:"bytes,1,opt,name=session,proto3" json:"session,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContextResponse) Reset()         { *m = ContextResponse{} }
func (m *ContextResponse) String() string { return proto.CompactTextString(m) }
func (*ContextResponse) ProtoMessage()    {}
func (*ContextResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{3}
}
func (m *ContextResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContextResponse.Unmarshal(m, b)
}
func (m *ContextResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContextResponse.Marshal(b, m, deterministic)
}
func (dst *ContextResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContextResponse.Merge(dst, src)
}
func (m *ContextResponse) XXX_Size() int {
	return xxx_messageInfo_ContextResponse.Size(m)
}
func (m *ContextResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ContextResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ContextResponse proto.InternalMessageInfo

func (m *ContextResponse) GetSession() *Session {
	if m != nil {
		return m.Session
	}
	return nil
}

type ListRequest struct {
	// Offset tells how many sessions should be skipped.
	Offset int64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	// Limit tells how many entries shuld be returned.
	// By default it's 10.
	Limit                int64    `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Query                *Query   `protobuf:"bytes,11,opt,name=query,proto3" json:"query,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRequest) Reset()         { *m = ListRequest{} }
func (m *ListRequest) String() string { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()    {}
func (*ListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{4}
}
func (m *ListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRequest.Unmarshal(m, b)
}
func (m *ListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRequest.Marshal(b, m, deterministic)
}
func (dst *ListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRequest.Merge(dst, src)
}
func (m *ListRequest) XXX_Size() int {
	return xxx_messageInfo_ListRequest.Size(m)
}
func (m *ListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRequest proto.InternalMessageInfo

func (m *ListRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *ListRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *ListRequest) GetQuery() *Query {
	if m != nil {
		return m.Query
	}
	return nil
}

type ListResponse struct {
	Sessions             []*Session `protobuf:"bytes,1,rep,name=sessions,proto3" json:"sessions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ListResponse) Reset()         { *m = ListResponse{} }
func (m *ListResponse) String() string { return proto.CompactTextString(m) }
func (*ListResponse) ProtoMessage()    {}
func (*ListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{5}
}
func (m *ListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListResponse.Unmarshal(m, b)
}
func (m *ListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListResponse.Marshal(b, m, deterministic)
}
func (dst *ListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListResponse.Merge(dst, src)
}
func (m *ListResponse) XXX_Size() int {
	return xxx_messageInfo_ListResponse.Size(m)
}
func (m *ListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListResponse proto.InternalMessageInfo

func (m *ListResponse) GetSessions() []*Session {
	if m != nil {
		return m.Sessions
	}
	return nil
}

type Query struct {
	ExpireAtFrom         *timestamp.Timestamp `protobuf:"bytes,1,opt,name=expire_at_from,json=expireAtFrom,proto3" json:"expire_at_from,omitempty"`
	ExpireAtTo           *timestamp.Timestamp `protobuf:"bytes,2,opt,name=expire_at_to,json=expireAtTo,proto3" json:"expire_at_to,omitempty"`
	RefreshToken         string               `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{6}
}
func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (dst *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(dst, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetExpireAtFrom() *timestamp.Timestamp {
	if m != nil {
		return m.ExpireAtFrom
	}
	return nil
}

func (m *Query) GetExpireAtTo() *timestamp.Timestamp {
	if m != nil {
		return m.ExpireAtTo
	}
	return nil
}

func (m *Query) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

type ExistsRequest struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExistsRequest) Reset()         { *m = ExistsRequest{} }
func (m *ExistsRequest) String() string { return proto.CompactTextString(m) }
func (*ExistsRequest) ProtoMessage()    {}
func (*ExistsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{7}
}
func (m *ExistsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExistsRequest.Unmarshal(m, b)
}
func (m *ExistsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExistsRequest.Marshal(b, m, deterministic)
}
func (dst *ExistsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExistsRequest.Merge(dst, src)
}
func (m *ExistsRequest) XXX_Size() int {
	return xxx_messageInfo_ExistsRequest.Size(m)
}
func (m *ExistsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExistsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExistsRequest proto.InternalMessageInfo

func (m *ExistsRequest) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

type StartRequest struct {
	Session              *Session `protobuf:"bytes,1,opt,name=session,proto3" json:"session,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartRequest) Reset()         { *m = StartRequest{} }
func (m *StartRequest) String() string { return proto.CompactTextString(m) }
func (*StartRequest) ProtoMessage()    {}
func (*StartRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{8}
}
func (m *StartRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartRequest.Unmarshal(m, b)
}
func (m *StartRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartRequest.Marshal(b, m, deterministic)
}
func (dst *StartRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartRequest.Merge(dst, src)
}
func (m *StartRequest) XXX_Size() int {
	return xxx_messageInfo_StartRequest.Size(m)
}
func (m *StartRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StartRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StartRequest proto.InternalMessageInfo

func (m *StartRequest) GetSession() *Session {
	if m != nil {
		return m.Session
	}
	return nil
}

type StartResponse struct {
	Session              *Session `protobuf:"bytes,1,opt,name=session,proto3" json:"session,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartResponse) Reset()         { *m = StartResponse{} }
func (m *StartResponse) String() string { return proto.CompactTextString(m) }
func (*StartResponse) ProtoMessage()    {}
func (*StartResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{9}
}
func (m *StartResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartResponse.Unmarshal(m, b)
}
func (m *StartResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartResponse.Marshal(b, m, deterministic)
}
func (dst *StartResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartResponse.Merge(dst, src)
}
func (m *StartResponse) XXX_Size() int {
	return xxx_messageInfo_StartResponse.Size(m)
}
func (m *StartResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StartResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StartResponse proto.InternalMessageInfo

func (m *StartResponse) GetSession() *Session {
	if m != nil {
		return m.Session
	}
	return nil
}

type AbandonRequest struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AbandonRequest) Reset()         { *m = AbandonRequest{} }
func (m *AbandonRequest) String() string { return proto.CompactTextString(m) }
func (*AbandonRequest) ProtoMessage()    {}
func (*AbandonRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{10}
}
func (m *AbandonRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AbandonRequest.Unmarshal(m, b)
}
func (m *AbandonRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AbandonRequest.Marshal(b, m, deterministic)
}
func (dst *AbandonRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AbandonRequest.Merge(dst, src)
}
func (m *AbandonRequest) XXX_Size() int {
	return xxx_messageInfo_AbandonRequest.Size(m)
}
func (m *AbandonRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AbandonRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AbandonRequest proto.InternalMessageInfo

func (m *AbandonRequest) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

type SetValueRequest struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetValueRequest) Reset()         { *m = SetValueRequest{} }
func (m *SetValueRequest) String() string { return proto.CompactTextString(m) }
func (*SetValueRequest) ProtoMessage()    {}
func (*SetValueRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{11}
}
func (m *SetValueRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetValueRequest.Unmarshal(m, b)
}
func (m *SetValueRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetValueRequest.Marshal(b, m, deterministic)
}
func (dst *SetValueRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetValueRequest.Merge(dst, src)
}
func (m *SetValueRequest) XXX_Size() int {
	return xxx_messageInfo_SetValueRequest.Size(m)
}
func (m *SetValueRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetValueRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetValueRequest proto.InternalMessageInfo

func (m *SetValueRequest) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *SetValueRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SetValueRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type SetValueResponse struct {
	Bag                  map[string]string `protobuf:"bytes,1,rep,name=bag,proto3" json:"bag,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SetValueResponse) Reset()         { *m = SetValueResponse{} }
func (m *SetValueResponse) String() string { return proto.CompactTextString(m) }
func (*SetValueResponse) ProtoMessage()    {}
func (*SetValueResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{12}
}
func (m *SetValueResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetValueResponse.Unmarshal(m, b)
}
func (m *SetValueResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetValueResponse.Marshal(b, m, deterministic)
}
func (dst *SetValueResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetValueResponse.Merge(dst, src)
}
func (m *SetValueResponse) XXX_Size() int {
	return xxx_messageInfo_SetValueResponse.Size(m)
}
func (m *SetValueResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetValueResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetValueResponse proto.InternalMessageInfo

func (m *SetValueResponse) GetBag() map[string]string {
	if m != nil {
		return m.Bag
	}
	return nil
}

type DeleteRequest struct {
	AccessToken          string               `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	ExpireAtFrom         *timestamp.Timestamp `protobuf:"bytes,2,opt,name=expire_at_from,json=expireAtFrom,proto3" json:"expire_at_from,omitempty"`
	ExpireAtTo           *timestamp.Timestamp `protobuf:"bytes,3,opt,name=expire_at_to,json=expireAtTo,proto3" json:"expire_at_to,omitempty"`
	RefreshToken         string               `protobuf:"bytes,4,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	SubjectId            string               `protobuf:"bytes,5,opt,name=subject_id,json=subjectId,proto3" json:"subject_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *DeleteRequest) Reset()         { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()    {}
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_session_5b66c0ae3e44b7ea, []int{13}
}
func (m *DeleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteRequest.Unmarshal(m, b)
}
func (m *DeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteRequest.Marshal(b, m, deterministic)
}
func (dst *DeleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteRequest.Merge(dst, src)
}
func (m *DeleteRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteRequest.Size(m)
}
func (m *DeleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteRequest proto.InternalMessageInfo

func (m *DeleteRequest) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *DeleteRequest) GetExpireAtFrom() *timestamp.Timestamp {
	if m != nil {
		return m.ExpireAtFrom
	}
	return nil
}

func (m *DeleteRequest) GetExpireAtTo() *timestamp.Timestamp {
	if m != nil {
		return m.ExpireAtTo
	}
	return nil
}

func (m *DeleteRequest) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

func (m *DeleteRequest) GetSubjectId() string {
	if m != nil {
		return m.SubjectId
	}
	return ""
}

func init() {
	proto.RegisterType((*Session)(nil), "mnemosynerpc.Session")
	proto.RegisterMapType((map[string]string)(nil), "mnemosynerpc.Session.BagEntry")
	proto.RegisterType((*GetRequest)(nil), "mnemosynerpc.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "mnemosynerpc.GetResponse")
	proto.RegisterType((*ContextResponse)(nil), "mnemosynerpc.ContextResponse")
	proto.RegisterType((*ListRequest)(nil), "mnemosynerpc.ListRequest")
	proto.RegisterType((*ListResponse)(nil), "mnemosynerpc.ListResponse")
	proto.RegisterType((*Query)(nil), "mnemosynerpc.Query")
	proto.RegisterType((*ExistsRequest)(nil), "mnemosynerpc.ExistsRequest")
	proto.RegisterType((*StartRequest)(nil), "mnemosynerpc.StartRequest")
	proto.RegisterType((*StartResponse)(nil), "mnemosynerpc.StartResponse")
	proto.RegisterType((*AbandonRequest)(nil), "mnemosynerpc.AbandonRequest")
	proto.RegisterType((*SetValueRequest)(nil), "mnemosynerpc.SetValueRequest")
	proto.RegisterType((*SetValueResponse)(nil), "mnemosynerpc.SetValueResponse")
	proto.RegisterMapType((map[string]string)(nil), "mnemosynerpc.SetValueResponse.BagEntry")
	proto.RegisterType((*DeleteRequest)(nil), "mnemosynerpc.DeleteRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SessionManagerClient is the client API for SessionManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SessionManagerClient interface {
	// Get retrieves session for given access token.
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	// Context works like Get but takes access token from metadata within context.
	// It expects "authorization" key to be present within metadata.
	Context(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ContextResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*wrappers.BoolValue, error)
	Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error)
	Abandon(ctx context.Context, in *AbandonRequest, opts ...grpc.CallOption) (*wrappers.BoolValue, error)
	SetValue(ctx context.Context, in *SetValueRequest, opts ...grpc.CallOption) (*SetValueResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*wrappers.Int64Value, error)
}

type sessionManagerClient struct {
	cc *grpc.ClientConn
}

func NewSessionManagerClient(cc *grpc.ClientConn) SessionManagerClient {
	return &sessionManagerClient{cc}
}

func (c *sessionManagerClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/mnemosynerpc.SessionManager/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) Context(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ContextResponse, error) {
	out := new(ContextResponse)
	err := c.cc.Invoke(ctx, "/mnemosynerpc.SessionManager/Context", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/mnemosynerpc.SessionManager/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*wrappers.BoolValue, error) {
	out := new(wrappers.BoolValue)
	err := c.cc.Invoke(ctx, "/mnemosynerpc.SessionManager/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error) {
	out := new(StartResponse)
	err := c.cc.Invoke(ctx, "/mnemosynerpc.SessionManager/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) Abandon(ctx context.Context, in *AbandonRequest, opts ...grpc.CallOption) (*wrappers.BoolValue, error) {
	out := new(wrappers.BoolValue)
	err := c.cc.Invoke(ctx, "/mnemosynerpc.SessionManager/Abandon", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) SetValue(ctx context.Context, in *SetValueRequest, opts ...grpc.CallOption) (*SetValueResponse, error) {
	out := new(SetValueResponse)
	err := c.cc.Invoke(ctx, "/mnemosynerpc.SessionManager/SetValue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*wrappers.Int64Value, error) {
	out := new(wrappers.Int64Value)
	err := c.cc.Invoke(ctx, "/mnemosynerpc.SessionManager/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SessionManagerServer is the server API for SessionManager service.
type SessionManagerServer interface {
	// Get retrieves session for given access token.
	Get(context.Context, *GetRequest) (*GetResponse, error)
	// Context works like Get but takes access token from metadata within context.
	// It expects "authorization" key to be present within metadata.
	Context(context.Context, *empty.Empty) (*ContextResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	Exists(context.Context, *ExistsRequest) (*wrappers.BoolValue, error)
	Start(context.Context, *StartRequest) (*StartResponse, error)
	Abandon(context.Context, *AbandonRequest) (*wrappers.BoolValue, error)
	SetValue(context.Context, *SetValueRequest) (*SetValueResponse, error)
	Delete(context.Context, *DeleteRequest) (*wrappers.Int64Value, error)
}

func RegisterSessionManagerServer(s *grpc.Server, srv SessionManagerServer) {
	s.RegisterService(&_SessionManager_serviceDesc, srv)
}

func _SessionManager_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemosynerpc.SessionManager/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_Context_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).Context(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemosynerpc.SessionManager/Context",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).Context(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemosynerpc.SessionManager/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemosynerpc.SessionManager/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).Exists(ctx, req.(*ExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemosynerpc.SessionManager/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).Start(ctx, req.(*StartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_Abandon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AbandonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).Abandon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemosynerpc.SessionManager/Abandon",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).Abandon(ctx, req.(*AbandonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_SetValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).SetValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemosynerpc.SessionManager/SetValue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).SetValue(ctx, req.(*SetValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemosynerpc.SessionManager/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SessionManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mnemosynerpc.SessionManager",
	HandlerType: (*SessionManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _SessionManager_Get_Handler,
		},
		{
			MethodName: "Context",
			Handler:    _SessionManager_Context_Handler,
		},
		{
			MethodName: "List",
			Handler:    _SessionManager_List_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _SessionManager_Exists_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _SessionManager_Start_Handler,
		},
		{
			MethodName: "Abandon",
			Handler:    _SessionManager_Abandon_Handler,
		},
		{
			MethodName: "SetValue",
			Handler:    _SessionManager_SetValue_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _SessionManager_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mnemosynerpc/session.proto",
}

func init() { proto.RegisterFile("mnemosynerpc/session.proto", fileDescriptor_session_5b66c0ae3e44b7ea) }

var fileDescriptor_session_5b66c0ae3e44b7ea = []byte{
	// 790 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0x4d, 0x4f, 0xeb, 0x46,
	0x14, 0xc5, 0x71, 0xbe, 0xb8, 0x4e, 0x02, 0x9a, 0xb6, 0xc8, 0x75, 0x0a, 0xa5, 0xae, 0xaa, 0xd2,
	0x8d, 0x4d, 0x43, 0x45, 0x3f, 0x84, 0x0a, 0x04, 0x52, 0x44, 0x3f, 0x16, 0x35, 0xa8, 0x8b, 0xaa,
	0x52, 0xe4, 0x84, 0x49, 0x70, 0x63, 0x7b, 0x8c, 0x67, 0x52, 0x48, 0xf7, 0xef, 0xdf, 0xbc, 0xf5,
	0x5b, 0xbc, 0x1f, 0xf6, 0xd6, 0x4f, 0x9e, 0x19, 0x87, 0xd8, 0x09, 0x3c, 0x02, 0xbb, 0x78, 0xee,
	0x99, 0x7b, 0xef, 0xb9, 0x73, 0xce, 0x0d, 0x18, 0x41, 0x88, 0x03, 0x42, 0x27, 0x21, 0x8e, 0xa3,
	0xbe, 0x4d, 0x31, 0xa5, 0x1e, 0x09, 0xad, 0x28, 0x26, 0x8c, 0xa0, 0xda, 0x6c, 0xcc, 0xf8, 0x7c,
	0x48, 0xc8, 0xd0, 0xc7, 0x36, 0x8f, 0xf5, 0xc6, 0x03, 0x9b, 0x79, 0x01, 0xa6, 0xcc, 0x0d, 0x22,
	0x01, 0x37, 0x9a, 0x79, 0x00, 0x0e, 0x22, 0x36, 0x91, 0xc1, 0xad, 0x7c, 0xf0, 0x36, 0x76, 0xa3,
	0x08, 0xc7, 0x54, 0xc4, 0xcd, 0xb7, 0x05, 0xa8, 0x5c, 0x88, 0xea, 0xe8, 0x0b, 0xa8, 0xb9, 0xfd,
	0x3e, 0xa6, 0xb4, 0xcb, 0xc8, 0x08, 0x87, 0xba, 0xb2, 0xad, 0xec, 0xac, 0x3a, 0x9a, 0x38, 0xbb,
	0x4c, 0x8e, 0xd0, 0x26, 0x00, 0x1d, 0xf7, 0xfe, 0xc5, 0x7d, 0xd6, 0xf5, 0xae, 0xf4, 0x02, 0x07,
	0xac, 0xca, 0x93, 0xf3, 0x2b, 0xf4, 0x15, 0x34, 0xd2, 0x70, 0xdf, 0xf7, 0x70, 0xc8, 0x74, 0x95,
	0x43, 0xea, 0xf2, 0xf4, 0x84, 0x1f, 0xa2, 0x5d, 0x50, 0x7b, 0xee, 0x50, 0x2f, 0x6e, 0xab, 0x3b,
	0x5a, 0x6b, 0xcb, 0x9a, 0xa5, 0x6b, 0xc9, 0x66, 0xac, 0xb6, 0x3b, 0xec, 0x84, 0x2c, 0x9e, 0x38,
	0x09, 0x14, 0x7d, 0x0f, 0xab, 0xf8, 0x2e, 0xf2, 0x62, 0xdc, 0x75, 0x99, 0x5e, 0xda, 0x56, 0x76,
	0xb4, 0x96, 0x61, 0x09, 0x6a, 0x56, 0x4a, 0xcd, 0xba, 0x4c, 0x07, 0xe3, 0x54, 0x05, 0xf8, 0x98,
	0xa1, 0x2f, 0xa1, 0x1e, 0xe3, 0x41, 0x8c, 0xe9, 0xb5, 0x24, 0x55, 0xe6, 0x0d, 0xd5, 0xe4, 0x21,
	0x67, 0x65, 0xec, 0x43, 0x35, 0x2d, 0x87, 0xd6, 0x41, 0x1d, 0xe1, 0x89, 0xe4, 0x9e, 0xfc, 0x44,
	0x1f, 0x43, 0xe9, 0x3f, 0xd7, 0x1f, 0x63, 0x49, 0x57, 0x7c, 0xfc, 0x54, 0xf8, 0x41, 0x31, 0x6d,
	0x80, 0x33, 0xcc, 0x1c, 0x7c, 0x33, 0xc6, 0x94, 0x3d, 0x61, 0x7c, 0xe6, 0xcf, 0xa0, 0xf1, 0x0b,
	0x34, 0x22, 0x21, 0xc5, 0xc8, 0x86, 0x8a, 0x7c, 0x79, 0x0e, 0xd6, 0x5a, 0x9f, 0x2c, 0x9c, 0x85,
	0x93, 0xa2, 0xcc, 0x36, 0xac, 0x9d, 0x90, 0x90, 0xe1, 0xbb, 0x17, 0xe4, 0xf0, 0x41, 0xfb, 0xdd,
	0xa3, 0xd3, 0xae, 0x37, 0xa0, 0x4c, 0x06, 0x03, 0x8a, 0x19, 0xbf, 0xae, 0x3a, 0xf2, 0x2b, 0x61,
	0xed, 0x7b, 0x81, 0xc7, 0x38, 0x6b, 0xd5, 0x11, 0x1f, 0xe8, 0x1b, 0x28, 0xdd, 0x8c, 0x71, 0x3c,
	0xd1, 0x35, 0x5e, 0xeb, 0xa3, 0x6c, 0xad, 0x3f, 0x93, 0x90, 0x23, 0x10, 0xbf, 0x16, 0xab, 0xea,
	0xba, 0x66, 0x1e, 0x43, 0x4d, 0x54, 0x93, 0xed, 0x7e, 0x0b, 0x55, 0xd9, 0x08, 0xd5, 0x15, 0xfe,
	0xfe, 0x0f, 0xf4, 0x3b, 0x85, 0x99, 0xaf, 0x15, 0x28, 0xf1, 0xcc, 0xe8, 0x08, 0x1a, 0x53, 0x15,
	0x74, 0x07, 0x31, 0x09, 0x24, 0xe5, 0xc7, 0xa4, 0x50, 0x4b, 0xa5, 0xf0, 0x4b, 0x4c, 0x02, 0x74,
	0x00, 0xb5, 0xfb, 0x0c, 0x8c, 0x70, 0x72, 0x8f, 0xdf, 0x87, 0xf4, 0xfe, 0x25, 0x99, 0x17, 0x93,
	0x3a, 0x2f, 0x26, 0xb3, 0x05, 0xf5, 0xce, 0x9d, 0x47, 0x19, 0x5d, 0x42, 0x17, 0x87, 0x50, 0xbb,
	0x60, 0x6e, 0x3c, 0x7d, 0x94, 0xa5, 0x1f, 0xf5, 0x08, 0xea, 0x32, 0xc1, 0x73, 0x65, 0xb1, 0x07,
	0x8d, 0xe3, 0x9e, 0x1b, 0x5e, 0x91, 0x70, 0x89, 0xbe, 0xff, 0x81, 0xb5, 0x0b, 0xcc, 0xfe, 0x4a,
	0x0c, 0xf1, 0xf4, 0x5b, 0xa9, 0xc5, 0x0a, 0x0b, 0x2c, 0xa6, 0xce, 0x58, 0xcc, 0x7c, 0xa5, 0xc0,
	0xfa, 0x7d, 0x7a, 0x49, 0xec, 0x47, 0xb1, 0x3b, 0x84, 0x76, 0xbe, 0xce, 0x93, 0xca, 0x82, 0xb3,
	0x4b, 0xe4, 0xd9, 0x36, 0x7f, 0xa7, 0x40, 0xfd, 0x14, 0xfb, 0x98, 0x2d, 0x43, 0x72, 0x5e, 0xab,
	0x85, 0x17, 0x6a, 0x55, 0x7d, 0x99, 0x56, 0x8b, 0xf3, 0x5a, 0xcd, 0xad, 0xf3, 0x52, 0x6e, 0x9d,
	0xb7, 0xde, 0x14, 0xa1, 0x21, 0x85, 0xf2, 0x87, 0x1b, 0xba, 0x43, 0x1c, 0xa3, 0x03, 0x50, 0xcf,
	0x30, 0x43, 0x7a, 0x76, 0xf0, 0xf7, 0x5b, 0xd0, 0xf8, 0x74, 0x41, 0x44, 0xbc, 0x86, 0xb9, 0x82,
	0xda, 0x50, 0x91, 0xfb, 0x0b, 0x6d, 0xcc, 0xf1, 0xe8, 0x24, 0x7f, 0x5b, 0xc6, 0x66, 0xf6, 0x7e,
	0x6e, 0xdd, 0x99, 0x2b, 0xe8, 0x10, 0x8a, 0xc9, 0x46, 0x41, 0xb9, 0x42, 0x33, 0x3b, 0xcd, 0x30,
	0x16, 0x85, 0xa6, 0x09, 0x4e, 0xa0, 0x2c, 0x0c, 0x8a, 0x9a, 0x59, 0x5c, 0xc6, 0xb6, 0xc6, 0xfc,
	0xa0, 0xdb, 0x84, 0xf8, 0x5c, 0x5f, 0x9c, 0x49, 0x89, 0x1b, 0x0e, 0xe5, 0x6a, 0xcd, 0xda, 0xd8,
	0x68, 0x2e, 0x8c, 0x4d, 0x1b, 0xe9, 0x40, 0x45, 0x5a, 0x0e, 0x7d, 0x96, 0x45, 0x66, 0x9d, 0xf8,
	0x81, 0x56, 0x7e, 0x83, 0x6a, 0x2a, 0x7c, 0xb4, 0xf9, 0x90, 0x21, 0x44, 0xa2, 0xad, 0xc7, 0xfd,
	0x62, 0xae, 0xa0, 0x53, 0x28, 0x0b, 0xa9, 0xe7, 0x87, 0x93, 0x31, 0x80, 0xd1, 0x9c, 0xeb, 0xe8,
	0x3c, 0x64, 0xfb, 0xdf, 0xc9, 0x96, 0xda, 0xad, 0xbf, 0x77, 0x87, 0x1e, 0xbb, 0x1e, 0xf7, 0xac,
	0x3e, 0x09, 0xec, 0xc8, 0x23, 0x2c, 0x1e, 0x91, 0x5b, 0xd7, 0xef, 0xff, 0x3f, 0x1e, 0xd9, 0xd3,
	0xb4, 0xf6, 0x6c, 0x81, 0x5e, 0x99, 0xa7, 0xda, 0x7b, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x5e, 0xbf,
	0x19, 0x91, 0x1a, 0x09, 0x00, 0x00,
}
