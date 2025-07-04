<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import { slide } from 'svelte/transition';
	import Badge from '../ui/badge/badge.svelte';
	import { CircleHelp } from '@lucide/svelte';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import {
		DnsProviderType,
		type DnsProvider,
		type DnsProviderConfig
	} from '$lib/gen/mantrae/v1/dns_provider_pb';
	import { dnsClient, utilClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { dnsProviderTypes } from '$lib/types';

	interface Props {
		data: DnsProvider[];
		item: DnsProvider;
		open?: boolean;
	}

	let { data = $bindable(), item = $bindable(), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			if (item.id) {
				await dnsClient.updateDnsProvider({
					id: item.id,
					name: item.name,
					type: item.type,
					config: item.config,
					isDefault: item.isDefault
				});
				toast.success('DNS Provider updated successfully');
			} else {
				await dnsClient.createDnsProvider({
					name: item.name,
					type: item.type,
					config: item.config,
					isDefault: item.isDefault
				});
				toast.success('DNS Provider created successfully');
			}

			// Refresh data
			let response = await dnsClient.listDnsProviders({
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.dnsProviders;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to save dnsProvider', {
				description: e.message
			});
		}
		open = false;
	};

	const handleDelete = async () => {
		if (!item.id) return;

		try {
			await dnsClient.deleteDnsProvider({ id: item.id });
			toast.success('EntryPoint deleted successfully');

			// Refresh data
			let response = await dnsClient.listDnsProviders({
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.dnsProviders;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete entry point', { description: e.message });
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[500px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{item?.id ? 'Edit' : 'Add'} DNS Provider</Dialog.Title>
			<Dialog.Description>Setup dns provider for automated dns records</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="flex flex-col gap-4">
			<div class="grid w-full grid-cols-3 gap-2">
				<div class="col-span-2 flex flex-col gap-2">
					<Label for="name">Name</Label>
					<Input id="name" bind:value={item.name} required placeholder="cloudflare" />
				</div>

				<div class="col-span-1 flex flex-col gap-2">
					<Label for="current" class="text-right">Type</Label>
					<Select.Root
						type="single"
						name="type"
						value={item.type?.toString()}
						onValueChange={(value) => (item.type = parseInt(value, 10))}
					>
						<Select.Trigger class="w-full">
							{dnsProviderTypes.find((t) => t.value === item.type)?.label ?? 'Select type'}
						</Select.Trigger>
						<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
							{#each dnsProviderTypes as t (t.value)}
								<Select.Item value={t.value.toString()} label={t.label}>
									{t.label}
								</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
			</div>

			<div class="flex items-center justify-between gap-2">
				<Label for="autoUpdate" class="flex flex-row items-center gap-1 text-sm font-medium">
					Set as default
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<CircleHelp size={16} />
							</Tooltip.Trigger>
							<Tooltip.Content align="start" class="w-64">
								<p>
									If enabled, this DNS provider will be used as the default DNS provider for all
									newly created routers.
								</p>
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</Label>
				<Tabs.Root
					class="flex flex-col gap-2"
					value={item.isDefault ? 'on' : 'off'}
					onValueChange={(value) => {
						if (item.isDefault === undefined) item.isDefault = value === 'on';
						else item.isDefault = value === 'on';
					}}
				>
					<div class="flex justify-end" transition:slide={{ duration: 200 }}>
						<Tabs.List class="h-8">
							<Tabs.Trigger value="on" class="px-2 py-0.5 font-bold">On</Tabs.Trigger>
							<Tabs.Trigger value="off" class="px-2 py-0.5 font-bold">Off</Tabs.Trigger>
						</Tabs.List>
					</div>
				</Tabs.Root>
			</div>

			{#if item.type === DnsProviderType.CLOUDFLARE}
				<div class="flex items-center justify-between gap-2">
					<Label for="autoUpdate" class="flex flex-row items-center gap-1 text-sm font-medium">
						Cloudflare Proxy
					</Label>
					<Tabs.Root
						class="flex flex-col gap-2"
						value={item.config?.proxied ? 'on' : 'off'}
						onValueChange={(value) => {
							if (item.config === undefined) item.config = {} as DnsProviderConfig;
							item.config.proxied = value === 'on';
						}}
					>
						<div class="flex justify-end" transition:slide={{ duration: 200 }}>
							<Tabs.List class="h-8">
								<Tabs.Trigger value="on" class="px-2 py-0.5 font-bold">On</Tabs.Trigger>
								<Tabs.Trigger value="off" class="px-2 py-0.5 font-bold">Off</Tabs.Trigger>
							</Tabs.List>
						</div>
					</Tabs.Root>
				</div>
			{/if}

			<div class="flex items-center justify-between gap-2">
				<Label for="autoUpdate" class="flex flex-row items-center gap-1 text-sm font-medium">
					Auto Update IP
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<CircleHelp size={16} />
							</Tooltip.Trigger>
							<Tooltip.Content align="start" class="w-64">
								<p>
									When enabled, Mantrae will automatically detect and use your server's public IP
									address. DNS records will be kept in sync as your IP changes.
								</p>
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</Label>
				<Tabs.Root
					class="flex flex-col gap-2"
					value={item.config?.autoUpdate ? 'on' : 'off'}
					onValueChange={(value) => {
						if (item.config === undefined) item.config = {} as DnsProviderConfig;
						item.config.autoUpdate = value === 'on';
					}}
				>
					<div class="flex justify-end" transition:slide={{ duration: 200 }}>
						<Tabs.List class="h-8">
							<Tabs.Trigger value="on" class="px-2 py-0.5 font-bold">On</Tabs.Trigger>
							<Tabs.Trigger value="off" class="px-2 py-0.5 font-bold">Off</Tabs.Trigger>
						</Tabs.List>
					</div>
				</Tabs.Root>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="ip">
					IP
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<CircleHelp size={16} />
							</Tooltip.Trigger>
							<Tooltip.Content align="start" class="w-64">
								<p>IP used for DNS records. Should be the IP of your Traefik instance.</p>
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</Label>

				{#if item.config?.autoUpdate}
					{#await utilClient.getPublicIP({}) then value}
						<div class="flex items-center gap-2">
							{#if value?.ipv4}
								<Badge variant="secondary">{value?.ipv4}</Badge>
							{/if}
							{#if value?.ipv6}
								<Badge variant="secondary">{value?.ipv6}</Badge>
							{/if}
						</div>
					{/await}
				{:else}
					<Input
						name="ip"
						type="text"
						value={item.config?.ip}
						oninput={(e) => {
							let input = e.target as HTMLInputElement;
							if (!input.value) {
								return;
							}
							if (item.config === undefined) item.config = {} as DnsProviderConfig;
							item.config.ip = input.value;
						}}
						placeholder="IP used for DNS records"
						required
					/>
				{/if}
			</div>

			<div class="flex flex-col gap-2">
				<Label for="key">API Key</Label>
				<PasswordInput
					value={item.config?.apiKey}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						if (!input.value) {
							return;
						}
						if (item.config === undefined) item.config = {} as DnsProviderConfig;
						item.config.apiKey = input.value;
					}}
				/>
			</div>
			{#if item.type === DnsProviderType.POWERDNS || item.type === DnsProviderType.TECHNITIUM}
				<div class="flex flex-col gap-2">
					<Label for="url">Endpoint</Label>
					<Input
						name="url"
						type="text"
						value={item.config?.apiUrl}
						oninput={(e) => {
							let input = e.target as HTMLInputElement;
							if (!input.value) {
								return;
							}
							if (item.config === undefined) item.config = {} as DnsProviderConfig;
							item.config.apiUrl = input.value;
						}}
						placeholder="Endpoint for {dnsProviderTypes.find((t) => t.value === item.type)?.label}"
						required
					/>
				</div>
			{/if}

			{#if item.type === DnsProviderType.TECHNITIUM}
				<div class="flex items-center justify-between gap-2">
					<Label for="zoneType">Zone Type</Label>
					<Tabs.Root
						class="flex flex-col gap-2"
						value={item.config?.zoneType || 'primary'}
						onValueChange={(value) => {
							if (item.config === undefined) item.config = {} as DnsProviderConfig;
							item.config.zoneType = value;
						}}
					>
						<div class="flex justify-end" transition:slide={{ duration: 200 }}>
							<Tabs.List class="h-8">
								<Tabs.Trigger value="primary" class="px-2 py-0.5 font-bold">Primary</Tabs.Trigger>
								<Tabs.Trigger value="forward" class="px-2 py-0.5 font-bold">Forward</Tabs.Trigger>
							</Tabs.List>
						</div>
					</Tabs.Root>
				</div>
			{/if}

			<Separator />

			<div class="flex w-full flex-row gap-2">
				{#if item.id}
					<Button type="button" variant="destructive" onclick={handleDelete} class="flex-1">
						Delete
					</Button>
				{/if}
				<Button type="submit" class="flex-1">
					{item.id ? 'Update' : 'Create'}
				</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
