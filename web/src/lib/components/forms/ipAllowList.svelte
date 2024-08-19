<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	middleware.ipAllowList = middleware.ipAllowList ?? {
		sourceRange: [],
		ipStrategy: { depth: 0, excludedIPs: [] }
	};
</script>

<ArrayInput bind:items={middleware.ipAllowList.sourceRange} label="Source Range" />
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="depth" class="text-right">Depth</Label>
	<Input
		id="depth"
		name="depth"
		type="text"
		bind:value={middleware.ipAllowList.ipStrategy.depth}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="IP Strategy"
	/>
</div>
<ArrayInput bind:items={middleware.ipAllowList.ipStrategy.excludedIPs} label="Excluded IPs" />
