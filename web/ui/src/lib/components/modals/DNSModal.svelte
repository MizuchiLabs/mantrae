<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import {
		DNSProviderType,
		type DNSProvider,
		type DNSProviderConfig
	} from '$lib/gen/mantrae/v1/dns_provider_pb';
	import { dnsProviderTypes } from '$lib/types';
	import { CircleQuestionMark } from '@lucide/svelte';
	import Badge from '../ui/badge/badge.svelte';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import { dns } from '$lib/api/dns.svelte';
	import { util } from '$lib/api/util.svelte';

	interface Props {
		data?: DNSProvider;
		open?: boolean;
	}
	let { data, open = $bindable(false) }: Props = $props();

	let dnsData = $state({} as DNSProvider);
	$effect(() => {
		if (data) dnsData = { ...data };
	});
	$effect(() => {
		if (!open) dnsData = {} as DNSProvider;
	});

	const currentIP = $derived(util.ip());
	const createMutation = dns.create();
	const updateMutation = dns.update();
	function onsubmit() {
		if (dnsData.id) {
			updateMutation.mutate({ ...dnsData });
		} else {
			createMutation.mutate({ ...dnsData });
		}
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[500px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{dnsData?.id ? 'Edit' : 'Add'} DNS Provider</Dialog.Title>
			<Dialog.Description>
				Configure automated DNS record management for your domains
			</Dialog.Description>
		</Dialog.Header>

		<form {onsubmit} class="space-y-4">
			<!-- Basic Configuration -->
			<div class="space-y-4">
				<div class="grid grid-cols-3 gap-2">
					<div class="col-span-2 space-y-2">
						<Label for="name" class="text-sm">Provider Name</Label>
						<Input id="name" bind:value={dnsData.name} required placeholder="e.g., Cloudflare" />
						<p class="text-xs text-muted-foreground">Friendly name for this provider</p>
					</div>

					<div class="space-y-2">
						<Label for="type" class="text-sm">Type</Label>
						<Select.Root
							type="single"
							name="type"
							value={dnsData.type?.toString()}
							onValueChange={(value) => (dnsData.type = parseInt(value, 10))}
						>
							<Select.Trigger>
								{dnsProviderTypes.find((t) => t.value === dnsData.type)?.label ?? 'Select type'}
							</Select.Trigger>
							<Select.Content class="no-scrollbar max-h-75 overflow-y-auto">
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
						<CustomSwitch bind:checked={dnsData.isDefault} size="md" />
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
							checked={dnsData.config?.autoUpdate}
							onCheckedChange={(value) => {
								if (dnsData.config === undefined) dnsData.config = {} as DNSProviderConfig;
								dnsData.config.autoUpdate = value;
							}}
							size="md"
						/>
					</div>

					<!-- Cloudflare Proxy -->
					{#if dnsData.type === DNSProviderType.DNS_PROVIDER_TYPE_CLOUDFLARE}
						<div class="flex items-center justify-between rounded-lg border p-3">
							<div class="space-y-1">
								<Label class="text-sm">Cloudflare Proxy</Label>
								<p class="text-xs text-muted-foreground">Enable Cloudflare's proxy service</p>
							</div>
							<CustomSwitch
								checked={dnsData.config?.proxied}
								onCheckedChange={(value) => {
									if (dnsData.config === undefined) dnsData.config = {} as DNSProviderConfig;
									dnsData.config.proxied = value;
								}}
								size="md"
							/>
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

				{#if dnsData.config?.autoUpdate}
					<div class="rounded-lg border p-3">
						<div class="space-y-2">
							<Label class="text-sm">Detected Public IP</Label>
							{#if currentIP.isSuccess}
								<div class="flex items-center gap-2">
									{#if currentIP.data?.ipv4}
										<Badge variant="secondary">{currentIP.data?.ipv4}</Badge>
									{/if}
									{#if currentIP.data?.ipv6}
										<Badge variant="secondary">{currentIP.data?.ipv6}</Badge>
									{/if}
								</div>
								<p class="text-xs text-muted-foreground">Automatically detected and updated</p>
							{/if}
						</div>
					</div>
				{:else}
					<div class="space-y-2">
						<Label for="ip" class="text-sm">IP Address</Label>
						<Input
							id="ip"
							name="ip"
							type="text"
							value={dnsData.config?.ip}
							oninput={(e) => {
								let input = e.target as HTMLInputElement;
								if (!input.value) return;
								if (dnsData.config === undefined) dnsData.config = {} as DNSProviderConfig;
								dnsData.config.ip = input.value;
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
							value={dnsData.config?.apiKey}
							oninput={(e) => {
								let input = e.target as HTMLInputElement;
								if (!input.value) return;
								if (dnsData.config === undefined) dnsData.config = {} as DNSProviderConfig;
								dnsData.config.apiKey = input.value;
							}}
							placeholder="Enter your API key"
						/>
						<p class="text-xs text-muted-foreground">API key from your DNS provider</p>
					</div>

					{#if dnsData.type === DNSProviderType.DNS_PROVIDER_TYPE_POWERDNS || dnsData.type === DNSProviderType.DNS_PROVIDER_TYPE_TECHNITIUM || dnsData.type === DNSProviderType.DNS_PROVIDER_TYPE_PIHOLE}
						<div class="space-y-2">
							<Label for="apiUrl" class="text-sm">API Endpoint</Label>
							<Input
								id="apiUrl"
								name="apiUrl"
								type="text"
								value={dnsData.config?.apiUrl}
								oninput={(e) => {
									let input = e.target as HTMLInputElement;
									if (!input.value) return;
									if (dnsData.config === undefined) dnsData.config = {} as DNSProviderConfig;
									dnsData.config.apiUrl = input.value;
								}}
								placeholder="https://dns.example.com/api"
								required
							/>
							<p class="text-xs text-muted-foreground">
								{dnsProviderTypes.find((t) => t.value === dnsData.type)?.label} server endpoint
							</p>
						</div>
					{/if}
				</div>
			</div>

			<Separator />

			<Button type="submit" class="w-full">{dnsData.id ? 'Update' : 'Create'}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
