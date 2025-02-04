<script lang="ts">
	import { Button } from '$lib/components/ui/button/index';
	import { safeClone } from '$lib/utils';
	import FormField from './FormField.svelte';
	import type { FieldMetadata } from '$lib/types/middlewares';
	import Separator from '../ui/separator/separator.svelte';
	import { loading } from '$lib/api';

	interface Props {
		data: Record<string, unknown>;
		metadata?: Record<string, FieldMetadata>;
		onSubmit: (data: Record<string, unknown>) => void;
		disabled?: boolean;
	}

	let { data = $bindable(), metadata = {}, onSubmit, disabled }: Props = $props();

	// Form state
	let formData = $state(safeClone(data));

	// Watch data changes and update formData
	$effect(() => {
		formData = safeClone(data);
		// console.log(formData);
	});

	// Handle form submission
	function handleSubmit(e: Event) {
		e.preventDefault();
		const submissionData = safeClone(formData);
		onSubmit(submissionData);
	}

	type FormFieldType = {
		key: string;
		path: string;
		value: unknown;
		metadata: FieldMetadata;
		type: string;
	};

	// Process object fields
	const fields = $derived(processFields(formData));
	function processFields(obj: Record<string, unknown>, parentKey = ''): FormFieldType[] {
		return Object.entries(obj).flatMap(([key, value]) => {
			const currentPath = parentKey ? `${parentKey}.${key}` : key;
			const fieldMetadata = metadata[currentPath] || {};

			if (value && typeof value === 'object' && !Array.isArray(value)) {
				return processFields(value as Record<string, unknown>, currentPath);
			}

			return [
				{
					key,
					path: currentPath,
					value,
					metadata: fieldMetadata,
					type: Array.isArray(value) ? 'array' : typeof value
				}
			];
		});
	}
</script>

<form onsubmit={handleSubmit}>
	<div class="grid gap-4">
		{#each fields as field}
			<FormField {...field} {disabled} bind:data={formData} />
		{/each}
	</div>

	<Separator class="my-4" />
	<Button type="submit" class="w-full" disabled={$loading}>Save</Button>
</form>
