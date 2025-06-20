<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import { ServiceType, type Service } from '$lib/gen/mantrae/v1/service_pb';
	import type { Service as HttpService } from '$lib/gen/tygo/dynamic';
	import { marshalConfig, unmarshalConfig } from '$lib/types';
	import { Plus, Trash } from '@lucide/svelte';

	interface Props {
		service: Service;
		router: Router;
	}
	let { service = $bindable(), router = $bindable() }: Props = $props();

	let config = $state(unmarshalConfig(service.config) as HttpService);

	$effect(() => {
		if (config) {
			if (!config.loadBalancer) config.loadBalancer = {};
			if (!config.loadBalancer.servers) config.loadBalancer.servers = [];
			if (config.loadBalancer.servers.length === 0) {
				config.loadBalancer.servers = [{ url: '' }];
			}

			service.config = marshalConfig(config);
		}
		if (router.name) {
			service.name = router.name;
		}
		if (router.type) {
			switch (router.type) {
				case RouterType.HTTP:
					service.type = ServiceType.HTTP;
					break;
				case RouterType.TCP:
					service.type = ServiceType.TCP;
					break;
				case RouterType.UDP:
					service.type = ServiceType.UDP;
					break;
			}
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
		}}
	>
		<Plus />
		Add Server
	</Button>
</div>
