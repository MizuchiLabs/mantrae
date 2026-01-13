<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { type Service } from '$lib/gen/mantrae/v1/service_pb';
	import { ChevronDown, Plus, Trash } from '@lucide/svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import { marshalConfig } from '$lib/types';
	import type {
		TCPServer,
		TCPServerHealthCheck,
		TCPServersLoadBalancer,
		TCPService
	} from '$lib/gen/zen/traefik-schemas';
	import { CustomSwitch } from '../ui/custom-switch';
	import { Separator } from '../ui/separator';
	import { serversTransportClient } from '$lib/api';
	import { profile } from '$lib/stores/profile';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';

	interface Props {
		service: Service;
	}
	let { service = $bindable() }: Props = $props();

	let servers = $state([{ address: '' }] as TCPServer[]);
	let serversTransport = $state('');
	let healthcheck = $state(false);
	let healthOptions = $state({
		port: 8080,
		interval: '10s',
		unhealthyInterval: '5s',
		timeout: '2s'
	} as TCPServerHealthCheck);

	function updateConfig() {
		const config = {} as TCPService;
		config.loadBalancer = {} as TCPServersLoadBalancer;
		config.loadBalancer.servers = servers;
		config.loadBalancer.serversTransport = serversTransport;
		config.loadBalancer.healthCheck = healthcheck ? healthOptions : undefined;
		service.config = marshalConfig(config);
	}

	$effect(() => {
		if (service.config) {
			let config = service.config as TCPService;
			servers = config.loadBalancer?.servers || [];
			serversTransport = config.loadBalancer?.serversTransport ?? '';
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
					bind:value={server.address}
					oninput={updateConfig}
					placeholder="127.0.0.1:8080"
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
		onclick={() => (servers = [...servers, { address: '' }])}
	>
		<Plus />
		Add Server
	</Button>

	<!-- Advanced Options -->
	<Collapsible.Root>
		<Collapsible.Trigger
			class="flex w-full items-center justify-between rounded-lg border p-3 hover:bg-muted/50"
		>
			<div class="flex items-center gap-2">
				<Label class="pointer-events-none text-sm font-medium">Advanced Options</Label>
				{#if healthcheck || serversTransport}
					<span class="rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
						{[healthcheck && 'Health', serversTransport && 'Transport'].filter(Boolean).join(', ')}
					</span>
				{/if}
			</div>
			<ChevronDown class="h-4 w-4 transition-transform [[data-state=open]>&]:rotate-180" />
		</Collapsible.Trigger>

		<Collapsible.Content>
			<div class="mt-2 flex flex-col gap-3 rounded-lg border p-3">
				<!-- Healthcheck -->
				<div class="flex items-center justify-between">
					<div>
						<Label class="text-sm">Healthcheck</Label>
						<p class="text-xs text-muted-foreground">Monitor backend health</p>
					</div>
					<CustomSwitch bind:checked={healthcheck} onCheckedChange={updateConfig} size="md" />
				</div>
				{#if healthcheck}
					<div class="grid grid-cols-3 gap-2">
						<Input
							bind:value={healthOptions.port}
							oninput={updateConfig}
							placeholder="Port (8080)"
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
							class="text-sm "
						/>
					</div>
				{/if}

				<Separator />

				<!-- Servers Transport -->
				{#await serversTransportClient.listServersTransports( { profileId: profile.id, type: ProtocolType.TCP } ) then value}
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
