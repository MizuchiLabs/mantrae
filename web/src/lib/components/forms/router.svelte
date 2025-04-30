<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import logo from '$lib/images/logo.svg';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { CircleCheck, Globe, Lock } from '@lucide/svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Toggle } from '$lib/components/ui/toggle';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api, dnsProviders, entrypoints, routers, middlewares, rdps } from '$lib/api';
	import { type Router } from '$lib/types/router';
	import { toast } from 'svelte-sonner';
	import type { RouterDNSProvider } from '$lib/types';
	import { source } from '$lib/stores/source';
	import Button from '../ui/button/button.svelte';

	interface Props {
		router: Router;
		mode: 'create' | 'edit';
	}

	let { router = $bindable(), mode }: Props = $props();

	// Computed properties
	let routerDNS: RouterDNSProvider[] = $derived(
		$rdps?.filter((item) => item.routerName === router.name)
	);
	let rdpNames = $derived(routerDNS ? routerDNS.map((item) => item.providerName) : []);
	let rdpIDs = $derived(routerDNS ? routerDNS.map((item) => item.providerId) : []);
	let routerProvider = $derived(router.name ? router.name?.split('@')[1] : 'http');
	let certResolvers = $derived([
		...new Set(
			$routers.filter((item) => item.tls?.certResolver).map((item) => item.tls?.certResolver)
		)
	]);

	async function handleDNSProviderChange(providerIDs: string[]) {
		try {
			await api.setRouterDNSProvider(providerIDs, router.name);
			toast.success('DNS Provider updated successfully');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to save dnsProvider', {
				description: e.message
			});
		}
	}

	function handleNameInput(e: Event) {
		router.name = (e.target as HTMLInputElement).value;
		if (router.name.includes('@')) {
			toast.warning("Avoid '@' in name – it's reserved for provider context");
			router.name = router.name.split('@')[0];
		}
	}

	let selectDNSOpen = $state(false);
	let dnsAnchor = $state<HTMLElement>(null!);
</script>

<Card.Root>
	<Card.Header class="flex flex-row items-center justify-between">
		<div>
			<Card.Title>{mode === 'create' ? 'Add' : 'Update'} Router</Card.Title>
			<Card.Description>
				{mode === 'create' ? 'Create a new router' : 'Edit existing router'}
			</Card.Description>
		</div>
		{#if $dnsProviders && mode === 'edit'}
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						<div bind:this={dnsAnchor}>
							<Button
								variant="ghost"
								size="sm"
								class="flex items-center gap-2"
								onclick={() => (selectDNSOpen = true)}
							>
								<Globe size={16} />
								<Badge>{rdpNames.length ? rdpNames.join(', ') : 'None'}</Badge>
							</Button>
						</div>
					</Tooltip.Trigger>
					<Tooltip.Content side="left" align="center">
						<p>Select DNS Provider</p>
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>

			<Select.Root
				type="multiple"
				value={rdpIDs.map((item) => item.toString())}
				onValueChange={handleDNSProviderChange}
				bind:open={selectDNSOpen}
			>
				<Select.Content customAnchor={dnsAnchor} align="end">
					{#each $dnsProviders as dns (dns.id)}
						<Select.Item value={dns.id.toString()} class="flex items-center gap-2">
							{dns.name} ({dns.type})
							{#if dns.isActive}
								<CircleCheck size="1rem" class="text-green-400" />
							{/if}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		{/if}
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		<!-- Provider Type Toggles -->
		{#if source.isLocal() && mode === 'create'}
			<div class="flex items-center justify-end gap-1 font-mono text-sm">
				<Toggle
					size="sm"
					pressed={router.protocol === 'http'}
					onPressedChange={() => (router.protocol = 'http')}
					class="font-bold data-[state=on]:bg-green-300 dark:data-[state=on]:text-black"
				>
					HTTP
				</Toggle>
				<Toggle
					size="sm"
					pressed={router.protocol === 'tcp'}
					onPressedChange={() => (router.protocol = 'tcp')}
					class="font-bold data-[state=on]:bg-blue-300 dark:data-[state=on]:text-black"
				>
					TCP
				</Toggle>
				<Toggle
					size="sm"
					pressed={router.protocol === 'udp'}
					onPressedChange={() => (router.protocol = 'udp')}
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
					oninput={handleNameInput}
					class="col-span-3"
					required
					disabled={!source.isLocal()}
				/>
				<!-- Icon based on provider -->
				<span
					class="pointer-events-none absolute inset-y-0 right-3 flex items-center text-gray-400"
				>
					{#if routerProvider === '' || routerProvider === 'http'}
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
			</div>
		</div>

		<!-- Entrypoints -->
		<div class="grid grid-cols-4 items-center gap-1">
			<Label class="mr-2 text-right">Entrypoints</Label>
			<Select.Root type="multiple" bind:value={router.entryPoints} disabled={!source.isLocal()}>
				<Select.Trigger class="col-span-3">
					{router.entryPoints?.length ? router.entryPoints.join(', ') : 'Select entrypoints'}
				</Select.Trigger>
				<Select.Content>
					{#each $entrypoints as ep (ep.name)}
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

		<!-- Only HTTP and TCP -->
		{#if router.protocol !== 'udp'}
			<!-- Middlewares -->
			<div class="grid grid-cols-4 items-center gap-1">
				<Label class="mr-2 text-right">Middlewares</Label>
				<Select.Root type="multiple" bind:value={router.middlewares} disabled={!source.isLocal()}>
					<Select.Trigger class="col-span-3">
						{router.middlewares?.length ? router.middlewares.join(', ') : 'Select middlewares'}
					</Select.Trigger>
					<Select.Content>
						{#each $middlewares.filter((m) => m.protocol === router.protocol) as middleware (middleware.name)}
							<Select.Item value={middleware.name}>
								{middleware.name}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>

			<!-- TLS -->
			<div class="grid grid-cols-4 items-center gap-1">
				<Label class="mr-2 text-right">TLS</Label>
				<Switch
					checked={!!router.tls}
					onCheckedChange={(checked) => {
						if (checked) {
							router.tls = {};
						} else {
							delete router.tls;
						}
					}}
					class="col-span-3"
					disabled={!source.isLocal()}
				/>
			</div>

			{#if router.protocol === 'http' && !!router.tls}
				<div class="grid grid-cols-4 items-center gap-1">
					<Label class="mr-2 text-right">Resolver</Label>
					<div class="col-span-3">
						<Input
							bind:value={router.tls.certResolver}
							class="mb-1"
							placeholder="Certificate resolver"
							disabled={!source.isLocal()}
						/>
						{#if source.isLocal()}
							<div class="flex flex-wrap gap-1">
								{#each certResolvers as resolver (resolver)}
									{#if resolver !== router.tls.certResolver}
										<Badge
											onclick={() => router.tls && (router.tls.certResolver = resolver)}
											class="cursor-pointer"
										>
											{resolver}
										</Badge>
									{/if}
								{/each}
							</div>
						{/if}
					</div>
				</div>
			{/if}

			{#if router.protocol === 'tcp' && !!router.tls}
				<div class="grid grid-cols-4 items-center gap-1">
					<Label class="mr-2 text-right">Passthrough</Label>
					<Switch
						checked={!!router.tls.passthrough}
						onCheckedChange={(checked) => (router.tls ? (router.tls.passthrough = checked) : null)}
						class="col-span-3"
						disabled={!source.isLocal()}
					/>
				</div>
			{/if}

			<!-- Rule -->
			<RuleEditor
				bind:rule={router.rule}
				bind:type={router.protocol}
				disabled={!source.isLocal()}
			/>
		{/if}
	</Card.Content>
</Card.Root>
