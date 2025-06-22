import type { DescService } from '@bufbuild/protobuf';
import { createClient, type Client } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { goto } from '$app/navigation';
import { token } from './stores/common';
import { user } from './stores/user';
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

// Global state variables
export const BACKEND_PORT = import.meta.env.PORT || 3000;
export const BASE_URL = import.meta.env.PROD ? '' : `http://127.0.0.1:${BACKEND_PORT}`;

export function useClient<T extends DescService>(
	service: T,
	fetch?: typeof window.fetch
): Client<T> {
	// Custom fetch function that adds the Authorization header
	const customFetch: typeof window.fetch = async (url, options) => {
		const headers = new Headers(options?.headers); // Get existing headers
		if (token.value) {
			headers.set('Authorization', 'Bearer ' + token.value); // Add the Authorization header
		}
		const customOptions = {
			...options,
			headers
		};
		return fetch ? fetch(url, customOptions) : window.fetch(url, customOptions); // Use custom fetch or default
	};

	const transport = createConnectTransport({ baseUrl: BASE_URL, fetch: customFetch });
	return createClient(service, transport);
}

export function logout() {
	token.value = null;
	user.clear();
	goto('/login');
}

export async function uploadBackup(event: Event) {
	const input = event.target as HTMLInputElement;
	if (!input.files?.length) return;

	const body = new FormData();
	body.append('file', input.files[0]);

	const headers = new Headers();
	headers.set('Authorization', 'Bearer ' + token.value);

	const response = await fetch(`${BASE_URL}/api/backup`, { method: 'POST', headers, body });
	if (!response.ok) {
		throw new Error('Failed to upload backup');
	}
	toast.success('Backup uploaded successfully');
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
export const utilClient = useClient(UtilService);
