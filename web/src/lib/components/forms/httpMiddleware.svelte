<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Middleware as HttpMiddleware } from '$lib/gen/tygo/dynamic';
	import { type Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
	import { unmarshalConfig, marshalConfig, HTTPMiddlewareKeys } from '$lib/types';

	let { middleware = $bindable() }: { middleware: Middleware } = $props();

	let config = $state(unmarshalConfig(middleware.config) as HttpMiddleware);
	let selectedType = $state('');

	$effect(() => {
		if (config) middleware.config = marshalConfig(config);
	});
</script>

<div class="flex flex-col gap-3">
	<!-- Middleware Type -->
	<div class="flex flex-col gap-2">
		<Label class="mr-2">Type</Label>
		<Select.Root type="single" bind:value={selectedType}>
			<Select.Trigger class="w-full">
				{selectedType
					? HTTPMiddlewareKeys.find((t) => t.value === selectedType)?.label
					: 'Select type'}
			</Select.Trigger>
			<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
				{#each HTTPMiddlewareKeys as t (t.value)}
					<Select.Item value={t.value}>{t.label}</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>

	{#if selectedType}
		<!-- dynamic form -->
	{/if}
</div>
