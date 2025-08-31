package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/store/schema"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SQL to Proto ---------------------------------------------------------------

func (p *Profile) ToProto() *mantraev1.Profile {
	return &mantraev1.Profile{
		Id:          p.ID,
		Name:        p.Name,
		Description: SafeString(p.Description),
		Token:       p.Token,
		CreatedAt:   SafeTimestamp(p.CreatedAt),
		UpdatedAt:   SafeTimestamp(p.UpdatedAt),
	}
}

func (e *EntryPoint) ToProto() *mantraev1.EntryPoint {
	return &mantraev1.EntryPoint{
		Id:        e.ID,
		ProfileId: e.ProfileID,
		Name:      e.Name,
		Address:   SafeString(e.Address),
		IsDefault: e.IsDefault,
		CreatedAt: SafeTimestamp(e.CreatedAt),
		UpdatedAt: SafeTimestamp(e.UpdatedAt),
	}
}

func (r *HttpRouter) ToProto() *mantraev1.Router {
	return &mantraev1.Router{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		AgentId:   SafeString(r.AgentID),
		Name:      r.Name,
		Config:    MustMarshalStruct(r.Config),
		Enabled:   r.Enabled,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}
}

func (r *TcpRouter) ToProto() *mantraev1.Router {
	return &mantraev1.Router{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		AgentId:   SafeString(r.AgentID),
		Name:      r.Name,
		Config:    MustMarshalStruct(r.Config),
		Enabled:   r.Enabled,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_TCP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}
}

func (r *UdpRouter) ToProto() *mantraev1.Router {
	return &mantraev1.Router{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		AgentId:   SafeString(r.AgentID),
		Name:      r.Name,
		Config:    MustMarshalStruct(r.Config),
		Enabled:   r.Enabled,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_UDP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}
}

func (s *HttpService) ToProto() *mantraev1.Service {
	return &mantraev1.Service{
		Id:        s.ID,
		ProfileId: s.ProfileID,
		AgentId:   SafeString(s.AgentID),
		Name:      s.Name,
		Config:    MustMarshalStruct(s.Config),
		Enabled:   s.Enabled,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP,
		CreatedAt: SafeTimestamp(s.CreatedAt),
		UpdatedAt: SafeTimestamp(s.UpdatedAt),
	}
}

func (s *TcpService) ToProto() *mantraev1.Service {
	return &mantraev1.Service{
		Id:        s.ID,
		ProfileId: s.ProfileID,
		AgentId:   SafeString(s.AgentID),
		Name:      s.Name,
		Config:    MustMarshalStruct(s.Config),
		Enabled:   s.Enabled,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_TCP,
		CreatedAt: SafeTimestamp(s.CreatedAt),
		UpdatedAt: SafeTimestamp(s.UpdatedAt),
	}
}

func (s *UdpService) ToProto() *mantraev1.Service {
	return &mantraev1.Service{
		Id:        s.ID,
		ProfileId: s.ProfileID,
		AgentId:   SafeString(s.AgentID),
		Name:      s.Name,
		Config:    MustMarshalStruct(s.Config),
		Enabled:   s.Enabled,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_UDP,
		CreatedAt: SafeTimestamp(s.CreatedAt),
		UpdatedAt: SafeTimestamp(s.UpdatedAt),
	}
}

func (m *HttpMiddleware) ToProto() *mantraev1.Middleware {
	return &mantraev1.Middleware{
		Id:        m.ID,
		ProfileId: m.ProfileID,
		AgentId:   SafeString(m.AgentID),
		Name:      m.Name,
		Config:    MustMarshalStruct(m.Config),
		Enabled:   m.Enabled,
		IsDefault: m.IsDefault,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP,
		CreatedAt: SafeTimestamp(m.CreatedAt),
		UpdatedAt: SafeTimestamp(m.UpdatedAt),
	}
}

func (m *TcpMiddleware) ToProto() *mantraev1.Middleware {
	return &mantraev1.Middleware{
		Id:        m.ID,
		ProfileId: m.ProfileID,
		AgentId:   SafeString(m.AgentID),
		Name:      m.Name,
		Config:    MustMarshalStruct(m.Config),
		Enabled:   m.Enabled,
		IsDefault: m.IsDefault,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_TCP,
		CreatedAt: SafeTimestamp(m.CreatedAt),
		UpdatedAt: SafeTimestamp(m.UpdatedAt),
	}
}

func (m *HttpServersTransport) ToProto() *mantraev1.ServersTransport {
	return &mantraev1.ServersTransport{
		Id:        m.ID,
		ProfileId: m.ProfileID,
		AgentId:   SafeString(m.AgentID),
		Name:      m.Name,
		Config:    MustMarshalStruct(m.Config),
		Enabled:   m.Enabled,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP,
		CreatedAt: SafeTimestamp(m.CreatedAt),
		UpdatedAt: SafeTimestamp(m.UpdatedAt),
	}
}

func (m *TcpServersTransport) ToProto() *mantraev1.ServersTransport {
	return &mantraev1.ServersTransport{
		Id:        m.ID,
		ProfileId: m.ProfileID,
		AgentId:   SafeString(m.AgentID),
		Name:      m.Name,
		Config:    MustMarshalStruct(m.Config),
		Enabled:   m.Enabled,
		Type:      mantraev1.ProtocolType_PROTOCOL_TYPE_TCP,
		CreatedAt: SafeTimestamp(m.CreatedAt),
		UpdatedAt: SafeTimestamp(m.UpdatedAt),
	}
}

func (u *User) ToProto() *mantraev1.User {
	return &mantraev1.User{
		Id:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Email:     SafeString(u.Email),
		Otp:       SafeString(u.Otp),
		OtpExpiry: SafeTimestamp(u.OtpExpiry),
		LastLogin: SafeTimestamp(u.LastLogin),
		CreatedAt: SafeTimestamp(u.CreatedAt),
		UpdatedAt: SafeTimestamp(u.UpdatedAt),
	}
}

func (a *Agent) ToProto() *mantraev1.Agent {
	containers := make([]*mantraev1.Container, 0)
	if a.Containers != nil {
		raw, ok := a.Containers.([]byte)
		if !ok {
			slog.Error("containers is not []byte", "type", fmt.Sprintf("%T", a.Containers))
		} else {
			if err := json.Unmarshal(raw, &containers); err != nil {
				slog.Error("failed to unmarshal agent containers", "error", err)
			}
		}
	}

	return &mantraev1.Agent{
		Id:         a.ID,
		ProfileId:  a.ProfileID,
		Hostname:   SafeString(a.Hostname),
		PublicIp:   SafeString(a.PublicIp),
		PrivateIp:  SafeString(a.PrivateIp),
		Containers: containers,
		ActiveIp:   SafeString(a.ActiveIp),
		Token:      a.Token,
		CreatedAt:  SafeTimestamp(a.CreatedAt),
		UpdatedAt:  SafeTimestamp(a.UpdatedAt),
	}
}

func (d *DnsProvider) ToProto() *mantraev1.DnsProvider {
	return &mantraev1.DnsProvider{
		Id:   d.ID,
		Name: d.Name,
		Type: mantraev1.DnsProviderType(d.Type),
		Config: &mantraev1.DnsProviderConfig{
			ApiKey:     d.Config.APIKey,
			ApiUrl:     d.Config.APIUrl,
			Ip:         d.Config.IP,
			Proxied:    d.Config.Proxied,
			AutoUpdate: d.Config.AutoUpdate,
			ZoneType:   d.Config.ZoneType,
		},
		IsDefault: d.IsDefault,
		CreatedAt: SafeTimestamp(d.CreatedAt),
		UpdatedAt: SafeTimestamp(d.UpdatedAt),
	}
}

func (t *TraefikInstance) ToProto() *mantraev1.TraefikInstance {
	return &mantraev1.TraefikInstance{
		Id:          t.ID,
		Name:        t.Name,
		Url:         t.Url,
		Tls:         t.Tls,
		EntryPoints: MustMarshalSlice(*t.Entrypoints),
		Overview:    MustMarshalStruct(t.Overview),
		Config:      MustMarshalStruct(t.Config),
		Version:     MustMarshalStruct(t.Version),
		CreatedAt:   SafeTimestamp(t.CreatedAt),
		UpdatedAt:   SafeTimestamp(t.UpdatedAt),
	}
}

func (a *ListAuditLogsRow) ToProto() *mantraev1.AuditLog {
	return &mantraev1.AuditLog{
		Id:          a.ID,
		Event:       a.Event,
		Details:     SafeString(a.Details),
		ProfileId:   SafeInt64(a.ProfileID),
		ProfileName: SafeString(a.ProfileName),
		UserId:      SafeString(a.UserID),
		UserName:    SafeString(a.UserName),
		AgentId:     SafeString(a.AgentID),
		AgentName:   SafeString(a.AgentName),
		CreatedAt:   SafeTimestamp(a.CreatedAt),
	}
}

// Proto to SQL ---------------------------------------------------------------

func (r *HttpRouter) FromProto(proto *mantraev1.Router) error {
	if proto.Type != mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP {
		return errors.New("invalid router type for HTTP router")
	}

	config, err := UnmarshalStruct[schema.HTTPRouter](proto.Config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	r.ID = proto.Id
	r.ProfileID = proto.ProfileId
	r.AgentID = StringPtr(proto.AgentId)
	r.Name = proto.Name
	r.Config = config
	r.Enabled = proto.Enabled
	r.CreatedAt = TimePtr(proto.CreatedAt.AsTime())
	r.UpdatedAt = TimePtr(proto.UpdatedAt.AsTime())
	return nil
}

func (a *Agent) FromProto(pb *mantraev1.Agent) error {
	raw, err := json.Marshal(pb.Containers)
	if err != nil {
		return fmt.Errorf("failed to marshal containers: %w", err)
	}

	a.ID = pb.Id
	a.ProfileID = pb.ProfileId
	a.Hostname = StringPtr(pb.Hostname)
	a.PublicIp = StringPtr(pb.PublicIp)
	a.PrivateIp = StringPtr(pb.PrivateIp)
	a.ActiveIp = StringPtr(pb.ActiveIp)
	a.Token = pb.Token
	a.Containers = raw
	return nil
}

// Common helper --------------------------------------------------------------

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func SafeInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func SafeInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func SafeFloat32(f *float32) float32 {
	if f == nil {
		return 0.0
	}
	return *f
}

func SafeFloat64(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

func SafeTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func TimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// JSON Objects marshalling and unmarshalling helper --------------------------

func UnmarshalStruct[T any](s *structpb.Struct) (*T, error) {
	// Marshal the proto Struct to JSON bytes
	data, err := s.MarshalJSON()
	if err != nil {
		return nil, err
	}

	// Unmarshal into your target struct
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func MarshalStruct[T any](s *T) (*structpb.Struct, error) {
	// Marshal the target struct to JSON bytes
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	// Unmarshal into your proto Struct
	var out structpb.Struct
	if err := out.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return &out, nil
}

func MustMarshalStruct[T any](s *T) *structpb.Struct {
	data, err := json.Marshal(s)
	if err != nil {
		slog.Error("failed to marshal struct", "error", err)
		return nil
	}

	var out structpb.Struct
	if err := out.UnmarshalJSON(data); err != nil {
		slog.Error("failed to unmarshal struct", "error", err)
		return nil
	}
	return &out
}

func MustUnmarshalStruct[T any](s *structpb.Struct) *T {
	data, err := s.MarshalJSON()
	if err != nil {
		slog.Error("failed to marshal struct", "error", err)
		return nil
	}

	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		slog.Error("failed to unmarshal struct", "error", err)
		return nil
	}
	return &out
}

// Slices marshalling and unmarshalling helper --------------------------------

func MustMarshalSlice[T any](s []T) *structpb.ListValue {
	data, err := json.Marshal(s)
	if err != nil {
		slog.Error("failed to marshal slice", "error", err)
		return nil
	}

	var lv structpb.ListValue
	if err := lv.UnmarshalJSON(data); err != nil {
		slog.Error("failed to unmarshal to ListValue", "error", err)
		return nil
	}
	return &lv
}

func MustUnmarshalSlice[T any](lv *structpb.ListValue) []T {
	data, err := lv.MarshalJSON()
	if err != nil {
		slog.Error("failed to marshal ListValue", "error", err)
		return nil
	}

	var out []T
	if err := json.Unmarshal(data, &out); err != nil {
		slog.Error("failed to unmarshal slice", "error", err)
		return nil
	}
	return out
}
