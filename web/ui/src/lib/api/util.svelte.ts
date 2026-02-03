import { useMutation, useQuery } from '$lib/query';
import { UtilService } from '$lib/gen/mantrae/v1/util_pb';
import { AuditLogService } from '$lib/gen/mantrae/v1/auditlog_pb';
import { profileID } from '$lib/store.svelte';
import { BackupService } from '$lib/gen/mantrae/v1/backup_pb';
import { toast } from 'svelte-sonner';

export const util = {
	// Queries
	ip: () => useQuery(UtilService.method.getPublicIP, {}, {}),
	version: () => useQuery(UtilService.method.getVersion, {}, { select: (res) => res.version }),
	config: () =>
		useQuery(
			UtilService.method.getDynamicConfig,
			{ profileId: profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.config }
		)
};

export const audit = {
	// Queries
	logs: (limit?: bigint, offset?: bigint) =>
		useQuery(
			AuditLogService.method.listAuditLogs,
			{ limit, offset },
			{ select: (res) => res.auditLogs }
		)
};

export const backup = {
	// Queries
	list: () => useQuery(BackupService.method.listBackups, {}, { select: (res) => res.backups }),

	// Mutations
	create: () =>
		useMutation(BackupService.method.createBackup, {
			onSuccess: () => toast.success('Backup created!')
		}),
	delete: () =>
		useMutation(BackupService.method.deleteBackup, {
			onSuccess: () => toast.success('Backup deleted!')
		}),
	restore: () =>
		useMutation(BackupService.method.restoreBackup, {
			onSuccess: () => toast.success('Backup restored!'),
			transform: (variables) => ({ ...variables, profileId: profileID.current })
		})
};
