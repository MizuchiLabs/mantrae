import type { Domain } from './tls';

export interface Router {
	// Common fields
	name: string;
	provider?: string;
	status?: string;
	routerType: string;
	dnsProvider?: string;

	entrypoints?: string[];
	middlewares?: string[];
	service: string;
	rule: string;
	priority?: number;
	tls: TLSConfig;
	tcpTLS?: TCPTLSConfig;
}

export function newRouter(): Router {
	return {
		name: '',
		provider: 'http',
		status: '',
		routerType: 'http',
		dnsProvider: '',
		entrypoints: [],
		middlewares: [],
		rule: '',
		service: '',
		priority: 0,
		tls: {
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

	// HTTP-specific fields
	loadBalancer?: ServersLoadBalancer;
	weighted?: WeightedRoundRobin;
	mirroring?: Mirroring;
	failover?: Failover;

	// TCP-specific fields
	tcpLoadBalancer?: TCPServersLoadBalancer;
	tcpWeighted?: TCPWeightedRoundRobin;

	// UDP-specific fields
	udpLoadBalancer?: UDPServersLoadBalancer;
	udpWeighted?: UDPWeightedRoundRobin;
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

export interface TLSConfig {
	options?: string;
	certResolver?: string;
	domains?: Domain[];
}

export interface TCPTLSConfig {
	passthrough?: boolean;
	options?: string;
	certResolver?: string;
	domains?: Domain[];
}

export interface ServersLoadBalancer {
	sticky?: Sticky;
	servers?: Server[];
	healthCheck?: ServerHealthCheck;
	passHostHeader?: boolean;
	responseForwarding?: ResponseForwarding;
	serversTransport?: string;
}

export interface TCPServersLoadBalancer {
	terminationDelay?: number;
	proxyProtocol?: ProxyProtocol;
	servers?: TCPServer[];
}

export interface UDPServersLoadBalancer {
	servers?: TCPServer[];
}

export interface WeightedRoundRobin {
	services?: WRRService[];
	sticky?: Sticky;
	healthCheck?: Record<string, unknown>;
}

export interface TCPWeightedRoundRobin {
	services?: WRRService[];
}

export interface UDPWeightedRoundRobin {
	services?: WRRService[];
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

export interface TCPServer {
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
