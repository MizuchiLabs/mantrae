<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { type DNSProvider, DNSProviderTypes, type PublicIP } from '$lib/types';
	import { Toggle } from '$lib/components/ui/toggle';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api, loading } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import { slide } from 'svelte/transition';
	import { onMount } from 'svelte';
	import Badge from '../ui/badge/badge.svelte';
	import { CircleHelp } from 'lucide-svelte';
	import PasswordInput from '../ui/password-input/password-input.svelte';

	interface Props {
		dns: DNSProvider | undefined;
		open?: boolean;
	}
	const defaultDNS: DNSProvider = {
		id: 0,
		name: '',
		type: DNSProviderTypes.CLOUDFLARE,
		isActive: false,
		config: {
			apiKey: '',
			apiUrl: '',
			traefikIp: '',
			proxied: false,
			autoUpdate: false,
			zoneType: 'primary'
		}
	};

	let { dns = $bindable(defaultDNS), open = $bindable(false) }: Props = $props();
	let currentIP: PublicIP | undefined = $state();
	const dnsProviders = Object.entries(DNSProviderTypes).map(([key, value]) => ({
		label: value.charAt(0).toUpperCase() + value.slice(1),
		value: key.toLowerCase()
	}));
	const handleSubmit = async () => {
		try {
			if (dns.id) {
				await api.updateDNSProvider(dns);
				toast.success('DNS Provider updated successfully');
			} else {
				await api.createDNSProvider(dns);
				toast.success('DNS Provider created successfully');
			}
			open = false;
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to save dnsProvider', {
				description: e.message
			});
		}
		dns = defaultDNS;
	};

	onMount(async () => {
		currentIP = await api.getIP();
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[550px]">
		<Dialog.Header>
			<Dialog.Title>{dns.id ? 'Edit' : 'Add'} DNS Provider</Dialog.Title>
			<Dialog.Description>Setup dns provider for automated dns records</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="space-y-4">
			<div class="mb-4 flex items-center justify-end gap-2">
				<Label for="is_active" class="text-right">Default</Label>
				<Switch name="is_active" bind:checked={dns.isActive} />
			</div>

			<div class="space-y-1 pb-2">
				<Label for="current" class="text-right">Type</Label>
				<Select.Root type="single" value={dns.type} onValueChange={(value) => (dns.type = value)}>
					<Select.Trigger class="col-span-3">
						{dns.type ? dns.type : 'Select type'}
					</Select.Trigger>
					<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
						{#each dnsProviders as type}
							<Select.Item value={type.value} label={type.label}>
								{type.label}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>

			<Separator />

			{#if dns.type === DNSProviderTypes.TECHNITIUM}
				<div class="flex items-center justify-end gap-2">
					<Toggle
						size="sm"
						pressed={dns.config.zoneType === 'primary'}
						onPressedChange={() => (dns.config.zoneType = 'primary')}
						class="font-bold data-[state=on]:bg-green-300 dark:data-[state=on]:text-black"
					>
						Primary
					</Toggle>
					<Toggle
						size="sm"
						pressed={dns.config.zoneType === 'forwarder'}
						onPressedChange={() => (dns.config.zoneType = 'forwarder')}
						class="font-bold data-[state=on]:bg-blue-300 dark:data-[state=on]:text-black"
					>
						Forwarder
					</Toggle>
				</div>
			{/if}

			<div class="flex items-center justify-between gap-2 py-2">
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
					value={dns.config.autoUpdate ? 'on' : 'off'}
					onValueChange={(value) => (dns.config.autoUpdate = value === 'on')}
					class="flex flex-col gap-2"
				>
					<div class="flex justify-end" transition:slide={{ duration: 200 }}>
						<Tabs.List class="h-8">
							<Tabs.Trigger value="on" class="px-2 py-0.5 font-bold">On</Tabs.Trigger>
							<Tabs.Trigger value="off" class="px-2 py-0.5 font-bold">Off</Tabs.Trigger>
						</Tabs.List>
					</div>
				</Tabs.Root>
			</div>

			<div class="space-y-1">
				<Label for="name">Name</Label>
				<Input
					name="name"
					type="text"
					bind:value={dns.name}
					placeholder="Name of the provider"
					required
				/>
			</div>

			<div class="flex flex-col gap-2 py-1">
				<Label for="traefikIp">Traefik IP</Label>
				{#if dns.config.autoUpdate}
					<div class="flex items-center gap-2">
						{#if currentIP?.ipv4}
							<Badge variant="secondary">{currentIP?.ipv4}</Badge>
						{/if}
						{#if currentIP?.ipv6}
							<Badge variant="secondary">{currentIP?.ipv6}</Badge>
						{/if}
					</div>
				{:else}
					<Input
						name="traefikIp"
						type="text"
						bind:value={dns.config.traefikIp}
						placeholder="IP of your Traefik instance"
						required
					/>
				{/if}
			</div>

			<div class="space-y-1">
				<Label for="key">API Key</Label>
				<PasswordInput bind:value={dns.config.apiKey} />
			</div>
			{#if dns.type === DNSProviderTypes.POWERDNS || dns.type === DNSProviderTypes.TECHNITIUM}
				<div class="space-y-1">
					<Label for="url">Endpoint</Label>
					<Input
						name="url"
						type="text"
						bind:value={dns.config.apiUrl}
						placeholder="Endpoint for {dns.type}"
						required
					/>
				</div>
			{/if}

			{#if dns.type === DNSProviderTypes.CLOUDFLARE}
				<div class="flex items-center gap-2">
					<Label for="proxied">Proxied?</Label>
					<Switch bind:checked={dns.config.proxied} />
				</div>
			{/if}

			<Separator />

			<Button type="submit" class="w-full" disabled={$loading}>{dns.id ? 'Update' : 'Save'}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
