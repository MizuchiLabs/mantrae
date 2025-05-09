import type { TraefikConfig } from '$lib/types';

export type Middleware = HTTPMiddleware | TCPMiddleware;
export type SupportedMiddleware = SupportedMiddlewareHTTP | SupportedMiddlewareTCP;
export type SupportedMiddlewareHTTP = keyof Omit<HTTPMiddleware, 'name' | 'protocol' | 'type'>;
export type SupportedMiddlewareTCP = keyof Omit<TCPMiddleware, 'name' | 'protocol' | 'type'>;

interface BaseMiddleware {
	name: string;
	type?: string;
}

// HTTP Middlewares ----------------------------------------------------------
export interface HTTPMiddleware extends BaseMiddleware {
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
export interface TCPMiddleware extends BaseMiddleware {
	protocol: 'tcp';
	type?: SupportedMiddlewareTCP;
	ipAllowList?: TCPIPAllowList;
	inFlightConn?: TCPInFlightConn;
}

export interface UpsertMiddlewareParams {
	name: string;
	protocol: 'http' | 'tcp';
	type?: SupportedMiddleware;
	middleware?: {
		[K in SupportedMiddlewareHTTP]?: HTTPMiddleware[K];
	};
	tcpMiddleware?: {
		[K in SupportedMiddlewareTCP]?: TCPMiddleware[K];
	};
}

export interface DeleteMiddlewareParams {
	profileId: number;
	name: string;
	protocol: 'http' | 'tcp';
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
type HTTPMiddlewareDefaults = {
	[K in SupportedMiddlewareHTTP]: K extends keyof HTTPMiddleware
		? Exclude<HTTPMiddleware[K], undefined>
		: never;
};
export function getDefaultValuesForType<T extends SupportedMiddlewareHTTP>(
	type: T
): HTTPMiddlewareDefaults[T] {
	const defaults: HTTPMiddlewareDefaults = {
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

	return defaults[type];
}

// Helper to get default values for TCP middleware types
type TCPMiddlewareDefaults = {
	[K in SupportedMiddlewareTCP]: K extends keyof TCPMiddleware
		? Exclude<TCPMiddleware[K], undefined>
		: never;
};
export function getTCPDefaultValuesForType<T extends SupportedMiddlewareTCP>(
	type: T
): TCPMiddlewareDefaults[T] {
	const defaults: TCPMiddlewareDefaults = {
		ipAllowList: {
			sourceRange: []
		},
		inFlightConn: {
			amount: 0
		}
	};

	return defaults[type];
}

// Metadata per field
export type FieldMetadata = {
	placeholder?: string;
	description?: string;
	examples?: string[];
	validation?: {
		min?: number;
		max?: number;
		pattern?: string;
	};
};

type HTTPMiddlewareFieldMetadata = {
	[K in SupportedMiddlewareHTTP]: {
		[Field in keyof Required<Exclude<HTTPMiddleware[K], undefined>>]: FieldMetadata;
	};
};

type TCPMiddlewareFieldMetadata = {
	[K in SupportedMiddlewareTCP]: {
		[Field in keyof Required<Exclude<TCPMiddleware[K], undefined>>]: FieldMetadata;
	};
};

export type MiddlewareFieldMetadata = HTTPMiddlewareFieldMetadata & TCPMiddlewareFieldMetadata;
export const middlewareFieldMetadata: MiddlewareFieldMetadata = {
	// Basic HTTP Middlewares
	addPrefix: {
		prefix: {
			placeholder: '/api',
			description: 'Adds a path prefix to the existing request path',
			examples: ['/api', '/v1', '/prefix']
		}
	},

	basicAuth: {
		users: {
			placeholder: 'test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/',
			description: 'List of authorized users in Basic Authentication format',
			examples: ['user:hashed-password']
		},
		usersFile: {
			placeholder: '/path/to/users',
			description: 'Path to the file containing authorized users'
		},
		realm: {
			placeholder: 'Traefik',
			description: 'Authentication realm to display in the browser'
		},
		removeHeader: {
			description: 'If true, removes the authorization header before forwarding the request'
		},
		headerField: {
			placeholder: 'X-WebAuth-User',
			description: 'Sets the header field for the authenticated user'
		}
	},

	digestAuth: {
		users: {
			placeholder: 'test:traefik:a2688e031edb4be6a3797f3882655c05',
			description: 'List of authorized users in Digest Authentication format',
			examples: ['user:realm:hashed-password']
		},
		usersFile: {
			placeholder: '/path/to/users',
			description: 'Path to the file containing authorized users'
		},
		removeHeader: {
			description: 'If true, removes the authorization header before forwarding the request'
		},
		realm: {
			placeholder: 'traefik',
			description: 'Authentication realm to display in the browser'
		},
		headerField: {
			placeholder: 'X-WebAuth-User',
			description: 'Sets the header field for the authenticated user'
		}
	},

	buffering: {
		maxRequestBodyBytes: {
			placeholder: '2097152',
			description: 'Maximum size of the request body in bytes (2MB default)',
			validation: { min: 0 }
		},
		memRequestBodyBytes: {
			placeholder: '1048576',
			description: 'Maximum size of the request body in memory (1MB default)',
			validation: { min: 0 }
		},
		maxResponseBodyBytes: {
			placeholder: '2097152',
			description: 'Maximum size of the response body in bytes (2MB default)',
			validation: { min: 0 }
		},
		memResponseBodyBytes: {
			placeholder: '1048576',
			description: 'Maximum size of the response body in memory (1MB default)',
			validation: { min: 0 }
		},
		retryExpression: {
			placeholder: 'IsNetworkError() && Attempts() <= 2',
			description: 'Retry request if expression matches'
		}
	},

	chain: {
		middlewares: {
			placeholder: 'auth-middleware, rate-limit',
			description: 'List of middleware names to be chained together',
			examples: ['auth-middleware', 'rate-limit', 'compress']
		}
	},

	circuitBreaker: {
		expression: {
			placeholder: 'NetworkErrorRatio() > 0.5',
			description: 'Expression that triggers the circuit breaker',
			examples: ['NetworkErrorRatio() > 0.5', 'ResponseCodeRatio(500, 600, 0, 600) > 0.5']
		},
		checkPeriod: {
			placeholder: '10s',
			description: 'Interval between successive checks',
			examples: ['10s', '1m', '1h']
		},
		fallbackDuration: {
			placeholder: '10s',
			description: 'Duration for which the circuit breaker stays open',
			examples: ['10s', '1m', '1h']
		},
		recoveryDuration: {
			placeholder: '10s',
			description: 'Duration for which the circuit breaker stays in recovery state',
			examples: ['10s', '1m', '1h']
		},
		responseCode: {
			placeholder: '503',
			description: 'Response code when circuit breaker is open',
			validation: { min: 100, max: 599 }
		}
	},

	compress: {
		excludedContentTypes: {
			placeholder: 'image/jpeg, image/png',
			description: 'List of content types to exclude from compression',
			examples: ['image/jpeg', 'image/png', 'application/pdf']
		},
		includeContentTypes: {
			placeholder: 'text/html, text/plain',
			description: 'List of content types to include for compression',
			examples: ['text/html', 'text/plain', 'application/json']
		},
		minResponseBodyBytes: {
			placeholder: '1024',
			description: 'Minimum response body size for compression',
			validation: { min: 0 }
		},
		defaultEncoding: {
			placeholder: 'gzip',
			description: 'Default compression encoding',
			examples: ['gzip', 'deflate']
		}
	},

	errorPage: {
		status: {
			placeholder: '500-599',
			description: 'HTTP status codes to match',
			examples: ['404', '500-599']
		},
		service: {
			placeholder: 'error-handler',
			description: 'Service to call when an error occurs'
		},
		query: {
			placeholder: '/error?code={status}&url={url}',
			description: 'Query string to use when calling error service'
		}
	},

	forwardAuth: {
		address: {
			placeholder: 'http://auth.example.com',
			description: 'Authentication server address'
		},
		tls: {
			ca: {
				placeholder: '/path/to/ca.crt',
				description: 'Certificate Authority certificate path'
			},
			cert: {
				placeholder: '/path/to/server.crt',
				description: 'Client certificate path'
			},
			key: {
				placeholder: '/path/to/server.key',
				description: 'Client certificate key path'
			},
			insecureSkipVerify: {
				description: 'Skip TLS certificate verification (not recommended)'
			}
		},
		trustForwardHeader: {
			description: 'Trust X-Forwarded-* headers from previous proxy'
		},
		authResponseHeaders: {
			placeholder: 'X-Auth-User, X-Auth-Role',
			description: 'Headers to forward from auth response'
		},
		authRequestHeaders: {
			placeholder: 'X-Auth-User, X-Auth-Role',
			description: 'Headers to forward to authentication server'
		}
	},

	rateLimit: {
		average: {
			placeholder: '100',
			description: 'Average requests per period',
			validation: { min: 1 }
		},
		period: {
			placeholder: '1m',
			description: 'Sampling period',
			examples: ['10s', '1m', '1h']
		},
		burst: {
			placeholder: '200',
			description: 'Maximum requests allowed to exceed the average',
			validation: { min: 1 }
		},
		sourceCriterion: {
			placeholder: 'request.host',
			description: 'Source criterion'
		}
	},

	headers: {
		customRequestHeaders: {
			placeholder: '{"X-Custom-Header": "value"}',
			description: 'Custom headers to add to request'
		},
		customResponseHeaders: {
			placeholder: '{"X-Response-Header": "value"}',
			description: 'Custom headers to add to response'
		},
		accessControlAllowOriginList: {
			placeholder: 'https://example.com',
			description: 'Allowed CORS origins',
			examples: ['https://example.com', '*']
		}
	},

	retry: {
		attempts: {
			placeholder: '3',
			description: 'Number of retry attempts',
			validation: { min: 1 }
		},
		initialInterval: {
			placeholder: '100ms',
			description: 'Initial retry interval',
			examples: ['100ms', '1s', '2s']
		}
	},

	inFlightReq: {
		amount: {
			placeholder: '100',
			description: 'Maximum number of allowed simultaneous requests',
			validation: { min: 1 }
		},
		sourceCriterion: {
			ipStrategy: {
				depth: {
					placeholder: '1',
					description: 'Number of IPs to take from X-Forwarded-For header',
					validation: { min: 0 }
				},
				excludedIPs: {
					placeholder: '127.0.0.1',
					description: 'IPs to exclude from X-Forwarded-For header',
					examples: ['127.0.0.1', '192.168.0.0/16']
				}
			},
			requestHeaderName: {
				placeholder: 'X-Real-IP',
				description: 'Header to use as source',
				examples: ['X-Real-IP', 'X-Forwarded-For']
			},
			requestHost: {
				description: 'Use request host as source'
			}
		}
	},

	passTLSClientCert: {
		pem: {
			description: 'Pass the PEM-formatted client certificate in a header'
		},
		info: {
			notAfter: {
				description: 'Add NotAfter info in a header'
			},
			notBefore: {
				description: 'Add NotBefore info in a header'
			},
			sans: {
				description: 'Add Subject Alternative Names info in a header'
			},
			serialNumber: {
				description: 'Add serial number info in a header'
			},
			subject: {
				country: {
					description: 'Add subject country info in a header'
				},
				province: {
					description: 'Add subject province info in a header'
				},
				locality: {
					description: 'Add subject locality info in a header'
				},
				organization: {
					description: 'Add subject organization info in a header'
				},
				organizationalUnit: {
					description: 'Add subject organizational unit info in a header'
				},
				commonName: {
					description: 'Add subject common name info in a header'
				},
				serialNumber: {
					description: 'Add subject serial number info in a header'
				},
				domainComponent: {
					description: 'Add subject domain component info in a header'
				}
			},
			issuer: {
				country: {
					description: 'Add issuer country info in a header'
				},
				province: {
					description: 'Add issuer province info in a header'
				},
				locality: {
					description: 'Add issuer locality info in a header'
				},
				organization: {
					description: 'Add issuer organization info in a header'
				},
				commonName: {
					description: 'Add issuer common name info in a header'
				},
				serialNumber: {
					description: 'Add issuer serial number info in a header'
				},
				domainComponent: {
					description: 'Add issuer domain component info in a header'
				}
			}
		}
	},

	redirectRegex: {
		regex: {
			placeholder: '^/redirect/(.*)$',
			description: 'Regular expression to match path for redirection',
			examples: ['^/redirect/(.*)$', '^/old-api/(.*)$']
		},
		replacement: {
			placeholder: '/new/$1',
			description: 'Replacement path for the redirection',
			examples: ['/new/$1', 'https://example.org/$1']
		},
		permanent: {
			description: 'Use permanent redirection (301 instead of 302)'
		}
	},

	redirectScheme: {
		scheme: {
			placeholder: 'https',
			description: 'Scheme to redirect to',
			examples: ['https', 'http']
		},
		port: {
			placeholder: '443',
			description: 'Port to redirect to',
			examples: ['443', '8443']
		},
		permanent: {
			description: 'Use permanent redirection (301 instead of 302)'
		}
	},

	replacePath: {
		path: {
			placeholder: '/new-path',
			description: 'Path to replace the matched URL',
			examples: ['/new-path', '/api/v2']
		}
	},

	replacePathRegex: {
		regex: {
			placeholder: '^/old-path/(.*)$',
			description: 'Regular expression to match path for replacement',
			examples: ['^/old-path/(.*)$', '^/api/v1/(.*)$']
		},
		replacement: {
			placeholder: '/new-path/$1',
			description: 'Replacement pattern for the path',
			examples: ['/new-path/$1', '/api/v2/$1']
		}
	},

	stripPrefix: {
		prefixes: {
			placeholder: '/api',
			description: 'List of path prefixes to strip from URL',
			examples: ['/api', '/api/v1', '/legacy']
		},
		forceSlash: {
			description: 'Force adding a trailing slash to the path'
		}
	},

	stripPrefixRegex: {
		regex: {
			placeholder: '^/api/.*',
			description: 'Regular expressions to match prefixes to strip from URL',
			examples: ['^/api/.*', '^/old/.*/api']
		}
	},

	// TCP Middlewares
	ipAllowList: {
		sourceRange: {
			placeholder: '192.168.1.0/24',
			description: 'IP ranges to allow',
			examples: ['192.168.1.0/24', '10.0.0.0/8']
		},
		ipStrategy: {
			description: 'IP strategy to use'
		}
	},
	inFlightConn: {
		amount: {
			placeholder: '100',
			description: 'Maximum amount of allowed simultaneous connections',
			validation: { min: 1 }
		}
	}
};

export function getMetadataForMiddleware(type: SupportedMiddleware) {
	return middlewareFieldMetadata[type] || {};
}
