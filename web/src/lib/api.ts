import type { DescService } from "@bufbuild/protobuf";
import { createClient, type Client } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ProfileService } from "./gen/mantrae/v1/profile_pb";
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
import { EventService } from "./gen/mantrae/v1/event_pb";

// Global state variables
export const BACKEND_PORT = import.meta.env.PORT || 3000;
export const BASE_URL = import.meta.env.PROD
	? ""
	: `http://127.0.0.1:${BACKEND_PORT}`;

export function useClient<T extends DescService>(
	service: T,
	customFetch?: typeof fetch,
): Client<T> {
	const headers = new Headers();
	headers.set("Content-Type", "application/json");

	// Wrap the fetch to always append headers & credentials
	const wrappedFetch: typeof fetch = (input, init = {}) => {
		const newHeaders = new Headers(init.headers || {});
		headers.forEach((value, key) => {
			newHeaders.set(key, value);
		});

		return (customFetch || fetch)(input, {
			...init,
			headers: newHeaders,
			credentials: "include",
		});
	};

	const transport = createConnectTransport({
		baseUrl: BASE_URL,
		fetch: wrappedFetch,
	});
	return createClient(service, transport);
}

export function handleOIDCLogin() {
	window.location.href = `${BASE_URL}/oidc/login`;
}

export async function upload(input: HTMLInputElement | null, endpoint: string) {
	if (!input?.files?.length) return;

	const body = new FormData();
	body.append("file", input.files[0]);

	const headers = new Headers();
	headers.set("Content-Type", "multipart/form-data");

	const response = await fetch(`${BASE_URL}/api/${endpoint}`, {
		method: "POST",
		headers,
		body,
		credentials: "include",
	});
	if (!response.ok) {
		throw new Error("Failed to upload");
	}
	toast.success("Uploaded successfully");
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
export const eventClient = useClient(EventService);
export const auditLogClient = useClient(AuditLogService);
