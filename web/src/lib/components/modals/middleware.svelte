<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { upsertMiddleware } from '$lib/api';
	import type { Middleware } from '$lib/types/middlewares';
	import MiddlewareForm from '../forms/middleware.svelte';

	export let middleware: Middleware;
	export let open = false;
	export let disabled = false;
	let originalName = middleware?.name;

	const update = async () => {
		if (middleware.name === '') return;
		await upsertMiddleware(originalName, middleware);
		originalName = middleware.name;
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content>
		<MiddlewareForm bind:middleware {disabled} />
		<Dialog.Close>
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
