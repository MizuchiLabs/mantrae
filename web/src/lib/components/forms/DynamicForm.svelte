<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import { MiddlewareSchema, TCPMiddlewareSchema } from '$lib/gen/zen/traefik-schemas';
	import type { JsonObject, JsonValue } from '@bufbuild/protobuf';
	import DynamicForm from './DynamicForm.svelte';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import { Plus, Trash } from '@lucide/svelte';
	import Button from '../ui/button/button.svelte';
	import YAML from 'yaml';
	import z from 'zod';

	interface Props {
		middlewareType: string;
		protocol: 'http' | 'tcp';
		data: JsonObject;
		onUpdate?: (data: JsonObject) => void;
	}

	let { middlewareType, protocol, data, onUpdate }: Props = $props();

	const schema = $derived.by(() => {
		const schemaMap = protocol === 'http' ? MiddlewareSchema.shape : TCPMiddlewareSchema.shape;
		const selectedSchema = schemaMap[middlewareType as keyof typeof schemaMap];
		return selectedSchema ? z.toJSONSchema(selectedSchema) : null;
	});

	const fields = $derived.by(() => {
		if (!schema?.properties) return [];
		return Object.entries(schema.properties).map(([fieldKey, fieldSchema]) => ({
			key: fieldKey,
			path: fieldKey,
			label: formatLabel(fieldKey),
			required: schema.required?.includes(fieldKey) ?? false,
			...(fieldSchema as {
				type: string;
				description?: string;
				properties?: Record<string, unknown>;
				items?: Record<string, unknown>;
				additionalProperties?: Record<string, unknown>;
			})
		}));
	});

	function formatLabel(key: string) {
		return key.replace(/([A-Z])/g, ' $1').replace(/^./, (s) => s.toUpperCase());
	}

	function updateValue(path: string, value: unknown): void {
		const pathArray = path.split('.');
		let current: JsonObject = { ...data };

		// Navigate to the parent object
		let target = current;
		for (let i = 0; i < pathArray?.length - 1; i++) {
			if (!target[pathArray[i]]) {
				target[pathArray[i]] = {};
			}
			target = target[pathArray[i]] as JsonObject;
		}

		// Set the final value
		const finalKey = pathArray[pathArray?.length - 1];
		target[finalKey] = value as JsonValue;

		// Update the data and notify parent
		data = current;
		onUpdate?.(data);
	}

	function getFieldValue(path: string): unknown {
		return path.split('.').reduce((obj, key) => obj?.[key], data);
	}

	function addArrayItem(path: string, itemType: string = 'string'): void {
		const currentArray = (getFieldValue(path) as unknown[]) ?? [];
		const defaultValue =
			itemType === 'string' ? '' : itemType === 'number' ? 0 : itemType === 'object' ? {} : '';
		updateValue(path, [...currentArray, defaultValue]);
	}

	function removeArrayItem(path: string, index: number): void {
		const currentArray = (getFieldValue(path) as unknown[]) ?? [];
		const newArray = currentArray.filter((_, i) => i !== index);
		updateValue(path, newArray);
	}

	function updateArrayItem(path: string, index: number, value: unknown): void {
		const currentArray = (getFieldValue(path) as unknown[]) ?? [];
		const newArray = [...currentArray];
		newArray[index] = value;
		updateValue(path, newArray);
	}

	function addRecordEntry(path: string): void {
		const currentRecord = (getFieldValue(path) as Record<string, unknown>) ?? {};
		updateValue(path, { ...currentRecord, '': '' });
	}

	function removeRecordEntry(path: string, key: string): void {
		const currentRecord = (getFieldValue(path) as Record<string, unknown>) ?? {};
		const { [key]: _, ...rest } = currentRecord;
		updateValue(path, rest);
	}

	function updateRecordKey(path: string, oldKey: string, newKey: string): void {
		const currentRecord = (getFieldValue(path) as Record<string, unknown>) ?? {};
		if (oldKey === newKey || currentRecord[newKey] !== undefined) return;

		const { [oldKey]: value, ...rest } = currentRecord;
		updateValue(path, { ...rest, [newKey]: value });
	}

	function updateRecordValue(path: string, key: string, value: unknown): void {
		const currentRecord = (getFieldValue(path) as Record<string, unknown>) ?? {};
		updateValue(path, { ...currentRecord, [key]: value });
	}

	function getArrayItemType(field: any): string {
		return field.items?.type ?? 'string';
	}

	function isRecordType(field: any): boolean {
		return field.type === 'object' && field.additionalProperties && !field.properties;
	}

	function isNestedObject(field: any): boolean {
		return field.type === 'object' && field.properties && !field.additionalProperties;
	}

	let yamlErrors = $state<string>('');
	function getYamlValue(): string {
		if (!data || typeof data !== 'object' || Object.keys(data).length === 0) {
			return '';
		}
		try {
			return YAML.stringify(data, { indent: 2 });
		} catch (error) {
			console.error('Failed to stringify YAML:', error);
			return '';
		}
	}

	function validateYaml(yamlText: string): void {
		yamlErrors = '';

		if (!yamlText.trim()) {
			data = {};
			onUpdate?.(data);
			return;
		}

		try {
			const parsed = YAML.parse(yamlText);
			if (parsed && typeof parsed === 'object') {
				data = parsed;
				onUpdate?.(data);
			} else {
				yamlErrors = 'YAML must be an object';
			}
		} catch (error) {
			yamlErrors = `Invalid YAML: ${error instanceof Error ? error.message : 'Unknown error'}`;
		}
	}

	function showDescription(fieldType: string): boolean {
		switch (fieldType) {
			case 'string':
			case 'number':
				return false;
			case 'boolean':
			case 'array':
			case 'object':
			case 'record':
				return true;
			default:
				return false;
		}
	}
</script>

{#if protocol && schema}
	<div class="flex flex-col gap-4 border-t pt-4">
		{#if middlewareType === 'plugin'}
			<div class="flex flex-col gap-2">
				<Textarea
					id="plugin"
					value={getYamlValue()}
					placeholder="Edit plugin configuration in YAML format..."
					class="min-h-32 font-mono text-sm"
					onchange={(e) => {
						let textarea = e.target as HTMLTextAreaElement;
						validateYaml(textarea.value);
					}}
				/>
				{#if yamlErrors}
					<p class="text-xs text-red-500">{yamlErrors}</p>
				{/if}
			</div>
		{:else}
			{#each fields as field (field.key)}
				<div
					class={`flex gap-2 ${field.type === 'boolean' ? 'flex-row items-center  justify-between rounded-lg border p-3' : 'flex-col'}`}
				>
					<div class="space-y-1">
						<Label class="text-sm font-medium">
							{field.label}
							{#if field.required}
								<span class="text-red-500">*</span>
							{/if}
						</Label>
						{#if field.description && showDescription(field.type)}
							<p class="text-muted-foreground text-xs">{field.description}</p>
						{/if}
					</div>

					{#if field.type === 'string'}
						<Input
							id={field.path}
							type="text"
							value={(getFieldValue(field.path) as string) ?? ''}
							placeholder={field.description}
							onchange={(e) => {
								let input = e.target as HTMLInputElement;
								return updateValue(field.path, input.value);
							}}
						/>
					{:else if field.type === 'number'}
						<Input
							id={field.path}
							type="number"
							value={(getFieldValue(field.path) as number) ?? ''}
							placeholder={field.description}
							onchange={(e) => {
								let input = e.target as HTMLInputElement;
								return updateValue(field.path, Number(input.value));
							}}
						/>
					{:else if field.type === 'boolean'}
						<CustomSwitch
							checked={(getFieldValue(field.path) as boolean) ?? false}
							onCheckedChange={(checked) => updateValue(field.path, checked)}
						/>
					{:else if field.type === 'array'}
						{@const arrayItems = (getFieldValue(field.path) as unknown[]) ?? []}
						{@const itemType = getArrayItemType(field)}

						<div class="flex flex-col gap-2 rounded-md border p-3">
							{#each arrayItems as item, index (index)}
								<div class="flex items-center gap-2">
									{#if itemType === 'string'}
										<Input
											type="text"
											value={item as string}
											placeholder="Enter value"
											onchange={(e) => {
												let input = e.target as HTMLInputElement;
												updateArrayItem(field.path, index, input.value);
											}}
										/>
									{:else if itemType === 'number'}
										<Input
											type="number"
											value={item as number}
											placeholder="Enter number"
											onchange={(e) => {
												let input = e.target as HTMLInputElement;
												updateArrayItem(field.path, index, Number(input.value));
											}}
										/>
									{:else if itemType === 'object'}
										<div class="flex-1 rounded border p-2">
											<DynamicForm
												middlewareType={`${field.key}.${index}`}
												{protocol}
												data={(item as JsonObject) ?? {}}
												onUpdate={(nestedData) => updateArrayItem(field.path, index, nestedData)}
											/>
										</div>
									{/if}
									<Button
										variant="ghost"
										size="icon"
										class="text-red-500"
										onclick={() => removeArrayItem(field.path, index)}
									>
										<Trash />
									</Button>
								</div>
							{/each}

							<Button
								variant="outline"
								size="sm"
								onclick={() => addArrayItem(field.path, itemType)}
							>
								<Plus />
								Add {field.label}
							</Button>
						</div>
					{:else if field.type === 'object' && isRecordType(field)}
						{@const recordData = (getFieldValue(field.path) as Record<string, unknown>) ?? {}}

						<div class="flex flex-col gap-2 rounded-md border p-3">
							{#each Object.entries(recordData) as [recordKey, recordValue] (recordKey)}
								<div class="flex items-center gap-2">
									<Input
										type="text"
										value={recordKey}
										placeholder="Key"
										class="w-1/3"
										onchange={(e) => {
											let input = e.target as HTMLInputElement;
											updateRecordKey(field.path, recordKey, input.value);
										}}
									/>
									<Input
										type="text"
										value={recordValue as string}
										placeholder="Value"
										class="flex-1"
										onchange={(e) => {
											let input = e.target as HTMLInputElement;
											updateRecordValue(field.path, recordKey, input.value);
										}}
									/>
									<Button
										variant="ghost"
										size="icon"
										class="text-red-500"
										onclick={() => removeRecordEntry(field.path, recordKey)}
									>
										<Trash size={16} />
									</Button>
								</div>
							{/each}

							<Button variant="outline" size="sm" onclick={() => addRecordEntry(field.path)}>
								<Plus size={16} />
								Add Entry
							</Button>
						</div>
					{:else if field.type === 'object' && isNestedObject(field)}
						<div class="rounded-md border p-3">
							<DynamicForm
								middlewareType={field.key}
								{protocol}
								data={(getFieldValue(field.path) as JsonObject) ?? {}}
								onUpdate={(nestedData) => updateValue(field.path, nestedData)}
							/>
						</div>
					{/if}
				</div>
			{/each}
		{/if}
	</div>
{/if}
