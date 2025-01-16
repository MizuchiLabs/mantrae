<script lang="ts">
	import { run } from 'svelte/legacy';

	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import type { Selected } from 'bits-ui';
	import { middlewares, profile } from '$lib/api';
	import { type Middleware } from '$lib/types/middlewares';
	import { LoadMiddlewareForm } from '../utils/middlewareModules';
	import { onMount, SvelteComponent } from 'svelte';
	import logo from '$lib/images/logo.svg';

	interface Props {
		middleware: Middleware;
		disabled?: boolean;
	}

	let { middleware = $bindable(), disabled = false }: Props = $props();

	const HTTPMiddlewareTypes: Selected<string>[] = [
		{ label: 'Rate Limit', value: 'rateLimit' },
		{ label: 'Headers', value: 'headers' },
		{ label: 'Compress', value: 'compress' },
		{ label: 'Retry', value: 'retry' },
		{ label: 'Whitelist', value: 'ipAllowList' },
		{ label: 'Basic Auth', value: 'basicAuth' },
		{ label: 'Forward Auth', value: 'forwardAuth' },
		{ label: 'Digest Auth', value: 'digestAuth' },
		{ label: 'Chain', value: 'chain' },
		{ label: 'Redirect Scheme', value: 'redirectScheme' },
		{ label: 'Redirect Regex', value: 'redirectRegex' },
		{ label: 'Add Prefix', value: 'addPrefix' },
		{ label: 'Strip Prefix', value: 'stripPrefix' },
		{ label: 'Strip Prefix Regex', value: 'stripPrefixRegex' },
		{ label: 'Replace Path', value: 'replacePath' },
		{ label: 'Replace Path Regex', value: 'replacePathRegex' },
		{ label: 'InFlightReq', value: 'inFlightReq' },
		{ label: 'Circuit Breaker', value: 'circuitBreaker' },
		{ label: 'Buffering', value: 'buffering' },
		{ label: 'Errors', value: 'errors' },
		{ label: 'Pass TLS Client Cert', value: 'passTLSClientCert' },
		{ label: 'Plugin', value: 'plugin' }
	];
	const TCPMiddlewareTypes: Selected<string>[] = [
		{ label: 'InFlightConn', value: 'inFlightConn' },
		{ label: 'Whitelist', value: 'ipAllowList' }
	];

	// Load the initial form component
	let isHTTP = $state(middleware.protocol == 'http' ? true : false);
	let middlewareType: Selected<string> | undefined = $state();

	let form: typeof SvelteComponent | null = $state(null);
	const setMiddlewareType = async (type: Selected<string> | undefined) => {
		if (!type || !$profile.id) return;
		// Delete previous middleware form
		if (middleware.type) {
			delete middleware.content[middleware.type];
		}

		// Set the middleware type
		middlewareType = type;
		middleware.profileId = $profile.id;
		middleware.type = type.value.toLowerCase();
		middleware.protocol = isHTTP ? 'http' : 'tcp';
		form = null;
		form = await LoadMiddlewareForm(middleware);
	};

	const setType = () => {
		if (middleware.type === '') {
			if (isHTTP) setMiddlewareType(HTTPMiddlewareTypes[0]);
			else setMiddlewareType(TCPMiddlewareTypes[0]);
			return;
		} else {
			if (isHTTP) {
				let middlewareType = HTTPMiddlewareTypes.find(
					(t) => t.value.toLowerCase() === middleware.type
				) ?? {
					label: 'Plugin',
					value: 'plugin'
				};
				setMiddlewareType(middlewareType);
			} else {
				setMiddlewareType(
					TCPMiddlewareTypes.find((t) => t.value.toLowerCase() === middleware.type)
				);
			}
		}
	};


	onMount(async () => {
		form = await LoadMiddlewareForm(middleware);
	});
	run(() => {
		isHTTP, setType();
	});
	// Check if middleware name is taken
	let nameTaken = $derived($middlewares.some((m) => m.id !== middleware.id && m.name === middleware.name));
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
		{#if middleware.provider === 'http' && middleware.type !== 'plugin'}
			<div class="flex items-center justify-end gap-2">
				<Label for="middleware-type">TCP</Label>
				<Switch id="middleware-type" bind:checked={isHTTP} />
				<Label for="middleware-type">HTTP</Label>
			</div>
		{/if}

		<!-- Type -->
		{#if !disabled}
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
		{/if}

		<!-- Name -->
		<div class="grid grid-cols-4 items-center gap-4 space-y-2">
			<Label for="name" class="text-right">Name</Label>
			<div class="relative col-span-3">
				<Input
					id="name"
					name="name"
					type="text"
					bind:value={middleware.name}
					placeholder="Name of the middleware"
					required
					{disabled}
				/>
				<!-- Icon based on provider -->
				{#if middleware.provider !== ''}
					<span
						class="pointer-events-none absolute inset-y-0 right-3 flex items-center text-gray-400"
					>
						{#if middleware.provider === 'http'}
							<img src={logo} alt="HTTP" width="20" />
						{/if}
						{#if middleware.provider === 'internal' || middleware.provider === 'file'}
							<iconify-icon icon="devicon:traefikproxy" height="20"></iconify-icon>
						{/if}
						{#if middleware.provider === 'docker' || middleware.provider === 'swarm'}
							<iconify-icon icon="logos:docker-icon" height="20"></iconify-icon>
						{/if}
						{#if middleware.provider === 'kubernetes' || middleware.provider === 'kubernetescrd'}
							<iconify-icon icon="logos:kubernetes" height="20"></iconify-icon>
						{/if}
						{#if middleware.provider === 'consul'}
							<iconify-icon icon="logos:consul" height="20"></iconify-icon>
						{/if}
						{#if middleware.provider === 'nomad'}
							<iconify-icon icon="logos:nomad-icon" height="20"></iconify-icon>
						{/if}
						{#if middleware.provider === 'kv'}
							<iconify-icon icon="logos:redis" height="20"></iconify-icon>
						{/if}
					</span>
				{/if}
			</div>
			{#if nameTaken}
				<div class="col-span-4 text-right text-sm text-red-500">Name already taken</div>
			{/if}
		</div>

		<!-- Dynamic Form -->
		{#if form !== null}
			{@const SvelteComponent_1 = form}
			<div class="mt-6 flex flex-col gap-2">
				<SvelteComponent_1 bind:middleware {disabled} />
			</div>
		{/if}
	</Card.Content>
</Card.Root>
