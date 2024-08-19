<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	middleware.rateLimit = middleware.rateLimit ?? {
		average: 0,
		period: '1s',
		burst: 1,
		sourceCriterion: {
			ipStrategy: { depth: 0, excludedIPs: [] },
			requestHeaderName: '',
			requestHost: false
		}
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="average" class="text-right">Average</Label>
	<Input
		id="average"
		name="average"
		type="number"
		bind:value={middleware.rateLimit.average}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Average"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="period" class="text-right">Period</Label>
	<Input
		id="period"
		name="period"
		type="text"
		bind:value={middleware.rateLimit.period}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Period"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="burst" class="text-right">Burst</Label>
	<Input
		id="burst"
		name="burst"
		type="number"
		bind:value={middleware.rateLimit.burst}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Burst"
	/>
</div>

<header class="border-b border-gray-200 py-2 font-bold">Source Criterion</header>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="depth" class="text-right">Depth</Label>
	<Input
		id="depth"
		name="depth"
		type="number"
		bind:value={middleware.rateLimit.sourceCriterion.ipStrategy.depth}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Depth"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-2">
	<Label for="request-header-name" class="text-right">Request Header Name</Label>
	<Input
		id="request-header-name"
		name="request-header-name"
		type="text"
		bind:value={middleware.rateLimit.sourceCriterion.requestHeaderName}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Request Header Name"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-2">
	<Label for="request-host" class="text-right">Request Host</Label>
	<Switch
		id="request-host"
		bind:checked={middleware.rateLimit.sourceCriterion.requestHost}
		class="col-span-3"
	/>
</div>
<ArrayInput
	bind:items={middleware.rateLimit.sourceCriterion.ipStrategy.excludedIPs}
	label="Excluded IPs"
/>
