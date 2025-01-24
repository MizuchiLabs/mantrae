<script lang="ts">
	import * as Form from '$lib/components/ui/form/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { Plus, Trash } from 'lucide-svelte';

	type Props = {
		schema: z.AnyZodObject;
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
		onSubmit: ({ formData }) => {
			// Convert FormData to a regular object
			const processedData = new FormData();

			// Handle array fields specially
			Object.entries(formData).forEach(([key, value]) => {
				if (Array.isArray(value)) {
					// For arrays, append each value with the same key
					value.forEach((item) => {
						processedData.append(key, item);
					});
				} else {
					processedData.append(key, value);
				}
			});

			onSubmit(processedData);
		}
	});

	const { form: formData, enhance, submitting } = form;

	function removeItem(index: number, fieldName: string) {
		// if optional able to set to undefined
		const isOptional = schema.shape[fieldName] instanceof z.ZodOptional;
		if ((index === 0 || $formData[fieldName].length <= 1) && !isOptional) return;
		$formData[fieldName] = [...($formData[fieldName] || [])];
		$formData[fieldName].splice(index, 1);
	}
	function addItem(fieldName: string) {
		$formData[fieldName] = [...($formData[fieldName] || []), ''];
	}
	function initializeArrays() {
		if (!subData) return;
		Object.entries(schema.shape).forEach(([fieldName, fieldSchema]) => {
			if (fieldSchema instanceof z.ZodArray) {
				if (!$formData[fieldName]) {
					$formData[fieldName] = [];
				} else if (!Array.isArray($formData[fieldName])) {
					$formData[fieldName] = [$formData[fieldName]];
				}
			}
		});
	}
	function getBaseType(fieldSchema: z.ZodTypeAny | unknown) {
		if (fieldSchema instanceof z.ZodOptional) {
			return fieldSchema._def.innerType;
		}
		return fieldSchema;
	}

	$effect(() => {
		initializeArrays();
	});
</script>

<form method="POST" use:enhance class="flex flex-col gap-2">
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
							name={fieldName}
							on:change={() => ($formData[fieldName] = $formData[fieldName])}
							disabled={$submitting}
						>
							<Select.Trigger {...props}>
								{Array.isArray($formData[fieldName])
									? $formData[fieldName].join(', ')
									: 'Select...'}
							</Select.Trigger>
							<Select.Content>
								{#each subData as item}
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
	{/each}

	<Form.Button type="submit" disabled={$submitting}>Save</Form.Button>
</form>
