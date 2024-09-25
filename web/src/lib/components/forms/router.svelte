<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { badgeVariants } from '$lib/components/ui/badge';
	import {
		routers,
		entrypoints,
		middlewares,
		provider,
		toggleEntrypoint,
		toggleMiddleware,
		toggleDNSProvider
	} from '$lib/api';
	import { newRouter, type Router } from '$lib/types/config';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import { z } from 'zod';
	import type { Selected } from 'bits-ui';

	export let router: Router;
	export let disabled = false;

	let errors: Record<any, string[] | undefined> = {};
	const formSchema = z.object({
		name: z.string({ required_error: 'Name is required' }).min(3).max(255),
		routerType: z
			.string()
			.toLowerCase()
			.regex(/^(http|tcp|udp)$/),
		tls: z.object({
			certResolver: z.string().trim().optional()
		})
	});
	export const validate = () => {
		try {
			formSchema.parse({
				name: router.name,
				routerType: router.routerType,
				tls: router.tls
			});
			errors = {};
			return true;
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
			return false;
		}
	};

	let routerType: Selected<string> | undefined = router.routerType
		? { label: router.routerType.toUpperCase(), value: router.routerType }
		: { label: 'HTTP', value: 'http' };
	const changeType = (serviceType: Selected<string> | undefined) => {
		if (serviceType === undefined) return;
		router = newRouter();
		router.routerType = serviceType.value;
		routerType = { label: serviceType.label || '', value: serviceType.value };
	};

	const getSelectedEntrypoints = (router: Router): Selected<unknown>[] => {
		let list = router?.entrypoints?.map((entrypoint) => {
			return { value: entrypoint, label: entrypoint };
		});
		return list ?? [];
	};
	const getSelectedMiddlewares = (router: Router): Selected<unknown>[] => {
		let list = router?.middlewares?.map((middleware) => {
			return { value: middleware, label: middleware };
		});
		return list ?? [];
	};
	const getSelectedDNSProvider = (router: Router): Selected<unknown> | undefined => {
		return router?.dnsProvider
			? { value: router.dnsProvider, label: router.dnsProvider }
			: undefined;
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
	let isNameTaken = false;
	let routerCompare = $routers.filter((r) => r.name !== router.name);
	$: isNameTaken = routerCompare.some(
		(r) => r.name.split('@')[0].toLowerCase() === router.name.split('@')[0].toLowerCase()
	);
</script>

<Card.Root>
	<Card.Header>
		<Card.Title class="flex items-center justify-between gap-1">
			<span>Router</span>
			<div>
				<Badge variant="secondary" class="bg-blue-400">
					Type: {router.routerType}
				</Badge>
				<Badge variant="secondary" class="bg-green-400">
					Provider: {router.provider}
				</Badge>
			</div>
		</Card.Title>
	</Card.Header>
	<Card.Content class="space-y-2">
		<!-- Type -->
		{#if router.provider === 'http'}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label for="current" class="mr-2 text-right">Type</Label>
				<Select.Root onSelectedChange={changeType} selected={routerType}>
					<Select.Trigger class="col-span-3">
						<Select.Value placeholder="Select a type" />
					</Select.Trigger>
					<Select.Content>
						<Select.Item value="http" label="HTTP">HTTP</Select.Item>
						<Select.Item value="tcp" label="TCP">TCP</Select.Item>
						<Select.Item value="udp" label="UDP">UDP</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>
		{/if}

		<!-- Name -->
		<div class="grid grid-cols-4 items-center gap-1">
			<Label for="name" class="mr-2 text-right">Name</Label>
			<Input
				id="name"
				name="name"
				type="text"
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				bind:value={router.name}
				placeholder="Name of the router"
				required
				{disabled}
			/>
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
			<ArrayInput bind:items={router.entrypoints} label="Entrypoints" placeholder="" {disabled} />
		{/if}

		<!-- Middlewares -->
		{#if router.provider === 'http'}
			<div class="grid grid-cols-4 items-center gap-1" class:hidden={router.routerType === 'udp'}>
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
							{#if router.routerType === middleware.middlewareType}
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
						<Select.Item value="" label="">None</Select.Item>
						{#each $provider as provider}
							<Select.Item value={provider.name} class="flex items-center gap-2">
								{provider.name} ({provider.type})
								{#if provider.is_active}
									<iconify-icon icon="fa6-solid:star" class="text-yellow-400" />
								{/if}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{/if}

		<!-- CertResolver -->
		{#if router.provider === 'http'}
			<div class:hidden={router.routerType === 'udp'} class="space-y-0.5">
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
							{#if router.tls.certResolver !== resolver}
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
		{#if router.routerType === 'http' || router.routerType === 'tcp'}
			<RuleEditor bind:rule={router.rule} {disabled} />
		{/if}
	</Card.Content>
</Card.Root>
