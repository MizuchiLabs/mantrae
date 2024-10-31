<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { upsertService, upsertRouter } from '$lib/api';
	import { type Router, type Service } from '$lib/types/config';
	import ServiceForm from '../forms/service.svelte';
	import RouterForm from '../forms/router.svelte';

	export let router: Router;
	export let service: Service;
	export let open = false;
	export let disabled = false;

	let routerForm: RouterForm;
	let serviceForm: ServiceForm;

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
	<Dialog.Trigger />
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
