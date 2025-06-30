import type { DescService } from '@bufbuild/protobuf';
import { createClient, type Client } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { ProfileService } from './gen/mantrae/v1/profile_pb';
import { UserService } from './gen/mantrae/v1/user_pb';
import { RouterService } from './gen/mantrae/v1/router_pb';
import { ServiceService } from './gen/mantrae/v1/service_pb';
import { MiddlewareService } from './gen/mantrae/v1/middleware_pb';
import { SettingService } from './gen/mantrae/v1/setting_pb';
import { BackupService } from './gen/mantrae/v1/backup_pb';
import { EntryPointService } from './gen/mantrae/v1/entry_point_pb';
import { DnsProviderService } from './gen/mantrae/v1/dns_provider_pb';
import { UtilService } from './gen/mantrae/v1/util_pb';
import { AgentService } from './gen/mantrae/v1/agent_pb';
import { toast } from 'svelte-sonner';
import { AuditLogService } from './gen/mantrae/v1/auditlog_pb';
import { profile } from './stores/profile';

// Global state variables
export const BACKEND_PORT = import.meta.env.PORT || 3000;
export const BASE_URL = import.meta.env.PROD ? '/' : `http://127.0.0.1:${BACKEND_PORT}`;

export function useClient<T extends DescService>(
	service: T,
	customFetch?: typeof fetch
): Client<T> {
	const wrappedFetch: typeof fetch = (input, init = {}) => {
		return (customFetch || fetch)(input, {
			...init,
			headers: new Headers(init.headers || {}),
			credentials: 'include'
		});
	};

	const transport = createConnectTransport({
		baseUrl: BASE_URL,
		fetch: wrappedFetch
	});
	return createClient(service, transport);
}

export function handleOIDCLogin() {
	window.location.href = `${BASE_URL}/oidc/login`;
}

export async function upload(input: HTMLInputElement | null, endpoint: string) {
	if (!input?.files?.length) return;

	const body = new FormData();
	body.append('file', input.files[0]);

	const response = await fetch(`${BASE_URL}/upload/${endpoint}`, {
		method: 'POST',
		body,
		credentials: 'include'
	});
	if (!response.ok) {
		throw new Error('Failed to upload');
	}
	toast.success('Uploaded successfully');
}

export async function getConfig(format: string) {
	if (!profile.isValid()) return '';

	const headers = new Headers();
	if (format === 'yaml') {
		headers.set('Accept', 'application/x-yaml');
	}

	try {
		const response = await fetch(`${BASE_URL}/api/${profile.name}`, {
			headers
		});
		if (!response.ok) return '';

		return await response.text();
	} catch (err) {
		const e = err as Error;
		toast.error('Failed to fetch config', { description: e.message });
	}
	return '';
}

// Clients
export const profileClient = useClient(ProfileService);
export const userClient = useClient(UserService);
export const entryPointClient = useClient(EntryPointService);
export const dnsClient = useClient(DnsProviderService);
export const agentClient = useClient(AgentService);
export const routerClient = useClient(RouterService);
export const serviceClient = useClient(ServiceService);
export const middlewareClient = useClient(MiddlewareService);
export const settingClient = useClient(SettingService);
export const backupClient = useClient(BackupService);
export const auditLogClient = useClient(AuditLogService);
export const utilClient = useClient(UtilService);
