<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Toggle } from '$lib/components/ui/toggle';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import {
		routers,
		entrypoints,
		middlewares,
		provider,
		toggleEntrypoint,
		toggleMiddleware,
		toggleDNSProvider,
		profile
	} from '$lib/api';
	import { newRouter, type Router } from '$lib/types/config';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import logo from '$lib/images/logo.svg';
	import { z } from 'zod';
	import type { Selected } from 'bits-ui';
	import autoAnimate from '@formkit/auto-animate';

	export let router: Router;
	export let disabled = false;

	const schema = z.object({
		name: z.string().trim().min(1, 'Name is required').max(255),
		provider: z.string().trim().nullish(),
		status: z.string().trim().nullish(),
		protocol: z
			.string()
			.trim()
			.toLowerCase()
			.regex(/^(http|tcp|udp)$/),
		dnsProvider: z.coerce.number().int().nullish(),
		entrypoints: z.array(z.string()).nullish(),
		middlewares: z.array(z.string()).nullish(),
		rule: z.string().trim().optional(),
		priority: z.coerce.number().int().nonnegative().optional(),
		tls: z
			.object({
				certResolver: z.string().trim().nullish()
			})
			.optional()
	});
	// .refine(
	// 	(data: Router) => {
	// 		// Conditionally check if `rule` is required for `http` or `tcp`
	// 		if (['http', 'tcp'].includes(data.protocol) && !data.rule) {
	// 			return false;
	// 		}
	// 		return true;
	// 	},
	// 	{
	// 		message: 'Rule is required for HTTP and TCP routers',
	// 		path: ['rule'] // This points to the 'rule' field
	// 	}
	//);

	let errors: Record<any, string[] | undefined> = {};
	export const validate = () => {
		try {
			schema.parse({ ...router });
			errors = {};
			return true;
		} catch (err: any) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
			return false;
		}
	};

	const changeRouterType = (value: string, oldRouter: Router) => {
		router = newRouter();
		router.profileId = $profile?.id ?? 0;
		router.name = oldRouter.name;
		router.entryPoints = oldRouter.entryPoints;
		router.middlewares = oldRouter.middlewares;
		router.dnsProvider = oldRouter.dnsProvider;
		router.priority = oldRouter.priority;
		router.tls = oldRouter.tls;
		router.protocol = value;
		if (value === 'udp') {
			router.rule = '';
		}
	};
	const getSelectedEntrypoints = (router: Router): Selected<unknown>[] => {
		let list = router?.entryPoints?.map((ep) => {
			return { value: ep, label: ep };
		});
		return list ?? [];
	};
	const getSelectedMiddlewares = (router: Router): Selected<unknown>[] => {
		let list = router?.middlewares?.map((middleware) => {
			return { value: middleware, label: middleware };
		});
		return list ?? [];
	};
	const getSelectedDNSProvider = (router: Router): Selected<number> | undefined => {
		let name = $provider.find((p) => p.id === router.dnsProvider)?.name;
		return router?.dnsProvider ? { value: router.dnsProvider, label: name } : undefined;
	};
	const getCertResolver = () => {
		const certResolvers = $routers
			.filter((item) => item.tls && item.tls.certResolver)
			.map((item) => item.tls.certResolver);

		// Use a Set to remove duplicates and return as an array
		return [...new Set(certResolvers)];
	};
	const setCertResolver = (resolver: string | undefined) => {
		router.tls = router.tls || {};
		router.tls.certResolver = resolver;
	};

	// Check if router name is taken unless self
	$: nameTaken = $routers.some((r) => r.id !== router.id && r.name === router.name);

	// Set default TLS settings
	$: router.tls = router.tls || {};
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Router</Card.Title>
		<Card.Description>Configure your router settings</Card.Description>
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		<!-- Protocol -->
		{#if router.provider === 'http'}
			<div class="flex items-center justify-end gap-1 font-mono text-sm" use:autoAnimate>
				<Toggle
					size="sm"
					pressed={router.protocol === 'http'}
					onPressedChange={() => changeRouterType('http', router)}
					class="font-bold data-[state=on]:bg-green-300  dark:data-[state=on]:text-black"
				>
					HTTP
				</Toggle>
				<Toggle
					size="sm"
					pressed={router.protocol === 'tcp'}
					onPressedChange={() => changeRouterType('tcp', router)}
					class="font-bold data-[state=on]:bg-blue-300 dark:data-[state=on]:text-black"
				>
					TCP
				</Toggle>
				<Toggle
					size="sm"
					pressed={router.protocol === 'udp'}
					onPressedChange={() => changeRouterType('udp', router)}
					class="font-bold data-[state=on]:bg-red-300 dark:data-[state=on]:text-black"
				>
					UDP
				</Toggle>
			</div>
		{/if}

		<!-- Name -->
		<div class="grid grid-cols-4 items-center gap-1">
			<Label for="name" class="mr-2 text-right">Name</Label>
			<div class="relative col-span-3">
				<Input
					id="name"
					name="name"
					type="text"
					bind:value={router.name}
					on:input={validate}
					placeholder="Name of the router"
					required
					{disabled}
				/>
				<!-- Icon based on provider -->
				{#if router.provider !== ''}
					<span
						class="pointer-events-none absolute inset-y-0 right-3 flex items-center text-gray-400"
					>
						{#if router.provider === 'http'}
							<img src={logo} alt="HTTP" width="20" />
						{/if}
						{#if router.provider === 'internal' || router.provider === 'file'}
							<iconify-icon icon="devicon:traefikproxy" height="20" />
						{/if}
						{#if router.provider?.includes('docker')}
							<iconify-icon icon="logos:docker-icon" height="20" />
						{/if}
						{#if router.provider?.includes('kubernetes')}
							<iconify-icon icon="logos:kubernetes" height="20" />
						{/if}
						{#if router.provider === 'consul'}
							<iconify-icon icon="logos:consul" height="20" />
						{/if}
						{#if router.provider === 'nomad'}
							<iconify-icon icon="logos:nomad-icon" height="20" />
						{/if}
						{#if router.provider === 'kv'}
							<iconify-icon icon="logos:redis" height="20" />
						{/if}
					</span>
				{/if}
			</div>
			{#if nameTaken}
				<div class="col-span-4 text-right text-sm text-red-500">Name already taken</div>
			{/if}
			{#if errors.name}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.name}</div>
			{/if}
		</div>

		<!-- Entrypoints -->
		{#if router.provider === 'http'}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label for="entrypoints" class="mr-2 text-right">Entrypoints</Label>
				<Select.Root
					multiple={true}
					selected={getSelectedEntrypoints(router)}
					onSelectedChange={(value) => toggleEntrypoint(router, value, false)}
				>
					<Select.Trigger class="col-span-3">
						<Select.Value placeholder="Select an entrypoint" />
					</Select.Trigger>
					<Select.Content>
						{#each $entrypoints || [] as entrypoint}
							<Select.Item value={entrypoint.name}>
								<div class="flex flex-row items-center gap-2">
									{entrypoint.name}
									{#if entrypoint.http}
										{#if 'tls' in entrypoint.http}
											<iconify-icon icon="fa6-solid:lock" class=" text-green-400" />
										{/if}
									{/if}
								</div>
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{:else}
			<ArrayInput bind:items={router.entryPoints} label="Entrypoints" placeholder="" {disabled} />
		{/if}

		<!-- Middlewares -->
		{#if router.provider === 'http'}
			<div class="grid grid-cols-4 items-center gap-1" class:hidden={router.protocol === 'udp'}>
				<Label for="middlewares" class="mr-2 text-right">Middlewares</Label>
				<Select.Root
					multiple={true}
					selected={getSelectedMiddlewares(router)}
					onSelectedChange={(value) => toggleMiddleware(router, value, false)}
				>
					<Select.Trigger class="col-span-3">
						<Select.Value placeholder="Select a middleware" />
					</Select.Trigger>
					<Select.Content>
						{#each $middlewares as middleware}
							{#if router.protocol === middleware.protocol}
								<Select.Item value={middleware.name}>
									{middleware.name}
								</Select.Item>
							{/if}
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{:else}
			<ArrayInput
				bind:items={router.middlewares}
				label="Middlewares"
				placeholder="None"
				{disabled}
			/>
		{/if}

		<!-- DNS Provider -->
		{#if $provider}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label for="provider" class="mr-2 text-right">DNS Provider</Label>
				<Select.Root
					selected={getSelectedDNSProvider(router)}
					onSelectedChange={(value) => toggleDNSProvider(router, value, false)}
				>
					<Select.Trigger class="col-span-3">
						<Select.Value placeholder="Select a dns provider" />
					</Select.Trigger>
					<Select.Content>
						<Select.Item value={0} label="">None</Select.Item>
						{#each $provider as provider}
							<Select.Item value={provider.id} class="flex items-center gap-2">
								{provider.name} ({provider.type})
								{#if provider.isActive}
									<iconify-icon icon="fa6-solid:star" class="text-yellow-400" />
								{/if}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{/if}

		<!-- CertResolver -->
		{#if router.provider === 'http' && router.protocol !== 'udp'}
			<div class="space-y-0.5">
				<div class="grid grid-cols-4 items-center gap-1">
					<Label for="certresolver" class="mr-2 text-right">CertResolver</Label>
					<Input
						id="certresolver"
						name="certresolver"
						type="text"
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						bind:value={router.tls.certResolver}
						placeholder="Certificate resolver"
					/>
				</div>
				<div class="grid grid-cols-4 flex-wrap items-center gap-1">
					<div></div>
					<div class="col-span-3">
						{#each getCertResolver() || [] as resolver}
							{#if router.tls?.certResolver !== resolver}
								<div on:click={() => setCertResolver(resolver)} aria-hidden>
									<Badge>{resolver}</Badge>
								</div>
							{/if}
						{/each}
					</div>
				</div>
			</div>
		{/if}

		<!-- Rule -->
		{#if router.protocol === 'http' || router.protocol === 'tcp'}
			{#if errors.rule}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.rule}</div>
			{/if}
			<RuleEditor bind:rule={router.rule} bind:type={router.protocol} {disabled} />
		{/if}
	</Card.Content>
</Card.Root>
