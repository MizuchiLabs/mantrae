import type { Domain } from './tls';

export interface Router {
	routerType: string; // http, tcp, udp
	entrypoints?: string[];
	middlewares?: string[];
	service: string;
	rule: string;
	priority?: number;
	tls?: TLSConfig;
	status?: string;
	name?: string;
	provider?: string;
}

export interface Service {
	serviceType: string; // http, tcp, udp
	loadBalancer?: ServersLoadBalancer;
	weighted?: WeightedRoundRobin;
	mirroring?: Mirroring; // Specific to HTTP, can be ignored for others
	failover?: Failover; // Specific to HTTP, can be ignored for others
	serverStatus?: Record<string, string>;
	status?: string;
	name?: string;
	provider?: string;
	type?: string;
}

export interface TLSConfig {
	passthrough?: boolean; // Specific to TCP
	options?: string;
	certResolver?: string;
	domains?: Domain[];
}

export interface ServersLoadBalancer {
	sticky?: Sticky;
	servers?: Server[];
	healthCheck?: ServerHealthCheck;
	passHostHeader?: boolean; // Specific to HTTP
	responseForwarding?: ResponseForwarding; // Specific to HTTP
	serversTransport?: string; // Specific to HTTP
	terminationDelay?: number; // Specific to TCP
	proxyProtocol?: ProxyProtocol; // Specific to TCP
}

export interface WeightedRoundRobin {
	services?: WRRService[];
	sticky?: Sticky;
	healthCheck?: Record<string, unknown>;
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

export interface RouterTLSConfig {
	options?: string;
	certResolver?: string;
	domains?: Domain[];
}

export interface Server {
	url?: string;
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
