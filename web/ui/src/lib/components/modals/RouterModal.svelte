<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { type Router } from '$lib/gen/mantrae/v1/router_pb';
	import { type Service } from '$lib/gen/mantrae/v1/service_pb';
	import { protocolTypes, unmarshalConfig } from '$lib/types';
	import { Bot, Server, StarIcon } from '@lucide/svelte';
	import type {
		Service as HTTPService,
		Router as HTTPRouter,
		TCPRouter,
		TCPService,
		UDPRouter,
		UDPService
	} from '$lib/gen/zen/traefik-schemas';
	import HTTPRouterForm from '../forms/HTTPRouterForm.svelte';
	import TCPRouterForm from '../forms/TCPRouterForm.svelte';
	import UDPRouterForm from '../forms/UDPRouterForm.svelte';
	import HTTPServiceForm from '../forms/HTTPServiceForm.svelte';
	import TCPServiceForm from '../forms/TCPServiceForm.svelte';
	import UDPServiceForm from '../forms/UDPServiceForm.svelte';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { router } from '$lib/api/router.svelte';
	import { dns } from '$lib/api/dns.svelte';
	import { service } from '$lib/api/service.svelte';

	interface Props {
		data?: Router;
		open?: boolean;
	}
	let { data, open = $bindable(false) }: Props = $props();

	let routerData = $state({} as Router);
	let serviceData = $state({} as Service);
	let getService = $derived(service.get(routerData.name, routerData.type));
	$effect(() => {
		if (data) routerData = { ...data };
	});
	$effect(() => {
		if (data?.id) {
			if (getService.isSuccess && getService.data) {
				serviceData = { ...getService.data };
			}
		}
	});
	$effect(() => {
		if (!open) {
			routerData = {} as Router;
			serviceData = {} as Service;
		}
	});

	const dnsList = dns.list();
	const createRouter = router.create();
	const updateRouter = router.update();
	const createService = service.create();
	const updateService = service.update();
	function onsubmit() {
		serviceData.name = routerData.name;
		serviceData.type = routerData.type;
		serviceData.enabled = routerData.enabled;
		if (routerData.id) {
			updateRouter.mutate({ ...routerData });
		} else {
			createRouter.mutate({ ...routerData });
		}
		if (serviceData.id) {
			updateService.mutate({ ...serviceData });
		} else {
			createService.mutate({ ...serviceData });
		}
		open = false;
	}

	// Get service preview for agent-managed router
	let servicePreview = $derived.by(() => {
		if (!serviceData?.config) return 'No service configured';

		if (serviceData.type === ProtocolType.HTTP) {
			const config = unmarshalConfig(serviceData.config) as HTTPService;
			const servers = config.loadBalancer?.servers || [];
			return servers.length > 0 ? servers[0].url : 'No servers';
		} else if (serviceData.type === ProtocolType.TCP) {
			const config = unmarshalConfig(serviceData.config) as TCPService;
			const servers = config.loadBalancer?.servers || [];
			return servers.length > 0 ? servers[0].address : 'No servers';
		} else if (serviceData.type === ProtocolType.UDP) {
			const config = unmarshalConfig(serviceData.config) as UDPService;
			const servers = config.loadBalancer?.servers || [];
			return servers.length > 0 ? servers[0].address : 'No servers';
		}
		return 'No service configured';
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] max-w-xl overflow-y-auto">
		{#if routerData.agentId}
			<div class="space-y-4">
				<div>
					<h2 class="text-l font-semibold tracking-tight">Agent-Managed Router</h2>
					<p class="text-sm text-muted-foreground">Automatically managed via Docker labels</p>
				</div>

				<Alert.Root variant="default" class="border-dashed border-blue-200 bg-blue-50">
					<Bot class="h-4 w-4 text-blue-600" />
					<Alert.Title class="text-blue-900">Configuration Managed by Agent</Alert.Title>
					<Alert.Description class="text-blue-800">
						This router is automatically managed by an agent. Configuration changes must be made
						through Docker labels. Only DNS provider assignments can be modified here.
					</Alert.Description>
				</Alert.Root>

				<!-- DNS Providers (Editable) -->
				{#if dnsList.isSuccess && dnsList.data.length}
					<Card.Root>
						<Card.Content>
							<div class="flex items-center justify-between">
								<div>
									<Card.Title class="flex items-center gap-2">DNS Providers</Card.Title>
									<Card.Description>Manage DNS providers for this router</Card.Description>
								</div>

								{#if routerData.type !== ProtocolType.UDP && dnsList.isSuccess}
									<Select.Root
										type="multiple"
										value={routerData.dnsProviders?.map((item) => item.id.toString())}
										onValueChange={(value) => {
											let providers =
												dnsList.data?.filter((p) => value.includes(p.id.toString())) ?? [];
											routerData.dnsProviders = providers;
											updateRouter.mutate({ ...routerData });
										}}
									>
										<Select.Trigger>
											{routerData.dnsProviders?.length > 0
												? routerData.dnsProviders.length + ' Selected'
												: 'None'}
										</Select.Trigger>
										<Select.Content align="end">
											{#each dnsList.data || [] as dns (dns.id)}
												<Select.Item value={dns.id.toString()} class="flex items-center gap-2">
													<span class="truncate">{dns.name}</span>
													{#if dns.isDefault}
														<StarIcon class="text-yellow-400" />
													{/if}
												</Select.Item>
											{/each}
										</Select.Content>
									</Select.Root>
								{/if}
							</div>
						</Card.Content>
					</Card.Root>
				{/if}

				<!-- Router Info (Read-only) -->
				<Card.Root class="bg-muted/30">
					<Card.Header>
						<Card.Title class="flex items-center gap-2 text-muted-foreground">
							<Server class="h-4 w-4" />
							Router Configuration
						</Card.Title>
						<Card.Description>View-only configuration managed by agent</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid grid-cols-2 gap-6">
							<div class="space-y-1">
								<Label class="text-muted-foreground">Router Name</Label>
								<p class="font-medium">{routerData.name || 'Not set'}</p>
							</div>
							<div class="space-y-1">
								<Label class="text-muted-foreground">Protocol</Label>
								<Badge variant="outline">
									{protocolTypes.find((t) => t.value === routerData.type)?.label || 'Unknown'}
								</Badge>
							</div>
							<div class="space-y-1">
								<Label class="text-muted-foreground">Service Endpoint</Label>
								<p class="text-sm font-medium">{servicePreview}</p>
							</div>
							<div class="space-y-1">
								<Label class="text-muted-foreground">Status</Label>
								<Badge variant={routerData.enabled ? 'default' : 'secondary'}>
									{routerData.enabled ? 'Enabled' : 'Disabled'}
								</Badge>
							</div>
						</div>
					</Card.Content>
				</Card.Root>

				<Button onclick={() => (open = false)} variant="outline" class="w-full">Close</Button>
			</div>
		{:else}
			<Tabs.Root value="router" class="mt-2 sm:mt-4">
				<Tabs.List class="grid w-full grid-cols-2">
					<Tabs.Trigger value="router">Router</Tabs.Trigger>
					<Tabs.Trigger value="service">Service</Tabs.Trigger>
				</Tabs.List>
				<Tabs.Content value="router" class="space-y-4">
					<Card.Root>
						<Card.Header
							class="space-y-2 sm:flex sm:flex-row sm:items-center sm:justify-between sm:space-y-0"
						>
							<div>
								<Card.Title class="flex items-center gap-2">Router Configuration</Card.Title>
								<Card.Description class="text-sm text-muted-foreground">
									Define how traffic is routed to your service
								</Card.Description>
							</div>

							<!-- DNS Providers -->
							{#if routerData.type !== ProtocolType.UDP && dnsList.isSuccess}
								<Select.Root
									type="multiple"
									value={routerData.dnsProviders?.map((item) => item.id.toString())}
									onValueChange={(value) => {
										let providers =
											dnsList.data?.filter((p) => value.includes(p.id.toString())) ?? [];
										routerData.dnsProviders = providers;
										updateRouter.mutate({ ...routerData });
									}}
								>
									<Select.Trigger>
										{routerData.dnsProviders?.length > 0
											? routerData.dnsProviders.length + ' Selected'
											: 'None'}
									</Select.Trigger>
									<Select.Content align="end">
										{#each dnsList.data || [] as dns (dns.id)}
											<Select.Item value={dns.id.toString()} class="flex items-center gap-2">
												<span class="truncate">{dns.name}</span>
												{#if dns.isDefault}
													<StarIcon class="text-yellow-400" />
												{/if}
											</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>
							{/if}
						</Card.Header>
						<Card.Content class="space-y-4">
							<div class="grid w-full grid-cols-1 gap-4 sm:grid-cols-3 sm:gap-2">
								<div
									class="flex flex-col gap-2 {routerData.id ? 'sm:col-span-3' : 'sm:col-span-2'}"
								>
									<Label for="name">Name</Label>
									<Input
										id="name"
										bind:value={routerData.name}
										required
										placeholder="Router Name"
										class="truncate"
									/>
								</div>

								{#if !routerData.id}
									<div class="flex flex-col gap-2 sm:col-span-1">
										<Label for="type">Protocol</Label>
										<Select.Root
											type="single"
											name="type"
											value={routerData.type?.toString()}
											onValueChange={(value) => {
												// Reset config
												routerData.type = parseInt(value, 10);
												switch (routerData.type) {
													case ProtocolType.HTTP:
														routerData.config = {} as HTTPRouter;
														break;
													case ProtocolType.TCP:
														routerData.config = {} as TCPRouter;
														break;
													case ProtocolType.UDP:
														routerData.config = {} as UDPRouter;
														break;
												}
											}}
										>
											<Select.Trigger class="w-full">
												<span class="truncate">
													{protocolTypes.find((t) => t.value === routerData.type)?.label ??
														'Select protocol'}
												</span>
											</Select.Trigger>
											<Select.Content>
												{#each protocolTypes as t (t.value)}
													<Select.Item value={t.value.toString()}>
														{t.label}
													</Select.Item>
												{/each}
											</Select.Content>
										</Select.Root>
									</div>
								{/if}
							</div>

							{#if routerData.type === ProtocolType.HTTP}
								<HTTPRouterForm bind:data={routerData} />
							{/if}
							{#if routerData.type === ProtocolType.TCP}
								<TCPRouterForm bind:data={routerData} />
							{/if}
							{#if routerData.type === ProtocolType.UDP}
								<UDPRouterForm bind:data={routerData} />
							{/if}
						</Card.Content>
					</Card.Root>
				</Tabs.Content>
				<Tabs.Content value="service" class="space-y-4">
					<Card.Root>
						<Card.Header>
							<Card.Title class="flex items-center gap-2">Service Configuration</Card.Title>
							<Card.Description>Configure backend servers and load balancing</Card.Description>
						</Card.Header>
						<Card.Content class="flex flex-col gap-3">
							{#if routerData.type === ProtocolType.HTTP}
								<HTTPServiceForm bind:service={serviceData} />
							{/if}
							{#if routerData.type === ProtocolType.TCP}
								<TCPServiceForm bind:service={serviceData} />
							{/if}
							{#if routerData.type === ProtocolType.UDP}
								<UDPServiceForm bind:service={serviceData} />
							{/if}
						</Card.Content>
					</Card.Root>
				</Tabs.Content>
			</Tabs.Root>

			<Button type="submit" class="w-full" onclick={onsubmit}
				>{routerData.id ? 'Update' : 'Create'}</Button
			>
		{/if}
	</Dialog.Content>
</Dialog.Root>
