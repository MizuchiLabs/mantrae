<script lang="ts">
	import { upsertRouter, upsertService } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { type Router, type Service } from '$lib/types/config';
	import RouterForm from '../forms/router.svelte';
	import ServiceForm from '../forms/service.svelte';

	interface Props {
		router: Router;
		service: Service;
		open?: boolean;
		disabled?: boolean;
	}

	let {
		router = $bindable(),
		service = $bindable(),
		open = $bindable(false),
		disabled = false
	}: Props = $props();

	let routerForm: RouterForm = $state();
	let serviceForm: ServiceForm = $state();

	const update = async () => {
		const rValid = routerForm.validate();
		const sValid = serviceForm.validate();
		if (rValid && sValid) {
			await upsertRouter(router);
			if (service.provider === 'http') {
				await upsertService(service);
			}
			open = false;
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
				<RouterForm bind:router {disabled} bind:this={routerForm} />
			</Tabs.Content>
			<Tabs.Content value="service">
				<ServiceForm bind:service bind:router {disabled} bind:this={serviceForm} />
			</Tabs.Content>
		</Tabs.Root>
		<Button class="w-full" on:click={() => update()}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
