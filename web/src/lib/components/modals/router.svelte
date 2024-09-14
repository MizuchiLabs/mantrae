<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { upsertRouter } from '$lib/api';
	import { type Router, type Service } from '$lib/types/config';
	import ServiceForm from '../forms/service.svelte';
	import RouterForm from '../forms/router.svelte';

	export let router: Router;
	export let service: Service;
	export let open = false;
	export let disabled = false;
	let originalName = router?.name;

	const update = async () => {
		if (router.name === '') return;
		if (service === undefined) return;
		await upsertRouter(originalName, router, service);
		originalName = router.name;
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content>
		<Tabs.Root value="router" class="mt-4">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="router">Router</Tabs.Trigger>
				<Tabs.Trigger value="service">Service</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="router">
				<RouterForm bind:router {disabled} />
			</Tabs.Content>
			<Tabs.Content value="service">
				<ServiceForm bind:service {disabled} />
			</Tabs.Content>
		</Tabs.Root>
		<Dialog.Close>
			<Button class="w-full" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
