<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { deleteProfile, updateProfile } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import type { Profile } from '$lib/types/dynamic';

	export let profile: Profile;
	let oldName = profile.name;
	let open = false;
	$: console.log(oldName);

	const update = async () => {
		try {
			await updateProfile(oldName, profile);
			toast.success('Success', {
				description: `Profile ${profile.name} has been updated`
			});
			oldName = profile.name;
		} catch (err: any) {
			toast.error('Failed to update profile', {
				description: err
			});
		}
		open = false;
	};

	const onKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			update();
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger class="flex flex-row justify-end">
		<Button variant="ghost" class="h-8 w-4 rounded-full bg-orange-400">
			<iconify-icon icon="fa6-solid:pencil" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Update profile</Dialog.Title>
		</Dialog.Header>
		<div class="grid gap-4 py-4" on:keydown={onKeydown} aria-hidden>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="name" class="text-right">Name</Label>
				<Input
					name="name"
					type="text"
					class="col-span-3"
					placeholder="Your profile name"
					bind:value={profile.name}
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="url" class="text-right">URL</Label>
				<Input
					name="url"
					type="text"
					class="col-span-3"
					bind:value={profile.instance.url}
					placeholder="URL of your instance"
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="username" class="text-right">Username</Label>
				<Input
					name="username"
					type="text"
					class="col-span-3"
					bind:value={profile.instance.username}
					placeholder="Username of your instance"
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="password" class="text-right">Password</Label>
				<Input
					name="password"
					type="password"
					class="col-span-3"
					bind:value={profile.instance.password}
					placeholder="Password of your instance"
					required
				/>
			</div>
		</div>
		<Dialog.Close class="flex w-full flex-row gap-2">
			<Button type="submit" class="w-full bg-red-400" on:click={() => deleteProfile(profile.name)}>
				Delete
			</Button>
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
