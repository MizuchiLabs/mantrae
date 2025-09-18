<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { type Service } from '$lib/gen/mantrae/v1/service_pb';
	import {
		type Cookie,
		type Service as HTTPService,
		type Server,
		type ServersLoadBalancer
	} from '$lib/gen/zen/traefik-schemas';
	import { Plus, Trash } from '@lucide/svelte';
	import { marshalConfig } from '$lib/types';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import { serversTransportClient } from '$lib/api';
	import { profile } from '$lib/stores/profile';

	interface Props {
		service: Service;
	}
	let { service = $bindable() }: Props = $props();

	let servers = $state([{ url: '' }] as Server[]);
	let passHostHeader = $state(true);
	let serversTransport = $state('');
	let sticky = $state(false);
	let cookie = $state({} as Cookie);

	function updateConfig() {
		const config = {} as HTTPService;
		config.loadBalancer = {} as ServersLoadBalancer;
		config.loadBalancer.servers = servers;
		config.loadBalancer.passHostHeader = passHostHeader;
		config.loadBalancer.sticky = sticky ? { cookie } : undefined;
		config.loadBalancer.serversTransport = serversTransport;
		service.config = marshalConfig(config);
	}

	$effect(() => {
		if (service.config) {
			let config = service.config as HTTPService;
			servers = config.loadBalancer?.servers || [];
			passHostHeader = config.loadBalancer?.passHostHeader ?? true;
			sticky = config.loadBalancer?.sticky?.cookie !== undefined;
			cookie = config.loadBalancer?.sticky?.cookie ?? {};
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

	<div class="flex items-center justify-between rounded-lg border p-3">
		<div class="space-y-1">
			<Label class="flex items-center gap-1 text-sm">Pass Host Header</Label>
			<p class="text-xs text-muted-foreground">Forward client host header to server</p>
		</div>
		<CustomSwitch bind:checked={passHostHeader} onCheckedChange={updateConfig} size="md" />
	</div>

	<div class="grid grid-cols-4 items-center gap-3 rounded-lg border p-3">
		<div class="col-span-3 space-y-1">
			<Label class="flex items-center gap-1 text-sm">Sticky Cookie</Label>
			<p class="text-xs text-muted-foreground">
				Defines a Set-Cookie header is set on the initial response to let the client know which
				server handles the first response.
			</p>
		</div>
		<div class="justify-self-end">
			<CustomSwitch bind:checked={sticky} onCheckedChange={updateConfig} size="md" />
		</div>
	</div>

	{#if sticky}
		<div class="flex flex-col gap-2">
			<Label for="cookie-name">Cookie Name</Label>
			<Input
				id="cookie-name"
				type="text"
				bind:value={cookie.name}
				oninput={updateConfig}
				placeholder="session"
			/>
		</div>
	{/if}

	{#await serversTransportClient.listServersTransports({ profileId: profile.id }) then value}
		<div class="flex flex-col gap-2">
			<Label for="servers-transport">Servers Transport</Label>
			<Select.Root type="single" bind:value={serversTransport} onValueChange={updateConfig}>
				<Select.Trigger class="w-full">
					<span class="truncate text-left">
						{serversTransport || 'Select servers transport'}
					</span>
				</Select.Trigger>
				<Select.Content>
					<Select.Item value="">
						<span class="truncate">Select servers transport</span>
					</Select.Item>
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
