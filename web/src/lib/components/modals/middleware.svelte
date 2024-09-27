<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { updateMiddleware } from '$lib/api';
	import type { Middleware } from '$lib/types/middlewares';
	import MiddlewareForm from '../forms/middleware.svelte';

	export let middleware: Middleware;
	export let open = false;
	export let disabled = false;
	let originalName = middleware?.name;

	const update = async () => {
		if (middleware.name === '') return;
		await updateMiddleware(middleware);
		originalName = middleware.name;
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<MiddlewareForm bind:middleware {disabled} />
		{#if disabled}
			<Button class="w-full" on:click={() => (open = false)}>Close</Button>
		{:else}
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		{/if}
	</Dialog.Content>
</Dialog.Root>
