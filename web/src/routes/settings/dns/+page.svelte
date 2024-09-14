<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button';
	import ProviderModal from '$lib/components/modals/providerModal.svelte';
	import powerdns from '$lib/images/powerdns.svg';
	import { deleteProvider, getProviders, provider } from '$lib/api';
	import { newProvider, type DNSProvider } from '$lib/types/base';
	import { onMount } from 'svelte';

	let dnsProvider: DNSProvider;
	let openModal = false;

	const createModal = () => {
		dnsProvider = newProvider();
		openModal = true;
	};
	const updateModal = (p: DNSProvider) => {
		dnsProvider = p;
		openModal = true;
	};

	onMount(async () => {
		if ($provider === undefined) {
			await getProviders();
		}
	});
</script>

<ProviderModal bind:dnsProvider bind:open={openModal} />

<div class="mt-4 flex flex-col gap-4 px-4 md:flex-row">
	<Button class="flex items-center gap-2 bg-red-400 text-black" on:click={createModal}>
		<span>Add Provider</span>
		<iconify-icon icon="fa6-solid:plus" />
	</Button>
</div>

<div class="flex flex-col gap-4 px-4 md:flex-row">
	{#if $provider}
		{#each $provider as p}
			<Card.Root class="w-full md:w-[400px]">
				<Card.Header>
					<Card.Title class="flex items-center justify-between gap-2">
						<span>{p.name}</span>
						<div class="flex items-center gap-2">
							{#if p.type === 'cloudflare'}
								<iconify-icon icon="devicon:cloudflare" width="26" />
							{:else if p.type === 'powerdns'}
								<img src={powerdns} alt="PowerDNS" width="26" />
							{/if}

							{#if p.is_active}
								<iconify-icon icon="fa6-solid:star" class="text-yellow-400" />
							{/if}
						</div>
					</Card.Title>
				</Card.Header>
				<Card.Content class="space-y-2"></Card.Content>
				<Card.Footer class="grid grid-cols-2 items-center gap-2">
					<Button variant="ghost" class="w-full bg-red-400" on:click={() => deleteProvider(p.id)}>
						Delete
					</Button>
					<Button variant="ghost" class="w-full bg-orange-400" on:click={() => updateModal(p)}>
						Update
					</Button>
				</Card.Footer>
			</Card.Root>
		{/each}
	{/if}
</div>
