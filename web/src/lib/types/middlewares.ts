import type { SupportedMiddleware } from '$lib/components/forms/mw_registry';
import type { BaseTraefikConfig } from '$lib/types';
import { z } from 'zod';

// HTTP Middlewares ----------------------------------------------------------
export interface Middleware {
	name: string;
	protocol: 'http' | 'tcp';
	type?: SupportedMiddleware;
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
	tcpIpAllowList?: TCPIPAllowList;
	tcpInFlightConn?: TCPInFlightConn;
	plugin?: Record<string, unknown>;
}

export interface UpsertMiddlewareParams {
	name: string;
	protocol: 'http' | 'tcp';
	type?: string;
	middleware?: Middleware;
	tcpMiddleware?: Middleware;
}

export function flattenMiddlewareData(config: BaseTraefikConfig): Middleware[] {
	const flatMiddleware: Middleware[] = [];
	if (!config) return flatMiddleware;

	Object.entries(config.middlewares || {}).forEach(([name, middleware]) => {
		if (!middleware) {
			flatMiddleware.push({
				name,
				protocol: 'http'
			});
			return;
		}
		const [type, details] = Object.entries(middleware)[0] || [undefined, {}];
		flatMiddleware.push({
			name,
			protocol: 'http',
			type,
			...details
		});
	});

	Object.entries(config.tcpMiddlewares || {}).forEach(([name, middleware]) => {
		if (!middleware) {
			flatMiddleware.push({
				name,
				protocol: 'http'
			});
			return;
		}
		const [type, details] = Object.entries(middleware)[0] || [undefined, {}];
		flatMiddleware.push({
			name,
			protocol: 'tcp',
			type,
			...details
		});
	});

	return flatMiddleware;
}

export const AddPrefixSchema = z.object({
	prefix: z.string().trim()
});
type AddPrefix = z.infer<typeof AddPrefixSchema>;

export const BasicAuthSchema = z.object({
	users: z.array(z.string().trim()),
	usersFile: z.string().trim().optional(),
	realm: z.string().trim().optional(),
	headerField: z.string().trim().optional(),
	removeHeader: z.boolean().optional()
});
type BasicAuth = z.infer<typeof BasicAuthSchema>;

export const DigestAuthSchema = z.object({
	users: z.array(z.string().trim()),
	usersFile: z.string().trim().optional(),
	removeHeader: z.boolean().optional(),
	realm: z.string().optional(),
	headerField: z.string().optional()
});
type DigestAuth = z.infer<typeof DigestAuthSchema>;

export const BufferingSchema = z.object({
	maxRequestBodyBytes: z.number().optional(),
	memRequestBodyBytes: z.number().optional(),
	maxResponseBodyBytes: z.number().optional(),
	memResponseBodyBytes: z.number().optional(),
	retryExpression: z.string().trim().optional()
});
type Buffering = z.infer<typeof BufferingSchema>;

export const ChainSchema = z.object({
	middlewares: z.array(z.string())
});
type Chain = z.infer<typeof ChainSchema>;

export const CircuitBreakerSchema = z.object({
	expression: z.string().optional(),
	checkPeriod: z.string().optional(),
	fallbackDuration: z.string().optional(),
	recoveryDuration: z.string().optional(),
	responseCode: z.number().optional()
});
type CircuitBreaker = z.infer<typeof CircuitBreakerSchema>;

export const CompressSchema = z.object({
	excludedContentTypes: z.array(z.string()).optional(),
	includeContentTypes: z.array(z.string()).optional(),
	minResponseBodyBytes: z.number().optional(),
	defaultEncoding: z.string().optional()
});
type Compress = z.infer<typeof CompressSchema>;

export const ErrorPageSchema = z.object({
	status: z.array(z.string()).optional(),
	service: z.string().optional(),
	query: z.string().optional()
});
type ErrorPage = z.infer<typeof ErrorPageSchema>;

export const ForwardAuthSchema = z.object({
	address: z.string().trim(),
	tls: z
		.object({
			ca: z.string().optional(),
			caOptional: z.boolean().optional(),
			cert: z.string().optional(),
			key: z.string().optional(),
			insecureSkipVerify: z.boolean().optional()
		})
		.optional(),
	trustForwardHeader: z.boolean().optional(),
	authResponseHeaders: z.array(z.string()).optional(),
	authResponseHeadersRegex: z.string().optional(),
	authRequestHeaders: z.array(z.string()).optional(),
	addAuthCookiesToResponse: z.array(z.string()).optional()
});
type ForwardAuth = z.infer<typeof ForwardAuthSchema>;

export const HeadersSchema = z.object({
	customRequestHeaders: z.record(z.string()).optional(),
	customResponseHeaders: z.record(z.string()).optional(),
	accessControlAllowCredentials: z.boolean().optional(),
	accessControlAllowHeaders: z.array(z.string()).optional(),
	accessControlAllowMethods: z.array(z.string()).optional(),
	accessControlAllowOriginList: z.array(z.string()).optional(),
	accessControlAllowOriginListRegex: z.array(z.string()).optional(),
	accessControlExposeHeaders: z.array(z.string()).optional(),
	accessControlMaxAge: z.number().optional(),
	addVaryHeader: z.boolean().optional(),
	allowedHosts: z.array(z.string()).optional(),
	hostsProxyHeaders: z.array(z.string()).optional(),
	sslProxyHeaders: z.record(z.string()).optional(),
	stsSeconds: z.number().optional(),
	stsIncludeSubdomains: z.boolean().optional(),
	stsPreload: z.boolean().optional(),
	forceSTSHeader: z.boolean().optional(),
	frameDeny: z.boolean().optional(),
	customFrameOptionsValue: z.string().optional(),
	contentTypeNosniff: z.boolean().optional(),
	browserXssFilter: z.boolean().optional(),
	customBrowserXSSValue: z.string().optional(),
	contentSecurityPolicy: z.string().optional(),
	publicKey: z.string().optional(),
	referrerPolicy: z.string().optional(),
	permissionsPolicy: z.string().optional(),
	isDevelopment: z.boolean().optional()
});
type Headers = z.infer<typeof HeadersSchema>;

export const IPAllowListSchema = z.object({
	sourceRange: z.array(z.string()).optional(),
	ipStrategy: z
		.object({
			depth: z.number().optional(),
			excludedIPs: z.array(z.string()).optional()
		})
		.optional()
});
type IPAllowList = z.infer<typeof IPAllowListSchema>;

export const InFlightReqSchema = z.object({
	amount: z.number(),
	sourceCriterion: z
		.object({
			ipStrategy: z
				.object({
					depth: z.number().optional(),
					excludedIPs: z.array(z.string()).optional()
				})
				.optional(),
			requestHeaderName: z.string().optional(),
			requestHost: z.boolean().optional()
		})
		.optional()
});
type InFlightReq = z.infer<typeof InFlightReqSchema>;

export const PassTLSClientCertSchema = z.object({
	pem: z.boolean().optional(),
	info: z
		.object({
			notAfter: z.boolean().optional(),
			notBefore: z.boolean().optional(),
			sans: z.boolean().optional(),
			serialNumber: z.boolean().optional(),
			subject: z
				.object({
					country: z.boolean().optional(),
					province: z.boolean().optional(),
					locality: z.boolean().optional(),
					organization: z.boolean().optional(),
					organizationalUnit: z.boolean().optional(),
					commonName: z.boolean().optional(),
					serialNumber: z.boolean().optional(),
					domainComponent: z.boolean().optional()
				})
				.optional(),
			issuer: z
				.object({
					country: z.boolean().optional(),
					province: z.boolean().optional(),
					locality: z.boolean().optional(),
					organization: z.boolean().optional(),
					commonName: z.boolean().optional(),
					serialNumber: z.boolean().optional(),
					domainComponent: z.boolean().optional()
				})
				.optional()
		})
		.optional()
});
type PassTLSClientCert = z.infer<typeof PassTLSClientCertSchema>;

export const RateLimitSchema = z.object({
	average: z.number().optional(),
	period: z.string().optional(),
	burst: z.number().optional(),
	sourceCriterion: z
		.object({
			ipStrategy: z
				.object({
					depth: z.number().optional(),
					excludedIPs: z.array(z.string()).optional()
				})
				.optional(),
			requestHeaderName: z.string().optional(),
			requestHost: z.boolean().optional()
		})
		.optional()
});
type RateLimit = z.infer<typeof RateLimitSchema>;

export const RedirectRegexSchema = z.object({
	regex: z.string().optional(),
	replacement: z.string().optional(),
	permanent: z.boolean().optional()
});
type RedirectRegex = z.infer<typeof RedirectRegexSchema>;

export const RedirectSchemeSchema = z.object({
	scheme: z.string().optional(),
	port: z.string().optional(),
	permanent: z.boolean().optional()
});
type RedirectScheme = z.infer<typeof RedirectSchemeSchema>;

export const ReplacePathSchema = z.object({
	path: z.string().trim()
});
type ReplacePath = z.infer<typeof ReplacePathSchema>;

export const ReplacePathRegexSchema = z.object({
	regex: z.string().optional(),
	replacement: z.string().optional()
});
type ReplacePathRegex = z.infer<typeof ReplacePathRegexSchema>;

export const RetrySchema = z.object({
	attempts: z.number().optional(),
	initialInterval: z.string().optional()
});
type Retry = z.infer<typeof RetrySchema>;

export const StripPrefixSchema = z.object({
	prefixes: z.array(z.string()).optional(),
	forceSlash: z.boolean().optional()
});
type StripPrefix = z.infer<typeof StripPrefixSchema>;

export const StripPrefixRegexSchema = z.object({
	regex: z.array(z.string()).optional()
});
type StripPrefixRegex = z.infer<typeof StripPrefixRegexSchema>;

// TCP Middlewares ------------------------------------------------------------
export const TCPIPAllowListSchema = z.object({
	sourceRange: z.array(z.string()).optional()
});
type TCPIPAllowList = z.infer<typeof TCPIPAllowListSchema>;

export const TCPInFlightConnSchema = z.object({
	amount: z.number().optional()
});
type TCPInFlightConn = z.infer<typeof TCPInFlightConnSchema>;
