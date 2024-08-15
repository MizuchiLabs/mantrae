<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	export let middleware: HttpMiddleware;
	middleware.retry = middleware.retry ?? {
		attempts: 0,
		initialInterval: ''
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="attempts" class="text-right">Attempts</Label>
	<Input
		id="attempts"
		name="attempts"
		type="number"
		bind:value={middleware.retry.attempts}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Attempts"
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
		placeholder="Initial Interval"
	/>
</div>
