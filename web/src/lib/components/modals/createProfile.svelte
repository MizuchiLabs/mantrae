<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { createProfile } from '$lib/api';
	import { newProfile } from '$lib/types/base';
	import HoverInfo from '../utils/hoverInfo.svelte';
	import { Plus } from 'lucide-svelte';

	let profile = newProfile();

	const create = async () => {
		if (profile.name === '') return;
		// Strip trailing slashes
		if (profile.url.endsWith('/')) {
			profile.url = profile.url.slice(0, -1);
		}
		await createProfile(profile);
		profile = newProfile();
	};

	const onKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			create();
		}
	};
</script>

<Dialog.Root>
	<Dialog.Trigger class="w-full">
		<div class="flex w-full flex-row items-center justify-between">
			<span>New Profile</span>
			<Button variant="ghost" class="h-8 w-8 rounded-full bg-green-400" size="icon">
				<Plus size="1rem" />
			</Button>
		</div>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[450px]">
		<Dialog.Header>
			<Dialog.Title>New profile</Dialog.Title>
			<Dialog.Description>Create a new profile to manage your Traefik instance.</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4 py-4" on:keydown={onKeydown} aria-hidden>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="name" class="text-right">Name</Label>
				<Input
					name="name"
					type="text"
					bind:value={profile.name}
					class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
					placeholder="Your profile name"
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="url" class="text-right">URL</Label>
				<Input
					name="url"
					type="text"
					class="col-span-3"
					bind:value={profile.url}
					placeholder="URL of your traefik instance"
					required
				/>
			</div>
			<div class="flex flex-row items-center justify-end gap-2">
				<Label for="tls" class="flex items-center gap-0.5">
					Verify Certificate
					<HoverInfo
						text="If your Traefik instance uses a self-signed certificate, you can enable/disable certificate verification here."
					/>
				</Label>
				<Checkbox name="tls" bind:checked={profile.tls} required />
			</div>

			<span class="mt-2 flex flex-row items-center gap-1 border-b border-gray-200 pb-2">
				<span class="font-bold">Basic Authentication</span>
				<HoverInfo
					text="If your Traefik instance requires basic authentication, you can enter your username and password here."
				/>
			</span>

			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="username" class="text-right">Username</Label>
				<Input
					name="username"
					type="text"
					class="col-span-3"
					bind:value={profile.username}
					placeholder="Basic auth username"
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="password" class="text-right">Password</Label>
				<Input
					name="password"
					type="password"
					class="col-span-3"
					bind:value={profile.password}
					placeholder="Basic auth password"
					required
				/>
			</div>
		</div>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => create()}>Create</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
