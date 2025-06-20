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
		item: Profile;
		open?: boolean;
	}

	let { item = $bindable(), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			if (item.id) {
				const result = await profileClient.updateProfile({
					name: item.name,
					description: item.description
				});
				toast.success('Profile updated successfully');
				if (result.profile) profileStore.value = result.profile;
			} else {
				const result = await profileClient.createProfile({
					name: item.name,
					description: item.description
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
		if (!item.id) return;

		try {
			await profileClient.deleteProfile({ id: item.id });
			toast.success('Profile deleted successfully');
			if (item.id === profileStore.value?.id) {
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
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[425px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{item?.id ? 'Edit' : 'Create'} Profile</Dialog.Title>
			<Dialog.Description>Configure your profile settings</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="flex flex-col gap-4">
			<div class="flex flex-col gap-2">
				<Label for="name">Name</Label>
				<Input id="name" bind:value={item.name} required placeholder="traefik-site" />
			</div>

			<div class="flex flex-col gap-2">
				<Label for="description">Description</Label>
				<Input id="description" bind:value={item.description} placeholder="Site description" />
			</div>

			<Separator />

			<div class="flex w-full flex-row gap-2">
				{#if item.id}
					<Button type="button" variant="destructive" onclick={handleDelete} class="flex-1">
						Delete
					</Button>
				{/if}
				<Button type="submit" class="flex-1" onclick={handleSubmit}>
					{item.id ? 'Update' : 'Create'}
				</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
