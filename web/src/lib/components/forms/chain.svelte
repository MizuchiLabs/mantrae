<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	export let middleware: Middleware;
	middleware.chain = middleware.chain ?? {
		middlewares: []
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="middlewares" class="text-right">Middlewares</Label>
	<Input
		id="middlewares"
		name="middlewares"
		type="text"
		bind:value={middleware.chain.middlewares}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Middlewares"
	/>
</div>
