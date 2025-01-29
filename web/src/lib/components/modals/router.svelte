<script lang="ts">
	import { api, profile, loading } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { type Router, type Service, type UpsertRouterParams } from '$lib/types/router';
	import { toast } from 'svelte-sonner';
	import RouterForm from '../forms/router.svelte';
	import ServiceForm from '../forms/service.svelte';
	import Separator from '../ui/separator/separator.svelte';

	interface Props {
		router: Router;
		service: Service;
		open?: boolean;
		mode: 'create' | 'edit';
	}

	let {
		router = $bindable(),
		service = $bindable(),
		open = $bindable(false),
		mode = 'create'
	}: Props = $props();

	const update = async () => {
		try {
			// Ensure proper name formatting and synchronization
			if (!router.name.includes('@')) {
				router.name = `${router.name}@http`;
			}
			if (service.loadBalancer?.servers?.length === 0) {
				toast.error('At least one server is required');
				return;
			}

			// Sync service name with router
			service.name = router.name;
			service.protocol = router.protocol;
			router.service = router.name;

			let params: UpsertRouterParams = {
				name: router.name,
				protocol: router.protocol
			};
			switch (router.protocol) {
				case 'http':
					params.router = router;
					params.service = service;
					break;
				case 'tcp':
					params.tcpRouter = router;
					params.tcpService = service;
					break;
				case 'udp':
					params.udpRouter = router;
					params.udpService = service;
					break;
			}

			await api.upsertRouter($profile.id, params);
			open = false;
			toast.success(`Router ${mode === 'create' ? 'created' : 'updated'} successfully`);
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(`Failed to ${mode} router`, {
				description: e.message
			});
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] max-w-xl overflow-y-auto">
		<Tabs.Root value="router" class="mt-4">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="router">Router</Tabs.Trigger>
				<Tabs.Trigger value="service">Service</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="router">
				<RouterForm bind:router {mode} />
			</Tabs.Content>
			<Tabs.Content value="service">
				<ServiceForm bind:service {router} {mode} />
			</Tabs.Content>
		</Tabs.Root>

		<Separator />

		<Button type="submit" onclick={update} disabled={$loading}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
