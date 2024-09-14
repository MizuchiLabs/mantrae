<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import type { Selected } from 'bits-ui';
	import { middlewares } from '$lib/api';
	import { newMiddleware, type Middleware } from '$lib/types/middlewares';
	import { LoadMiddlewareForm } from '../utils/middlewareModules';
	import { onMount, SvelteComponent } from 'svelte';

	export let middleware: Middleware;
	export let disabled = false;

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
	let middlewareType: Selected<string> | undefined = HTTPMiddlewareTypes.find(
		(t) => t.value.toLowerCase() === middleware.type
	);

	let form: typeof SvelteComponent | null = null;
	const setMiddlewareType = async (type: Selected<string> | undefined) => {
		if (type === undefined) return;
		middlewareType = type;
		middleware = newMiddleware();
		middleware.type = type.value.toLowerCase();
		form = null;
		form = await LoadMiddlewareForm(middleware);
	};

	let isHTTP = middleware.middlewareType == 'http' ? true : false;
	$: isHTTP, setType();
	const setType = () => {
		middleware.middlewareType = isHTTP ? 'http' : 'tcp';
		if (isHTTP) middlewareType = HTTPMiddlewareTypes[0];
		else middlewareType = TCPMiddlewareTypes[0];
	};

	// Check if middleware name is taken
	let isNameTaken = false;
	$: isNameTaken = $middlewares.some((m) => m.name === middleware.name + '@' + middleware.provider);

	onMount(async () => {
		form = await LoadMiddlewareForm(middleware);
		if (middleware.type === '') {
			setMiddlewareType(HTTPMiddlewareTypes[0]);
		}
	});
</script>

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
	</Card.Header>
	<Card.Content>
		<div class="flex items-center justify-end gap-2">
			<Label for="middleware-type">TCP</Label>
			<Switch id="middleware-type" bind:checked={isHTTP} />
			<Label for="middleware-type">HTTP</Label>
		</div>

		<!-- Type -->
		{#if !disabled}
			<div class="grid grid-cols-4 items-center gap-4 space-y-2">
				<Label for="current" class="text-right">Type</Label>
				<Select.Root onSelectedChange={setMiddlewareType} selected={middlewareType}>
					<Select.Trigger class="col-span-3">
						<Select.Value placeholder="Select a type" />
					</Select.Trigger>
					<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
						{#if middleware.middlewareType === 'http'}
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
		{/if}

		<!-- Name -->
		<div class="grid grid-cols-4 items-center gap-4 space-y-2">
			<Label for="name" class="text-right">Name</Label>
			<Input
				id="name"
				name="name"
				type="text"
				bind:value={middleware.name}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="Name of the middleware"
				required
				{disabled}
			/>
		</div>

		<!-- Dynamic Form -->
		{#if form !== null}
			<div class="mt-6 space-y-2">
				<svelte:component this={form} bind:middleware {disabled} />
			</div>
		{/if}
	</Card.Content>
</Card.Root>
