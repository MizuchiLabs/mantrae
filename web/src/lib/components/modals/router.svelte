<script lang="ts">
	import { api, profile } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { type Router, type Service } from '$lib/types/router';
	import { toast } from 'svelte-sonner';
	import RouterForm from '../forms/router.svelte';
	import ServiceForm from '../forms/service.svelte';

	interface Props {
		router: Router | undefined;
		service: Service | undefined;
		open?: boolean;
		disabled?: boolean;
	}

	let {
		router = $bindable({} as Router),
		service = $bindable({} as Service),
		open = $bindable(false),
		disabled = false
	}: Props = $props();

	const update = async () => {
		try {
			await api.upsertRouter($profile.id, router);
			open = false;
		} catch (e) {
			toast.error('Failed to save router', {
				description: e.message
			});
		}
	};

	// Set defaults
	$effect(() => {
		router.type = router.type || 'http';
		router.tls = router.tls || {};
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] max-w-xl overflow-y-auto">
		<Tabs.Root value="router" class="mt-4">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="router">Router</Tabs.Trigger>
				<Tabs.Trigger value="service">Service</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="router">
				<RouterForm {router} {disabled} />
			</Tabs.Content>
			<Tabs.Content value="service">
				<ServiceForm {service} {router} {disabled} />
			</Tabs.Content>
		</Tabs.Root>
		<Button class="w-full" onclick={() => update()}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
