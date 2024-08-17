<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	export let middleware: Middleware;
	middleware.addPrefix = middleware.addPrefix ?? { prefix: '' };
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="prefix" class="text-right">Prefix</Label>
	<Input
		type="text"
		id="prefix"
		placeholder="Enter prefix"
		bind:value={middleware.addPrefix.prefix}
		class="col-span-3"
	/>
</div>
