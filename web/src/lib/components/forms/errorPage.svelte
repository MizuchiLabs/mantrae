<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.errors = {
		status: [],
		service: '',
		query: '',
		...middleware.errors
	};
</script>

{#if middleware.errors}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="service" class="text-right">Service</Label>
		<Input
			id="service"
			name="service"
			type="text"
			bind:value={middleware.errors.service}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="serviceError"
			{disabled}
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="query" class="text-right">Query</Label>
		<Input
			id="query"
			name="query"
			type="text"
			bind:value={middleware.errors.query}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="/status.html"
			{disabled}
		/>
	</div>
	<ArrayInput bind:items={middleware.errors.status} label="Status" placeholder="500" {disabled} />
{/if}
