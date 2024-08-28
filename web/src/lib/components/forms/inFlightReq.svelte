<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.inFlightReq = {
		sourceCriterion: {
			ipStrategy: { excludedIPs: [] },
			requestHeaderName: '',
			requestHost: false
		},
		...middleware.inFlightReq
	};
</script>

{#if middleware.inFlightReq}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="amount" class="text-right">Amount</Label>
		<Input
			id="amount"
			name="amount"
			type="number"
			bind:value={middleware.inFlightReq.amount}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="50"
			{disabled}
		/>
	</div>

	<header class="border-b border-gray-200 py-2 font-bold">Source Criterion</header>
	{#if middleware.inFlightReq.sourceCriterion && middleware.inFlightReq.sourceCriterion.ipStrategy}
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="depth" class="text-right">Depth</Label>
			<Input
				id="depth"
				name="depth"
				type="number"
				bind:value={middleware.inFlightReq.sourceCriterion.ipStrategy.depth}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="0"
				{disabled}
			/>
		</div>
		<div class="grid grid-cols-4 items-center gap-2">
			<Label for="request-header-name" class="text-right">Request Header Name</Label>
			<Input
				id="request-header-name"
				name="request-header-name"
				type="text"
				bind:value={middleware.inFlightReq.sourceCriterion.requestHeaderName}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="X-CustomHeader"
				{disabled}
			/>
		</div>
		<div class="grid grid-cols-4 items-center gap-2">
			<Label for="request-host" class="text-right">Request Host</Label>
			<Switch
				id="request-host"
				bind:checked={middleware.inFlightReq.sourceCriterion.requestHost}
				class="col-span-3"
				{disabled}
			/>
		</div>
		<ArrayInput
			bind:items={middleware.inFlightReq.sourceCriterion.ipStrategy.excludedIPs}
			label="Excluded IPs"
			placeholder="192.168.1.1"
			{disabled}
		/>
	{/if}
{/if}
