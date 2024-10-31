<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="force-slash" class="text-right">Force Slash</Label>
	<Switch
		id="force-slash"
		bind:checked={middleware.content.forceSlash}
		class="col-span-3"
		{disabled}
	/>
</div>
<ArrayInput
	bind:items={middleware.content.prefixes}
	label="Prefixes"
	placeholder="/foo"
	{disabled}
/>
