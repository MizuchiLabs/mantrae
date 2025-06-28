package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func AuditLogsToProto(logs []db.AuditLog) []*mantraev1.AuditLog {
	var auditLogs []*mantraev1.AuditLog
	for _, l := range logs {
		auditLogs = append(auditLogs, &mantraev1.AuditLog{
			Id:        l.ID,
			ProfileId: l.ProfileID,
			UserId:    SafeString(l.UserID),
			AgentId:   SafeString(l.AgentID),
			Event:     l.Event,
			Details:   SafeString(l.Details),
			CreatedAt: SafeTimestamp(l.CreatedAt),
		})
	}
	return auditLogs
}
