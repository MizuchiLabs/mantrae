<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { RouterTCPTLSConfig, TCPRouter } from '$lib/gen/tygo/dynamic';
	import { Star } from '@lucide/svelte';
	import { entryPointClient, middlewareClient, routerClient } from '$lib/api';
	import { MiddlewareType } from '$lib/gen/mantrae/v1/middleware_pb';
	import { unmarshalConfig, marshalConfig } from '$lib/types';
	import { onMount } from 'svelte';
	import { profile } from '$lib/stores/profile';

	let { router = $bindable() }: { router: Router } = $props();

	let config = $state(unmarshalConfig(router.config) as TCPRouter);
	let certResolvers: string[] = $state([]);

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
					let tmp = unmarshalConfig(r.config) as TCPRouter;
					if (!tmp?.tls?.certResolver) return false;
					return true;
				})
				.map((r) => {
					let tmp = unmarshalConfig(r.config) as TCPRouter;
					return tmp.tls?.certResolver ?? '';
				})
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
				{config.entryPoints?.join(', ') || 'Select entrypoints'}
			</Select.Trigger>
			<Select.Content>
				{#await entryPointClient.listEntryPoints( { profileId: profile.id, limit: -1n, offset: 0n } ) then value}
					{#each value.entryPoints as e (e.id)}
						<Select.Item value={e.name}>
							<div class="flex items-center gap-2">
								{e.name}
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
				{config.middlewares?.join(', ') || 'Select middlewares'}
			</Select.Trigger>
			<Select.Content>
				{#await middlewareClient.listMiddlewares( { profileId: profile.id, type: MiddlewareType.HTTP, limit: -1n, offset: 0n } ) then value}
					{#each value.middlewares as middleware (middleware.name)}
						<Select.Item value={middleware.name}>
							{middleware.name}
						</Select.Item>
					{/each}
				{/await}
			</Select.Content>
		</Select.Root>
	</div>

	<div class="flex flex-col gap-2">
		<Label class="mr-2">Resolver</Label>
		<div class="col-span-3">
			<Input
				value={config.tls?.certResolver}
				placeholder="Certificate resolver"
				oninput={(e) => {
					const input = e.target as HTMLInputElement;
					if (!input.value) {
						delete config.tls;
						return;
					}

					if (!config.tls) config.tls = {} as RouterTCPTLSConfig;
					config.tls.certResolver = input.value;
				}}
			/>
			<div class="flex flex-wrap gap-1">
				{#each certResolvers as resolver (resolver)}
					{#if resolver !== config.tls?.certResolver}
						<Badge
							onclick={() => {
								if (!config.tls) config.tls = {} as RouterTCPTLSConfig;
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

	<!-- Rule -->
	{#if router.type === RouterType.TCP}
		<RuleEditor bind:rule={config.rule} bind:type={router.type} />
	{/if}
</div>
