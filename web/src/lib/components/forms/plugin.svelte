<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.plugin = { ...middleware.plugin };
</script>

{#if middleware.plugin}
	<div class="grid grid-cols-4 items-center gap-2">
		<Label for="key" class="text-right">Key</Label>
		<Input
			id="request-header-name"
			name="request-header-name"
			type="text"
			bind:value={middleware.plugin.key}
			placeholder="username"
			{disabled}
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-2">
		<Label for="value" class="text-right">Value</Label>
		<!-- <Textarea -->
		<!-- 	id="value" -->
		<!-- 	name="value" -->
		<!-- 	bind:value={middleware.plugin[key]} -->
		<!-- 	{disabled} -->
		<!-- /> -->
	</div>
{/if}
