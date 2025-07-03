<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { Router as HTTPRouter, RouterTLSConfig } from '$lib/gen/zen/traefik-schemas';
	import { Star } from '@lucide/svelte';
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
	});
</script>

<div class="flex flex-col gap-3">
	<!-- Entrypoints -->
	<div class="flex flex-col gap-2">
		<Label class="mr-2">Entrypoints</Label>
		<Select.Root type="multiple" bind:value={config.entryPoints}>
			<Select.Trigger class="w-full">
				<span class="truncate text-left">
					{formatArrayDisplay(config.entryPoints) || 'Select entrypoints'}
				</span>
			</Select.Trigger>
			<Select.Content>
				{#await entryPointClient.listEntryPoints( { profileId: profile.id, limit: -1n, offset: 0n } ) then value}
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
				{/await}
			</Select.Content>
		</Select.Root>
	</div>

	<!-- Middlewares -->
	<div class="flex flex-col gap-2">
		<Label class="mr-2">Middlewares</Label>
		<Select.Root type="multiple" bind:value={config.middlewares}>
			<Select.Trigger class="w-full">
				<span class="truncate text-left">
					{formatArrayDisplay(config.middlewares) || 'Select middlewares'}
				</span>
			</Select.Trigger>
			<Select.Content>
				{#await middlewareClient.listMiddlewares( { profileId: profile.id, type: MiddlewareType.HTTP, limit: -1n, offset: 0n } ) then value}
					{#each value.middlewares as middleware (middleware.name)}
						<Select.Item value={middleware.name}>
							<span class="truncate">{middleware.name}</span>
						</Select.Item>
					{/each}
				{/await}
			</Select.Content>
		</Select.Root>
	</div>

	<!-- TLS Configuration -->
	<div class="flex flex-col gap-2">
		<Label class="mr-2">Certificate Resolver</Label>
		<div class="col-span-3">
			<Input
				value={config.tls?.certResolver}
				placeholder="Certificate resolver"
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
			<div class="mt-2 flex max-h-20 flex-wrap gap-1 overflow-y-auto">
				{#each certResolvers as resolver (resolver)}
					{#if resolver !== config.tls?.certResolver}
						<Badge
							onclick={() => {
								if (!config.tls) config.tls = {} as RouterTLSConfig;
								config.tls.certResolver = resolver;
							}}
							class="max-w-32 cursor-pointer truncate text-xs"
						>
							{resolver}
						</Badge>
					{/if}
				{/each}
			</div>
		</div>
	</div>

	<!-- Rule -->
	{#if router.type === RouterType.HTTP}
		<RuleEditor bind:rule={config.rule} bind:type={router.type} />
	{/if}
</div>
