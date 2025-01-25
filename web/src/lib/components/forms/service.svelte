<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Router, type Service } from '$lib/types/router';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	interface Props {
		service: Service;
		router: Router;
	}

	let { service = $bindable({} as Service), router }: Props = $props();

	let routerProvider = $derived(router.name ? router.name.split('@')[1] : 'http');
	let disabled = $derived(routerProvider !== 'http');

	let passHostHeader = $derived(service.loadBalancer?.passHostHeader ?? true);
	let servers = $state(
		service.loadBalancer?.servers
			?.map((s) => (router.protocol === 'http' ? s.url : s.address))
			.filter((s): s is string => s !== undefined) ?? []
	);

	function updateServers(newServers: string[]) {
		if (!service.loadBalancer) {
			service.loadBalancer = { servers: [] };
		}

		service.protocol = router.protocol;
		service.loadBalancer.servers = newServers.map((s) =>
			router.protocol === 'http' ? { url: s } : { address: s }
		);
	}

	function updatePassHostHeader(value: boolean) {
		if (!service.loadBalancer) {
			service.loadBalancer = { servers: [] };
		}
		service.loadBalancer.passHostHeader = value;
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Service Configuration</Card.Title>
		<Card.Description>
			Configure your
			<b>{router.service}</b>
			settings
		</Card.Description>
	</Card.Header>

	<Card.Content class="flex flex-col gap-2">
		{#if router.protocol === 'http'}
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="passHostHeader" class="text-right">Pass Host Header</Label>
				<Switch
					id="passHostHeader"
					class="col-span-3"
					checked={passHostHeader}
					onCheckedChange={updatePassHostHeader}
					{disabled}
				/>
			</div>
		{/if}

		<ArrayInput
			bind:items={servers}
			on:update={({ detail }) => updateServers(detail)}
			label="Servers"
			placeholder={router.protocol === 'http' ? 'http://192.168.1.1:8080' : '192.168.1.1:8080'}
			{disabled}
		/>
	</Card.Content>
</Card.Root>
