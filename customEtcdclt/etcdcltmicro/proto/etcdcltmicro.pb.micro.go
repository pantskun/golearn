// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: etcdcltmicro.proto

package etcdcltmicro

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v3/api"
	client "github.com/micro/go-micro/v3/client"
	server "github.com/micro/go-micro/v3/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for EtcdcltMicro service

func NewEtcdcltMicroEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for EtcdcltMicro service

type EtcdcltMicroService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type etcdcltMicroService struct {
	c    client.Client
	name string
}

func NewEtcdcltMicroService(name string, c client.Client) EtcdcltMicroService {
	return &etcdcltMicroService{
		c:    c,
		name: name,
	}
}

func (c *etcdcltMicroService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "EtcdcltMicro.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EtcdcltMicro service

type EtcdcltMicroHandler interface {
	Call(context.Context, *Request, *Response) error
}

func RegisterEtcdcltMicroHandler(s server.Server, hdlr EtcdcltMicroHandler, opts ...server.HandlerOption) error {
	type etcdcltMicro interface {
		Call(ctx context.Context, in *Request, out *Response) error
	}
	type EtcdcltMicro struct {
		etcdcltMicro
	}
	h := &etcdcltMicroHandler{hdlr}
	return s.Handle(s.NewHandler(&EtcdcltMicro{h}, opts...))
}

type etcdcltMicroHandler struct {
	EtcdcltMicroHandler
}

func (h *etcdcltMicroHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.EtcdcltMicroHandler.Call(ctx, in, out)
}