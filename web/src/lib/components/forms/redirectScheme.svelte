<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.redirectScheme = {
		scheme: '',
		port: '',
		permanent: false,
		...middleware.redirectScheme
	};
</script>

{#if middleware.redirectScheme}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="permanent" class="text-right">Permanent</Label>
		<Switch
			id="permanent"
			bind:checked={middleware.redirectScheme.permanent}
			class="col-span-3"
			{disabled}
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="scheme" class="text-right">Scheme</Label>
		<Input
			id="scheme"
			name="scheme"
			type="text"
			bind:value={middleware.redirectScheme.scheme}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="https"
			{disabled}
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="port" class="text-right">Port</Label>
		<Input
			id="port"
			name="port"
			type="text"
			bind:value={middleware.redirectScheme.port}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="443"
			{disabled}
		/>
	</div>
{/if}
