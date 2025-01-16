export interface Router {
	entryPoints?: string[];
	middlewares?: string[];
	service: string;
	rule: string;
	ruleSyntax?: string;
	priority?: number;
	tls: TLS;
	observability?: Observability;
	status?: string;
	using?: string[];
}

export function newRouter(): Router {
	return {
		entryPoints: [],
		middlewares: [],
		rule: '',
		ruleSyntax: 'v3',
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
	loadBalancer?: ServersLoadBalancer;
	weighted?: WeightedRoundRobin;
	mirroring?: Mirroring;
	failover?: Failover;
	status?: string;
	usedBy?: string[];
	serverStatus?: Record<string, string>;
}

export function newService(): Service {
	return {
		loadBalancer: {
			servers: []
		}
	};
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
