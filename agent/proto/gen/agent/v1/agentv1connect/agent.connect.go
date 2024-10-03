// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: agent/v1/agent.proto

package agentv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// AgentServiceName is the fully-qualified name of the AgentService service.
	AgentServiceName = "agent.v1.AgentService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AgentServiceGetContainerProcedure is the fully-qualified name of the AgentService's GetContainer
	// RPC.
	AgentServiceGetContainerProcedure = "/agent.v1.AgentService/GetContainer"
	// AgentServiceRefreshTokenProcedure is the fully-qualified name of the AgentService's RefreshToken
	// RPC.
	AgentServiceRefreshTokenProcedure = "/agent.v1.AgentService/RefreshToken"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	agentServiceServiceDescriptor            = v1.File_agent_v1_agent_proto.Services().ByName("AgentService")
	agentServiceGetContainerMethodDescriptor = agentServiceServiceDescriptor.Methods().ByName("GetContainer")
	agentServiceRefreshTokenMethodDescriptor = agentServiceServiceDescriptor.Methods().ByName("RefreshToken")
)

// AgentServiceClient is a client for the agent.v1.AgentService service.
type AgentServiceClient interface {
	GetContainer(context.Context, *connect.Request[v1.GetContainerRequest]) (*connect.Response[v1.GetContainerResponse], error)
	RefreshToken(context.Context, *connect.Request[v1.RefreshTokenRequest]) (*connect.Response[v1.RefreshTokenResponse], error)
}

// NewAgentServiceClient constructs a client for the agent.v1.AgentService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAgentServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AgentServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &agentServiceClient{
		getContainer: connect.NewClient[v1.GetContainerRequest, v1.GetContainerResponse](
			httpClient,
			baseURL+AgentServiceGetContainerProcedure,
			connect.WithSchema(agentServiceGetContainerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		refreshToken: connect.NewClient[v1.RefreshTokenRequest, v1.RefreshTokenResponse](
			httpClient,
			baseURL+AgentServiceRefreshTokenProcedure,
			connect.WithSchema(agentServiceRefreshTokenMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// agentServiceClient implements AgentServiceClient.
type agentServiceClient struct {
	getContainer *connect.Client[v1.GetContainerRequest, v1.GetContainerResponse]
	refreshToken *connect.Client[v1.RefreshTokenRequest, v1.RefreshTokenResponse]
}

// GetContainer calls agent.v1.AgentService.GetContainer.
func (c *agentServiceClient) GetContainer(ctx context.Context, req *connect.Request[v1.GetContainerRequest]) (*connect.Response[v1.GetContainerResponse], error) {
	return c.getContainer.CallUnary(ctx, req)
}

// RefreshToken calls agent.v1.AgentService.RefreshToken.
func (c *agentServiceClient) RefreshToken(ctx context.Context, req *connect.Request[v1.RefreshTokenRequest]) (*connect.Response[v1.RefreshTokenResponse], error) {
	return c.refreshToken.CallUnary(ctx, req)
}

// AgentServiceHandler is an implementation of the agent.v1.AgentService service.
type AgentServiceHandler interface {
	GetContainer(context.Context, *connect.Request[v1.GetContainerRequest]) (*connect.Response[v1.GetContainerResponse], error)
	RefreshToken(context.Context, *connect.Request[v1.RefreshTokenRequest]) (*connect.Response[v1.RefreshTokenResponse], error)
}

// NewAgentServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAgentServiceHandler(svc AgentServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	agentServiceGetContainerHandler := connect.NewUnaryHandler(
		AgentServiceGetContainerProcedure,
		svc.GetContainer,
		connect.WithSchema(agentServiceGetContainerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	agentServiceRefreshTokenHandler := connect.NewUnaryHandler(
		AgentServiceRefreshTokenProcedure,
		svc.RefreshToken,
		connect.WithSchema(agentServiceRefreshTokenMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/agent.v1.AgentService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AgentServiceGetContainerProcedure:
			agentServiceGetContainerHandler.ServeHTTP(w, r)
		case AgentServiceRefreshTokenProcedure:
			agentServiceRefreshTokenHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAgentServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAgentServiceHandler struct{}

func (UnimplementedAgentServiceHandler) GetContainer(context.Context, *connect.Request[v1.GetContainerRequest]) (*connect.Response[v1.GetContainerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("agent.v1.AgentService.GetContainer is not implemented"))
}

func (UnimplementedAgentServiceHandler) RefreshToken(context.Context, *connect.Request[v1.RefreshTokenRequest]) (*connect.Response[v1.RefreshTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("agent.v1.AgentService.RefreshToken is not implemented"))
}
