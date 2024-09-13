<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { deleteMiddleware, upsertMiddleware } from '$lib/api';
	import type { Middleware } from '$lib/types/middlewares';
	import MiddlewareForm from '../forms/middleware.svelte';

	export let middleware: Middleware;
	export let open = false;
	export let disabled = false;
	let originalName = middleware.name;

	const update = async () => {
		if (middleware.name === '') return;
		await upsertMiddleware(originalName, middleware);
		originalName = middleware.name;
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content class="sm:max-w-[600px]">
		<MiddlewareForm bind:middleware {disabled} />
		<Dialog.Close class="grid grid-cols-2 items-center justify-between gap-2">
			<Button class="bg-red-400" on:click={() => deleteMiddleware(middleware.name)}>Delete</Button>
			<Button type="submit" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
