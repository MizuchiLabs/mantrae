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
	PluginSchema,
	type HTTPMiddleware,
	type TCPMiddleware
} from '$lib/types/middlewares';

export type ZodObjectOrRecord = AnyZodObject | z.ZodRecord<any, any>;

// Split middleware types into HTTP and TCP
export type SupportedMiddlewareHTTP = keyof Omit<HTTPMiddleware, 'name' | 'protocol' | 'type'>;

export type SupportedMiddlewareTCP = keyof Omit<TCPMiddleware, 'name' | 'protocol' | 'type'>;

export type SupportedMiddleware = SupportedMiddlewareHTTP | SupportedMiddlewareTCP;

export const HTTPMiddlewareSchemaMap: Record<SupportedMiddlewareHTTP, ZodObjectOrRecord> = {
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
	plugin: PluginSchema
};
export const TCPMiddlewareSchemaMap: Record<SupportedMiddlewareTCP, ZodObjectOrRecord> = {
	ipAllowList: TCPIPAllowListSchema,
	inFlightConn: TCPInFlightConnSchema
};

// Combined schema map
export const MiddlewareSchemaMap: Record<SupportedMiddleware, ZodObjectOrRecord> = {
	...HTTPMiddlewareSchemaMap,
	...TCPMiddlewareSchemaMap
};

export const GetSchema = (type: SupportedMiddleware | undefined) => {
	if (!type) return z.object({});
	return MiddlewareSchemaMap[type];
};

// Split middleware types for UI
export const HTTPMiddlewareTypes = Object.keys(HTTPMiddlewareSchemaMap).map((key) => ({
	value: key,
	label: key
		.split(/(?=[A-Z])/)
		.map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
		.join(' ')
}));

export const TCPMiddlewareTypes = Object.keys(TCPMiddlewareSchemaMap).map((key) => ({
	value: key,
	label: key
		.split(/(?=[A-Z])/)
		.map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
		.join(' ')
}));

export const MiddlewareTypes = [...HTTPMiddlewareTypes, ...TCPMiddlewareTypes];
