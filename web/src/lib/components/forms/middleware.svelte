<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { profile } from '$lib/api';
	import type { Middleware } from '$lib/types/middlewares';
	import logo from '$lib/images/logo.svg';
	import Testform from './testform.svelte';
	import type { Component, SvelteComponent } from 'svelte';

	interface Props {
		middleware: Middleware;
		mode: 'create' | 'edit';
		disabled?: boolean;
	}

	let { middleware = $bindable({} as Middleware), mode, disabled = false }: Props = $props();

	// Computed properties
	let middlewareProvider = $derived(middleware.name ? middleware.name.split('@')[1] : 'http');
	let isHttpProvider = $derived(middlewareProvider === 'http' || !middlewareProvider);

	const middlewareTypes: Record<string, { value: string; label: string; form: any | undefined }[]> =
		{
			http: [
				{ value: 'rateLimit', label: 'Rate Limit', form: Testform },
				{ value: 'headers', label: 'Headers' }
				// { value: 'compress', label: 'Compress' },
				// { value: 'retry', label: 'Retry' },
				// { value: 'ipAllowList', label: 'Whitelist' },
				// { value: 'basicAuth', label: 'Basic Auth' },
				// { value: 'forwardAuth', label: 'Forward Auth' },
				// { value: 'digestAuth', label: 'Digest Auth' },
				// { value: 'chain', label: 'Chain' },
				// { value: 'redirectScheme', label: 'Redirect Scheme' },
				// { value: 'redirectRegex', label: 'Redirect Regex' },
				// { value: 'addPrefix', label: 'Add Prefix' },
				// { value: 'stripPrefix', label: 'Strip Prefix' },
				// { value: 'stripPrefixRegex', label: 'Strip Prefix Regex' },
				// { value: 'replacePath', label: 'Replace Path' },
				// { value: 'replacePathRegex', label: 'Replace Path Regex' },
				// { value: 'inFlightReq', label: 'InFlightReq' },
				// { value: 'circuitBreaker', label: 'Circuit Breaker' },
				// { value: 'buffering', label: 'Buffering' },
				// { value: 'errors', label: 'Errors' },
				// { value: 'passTLSClientCert', label: 'Pass TLS Client Cert' },
				// { value: 'plugin', label: 'Plugin' }
			],
			tcp: [
				{ value: 'inFlightConn', label: 'InFlightConn', form: Testform },
				{ value: 'ipAllowList', label: 'Whitelist', form: null }
			]
		};

	let isHTTP = $state(middleware.protocol === 'http');
	let MiddlewareForm = $state(null);

	async function handleMiddlewareTypeChange(value: string) {
		if (!value || !$profile.id) return;

		middleware.type = isHTTP ? 'http' : 'tcp';
		middleware.protocol = isHTTP ? 'http' : 'tcp';
		MiddlewareForm = middlewareTypes[isHTTP ? 'http' : 'tcp'].find((t) => t.value === value)?.form;
	}
</script>

<Card.Root class="mt-4">
	<Card.Header>
		<Card.Title>{mode === 'create' ? 'Add' : 'Update'} Middleware</Card.Title>
		<Card.Description>
			{mode === 'create' ? 'Create a new middleware' : 'Edit existing middleware'}
		</Card.Description>
	</Card.Header>
	<Card.Content class="flex flex-col gap-4">
		<!-- Protocol Switch -->
		{#if isHttpProvider}
			<div class="flex items-center justify-end gap-2 pb-2">
				<Switch id="protocol-type" bind:checked={isHTTP} />
				<Label for="protocol-type">{isHTTP ? 'HTTP' : 'TCP'}</Label>
			</div>
		{/if}

		<!-- Type Selector -->
		<div class="grid grid-cols-4 items-center gap-2">
			<Label for="type" class="text-right">Type</Label>
			<Select.Root type="single" value={middleware.type} onValueChange={handleMiddlewareTypeChange}>
				<Select.Trigger class="col-span-3">
					{middleware.type
						? middlewareTypes[isHTTP ? 'http' : 'tcp'].find((t) => t.value === middleware.type)
								?.label
						: 'Select type'}
				</Select.Trigger>
				<Select.Content>
					{#each middlewareTypes[isHTTP ? 'http' : 'tcp'] as type}
						<Select.Item value={type.value}>
							{type.label}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>

		<!-- Name Input -->
		<div class="grid grid-cols-4 items-center gap-2">
			<Label for="name" class="text-right">Name</Label>
			<div class="relative col-span-3">
				<Input
					id="name"
					bind:value={middleware.name}
					placeholder="Middleware name"
					{disabled}
					required
				/>
				{#if middlewareProvider}
					<span
						class="pointer-events-none absolute inset-y-0 right-3 flex items-center text-gray-400"
					>
						{#if isHttpProvider}
							<img src={logo} alt="HTTP" width="20" />
						{:else if middlewareProvider.includes('docker')}
							<iconify-icon icon="logos:docker-icon" height="20"></iconify-icon>
						{:else if middlewareProvider.includes('kubernetes')}
							<iconify-icon icon="logos:kubernetes" height="20"></iconify-icon>
						{:else if middlewareProvider === 'consul'}
							<iconify-icon icon="logos:consul" height="20"></iconify-icon>
						{:else if middlewareProvider === 'nomad'}
							<iconify-icon icon="logos:nomad-icon" height="20"></iconify-icon>
						{:else if middlewareProvider === 'kv'}
							<iconify-icon icon="logos:redis" height="20"></iconify-icon>
						{/if}
					</span>
				{/if}
			</div>
		</div>

		<!-- Dynamic Form -->
		{#if MiddlewareForm}
			<div class="mt-6 flex flex-col gap-2">
				<MiddlewareForm bind:middleware {disabled} />
			</div>
		{/if}
	</Card.Content>
</Card.Root>
