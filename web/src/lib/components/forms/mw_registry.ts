import { z, type AnyZodObject } from 'zod';
import {
	AddPrefixSchema,
	BasicAuthSchema,
	DigestAuthSchema,
	BufferingSchema,
	ChainSchema,
	CircuitBreakerSchema,
	CompressSchema,
	ErrorPageSchema,
	ForwardAuthSchema,
	HeadersSchema,
	IPAllowListSchema,
	InFlightReqSchema,
	PassTLSClientCertSchema,
	RateLimitSchema,
	RedirectRegexSchema,
	RedirectSchemeSchema,
	ReplacePathSchema,
	ReplacePathRegexSchema,
	RetrySchema,
	StripPrefixSchema,
	StripPrefixRegexSchema,
	TCPIPAllowListSchema,
	TCPInFlightConnSchema,
	type Middleware
} from '$lib/types/middlewares';

// Create a mapping of SupportedMiddleware keys to their corresponding schemas
export const MiddlewareSchemaMap: Record<SupportedMiddleware, AnyZodObject> = {
	addPrefix: AddPrefixSchema,
	basicAuth: BasicAuthSchema,
	digestAuth: DigestAuthSchema,
	buffering: BufferingSchema,
	chain: ChainSchema,
	circuitBreaker: CircuitBreakerSchema,
	compress: CompressSchema,
	errorPage: ErrorPageSchema,
	forwardAuth: ForwardAuthSchema,
	headers: HeadersSchema,
	ipAllowList: IPAllowListSchema,
	inFlightReq: InFlightReqSchema,
	passTLSClientCert: PassTLSClientCertSchema,
	rateLimit: RateLimitSchema,
	redirectRegex: RedirectRegexSchema,
	redirectScheme: RedirectSchemeSchema,
	replacePath: ReplacePathSchema,
	replacePathRegex: ReplacePathRegexSchema,
	retry: RetrySchema,
	stripPrefix: StripPrefixSchema,
	stripPrefixRegex: StripPrefixRegexSchema,
	tcpIpAllowList: TCPIPAllowListSchema,
	tcpInFlightConn: TCPInFlightConnSchema,
	plugin: z.object({ plugin: z.string() }) // Generic handling for plugins
};
export const GetSchema = (type: SupportedMiddleware | undefined) => {
	if (!type) return z.object({});
	return MiddlewareSchemaMap[type as SupportedMiddleware];
};

// Type definition for SupportedMiddleware
export type SupportedMiddleware = keyof Omit<Middleware, 'name' | 'protocol' | 'type'>;
export const MiddlewareTypes = Object.keys(MiddlewareSchemaMap).map((key) => ({
	value: key,
	label: key
		.split(/(?=[A-Z])/) // Split on capital letters
		.map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
		.join(' ')
}));
