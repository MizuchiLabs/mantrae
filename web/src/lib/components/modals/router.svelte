<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { toast } from 'svelte-sonner';
	import RouterForm from '../forms/router.svelte';
	import ServiceForm from '../forms/service.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import type { Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { Service } from '$lib/gen/mantrae/v1/service_pb';

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
			if (service.loadBalancer?.servers?.length === 0) {
				toast.error('At least one server is required');
				return;
			}

			const protocol = (router.protocol = service.protocol = router.protocol);
			router.service = router.name;
			service.name = router.name;

			const params: UpsertRouterParams = {
				name: router.name,
				protocol,
				...(protocol === 'http'
					? {
							router,
							service
						}
					: {
							[`${protocol}Router`]: router,
							[`${protocol}Service`]: service
						})
			};
			await api.upsertRouter(params);
			toast.success(`Router ${mode === 'create' ? 'created' : 'updated'} successfully`);
			open = false;
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
				<ServiceForm bind:service bind:router />
			</Tabs.Content>
		</Tabs.Root>

		<Separator />

		<Button type="submit" onclick={update}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
