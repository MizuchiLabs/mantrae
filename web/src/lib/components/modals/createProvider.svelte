<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Switch } from '$lib/components/ui/switch';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button';
	import { newProvider, type DNSProvider } from '$lib/types/provider';
	import type { Selected } from 'bits-ui';
	import { createProvider } from '$lib/api';

	let provider: DNSProvider = newProvider();
	const providerTypes: Selected<string>[] = [
		{ label: 'Cloudflare', value: 'cloudflare' },
		{ label: 'PowerDNS', value: 'powerdns' }
	];

	const create = async () => {
		if (
			provider.name === '' ||
			provider.type === '' ||
			provider.api_key === '' ||
			provider.external_ip === ''
		)
			return;
		console.log(provider);
		await createProvider(provider);
		provider = newProvider();
		providerType = providerTypes[0];
	};

	let providerType: Selected<string> | undefined = providerTypes[0];
	const setProviderType = async (type: Selected<string> | undefined) => {
		if (type === undefined) return;
		provider.type = type.value.toLowerCase();
	};
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<div class="flex w-full flex-row items-center justify-between">
			<Button class="flex items-center gap-2 bg-red-400 text-black">
				<span>Add Provider</span>
				<iconify-icon icon="fa6-solid:plus" />
			</Button>
		</div>
	</Dialog.Trigger>
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[500px]">
		<Card.Root class="mt-4">
			<Card.Header>
				<Card.Title>DNS Provider</Card.Title>
				<Card.Description>Add a new DNS provider.</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-2">
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
						placeholder="Name"
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						bind:value={provider.name}
						required
					/>
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="name" class="text-right">External IP</Label>
					<Input
						name="externalIP"
						type="text"
						placeholder="The public IP address of the traefik instance"
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						bind:value={provider.external_ip}
						required
					/>
				</div>
				{#if provider.type === 'powerdns'}
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="url" class="text-right">PowerDNS URL</Label>
						<Input
							name="url"
							type="text"
							placeholder="http://127.0.0.1:8081"
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							bind:value={provider.api_url}
							required
						/>
					</div>
				{/if}
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="key" class="text-right">API Key</Label>
					<Input
						name="key"
						type="password"
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						bind:value={provider.api_key}
						placeholder="API Key of the provider"
						required
					/>
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="is_active" class="text-right">Default</Label>
					<Switch name="is_active" class="col-span-3" bind:checked={provider.is_active} required />
				</div>
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => create()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
