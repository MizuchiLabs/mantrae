import type { TraefikConfig } from '$lib/types';

export type Middleware = HTTPMiddleware | TCPMiddleware;
export type SupportedMiddleware = SupportedMiddlewareHTTP | SupportedMiddlewareTCP;
export type SupportedMiddlewareHTTP = keyof Omit<HTTPMiddleware, 'name' | 'protocol' | 'type'>;
export type SupportedMiddlewareTCP = keyof Omit<TCPMiddleware, 'name' | 'protocol' | 'type'>;

// HTTP Middlewares ----------------------------------------------------------
export interface HTTPMiddleware {
	name: string;
	protocol: 'http';
	type?: SupportedMiddlewareHTTP;
	addPrefix?: AddPrefix;
	basicAuth?: BasicAuth;
	digestAuth?: DigestAuth;
	buffering?: Buffering;
	chain?: Chain;
	circuitBreaker?: CircuitBreaker;
	compress?: Compress;
	errorPage?: ErrorPage;
	forwardAuth?: ForwardAuth;
	headers?: Headers;
	ipAllowList?: IPAllowList;
	inFlightReq?: InFlightReq;
	passTLSClientCert?: PassTLSClientCert;
	rateLimit?: RateLimit;
	redirectRegex?: RedirectRegex;
	redirectScheme?: RedirectScheme;
	replacePath?: ReplacePath;
	replacePathRegex?: ReplacePathRegex;
	retry?: Retry;
	stripPrefix?: StripPrefix;
	stripPrefixRegex?: StripPrefixRegex;
	plugin?: Record<string, Record<string, unknown>>;
}

// TCP Middlewares -----------------------------------------------------------
export interface TCPMiddleware {
	name: string;
	protocol: 'tcp';
	type?: SupportedMiddlewareTCP;
	ipAllowList?: TCPIPAllowList;
	inFlightConn?: TCPInFlightConn;
}

export interface UpsertMiddlewareParams {
	name: string;
	protocol: 'http' | 'tcp';
	type?: SupportedMiddleware;
	middleware?: HTTPMiddleware;
	tcpMiddleware?: TCPMiddleware;
}

export function flattenMiddlewareData(configs: TraefikConfig[]): Middleware[] {
	const flatMiddleware: Middleware[] = [];
	if (!configs) return flatMiddleware;

	for (const base of configs) {
		const config = base.config;
		if (!config) continue;
		Object.entries(config.middlewares || {}).forEach(([name, middleware]) => {
			if (middleware) {
				const [type, details] = Object.entries(middleware)[0] || [undefined, {}];
				flatMiddleware.push({
					name,
					protocol: 'http',
					type,
					...details
				});
			} else {
				flatMiddleware.push({
					name,
					protocol: 'http'
				});
			}
		});

		Object.entries(config.tcpMiddlewares || {}).forEach(([name, middleware]) => {
			if (middleware) {
				const [type, details] = Object.entries(middleware)[0] || [undefined, {}];
				flatMiddleware.push({
					name,
					protocol: 'tcp',
					type,
					...details
				});
			} else {
				flatMiddleware.push({
					name,
					protocol: 'tcp'
				});
			}
		});
	}

	return flatMiddleware;
}

// Middleware keys (excluding name, protocol, and type)
export const HTTPMiddlewareKeys = [
	{ value: 'addPrefix', label: 'Add Prefix' },
	{ value: 'basicAuth', label: 'Basic Auth' },
	{ value: 'digestAuth', label: 'Digest Auth' },
	{ value: 'buffering', label: 'Buffering' },
	{ value: 'chain', label: 'Chain' },
	{ value: 'circuitBreaker', label: 'Circuit Breaker' },
	{ value: 'compress', label: 'Compress' },
	{ value: 'errorPage', label: 'Error Page' },
	{ value: 'forwardAuth', label: 'Forward Auth' },
	{ value: 'headers', label: 'Headers' },
	{ value: 'ipAllowList', label: 'IP Allow List' },
	{ value: 'inFlightReq', label: 'In-Flight Request' },
	{ value: 'passTLSClientCert', label: 'Pass TLS Client Cert' },
	{ value: 'rateLimit', label: 'Rate Limit' },
	{ value: 'redirectRegex', label: 'Redirect Regex' },
	{ value: 'redirectScheme', label: 'Redirect Scheme' },
	{ value: 'replacePath', label: 'Replace Path' },
	{ value: 'replacePathRegex', label: 'Replace Path Regex' },
	{ value: 'retry', label: 'Retry' },
	{ value: 'stripPrefix', label: 'Strip Prefix' },
	{ value: 'stripPrefixRegex', label: 'Strip Prefix Regex' },
	{ value: 'plugin', label: 'Plugin' }
];

export const TCPMiddlewareKeys = [
	{ value: 'ipAllowList', label: 'IP Allow List' },
	{ value: 'inFlightConn', label: 'In-Flight Connection' }
];

export interface AddPrefix {
	prefix?: string;
}

export interface BasicAuth {
	users?: string[];
	usersFile?: string;
	realm?: string;
	removeHeader?: boolean;
	headerField?: string;
}

export interface DigestAuth {
	users?: string[];
	usersFile?: string;
	removeHeader?: boolean;
	realm?: string;
	headerField?: string;
}

export interface Buffering {
	maxRequestBodyBytes?: number;
	memRequestBodyBytes?: number;
	maxResponseBodyBytes?: number;
	memResponseBodyBytes?: number;
	retryExpression?: string;
}

export interface Chain {
	middlewares?: string[];
}

export interface CircuitBreaker {
	expression?: string;
	checkPeriod?: string;
	fallbackDuration?: string;
	recoveryDuration?: string;
	responseCode?: number;
}

export interface Compress {
	excludedContentTypes?: string[];
	includeContentTypes?: string[];
	minResponseBodyBytes?: number;
	defaultEncoding?: string;
}

export interface ErrorPage {
	status?: string[];
	service?: string;
	query?: string;
}

export interface ForwardAuth {
	address?: string;
	tls?: {
		ca?: string;
		caOptional?: boolean;
		cert?: string;
		key?: string;
		insecureSkipVerify?: boolean;
	};
	trustForwardHeader?: boolean;
	authResponseHeaders?: string[];
	authResponseHeadersRegex?: string;
	authRequestHeaders?: string[];
	addAuthCookiesToResponse?: string[];
}

export interface Headers {
	customRequestHeaders?: Record<string, string>;
	customResponseHeaders?: Record<string, string>;
	accessControlAllowCredentials?: boolean;
	accessControlAllowHeaders?: string[];
	accessControlAllowMethods?: string[];
	accessControlAllowOriginList?: string[];
	accessControlAllowOriginListRegex?: string[];
	accessControlExposeHeaders?: string[];
	accessControlMaxAge?: number;
	addVaryHeader?: boolean;
	allowedHosts?: string[];
	hostsProxyHeaders?: string[];
	sslProxyHeaders?: Record<string, string>;
	stsSeconds?: number;
	stsIncludeSubdomains?: boolean;
	stsPreload?: boolean;
	forceSTSHeader?: boolean;
	frameDeny?: boolean;
	customFrameOptionsValue?: string;
	contentTypeNosniff?: boolean;
	browserXssFilter?: boolean;
	customBrowserXSSValue?: string;
	contentSecurityPolicy?: string;
	publicKey?: string;
	referrerPolicy?: string;
	permissionsPolicy?: string;
	isDevelopment?: boolean;
}

export interface IPAllowList {
	sourceRange?: string[];
	ipStrategy?: IPStrategy;
}

export interface IPStrategy {
	depth?: number;
	excludedIPs?: string[];
}

export interface InFlightReq {
	amount?: number;
	sourceCriterion?: SourceCriterion;
}

export interface PassTLSClientCert {
	pem?: boolean;
	info?: TLSClientCertificateInfo;
}

export interface RateLimit {
	average?: number;
	period?: string;
	burst?: number;
	sourceCriterion?: SourceCriterion;
}

export interface RedirectRegex {
	regex?: string;
	replacement?: string;
	permanent?: boolean;
}

export interface RedirectScheme {
	scheme?: string;
	port?: string;
	permanent?: boolean;
}

export interface ReplacePath {
	path?: string;
}

export interface ReplacePathRegex {
	regex?: string;
	replacement?: string;
}

export interface Retry {
	attempts?: number;
	initialInterval?: string;
}

export interface SourceCriterion {
	ipStrategy?: IPStrategy;
	requestHeaderName?: string;
	requestHost?: boolean;
}

export interface StripPrefix {
	prefixes?: string[];
	forceSlash?: boolean;
}

export interface StripPrefixRegex {
	regex?: string[];
}

export interface TLSClientCertificateInfo {
	notAfter?: boolean;
	notBefore?: boolean;
	sans?: boolean;
	serialNumber?: boolean;
	subject?: TLSClientCertificateSubjectDNInfo;
	issuer?: TLSClientCertificateIssuerDNInfo;
}

export interface TLSClientCertificateIssuerDNInfo {
	country?: boolean;
	province?: boolean;
	locality?: boolean;
	organization?: boolean;
	commonName?: boolean;
	serialNumber?: boolean;
	domainComponent?: boolean;
}

export interface TLSClientCertificateSubjectDNInfo {
	country?: boolean;
	province?: boolean;
	locality?: boolean;
	organization?: boolean;
	organizationalUnit?: boolean;
	commonName?: boolean;
	serialNumber?: boolean;
	domainComponent?: boolean;
}

// TCP Middlewares ------------------------------------------------------------
export interface TCPIPAllowList {
	sourceRange?: string[];
}

export interface TCPInFlightConn {
	amount?: number;
}

// Helper to get default values for HTTP middleware types
export function getDefaultValuesForType(type: keyof HTTPMiddleware): Record<string, any> {
	const defaults: Record<string, any> = {
		addPrefix: {
			prefix: ''
		},
		basicAuth: {
			users: [],
			usersFile: '',
			realm: '',
			removeHeader: false,
			headerField: ''
		},
		digestAuth: {
			users: [],
			usersFile: '',
			removeHeader: false,
			realm: '',
			headerField: ''
		},
		buffering: {
			maxRequestBodyBytes: 0,
			memRequestBodyBytes: 0,
			maxResponseBodyBytes: 0,
			memResponseBodyBytes: 0,
			retryExpression: ''
		},
		chain: {
			middlewares: []
		},
		circuitBreaker: {
			expression: '',
			checkPeriod: '',
			fallbackDuration: '',
			recoveryDuration: '',
			responseCode: 0
		},
		compress: {
			excludedContentTypes: [],
			includeContentTypes: [],
			minResponseBodyBytes: 0,
			defaultEncoding: ''
		},
		errorPage: {
			status: [],
			service: '',
			query: ''
		},
		forwardAuth: {
			address: '',
			tls: {
				ca: '',
				caOptional: false,
				cert: '',
				key: '',
				insecureSkipVerify: false
			},
			trustForwardHeader: false,
			authResponseHeaders: [],
			authResponseHeadersRegex: '',
			authRequestHeaders: [],
			addAuthCookiesToResponse: []
		},
		headers: {
			customRequestHeaders: {},
			customResponseHeaders: {},
			accessControlAllowCredentials: false,
			accessControlAllowHeaders: [],
			accessControlAllowMethods: [],
			accessControlAllowOriginList: [],
			accessControlAllowOriginListRegex: [],
			accessControlExposeHeaders: [],
			accessControlMaxAge: 0,
			addVaryHeader: false,
			allowedHosts: [],
			hostsProxyHeaders: [],
			sslProxyHeaders: {},
			stsSeconds: 0,
			stsIncludeSubdomains: false,
			stsPreload: false,
			forceSTSHeader: false,
			frameDeny: false,
			customFrameOptionsValue: '',
			contentTypeNosniff: false,
			browserXssFilter: false,
			customBrowserXSSValue: '',
			contentSecurityPolicy: '',
			publicKey: '',
			referrerPolicy: '',
			permissionsPolicy: '',
			isDevelopment: false
		},
		ipAllowList: {
			sourceRange: [],
			ipStrategy: {
				depth: 0,
				excludedIPs: []
			}
		},
		inFlightReq: {
			amount: 0,
			sourceCriterion: {
				ipStrategy: {
					depth: 0,
					excludedIPs: []
				},
				requestHeaderName: '',
				requestHost: false
			}
		},
		passTLSClientCert: {
			pem: false,
			info: {
				notAfter: false,
				notBefore: false,
				sans: false,
				serialNumber: false,
				subject: {
					country: false,
					province: false,
					locality: false,
					organization: false,
					organizationalUnit: false,
					commonName: false,
					serialNumber: false,
					domainComponent: false
				},
				issuer: {
					country: false,
					province: false,
					locality: false,
					organization: false,
					commonName: false,
					serialNumber: false,
					domainComponent: false
				}
			}
		},
		rateLimit: {
			average: 0,
			period: '',
			burst: 0,
			sourceCriterion: {
				ipStrategy: {
					depth: 0,
					excludedIPs: []
				},
				requestHeaderName: '',
				requestHost: false
			}
		},
		redirectRegex: {
			regex: '',
			replacement: '',
			permanent: false
		},
		redirectScheme: {
			scheme: '',
			port: '',
			permanent: false
		},
		replacePath: {
			path: ''
		},
		replacePathRegex: {
			regex: '',
			replacement: ''
		},
		retry: {
			attempts: 0,
			initialInterval: ''
		},
		stripPrefix: {
			prefixes: [],
			forceSlash: false
		},
		stripPrefixRegex: {
			regex: []
		},
		plugin: {}
	};

	return defaults[type] || {};
}

// Helper to get default values for TCP middleware types
export function getTCPDefaultValuesForType(type: keyof TCPMiddleware): Record<string, any> {
	const defaults: Record<string, any> = {
		ipAllowList: {
			sourceRange: []
		},
		inFlightConn: {
			amount: 0
		}
	};

	return defaults[type] || {};
}
