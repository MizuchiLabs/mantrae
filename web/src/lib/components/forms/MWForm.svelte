<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { GetSchema } from '$lib/components/forms/mw_registry';
	import type { SupportedMiddleware } from '$lib/components/forms/mw_registry';
	import { z } from 'zod';
	import { Plus, Trash } from 'lucide-svelte';
	import Self from './MWForm.svelte';

	interface Props {
		type: SupportedMiddleware;
		values: Record<string, unknown>;
		onChange: (values: Record<string, unknown>) => void;
		schema?: z.ZodType;
	}

	let { type, values, onChange, schema }: Props = $props();

	const currentSchema = schema || GetSchema(type);

	type FieldConfig = {
		type: 'string' | 'number' | 'boolean' | 'array' | 'object';
		key: string;
		path: string;
		optional?: boolean;
		fields?: FieldConfig[];
		itemSchema?: z.ZodTypeAny;
	};

	function getBaseType(schema: z.ZodTypeAny): z.ZodTypeAny {
		if (schema instanceof z.ZodOptional) {
			return getBaseType(schema.unwrap());
		}
		if (schema instanceof z.ZodNullable) {
			return getBaseType(schema.unwrap());
		}
		if (schema instanceof z.ZodDefault) {
			return getBaseType(schema._def.innerType);
		}
		return schema;
	}

	function generateFormFields(schema: z.ZodType, path = ''): FieldConfig[] {
		if (!(schema instanceof z.ZodObject)) return [];

		return Object.entries(schema.shape).map(([key, value]): FieldConfig => {
			const fullPath = path ? `${path}.${key}` : key;
			const baseType = getBaseType(value);
			const optional = value instanceof z.ZodOptional;

			if (baseType instanceof z.ZodArray) {
				return {
					type: 'array',
					key,
					path: fullPath,
					optional,
					itemSchema: baseType.element
				};
			}

			if (baseType instanceof z.ZodObject) {
				return {
					type: 'object',
					key,
					path: fullPath,
					optional,
					fields: generateFormFields(baseType, fullPath),
					itemSchema: baseType
				};
			}

			if (baseType instanceof z.ZodBoolean) {
				return {
					type: 'boolean',
					key,
					path: fullPath,
					optional
				};
			}

			if (baseType instanceof z.ZodNumber) {
				return {
					type: 'number',
					key,
					path: fullPath,
					optional
				};
			}

			return {
				type: 'string',
				key,
				path: fullPath,
				optional
			};
		});
	}
	$effect(() => {
		// Initialize empty values when type changes
		if (type && Object.keys(values).length === 0) {
			onChange({});
		}
	});
	function getFieldValue(field: FieldConfig): unknown {
		const value = values[field.key];

		if (value === undefined) {
			if (field.type === 'array') return [];
			if (field.type === 'boolean') return false;
			if (field.type === 'number') return 0;
			if (field.type === 'object') return {};
			return '';
		}

		return value;
	}

	function updateNestedValue(obj: Record<string, unknown>, path: string[], value: unknown): void {
		const [key, ...rest] = path;
		if (rest.length === 0) {
			obj[key] = value;
		} else {
			obj[key] = obj[key] || {};
			if (typeof obj[key] === 'object' && obj[key] !== null) {
				updateNestedValue(obj[key] as Record<string, unknown>, rest, value);
			}
		}
	}

	function handleValueChange(path: string, value: unknown): void {
		const pathParts = path.split('.');
		const newValues = { ...values };
		updateNestedValue(newValues, pathParts, value);
		onChange(newValues);
	}

	function formatLabel(key: string): string {
		return key
			.split(/(?=[A-Z])/)
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
			.join(' ');
	}
</script>

<div class="space-y-4">
	{#if currentSchema instanceof z.ZodObject}
		{#each generateFormFields(currentSchema) as field}
			{#if field.type === 'object'}
				<div class="rounded-lg border p-4">
					<Label class="mb-2 block font-semibold">{formatLabel(field.key)}</Label>
					<Self
						{type}
						values={(values[field.key] as Record<string, unknown>) || {}}
						onChange={(newValues) => handleValueChange(field.path, newValues)}
					/>
				</div>
			{:else if field.type === 'string'}
				<div class="flex flex-col gap-2">
					<Label for={field.path}>{formatLabel(field.key)}</Label>
					<Input
						type="text"
						id={field.path}
						name={field.path}
						value={getFieldValue(field) as string}
						oninput={(e) => handleValueChange(field.path, e.currentTarget.value)}
					/>
				</div>
			{:else if field.type === 'number'}
				<div class="flex flex-col gap-2">
					<Label for={field.path}>{formatLabel(field.key)}</Label>
					<Input
						type="number"
						id={field.path}
						name={field.path}
						value={getFieldValue(field) as number}
						oninput={(e) => handleValueChange(field.path, Number(e.currentTarget.value))}
					/>
				</div>
			{:else if field.type === 'boolean'}
				<div class="flex items-center gap-2">
					<Switch
						id={field.path}
						name={field.path}
						checked={getFieldValue(field) as boolean}
						onCheckedChange={(checked) => handleValueChange(field.path, checked)}
					/>
					<Label for={field.path}>{formatLabel(field.key)}</Label>
				</div>
			{:else if field.type === 'array'}
				<div class="flex flex-col gap-2">
					<Label>{formatLabel(field.key)}</Label>
					{#each (values[field.key] as unknown[]) || [] as item, i}
						<div class="flex gap-2">
							<Input
								type="text"
								value={item as string}
								oninput={(e) => {
									const newArray = [...((values[field.key] as unknown[]) || [])];
									newArray[i] = e.currentTarget.value;
									handleValueChange(field.path, newArray);
								}}
							/>
							<Button
								variant="ghost"
								size="icon"
								class="text-red-500 hover:text-red-600"
								onclick={() => {
									const newArray = (values[field.key] as unknown[]).filter(
										(_, index) => index !== i
									);
									handleValueChange(field.path, newArray);
								}}
							>
								<Trash />
							</Button>
						</div>
					{/each}
					<Button
						variant="secondary"
						class="flex items-center gap-2"
						onclick={() => {
							const newArray = [...(getFieldValue(field) as unknown[]), ''];
							handleValueChange(field.path, newArray);
						}}
					>
						<Plus />
						Add {field.key}
					</Button>
				</div>
			{/if}
		{/each}
	{/if}
</div>
