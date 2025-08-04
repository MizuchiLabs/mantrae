<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { RouterTCPTLSConfig, TCPRouter } from '$lib/gen/zen/traefik-schemas';
	import { entryPoints, middlewares, routers } from '$lib/stores/realtime';
	import { marshalConfig, unmarshalConfig } from '$lib/types';
	import { Star } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { SvelteSet } from 'svelte/reactivity';
	import RuleEditor from '../utils/ruleEditor.svelte';

	let { router = $bindable() }: { router: Router } = $props();

	let config = $state(unmarshalConfig(router.config) as TCPRouter);
	let certResolvers = new SvelteSet();

	$effect(() => {
		if (config) router.config = marshalConfig(config);
	});

	onMount(async () => {
		$routers.forEach((r) => {
			if (r.type === router.type) {
				let tmp = unmarshalConfig(r.config) as TCPRouter;
				if (tmp?.tls?.certResolver) certResolvers.add(tmp.tls?.certResolver);
			}
		});

		// Set default entrypoint
		let defaultEntryPoint = $entryPoints.find((e) => e.isDefault);
		if (defaultEntryPoint) config.entryPoints = [defaultEntryPoint.name];

		let defaultMiddleware = $middlewares.find((m) => m.isDefault && m.type === ProtocolType.TCP);
		if (defaultMiddleware) config.middlewares = [defaultMiddleware.name];
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
				{#each $entryPoints || [] as e (e.id)}
					<Select.Item value={e.name}>
						<div class="flex items-center gap-2">
							{e.name}
							{#if e.isDefault}
								<Star size="1rem" class="text-yellow-300" />
							{/if}
						</div>
					</Select.Item>
				{/each}
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
				{#each $middlewares || [] as middleware (middleware.id)}
					<Select.Item value={middleware.name}>
						{middleware.name}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>

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
									if (resolver) config.tls.certResolver = resolver.toString();
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
	{#if router.type === ProtocolType.TCP}
		<RuleEditor bind:rule={config.rule} bind:type={router.type} />
	{/if}
</div>
