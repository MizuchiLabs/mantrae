<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { deleteMiddleware, upsertMiddleware, middlewares } from '$lib/api';
	import type { Middleware } from '$lib/types/middlewares';
	import { LoadMiddlewareForm } from '../utils/middlewareModules';
	import { onMount, type SvelteComponent } from 'svelte';

	export let middleware: Middleware;
	let originalName = middleware.name;
	let middlewareCompare = $middlewares.filter((m) => m.name !== middleware.name);

	let open = false;
	const update = async () => {
		if (middleware.name === '' || isNameTaken) return;
		await upsertMiddleware(originalName, middleware);
		originalName = middleware.name;
		open = false;
	};

	// Check if middleware name is taken unless self
	let isNameTaken = false;
	$: isNameTaken = middlewareCompare.some((m) => m.name === middleware.name);

	const onKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			update();
		}
	};

	let form: typeof SvelteComponent | null = null;
	onMount(async () => {
		form = await LoadMiddlewareForm(middleware);
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		<Button variant="ghost" class="h-8 w-4 rounded-full bg-orange-400">
			<iconify-icon icon="fa6-solid:pencil" />
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
						bind:value={middleware.name}
						on:keydown={onKeydown}
						class={isNameTaken
							? 'col-span-3 border-red-400 focus-visible:ring-0 focus-visible:ring-offset-0'
							: 'col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0'}
						placeholder="Name of the middleware"
						required
					/>
				</div>
				{#if form !== null}
					<div class="mt-6 space-y-2">
						<svelte:component this={form} bind:middleware />
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="grid grid-cols-2 items-center justify-between gap-2">
			<Button class="bg-red-400" on:click={() => deleteMiddleware(middleware.name)}>Delete</Button>
			<Button type="submit" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
