<script lang="ts">
	import { upsertProvider } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import type { DNSProvider } from '$lib/types/base';
	import DNSForm from '../forms/dns.svelte';

	interface Props {
		dnsProvider: DNSProvider;
		open?: boolean;
	}

	let { dnsProvider = $bindable(), open = $bindable(false) }: Props = $props();
	let dnsForm: DNSForm = $state();

	const update = async () => {
		const valid = dnsForm.validate();
		if (valid) {
			await upsertProvider(dnsProvider);
			open = false;
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[550px]">
		<DNSForm bind:provider={dnsProvider} bind:this={dnsForm} />
		<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
