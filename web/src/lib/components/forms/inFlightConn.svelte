<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	export let middleware: Middleware;
	middleware.inFlightConn = { ...middleware.inFlightConn };
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="amount" class="text-right">Amount</Label>
	{#if middleware.inFlightConn}
		<Input
			id="amount"
			name="amount"
			type="number"
			bind:value={middleware.inFlightConn.amount}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="50"
		/>
	{/if}
</div>
