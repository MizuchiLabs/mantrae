<script lang="ts">
	import { Button } from '$lib/components/ui/button/index';
	import { safeClone } from '$lib/utils';
	import FormField from './FormField.svelte';

	interface Props {
		data: Record<string, unknown>;
		onSubmit: (data: Record<string, unknown>) => void;
		disabled?: boolean;
	}

	let { data = $bindable(), onSubmit, disabled }: Props = $props();

	// Form state
	let formData = $state(safeClone(data));

	// Watch data changes and update formData
	$effect(() => {
		formData = safeClone(data);
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
		type: string;
	};

	// Process object fields
	function processFields(obj: Record<string, unknown>, parentKey = ''): FormFieldType[] {
		return Object.entries(obj).flatMap(([key, value]) => {
			const currentPath = parentKey ? `${parentKey}.${key}` : key;

			if (value && typeof value === 'object' && !Array.isArray(value)) {
				return processFields(value as Record<string, unknown>, currentPath);
			}

			return [
				{
					key,
					path: currentPath,
					value,
					type: Array.isArray(value) ? 'array' : typeof value
				}
			];
		});
	}

	const fields = $derived(processFields(formData));
</script>

<form onsubmit={handleSubmit}>
	<div class="grid gap-4">
		{#each fields as field}
			<FormField {...field} {disabled} bind:data={formData} />
		{/each}
	</div>

	<Button type="submit" class="mt-4 w-full">Save</Button>
</form>
