<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import DNSModal from '$lib/components/modals/dns.svelte';
	import powerdns from '$lib/images/powerdns.svg';
	import { deleteProvider, getProviders, provider, profile } from '$lib/api';
	import { newProvider, type DNSProvider } from '$lib/types/base';
	import { onMount } from 'svelte';
	import { CircleCheck, Globe, Plus, Sparkles, Star } from 'lucide-svelte';

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
		await getProviders();
	});
</script>

<DNSModal bind:dnsProvider bind:open={openModal} />

<div class="mt-4 flex flex-col gap-4 px-4 md:flex-row">
	<Button class="flex items-center gap-2 bg-red-400 text-black" on:click={createModal}>
		<span>Add Provider</span>
		<Plus size="1rem" />
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
							{#if p.isActive}
								<CircleCheck class="text-green-400" />
							{/if}
						</div>
					</Card.Title>
				</Card.Header>
				<Card.Content class="flex flex-col gap-2">
					<div class="flex items-center gap-2">
						<span class="col-span-1 font-mono">Type:</span>
						<Badge variant="secondary">{p.type}</Badge>
					</div>
					<div class="flex items-center gap-2">
						<span class="col-span-1 font-mono">IP:</span>
						<Badge variant="secondary">{p.externalIp}</Badge>
					</div>
					{#if p.type === 'cloudflare'}
						<div class="flex items-center gap-2">
							<span class="col-span-1 font-mono">Proxied:</span>
							{#if p.proxied}
								<Badge variant="secondary" class="bg-green-300">Yes</Badge>
							{:else}
								<Badge variant="secondary" class="bg-red-300">No</Badge>
							{/if}
						</div>
					{/if}
					{#if p.type === 'powerdns' || p.type === 'technitium'}
						<div class="flex items-center gap-2">
							<span class="col-span-1 font-mono">Endpoint:</span>
							<Badge variant="secondary">{p.apiUrl}</Badge>
						</div>
					{/if}
				</Card.Content>
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
	{:else}
		<div class="flex h-[calc(75vh)] w-full flex-col items-center justify-center gap-2">
			<Globe size="8rem" />
			<span class="text-xl font-semibold">No providers configured</span>
		</div>
	{/if}
</div>
