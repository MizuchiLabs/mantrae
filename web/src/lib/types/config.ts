import type { Domain } from './tls';

export interface Router {
	// Common fields
	name: string;
	provider?: string;
	status?: string;
	routerType: string;
	errorState?: ErrorState;
	dnsProvider?: number;

	entrypoints?: string[];
	middlewares?: string[];
	service: string;
	rule: string;
	priority?: number;
	tls: TLSConfig;
}

export function newRouter(): Router {
	return {
		name: '',
		provider: 'http',
		status: '',
		routerType: 'http',
		dnsProvider: 0,
		entrypoints: [],
		middlewares: [],
		rule: '',
		service: '',
		priority: 0,
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
	name: string;
	provider?: string;
	type?: string;
	status?: string;
	serviceType: string;
	serverStatus?: Record<string, string>;

	loadBalancer?: ServersLoadBalancer;
	weighted?: WeightedRoundRobin;
	mirroring?: Mirroring;
	failover?: Failover;
}

export function newService(): Service {
	return {
		name: '',
		provider: 'http',
		type: '',
		status: '',
		serviceType: 'http',
		serverStatus: {}
	};
}

export interface ErrorState {
	ssl?: string;
	dns?: string;
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
