<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { createProfile } from '$lib/api';
	import { newProfile } from '$lib/types/base';

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
			<Button variant="ghost" class="h-8 w-4 rounded-full bg-green-400">
				<iconify-icon icon="fa6-solid:plus" />
			</Button>
		</div>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[450px]">
		<Dialog.Header>
			<Dialog.Title>New profile</Dialog.Title>
			<Dialog.Description>Create a new profile to manage your Traefik clients.</Dialog.Description>
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
					placeholder="URL of your traefik client"
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="username" class="text-right">Username</Label>
				<Input
					name="username"
					type="text"
					class="col-span-3"
					bind:value={profile.username}
					placeholder="Username of your client"
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
					placeholder="Password of your client"
					required
				/>
			</div>
			<div class="flex items-center justify-end gap-4">
				<Label for="tls" class="text-right">Verify Certificate?</Label>
				<Switch name="tls" bind:checked={profile.tls} required />
			</div>
		</div>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => create()}>Add profile</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
