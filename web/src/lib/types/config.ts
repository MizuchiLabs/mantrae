import type { Domain } from './tls';

export interface Router {
	// Common fields
	id: string;
	profileId: number;
	name: string;
	provider?: string;
	protocol: string;
	status?: string;
	errors?: Errors;
	agentId?: string;
	dnsProvider?: number;

	entryPoints?: string[];
	middlewares?: string[];
	service: string;
	rule: string;
	ruleSyntax?: string;
	priority?: number;
	tls: TLSConfig;
}

export function newRouter(): Router {
	return {
		id: '',
		profileId: 0,
		name: '',
		provider: 'http',
		status: '',
		protocol: 'http',
		entryPoints: [],
		middlewares: [],
		rule: '',
		service: '',
		tls: {
			passthrough: false,
			options: '',
			certResolver: '',
			domains: []
		}
	};
}

export interface Service {
	// Common fields
	id: string;
	profileId: number;
	name: string;
	provider?: string;
	type?: string;
	status?: string;
	protocol: string;
	agentId?: string;
	serverStatus?: Record<string, string>;

	loadBalancer?: ServersLoadBalancer;
	weighted?: WeightedRoundRobin;
	mirroring?: Mirroring;
	failover?: Failover;
}

export function newService(): Service {
	return {
		id: '',
		profileId: 0,
		name: '',
		provider: 'http',
		type: '',
		status: '',
		protocol: 'http',
		serverStatus: {}
	};
}

export interface Errors {
	ssl?: string;
	dns?: string;
	agent?: string;
}

export interface TLSConfig {
	passthrough?: boolean; // TCP Router only
	options?: string;
	certResolver?: string;
	domains?: Domain[];
}

export interface ServersLoadBalancer {
	servers?: Server[]; // for every service type
	sticky?: Sticky;
	healthCheck?: ServerHealthCheck;
	passHostHeader?: boolean;
	responseForwarding?: ResponseForwarding;
	serversTransport?: string;
	terminationDelay?: number; // TCP only
	proxyProtocol?: ProxyProtocol; // TCP only
}

export interface WeightedRoundRobin {
	services?: WRRService[]; // for every service type
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
