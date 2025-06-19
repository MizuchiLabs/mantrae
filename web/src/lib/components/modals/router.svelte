<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Toggle } from '$lib/components/ui/toggle';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import RouterForm from '../forms/router.svelte';
	import ServiceForm from '../forms/service.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { Service } from '$lib/gen/mantrae/v1/service_pb';
	import { routerClient, serviceClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { Star, Globe, CircleCheck } from '@lucide/svelte';
	import { routerTypes } from '$lib/types';
	import HttpRouter from '../forms/httpRouter.svelte';
	import { profile } from '$lib/stores/profile';
	import { pageIndex, pageSize } from '$lib/stores/common';

	interface Props {
		data: Router[];
		router: Router;
		open?: boolean;
	}

	let { data = $bindable(), router = $bindable(), open = $bindable(false) }: Props = $props();
	// let service = $state({} as Service);
	// function getService(router: Router): Service {
	//     if (!router.id) return {} as Service;
	//
	//     const httpServices = $derived(serviceClient.listServices({ profileId: profile.id, type: ServiceType.HTTP }));
	//     const tcpServices = $derived(serviceClient.listServices({ profileId: profile.id, type: ServiceType.TCP }));
	//     const udpServices = $derived(serviceClient.listServices({ profileId: profile.id, type: ServiceType.UDP }));
	//         switch (router.type) {
	//             case RouterType.HTTP:
	//                 return $derived(httpServices.filter((s) => s.id === router.serviceId));
	//             case RouterType.TCP:
	//                 return $derived(tcpServices.filter((s) => s.id === router.serviceId));
	//             case RouterType.UDP:
	//                 return $derived(udpServices.filter((s) => s.id === router.serviceId));
	//             default:
	//                 return {} as Service;
	//         }
	// }

	const update = async () => {
		if (!profile.id) return;

		try {
			if (router.id) {
				await routerClient.updateRouter({
					id: router.id,
					name: router.name,
					config: router.config,
					enabled: router.enabled,
					type: router.type
				});
				toast.success(`Router ${router.name} updated successfully`);
			} else {
				await routerClient.createRouter({
					profileId: profile.id,
					name: router.name,
					config: router.config,
					enabled: router.enabled,
					type: router.type
				});
				toast.success(`Router ${router.name} created successfully`);
			}

			// if (service.id) {
			// 	await serviceClient.updateService({
			// 		id: service.id,
			// 		name: service.name,
			// 		config: service.config,
			// 		type: service.type
			// 	});
			// 	toast.success(`Service ${service.name} updated successfully`);
			// } else {
			// 	await serviceClient.createService({
			// 		name: service.name,
			// 		config: service.config,
			// 		type: service.type
			// 	});
			// 	toast.success(`Service ${service.name} created successfully`);
			// }

			// Refresh data
			const response = await routerClient.listRouters({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.routers;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error(`Failed to ${router.id ? 'update' : 'save'} router`, {
				description: e.message
			});
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] max-w-xl overflow-y-auto">
		<Tabs.Root value="router" class="mt-4">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="router">Router</Tabs.Trigger>
				<Tabs.Trigger value="service">Service</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="router">
				<Card.Root>
					<Card.Header class="flex flex-row items-center justify-between">
						<div>
							<Card.Title>{router.id ? 'Update' : 'Create'} Router</Card.Title>
							<Card.Description>
								{router.id ? 'Update existing router' : 'Create a new router'}
							</Card.Description>
						</div>

						<Select.Root
							type="single"
							name="router_type"
							value={router.type?.toString()}
							onValueChange={(value) => (router.type = parseInt(value, 10))}
						>
							<Select.Trigger class="w-[120px]">
								{routerTypes.find((t) => t.value === router.type)?.label ?? 'Select type'}
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

						<!-- DNS Providers -->
						<!-- {#await dnsClient.listDnsProviders({}) then value} -->
						<!-- 	<Tooltip.Provider> -->
						<!-- 		<Tooltip.Root> -->
						<!-- 			<Tooltip.Trigger> -->
						<!-- 				<div bind:this={dnsAnchor}> -->
						<!-- 					<Button -->
						<!-- 						variant="ghost" -->
						<!-- 						size="sm" -->
						<!-- 						class="flex items-center gap-2" -->
						<!-- 						onclick={() => (selectDNSOpen = true)} -->
						<!-- 					> -->
						<!-- 						<Globe size={16} /> -->
						<!-- 						<Badge -->
						<!-- 							>{value.dnsProviders.length -->
						<!-- 								? value.dnsProviders.join(', ') -->
						<!-- 								: 'None'}</Badge -->
						<!-- 						> -->
						<!-- 					</Button> -->
						<!-- 				</div> -->
						<!-- 			</Tooltip.Trigger> -->
						<!-- 			<Tooltip.Content side="left" align="center"> -->
						<!-- 				<p>Select DNS Provider</p> -->
						<!-- 			</Tooltip.Content> -->
						<!-- 		</Tooltip.Root> -->
						<!-- 	</Tooltip.Provider> -->
						<!---->
						<!-- 	<Select.Root -->
						<!-- 		type="multiple" -->
						<!-- 		value={value.dnsProviders.map((item) => item.id.toString())} -->
						<!-- 		onValueChange={handleDNSProviderChange} -->
						<!-- 		bind:open={selectDNSOpen} -->
						<!-- 	> -->
						<!-- 		<Select.Content customAnchor={dnsAnchor} align="end"> -->
						<!-- 			{#each value.dnsProviders as dns (dns.id)} -->
						<!-- 				<Select.Item value={dns.id.toString()} class="flex items-center gap-2"> -->
						<!-- 					{dns.name} ({dns.type}) -->
						<!-- 					{#if dns.isActive} -->
						<!-- 						<CircleCheck size="1rem" class="text-green-400" /> -->
						<!-- 					{/if} -->
						<!-- 				</Select.Item> -->
						<!-- 			{/each} -->
						<!-- 		</Select.Content> -->
						<!-- 	</Select.Root> -->
						<!-- {/await} -->
					</Card.Header>
					<Card.Content class="flex flex-col gap-2">
						<div class="flex flex-col gap-2">
							<Label for="name">Name</Label>
							<Input id="name" bind:value={router.name} placeholder="Router" required />
						</div>
						{#if router.type === RouterType.HTTP}
							<HttpRouter bind:router />
						{/if}
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
			<Tabs.Content value="service">
				<!-- <ServiceForm bind:service bind:router /> -->
			</Tabs.Content>
		</Tabs.Root>

		<Separator />

		<Button type="submit" onclick={update}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
