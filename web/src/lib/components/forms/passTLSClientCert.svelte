<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.content = {
		pem: true
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="pem" class="text-right">Pem</Label>
	<Switch id="remove-header" bind:checked={middleware.content.pem} class="col-span-3" {disabled} />
</div>
