import type { TraefikConfig } from '$lib/types';

export interface Router {
	name: string;
	protocol: 'http' | 'tcp' | 'udp';
	entryPoints?: string[];
	middlewares?: string[];
	service: string;
	rule?: string;
	ruleSyntax?: string;
	priority?: number;
	tls?: TLS;
	observability?: Observability;
	status?: string;
}

export interface UpsertRouterParams {
	name: string;
	protocol: 'http' | 'tcp' | 'udp';
	router?: Router;
	tcpRouter?: Router;
	udpRouter?: Router;
	service?: Service;
	tcpService?: Service;
	udpService?: Service;
}

export interface DeleteRouterParams {
	profileId: number;
	name: string;
	protocol: 'http' | 'tcp' | 'udp';
}

// Transform function to flatten the router data
export function flattenRouterData(configs: TraefikConfig[]): Router[] {
	const flatRouters: Router[] = [];
	if (!configs) return flatRouters;

	for (const base of configs) {
		const config = base.config;
		if (!config) continue;
		// Process HTTP Routers
		Object.entries(config.routers || {}).forEach(([name, router]) => {
			flatRouters.push({
				name,
				protocol: 'http',
				entryPoints: router.entryPoints || [],
				middlewares: router.middlewares || [],
				service: router.service,
				rule: router.rule,
				ruleSyntax: router.ruleSyntax,
				priority: router.priority,
				tls: router.tls,
				observability: router.observability,
				status: router.status
			});
		});

		// Process TCP Routers
		Object.entries(config.tcpRouters || {}).forEach(([name, router]) => {
			flatRouters.push({
				name,
				protocol: 'tcp',
				entryPoints: router.entryPoints || [],
				middlewares: router.middlewares || [],
				service: router.service,
				rule: router.rule,
				ruleSyntax: router.ruleSyntax,
				priority: router.priority,
				tls: router.tls,
				observability: router.observability,
				status: router.status
			});
		});

		// Process UDP Routers
		Object.entries(config.udpRouters || {}).forEach(([name, router]) => {
			flatRouters.push({
				name,
				protocol: 'udp',
				entryPoints: router.entryPoints || [],
				service: router.service,
				observability: router.observability,
				status: router.status
			});
		});
	}

	return flatRouters;
}

export interface Service {
	name: string;
	protocol: 'http' | 'tcp' | 'udp';
	loadBalancer?: ServersLoadBalancer;
	weighted?: WeightedRoundRobin;
	mirroring?: Mirroring;
	failover?: Failover;
	status?: string;
	serverStatus?: Record<string, string>;
}

export function newService(): Service {
	return {
		name: '',
		protocol: 'http',
		loadBalancer: {
			servers: []
		}
	};
}

export function flattenServiceData(configs: TraefikConfig[]): Service[] {
	const flatServices: Service[] = [];
	if (!configs) return flatServices;

	for (const base of configs) {
		const config = base.config;
		if (!config) continue;
		// Process HTTP Services
		Object.entries(config.services || {}).forEach(([name, service]) => {
			flatServices.push({
				name,
				protocol: 'http',
				loadBalancer: service.loadBalancer,
				weighted: service.weighted,
				mirroring: service.mirroring,
				failover: service.failover,
				status: service.status,
				serverStatus: service.serverStatus
			});
		});

		// Process TCP Services
		Object.entries(config.tcpServices || {}).forEach(([name, service]) => {
			flatServices.push({
				name,
				protocol: 'tcp',
				loadBalancer: service.loadBalancer,
				weighted: service.weighted,
				status: service.status,
				serverStatus: service.serverStatus
			});
		});

		// Process UDP Services
		Object.entries(config.udpServices || {}).forEach(([name, service]) => {
			flatServices.push({
				name,
				protocol: 'udp',
				loadBalancer: service.loadBalancer,
				weighted: service.weighted,
				status: service.status,
				serverStatus: service.serverStatus
			});
		});
	}

	return flatServices;
}

export interface TLS {
	passthrough?: boolean; // TCP Router only
	options?: string;
	certResolver?: string;
	domains?: Domain[];
}

export interface Domain {
	main?: string;
	sans?: string[];
}

export interface Observability {
	accessLogs: boolean;
	tracing: boolean;
	metrics: boolean;
}

export interface ServersLoadBalancer {
	servers?: Server[]; // for every service protocol
	sticky?: Sticky;
	healthCheck?: ServerHealthCheck;
	passHostHeader?: boolean;
	responseForwarding?: ResponseForwarding;
	serversTransport?: string;
	terminationDelay?: number; // TCP only
	proxyProtocol?: ProxyProtocol; // TCP only
}

export interface WeightedRoundRobin {
	services?: WRRService[]; // for every service protocol
	sticky?: Sticky; // HTTP only
	healthCheck?: Record<string, unknown>; // HTTP only
}

export interface Cookie {
	name?: string;
	secure?: boolean;
	httpOnly?: boolean;
	sameSite?: string;
}

export interface Failover {
	service?: string;
	fallback?: string;
	healthCheck?: Record<string, unknown>;
}

export interface ForwardingTimeouts {
	dialTimeout?: string;
	responseHeaderTimeout?: string;
	idleConnTimeout?: string;
	readIdleTimeout?: string;
	pingTimeout?: string;
}

export interface MirrorService {
	name?: string;
	percent?: number;
}

export interface Mirroring {
	service?: string;
	maxBodySize?: number;
	mirrors?: MirrorService[];
	healthCheck?: Record<string, unknown>;
}

export interface ResponseForwarding {
	flushInterval?: string;
}

export interface Server {
	url?: string;
	address?: string;
}

export interface ServerHealthCheck {
	scheme?: string;
	path?: string;
	method?: string;
	port?: number;
	interval?: string;
	timeout?: string;
	hostname?: string;
	followRedirects?: boolean;
	headers?: Record<string, string>;
}

export interface Sticky {
	cookie?: Cookie;
}

export interface WRRService {
	name?: string;
	weight?: number;
}

export interface ProxyProtocol {
	version?: number;
}
