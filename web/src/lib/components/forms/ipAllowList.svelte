<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.ipAllowList = {
		sourceRange: [],
		ipStrategy: { excludedIPs: [] },
		...middleware.ipAllowList
	};
</script>

{#if middleware.ipAllowList && middleware.ipAllowList.ipStrategy}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="depth" class="text-right">Depth</Label>
		<Input
			id="depth"
			name="depth"
			type="text"
			bind:value={middleware.ipAllowList.ipStrategy.depth}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="0"
			{disabled}
		/>
	</div>
	<ArrayInput
		bind:items={middleware.ipAllowList.sourceRange}
		label="Source Range"
		placeholder="192.168.1.1/32"
		{disabled}
	/>
	<ArrayInput
		bind:items={middleware.ipAllowList.ipStrategy.excludedIPs}
		label="Excluded IPs"
		placeholder="192.168.1.1"
		{disabled}
	/>
{/if}
