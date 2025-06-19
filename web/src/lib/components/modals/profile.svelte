<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import { profileClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { profile as profileStore } from '$lib/stores/profile';

	interface Props {
		profile: Profile;
		open?: boolean;
	}

	let { profile = $bindable(), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			if (profile.id) {
				const result = await profileClient.updateProfile({
					name: profile.name,
					description: profile.description
				});
				toast.success('Profile updated successfully');
				if (result.profile) profileStore.value = result.profile;
			} else {
				const result = await profileClient.createProfile({
					name: profile.name,
					description: profile.description
				});
				toast.success('Profile created successfully');
				if (result.profile) profileStore.value = result.profile;
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to save profile', { description: e.message });
		}
		open = false;
	};

	const handleDelete = async () => {
		if (!profile.id) return;

		try {
			await profileClient.deleteProfile({ id: profile.id });
			toast.success('Profile deleted successfully');
			if (profile.id === profileStore.value?.id) {
				let profiles = await profileClient.listProfiles({ limit: -1n, offset: 0n });
				profileStore.value = profiles.profiles[0];
			} else {
				profileStore.value = {} as Profile;
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete profile', { description: e.message });
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>{profile.id ? 'Edit' : 'Create'} Profile</Dialog.Title>
			<Dialog.Description>Configure your profile settings</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="flex flex-col gap-4">
			<div>
				<Label for="name">Name</Label>
				<Input id="name" bind:value={profile.name} required placeholder="traefik-site" />
			</div>

			<div>
				<Label for="description">Description</Label>
				<Input id="description" bind:value={profile.description} placeholder="Site description" />
			</div>

			<Separator />

			<div class="flex w-full flex-row-reverse gap-2">
				<Button type="submit" class="w-full">{profile.id ? 'Update' : 'Create'}</Button>
				{#if profile.id}
					<Button type="button" variant="destructive" class="w-full" onclick={handleDelete}>
						Delete
					</Button>
				{/if}
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
