import type { ZodSchema, ZodTypeAny } from 'zod';

export interface FormField {
	type: 'string' | 'number' | 'boolean' | 'array' | 'object' | 'record' | 'plugin';
	label: string;
	optional: boolean;
	description?: string;
	arrayItemType?: FormField;
	nestedSchema?: ZodSchema;
	recordValueType?: FormField;
}

export function extractSchemaFields(schema: ZodSchema): Record<string, FormField> {
	let actualSchema = schema;

	// Unwrap ZodOptional if present
	if ((schema as any)._def.typeName === 'ZodOptional') {
		actualSchema = (schema as any)._def.innerType;
	}

	// Handle the case where the entire schema is a ZodRecord (like for plugins)
	const schemaType = (actualSchema as any)._def.typeName;

	if (schemaType === 'ZodRecord') {
		// For ZodRecord schemas, create a single plugin field
		return {
			plugin: {
				type: 'plugin',
				label: 'Plugin Configuration',
				optional: false,
				description: 'YAML configuration for plugins'
			}
		};
	}

	// Check if it's a ZodObject
	if ((actualSchema as any)._def.typeName !== 'ZodObject') {
		console.error('Expected ZodObject, got:', (actualSchema as any)._def.typeName);
		return {};
	}

	const shape = (actualSchema as any)._def.shape();
	const fields: Record<string, FormField> = {};

	for (const [key, zodType] of Object.entries(shape)) {
		fields[key] = parseZodType(zodType as ZodTypeAny, key);
	}

	return fields;
}

function parseZodType(zodType: ZodTypeAny, key: string): FormField {
	const def = zodType._def;
	let type = def.typeName;
	let optional = false;
	let currentType = zodType;

	// Handle optional fields
	if (type === 'ZodOptional') {
		optional = true;
		currentType = def.innerType;
		type = currentType._def.typeName;
	}

	// Special handling for plugin field - handle it as plugin type regardless of Zod type
	if (key === 'plugin') {
		return {
			type: 'plugin',
			label: formatLabel(key),
			optional,
			description: 'YAML or JSON configuration for plugins'
		};
	}

	// Map Zod types to form field types
	switch (type) {
		case 'ZodString':
			return {
				type: 'string',
				label: formatLabel(key),
				optional,
				description: def.description
			};
		case 'ZodNumber':
			return {
				type: 'number',
				label: formatLabel(key),
				optional,
				description: def.description
			};
		case 'ZodBoolean':
			return {
				type: 'boolean',
				label: formatLabel(key),
				optional,
				description: def.description
			};
		case 'ZodArray':
			return {
				type: 'array',
				label: formatLabel(key),
				optional,
				description: def.description,
				arrayItemType: parseZodType(currentType._def.type, `${key}Item`)
			};
		case 'ZodObject':
			return {
				type: 'object',
				label: formatLabel(key),
				optional,
				description: def.description,
				nestedSchema: currentType
			};
		case 'ZodRecord':
			return {
				type: 'record',
				label: formatLabel(key),
				optional,
				description: def.description,
				recordValueType: parseZodType(currentType._def.valueType, `${key}Value`)
			};
		default:
			return {
				type: 'string',
				label: formatLabel(key),
				optional,
				description: def.description
			};
	}
}

function formatLabel(key: string): string {
	return key.replace(/([A-Z])/g, ' $1').replace(/^./, (s) => s.toUpperCase());
}
