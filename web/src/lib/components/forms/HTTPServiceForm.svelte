<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ServiceType, type Service } from '$lib/gen/mantrae/v1/service_pb';
	import type { Service as HTTPService } from '$lib/gen/zen/traefik-schemas';
	import { Plus, Trash } from '@lucide/svelte';
	import { marshalConfig } from '$lib/types';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';

	interface Props {
		service: Service;
	}
	let { service = $bindable() }: Props = $props();

	let config = $state({} as HTTPService);

	$effect(() => {
		if (service.config) {
			config = service.config as HTTPService;
		}
		if (!config.loadBalancer) config.loadBalancer = { passHostHeader: true };
		if (!config.loadBalancer.servers) config.loadBalancer.servers = [];
		if (config.loadBalancer.servers.length === 0) {
			config.loadBalancer.servers = [{ url: '' }];
		}
	});
</script>

<div class="flex flex-col gap-3">
	<div class="flex items-center justify-between rounded-lg border p-3">
		<div class="space-y-1">
			<Label class="flex items-center gap-1 text-sm">Pass Host Header</Label>
			<p class="text-muted-foreground text-xs">Forward client host header to server</p>
		</div>
		<CustomSwitch
			checked={config.loadBalancer?.passHostHeader ?? true}
			onCheckedChange={(value) => {
				if (!config.loadBalancer) config.loadBalancer = { passHostHeader: true };
				config.loadBalancer.passHostHeader = value;
				service.config = marshalConfig(config);
			}}
			size="md"
		/>
	</div>

	<div class="flex flex-col gap-2">
		<Label for="servers">Server Endpoints</Label>
		{#each config.loadBalancer?.servers || [] as server, i (i)}
			<div class="flex gap-2">
				<Input
					type="text"
					value={server.url}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						server.url = input.value;
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
						if (!config.loadBalancer) config.loadBalancer = { passHostHeader: true };
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
			if (!config.loadBalancer) config.loadBalancer = { passHostHeader: true };
			if (!config.loadBalancer.servers) config.loadBalancer.servers = [];
			config.loadBalancer.servers = [...config.loadBalancer.servers, { url: '' }];
			service.config = marshalConfig(config);
		}}
	>
		<Plus />
		Add Server
	</Button>
</div>
