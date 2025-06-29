<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ServiceType, type Service } from '$lib/gen/mantrae/v1/service_pb';
	import { Plus, Trash } from '@lucide/svelte';
	import { marshalConfig } from '$lib/types';
	import type { UDPService } from '$lib/gen/zen/traefik-schemas';

	interface Props {
		service: Service;
	}
	let { service = $bindable() }: Props = $props();

	let config = $state({} as UDPService);

	$effect(() => {
		if (service.config) {
			config = service.config as UDPService;
		}
		if (!config.loadBalancer) config.loadBalancer = {};
		if (!config.loadBalancer.servers) config.loadBalancer.servers = [];
		if (config.loadBalancer.servers.length === 0) {
			config.loadBalancer.servers = [{ address: '' }];
		}
	});
</script>

<div class="flex flex-col gap-3">
	<div class="flex flex-col gap-2">
		<Label for="servers">Server Endpoints</Label>
		{#each config.loadBalancer?.servers || [] as server, i (i)}
			<div class="flex gap-2">
				<Input
					type="text"
					value={server.address}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						server.address = input.value;
						service.config = marshalConfig(config);
					}}
					placeholder={service.type === ServiceType.HTTP
						? 'http://127.0.0.1:8080'
						: '127.0.0.1:8080'}
				/>
				<Button
					variant="ghost"
					size="icon"
					type="button"
					class="text-red-500"
					onclick={() => {
						if (i === 0) return;
						if (!config.loadBalancer) config.loadBalancer = {};
						if (!config.loadBalancer.servers) config.loadBalancer.servers = [];
						config.loadBalancer.servers = config.loadBalancer.servers.filter((_, j) => j !== i);
						service.config = marshalConfig(config);
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
		onclick={() => {
			if (!config.loadBalancer) config.loadBalancer = {};
			if (!config.loadBalancer.servers) config.loadBalancer.servers = [];
			config.loadBalancer.servers = [...config.loadBalancer.servers, { address: '' }];
			service.config = marshalConfig(config);
		}}
	>
		<Plus />
		Add Server
	</Button>
</div>
