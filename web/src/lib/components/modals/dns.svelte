<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import type { DNSProvider } from '$lib/types/base';
	import { updateProvider, createProvider, provider } from '$lib/api';
	import DNSForm from '../forms/dns.svelte';

	export let dnsProvider: DNSProvider;
	export let open = false;

	const update = async () => {
		if (
			dnsProvider.name === '' ||
			dnsProvider.type === '' ||
			dnsProvider.api_key === '' ||
			dnsProvider.external_ip === ''
		)
			return;

		// PowerDNS requires the API URL to start with http:// or https://
		if (
			dnsProvider.type === 'powerdns' &&
			!dnsProvider.api_url?.startsWith('http://') &&
			!dnsProvider.api_url?.startsWith('https://')
		) {
			dnsProvider.api_url = 'http://' + dnsProvider.api_url;
		}

		if ($provider?.find((p) => p.name === dnsProvider.name)) {
			await updateProvider(dnsProvider);
		} else {
			await createProvider(dnsProvider);
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[500px]">
		<DNSForm bind:provider={dnsProvider} />
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
