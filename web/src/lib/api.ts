import type { DescService } from "@bufbuild/protobuf";
import { createClient, type Client } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ProfileService, type Profile } from "./gen/mantrae/v1/profile_pb";
import { UserService } from "./gen/mantrae/v1/user_pb";
import { RouterService } from "./gen/mantrae/v1/router_pb";
import { ServiceService } from "./gen/mantrae/v1/service_pb";
import { MiddlewareService } from "./gen/mantrae/v1/middleware_pb";
import { SettingService } from "./gen/mantrae/v1/setting_pb";
import { BackupService } from "./gen/mantrae/v1/backup_pb";
import { EntryPointService } from "./gen/mantrae/v1/entry_point_pb";
import { DnsProviderService } from "./gen/mantrae/v1/dns_provider_pb";
import { UtilService } from "./gen/mantrae/v1/util_pb";
import { AgentService } from "./gen/mantrae/v1/agent_pb";
import { toast } from "svelte-sonner";
import { AuditLogService } from "./gen/mantrae/v1/auditlog_pb";
import { profile } from "./stores/profile";
import { baseURL } from "./stores/common";
import { ServersTransportService } from "./gen/mantrae/v1/servers_transport_pb";

export function useClient<T extends DescService>(
	service: T,
	customFetch?: typeof fetch,
): Client<T> {
	const wrappedFetch: typeof fetch = (input, init = {}) => {
		return (customFetch || fetch)(input, {
			...init,
			headers: new Headers(init.headers || {}),
			credentials: "include",
		});
	};

	if (!baseURL.value) throw new Error("Base URL not set");

	const transport = createConnectTransport({
		baseUrl: baseURL.value,
		fetch: wrappedFetch,
	});
	return createClient(service, transport);
}

// Basic health check function
export async function checkHealth(
	customFetch?: typeof fetch,
): Promise<boolean> {
	try {
		if (!baseURL.value) throw new Error("Base URL not set");
		if (!customFetch) customFetch = fetch;
		const res = await customFetch(`${baseURL.value}/healthz`, {
			method: "GET",
		});
		return res.ok;
	} catch {
		return false;
	}
}

export function handleOIDCLogin() {
	window.location.href = `${baseURL.value}/oidc/login`;
}

export async function upload(input: HTMLInputElement | null, endpoint: string) {
	if (!input?.files?.length) return;

	const body = new FormData();
	body.append("file", input.files[0]);

	const response = await fetch(`${baseURL.value}/upload/${endpoint}`, {
		method: "POST",
		body,
		credentials: "include",
	});
	if (!response.ok) {
		throw new Error("Failed to upload");
	}
	toast.success("Uploaded successfully");
}

// Get dynamic traefik config
export async function getConfig(format: string) {
	if (!profile.isValid() || !profile.token) return "";

	const headers = new Headers();
	// headers.set("Mantrae-Traefik-Token", profile.token);
	if (format === "yaml") {
		headers.set("Accept", "application/x-yaml");
	}

	try {
		const response = await fetch(
			`${baseURL.value}/api/${profile.name}?token=${profile.token}`,
			{
				headers,
			},
		);
		if (!response.ok) return "";

		return await response.text();
	} catch (err) {
		const e = err as Error;
		toast.error("Failed to fetch config", { description: e.message });
	}
	return "";
}

// Build traefik connection string
export async function buildConnectionString(p: Profile) {
	const item = p ?? profile?.value;
	if (!item) return "";
	const serverUrl = await settingClient.getSetting({ key: "server_url" });
	if (!serverUrl.value) return "";

	return `${serverUrl.value}/api/${item.name}?token=${item.token}`;
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
export const serversTransportClient = useClient(ServersTransportService);
export const settingClient = useClient(SettingService);
export const backupClient = useClient(BackupService);
export const auditLogClient = useClient(AuditLogService);
export const utilClient = useClient(UtilService);
