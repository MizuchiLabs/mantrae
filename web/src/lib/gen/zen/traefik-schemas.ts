// This file is auto-generated via `zen.StructToZodSchema`.
// Do not edit manually.

import { z } from "zod";

export const DomainSchema = z.object({
	main: z.string().optional(),
	sans: z.string().array().optional(),
});
export type Domain = z.infer<typeof DomainSchema>;

export const RouterTLSConfigSchema = z.object({
	options: z.string().optional(),
	certResolver: z.string().optional(),
	domains: DomainSchema.array().optional(),
});
export type RouterTLSConfig = z.infer<typeof RouterTLSConfigSchema>;

export const RouterObservabilityConfigSchema = z.object({
	accessLogs: z.boolean().optional(),
	tracing: z.boolean().optional(),
	metrics: z.boolean().optional(),
});
export type RouterObservabilityConfig = z.infer<
	typeof RouterObservabilityConfigSchema
>;

export const RouterSchema = z.object({
	entryPoints: z.string().array().optional(),
	middlewares: z.string().array().optional(),
	service: z.string().optional(),
	rule: z.string().optional(),
	ruleSyntax: z.string().optional(),
	priority: z.number().optional(),
	tls: RouterTLSConfigSchema.optional(),
	observability: RouterObservabilityConfigSchema.optional(),
});
export type Router = z.infer<typeof RouterSchema>;

export const RouterTCPTLSConfigSchema = z.object({
	passthrough: z.boolean(),
	options: z.string().optional(),
	certResolver: z.string().optional(),
	domains: DomainSchema.array().optional(),
});
export type RouterTCPTLSConfig = z.infer<typeof RouterTCPTLSConfigSchema>;

export const TCPRouterSchema = z.object({
	entryPoints: z.string().array().optional(),
	middlewares: z.string().array().optional(),
	service: z.string().optional(),
	rule: z.string().optional(),
	ruleSyntax: z.string().optional(),
	priority: z.number().optional(),
	tls: RouterTCPTLSConfigSchema.optional(),
});
export type TCPRouter = z.infer<typeof TCPRouterSchema>;

export const CookieSchema = z.object({
	name: z.string().optional(),
	secure: z.boolean().optional(),
	httpOnly: z.boolean().optional(),
	sameSite: z.string().optional(),
	maxAge: z.number().optional(),
	path: z.string().optional(),
	domain: z.string().optional(),
});
export type Cookie = z.infer<typeof CookieSchema>;

export const StickySchema = z.object({
	cookie: CookieSchema.optional(),
});
export type Sticky = z.infer<typeof StickySchema>;

export const ServerSchema = z.object({
	url: z.string().optional(),
	weight: z.number().optional(),
	preservePath: z.boolean().optional(),
	fenced: z.boolean().optional(),
});
export type Server = z.infer<typeof ServerSchema>;

export const ServerHealthCheckSchema = z.object({
	scheme: z.string().optional(),
	mode: z.string().optional(),
	path: z.string().optional(),
	method: z.string().optional(),
	status: z.number().optional(),
	port: z.number().optional(),
	interval: z.number().optional(),
	timeout: z.number().optional(),
	hostname: z.string().optional(),
	followRedirects: z.boolean().optional(),
	headers: z.record(z.string(), z.string()).optional(),
});
export type ServerHealthCheck = z.infer<typeof ServerHealthCheckSchema>;

export const ResponseForwardingSchema = z.object({
	flushInterval: z.number().optional(),
});
export type ResponseForwarding = z.infer<typeof ResponseForwardingSchema>;

export const ServersLoadBalancerSchema = z.object({
	sticky: StickySchema.optional(),
	servers: ServerSchema.array().optional(),
	strategy: z.string().optional(),
	healthCheck: ServerHealthCheckSchema.optional(),
	passHostHeader: z.boolean().nullable(),
	responseForwarding: ResponseForwardingSchema.optional(),
	serversTransport: z.string().optional(),
});
export type ServersLoadBalancer = z.infer<typeof ServersLoadBalancerSchema>;

export const WRRServiceSchema = z.object({
	name: z.string().optional(),
	weight: z.number().optional(),
});
export type WRRService = z.infer<typeof WRRServiceSchema>;

export const HealthCheckSchema = z.object({});
export type HealthCheck = z.infer<typeof HealthCheckSchema>;

export const WeightedRoundRobinSchema = z.object({
	services: WRRServiceSchema.array().optional(),
	sticky: StickySchema.optional(),
	healthCheck: HealthCheckSchema.optional(),
});
export type WeightedRoundRobin = z.infer<typeof WeightedRoundRobinSchema>;

export const MirrorServiceSchema = z.object({
	name: z.string().optional(),
	percent: z.number().optional(),
});
export type MirrorService = z.infer<typeof MirrorServiceSchema>;

export const MirroringSchema = z.object({
	service: z.string().optional(),
	mirrorBody: z.boolean().optional(),
	maxBodySize: z.number().optional(),
	mirrors: MirrorServiceSchema.array().optional(),
	healthCheck: HealthCheckSchema.optional(),
});
export type Mirroring = z.infer<typeof MirroringSchema>;

export const FailoverSchema = z.object({
	service: z.string().optional(),
	fallback: z.string().optional(),
	healthCheck: HealthCheckSchema.optional(),
});
export type Failover = z.infer<typeof FailoverSchema>;

export const ServiceSchema = z.object({
	loadBalancer: ServersLoadBalancerSchema.optional(),
	weighted: WeightedRoundRobinSchema.optional(),
	mirroring: MirroringSchema.optional(),
	failover: FailoverSchema.optional(),
});
export type Service = z.infer<typeof ServiceSchema>;

export const AddPrefixSchema = z.object({
	prefix: z.string().optional().describe("/prefix"),
});
export type AddPrefix = z.infer<typeof AddPrefixSchema>;

export const StripPrefixSchema = z.object({
	prefixes: z.string().array().optional(),
	forceSlash: z.boolean().optional(),
});
export type StripPrefix = z.infer<typeof StripPrefixSchema>;

export const StripPrefixRegexSchema = z.object({
	regex: z.string().array().optional(),
});
export type StripPrefixRegex = z.infer<typeof StripPrefixRegexSchema>;

export const ReplacePathSchema = z.object({
	path: z.string().optional(),
});
export type ReplacePath = z.infer<typeof ReplacePathSchema>;

export const ReplacePathRegexSchema = z.object({
	regex: z.string().optional(),
	replacement: z.string().optional(),
});
export type ReplacePathRegex = z.infer<typeof ReplacePathRegexSchema>;

export const ChainSchema = z.object({
	middlewares: z.string().array().optional(),
});
export type Chain = z.infer<typeof ChainSchema>;

export const IPStrategySchema = z.object({
	depth: z.number().optional(),
	excludedIPs: z.string().array().optional(),
	ipv6Subnet: z.number().optional(),
});
export type IPStrategy = z.infer<typeof IPStrategySchema>;

export const IPWhiteListSchema = z.object({
	sourceRange: z.string().array().optional(),
	ipStrategy: IPStrategySchema.optional(),
});
export type IPWhiteList = z.infer<typeof IPWhiteListSchema>;

export const IPAllowListSchema = z.object({
	sourceRange: z.string().array().optional(),
	ipStrategy: IPStrategySchema.optional(),
	rejectStatusCode: z.number().optional(),
});
export type IPAllowList = z.infer<typeof IPAllowListSchema>;

export const HeadersSchema = z.object({
	customRequestHeaders: z.record(z.string(), z.string()).optional(),
	customResponseHeaders: z.record(z.string(), z.string()).optional(),
	accessControlAllowCredentials: z.boolean().optional(),
	accessControlAllowHeaders: z.string().array().optional(),
	accessControlAllowMethods: z.string().array().optional(),
	accessControlAllowOriginList: z.string().array().optional(),
	accessControlAllowOriginListRegex: z.string().array().optional(),
	accessControlExposeHeaders: z.string().array().optional(),
	accessControlMaxAge: z.number().optional(),
	addVaryHeader: z.boolean().optional(),
	allowedHosts: z.string().array().optional(),
	hostsProxyHeaders: z.string().array().optional(),
	sslProxyHeaders: z.record(z.string(), z.string()).optional(),
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
	contentSecurityPolicyReportOnly: z.string().optional(),
	publicKey: z.string().optional(),
	referrerPolicy: z.string().optional(),
	permissionsPolicy: z.string().optional(),
	isDevelopment: z.boolean().optional(),
	featurePolicy: z.string().optional(),
	sslRedirect: z.boolean().optional(),
	sslTemporaryRedirect: z.boolean().optional(),
	sslHost: z.string().optional(),
	sslForceHost: z.boolean().optional(),
});
export type Headers = z.infer<typeof HeadersSchema>;

export const ErrorPageSchema = z.object({
	status: z.string().array().optional(),
	statusRewrites: z.record(z.string(), z.number()).optional(),
	service: z.string().optional(),
	query: z.string().optional(),
});
export type ErrorPage = z.infer<typeof ErrorPageSchema>;

export const SourceCriterionSchema = z.object({
	ipStrategy: IPStrategySchema.optional(),
	requestHeaderName: z.string().optional(),
	requestHost: z.boolean().optional(),
});
export type SourceCriterion = z.infer<typeof SourceCriterionSchema>;

export const ClientTLSSchema = z.object({
	ca: z.string().optional(),
	cert: z.string().optional(),
	key: z.string().optional(),
	insecureSkipVerify: z.boolean().optional(),
});
export type ClientTLS = z.infer<typeof ClientTLSSchema>;

export const RedisSchema = z.object({
	endpoints: z.string().array().optional(),
	tls: ClientTLSSchema.optional(),
	username: z.string().optional(),
	password: z.string().optional(),
	db: z.number().optional(),
	poolSize: z.number().optional(),
	minIdleConns: z.number().optional(),
	maxActiveConns: z.number().optional(),
	readTimeout: z.number().optional(),
	writeTimeout: z.number().optional(),
	dialTimeout: z.number().optional(),
});
export type Redis = z.infer<typeof RedisSchema>;

export const RateLimitSchema = z.object({
	average: z.number().optional(),
	period: z.number().optional(),
	burst: z.number().optional(),
	sourceCriterion: SourceCriterionSchema.optional(),
	redis: RedisSchema.optional(),
});
export type RateLimit = z.infer<typeof RateLimitSchema>;

export const RedirectRegexSchema = z.object({
	regex: z.string().optional(),
	replacement: z.string().optional(),
	permanent: z.boolean().optional(),
});
export type RedirectRegex = z.infer<typeof RedirectRegexSchema>;

export const RedirectSchemeSchema = z.object({
	scheme: z.string().optional(),
	port: z.string().optional(),
	permanent: z.boolean().optional(),
});
export type RedirectScheme = z.infer<typeof RedirectSchemeSchema>;

export const BasicAuthSchema = z.object({
	users: z
		.string()
		.array()
		.optional()
		.describe(
			"List of users in the format `user:password` (will be encrypted)",
		),

	usersFile: z.string().optional().describe("/etc/traefik/usersfile"),
	realm: z.string().optional().describe("Traefik Basic Auth"),
	removeHeader: z.boolean().optional().describe("Remove Authorization header"),
	headerField: z
		.string()
		.optional()
		.describe("Custom header name (default: Authorization)"),
});
export type BasicAuth = z.infer<typeof BasicAuthSchema>;

export const DigestAuthSchema = z.object({
	users: z
		.string()
		.array()
		.optional()
		.describe("user:realm:password (will be encrypted)"),
	usersFile: z.string().optional().describe("/etc/traefik/users.digest"),
	removeHeader: z.boolean().optional().describe("Remove Authorization header"),
	realm: z.string().optional().describe("Traefik Digest Auth"),
	headerField: z
		.string()
		.optional()
		.describe("Custom header name (default: Authorization)"),
});

export type DigestAuth = z.infer<typeof DigestAuthSchema>;

export const ForwardAuthSchema = z.object({
	address: z.string().optional(),
	tls: ClientTLSSchema.optional(),
	trustForwardHeader: z.boolean().optional(),
	authResponseHeaders: z.string().array().optional(),
	authResponseHeadersRegex: z.string().optional(),
	authRequestHeaders: z.string().array().optional(),
	addAuthCookiesToResponse: z.string().array().optional(),
	headerField: z.string().optional(),
	forwardBody: z.boolean().optional(),
	maxBodySize: z.number().optional(),
	preserveLocationHeader: z.boolean().optional(),
	preserveRequestMethod: z.boolean().optional(),
});
export type ForwardAuth = z.infer<typeof ForwardAuthSchema>;

export const InFlightReqSchema = z.object({
	amount: z.number().optional(),
	sourceCriterion: SourceCriterionSchema.optional(),
});
export type InFlightReq = z.infer<typeof InFlightReqSchema>;

export const BufferingSchema = z.object({
	maxRequestBodyBytes: z.number().optional(),
	memRequestBodyBytes: z.number().optional(),
	maxResponseBodyBytes: z.number().optional(),
	memResponseBodyBytes: z.number().optional(),
	retryExpression: z.string().optional(),
});
export type Buffering = z.infer<typeof BufferingSchema>;

export const CircuitBreakerSchema = z.object({
	expression: z.string().optional(),
	checkPeriod: z.number().optional(),
	fallbackDuration: z.number().optional(),
	recoveryDuration: z.number().optional(),
	responseCode: z.number().optional(),
});
export type CircuitBreaker = z.infer<typeof CircuitBreakerSchema>;

export const CompressSchema = z.object({
	excludedContentTypes: z.string().array().optional(),
	includedContentTypes: z.string().array().optional(),
	minResponseBodyBytes: z.number().optional(),
	encodings: z.string().array().optional(),
	defaultEncoding: z.string().optional(),
});
export type Compress = z.infer<typeof CompressSchema>;

export const TLSClientCertificateSubjectDNInfoSchema = z.object({
	country: z.boolean().optional(),
	province: z.boolean().optional(),
	locality: z.boolean().optional(),
	organization: z.boolean().optional(),
	organizationalUnit: z.boolean().optional(),
	commonName: z.boolean().optional(),
	serialNumber: z.boolean().optional(),
	domainComponent: z.boolean().optional(),
});
export type TLSClientCertificateSubjectDNInfo = z.infer<
	typeof TLSClientCertificateSubjectDNInfoSchema
>;

export const TLSClientCertificateIssuerDNInfoSchema = z.object({
	country: z.boolean().optional(),
	province: z.boolean().optional(),
	locality: z.boolean().optional(),
	organization: z.boolean().optional(),
	commonName: z.boolean().optional(),
	serialNumber: z.boolean().optional(),
	domainComponent: z.boolean().optional(),
});
export type TLSClientCertificateIssuerDNInfo = z.infer<
	typeof TLSClientCertificateIssuerDNInfoSchema
>;

export const TLSClientCertificateInfoSchema = z.object({
	notAfter: z.boolean().optional(),
	notBefore: z.boolean().optional(),
	sans: z.boolean().optional(),
	serialNumber: z.boolean().optional(),
	subject: TLSClientCertificateSubjectDNInfoSchema.optional(),
	issuer: TLSClientCertificateIssuerDNInfoSchema.optional(),
});
export type TLSClientCertificateInfo = z.infer<
	typeof TLSClientCertificateInfoSchema
>;

export const PassTLSClientCertSchema = z.object({
	pem: z.boolean().optional(),
	info: TLSClientCertificateInfoSchema.optional(),
});
export type PassTLSClientCert = z.infer<typeof PassTLSClientCertSchema>;

export const RetrySchema = z.object({
	attempts: z.number().optional(),
	initialInterval: z.number().optional(),
});
export type Retry = z.infer<typeof RetrySchema>;

export const ContentTypeSchema = z.object({
	autoDetect: z.boolean().optional(),
});
export type ContentType = z.infer<typeof ContentTypeSchema>;

export const GrpcWebSchema = z.object({
	allowOrigins: z.string().array().optional(),
});
export type GrpcWeb = z.infer<typeof GrpcWebSchema>;

export const HeaderModifierSchema = z.object({
	set: z.record(z.string(), z.string()).optional(),
	add: z.record(z.string(), z.string()).optional(),
	remove: z.string().array().optional(),
});
export type HeaderModifier = z.infer<typeof HeaderModifierSchema>;

export const RequestRedirectSchema = z.object({
	scheme: z.string().optional(),
	hostname: z.string().optional(),
	port: z.string().optional(),
	path: z.string().optional(),
	pathPrefix: z.string().optional(),
	statusCode: z.number().optional(),
});
export type RequestRedirect = z.infer<typeof RequestRedirectSchema>;

export const URLRewriteSchema = z.object({
	hostname: z.string().optional(),
	path: z.string().optional(),
	pathPrefix: z.string().optional(),
});
export type URLRewrite = z.infer<typeof URLRewriteSchema>;

export const MiddlewareSchema = z.object({
	addPrefix: AddPrefixSchema.optional(),
	stripPrefix: StripPrefixSchema.optional(),
	stripPrefixRegex: StripPrefixRegexSchema.optional(),
	replacePath: ReplacePathSchema.optional(),
	replacePathRegex: ReplacePathRegexSchema.optional(),
	chain: ChainSchema.optional(),
	ipWhiteList: IPWhiteListSchema.optional(),
	ipAllowList: IPAllowListSchema.optional(),
	headers: HeadersSchema.optional(),
	errors: ErrorPageSchema.optional(),
	rateLimit: RateLimitSchema.optional(),
	redirectRegex: RedirectRegexSchema.optional(),
	redirectScheme: RedirectSchemeSchema.optional(),
	basicAuth: BasicAuthSchema.optional(),
	digestAuth: DigestAuthSchema.optional(),
	forwardAuth: ForwardAuthSchema.optional(),
	inFlightReq: InFlightReqSchema.optional(),
	buffering: BufferingSchema.optional(),
	circuitBreaker: CircuitBreakerSchema.optional(),
	compress: CompressSchema.optional(),
	passTLSClientCert: PassTLSClientCertSchema.optional(),
	retry: RetrySchema.optional(),
	contentType: ContentTypeSchema.optional(),
	grpcWeb: GrpcWebSchema.optional(),
	plugin: z.record(z.string(), z.record(z.string(), z.any())).optional(),
	requestHeaderModifier: HeaderModifierSchema.optional(),
	responseHeaderModifier: HeaderModifierSchema.optional(),
	requestRedirect: RequestRedirectSchema.optional(),
	URLRewrite: URLRewriteSchema.optional(),
});
export type Middleware = z.infer<typeof MiddlewareSchema>;

export const TCPInFlightConnSchema = z.object({
	amount: z.number().optional(),
});
export type TCPInFlightConn = z.infer<typeof TCPInFlightConnSchema>;

export const TCPIPWhiteListSchema = z.object({
	sourceRange: z.string().array().optional(),
});
export type TCPIPWhiteList = z.infer<typeof TCPIPWhiteListSchema>;

export const TCPIPAllowListSchema = z.object({
	sourceRange: z.string().array().optional(),
});
export type TCPIPAllowList = z.infer<typeof TCPIPAllowListSchema>;

export const TCPMiddlewareSchema = z.object({
	inFlightConn: TCPInFlightConnSchema.optional(),
	ipWhiteList: TCPIPWhiteListSchema.optional(),
	ipAllowList: TCPIPAllowListSchema.optional(),
});
export type TCPMiddleware = z.infer<typeof TCPMiddlewareSchema>;

export const CertificateSchema = z.object({
	certFile: z.string().optional(),
	keyFile: z.string().optional(),
});
export type Certificate = z.infer<typeof CertificateSchema>;

export const ForwardingTimeoutsSchema = z.object({
	dialTimeout: z.string().optional(),
	responseHeaderTimeout: z.string().optional(),
	idleConnTimeout: z.string().optional(),
	readIdleTimeout: z.string().optional(),
	pingTimeout: z.string().optional(),
});
export type ForwardingTimeouts = z.infer<typeof ForwardingTimeoutsSchema>;

export const SpiffeSchema = z.object({
	ids: z.string().array().optional(),
	trustDomain: z.string().optional(),
});
export type Spiffe = z.infer<typeof SpiffeSchema>;

export const ServersTransportSchema = z.object({
	serverName: z.string().optional(),
	insecureSkipVerify: z.boolean().optional(),
	rootCAs: z.string().array().optional(),
	certificates: CertificateSchema.array().optional(),
	maxIdleConnsPerHost: z.number().optional(),
	forwardingTimeouts: ForwardingTimeoutsSchema.optional(),
	disableHTTP2: z.boolean().optional(),
	peerCertURI: z.string().optional(),
	spiffe: SpiffeSchema.optional(),
});
export type ServersTransport = z.infer<typeof ServersTransportSchema>;

export const UDPRouterSchema = z.object({
	entryPoints: z.string().array().optional(),
	service: z.string().optional(),
});
export type UDPRouter = z.infer<typeof UDPRouterSchema>;

export const ProxyProtocolSchema = z.object({
	version: z.number().optional(),
});
export type ProxyProtocol = z.infer<typeof ProxyProtocolSchema>;

export const TCPServerSchema = z.object({
	address: z.string().optional(),
	tls: z.boolean().optional(),
});
export type TCPServer = z.infer<typeof TCPServerSchema>;

export const TCPServersLoadBalancerSchema = z.object({
	proxyProtocol: ProxyProtocolSchema.optional(),
	servers: TCPServerSchema.array().optional(),
	serversTransport: z.string().optional(),
	terminationDelay: z.number().optional(),
});
export type TCPServersLoadBalancer = z.infer<
	typeof TCPServersLoadBalancerSchema
>;

export const TCPWRRServiceSchema = z.object({
	name: z.string().optional(),
	weight: z.number().optional(),
});
export type TCPWRRService = z.infer<typeof TCPWRRServiceSchema>;

export const TCPWeightedRoundRobinSchema = z.object({
	services: TCPWRRServiceSchema.array().optional(),
});
export type TCPWeightedRoundRobin = z.infer<typeof TCPWeightedRoundRobinSchema>;

export const TCPServiceSchema = z.object({
	loadBalancer: TCPServersLoadBalancerSchema.optional(),
	weighted: TCPWeightedRoundRobinSchema.optional(),
});
export type TCPService = z.infer<typeof TCPServiceSchema>;

export const UDPServerSchema = z.object({
	address: z.string().optional(),
});
export type UDPServer = z.infer<typeof UDPServerSchema>;

export const UDPServersLoadBalancerSchema = z.object({
	servers: UDPServerSchema.array().optional(),
});
export type UDPServersLoadBalancer = z.infer<
	typeof UDPServersLoadBalancerSchema
>;

export const UDPWRRServiceSchema = z.object({
	name: z.string().optional(),
	weight: z.number().optional(),
});
export type UDPWRRService = z.infer<typeof UDPWRRServiceSchema>;

export const UDPWeightedRoundRobinSchema = z.object({
	services: UDPWRRServiceSchema.array().optional(),
});
export type UDPWeightedRoundRobin = z.infer<typeof UDPWeightedRoundRobinSchema>;

export const UDPServiceSchema = z.object({
	loadBalancer: UDPServersLoadBalancerSchema.optional(),
	weighted: UDPWeightedRoundRobinSchema.optional(),
});
export type UDPService = z.infer<typeof UDPServiceSchema>;

export const TLSClientConfigSchema = z.object({
	serverName: z.string().optional(),
	insecureSkipVerify: z.boolean().optional(),
	rootCAs: z.string().array().optional(),
	certificates: CertificateSchema.array().optional(),
	peerCertURI: z.string().optional(),
	spiffe: SpiffeSchema.optional(),
});
export type TLSClientConfig = z.infer<typeof TLSClientConfigSchema>;

export const TCPServersTransportSchema = z.object({
	dialKeepAlive: z.string().optional(),
	dialTimeout: z.string().optional(),
	terminationDelay: z.string().optional(),
	tls: TLSClientConfigSchema.optional(),
});
export type TCPServersTransport = z.infer<typeof TCPServersTransportSchema>;
