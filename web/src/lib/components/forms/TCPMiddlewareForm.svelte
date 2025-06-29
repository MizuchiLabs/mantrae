<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
	import { unmarshalConfig, marshalConfig } from '$lib/types';
	import { TCPMiddlewareSchema, type TCPMiddleware } from '$lib/gen/zen/traefik-schemas';
	import DynamicForm from './DynamicForm.svelte';

	let { middleware = $bindable() }: { middleware: Middleware } = $props();

	let config = $state(unmarshalConfig(middleware.config) as TCPMiddleware);
	let selectedType = $derived(config ? Object.keys(config)[0] : '');

	$effect(() => {
		if (config) middleware.config = marshalConfig(config);
	});

	const middlewareTypes = Object.keys(TCPMiddlewareSchema.shape).map((key) => ({
		value: key,
		label: key.replace(/([A-Z])/g, ' $1').replace(/^./, (s) => s.toUpperCase()) // dumb prettifier
	}));
</script>

<div class="flex flex-col gap-3">
	<!-- Middleware Type -->
	<div class="flex flex-col gap-2">
		<Label class="mr-2">Type</Label>
		<Select.Root type="single" bind:value={selectedType}>
			<Select.Trigger class="w-full">
				{selectedType
					? middlewareTypes.find((t) => t.value === selectedType)?.label
					: 'Select type'}
			</Select.Trigger>
			<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
				{#each middlewareTypes as t (t.value)}
					<Select.Item value={t.value}>{t.label}</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>

	{#if selectedType}
		<DynamicForm
			schema={TCPMiddlewareSchema.shape[selectedType as keyof typeof TCPMiddlewareSchema.shape]}
			data={(config[selectedType as keyof TCPMiddleware] as Record<string, unknown>) || {}}
			onUpdate={(updatedData) => {
				config = { [selectedType]: updatedData } as TCPMiddleware;
			}}
		/>
	{/if}
</div>
