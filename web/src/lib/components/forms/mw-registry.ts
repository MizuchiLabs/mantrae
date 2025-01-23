import type { AddPrefix, BasicAuth } from '$lib/types/middlewares';

type FieldType = 'text' | 'number' | 'boolean' | 'array' | 'object' | 'select';

interface FieldDefinition {
	type: FieldType;
	label: string;
	description: string;
}

interface MiddlewareConfig {
	label: string;
	description: string;
	protocol: 'http' | 'tcp';
	fields: Record<string, FieldDefinition>;
}

function toTitleCase(str: string): string {
	return str.charAt(0).toUpperCase() + str.slice(1);
}

function inferFieldType(value: unknown): FieldType {
	if (Array.isArray(value)) return 'array';
	if (typeof value === 'boolean') return 'boolean';
	if (typeof value === 'number') return 'number';
	if (typeof value === 'object' && value !== null) return 'object';
	return 'text';
}

function createFieldDefinition(key: string, type: FieldType): FieldDefinition {
	return {
		type,
		label: toTitleCase(key),
		description: `${toTitleCase(key)} for the middleware`
	};
}

function generateFields<T extends object>(type: T): Record<string, FieldDefinition> {
	const fields: Record<string, FieldDefinition> = {};

	Object.keys(type as object).forEach((key) => {
		const fieldType = inferFieldType((type as Record<string, unknown>)[key]);
		fields[key] = createFieldDefinition(key, fieldType);
	});

	return fields;
}

type SupportedMiddleware = 'addPrefix' | 'basicAuth';
export const MIDDLEWARE_REGISTRY: Record<SupportedMiddleware, MiddlewareConfig> = {
	addPrefix: {
		label: 'Add Prefix',
		description: 'Adds a path prefix to the request URL',
		protocol: 'http',
		fields: generateFields<AddPrefix>({} as AddPrefix)
	},
	basicAuth: {
		label: 'Basic Authentication',
		description: 'Adds Basic Authentication to your services',
		protocol: 'http',
		fields: generateFields<BasicAuth>({} as BasicAuth)
	}
};

export function getMiddlewareConfig(type: SupportedMiddleware): MiddlewareConfig {
	return MIDDLEWARE_REGISTRY[type];
}
