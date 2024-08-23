<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { deleteMiddleware, profile, updateMiddleware, middlewares } from '$lib/api';
	import type { Middleware } from '$lib/types/middlewares';

	export let middleware: Middleware;
	let name = middleware.name.split('@')[0];
	let middlewareCompare = $middlewares.filter((m) => m.name !== middleware.name);

	const update = () => {
		if (middleware.name === '' || isNameTaken) return;
		let oldName = middleware.name;
		middleware.name = name + '@' + middleware.provider;
		updateMiddleware($profile, middleware, oldName);
	};

	// Dynamic imports based on middleware type
	// TODO: move to a separate file
	const middlewareForms = {
		addPrefix: () => import('$lib/components/forms/addPrefix.svelte'),
		stripPrefix: () => import('$lib/components/forms/stripPrefix.svelte'),
		stripPrefixRegex: () => import('$lib/components/forms/stripPrefixRegex.svelte'),
		replacePath: () => import('$lib/components/forms/replacePath.svelte'),
		replacePathRegex: () => import('$lib/components/forms/replacePathRegex.svelte'),
		chain: () => import('$lib/components/forms/chain.svelte'),
		ipAllowList: () => import('$lib/components/forms/ipAllowList.svelte'),
		headers: () => import('$lib/components/forms/headers.svelte'),
		errors: () => import('$lib/components/forms/errorPage.svelte'),
		rateLimit: () => import('$lib/components/forms/rateLimit.svelte'),
		redirectRegex: () => import('$lib/components/forms/redirectRegex.svelte'),
		redirectScheme: () => import('$lib/components/forms/redirectScheme.svelte'),
		basicAuth: () => import('$lib/components/forms/basicAuth.svelte'),
		digestAuth: () => import('$lib/components/forms/digestAuth.svelte'),
		forwardAuth: () => import('$lib/components/forms/forwardAuth.svelte'),
		inFlightReq: () => import('$lib/components/forms/inFlightReq.svelte'),
		buffering: () => import('$lib/components/forms/buffering.svelte'),
		circuitBreaker: () => import('$lib/components/forms/circuitBreaker.svelte'),
		compress: () => import('$lib/components/forms/compress.svelte'),
		// passTLSClientCert: () => import('$lib/components/forms/passTLSClientCert.svelte'),
		retry: () => import('$lib/components/forms/retry.svelte'),

		// TCP-specific
		inFlightConn: () => import('$lib/components/forms/inFlightConn.svelte'),
		tcpIpAllowList: () => import('$lib/components/forms/tcpIpAllowList.svelte')
	};

	let MiddlewareFormComponent: any = null;
	const loadMiddlewareFormComponent = async () => {
		const moduleKey = Object.keys(middlewareForms).find(
			(key) => key.toLowerCase() === middleware.type?.toLowerCase()
		) as keyof typeof middlewareForms | undefined;
		if (moduleKey) {
			const module = await middlewareForms[moduleKey]();
			MiddlewareFormComponent = module.default;
		} else {
			MiddlewareFormComponent = null;
		}
	};
	loadMiddlewareFormComponent();

	// Check if middleware name is taken unless self
	let isNameTaken = false;
	$: isNameTaken = middlewareCompare.some((m) => m.name === name + '@' + middleware.provider);

	const onKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			update();
		}
	};
</script>

<Dialog.Root>
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
						bind:value={name}
						on:keydown={onKeydown}
						class={isNameTaken
							? 'col-span-3 border-red-400 focus-visible:ring-0 focus-visible:ring-offset-0'
							: 'col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0'}
						placeholder="Name of the middleware"
						required
					/>
				</div>
				{#if MiddlewareFormComponent}
					<div class="mt-6 space-y-2">
						<svelte:component this={MiddlewareFormComponent} bind:middleware />
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="grid grid-cols-2 items-center justify-between gap-2">
			<Button class="bg-red-400" on:click={() => deleteMiddleware($profile, middleware.name)}
				>Delete</Button
			>
			<Button type="submit" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
