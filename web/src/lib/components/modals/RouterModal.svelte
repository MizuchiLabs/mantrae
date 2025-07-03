<script lang="ts">
	import { dnsClient, routerClient, serviceClient } from '$lib/api';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import { ServiceType, type Service } from '$lib/gen/mantrae/v1/service_pb';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { profile } from '$lib/stores/profile';
	import { routerTypes } from '$lib/types';
	import { ConnectError } from '@connectrpc/connect';
	import { CircleCheck, Globe } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import HTTPRouterForm from '../forms/HTTPRouterForm.svelte';
	import TCPRouterForm from '../forms/TCPRouterForm.svelte';
	import UDPRouterForm from '../forms/UDPRouterForm.svelte';
	import HTTPServiceForm from '../forms/HTTPServiceForm.svelte';
	import TCPServiceForm from '../forms/TCPServiceForm.svelte';
	import UDPServiceForm from '../forms/UDPServiceForm.svelte';

	interface Props {
		data: Router[];
		item: Router;
		open?: boolean;
	}

	let { data = $bindable(), item = $bindable(), open = $bindable(false) }: Props = $props();
	let service = $state({} as Service);

	$effect(() => {
		if (!open) {
			service = {} as Service;
		}
	});
	$effect(() => {
		if (item.id && open) {
			serviceClient
				.getServiceByRouter({
					name: item.name,
					type: item.type
				})
				.then((data) => {
					service = data.service ?? ({} as Service);
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
		if (item.type) {
			switch (item.type) {
				case RouterType.HTTP:
					service.type = ServiceType.HTTP;
					break;
				case RouterType.TCP:
					service.type = ServiceType.TCP;
					break;
				case RouterType.UDP:
					service.type = ServiceType.UDP;
					break;
			}
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
					type: service.type
				});
				toast.success(`Service ${service.name} updated successfully`);
			} else {
				await serviceClient.createService({
					profileId: profile.id,
					name: service.name,
					config: service.config,
					type: service.type
				});
				toast.success(`Service ${service.name} created successfully`);
			}

			// Refresh data
			const response = await routerClient.listRouters({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.routers;
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

			// Refresh data
			let response = await routerClient.listRouters({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.routers;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
		open = false;
	};

	let dnsAnchor = $state({} as HTMLElement);
	let selectDNSOpen = $state(false);

	async function handleDNSProviderChange(value: string[]) {
		if (value.length === 0) {
			item.dnsProviders = [];
			return;
		}
		const result = await dnsClient.listDnsProviders({ limit: -1n, offset: 0n });
		item.dnsProviders = result.dnsProviders.filter((p) => value.includes(p.id.toString()));
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] max-w-xl overflow-y-auto">
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
						<div class="space-y-1">
							<Card.Title class="text-lg sm:text-xl"
								>{item.id ? 'Update' : 'Create'} Router</Card.Title
							>
							<Card.Description class="text-sm">
								{item.id ? 'Update existing router' : 'Create a new router'}
							</Card.Description>
						</div>

						<!-- DNS Providers -->
						{#await dnsClient.listDnsProviders({ limit: -1n, offset: 0n }) then value}
							{#if value.dnsProviders.length > 0}
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
													<Badge>
														{item.dnsProviders?.length > 0
															? item.dnsProviders?.map((p) => p.name).join(', ')
															: 'None'}
													</Badge>
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
									value={item.dnsProviders?.map((item) => item.id.toString())}
									onValueChange={handleDNSProviderChange}
									bind:open={selectDNSOpen}
								>
									<Select.Content customAnchor={dnsAnchor} align="end">
										{#each value.dnsProviders as dns (dns.id)}
											<Select.Item value={dns.id.toString()} class="flex items-center gap-2">
												<span class="truncate">{dns.name}</span>
												{#if dns.isActive}
													<CircleCheck size="1rem" class="text-green-400" />
												{/if}
											</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>
							{/if}
						{/await}
					</Card.Header>
					<Card.Content class="space-y-4">
						<div class="grid w-full grid-cols-1 gap-4 sm:grid-cols-3 sm:gap-2">
							<div class="flex flex-col gap-2 sm:col-span-2">
								<Label for="name">Name</Label>
								<Input
									id="name"
									bind:value={item.name}
									required
									placeholder="Router Name"
									class="truncate"
								/>
							</div>

							<div class="flex flex-col gap-2 sm:col-span-1">
								<Label for="type" class="text-right">Protocol</Label>
								<Select.Root
									type="single"
									name="type"
									value={item.type?.toString()}
									onValueChange={(value) => (item.type = parseInt(value, 10))}
								>
									<Select.Trigger class="w-full">
										<span class="truncate">
											{routerTypes.find((t) => t.value === item.type)?.label ?? 'Select type'}
										</span>
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											<Select.Label>Router Type</Select.Label>
											{#each routerTypes as t (t.value)}
												<Select.Item value={t.value.toString()} label={t.label}>
													{t.label}
												</Select.Item>
											{/each}
										</Select.Group>
									</Select.Content>
								</Select.Root>
							</div>
						</div>

						{#if item.type === RouterType.HTTP}
							<HTTPRouterForm bind:router={item} />
						{/if}
						{#if item.type === RouterType.TCP}
							<TCPRouterForm bind:router={item} />
						{/if}
						{#if item.type === RouterType.UDP}
							<UDPRouterForm bind:router={item} />
						{/if}
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
			<Tabs.Content value="service" class="space-y-4">
				<Card.Root>
					<Card.Header>
						<Card.Title class="text-lg sm:text-xl">
							{item.id ? 'Update' : 'Create'}
							Service
						</Card.Title>
						<Card.Description class="text-sm">
							{item.id ? 'Update existing service' : 'Create a new service'}
						</Card.Description>
					</Card.Header>
					<Card.Content class="flex flex-col gap-3">
						{#if item.type === RouterType.HTTP}
							<HTTPServiceForm bind:service />
						{/if}
						{#if item.type === RouterType.TCP}
							<TCPServiceForm bind:service />
						{/if}
						{#if item.type === RouterType.UDP}
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
	</Dialog.Content>
</Dialog.Root>
