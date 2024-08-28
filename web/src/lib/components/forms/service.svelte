<script lang="ts">
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Service } from '$lib/types/config';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let service: Service;
	export let disabled = false;
	let servers = service.loadBalancer?.servers?.map((s) => s.url ?? '') ?? [];
	let tcpServers = service.tcpLoadBalancer?.servers?.map((s) => s.address ?? '') ?? [];
	let udpServers = service.udpLoadBalancer?.servers?.map((s) => s.address ?? '') ?? [];
	$: servers, tcpServers, udpServers, onChange();

	const onChange = () => {
		if (service.loadBalancer === undefined) service.loadBalancer = { servers: [] };
		service.loadBalancer.servers = servers.map((s) => {
			return { url: s };
		});
		if (service.tcpLoadBalancer === undefined) service.tcpLoadBalancer = { servers: [] };
		service.tcpLoadBalancer.servers = tcpServers.map((s) => {
			return { address: s };
		});
		if (service.udpLoadBalancer === undefined) service.udpLoadBalancer = { servers: [] };
		service.udpLoadBalancer.servers = udpServers.map((s) => {
			return { address: s };
		});
	};
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Service</Card.Title>
		<Card.Description>
			Make changes to your Service here. Click save when you're done.
		</Card.Description>
	</Card.Header>
	<Card.Content class="space-y-2">
		{#if service.loadBalancer}
			{#if service.serviceType === 'http'}
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="passHostHeader" class="text-right">Pass Host Header</Label>
					<Switch
						id="passHostHeader"
						class="col-span-3"
						bind:checked={service.loadBalancer.passHostHeader}
						{disabled}
					/>
				</div>
				<ArrayInput bind:items={servers} label="Servers" placeholder="192.168.1.1" {disabled} />
			{/if}
			{#if service.serviceType === 'tcp'}
				<ArrayInput bind:items={tcpServers} label="Servers" placeholder="192.168.1.1" {disabled} />
			{/if}
			{#if service.serviceType === 'udp'}
				<ArrayInput bind:items={udpServers} label="Servers" placeholder="192.168.1.1" {disabled} />
			{/if}
		{/if}
	</Card.Content>
</Card.Root>
