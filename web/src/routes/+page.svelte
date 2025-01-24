<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import BaseForm from '$lib/components/forms/BaseForm.svelte';
	import { z } from 'zod';
	import {
		type SupportedMiddleware,
		GetSchema,
		MiddlewareTypes
	} from '$lib/components/forms/mw_registry';

	const onSubmit = (data: FormData) => {
		console.log('Form submitted:', data.get('name'));
	};

	let selectedType: SupportedMiddleware | undefined = $state();
	let schema = $state<z.AnyZodObject>();
	let data = $state<Record<string, any>>({} as Record<string, any>);

	const handleSelect = (value: string) => {
		selectedType = value as SupportedMiddleware;
		schema = GetSchema(selectedType);
		data = {};
	};
</script>

<div class="flex flex-col items-center gap-2">
	<Select.Root type="single" bind:value={selectedType} onValueChange={handleSelect}>
		<Select.Trigger class="w-[380px]">
			{selectedType
				? MiddlewareTypes.find((t) => t.value === selectedType)?.label
				: 'Select a middleware type'}
		</Select.Trigger>
		<Select.Content>
			{#each MiddlewareTypes as type}
				<Select.Item value={type.value}>{type.label}</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>

	{#if schema}
		<BaseForm {schema} {data} {onSubmit} />
	{/if}
</div>
