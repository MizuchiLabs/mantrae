<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { Router as HTTPRouter, RouterTLSConfig } from '$lib/gen/zen/traefik-schemas';
	import { AlertCircle, ExternalLink, Plus, Star } from '@lucide/svelte';
	import { entryPointClient, middlewareClient, routerClient } from '$lib/api';
	import { MiddlewareType } from '$lib/gen/mantrae/v1/middleware_pb';
	import { unmarshalConfig, marshalConfig } from '$lib/types';
	import { profile } from '$lib/stores/profile';
	import { formatArrayDisplay } from '$lib/utils';
	import { onMount } from 'svelte';

	let { router = $bindable() }: { router: Router } = $props();

	let certResolvers: string[] = $state([]);
	let config = $state<HTTPRouter>(unmarshalConfig(router.config) as HTTPRouter);

	$effect(() => {
		if (config) router.config = marshalConfig(config);
	});

	onMount(async () => {
		// Get cert resolvers
		const response = await routerClient.listRouters({
			profileId: router.profileId ?? profile.id,
			type: router.type,
			limit: -1n,
			offset: 0n
		});
		const resolverSet = new Set(
			response.routers
				.filter((r) => {
					let tmp = unmarshalConfig(r.config) as HTTPRouter;
					if (!tmp?.tls?.certResolver) return false;
					return true;
				})
				.map((r) => {
					let tmp = unmarshalConfig(r.config) as HTTPRouter;
					return tmp.tls?.certResolver ?? '';
				})
				.filter(Boolean)
		);
		certResolvers = Array.from(resolverSet);

		// Set default entrypoint
		entryPointClient
			.listEntryPoints({ profileId: profile.id, limit: -1n, offset: 0n })
			.then((data) => {
				let defaultEntryPoint = data.entryPoints.find((e) => e.isDefault);
				if (defaultEntryPoint) config.entryPoints = [defaultEntryPoint.name];
			});
		middlewareClient
			.listMiddlewares({ profileId: profile.id, type: MiddlewareType.HTTP, limit: -1n, offset: 0n })
			.then((data) => {
				let defaultMiddleware = data.middlewares.find((m) => m.isDefault);
				if (defaultMiddleware) config.middlewares = [defaultMiddleware.name];
			});
	});
</script>

<div class="flex flex-col gap-3">
	<!-- Entrypoints -->
	{#await entryPointClient.listEntryPoints( { profileId: profile.id, limit: -1n, offset: 0n } ) then value}
		{#if !value.entryPoints.length}
			<Alert.Root class="border-dashed">
				<AlertCircle class="h-4 w-4" />
				<Alert.Title>No entrypoints found</Alert.Title>
				<Alert.Description class="flex items-center justify-between">
					<span>Create an entrypoint to get started.</span>
					<Button
						variant="outline"
						size="sm"
						href="/entrypoints"
						class="ml-4 flex shrink-0 items-center gap-2"
					>
						<Plus />
						Create Entrypoint
						<ExternalLink />
					</Button>
				</Alert.Description>
			</Alert.Root>
		{:else}
			<div class="flex flex-col gap-2">
				<Label class="mr-2">Entrypoints</Label>
				<Select.Root type="multiple" bind:value={config.entryPoints}>
					<Select.Trigger class="w-full">
						<span class="truncate text-left">
							{formatArrayDisplay(config.entryPoints) || 'Select entrypoints'}
						</span>
					</Select.Trigger>
					<Select.Content>
						{#each value.entryPoints as e (e.id)}
							<Select.Item value={e.name}>
								<div class="flex items-center gap-2">
									<span class="truncate">{e.name}</span>
									{#if e.isDefault}
										<Star size="1rem" class="text-yellow-300" />
									{/if}
								</div>
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{/if}
	{/await}

	<!-- Middlewares -->
	{#await middlewareClient.listMiddlewares( { profileId: profile.id, type: MiddlewareType.HTTP, limit: -1n, offset: 0n } ) then value}
		{#if !value.middlewares.length}
			<Alert.Root class="border-dashed">
				<AlertCircle class="h-4 w-4" />
				<Alert.Title>No middlewares found</Alert.Title>
				<Alert.Description class="flex items-center justify-between">
					<span>Create middlewares to add authentication, rate limiting, and more.</span>
					<Button
						variant="outline"
						size="sm"
						href="/middlewares"
						class="ml-4 flex shrink-0 items-center gap-2"
					>
						<Plus />
						Create Middleware
						<ExternalLink />
					</Button>
				</Alert.Description>
			</Alert.Root>
		{:else}
			<div class="flex flex-col gap-2">
				<Label class="mr-2">Middlewares</Label>
				<Select.Root type="multiple" bind:value={config.middlewares}>
					<Select.Trigger class="w-full" disabled={!value.middlewares.length}>
						<span class="truncate text-left">
							{formatArrayDisplay(config.middlewares) || 'Select middlewares'}
						</span>
					</Select.Trigger>
					<Select.Content>
						{#each value.middlewares as middleware (middleware.name)}
							<Select.Item value={middleware.name}>
								<span class="truncate">{middleware.name}</span>
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{/if}
	{/await}

	<div class="grid w-full grid-cols-1 gap-4 sm:grid-cols-3 sm:gap-2">
		<!-- TLS Configuration -->
		<div class="flex flex-col gap-2 sm:col-span-2">
			<Label for="certResolver" class="mr-2">Certificate Resolver</Label>
			<div class="col-span-3">
				<Input
					value={config.tls?.certResolver}
					name="certResolver"
					placeholder="letsencrypt"
					class="truncate"
					oninput={(e) => {
						const input = e.target as HTMLInputElement;
						if (!input.value) {
							delete config.tls;
							return;
						}

						if (!config.tls) config.tls = {} as RouterTLSConfig;
						config.tls.certResolver = input.value;
					}}
				/>
				<div class="flex flex-wrap gap-1">
					{#each certResolvers as resolver (resolver)}
						{#if resolver !== config.tls?.certResolver}
							<Badge
								onclick={() => {
									if (!config.tls) config.tls = {} as RouterTLSConfig;
									config.tls.certResolver = resolver;
								}}
								class="mt-1 cursor-pointer"
							>
								{resolver}
							</Badge>
						{/if}
					{/each}
				</div>
			</div>
		</div>

		<!-- Priority -->
		<div class="flex flex-col gap-2 sm:col-span-1">
			<Label for="priority" class="mr-2">Priority</Label>
			<Input
				id="priority"
				type="number"
				bind:value={config.priority}
				placeholder="0"
				min="0"
				max="1000"
			/>
		</div>
	</div>

	<!-- Rule -->
	{#if router.type === RouterType.HTTP}
		<RuleEditor bind:rule={config.rule} bind:type={router.type} />
	{/if}
</div>
