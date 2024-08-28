<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.retry = { initialInterval: '', ...middleware.retry };
</script>

{#if middleware.retry}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="attempts" class="text-right">Attempts</Label>
		<Input
			id="attempts"
			name="attempts"
			type="number"
			bind:value={middleware.retry.attempts}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="3"
			{disabled}
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="initial-interval" class="text-right">Initial Interval</Label>
		<Input
			id="initial-interval"
			name="initial-interval"
			type="text"
			bind:value={middleware.retry.initialInterval}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="100ms"
			{disabled}
		/>
	</div>
{/if}
