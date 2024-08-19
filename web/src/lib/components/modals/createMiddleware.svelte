<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import type { Selected } from 'bits-ui';
	import { profile, updateMiddleware } from '$lib/api';
	import { newMiddleware } from '$lib/types/middlewares';

	let middleware = newMiddleware();
	let isHTTP = middleware.middlewareType === 'http';

	const create = () => {
		if (middleware.type === '') return;
		middleware.name = middleware.name + '@' + middleware.provider;
		if (isHTTP) middleware.middlewareType = 'http';
		else middleware.middlewareType = 'tcp';

		updateMiddleware($profile, middleware, middleware.name);
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
		{ label: 'Retry', value: 'retry' },
		{ label: 'Content Type', value: 'contentType' }
	];
	const TCPMiddlewareTypes: Selected<string>[] = [
		{ label: 'In Flight Conn', value: 'inFlightConn' },
		{ label: 'IP Allow List', value: 'tcpIpAllowList' }
	];

	// Dynamic imports based on middleware type
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
		contentType: () => import('$lib/components/forms/contentType.svelte'),

		// TCP-specific
		inFlightConn: () => import('$lib/components/forms/inFlightConn.svelte'),
		tcpIpAllowList: () => import('$lib/components/forms/tcpIpAllowList.svelte')
	};

	let MiddlewareFormComponent: any = null;
	let middlewareType: Selected<string> | undefined = HTTPMiddlewareTypes[0];

	const changeMiddlewareType = async (serviceType: Selected<string> | undefined) => {
		if (serviceType === undefined) return;
		middlewareType = { label: serviceType.label || '', value: serviceType.value };
		middleware.type = serviceType.value.toLowerCase();
		await loadMiddlewareFormComponent(serviceType);
	};

	const loadMiddlewareFormComponent = async (serviceType: Selected<string> | undefined) => {
		if (serviceType && middlewareForms[serviceType.value as keyof typeof middlewareForms]) {
			middleware.type = serviceType.value.toLowerCase();
			const module = await middlewareForms[serviceType.value as keyof typeof middlewareForms]();
			MiddlewareFormComponent = module.default;
		} else {
			MiddlewareFormComponent = null;
		}
	};

	// Load the initial form component
	loadMiddlewareFormComponent(middlewareType);
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
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[520px]">
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
					<Select.Root onSelectedChange={changeMiddlewareType} selected={middlewareType}>
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
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
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
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => create()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
