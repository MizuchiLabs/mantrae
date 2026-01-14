package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
)

// AuditEvent represents an auditable operation
type AuditEvent struct {
	ProfileID *int64
	Event     string
	Details   string
	UserID    *string
	AgentID   *string
}

// NewAuditInterceptor automatically logs CRUD operations
func NewAuditInterceptor(app *config.App) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Execute the actual request
			resp, err := next(ctx, req)

			// Only audit on successful operations
			if err == nil {
				params := &db.CreateAuditLogParams{}
				params.UserID = GetUserIDFromContext(ctx)
				params.AgentID = GetAgentIDFromContext(ctx)
				if auditEvent := extractAuditEvent(req, resp); auditEvent != nil {
					if auditEvent.Details == "" || auditEvent.Event == "" {
						slog.Warn("audit event is missing details or event", "event", auditEvent)
						return resp, err
					}
					// Log audit event asynchronously to avoid blocking the response
					go func() {
						params.ProfileID = auditEvent.ProfileID
						params.Event = auditEvent.Event
						if auditEvent.Details != "" {
							params.Details = &auditEvent.Details
						}
						if err = app.Conn.Query.CreateAuditLog(context.Background(), params); err != nil {
							slog.Error("failed to create audit log", "error", err)
						}
					}()
				}
			}

			return resp, err
		}
	}
}

// extractAuditEvent extracts audit information from request/response
func extractAuditEvent(req connect.AnyRequest, resp connect.AnyResponse) *AuditEvent {
	procedure := req.Spec().Procedure
	parts := strings.Split(procedure, "/")
	if len(parts) < 3 {
		return nil
	}

	service := parts[1]
	method := parts[2]

	eventType := mapMethodToEvent(method)
	if eventType == "" {
		return nil // Not an auditable operation
	}

	// Extract profile ID and details based on service and method
	profileID, details := extractProfileAndDetails(service, method, req, resp)
	return &AuditEvent{
		ProfileID: profileID,
		Event:     fmt.Sprintf("%s.%s", getResourceType(service), eventType),
		Details:   details,
	}
}

// mapMethodToEvent maps gRPC method names to audit event types
func mapMethodToEvent(method string) string {
	switch {
	case strings.HasPrefix(method, "Create"):
		return "create"
	case strings.HasPrefix(method, "Update"):
		return "update"
	case strings.HasPrefix(method, "Delete"):
		return "delete"
	default:
		return ""
	}
}

// getResourceType extracts resource type from service name
func getResourceType(service string) string {
	switch {
	case strings.Contains(service, "ProfileService"):
		return "profile"
	case strings.Contains(service, "RouterService"):
		return "router"
	case strings.Contains(service, "ServiceService"):
		return "service"
	case strings.Contains(service, "MiddlewareService"):
		return "middleware"
	case strings.Contains(service, "EntryPointService"):
		return "entrypoint"
	case strings.Contains(service, "DNSProviderService"):
		return "dns_provider"
	case strings.Contains(service, "AgentService"):
		return "agent"
	case strings.Contains(service, "UserService"):
		return "user"
	default:
		return "unknown"
	}
}

// extractProfileAndDetails extracts profile ID and operation details
func extractProfileAndDetails(
	service, method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch service {
	case "mantrae.v1.ProfileService":
		return extractProfileServiceDetails(method, req, resp)
	case "mantrae.v1.RouterService":
		return extractRouterServiceDetails(method, req, resp)
	case "mantrae.v1.ServiceService":
		return extractServiceServiceDetails(method, req, resp)
	case "mantrae.v1.MiddlewareService":
		return extractMiddlewareServiceDetails(method, req, resp)
	case "mantrae.v1.EntryPointService":
		return extractEntryPointServiceDetails(method, req, resp)
	case "mantrae.v1.DNSProviderService":
		return extractDNSProviderServiceDetails(method, req, resp)
	case "mantrae.v1.AgentService":
		return extractAgentServiceDetails(method, req, resp)
	case "mantrae.v1.UserService":
		return extractUserServiceDetails(method, req, resp)
	default:
		return nil, ""
	}
}

func extractProfileServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch method {
	case "CreateProfile":
		if createReq, ok := req.Any().(*mantraev1.CreateProfileRequest); ok {
			if createResp, ok := resp.Any().(*mantraev1.CreateProfileResponse); ok {
				return &createResp.Profile.Id, fmt.Sprintf(
					"Created profile '%s' (ID: %d)",
					createReq.Name,
					createResp.Profile.Id,
				)
			}
		}
	case "UpdateProfile":
		if updateReq, ok := req.Any().(*mantraev1.UpdateProfileRequest); ok {
			return &updateReq.Id, fmt.Sprintf(
				"Updated profile to name '%s' (ID: %d)",
				updateReq.Name,
				updateReq.Id,
			)
		}
	case "DeleteProfile":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteProfileRequest); ok {
			return &deleteReq.Id, fmt.Sprintf("Deleted profile (ID: %d)", deleteReq.Id)
		}
	}
	return nil, ""
}

func extractRouterServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch method {
	case "CreateRouter":
		if createReq, ok := req.Any().(*mantraev1.CreateRouterRequest); ok {
			return &createReq.ProfileId, fmt.Sprintf(
				"Created router '%s' under profile ID %d",
				createReq.Name,
				createReq.ProfileId,
			)
		}
	case "UpdateRouter":
		if updateReq, ok := req.Any().(*mantraev1.UpdateRouterRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateRouterResponse); ok {
				return &updateResp.Router.ProfileId, fmt.Sprintf(
					"Updated router '%s' (ID: %s)",
					updateReq.Name,
					updateResp.Router.Id,
				)
			}
		}
	case "DeleteRouter":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteRouterRequest); ok {
			return nil, fmt.Sprintf("Deleted router (ID: %s)", deleteReq.Id)
		}
	}
	return nil, ""
}

func extractServiceServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch method {
	case "CreateService":
		if createReq, ok := req.Any().(*mantraev1.CreateServiceRequest); ok {
			return &createReq.ProfileId, fmt.Sprintf(
				"Created service '%s' under profile ID %d",
				createReq.Name,
				createReq.ProfileId,
			)
		}
	case "UpdateService":
		if updateReq, ok := req.Any().(*mantraev1.UpdateServiceRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateServiceResponse); ok {
				return &updateResp.Service.ProfileId, fmt.Sprintf(
					"Updated service '%s' (ID: %s)",
					updateReq.Name,
					updateResp.Service.Id,
				)
			}
		}
	case "DeleteService":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteServiceRequest); ok {
			return nil, fmt.Sprintf("Deleted service (ID: %s)", deleteReq.Id)
		}
	}
	return nil, ""
}

func extractMiddlewareServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch method {
	case "CreateMiddleware":
		if createReq, ok := req.Any().(*mantraev1.CreateMiddlewareRequest); ok {
			return &createReq.ProfileId, fmt.Sprintf(
				"Created middleware '%s' under profile ID %d",
				createReq.Name,
				createReq.ProfileId,
			)
		}
	case "UpdateMiddleware":
		if updateReq, ok := req.Any().(*mantraev1.UpdateMiddlewareRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateMiddlewareResponse); ok {
				return &updateResp.Middleware.ProfileId, fmt.Sprintf(
					"Updated middleware '%s' (ID: %s)",
					updateReq.Name,
					updateResp.Middleware.Id,
				)
			}
		}
	case "DeleteMiddleware":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteMiddlewareRequest); ok {
			return nil, fmt.Sprintf("Deleted middleware (ID: %s)", deleteReq.Id)
		}
	}
	return nil, ""
}

func extractEntryPointServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch method {
	case "CreateEntryPoint":
		if createReq, ok := req.Any().(*mantraev1.CreateEntryPointRequest); ok {
			return &createReq.ProfileId, fmt.Sprintf(
				"Created entrypoint '%s' under profile ID %d",
				createReq.Name,
				createReq.ProfileId,
			)
		}
	case "UpdateEntryPoint":
		if updateReq, ok := req.Any().(*mantraev1.UpdateEntryPointRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateEntryPointResponse); ok {
				return &updateResp.EntryPoint.ProfileId, fmt.Sprintf(
					"Updated entrypoint '%s' (ID: %s)",
					updateReq.Name,
					updateResp.EntryPoint.Id,
				)
			}
		}
	case "DeleteEntryPoint":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteEntryPointRequest); ok {
			return nil, fmt.Sprintf("Deleted entrypoint (ID: %s)", deleteReq.Id)
		}
	}
	return nil, ""
}

func extractDNSProviderServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch method {
	case "CreateDnsProvider":
		if createReq, ok := req.Any().(*mantraev1.CreateDNSProviderRequest); ok {
			return nil, fmt.Sprintf("Created DNS provider '%s'", createReq.Name)
		}
	case "UpdateDnsProvider":
		if updateReq, ok := req.Any().(*mantraev1.UpdateDNSProviderRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateDNSProviderResponse); ok {
				return nil, fmt.Sprintf(
					"Updated DNS provider '%s' (ID: %s)",
					updateReq.Name,
					updateResp.DnsProvider.Id,
				)
			}
		}
	case "DeleteDnsProvider":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteDNSProviderRequest); ok {
			return nil, fmt.Sprintf("Deleted DNS provider (ID: %s)", deleteReq.Id)
		}
	}
	return nil, ""
}

func extractAgentServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch method {
	case "CreateAgent":
		if createReq, ok := req.Any().(*mantraev1.CreateAgentRequest); ok {
			return &createReq.ProfileId, fmt.Sprintf(
				"Added agent under profile ID %d",
				createReq.ProfileId,
			)
		}
	case "UpdateAgent":
		if updateResp, ok := resp.Any().(*mantraev1.UpdateAgentResponse); ok {
			return &updateResp.Agent.ProfileId, fmt.Sprintf(
				"Updated agent (ID: %s)",
				updateResp.Agent.Id,
			)
		}
	case "DeleteAgent":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteAgentRequest); ok {
			return nil, fmt.Sprintf("Deleted agent (ID: %s)", deleteReq.Id)
		}
	}
	return nil, ""
}

func extractUserServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (*int64, string) {
	switch method {
	case "CreateUser":
		if createReq, ok := req.Any().(*mantraev1.CreateUserRequest); ok {
			return nil, fmt.Sprintf("Created user '%s'", createReq.Username)
		}
	case "UpdateUser":
		if updateReq, ok := req.Any().(*mantraev1.UpdateUserRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateUserResponse); ok {
				return nil, fmt.Sprintf(
					"Updated user '%s' (ID: %s)",
					updateReq.Username,
					updateResp.User.Id,
				)
			}
		}
	case "DeleteUser":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteUserRequest); ok {
			return nil, fmt.Sprintf("Deleted user (ID: %s)", deleteReq.Id)
		}
	}
	return nil, ""
}
