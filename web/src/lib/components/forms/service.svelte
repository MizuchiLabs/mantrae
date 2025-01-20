<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Router, type Service } from '$lib/types/router';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	interface Props {
		service: Service | undefined;
		router: Router | undefined;
		disabled?: boolean;
	}

	let {
		service = $bindable({} as Service),
		router = $bindable({} as Router),
		disabled = false
	}: Props = $props();

	let passHostHeader = $state(service?.loadBalancer?.passHostHeader ?? true);
	let servers: string[] = $state([]);

	const update = () => {
		if (service.loadBalancer === undefined) service.loadBalancer = { servers: [] };
		service.type = router.type;
		switch (service.type) {
			case 'http':
				service.loadBalancer.servers = servers.map((s) => {
					return { url: s };
				});
				break;
			case 'tcp':
			case 'udp':
				service.loadBalancer.servers = servers.map((s) => {
					return { address: s };
				});
				break;
		}
	};
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Service</Card.Title>
		<Card.Description>Configure your {router.type} service</Card.Description>
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		{#if service.type === 'http'}
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="passHostHeader" class="text-right">Pass Host Header</Label>
				<Switch
					id="passHostHeader"
					class="col-span-3"
					checked={passHostHeader}
					onCheckedChange={() => {
						if (service.loadBalancer === undefined) service.loadBalancer = { servers: [] };
						return (service.loadBalancer.passHostHeader = true);
					}}
					{disabled}
				/>
			</div>
		{/if}
		<ArrayInput
			bind:items={servers}
			label="Servers"
			on:update={update}
			placeholder="http://192.168.1.1:8080"
			{disabled}
		/>
	</Card.Content>
</Card.Root>
