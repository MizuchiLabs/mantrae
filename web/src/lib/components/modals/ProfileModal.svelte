<script lang="ts">
	import { buildConnectionString, profileClient } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import { profile as profileStore } from '$lib/stores/profile';
	import { ConnectError } from '@connectrpc/connect';
	import { RotateCcw, Trash2 } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import ConfirmButton from '../ui/confirm-button/confirm-button.svelte';
	import { CopyInput } from '../ui/input-group';

	interface Props {
		item: Profile;
		open?: boolean;
	}

	let { item = $bindable(), open = $bindable(false) }: Props = $props();

	const onsubmit = async () => {
		try {
			if (item.id) {
				const response = await profileClient.updateProfile({
					id: item.id,
					name: item.name,
					description: item.description
				});
				toast.success(`Profile ${response.profile?.name} updated successfully`);
				if (response.profile) profileStore.value = response.profile;
			} else {
				const response = await profileClient.createProfile({
					name: item.name,
					description: item.description
				});
				toast.success(`Profile ${response.profile?.name} created successfully`);
				if (response.profile) profileStore.value = response.profile;
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to save profile', { description: e.message });
		}
		open = false;
	};

	const deleteProfile = async () => {
		if (!item.id) return;

		try {
			await profileClient.deleteProfile({ id: item.id });
			toast.success('Profile deleted successfully');
			if (item.id === profileStore.value?.id) {
				let response = await profileClient.listProfiles({});
				if (response.profiles.length === 0) {
					profileStore.value = {} as Profile;
					return;
				}
				profileStore.value = response.profiles[0];
			} else {
				profileStore.value = {} as Profile;
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete profile', { description: e.message });
		}
		open = false;
	};

	const regenerate = async () => {
		try {
			const response = await profileClient.updateProfile({
				id: item.id,
				name: item.name,
				description: item.description,
				regenerateToken: true
			});
			if (!response.profile) throw new Error('Failed to regenerate token');
			toast.success(`Token regenerated successfully`);
			item = response.profile;
			profileStore.value = response.profile;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to regenerate token', { description: e.message });
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-100 overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{item?.id ? 'Edit' : 'Create'} Profile</Dialog.Title>
			<Dialog.Description>Configure your profile settings</Dialog.Description>
		</Dialog.Header>

		<form {onsubmit} class="space-y-4">
			<div class="space-y-2">
				<Label for="name" class="text-sm font-medium">Name</Label>
				<Input
					id="name"
					bind:value={item.name}
					placeholder="traefik-site"
					class="transition-colors"
				/>
				<p class="text-xs text-muted-foreground">A descriptive name for this profile</p>
			</div>

			<div class="space-y-2">
				<Label for="description" class="text-sm font-medium">Description</Label>
				<Input id="description" bind:value={item.description} placeholder="Site description" />
				<p class="text-xs text-muted-foreground">Optional description for this profile</p>
			</div>

			{#if item.id}
				<div class="space-y-2">
					<Label for="token" class="text-sm font-medium">Connection Token</Label>
					<div class="flex gap-2">
						<CopyInput value={item.token} readonly />
						<Button variant="outline" size="icon" onclick={regenerate} title="Regenerate token">
							<RotateCcw class="h-4 w-4" />
						</Button>
					</div>
					{#await buildConnectionString(item) then value}
						<!-- <CopyInput id="token" {value} readonly /> -->
						<p class="text-xs text-muted-foreground">
							Used in the connection URL to connect to this profile with your traefik instance
							<span class="underline">
								{value}
							</span>
						</p>
					{/await}
				</div>
			{/if}

			<Separator />

			<div class="flex w-full flex-row gap-2">
				{#if item.id}
					<ConfirmButton
						title="Delete Profile"
						description="This profile and all associated data will be permanently deleted."
						confirmLabel="Delete"
						cancelLabel="Cancel"
						icon={Trash2}
						class="text-destructive"
						onclick={deleteProfile}
					/>
				{/if}
				<Button type="submit" class="flex-1">{item.id ? 'Update' : 'Create'}</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
