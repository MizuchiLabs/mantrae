<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Switch } from '$lib/components/ui/switch';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { DNSProvider } from '$lib/types/provider';
	import { updateProvider } from '$lib/api';

	export let p: DNSProvider;
	let open = false;

	const update = async () => {
		if (p.name === '' || p.type === '' || p.api_key === '' || p.external_ip === '') return;
		// check if url starts with http:// or https://
		if (
			p.type === 'powerdns' &&
			!p.api_url?.startsWith('http://') &&
			!p.api_url?.startsWith('https://')
		) {
			p.api_url = 'http://' + p.api_url;
		}
		await updateProvider(p);
		open = false;
	};

	const onKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			update();
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		<Button class="w-full bg-orange-400 text-black">
			<span>Edit</span>
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[500px]">
		<Card.Root class="mt-4">
			<Card.Header>
				<Card.Title class="flex items-center justify-between gap-2">
					<span>DNS Provider</span>
					<div class="flex items-center gap-2">
						<Badge variant="secondary" class="bg-blue-400">
							{p.type}
						</Badge>
						{#if p.is_active}
							<iconify-icon icon="fa6-solid:star" class="text-yellow-400" />
						{/if}
					</div>
				</Card.Title>
				<Card.Description>Update your DNS provider.</Card.Description>
			</Card.Header>
			<Card.Content>
				<div class="space-y-2" on:keydown={onKeydown} aria-hidden>
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="name" class="text-right">Name</Label>
						<Input
							name="name"
							type="text"
							bind:value={p.name}
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							placeholder="Your profile name"
							required
						/>
					</div>
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="externalIP" class="text-right">External IP</Label>
						<Input
							name="externalIP"
							type="text"
							placeholder="The public IP address of the traefik instance"
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							bind:value={p.external_ip}
							required
						/>
					</div>
					{#if p.type === 'powerdns'}
						<div class="grid grid-cols-4 items-center gap-4">
							<Label for="url" class="text-right">PowerDNS URL</Label>
							<Input
								name="url"
								type="text"
								placeholder="http://127.0.0.1:8081"
								class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
								bind:value={p.api_url}
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
							bind:value={p.api_key}
							placeholder="API Key of the provider"
							required
						/>
					</div>
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="is_active" class="text-right">Default</Label>
						<Switch name="is_active" class="col-span-3" bind:checked={p.is_active} required />
					</div>
				</div>
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => update()}>Update</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
