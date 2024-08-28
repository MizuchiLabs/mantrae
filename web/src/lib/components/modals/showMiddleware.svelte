<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import type { Middleware } from '$lib/types/middlewares';
	import { LoadMiddlewareForm } from '../utils/middlewareModules';
	import { onMount, SvelteComponent } from 'svelte';

	export let middleware: Middleware;
	let name = middleware.name.split('@')[0];

	let form: typeof SvelteComponent | null = null;
	onMount(async () => {
		form = await LoadMiddlewareForm(middleware);
	});
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<Button variant="ghost" class="h-8 w-4 rounded-full bg-green-400">
			<iconify-icon icon="fa6-solid:eye" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[600px]">
		<Card.Root class="mt-4">
			<Card.Header>
				<Card.Title class="flex items-center justify-between gap-2">
					<span>Middleware</span>
					<div>
						<Badge variant="secondary" class="bg-blue-400">
							Type: {middleware.type}
						</Badge>
						<Badge variant="secondary" class="bg-green-400">
							Provider: {middleware.provider}
						</Badge>
					</div>
				</Card.Title>
				<Card.Description>
					Make changes to your Middleware here. Click save when you're done.
				</Card.Description>
			</Card.Header>
			<Card.Content>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="name" class="text-right">Name</Label>
					<Input
						id="name"
						name="name"
						type="text"
						value={name}
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						placeholder="Name of the middleware"
						disabled
					/>
				</div>
				{#if form !== null}
					<div class="mt-6 space-y-2">
						<svelte:component this={form} bind:middleware disabled={true} />
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button class="w-full">Close</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
