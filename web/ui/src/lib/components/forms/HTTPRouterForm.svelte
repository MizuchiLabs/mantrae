<script lang="ts">
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { Router as HTTPRouter, RouterTLSConfig } from '$lib/gen/zen/traefik-schemas';
	import { marshalConfig, unmarshalConfig } from '$lib/types';
	import { formatArrayDisplay } from '$lib/utils';
	import { CircleAlert, ExternalLink, Plus, Star } from '@lucide/svelte';
	import { SvelteSet } from 'svelte/reactivity';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import { CustomSwitch } from '../ui/custom-switch';
	import { router } from '$lib/api/router.svelte';
	import { entrypoint } from '$lib/api/entrypoints.svelte';
	import { middleware } from '$lib/api/middleware.svelte';

	interface Props {
		data: Router;
	}
	let { data = $bindable() }: Props = $props();

	let config = $state<HTTPRouter>(unmarshalConfig(data.config) as HTTPRouter);
	let certResolvers = new SvelteSet();

	const routerList = router.list();
	const epList = entrypoint.list();
	const mwList = middleware.list();

	$effect(() => {
		if (config) data.config = marshalConfig(config);
	});
	$effect(() => {
		if (routerList.isSuccess && routerList.data) {
			routerList.data.forEach((r) => {
				if (r.type === ProtocolType.HTTP) {
					let tmp = unmarshalConfig(r.config) as HTTPRouter;
					if (tmp?.tls?.certResolver) certResolvers.add(tmp.tls?.certResolver);
				}
			});
		}
	});
	$effect(() => {
		if (epList.isSuccess && epList.data && !data.id) {
			config.entryPoints = [epList.data.find((ep) => ep.isDefault)?.name ?? ''];
		}
	});
	$effect(() => {
		if (mwList.isSuccess && mwList.data && !data.id) {
			config.middlewares = [
				mwList.data.find((mw) => mw.isDefault && mw.type === data.type)?.name ?? ''
			];
		}
	});
</script>

<div class="flex flex-col gap-3">
	<!-- Entrypoints -->
	{#if epList.data?.length === 0}
		<Alert.Root class="border-dashed">
			<CircleAlert class="h-4 w-4" />
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
					{#each epList.data || [] as e (e.id)}
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

	<!-- Middlewares -->
	{#if mwList.data?.length === 0}
		<Alert.Root class="border-dashed">
			<CircleAlert class="h-4 w-4" />
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
				<Select.Trigger class="w-full" disabled={mwList.data?.length === 0}>
					<span class="truncate text-left">
						{formatArrayDisplay(config.middlewares) || 'Select middlewares'}
					</span>
				</Select.Trigger>
				<Select.Content>
					{#each mwList.data || [] as middleware (middleware.id)}
						<Select.Item value={middleware.name}>
							<span class="truncate">{middleware.name}</span>
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
	{/if}

	<div class="flex items-center justify-between rounded-lg border p-3">
		<div class="space-y-1">
			<Label class="text-sm">Enable TLS</Label>
			<p class="text-xs text-muted-foreground">Automatically enabled when a certResolver is set</p>
		</div>

		<CustomSwitch
			checked={config.tls !== undefined}
			onCheckedChange={(checked) => {
				if (!config.tls && checked) config.tls = {} as RouterTLSConfig;
				if (!checked) delete config.tls;
			}}
		/>
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
	{#if data.type === ProtocolType.HTTP}
		<RuleEditor bind:rule={config.rule} bind:type={data.type} />
	{/if}
</div>
