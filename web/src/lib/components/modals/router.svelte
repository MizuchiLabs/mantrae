<script lang="ts">
	import { api, profile } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { type Router, type Service, type UpsertRouterParams } from '$lib/types/router';
	import { toast } from 'svelte-sonner';
	import RouterForm from '../forms/router.svelte';
	import ServiceForm from '../forms/service.svelte';

	interface Props {
		router?: Router;
		service?: Service;
		open?: boolean;
		mode: 'create' | 'edit';
	}

	const defaultRouter: Router = {
		name: '',
		type: 'http',
		tls: {},
		entryPoints: [],
		middlewares: [],
		rule: '',
		service: ''
	};

	const defaultService: Service = {
		name: '',
		type: 'http',
		loadBalancer: {
			servers: [],
			passHostHeader: true
		}
	};

	let {
		router = $bindable(defaultRouter),
		service = $bindable(defaultService),
		open = $bindable(false),
		mode = 'create'
	}: Props = $props();

	const update = async () => {
		try {
			// Ensure proper name formatting and synchronization
			if (!router.name.includes('@')) {
				router.name = `${router.name}@http`;
			}

			// Sync service name with router
			service.name = router.name;
			router.service = router.name;

			let params: UpsertRouterParams = {
				name: router.name,
				type: router.type
			};
			switch (router.type) {
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
				<RouterForm
					bind:router
					{mode}
					disabled={mode === 'edit' && router.name?.split('@')[1] !== 'http'}
				/>
			</Tabs.Content>
			<Tabs.Content value="service">
				<ServiceForm
					bind:service
					{router}
					disabled={mode === 'edit' && router.name?.split('@')[1] !== 'http'}
				/>
			</Tabs.Content>
		</Tabs.Root>
		<Button class="w-full" onclick={() => update()}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
