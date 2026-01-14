<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import { type Service } from '$lib/gen/mantrae/v1/service_pb';
	import {
		type Cookie,
		type Service as HTTPService,
		type Server,
		type ServerHealthCheck,
		type ServersLoadBalancer
	} from '$lib/gen/zen/traefik-schemas';
	import { ChevronDown, Plus, Trash } from '@lucide/svelte';
	import { marshalConfig } from '$lib/types';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import { serversTransportClient } from '$lib/api';
	import { profile } from '$lib/stores/profile';
	import { Separator } from '../ui/separator';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';

	interface Props {
		service: Service;
	}
	let { service = $bindable() }: Props = $props();

	let servers = $state([{ url: '' }] as Server[]);
	let passHostHeader = $state(true);
	let serversTransport = $state('');
	let sticky = $state(false);
	let cookie = $state({} as Cookie);
	let healthcheck = $state(false);
	let healthOptions = $state({
		path: '/health',
		status: 200,
		interval: '10s',
		unhealthyInterval: '5s',
		timeout: '2s',
		followRedirects: false
	} as ServerHealthCheck);

	function updateConfig() {
		const config = {} as HTTPService;
		config.loadBalancer = {} as ServersLoadBalancer;
		config.loadBalancer.servers = servers;
		config.loadBalancer.passHostHeader = passHostHeader;
		config.loadBalancer.sticky = sticky ? { cookie } : undefined;
		config.loadBalancer.serversTransport = serversTransport;
		config.loadBalancer.healthCheck = healthcheck ? healthOptions : undefined;
		service.config = marshalConfig(config);
	}

	$effect(() => {
		if (service.config) {
			let config = service.config as HTTPService;
			servers = config.loadBalancer?.servers || [];
			passHostHeader = config.loadBalancer?.passHostHeader ?? true;
			sticky = !!config.loadBalancer?.sticky?.cookie;
			cookie = config.loadBalancer?.sticky?.cookie ?? {};
			serversTransport = config.loadBalancer?.serversTransport ?? '';
			healthcheck = !!config.loadBalancer?.healthCheck;
			healthOptions = config.loadBalancer?.healthCheck ?? {};
		}
	});
</script>

<div class="flex flex-col gap-3">
	<div class="flex flex-col gap-2">
		<Label for="servers">Server Endpoints</Label>
		{#each servers || [] as server, i (i)}
			<div class="flex gap-2">
				<Input
					type="text"
					bind:value={server.url}
					oninput={updateConfig}
					placeholder="http://127.0.0.1:8080"
				/>
				<Button
					variant="ghost"
					size="icon"
					type="button"
					class="text-red-500"
					onclick={() => {
						if (i === 0) return;
						servers = servers.filter((_, j) => j !== i);
					}}
				>
					<Trash />
				</Button>
			</div>
		{/each}
	</div>
	<Button
		type="button"
		variant="outline"
		class="w-full"
		onclick={() => (servers = [...servers, { url: '' }])}
	>
		<Plus />
		Add Server
	</Button>

	<!-- TODO: Test other array inputs -->
	<!-- <div class="flex flex-col gap-2"> -->
	<!-- 	<Label for="servers">Server Endpoints</Label> -->
	<!-- 	<ArrayInput -->
	<!-- 		bind:values={servers} -->
	<!-- 		placeholder="http://127.0.0.1:8080" -->
	<!-- 		getLabel={(server) => server.url ?? ''} -->
	<!-- 		getValue={(server) => server.url ?? ''} -->
	<!-- 		createItem={(url) => ({ url })} -->
	<!-- 		onchange={updateConfig} -->
	<!-- 	/> -->
	<!-- </div> -->

	<!-- Advanced Options  -->
	<Collapsible.Root>
		<Collapsible.Trigger
			class="flex w-full items-center justify-between rounded-lg border p-3 hover:bg-muted/50"
		>
			<div class="flex items-center gap-2">
				<Label class="pointer-events-none text-sm font-medium">Advanced Options</Label>
				{#if sticky || healthcheck || serversTransport}
					<span class="rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
						{[sticky && 'Sticky', healthcheck && 'Health', serversTransport && 'Transport']
							.filter(Boolean)
							.join(', ')}
					</span>
				{/if}
			</div>
			<ChevronDown class="h-4 w-4 transition-transform [[data-state=open]>&]:rotate-180" />
		</Collapsible.Trigger>

		<Collapsible.Content>
			<div class="mt-2 flex flex-col gap-3 rounded-lg border p-3">
				<!-- Pass Host Header -->
				<div class="flex items-center justify-between">
					<div>
						<Label class="text-sm">Pass Host Header</Label>
						<p class="text-xs text-muted-foreground">Forward client host header</p>
					</div>
					<CustomSwitch bind:checked={passHostHeader} onCheckedChange={updateConfig} size="md" />
				</div>

				<Separator />

				<!-- Sticky Cookie -->
				<div class="flex items-center justify-between">
					<div>
						<Label class="text-sm">Sticky Cookie</Label>
						<p class="text-xs text-muted-foreground">Session persistence via cookie</p>
					</div>
					<CustomSwitch bind:checked={sticky} onCheckedChange={updateConfig} size="md" />
				</div>
				{#if sticky}
					<Input
						type="text"
						bind:value={cookie.name}
						oninput={updateConfig}
						placeholder="Cookie name (e.g., session)"
						class="text-sm"
					/>
				{/if}

				<Separator />

				<!-- Healthcheck -->
				<div class="flex items-center justify-between">
					<div>
						<Label class="text-sm">Healthcheck</Label>
						<p class="text-xs text-muted-foreground">Monitor backend health</p>
					</div>
					<CustomSwitch bind:checked={healthcheck} onCheckedChange={updateConfig} size="md" />
				</div>
				{#if healthcheck}
					<div class="grid grid-cols-2 gap-2">
						<Input
							bind:value={healthOptions.path}
							oninput={updateConfig}
							placeholder="Path (/health)"
							class="text-sm"
						/>
						<Input
							bind:value={healthOptions.interval}
							oninput={updateConfig}
							placeholder="Interval (10s)"
							class="text-sm"
						/>
						<Input
							bind:value={healthOptions.timeout}
							oninput={updateConfig}
							placeholder="Timeout (2s)"
							class="text-sm"
						/>
						<Input
							type="number"
							bind:value={healthOptions.status}
							oninput={updateConfig}
							placeholder="Status (200)"
							class="text-sm"
						/>
					</div>
				{/if}

				<Separator />

				<!-- Servers Transport -->
				{#await serversTransportClient.listServersTransports( { profileId: profile.id, type: ProtocolType.HTTP } ) then value}
					<div class="flex flex-col gap-2">
						<Label class="text-sm">Servers Transport</Label>
						<Select.Root type="single" bind:value={serversTransport} onValueChange={updateConfig}>
							<Select.Trigger class="w-full">
								<span class="truncate text-left text-sm">
									{serversTransport || 'Default'}
								</span>
							</Select.Trigger>
							<Select.Content>
								<Select.Item value="">Default</Select.Item>
								{#each value.serversTransports || [] as transport (transport.id)}
									<Select.Item value={transport.name}>
										<span class="truncate">{transport.name}</span>
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				{/await}
			</div>
		</Collapsible.Content>
	</Collapsible.Root>
</div>
