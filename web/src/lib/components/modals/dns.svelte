<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { type DNSProvider, DNSProviderTypes } from '$lib/types';
	import { Toggle } from '$lib/components/ui/toggle';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api, loading } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';

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
			zoneType: 'primary'
		}
	};

	let { dns = $bindable(defaultDNS), open = $bindable(false) }: Props = $props();

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
			<div class="space-y-1">
				<Label for="traefikIp">Traefik IP</Label>
				<Input
					name="traefikIp"
					type="text"
					bind:value={dns.config.traefikIp}
					placeholder="IP of your Traefik instance"
					required
				/>
			</div>

			<div class="space-y-1">
				<Label for="key">API Key</Label>
				<Input
					name="key"
					type="text"
					bind:value={dns.config.apiKey}
					placeholder="API Key for {dns.type}"
					required
				/>
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
