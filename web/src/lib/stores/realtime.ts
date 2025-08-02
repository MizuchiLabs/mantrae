import {
	EventAction,
	type EventStreamResponse,
} from "$lib/gen/mantrae/v1/util_pb";
import {
	utilClient,
	entryPointClient,
	middlewareClient,
	dnsClient,
	profileClient,
	userClient,
	agentClient,
	routerClient,
	serviceClient,
	serversTransportClient,
	traefikClient,
} from "$lib/api";
import { profile } from "./profile";
import { writable } from "svelte/store";
import type { ConnectError } from "@connectrpc/connect";
import type { Profile } from "$lib/gen/mantrae/v1/profile_pb";
import type { User } from "$lib/gen/mantrae/v1/user_pb";
import type { Agent } from "$lib/gen/mantrae/v1/agent_pb";
import type { EntryPoint } from "$lib/gen/mantrae/v1/entry_point_pb";
import type { Router } from "$lib/gen/mantrae/v1/router_pb";
import type { Service } from "$lib/gen/mantrae/v1/service_pb";
import type { Middleware } from "$lib/gen/mantrae/v1/middleware_pb";
import type { ServersTransport } from "$lib/gen/mantrae/v1/servers_transport_pb";
import type { DnsProvider } from "$lib/gen/mantrae/v1/dns_provider_pb";
import type { TraefikInstance } from "$lib/gen/mantrae/v1/traefik_instance_pb";

export const profiles = writable<Profile[]>([]);
export const users = writable<User[]>([]);
export const agents = writable<Agent[]>([]);
export const entryPoints = writable<EntryPoint[]>([]);
export const routers = writable<Router[]>([]);
export const services = writable<Service[]>([]);
export const middlewares = writable<Middleware[]>([]);
export const serversTransports = writable<ServersTransport[]>([]);
export const dnsProviders = writable<DnsProvider[]>([]);
export const traefikInstances = writable<TraefikInstance[]>([]);

let currentStream: AbortController | null = null;

export async function subscribe() {
	if (!profile.isValid()) return;
	await preload();

	// Cleanup previous connection
	if (currentStream) {
		currentStream.abort();
	}

	currentStream = new AbortController();
	try {
		const stream = utilClient.eventStream(
			{ profileId: profile.id },
			{ signal: currentStream.signal },
		);

		for await (const event of stream) {
			handleEvent(event);
		}
	} catch (error) {
		const err = error as ConnectError;
		if (err.message.includes("canceled") || err.message.includes("aborted")) {
			return;
		}
		console.error("Event stream error:", err.message);
		// Retry connection after delay
		setTimeout(subscribe, 5000);
	}
}

export function unsubscribe() {
	if (currentStream) {
		currentStream.abort();
		currentStream = null;
	}
}

function handleEvent(event: EventStreamResponse) {
	const action = event.action;
	switch (action) {
		case EventAction.CREATED:
			switch (event.data.case) {
				case "profile": {
					const data = event.data.value;
					if (!data) return;
					profiles.update((p) => p.concat(data));
					break;
				}
				case "user": {
					const data = event.data.value;
					if (!data) return;
					users.update((u) => u.concat(data));
					break;
				}
				case "agent": {
					const data = event.data.value;
					if (!data) return;
					agents.update((a) => a.concat(data));
					break;
				}
				case "entryPoint": {
					const data = event.data.value;
					if (!data) return;
					entryPoints.update((e) => e.concat(data));
					break;
				}
				case "router": {
					const data = event.data.value;
					if (!data) return;
					routers.update((r) => r.concat(data));
					break;
				}
				case "service": {
					const data = event.data.value;
					if (!data) return;
					services.update((s) => s.concat(data));
					break;
				}
				case "middleware": {
					const data = event.data.value;
					if (!data) return;
					middlewares.update((m) => m.concat(data));
					break;
				}
				case "serversTransport": {
					const data = event.data.value;
					if (!data) return;
					serversTransports.update((s) => s.concat(data));
					break;
				}
				case "dnsProvider": {
					const data = event.data.value;
					if (!data) return;
					dnsProviders.update((d) => d.concat(data));
					break;
				}
				case "traefikInstance": {
					const data = event.data.value;
					if (!data) return;
					traefikInstances.update((t) => t.concat(data));
					break;
				}
			}
			break;
		case EventAction.UPDATED:
			switch (event.data.case) {
				case "profile": {
					const data = event.data.value;
					if (!data) return;
					profiles.update((p) => p.map((p) => (p.id === data.id ? data : p)));
					break;
				}
				case "user": {
					const data = event.data.value;
					if (!data) return;
					users.update((u) => u.map((u) => (u.id === data.id ? data : u)));
					break;
				}
				case "agent": {
					const data = event.data.value;
					if (!data) return;
					agents.update((a) => a.map((a) => (a.id === data.id ? data : a)));
					break;
				}
				case "entryPoint": {
					// Refetch since we might update multiple entry points
					entryPointClient
						.listEntryPoints({ profileId: profile.id })
						.then((response) => {
							entryPoints.set(response.entryPoints);
						});
					break;
				}
				case "router": {
					const data = event.data.value;
					if (!data) return;
					routers.update((r) => r.map((r) => (r.id === data.id ? data : r)));
					break;
				}
				case "service": {
					const data = event.data.value;
					if (!data) return;
					services.update((s) => s.map((s) => (s.id === data.id ? data : s)));
					break;
				}
				case "middleware": {
					middlewareClient
						.listMiddlewares({ profileId: profile.id })
						.then((response) => {
							middlewares.set(response.middlewares);
						});
					break;
				}
				case "serversTransport": {
					const data = event.data.value;
					if (!data) return;
					serversTransports.update((s) =>
						s.map((s) => (s.id === data.id ? data : s)),
					);
					break;
				}
				case "dnsProvider": {
					dnsClient.listDnsProviders({}).then((response) => {
						dnsProviders.set(response.dnsProviders);
					});
					break;
				}
				case "traefikInstance": {
					const data = event.data.value;
					if (!data) return;
					traefikInstances.update((t) =>
						t.map((t) => (t.id === data.id ? data : t)),
					);
					break;
				}
			}
			break;
		case EventAction.DELETED:
			switch (event.data.case) {
				case "profile": {
					const data = event.data.value;
					if (!data) return;
					profiles.update((p) => p.filter((p) => p.id !== data.id));
					break;
				}
				case "user": {
					const data = event.data.value;
					if (!data) return;
					users.update((u) => u.filter((u) => u.id !== data.id));
					break;
				}
				case "agent": {
					const data = event.data.value;
					if (!data) return;
					agents.update((a) => a.filter((a) => a.id !== data.id));
					break;
				}
				case "entryPoint": {
					const data = event.data.value;
					if (!data) return;
					entryPoints.update((e) => e.filter((e) => e.id !== data.id));
					break;
				}
				case "router": {
					const data = event.data.value;
					if (!data) return;
					routers.update((r) => r.filter((r) => r.id !== data.id));
					break;
				}
				case "service": {
					const data = event.data.value;
					if (!data) return;
					services.update((s) => s.filter((s) => s.id !== data.id));
					break;
				}
				case "middleware": {
					const data = event.data.value;
					if (!data) return;
					middlewares.update((m) => m.filter((m) => m.id !== data.id));
					break;
				}
				case "serversTransport": {
					const data = event.data.value;
					if (!data) return;
					serversTransports.update((s) => s.filter((s) => s.id !== data.id));
					break;
				}
				case "dnsProvider": {
					const data = event.data.value;
					if (!data) return;
					dnsProviders.update((d) => d.filter((d) => d.id !== data.id));
					break;
				}
				case "traefikInstance": {
					const data = event.data.value;
					if (!data) return;
					traefikInstances.update((t) => t.filter((t) => t.id !== data.id));
					break;
				}
			}
			break;
	}
}

// Preload data
async function preload() {
	// Global
	profileClient.listProfiles({}).then((response) => {
		profiles.set(response.profiles);
	});
	userClient.listUsers({}).then((response) => {
		users.set(response.users);
	});
	dnsClient.listDnsProviders({}).then((response) => {
		dnsProviders.set(response.dnsProviders);
	});

	// Profile specific
	agentClient.listAgents({ profileId: profile.id }).then((response) => {
		agents.set(response.agents);
	});
	entryPointClient
		.listEntryPoints({ profileId: profile.id })
		.then((response) => {
			entryPoints.set(response.entryPoints);
		});
	routerClient.listRouters({ profileId: profile.id }).then((response) => {
		routers.set(response.routers);
	});
	serviceClient.listServices({ profileId: profile.id }).then((response) => {
		services.set(response.services);
	});
	middlewareClient
		.listMiddlewares({ profileId: profile.id })
		.then((response) => {
			middlewares.set(response.middlewares);
		});
	serversTransportClient
		.listServersTransports({ profileId: profile.id })
		.then((response) => {
			serversTransports.set(response.serversTransports);
		});
	traefikClient
		.listTraefikInstances({ profileId: profile.id })
		.then((response) => {
			traefikInstances.set(response.traefikInstances);
		});
}
