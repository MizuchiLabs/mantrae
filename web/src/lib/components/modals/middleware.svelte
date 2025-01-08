<script lang="ts">
	import { upsertMiddleware } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import type { Middleware } from '$lib/types/middlewares';
	import MiddlewareForm from '../forms/middleware.svelte';
	import { cleanEmptyObjects } from '../utils/validation';

	export let middleware: Middleware;
	export let open = false;
	export let disabled = false;

	const update = async () => {
		if (middleware.name === '') return;
		cleanEmptyObjects(middleware.content);
		await upsertMiddleware(middleware);
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<MiddlewareForm bind:middleware {disabled} />
		{#if disabled}
			<Button class="w-full" on:click={() => (open = false)}>Close</Button>
		{:else}
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		{/if}
	</Dialog.Content>
</Dialog.Root>
