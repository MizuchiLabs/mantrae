import { DnsProviderType } from "./gen/mantrae/v1/dns_provider_pb";
import { MiddlewareType } from "./gen/mantrae/v1/middleware_pb";
import { RouterType } from "./gen/mantrae/v1/router_pb";
import { ServiceType } from "./gen/mantrae/v1/service_pb";
import type { JsonObject } from "@bufbuild/protobuf";

// Parse protobuf config
export function unmarshalConfig<T>(json: JsonObject | undefined): T {
	if (!json) return {} as T;
	const str = JSON.stringify(json);
	return JSON.parse(str) as T;
}
export function marshalConfig<T>(config: T): JsonObject {
	const str = JSON.stringify(config);
	return JSON.parse(str) as JsonObject;
}

// Convert enum to select options
export const routerTypes = Object.keys(RouterType)
	.filter((key) => isNaN(Number(key)) && key !== "UNSPECIFIED")
	.map((key) => ({
		label: key.toUpperCase(),
		value: RouterType[key as keyof typeof RouterType],
	}));

export const serviceTypes = Object.keys(ServiceType)
	.filter((key) => isNaN(Number(key)) && key !== "UNSPECIFIED")
	.map((key) => ({
		label: key.toUpperCase(),
		value: ServiceType[key as keyof typeof ServiceType],
	}));

export const middlewareTypes = Object.keys(MiddlewareType)
	.filter((key) => isNaN(Number(key)) && key !== "UNSPECIFIED")
	.map((key) => ({
		label: key.toUpperCase(),
		value: MiddlewareType[key as keyof typeof MiddlewareType],
	}));

export const dnsProviderTypes = Object.keys(DnsProviderType)
	.filter((key) => isNaN(Number(key)) && key !== "UNSPECIFIED")
	.map((key) => ({
		label: key
			.replace("DNS_PROVIDER_TYPE_", "")
			.toLowerCase()
			.replace(/^\w/, (c) => c.toUpperCase()),
		value: DnsProviderType[key as keyof typeof DnsProviderType],
	}));

export const HTTPMiddlewareKeys = [
	{ value: "addPrefix", label: "Add Prefix" },
	{ value: "basicAuth", label: "Basic Auth" },
	{ value: "digestAuth", label: "Digest Auth" },
	{ value: "buffering", label: "Buffering" },
	{ value: "chain", label: "Chain" },
	{ value: "circuitBreaker", label: "Circuit Breaker" },
	{ value: "compress", label: "Compress" },
	{ value: "errorPage", label: "Error Page" },
	{ value: "forwardAuth", label: "Forward Auth" },
	{ value: "headers", label: "Headers" },
	{ value: "ipAllowList", label: "IP Allow List" },
	{ value: "inFlightReq", label: "In-Flight Request" },
	{ value: "passTLSClientCert", label: "Pass TLS Client Cert" },
	{ value: "rateLimit", label: "Rate Limit" },
	{ value: "redirectRegex", label: "Redirect Regex" },
	{ value: "redirectScheme", label: "Redirect Scheme" },
	{ value: "replacePath", label: "Replace Path" },
	{ value: "replacePathRegex", label: "Replace Path Regex" },
	{ value: "retry", label: "Retry" },
	{ value: "stripPrefix", label: "Strip Prefix" },
	{ value: "stripPrefixRegex", label: "Strip Prefix Regex" },
	{ value: "plugin", label: "Plugin" },
];

export const TCPMiddlewareKeys = [
	{ value: "ipAllowList", label: "IP Allow List" },
	{ value: "inFlightConn", label: "In-Flight Connection" },
];
