<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Provider } from '$lib/types/provider';
	import { updateProvider, provider } from '$lib/api';

	export let p: Provider;
	let open = false;
	let oldName = p.name;
	let providerCompare = Object.keys($provider).filter((prov) => prov !== p.name);

	const update = async () => {
		if (p.name === '' || p.type === '' || p.key === '' || p.externalIP === '') return;
		// check if url starts with http:// or https://
		if (p.type === 'powerdns' && !p.url?.startsWith('http://') && !p.url?.startsWith('https://')) {
			p.url = 'http://' + p.url;
		}
		updateProvider(oldName, p);
		oldName = p.name;
		open = false;
	};

	// Check if provider name is taken
	let isNameTaken = false;
	$: isNameTaken = providerCompare.some((prov) => prov === p.name);

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
				<Card.Title>DNS Provider</Card.Title>
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
							class={isNameTaken
								? 'col-span-3 border-red-400 focus-visible:ring-0 focus-visible:ring-offset-0'
								: 'col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0'}
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
							bind:value={p.externalIP}
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
								bind:value={p.url}
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
							bind:value={p.key}
							placeholder="API Key of the provider"
							required
						/>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => update()}>Update</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
