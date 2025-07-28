<script lang="ts">
	import type { JsonObject } from '@bufbuild/protobuf';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
	import { MiddlewareSchema } from '$lib/gen/zen/traefik-schemas';
	import { unmarshalConfig, marshalConfig } from '$lib/types';
	import DynamicForm from './DynamicForm.svelte';

	let { middleware = $bindable() }: { middleware: Middleware } = $props();

	let middlewareType = $state('');
	$effect(() => {
		if (middleware.config) middlewareType = Object.keys(middleware.config)[0];
	});

	const data = $derived.by(() => {
		if (!middlewareType || !middleware.config?.[middlewareType]) {
			return {};
		}
		return marshalConfig(middleware.config[middlewareType]);
	});

	const middlewareTypes = Object.keys(MiddlewareSchema.shape).map((key) => ({
		value: key,
		label: key.replace(/([A-Z])/g, ' $1').replace(/^./, (s) => s.toUpperCase()) // dumb prettifier
	}));

	function onUpdate(data: JsonObject) {
		if (!middleware.config) middleware.config = {};
		// Clear other middleware types when switching
		middleware.config = {};
		middleware.config[middlewareType] = unmarshalConfig(data);
	}
</script>

<div class="flex flex-col gap-3">
	<!-- Middleware Type -->
	<div class="flex flex-col gap-2">
		<Label class="mr-2">Type</Label>
		<Select.Root type="single" bind:value={middlewareType}>
			<Select.Trigger class="w-full">
				{middlewareType
					? middlewareTypes.find((t) => t.value === middlewareType)?.label
					: 'Select type'}
			</Select.Trigger>
			<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
				{#each middlewareTypes as t (t.value)}
					<Select.Item value={t.value}>{t.label}</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>

	<DynamicForm {middlewareType} protocol="http" {data} {onUpdate} />
</div>
