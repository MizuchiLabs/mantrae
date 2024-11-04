<script lang="ts">
	import * as Select from '$lib/components/ui/select';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Switch } from '$lib/components/ui/switch';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { DNSProvider } from '$lib/types/base';
	import { getPublicIP } from '$lib/api';
	import { Copy, Eye, EyeOff } from 'lucide-svelte';
	import type { Selected } from 'bits-ui';
	import HoverInfo from '../utils/hoverInfo.svelte';
	import { Toggle } from '../ui/toggle';
	import autoAnimate from '@formkit/auto-animate';

	export let provider: DNSProvider;
	const providerTypes: Selected<string>[] = [
		{ label: 'Cloudflare', value: 'cloudflare' },
		{ label: 'PowerDNS', value: 'powerdns' },
		{ label: 'Technitium', value: 'technitium' }
	];

	let providerType: Selected<string> | undefined = providerTypes.find(
		(t) => t.value === provider.type
	);
	const setProviderType = async (type: Selected<string> | undefined) => {
		if (type === undefined) return;
		provider.type = type.value.toLowerCase();
		if (provider.type === 'technitium') {
			provider.zoneType = 'primary';
		}
	};

	let showAPIKey = false;
</script>

<Card.Root class="mt-4">
	<Card.Header>
		<Card.Title class="flex items-center justify-between gap-2">
			<span>DNS Provider</span>
			<div class="flex items-center gap-2">
				<Badge variant="secondary" class="bg-blue-400">
					{provider.type}
				</Badge>
				{#if provider.isActive}
					<iconify-icon icon="fa6-solid:star" class="text-yellow-400" />
				{/if}
			</div>
		</Card.Title>
	</Card.Header>
	<Card.Content class="space-y-2">
		<div class="mb-4 flex items-center justify-end gap-2">
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Label for="is_active" class="text-right">Default</Label>
					<Switch name="is_active" bind:checked={provider.isActive} required />
				</Tooltip.Trigger>
				<Tooltip.Content class="max-w-sm">
					<p>Sets this provider as the default, any new router created will use this provider</p>
				</Tooltip.Content>
			</Tooltip.Root>
		</div>

		{#if provider.type === 'cloudflare'}
			<div class="flex items-center justify-end gap-2">
				<Label for="proxied" class="text-right">Proxied</Label>
				<Switch name="proxied" bind:checked={provider.proxied} required />
			</div>
		{/if}

		{#if provider.type === 'technitium'}
			<div class="flex items-center justify-end gap-1 font-mono text-sm" use:autoAnimate>
				<Toggle
					size="sm"
					pressed={provider.zoneType === 'primary'}
					onPressedChange={() => (provider.zoneType = 'primary')}
					class="font-bold data-[state=on]:bg-green-300  dark:data-[state=on]:text-black"
				>
					Primary
				</Toggle>
				<Toggle
					size="sm"
					pressed={provider.zoneType === 'forwarder'}
					onPressedChange={() => (provider.zoneType = 'forwarder')}
					class="font-bold data-[state=on]:bg-blue-300 dark:data-[state=on]:text-black"
				>
					Forwarder
				</Toggle>
			</div>
		{/if}

		<div class="grid grid-cols-4 items-center gap-4 space-y-2">
			<Label for="current" class="text-right">Type</Label>
			<Select.Root onSelectedChange={setProviderType} selected={providerType}>
				<Select.Trigger class="col-span-3">
					<Select.Value placeholder="Select a type" />
				</Select.Trigger>
				<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
					{#each providerTypes as type}
						<Select.Item value={type.value} label={type.label}>
							{type.label}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>

		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="name" class="text-right">Name</Label>
			<Input
				name="name"
				type="text"
				bind:value={provider.name}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="Name of the provider"
				required
			/>
		</div>
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="external-ip" class="col-span-1 flex items-center justify-end gap-0.5">
				IP Address
				<HoverInfo text="Use the public IP of your Traefik instance. (No schema or port)" />
			</Label>
			<div class="col-span-3 flex flex-row items-center justify-end gap-1">
				<Input
					name="external-ip"
					type="text"
					placeholder="Public IP address of Traefik"
					bind:value={provider.externalIp}
					class="pr-10"
					required
				/>
				<Tooltip.Root openDelay={500}>
					<Tooltip.Trigger class="absolute">
						<Button
							variant="ghost"
							size="icon"
							on:click={async () => {
								provider.externalIp = await getPublicIP();
							}}
							class="hover:bg-transparent hover:text-red-400"
						>
							<Copy size="1rem" />
						</Button>
					</Tooltip.Trigger>
					<Tooltip.Content side="top" align="center" class="max-w-sm">
						Use the external IP of your Traefik instance
					</Tooltip.Content>
				</Tooltip.Root>
			</div>
		</div>
		{#if provider.type === 'powerdns' || provider.type === 'technitium'}
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="url" class="col-span-1 flex items-center justify-end gap-0.5">
					Endpoint
					<HoverInfo text="The API endpoint of the provider" />
				</Label>
				<Input
					name="url"
					type="text"
					placeholder={provider.type === 'powerdns'
						? 'http://127.0.0.1:8081'
						: 'http://127.0.0.1:5380'}
					class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
					bind:value={provider.apiUrl}
					required
				/>
			</div>
		{/if}
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="key" class="text-right">API Key</Label>
			<div class="col-span-3 flex flex-row items-center justify-end gap-1">
				{#if showAPIKey}
					<Input
						name="key"
						type="text"
						bind:value={provider.apiKey}
						placeholder="API Key of the provider"
						class="pr-10"
						required
					/>
				{:else}
					<Input
						name="key"
						type="password"
						bind:value={provider.apiKey}
						placeholder="API Key of the provider"
						class="pr-10"
						required
					/>
				{/if}
				<Button
					variant="ghost"
					size="icon"
					class="absolute hover:bg-transparent hover:text-red-400"
					on:click={() => (showAPIKey = !showAPIKey)}
				>
					{#if showAPIKey}
						<Eye size="1rem" />
					{:else}
						<EyeOff size="1rem" />
					{/if}
				</Button>
			</div>
		</div>
	</Card.Content>
</Card.Root>
