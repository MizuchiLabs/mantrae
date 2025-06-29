package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

// AuditEvent represents an auditable operation
type AuditEvent struct {
	ProfileID int64
	Event     string
	Details   string
}

// AuditInterceptor automatically logs CRUD operations
func NewAuditInterceptor(app *config.App) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Execute the actual request
			resp, err := next(ctx, req)

			// Only audit on successful operations
			if err == nil {
				if auditEvent := extractAuditEvent(req, resp); auditEvent != nil {
					// Log audit event asynchronously to avoid blocking the response
					go func() {
						if auditErr := createAuditLog(context.Background(), app.Conn.GetQuery(), *auditEvent); auditErr != nil {
							slog.Error("failed to create audit log", "error", auditErr)
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

	service := parts[1] // e.g., "mantrae.v1.ProfileService"
	method := parts[2]  // e.g., "CreateProfile"

	// Map method names to audit events
	eventType := mapMethodToEvent(method)
	if eventType == "" {
		return nil // Not an auditable operation
	}

	// Extract profile ID and details based on service and method
	profileID, details := extractProfileAndDetails(service, method, req, resp)
	if profileID == 0 {
		return nil
	}

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
		return "" // Not auditable
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
	case strings.Contains(service, "DnsProviderService"):
		return "dns_provider"
	case strings.Contains(service, "AgentService"):
		return "agent"
	case strings.Contains(service, "UserService"):
		return "user"
	case strings.Contains(service, "SettingService"):
		return "setting"
	default:
		return "unknown"
	}
}

// extractProfileAndDetails extracts profile ID and operation details
func extractProfileAndDetails(
	service, method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (int64, string) {
	switch service {
	case "mantrae.v1.ProfileService":
		return extractProfileServiceDetails(method, req, resp)
	case "mantrae.v1.RouterService":
		return extractRouterServiceDetails(method, req, resp)
	case "mantrae.v1.MiddlewareService":
		return extractMiddlewareServiceDetails(method, req, resp)
	case "mantrae.v1.ServiceService":
		return extractServiceServiceDetails(method, req, resp)
	// Add other services as needed
	default:
		return 0, ""
	}
}

func extractProfileServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (int64, string) {
	switch method {
	case "CreateProfile":
		if createReq, ok := req.Any().(*mantraev1.CreateProfileRequest); ok {
			if createResp, ok := resp.Any().(*mantraev1.CreateProfileResponse); ok {
				return createResp.Profile.Id, fmt.Sprintf(
					"Profile created with name %s",
					createReq.Name,
				)
			}
		}
	case "UpdateProfile":
		if updateReq, ok := req.Any().(*mantraev1.UpdateProfileRequest); ok {
			return updateReq.Id, fmt.Sprintf("Profile updated with name %s", updateReq.Name)
		}
	case "DeleteProfile":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteProfileRequest); ok {
			return deleteReq.Id, fmt.Sprintf("Profile deleted with ID %d", deleteReq.Id)
		}
	}
	return 0, ""
}

func extractRouterServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (int64, string) {
	switch method {
	case "CreateRouter":
		if createReq, ok := req.Any().(*mantraev1.CreateRouterRequest); ok {
			return createReq.ProfileId, fmt.Sprintf("Router created with name %s", createReq.Name)
		}
	case "UpdateRouter":
		if updateReq, ok := req.Any().(*mantraev1.UpdateRouterRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateRouterResponse); ok {
				return updateResp.Router.ProfileId, fmt.Sprintf(
					"Router updated with name %s",
					updateReq.Name,
				)
			}
		}
	case "DeleteRouter":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteRouterRequest); ok {
			// For delete, we need to get profile ID from the response or make a pre-delete query
			// This is a limitation - we might need the profile ID from the request context
			return 0, fmt.Sprintf("Router deleted with ID %d", deleteReq.Id)
		}
	}
	return 0, ""
}

func extractServiceServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (int64, string) {
	switch method {
	case "CreateService":
		if createReq, ok := req.Any().(*mantraev1.CreateServiceRequest); ok {
			return createReq.ProfileId, fmt.Sprintf(
				"Service created with name %s",
				createReq.Name,
			)
		}
	case "UpdateService":
		if updateReq, ok := req.Any().(*mantraev1.UpdateServiceRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateServiceResponse); ok {
				return updateResp.Service.ProfileId, fmt.Sprintf(
					"Service updated with name %s",
					updateReq.Name,
				)
			}
		}
	case "DeleteService":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteServiceRequest); ok {
			// For delete, we need to get profile ID from the response or make a pre-delete query
			// This is a limitation - we might need the profile ID from the request context
			return 0, fmt.Sprintf("Service deleted with ID %d", deleteReq.Id)
		}
	}
	return 0, ""
}

func extractMiddlewareServiceDetails(
	method string,
	req connect.AnyRequest,
	resp connect.AnyResponse,
) (int64, string) {
	switch method {
	case "CreateMiddleware":
		if createReq, ok := req.Any().(*mantraev1.CreateMiddlewareRequest); ok {
			return createReq.ProfileId, fmt.Sprintf(
				"Middleware created with name %s",
				createReq.Name,
			)
		}
	case "UpdateMiddleware":
		if updateReq, ok := req.Any().(*mantraev1.UpdateMiddlewareRequest); ok {
			if updateResp, ok := resp.Any().(*mantraev1.UpdateMiddlewareResponse); ok {
				return updateResp.Middleware.ProfileId, fmt.Sprintf(
					"Middleware updated with name %s",
					updateReq.Name,
				)
			}
		}
	case "DeleteMiddleware":
		if deleteReq, ok := req.Any().(*mantraev1.DeleteMiddlewareRequest); ok {
			return 0, fmt.Sprintf("Middleware deleted with ID %d", deleteReq.Id)
		}
	}
	return 0, ""
}

// createAuditLog creates an audit log entry
func createAuditLog(ctx context.Context, q *db.Queries, event AuditEvent) error {
	var params db.CreateAuditLogParams
	params.ProfileID = event.ProfileID
	params.Event = event.Event
	if event.Details != "" {
		params.Details = &event.Details
	}

	// Extract user/agent context if available
	if valUserID := ctx.Value(AuthUserIDKey); valUserID != nil {
		if userID, ok := valUserID.(string); ok && userID != "" {
			params.UserID = &userID
		}
	}
	if valAgentID := ctx.Value(AuthAgentIDKey); valAgentID != nil {
		if agentID, ok := valAgentID.(string); ok && agentID != "" {
			params.AgentID = &agentID
		}
	}

	return q.CreateAuditLog(ctx, params)
}
