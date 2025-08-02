<script lang="ts">
	import { routerClient, serviceClient } from '$lib/api';
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
	import { profile } from '$lib/stores/profile';
	import { protocolTypes, unmarshalConfig } from '$lib/types';
	import { ConnectError } from '@connectrpc/connect';
	import { Bot, Server } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
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
	import DnsProviderSelect from '../forms/DNSProviderSelect.svelte';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { dnsProviders } from '$lib/stores/realtime';

	interface Props {
		item: Router;
		open?: boolean;
	}

	let { item = $bindable(), open = $bindable(false) }: Props = $props();
	let service = $state({} as Service);
	let hasLoadedService = $state(false);

	// Reset state when modal closes
	$effect(() => {
		if (!open) {
			service = {} as Service;
			hasLoadedService = false;
		}
	});

	// Only fetch service when modal opens with existing router (once)
	$effect(() => {
		if (item.id && open && !hasLoadedService) {
			hasLoadedService = true;
			serviceClient
				.getService({
					profileId: profile.id,
					type: item.type,
					identifier: {
						value: item.name,
						case: 'name'
					}
				})
				.then((data) => {
					service = data.service ?? ({} as Service);
				})
				.catch(() => {
					service = {} as Service;
				});
		}
	});
	$effect(() => {
		if (item.profileId) {
			service.profileId = item.profileId;
		}
		if (item.name) {
			service.name = item.name;
		}
		if (item.enabled) {
			service.enabled = item.enabled;
		}
		if (item.type) {
			service.type = item.type;
		}
	});

	const handleSubmit = async () => {
		if (!profile.id) return;

		try {
			if (item.id) {
				await routerClient.updateRouter({
					id: item.id,
					name: item.name,
					config: item.config,
					enabled: item.enabled,
					type: item.type,
					dnsProviders: item.dnsProviders
				});
				toast.success(`Router ${item.name} updated successfully`);
			} else {
				await routerClient.createRouter({
					profileId: profile.id,
					name: item.name,
					config: item.config,
					enabled: item.enabled,
					type: item.type
				});
				toast.success(`Router ${item.name} created successfully`);
			}

			if (service.id) {
				await serviceClient.updateService({
					id: service.id,
					name: service.name,
					config: service.config,
					type: service.type,
					enabled: service.enabled
				});
			} else {
				await serviceClient.createService({
					profileId: profile.id,
					name: service.name,
					config: service.config,
					type: service.type,
					enabled: service.enabled
				});
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error(`Failed to ${item.id ? 'update' : 'save'} router`, {
				description: e.message
			});
		}
		open = false;
	};

	const handleDelete = async () => {
		if (!item.id || !item.type) return;

		try {
			await routerClient.deleteRouter({ id: item.id, type: item.type });
			toast.success('Router deleted successfully');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
		open = false;
	};

	// Get service preview for agent-managed router
	const getServicePreview = () => {
		if (!service?.config) return 'No service configured';

		if (service.type === ProtocolType.HTTP) {
			const config = unmarshalConfig(service.config) as HTTPService;
			const servers = config.loadBalancer?.servers || [];
			return servers.length > 0 ? servers[0].url : 'No servers';
		} else if (service.type === ProtocolType.TCP) {
			const config = unmarshalConfig(service.config) as TCPService;
			const servers = config.loadBalancer?.servers || [];
			return servers.length > 0 ? servers[0].address : 'No servers';
		} else if (service.type === ProtocolType.UDP) {
			const config = unmarshalConfig(service.config) as UDPService;
			const servers = config.loadBalancer?.servers || [];
			return servers.length > 0 ? servers[0].address : 'No servers';
		}
		return 'No service configured';
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] max-w-xl overflow-y-auto">
		{#if item.agentId}
			<div class="space-y-4">
				<div>
					<h2 class="text-l font-semibold tracking-tight">Agent-Managed Router</h2>
					<p class="text-muted-foreground text-sm">Automatically managed via Docker labels</p>
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
				{#if $dnsProviders.length > 0}
					<Card.Root>
						<Card.Header>
							<div class="flex items-center justify-between">
								<div>
									<Card.Title class="flex items-center gap-2">DNS Providers</Card.Title>
									<Card.Description>Manage DNS providers for this router</Card.Description>
								</div>
								<DnsProviderSelect
									bind:item
									disabled={item.type === ProtocolType.UDP || !item.id}
								/>
							</div>
						</Card.Header>
						{#if item.dnsProviders?.length > 0}
							<Card.Content>
								<div class="flex flex-wrap gap-2">
									{#each item.dnsProviders as provider (provider.id)}
										<Badge variant="secondary">
											{provider.name}
										</Badge>
									{/each}
								</div>
							</Card.Content>
						{/if}
					</Card.Root>
				{/if}

				<!-- Router Info (Read-only) -->
				<Card.Root class="bg-muted/30">
					<Card.Header>
						<Card.Title class="text-muted-foreground flex items-center gap-2">
							<Server class="h-4 w-4" />
							Router Configuration
						</Card.Title>
						<Card.Description>View-only configuration managed by agent</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid grid-cols-2 gap-6">
							<div class="space-y-1">
								<Label class="text-muted-foreground">Router Name</Label>
								<p class="font-medium">{item.name || 'Not set'}</p>
							</div>
							<div class="space-y-1">
								<Label class="text-muted-foreground">Protocol</Label>
								<Badge variant="outline">
									{protocolTypes.find((t) => t.value === item.type)?.label || 'Unknown'}
								</Badge>
							</div>
							<div class="space-y-1">
								<Label class="text-muted-foreground">Service Endpoint</Label>
								<p class="text-sm font-medium">{getServicePreview()}</p>
							</div>
							<div class="space-y-1">
								<Label class="text-muted-foreground">Status</Label>
								<Badge variant={item.enabled ? 'default' : 'secondary'}>
									{item.enabled ? 'Enabled' : 'Disabled'}
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
								<Card.Description class="text-muted-foreground text-sm">
									Define how traffic is routed to your service
								</Card.Description>
							</div>

							<!-- DNS Providers -->
							<DnsProviderSelect bind:item disabled={item.type === ProtocolType.UDP || !item.id} />
						</Card.Header>
						<Card.Content class="space-y-4">
							<div class="grid w-full grid-cols-1 gap-4 sm:grid-cols-3 sm:gap-2">
								<div class="flex flex-col gap-2 {item.id ? 'sm:col-span-3' : 'sm:col-span-2'}">
									<Label for="name">Name</Label>
									<Input
										id="name"
										bind:value={item.name}
										required
										placeholder="Router Name"
										class="truncate"
									/>
								</div>

								{#if !item.id}
									<div class="flex flex-col gap-2 sm:col-span-1">
										<Label for="type">Protocol</Label>
										<Select.Root
											type="single"
											name="type"
											value={item.type?.toString()}
											onValueChange={(value) => {
												// Reset config
												item.type = parseInt(value, 10);
												switch (item.type) {
													case ProtocolType.HTTP:
														item.config = {} as HTTPRouter;
														break;
													case ProtocolType.TCP:
														item.config = {} as TCPRouter;
														break;
													case ProtocolType.UDP:
														item.config = {} as UDPRouter;
														break;
												}
											}}
										>
											<Select.Trigger class="w-full">
												<span class="truncate">
													{protocolTypes.find((t) => t.value === item.type)?.label ??
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

							{#if item.type === ProtocolType.HTTP}
								<HTTPRouterForm bind:router={item} />
							{/if}
							{#if item.type === ProtocolType.TCP}
								<TCPRouterForm bind:router={item} />
							{/if}
							{#if item.type === ProtocolType.UDP}
								<UDPRouterForm bind:router={item} />
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
							{#if item.type === ProtocolType.HTTP}
								<HTTPServiceForm bind:service />
							{/if}
							{#if item.type === ProtocolType.TCP}
								<TCPServiceForm bind:service />
							{/if}
							{#if item.type === ProtocolType.UDP}
								<UDPServiceForm bind:service />
							{/if}
						</Card.Content>
					</Card.Root>
				</Tabs.Content>
			</Tabs.Root>

			<div class="flex w-full flex-col gap-2 sm:flex-row">
				{#if item.id}
					<Button type="button" variant="destructive" onclick={handleDelete} class="flex-1">
						Delete
					</Button>
				{/if}
				<Button type="submit" class="flex-1 text-sm" onclick={handleSubmit}>
					{item.id ? 'Update' : 'Create'}
				</Button>
			</div>
		{/if}
	</Dialog.Content>
</Dialog.Root>
