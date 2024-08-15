<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: HttpMiddleware;
	middleware.stripPrefix = middleware.stripPrefix ?? {
		prefixes: [''],
		forceSlash: false
	};
</script>

<ArrayInput bind:items={middleware.stripPrefix.prefixes} label="Prefixes" />
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="force-slash" class="text-right">Force Slash</Label>
	<Switch
		id="force-slash"
		bind:checked={middleware.stripPrefix.forceSlash}
		class="col-span-3 justify-self-end"
	/>
</div>
