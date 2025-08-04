<script lang="ts">
	import { dnsClient, utilClient } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as ToggleGroup from '$lib/components/ui/toggle-group/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import {
		DnsProviderType,
		type DnsProvider,
		type DnsProviderConfig
	} from '$lib/gen/mantrae/v1/dns_provider_pb';
	import { dnsProviderTypes } from '$lib/types';
	import { ConnectError } from '@connectrpc/connect';
	import { CircleQuestionMark } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import Badge from '../ui/badge/badge.svelte';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import Separator from '../ui/separator/separator.svelte';

	interface Props {
		item: DnsProvider;
		open?: boolean;
	}

	let { item = $bindable(), open = $bindable(false) }: Props = $props();

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
			<Dialog.Description>
				Configure automated DNS record management for your domains
			</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="space-y-4">
			<!-- Basic Configuration -->
			<div class="space-y-4">
				<div class="grid grid-cols-3 gap-2">
					<div class="col-span-2 space-y-2">
						<Label for="name" class="text-sm">Provider Name</Label>
						<Input id="name" bind:value={item.name} required placeholder="e.g., Cloudflare" />
						<p class="text-xs text-muted-foreground">Friendly name for this provider</p>
					</div>

					<div class="space-y-2">
						<Label for="type" class="text-sm">Type</Label>
						<Select.Root
							type="single"
							name="type"
							value={item.type?.toString()}
							onValueChange={(value) => (item.type = parseInt(value, 10))}
						>
							<Select.Trigger>
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
			</div>

			<!-- Provider Settings -->
			<div class="space-y-4">
				<div class="space-y-2">
					<Label class="text-sm font-medium">Provider Settings</Label>
					<p class="text-xs text-muted-foreground">Configure how this DNS provider should behave</p>
				</div>

				<div class="space-y-4">
					<!-- Default Provider -->
					<div class="flex items-center justify-between rounded-lg border p-3">
						<div class="space-y-1">
							<Label class="flex items-center gap-1 text-sm">
								Default Provider
								<Tooltip.Provider>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<CircleQuestionMark size={14} />
										</Tooltip.Trigger>
										<Tooltip.Content align="start" class="w-64">
											<p>
												If enabled, this DNS provider will be used as the default DNS provider for
												all newly created routers.
											</p>
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							</Label>
							<p class="text-xs text-muted-foreground">Use for new routers by default</p>
						</div>
						<CustomSwitch bind:checked={item.isDefault} size="md" />
					</div>

					<!-- Auto Update IP -->
					<div class="flex items-center justify-between rounded-lg border p-3">
						<div class="space-y-1">
							<Label class="flex items-center gap-1 text-sm">
								Auto Update IP
								<Tooltip.Provider>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<CircleQuestionMark size={14} />
										</Tooltip.Trigger>
										<Tooltip.Content align="start" class="w-64">
											<p>
												When enabled, Mantrae will automatically detect and use your server's public
												IP address. DNS records will be kept in sync as your IP changes.
											</p>
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							</Label>
							<p class="text-xs text-muted-foreground">Automatically sync with public IP</p>
						</div>
						<CustomSwitch
							checked={item.config?.autoUpdate}
							onCheckedChange={(value) => {
								if (item.config === undefined) item.config = {} as DnsProviderConfig;
								item.config.autoUpdate = value;
							}}
							size="md"
						/>
					</div>

					<!-- Cloudflare Proxy -->
					{#if item.type === DnsProviderType.CLOUDFLARE}
						<div class="flex items-center justify-between rounded-lg border p-3">
							<div class="space-y-1">
								<Label class="text-sm">Cloudflare Proxy</Label>
								<p class="text-xs text-muted-foreground">Enable Cloudflare's proxy service</p>
							</div>
							<CustomSwitch
								checked={item.config?.proxied}
								onCheckedChange={(value) => {
									if (item.config === undefined) item.config = {} as DnsProviderConfig;
									item.config.proxied = value;
								}}
								size="md"
							/>
						</div>
					{/if}

					<!-- Technitium Zone Type -->
					{#if item.type === DnsProviderType.TECHNITIUM}
						<div class="flex items-center justify-between rounded-lg border p-3">
							<div class="space-y-1">
								<Label class="text-sm">Zone Type</Label>
								<p class="text-xs text-muted-foreground">DNS zone configuration type</p>
							</div>
							<ToggleGroup.Root
								type="single"
								size="sm"
								value={item.config?.zoneType}
								onValueChange={(value) => {
									if (item.config === undefined) item.config = {} as DnsProviderConfig;
									item.config.zoneType = value;
								}}
								class="border"
							>
								<ToggleGroup.Item value="primary" aria-label="Toggle primary">
									<span class="text-xs">Primary</span>
								</ToggleGroup.Item>
								<ToggleGroup.Item value="forwarder" aria-label="Toggle forwarder" class="px-2">
									<span class="text-xs">Forwarder</span>
								</ToggleGroup.Item>
							</ToggleGroup.Root>
						</div>
					{/if}
				</div>
			</div>

			<Separator />

			<div class="space-y-4">
				<div class="space-y-2">
					<Label class="flex items-center gap-1 text-sm font-medium">
						Network Configuration
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger>
									<CircleQuestionMark size={14} />
								</Tooltip.Trigger>
								<Tooltip.Content align="start" class="w-64">
									<p>IP used for DNS records. Should be the IP of your Traefik instance.</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</Label>
					<p class="text-xs text-muted-foreground">
						Configure the IP address for DNS record creation
					</p>
				</div>

				{#if item.config?.autoUpdate}
					<div class="rounded-lg border p-3">
						<div class="space-y-2">
							<Label class="text-sm">Detected Public IP</Label>
							{#await utilClient.getPublicIP({}) then value}
								<div class="flex items-center gap-2">
									{#if value?.ipv4}
										<Badge variant="secondary">{value?.ipv4}</Badge>
									{/if}
									{#if value?.ipv6}
										<Badge variant="secondary">{value?.ipv6}</Badge>
									{/if}
								</div>
								<p class="text-xs text-muted-foreground">Automatically detected and updated</p>
							{/await}
						</div>
					</div>
				{:else}
					<div class="space-y-2">
						<Label for="ip" class="text-sm">IP Address</Label>
						<Input
							id="ip"
							name="ip"
							type="text"
							value={item.config?.ip}
							oninput={(e) => {
								let input = e.target as HTMLInputElement;
								if (!input.value) return;
								if (item.config === undefined) item.config = {} as DnsProviderConfig;
								item.config.ip = input.value;
							}}
							placeholder="Enter IP address for DNS records"
							required
						/>
						<p class="text-xs text-muted-foreground">Static IP address for DNS record creation</p>
					</div>
				{/if}
			</div>

			<Separator />

			<!-- Authentication -->
			<div class="space-y-4">
				<div class="space-y-2">
					<Label class="text-sm font-medium">Authentication</Label>
					<p class="text-xs text-muted-foreground">
						Provide credentials to access your DNS provider's API
					</p>
				</div>

				<div class="space-y-4">
					<div class="space-y-2">
						<Label for="apiKey" class="text-sm">API Key</Label>
						<PasswordInput
							id="apiKey"
							value={item.config?.apiKey}
							oninput={(e) => {
								let input = e.target as HTMLInputElement;
								if (!input.value) return;
								if (item.config === undefined) item.config = {} as DnsProviderConfig;
								item.config.apiKey = input.value;
							}}
							placeholder="Enter your API key"
						/>
						<p class="text-xs text-muted-foreground">API key from your DNS provider</p>
					</div>

					{#if item.type === DnsProviderType.POWERDNS || item.type === DnsProviderType.TECHNITIUM}
						<div class="space-y-2">
							<Label for="apiUrl" class="text-sm">API Endpoint</Label>
							<Input
								id="apiUrl"
								name="apiUrl"
								type="text"
								value={item.config?.apiUrl}
								oninput={(e) => {
									let input = e.target as HTMLInputElement;
									if (!input.value) return;
									if (item.config === undefined) item.config = {} as DnsProviderConfig;
									item.config.apiUrl = input.value;
								}}
								placeholder="https://dns.example.com/api"
								required
							/>
							<p class="text-xs text-muted-foreground">
								{dnsProviderTypes.find((t) => t.value === item.type)?.label} server endpoint
							</p>
						</div>
					{/if}
				</div>
			</div>

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
