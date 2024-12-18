<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import type { DNSProvider } from '$lib/types/base';
	import { upsertProvider } from '$lib/api';
	import DNSForm from '../forms/dns.svelte';

	export let dnsProvider: DNSProvider;
	export let open = false;
	let dnsForm: DNSForm;

	const update = async () => {
		const valid = dnsForm.validate();
		if (valid) {
			await upsertProvider(dnsProvider);
			open = false;
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[550px]">
		<DNSForm bind:provider={dnsProvider} bind:this={dnsForm} />
		<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
