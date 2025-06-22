<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ServiceType, type Service } from '$lib/gen/mantrae/v1/service_pb';
	import type { Service as HttpService } from '$lib/gen/tygo/dynamic';
	import { Plus, Trash } from '@lucide/svelte';
	import { marshalConfig } from '$lib/types';

	interface Props {
		service: Service;
	}
	let { service = $bindable() }: Props = $props();

	let config = $state({} as HttpService);

	$effect(() => {
		if (service.config) {
			config = service.config as HttpService;
		}
		if (!config.loadBalancer) config.loadBalancer = {};
		if (!config.loadBalancer.servers) config.loadBalancer.servers = [];
		if (config.loadBalancer.servers.length === 0) {
			config.loadBalancer.servers = [{ url: '' }];
		}
	});
</script>

<div class="flex flex-col gap-3">
	<div class="flex flex-col gap-2">
		<Label for="passHostHeader" class="text-right">Pass Host Header</Label>
		<Switch
			id="passHostHeader"
			class="col-span-3"
			checked={config.loadBalancer?.passHostHeader ?? true}
			onCheckedChange={(value) => {
				if (!config.loadBalancer) config.loadBalancer = {};
				config.loadBalancer.passHostHeader = value;
				service.config = marshalConfig(config);
			}}
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
			config.loadBalancer.servers = [...config.loadBalancer.servers, { url: '' }];
			service.config = marshalConfig(config);
		}}
	>
		<Plus />
		Add Server
	</Button>
</div>
