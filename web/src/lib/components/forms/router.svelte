<script lang="ts">
	import RuleEditor from '../utils/ruleEditor.svelte';
	import logo from '$lib/images/logo.svg';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { CircleCheck, Lock } from 'lucide-svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Toggle } from '$lib/components/ui/toggle';
	import { api, dnsProviders, entrypoints, routers, middlewares, traefik, rdps } from '$lib/api';
	import { onMount } from 'svelte';
	import { type Router } from '$lib/types/router';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { toast } from 'svelte-sonner';
	import type { RouterDNSProvider } from '$lib/types';

	interface Props {
		router: Router;
		mode: 'create' | 'edit';
	}

	let { router = $bindable(), mode }: Props = $props();

	// Computed properties
	let routerDNS: RouterDNSProvider = $derived(
		$rdps?.filter((item) => item.routerName === router.name)[0]
	);
	let rdpName = $derived(routerDNS ? routerDNS.providerName : '');
	let routerProvider = $derived(router.name ? router.name?.split('@')[1] : 'http');
	let isHttpProvider = $derived(routerProvider === 'http' || !routerProvider);
	let isHttpType = $derived(router.protocol === 'http');
	let disabled = $derived(routerProvider !== 'http' && mode === 'edit');
	let certResolvers = $derived([
		...new Set(
			$routers.filter((item) => item.tls?.certResolver).map((item) => item.tls?.certResolver)
		)
	]);

	// let routerDNSProvider = $derived(router.dnsProvider);
	async function handleDNSProviderChange(providerId: string) {
		try {
			if (providerId === '') {
				await api.deleteRouterDNSProvider($traefik[0].id, router.name);
				toast.success('DNS Provider deleted successfully');
			} else {
				if (providerId === undefined) return;
				await api.setRouterDNSProvider($traefik[0].id, parseInt(providerId), router.name);
				toast.success('DNS Provider updated successfully');
			}
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to save dnsProvider', {
				description: e.message
			});
		}
	}

	onMount(async () => {
		await api.listDNSProviders();
	});
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{mode === 'create' ? 'Add' : 'Update'} Router</Card.Title>
		<Card.Description>
			{mode === 'create' ? 'Create a new router' : 'Edit existing router'}
		</Card.Description>
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		<!-- Provider Type Toggles -->
		{#if isHttpProvider}
			<div class="flex items-center justify-end gap-1 font-mono text-sm">
				<Toggle
					size="sm"
					pressed={router.protocol === 'http'}
					onPressedChange={() => (router.protocol = 'http')}
					{disabled}
					class="font-bold data-[state=on]:bg-green-300 dark:data-[state=on]:text-black"
				>
					HTTP
				</Toggle>
				<Toggle
					size="sm"
					pressed={router.protocol === 'tcp'}
					onPressedChange={() => (router.protocol = 'tcp')}
					{disabled}
					class="font-bold data-[state=on]:bg-blue-300 dark:data-[state=on]:text-black"
				>
					TCP
				</Toggle>
				<Toggle
					size="sm"
					pressed={router.protocol === 'udp'}
					onPressedChange={() => (router.protocol = 'udp')}
					{disabled}
					class="font-bold data-[state=on]:bg-red-300 dark:data-[state=on]:text-black"
				>
					UDP
				</Toggle>
			</div>
		{/if}

		<!-- Name -->
		<div class="grid grid-cols-4 items-center gap-1">
			<Label for="name" class="mr-2 text-right">Name</Label>
			<div class="relative col-span-3">
				<Input
					id="name"
					bind:value={router.name}
					placeholder="Router name"
					oninput={() => (router.name = router.name?.split('@')[0])}
					class="col-span-3"
					required
					disabled={disabled || mode === 'edit'}
				/>
				<!-- Icon based on provider -->
				{#if routerProvider !== ''}
					<span
						class="pointer-events-none absolute inset-y-0 right-3 flex items-center text-gray-400"
					>
						{#if isHttpProvider}
							<img src={logo} alt="HTTP" width="20" />
						{/if}
						{#if routerProvider === 'internal' || routerProvider === 'file'}
							<iconify-icon icon="devicon:traefikproxy" height="20"></iconify-icon>
						{/if}
						{#if routerProvider?.includes('docker')}
							<iconify-icon icon="logos:docker-icon" height="20"></iconify-icon>
						{/if}
						{#if routerProvider?.includes('kubernetes')}
							<iconify-icon icon="logos:kubernetes" height="20"></iconify-icon>
						{/if}
						{#if routerProvider === 'consul'}
							<iconify-icon icon="logos:consul" height="20"></iconify-icon>
						{/if}
						{#if routerProvider === 'nomad'}
							<iconify-icon icon="logos:nomad-icon" height="20"></iconify-icon>
						{/if}
						{#if routerProvider === 'kv'}
							<iconify-icon icon="logos:redis" height="20"></iconify-icon>
						{/if}
					</span>
				{/if}
			</div>
		</div>

		<!-- Entrypoints -->
		{#if isHttpProvider}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label class="mr-2 text-right">Entrypoints</Label>
				<Select.Root type="multiple" bind:value={router.entryPoints} {disabled}>
					<Select.Trigger class="col-span-3">
						{router.entryPoints?.length ? router.entryPoints.join(', ') : 'Select entrypoints'}
					</Select.Trigger>
					<Select.Content>
						{#each $entrypoints as ep}
							<Select.Item value={ep.name}>
								<div class="flex items-center gap-2">
									{ep.name}
									{#if ep.http?.tls}
										<Lock size="1rem" class="text-green-400" />
									{/if}
								</div>
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{/if}

		<!-- Middlewares -->
		{#if router.protocol !== 'udp'}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label class="mr-2 text-right">Middlewares</Label>
				<Select.Root type="multiple" bind:value={router.middlewares} {disabled}>
					<Select.Trigger class="col-span-3">
						{router.middlewares?.length ? router.middlewares.join(', ') : 'Select middlewares'}
					</Select.Trigger>
					<Select.Content>
						{#each $middlewares.filter((m) => m.protocol === router.protocol) as middleware}
							<Select.Item value={middleware.name}>
								{middleware.name}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{/if}

		<!-- TLS Configuration -->
		{#if isHttpType && router.protocol !== 'udp' && router.tls}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label class="mr-2 text-right">Cert Resolver</Label>
				<div class="col-span-3">
					<Input
						bind:value={router.tls.certResolver}
						placeholder="Certificate resolver"
						{disabled}
					/>
					<div class="flex flex-wrap gap-1">
						{#each certResolvers as resolver}
							{#if resolver !== router.tls.certResolver}
								<Badge
									onclick={() => !disabled && router.tls && (router.tls.certResolver = resolver)}
									class={disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'}
								>
									{resolver}
								</Badge>
							{/if}
						{/each}
					</div>
				</div>
			</div>
		{/if}

		<!-- DNS Provider -->
		{#if $dnsProviders && mode === 'edit'}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label for="provider" class="mr-2 text-right">DNS Provider</Label>
				<Select.Root type="single" value={rdpName} onValueChange={handleDNSProviderChange}>
					<Select.Trigger class="col-span-3">
						{rdpName ? rdpName : 'Select DNS provider'}
					</Select.Trigger>
					<Select.Content>
						<Select.Item value="" label="">None</Select.Item>
						{#each $dnsProviders as dns}
							<Select.Item value={dns.id.toString()} class="flex items-center gap-2">
								{dns.name} ({dns.type})
								{#if dns.isActive}
									<CircleCheck size="1rem" class="text-green-400" />
								{/if}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{/if}

		<!-- Rule -->
		{#if router.protocol === 'http' || router.protocol === 'tcp'}
			<RuleEditor bind:rule={router.rule} bind:type={router.protocol} {disabled} />
		{/if}
	</Card.Content>
</Card.Root>
