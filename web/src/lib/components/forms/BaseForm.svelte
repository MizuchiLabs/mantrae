<script lang="ts">
	import * as Form from '$lib/components/ui/form/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { Plus, Trash } from 'lucide-svelte';
	import Textarea from '../ui/textarea/textarea.svelte';
	import type { ZodObjectOrRecord } from './mw_registry';

	type Props = {
		schema: ZodObjectOrRecord;
		// eslint-disable-next-line
		data: Record<string, any>;
		subData?: string[];
		subMultiple?: boolean;
		onSubmit: (data: FormData) => void;
	};

	let { schema, data, subData, subMultiple, onSubmit }: Props = $props();

	const form = superForm(data, {
		id: schema._def.typeName || 'form', // Unique ID per schema
		SPA: true,
		resetForm: true,
		dataType: 'json',
		validators: zodClient(schema),
		validationMethod: 'auto',
		onSubmit: ({ formData }) => {
			onSubmit(formData);
		}
	});

	const { form: formData, enhance, submitting } = form;

	// Helper to get nested field path
	function getFieldPath(parent: string, field: string) {
		return parent ? `${parent}.${field}` : field;
	}

	function removeItem(index: number, fieldPath: string) {
		const [parent, field] = fieldPath.includes('.') ? fieldPath.split('.') : [null, fieldPath];

		if (parent) {
			if (!$formData[parent]) return;
			if (index === 0 || $formData[parent][field].length <= 1) return;
			$formData[parent][field] = [...$formData[parent][field]];
			$formData[parent][field].splice(index, 1);
		} else {
			if (index === 0 || $formData[field].length <= 1) return;
			$formData[field] = [...$formData[field]];
			$formData[field].splice(index, 1);
		}
	}

	function addItem(fieldPath: string) {
		const [parent, field] = fieldPath.includes('.') ? fieldPath.split('.') : [null, fieldPath];

		if (parent) {
			if (!$formData[parent]) $formData[parent] = {};
			if (!$formData[parent][field]) $formData[parent][field] = [];
			$formData[parent][field] = [...$formData[parent][field], ''];
		} else {
			if (!$formData[field]) $formData[field] = [];
			$formData[field] = [...$formData[field], ''];
		}
	}

	function initializeFormData() {
		if (schema instanceof z.ZodObject) {
			Object.entries(schema.shape).forEach(([fieldName, fieldSchema]) => {
				const baseSchema = getBaseType(fieldSchema);

				// Initialize optional objects with empty object
				if (fieldSchema instanceof z.ZodOptional && baseSchema instanceof z.ZodObject) {
					if (!$formData[fieldName]) {
						$formData[fieldName] = {};
						// Initialize all properties of the optional object
						Object.entries(baseSchema.shape).forEach(([subField, subSchema]) => {
							const subBaseType = getBaseType(subSchema);
							if (subBaseType instanceof z.ZodString) {
								$formData[fieldName][subField] = '';
							} else if (subBaseType instanceof z.ZodBoolean) {
								$formData[fieldName][subField] = false;
							} else if (subBaseType instanceof z.ZodArray) {
								$formData[fieldName][subField] = [];
							}
						});
					}
				}
				// Initialize arrays
				if (baseSchema instanceof z.ZodArray) {
					if (!$formData[fieldName] && baseSchema instanceof z.ZodOptional) {
						$formData[fieldName] = [];
					} else if (!Array.isArray($formData[fieldName])) {
						$formData[fieldName] = [$formData[fieldName]];
					}
				}

				// Initialize records
				if (fieldSchema instanceof z.ZodRecord && !$formData[fieldName]) {
					$formData[fieldName] = {};
				}
			});
		}
		if (schema instanceof z.ZodRecord) {
			if (!$formData) $formData = {};
		}
	}

	function getBaseType(fieldSchema: z.ZodTypeAny | unknown) {
		if (fieldSchema instanceof z.ZodOptional || fieldSchema instanceof z.ZodDefault) {
			return fieldSchema._def.innerType;
		}
		return fieldSchema;
	}

	function formatPluginConfig(pluginData: Record<string, unknown>): string {
		if (!pluginData) return '';

		const knownKeys = ['name', 'protocol', 'type'];
		const configKey = Object.keys(pluginData).find((key) => !knownKeys.includes(key));

		return configKey ? JSON.stringify(pluginData[configKey], null, 2) : '';
	}

	$effect(() => {
		initializeFormData();
	});
</script>

<form method="POST" use:enhance class="flex flex-col gap-2">
	{#if schema instanceof z.ZodObject}
		{#each Object.entries(schema.shape) as [fieldName, fieldSchema]}
			{@const baseType = getBaseType(fieldSchema)}
			{#if baseType instanceof z.ZodString}
				<Form.Field {form} name={fieldName}>
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>{fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}</Form.Label>
							<Input {...props} bind:value={$formData[fieldName]} disabled={$submitting} />
						{/snippet}
					</Form.Control>
					<Form.FieldErrors />
				</Form.Field>
			{:else if baseType instanceof z.ZodNumber || baseType instanceof z.ZodBigInt}
				<Form.Field {form} name={fieldName}>
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>{fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}</Form.Label>
							<Input
								{...props}
								type="number"
								bind:value={$formData[fieldName]}
								disabled={$submitting}
							/>
						{/snippet}
					</Form.Control>
					<Form.FieldErrors />
				</Form.Field>
			{:else if baseType instanceof z.ZodBoolean}
				<Form.Field {form} name={fieldName}>
					<Form.Control>
						{#snippet children({ props })}
							<div class="flex items-center gap-2">
								<Form.Label>{fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}</Form.Label>
								<Checkbox {...props} checked={$formData[fieldName]} disabled={$submitting} />
							</div>
						{/snippet}
					</Form.Control>
					<Form.FieldErrors />
				</Form.Field>
			{:else if baseType instanceof z.ZodArray && !subData}
				<Form.Fieldset {form} name={fieldName}>
					<Form.Legend>{fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}</Form.Legend>
					<!-- eslint-disable-next-line -->
					{#each $formData[fieldName] || [] as _i, i}
						<Form.ElementField {form} name="{fieldName}[{i}]">
							<Form.Control>
								{#snippet children({ props })}
									<div class="flex gap-2">
										<Form.Label class="sr-only">URL {i + 1}</Form.Label>
										<Input type="text" {...props} bind:value={$formData[fieldName][i]} />
										<Form.Button
											variant="ghost"
											type="button"
											size="icon"
											onclick={() => removeItem(i, fieldName)}
										>
											<Trash />
										</Form.Button>
									</div>
								{/snippet}
							</Form.Control>
							<Form.FieldErrors />
						</Form.ElementField>
					{/each}
					<Form.Button type="button" variant="outline" onclick={() => addItem(fieldName)}>
						<Plus />
						Add {fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}
					</Form.Button>
					<Form.FieldErrors />
				</Form.Fieldset>
			{:else if baseType instanceof z.ZodArray && subData}
				<Form.Field {form} name={fieldName}>
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>{fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}</Form.Label>
							<Select.Root
								type={subMultiple ? 'multiple' : 'single'}
								bind:value={$formData[fieldName]}
								name={`${fieldName}[]`}
								on:change={(event) => {
									if (subMultiple) {
										$formData[fieldName] = Array.isArray(event.detail)
											? event.detail
											: event.detail
												? [event.detail]
												: [];
									} else {
										$formData[fieldName] = event.detail;
									}
								}}
								disabled={$submitting}
							>
								<Select.Trigger {...props}>
									{$formData[fieldName]?.length > 0 ? $formData[fieldName].join(', ') : 'Select...'}
								</Select.Trigger>
								<Select.Content>
									{#each subData ?? [] as item}
										<Select.Item value={item}>{item}</Select.Item>
									{/each}
								</Select.Content>
							</Select.Root>
						{/snippet}
					</Form.Control>
					<Form.FieldErrors />
				</Form.Field>
			{:else if baseType instanceof z.ZodEnum}
				<Form.Field {form} name={fieldName}>
					<Form.Control>
						<Form.Label>{fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}</Form.Label>
						<Select.Root type="single" bind:value={$formData[fieldName]} disabled={$submitting}>
							<Select.Trigger>
								{$formData[fieldName] || 'Select...'}
							</Select.Trigger>
							<Select.Content>
								{#each Object.entries(baseType._def.values) as [value, label]}
									<Select.Item {value}>{label}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</Form.Control>
					<Form.FieldErrors />
				</Form.Field>
			{/if}

			<!-- Objects/Records (recursive) -->
			{#if baseType instanceof z.ZodObject}
				<Form.Fieldset {form} name={fieldName} class="py-4">
					<Form.Legend>{fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}</Form.Legend>
					<div class="border-l-2 border-gray-200 pl-4">
						{#if $formData[fieldName] !== undefined}
							{#each Object.entries(baseType.shape || {}) as [subFieldName, subFieldSchema]}
								{@const subFieldPath = getFieldPath(fieldName, subFieldName)}
								{@const subBaseType = getBaseType(subFieldSchema)}

								{#if subBaseType instanceof z.ZodString}
									<Form.Field {form} name={subFieldPath}>
										<Form.Control>
											{#snippet children({ props })}
												<Form.Label
													>{subFieldName.charAt(0).toUpperCase() +
														subFieldName.slice(1)}</Form.Label
												>
												<Input
													{...props}
													bind:value={$formData[fieldName][subFieldName]}
													disabled={$submitting}
												/>
											{/snippet}
										</Form.Control>
										<Form.FieldErrors />
									</Form.Field>
								{:else if subBaseType instanceof z.ZodNumber}
									<Form.Field {form} name={subFieldPath}>
										<Form.Control>
											{#snippet children({ props })}
												<Form.Label
													>{subFieldName.charAt(0).toUpperCase() +
														subFieldName.slice(1)}</Form.Label
												>
												<Input
													{...props}
													type="number"
													bind:value={$formData[fieldName][subFieldName]}
													disabled={$submitting}
												/>
											{/snippet}
										</Form.Control>
										<Form.FieldErrors />
									</Form.Field>
								{:else if subBaseType instanceof z.ZodBoolean}
									<Form.Field {form} name={subFieldPath}>
										<Form.Control>
											{#snippet children({ props })}
												<div class="flex items-center gap-2">
													<Form.Label
														>{subFieldName.charAt(0).toUpperCase() +
															subFieldName.slice(1)}</Form.Label
													>
													<Checkbox
														{...props}
														checked={$formData[fieldName][subFieldName]}
														disabled={$submitting}
													/>
												</div>
											{/snippet}
										</Form.Control>
										<Form.FieldErrors />
									</Form.Field>
								{:else if subBaseType instanceof z.ZodRecord}
									<Form.Field {form} name={subFieldPath}>
										<Form.Control>
											{#snippet children({ props })}
												<div class="flex items-center gap-2">
													<Form.Label
														>{subFieldName.charAt(0).toUpperCase() +
															subFieldName.slice(1)}</Form.Label
													>
													<Textarea
														{...props}
														bind:value={$formData[fieldName][subFieldName]}
														disabled={$submitting}
													/>
												</div>
											{/snippet}
										</Form.Control>
										<Form.FieldErrors />
									</Form.Field>
								{/if}
							{/each}
						{/if}
					</div>
				</Form.Fieldset>
			{/if}
		{/each}
	{/if}

	{#if schema instanceof z.ZodRecord}
		<Form.Field {form} name="record">
			<Form.Control>
				{#snippet children({ props })}
					<div class="flex items-center gap-2">
						<Textarea
							{...props}
							value={formatPluginConfig($formData)}
							rows={10}
							disabled={$submitting}
						/>
					</div>
				{/snippet}
			</Form.Control>
			<Form.FieldErrors />
		</Form.Field>
	{/if}

	<Form.Button type="submit" disabled={$submitting}>Save</Form.Button>
</form>
