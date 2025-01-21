<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Router, type Service } from '$lib/types/router';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	interface Props {
		service: Service;
		router: Router;
		disabled?: boolean;
	}

	let { service = $bindable({} as Service), router, disabled = false }: Props = $props();

	let passHostHeader = $derived(service.loadBalancer?.passHostHeader ?? true);
	let servers = $state(
		service.loadBalancer?.servers
			?.map((s) => (router.type === 'http' ? s.url : s.address))
			.filter((s): s is string => s !== undefined) ?? []
	);

	function updateServers(newServers: string[]) {
		if (!service.loadBalancer) {
			service.loadBalancer = { servers: [] };
		}

		service.type = router.type;
		service.loadBalancer.servers = newServers.map((s) =>
			router.type === 'http' ? { url: s } : { address: s }
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
		<Card.Description>Configure your {router.type} service settings</Card.Description>
	</Card.Header>

	<Card.Content class="flex flex-col gap-2">
		{#if router.type === 'http'}
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
			placeholder={router.type === 'http' ? 'http://192.168.1.1:8080' : '192.168.1.1:8080'}
			{disabled}
		/>
	</Card.Content>
</Card.Root>
