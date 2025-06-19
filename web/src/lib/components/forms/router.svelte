<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import logo from '$lib/images/logo.svg';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { CircleCheck, Globe, Lock, Star } from '@lucide/svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Toggle } from '$lib/components/ui/toggle';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { toast } from 'svelte-sonner';
	import type { RouterDNSProvider } from '$lib/types';
	import Button from '../ui/button/button.svelte';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { Router as HttpRouter } from '$lib/gen/tygo/dynamic';
	import type { TCPRouter as TcpRouter } from '$lib/gen/tygo/dynamic';
	import type { UDPRouter as UdpRouter } from '$lib/gen/tygo/dynamic';
	import { dnsClient, entryPointClient, middlewareClient } from '$lib/api';

	let { router = $bindable() }: { router: Router } = $props();
	// let config = $derived(parseRouterConfig(router));
	let config = $derived.by(() => parseRouterConfig(router));
    let httpConfig = $state({} as HttpRouter);
    let tcpConfig = $state({} as TcpRouter);
    let udpConfig = $state({} as UdpRouter);

	type TypedRouterConfig = HttpRouter | TcpRouter | UdpRouter;
	function parseRouterConfig(router: Router): TypedRouterConfig {
		const config = router.config ?? {};

		switch (router.type) {
			case RouterType.HTTP:
				return config as HttpRouter;
			case RouterType.TCP:
				return config as TcpRouter;
			case RouterType.UDP:
				return config as UdpRouter;
			default:
				throw new Error('Unsupported router type');
		}
	}

	// Computed properties
	// let routerDNS: RouterDNSProvider[] = $derived(
	// 	$rdps?.filter((item) => item.routerName === router.name)
	// );
	// let rdpNames = $derived(routerDNS ? routerDNS.map((item) => item.providerName) : []);
	// let rdpIDs = $derived(routerDNS ? routerDNS.map((item) => item.providerId) : []);
	// let routerProvider = $derived(router.name ? router.name?.split('@')[1] : 'http');
	// let certResolvers = $derived([
	// 	...new Set(
	// 		$routers.filter((item) => item.tls?.certResolver).map((item) => item.tls?.certResolver)
	// 	)
	// ]);

	// Convert enum to select options
	const dnsTypes = Object.entries(DnsProviderType).map(([key, value]) => ({
		label: key.charAt(0).toUpperCase() + key.slice(1).toLowerCase(),
		value
	}));
	// Parsing select string value back to number (enum)
	function handleTypeChange(value: string) {
		dns.type = parseInt(value, 10);
	}

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
			<Card.Title>{router.id ? 'Update' : 'Create'} Router</Card.Title>
			<Card.Description>
				{router.id ? 'Update existing router' : 'Create a new router'}
			</Card.Description>
		</div>

		<!-- DNS Providers -->
		{#await dnsClient.listDnsProviders({}) then value}
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
								<Badge>{value.dnsProviders.length ? value.dnsProviders.join(', ') : 'None'}</Badge>
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
				value={value.dnsProviders.map((item) => item.id.toString())}
				onValueChange={handleDNSProviderChange}
				bind:open={selectDNSOpen}
			>
				<Select.Content customAnchor={dnsAnchor} align="end">
					{#each value.dnsProviders as dns (dns.id)}
						<Select.Item value={dns.id.toString()} class="flex items-center gap-2">
							{dns.name} ({dns.type})
							{#if dns.isActive}
								<CircleCheck size="1rem" class="text-green-400" />
							{/if}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		{/await}
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		<!-- Provider Type Toggles -->
		<div class="flex items-center justify-end gap-1 font-mono text-sm">
			<Toggle
				size="sm"
				pressed={router.type === RouterType.HTTP}
				onPressedChange={() => (router.type = RouterType.HTTP)}
				class="font-bold data-[state=on]:bg-green-300 dark:data-[state=on]:text-black"
			>
				HTTP
			</Toggle>
			<Toggle
				size="sm"
				pressed={router.type === RouterType.TCP}
				onPressedChange={() => (router.type = RouterType.TCP)}
				class="font-bold data-[state=on]:bg-blue-300 dark:data-[state=on]:text-black"
			>
				TCP
			</Toggle>
			<Toggle
				size="sm"
				pressed={router.type === RouterType.UDP}
				onPressedChange={() => (router.type = RouterType.UDP)}
				class="font-bold data-[state=on]:bg-red-300 dark:data-[state=on]:text-black"
			>
				UDP
			</Toggle>
		</div>

		<!-- Name -->
		<div class="grid grid-cols-4 items-center gap-1">
			<Label for="name" class="mr-2 text-right">Name</Label>
			<div class="relative col-span-3">
				<Input
					id="name"
					bind:value={router.name}
					placeholder="Router name"
					class="col-span-3"
					required
				/>
			</div>
		</div>

		<!-- Only HTTP and TCP -->

		{#if router.type !== RouterType.UDP}
			<!-- Entrypoints -->
			<div class="grid grid-cols-4 items-center gap-1">
				<Label class="mr-2 text-right">Entrypoints</Label>
				<Select.Root type="multiple" bind:value={config.entryPoints}>
					<Select.Trigger class="col-span-3">
						{config.entryPoints?.join(', ') || 'Select entrypoints'}
					</Select.Trigger>
					<Select.Content>
						{#await entryPointClient.listEntryPoints({}) then value}
							{#each value.entryPoints as e (e.id)}
								<Select.Item value={e.name}>
									<div class="flex items-center gap-2">
										{e.name}
										{#if e.isDefault}
											<Star size="1rem" class="text-yellow-300" />
										{/if}
									</div>
								</Select.Item>
							{/each}
						{/await}
					</Select.Content>
				</Select.Root>
			</div>

			<!-- Middlewares -->
			<div class="grid grid-cols-4 items-center gap-1">
				<Label class="mr-2 text-right">Middlewares</Label>
				<Select.Root type="multiple" bind:value={config.middlewares}>
					<Select.Trigger class="col-span-3">
						{config.middlewares?.join(', ') || 'Select middlewares'}
					</Select.Trigger>
					<Select.Content>
						<!-- TODO: Filter by protocol -->
						{#await middlewareClient.listMiddlewares({}) then value}
							{#each value.middlewares as middleware (middleware.name)}
								<Select.Item value={middleware.name}>
									{middleware.name}
								</Select.Item>
							{/each}
						{/await}
					</Select.Content>
				</Select.Root>
			</div>

			<!-- <!-- TLS --> -->
			<!-- <div class="grid grid-cols-4 items-center gap-1"> -->
			<!-- 	<Label class="mr-2 text-right">TLS</Label> -->
			<!-- 	<Switch -->
			<!-- 		bind:checked={config.tls?.passthrough} -->
			<!-- 		onCheckedChange={(checked) => { -->
			<!-- 			if (checked) { -->
			<!-- 				router.tls = {}; -->
			<!-- 			} else { -->
			<!-- 				delete router.tls; -->
			<!-- 			} -->
			<!-- 		}} -->
			<!-- 		class="col-span-3" -->
			<!-- 		disabled={!source.isLocal()} -->
			<!-- 	/> -->
			<!-- </div> -->
			<!---->
			<!-- {#if router.protocol === 'http' && !!router.tls} -->
			<!-- 	<div class="grid grid-cols-4 items-center gap-1"> -->
			<!-- 		<Label class="mr-2 text-right">Resolver</Label> -->
			<!-- 		<div class="col-span-3"> -->
			<!-- 			<Input -->
			<!-- 				bind:value={router.tls.certResolver} -->
			<!-- 				class="mb-1" -->
			<!-- 				placeholder="Certificate resolver" -->
			<!-- 			/> -->
			<!-- 			{#if source.isLocal()} -->
			<!-- 				<div class="flex flex-wrap gap-1"> -->
			<!-- 					{#each certResolvers as resolver (resolver)} -->
			<!-- 						{#if resolver !== router.tls.certResolver} -->
			<!-- 							<Badge -->
			<!-- 								onclick={() => router.tls && (router.tls.certResolver = resolver)} -->
			<!-- 								class="cursor-pointer" -->
			<!-- 							> -->
			<!-- 								{resolver} -->
			<!-- 							</Badge> -->
			<!-- 						{/if} -->
			<!-- 					{/each} -->
			<!-- 				</div> -->
			<!-- 			{/if} -->
			<!-- 		</div> -->
			<!-- 	</div> -->
			<!-- {/if} -->
			<!---->
			<!-- {#if router.protocol === 'tcp' && !!router.tls} -->
			<!-- 	<div class="grid grid-cols-4 items-center gap-1"> -->
			<!-- 		<Label class="mr-2 text-right">Passthrough</Label> -->
			<!-- 		<Switch -->
			<!-- 			checked={!!router.tls.passthrough} -->
			<!-- 			onCheckedChange={(checked) => (router.tls ? (router.tls.passthrough = checked) : null)} -->
			<!-- 			class="col-span-3" -->
			<!-- 			disabled={!source.isLocal()} -->
			<!-- 		/> -->
			<!-- 	</div> -->
			<!-- {/if} -->
			<!---->
			<!-- <!-- Rule --> -->
			<!-- <RuleEditor -->
			<!-- 	bind:rule={router.rule} -->
			<!-- 	bind:type={router.protocol} -->
			<!-- 	disabled={!source.isLocal()} -->
			<!-- /> -->
		{/if}
	</Card.Content>
</Card.Root>
