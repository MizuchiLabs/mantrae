import type { ClientTLS } from './tls';

// HTTP  Middlewares ----------------------------------------------------------
export interface Middleware {
	// Common fields
	name: string;
	provider?: string;
	type?: string;
	status?: string;
	middlewareType: string;

	// HTTP-specific fields
	addPrefix?: AddPrefix;
	stripPrefix?: StripPrefix;
	stripPrefixRegex?: StripPrefixRegex;
	replacePath?: ReplacePath;
	replacePathRegex?: ReplacePathRegex;
	chain?: Chain;
	ipAllowList?: IPAllowList; // also for tcp
	headers?: Headers;
	errors?: ErrorPage;
	rateLimit?: RateLimit;
	redirectRegex?: RedirectRegex;
	redirectScheme?: RedirectScheme;
	basicAuth?: BasicAuth;
	digestAuth?: DigestAuth;
	forwardAuth?: ForwardAuth;
	inFlightReq?: InFlightReq;
	buffering?: Buffering;
	circuitBreaker?: CircuitBreaker;
	compress?: Compress;
	passTLSClientCert?: PassTLSClientCert;
	retry?: Retry;

	// TCP-specific fields
	inFlightConn?: TCPInFlightConn;
}

export function newMiddleware(): Middleware {
	return {
		name: '',
		provider: 'http',
		type: '',
		status: '',
		middlewareType: 'http'
	};
}

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
	tls?: ClientTLS;
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
