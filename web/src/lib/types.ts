import type { JsonObject } from '@bufbuild/protobuf';
import type { IconProps } from '@lucide/svelte';
import type { Component } from 'svelte';
import { DNSProviderType } from './gen/mantrae/v1/dns_provider_pb';
import { ProtocolType } from './gen/mantrae/v1/protocol_pb';

export type IconComponent = Component<IconProps, Record<string, never>, ''>;

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
export const protocolTypes = Object.keys(ProtocolType)
	.filter((key) => isNaN(Number(key)) && key !== 'UNSPECIFIED')
	.map((key) => ({
		label: key.toUpperCase(),
		value: ProtocolType[key as keyof typeof ProtocolType]
	}));

export const dnsProviderTypes = Object.keys(DNSProviderType)
	.filter((key) => isNaN(Number(key)) && key !== 'UNSPECIFIED')
	.map((key) => ({
		label: key
			.replace('DNS_PROVIDER_TYPE_', '')
			.toLowerCase()
			.replace(/^\w/, (c) => c.toUpperCase()),
		value: DNSProviderType[key as keyof typeof DNSProviderType]
	}));
