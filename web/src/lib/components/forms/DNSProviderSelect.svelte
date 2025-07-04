<script lang="ts">
	import { dnsClient, routerClient } from '$lib/api';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { type Router } from '$lib/gen/mantrae/v1/router_pb';
	import { ConnectError } from '@connectrpc/connect';
	import { CircleCheck, Globe } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	interface Props {
		data: Router[];
		item: Router;
		disabled?: boolean;
	}

	let { data = $bindable(), item = $bindable(), disabled = $bindable(false) }: Props = $props();

	let dnsAnchor = $state({} as HTMLElement);
	let selectDNSOpen = $state(false);

	async function handleDNSProviderChange(value: string[]) {
		if (value.length === 0) item.dnsProviders = [];
		try {
			const result = await dnsClient.listDnsProviders({ limit: -1n, offset: 0n });
			item.dnsProviders = result.dnsProviders.filter((p) => value.includes(p.id.toString()));
			await routerClient.updateRouter({
				id: item.id,
				name: item.name,
				config: item.config,
				enabled: item.enabled,
				type: item.type,
				dnsProviders: item.dnsProviders
			});
			toast.success(`Router ${item.name} updated successfully`);
		} catch (e) {
			let error = ConnectError.from(e);
			toast.error(`Failed to update router ${item.name}: ${error.message}`);
		}
	}
</script>

{#if !disabled}
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
							{#if dns.isDefault}
								<CircleCheck size="1rem" class="text-green-400" />
							{/if}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		{/if}
	{/await}
{/if}
