package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func AuditLogsToProto(logs []db.ListAuditLogsRow) []*mantraev1.AuditLog {
	var auditLogs []*mantraev1.AuditLog
	for _, l := range logs {
		auditLogs = append(auditLogs, &mantraev1.AuditLog{
			Id:          l.ID,
			ProfileId:   SafeInt64(l.ProfileID),
			ProfileName: SafeString(l.ProfileName),
			UserId:      SafeString(l.UserID),
			UserName:    SafeString(l.UserName),
			AgentId:     SafeString(l.AgentID),
			AgentName:   SafeString(l.AgentName),
			Event:       l.Event,
			Details:     SafeString(l.Details),
			CreatedAt:   SafeTimestamp(l.CreatedAt),
		})
	}
	return auditLogs
}
