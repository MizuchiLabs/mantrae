<script lang="ts">
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Service } from '$lib/types/config';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import { onMount } from 'svelte';

	export let service: Service;
	export let disabled = false;
	let passHostHeader = service?.loadBalancer?.passHostHeader ?? true;
	let servers: string[] = [];

	const update = () => {
		if (service.loadBalancer === undefined) service.loadBalancer = { servers: [] };
		if (service.tcpLoadBalancer === undefined) service.tcpLoadBalancer = { servers: [] };
		if (service.udpLoadBalancer === undefined) service.udpLoadBalancer = { servers: [] };

		switch (service.serviceType) {
			case 'http':
				service.loadBalancer.servers = servers.map((s) => {
					return { url: s };
				});
				break;
			case 'tcp':
				service.tcpLoadBalancer.servers = servers.map((s) => {
					return { address: s };
				});
				break;
			case 'udp':
				service.udpLoadBalancer.servers = servers.map((s) => {
					return { address: s };
				});
				break;
		}
	};

	const onSwitchChange = () => {
		passHostHeader = !passHostHeader;
		if (service.loadBalancer === undefined) service.loadBalancer = { servers: [] };
		service.loadBalancer.passHostHeader = passHostHeader;
	};

	onMount(() => {
		switch (service.serviceType) {
			case 'http':
				servers = service.loadBalancer?.servers?.map((s) => s.url ?? '') ?? [''];
				break;
			case 'tcp':
				servers = service.tcpLoadBalancer?.servers?.map((s) => s.address ?? '') ?? [''];
				break;
			case 'udp':
				servers = service.udpLoadBalancer?.servers?.map((s) => s.address ?? '') ?? [''];
				break;
		}
	});
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Service</Card.Title>
		<Card.Description>
			Make changes to your Service here. Click save when you're done.
		</Card.Description>
	</Card.Header>
	<Card.Content class="space-y-2">
		{#if service.serviceType === 'http'}
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="passHostHeader" class="text-right">Pass Host Header</Label>
				<Switch
					id="passHostHeader"
					class="col-span-3"
					bind:checked={passHostHeader}
					onCheckedChange={onSwitchChange}
					{disabled}
				/>
			</div>
		{/if}
		<ArrayInput
			bind:items={servers}
			label="Servers"
			placeholder="192.168.1.1"
			on:update={update}
			{disabled}
		/>
	</Card.Content>
</Card.Root>
