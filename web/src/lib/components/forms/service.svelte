<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { type Router, type Service } from '$lib/types/router';
	import { Plus, Trash } from 'lucide-svelte';

	interface Props {
		service: Service;
		router: Router;
		mode: 'create' | 'edit';
	}

	let { service = $bindable(), router, mode }: Props = $props();

	let routerProvider = $derived(router.name ? router.name?.split('@')[1] : 'http');
	let disabled = $derived(routerProvider !== 'http' && mode === 'edit');

	let passHostHeader = $state(true);
	let servers = $state(['']);

	function update() {
		if (!service.loadBalancer) {
			service.loadBalancer = { servers: [] };
		}

		service.loadBalancer.servers = servers.map((s) =>
			router.protocol === 'http' ? { url: s } : { address: s }
		);
		service.loadBalancer.passHostHeader = passHostHeader;
	}

	function addItem() {
		servers = [...servers, ''];
	}
	function removeItem(index: number) {
		if (index < 1) return;
		servers = servers.filter((_, i) => i !== index);
		update();
	}

	$effect(() => {
		if (service.loadBalancer?.servers) {
			servers = service.loadBalancer.servers
				.map((s) => (router.protocol === 'http' ? s.url : s.address))
				.filter((s): s is string => s !== undefined);
		}
		if (service.loadBalancer?.passHostHeader !== undefined) {
			passHostHeader = service.loadBalancer.passHostHeader;
		}
	});
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Service Configuration</Card.Title>
		<Card.Description>
			Configure your
			<b>{router.service?.split('@')[0]}</b>
			service settings
		</Card.Description>
	</Card.Header>

	<Card.Content>
		<form class="space-y-4">
			{#if router.protocol === 'http'}
				<div class="flex flex-row items-center gap-4">
					<Label for="passHostHeader" class="text-right">Pass Host Header</Label>
					<Switch
						id="passHostHeader"
						class="col-span-3"
						bind:checked={passHostHeader}
						onCheckedChange={update}
						{disabled}
					/>
				</div>
			{/if}

			<div class="flex w-full flex-col gap-2">
				<Label for="servers">Server Endpoints</Label>
				<!-- eslint-disable-next-line -->
				{#each servers || [] as _, i}
					<div class="flex gap-2">
						<Input
							type="text"
							bind:value={servers[i]}
							placeholder={router.protocol === 'http' ? 'http://127.0.0.1:8080' : '127.0.0.1:8080'}
							oninput={update}
						/>
						<Button variant="ghost" type="button" size="icon" onclick={() => removeItem(i)}>
							<Trash />
						</Button>
					</div>
				{/each}
			</div>
			<Button type="button" variant="outline" onclick={addItem}>
				<Plus />
				Add Server
			</Button>
		</form>
	</Card.Content>
</Card.Root>
