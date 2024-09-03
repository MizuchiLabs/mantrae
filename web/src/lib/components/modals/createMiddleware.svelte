<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import type { Selected } from 'bits-ui';
	import { upsertMiddleware, middlewares } from '$lib/api';
	import { newMiddleware } from '$lib/types/middlewares';
	import { LoadMiddlewareForm } from '../utils/middlewareModules';
	import { onMount, SvelteComponent } from 'svelte';

	let middleware = newMiddleware();
	let isHTTP = middleware.middlewareType === 'http';

	const create = async () => {
		if (middleware.type === '' || middleware.name === '' || isNameTaken) return;
		if (isHTTP) middleware.middlewareType = 'http';
		else middleware.middlewareType = 'tcp';

		await upsertMiddleware(middleware.name, middleware);
		middleware = newMiddleware();
		middlewareType = HTTPMiddlewareTypes[0];
	};

	const HTTPMiddlewareTypes: Selected<string>[] = [
		{ label: 'Add Prefix', value: 'addPrefix' },
		{ label: 'Strip Prefix', value: 'stripPrefix' },
		{ label: 'Strip Prefix Regex', value: 'stripPrefixRegex' },
		{ label: 'Replace Path', value: 'replacePath' },
		{ label: 'Replace Path Regex', value: 'replacePathRegex' },
		{ label: 'Chain', value: 'chain' },
		{ label: 'IP Allow List', value: 'ipAllowList' },
		{ label: 'Headers', value: 'headers' },
		{ label: 'Errors', value: 'errors' },
		{ label: 'Rate Limit', value: 'rateLimit' },
		{ label: 'Redirect Regex', value: 'redirectRegex' },
		{ label: 'Redirect Scheme', value: 'redirectScheme' },
		{ label: 'Basic Auth', value: 'basicAuth' },
		{ label: 'Digest Auth', value: 'digestAuth' },
		{ label: 'Forward Auth', value: 'forwardAuth' },
		{ label: 'InFlightReq', value: 'inFlightReq' },
		{ label: 'Buffering', value: 'buffering' },
		{ label: 'Circuit Breaker', value: 'circuitBreaker' },
		{ label: 'Compress', value: 'compress' },
		//{ label: 'Pass TLS Client Cert', value: 'passTLSClientCert' },
		{ label: 'Retry', value: 'retry' }
	];
	const TCPMiddlewareTypes: Selected<string>[] = [
		{ label: 'In Flight Conn', value: 'inFlightConn' },
		{ label: 'IP Allow List', value: 'tcpIpAllowList' }
	];

	// Load the initial form component
	let form: typeof SvelteComponent | null = null;
	let middlewareType: Selected<string> | undefined = HTTPMiddlewareTypes[0];
	const setMiddlewareType = async (type: Selected<string> | undefined) => {
		if (type === undefined) return;
		middlewareType = type;
		middleware = newMiddleware();
		middleware.type = type.value.toLowerCase();
		form = await LoadMiddlewareForm(middleware);
	};

	// Check if middleware name is taken
	let isNameTaken = false;
	$: isNameTaken = $middlewares.some((m) => m.name === middleware.name + '@' + middleware.provider);

	onMount(async () => {
		form = await LoadMiddlewareForm(middleware);
		await setMiddlewareType(middlewareType);
	});
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<div class="flex w-full flex-row items-center justify-between">
			<Button class="flex items-center gap-2 bg-red-400 text-black">
				<span>Add Middleware</span>
				<iconify-icon icon="fa6-solid:plus" />
			</Button>
		</div>
	</Dialog.Trigger>
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[600px]">
		<Card.Root class="mt-4">
			<Card.Header>
				<Card.Title>Middleware</Card.Title>
				<Card.Description>
					Make changes to your Middleware here. Click save when you're done.
				</Card.Description>
			</Card.Header>
			<Card.Content>
				<div class="flex justify-end space-y-2">
					<div class="flex items-center space-x-2">
						<Label for="middleware-type">TCP</Label>
						<Switch id="middleware-type" bind:checked={isHTTP} />
						<Label for="middleware-type">HTTP</Label>
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-4 space-y-2">
					<Label for="current" class="text-right">Type</Label>
					<Select.Root onSelectedChange={setMiddlewareType} selected={middlewareType}>
						<Select.Trigger class="col-span-3">
							<Select.Value placeholder="Select a type" />
						</Select.Trigger>
						<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
							{#if isHTTP}
								{#each HTTPMiddlewareTypes as type}
									<Select.Item value={type.value} label={type.label}>
										{type.label}
									</Select.Item>
								{/each}
							{:else}
								{#each TCPMiddlewareTypes as type}
									<Select.Item value={type.value} label={type.label}>
										{type.label}
									</Select.Item>
								{/each}
							{/if}
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid grid-cols-4 items-center gap-4 space-y-2">
					<Label for="name" class="text-right">Name</Label>
					<Input
						id="name"
						name="name"
						type="text"
						bind:value={middleware.name}
						class={isNameTaken
							? 'col-span-3 border-red-400 focus-visible:ring-0 focus-visible:ring-offset-0'
							: 'col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0'}
						placeholder="Name of the middleware"
						required
					/>
				</div>
				{#if form !== null}
					{#if middleware.type === 'basicauth' || middleware.type === 'digestauth'}
						<header class="mt-4 text-right text-sm font-semibold">
							Password will be hashed automatically.<br /> You will not be able to see the password again!
						</header>
					{/if}
					<div class="mt-6 space-y-2">
						<svelte:component this={form} bind:middleware />
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => create()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
