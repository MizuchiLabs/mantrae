<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import { profileClient } from '$lib/api';

	interface Props {
		profile?: Profile;
		open?: boolean;
	}

	let { profile = $bindable({} as Profile), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			if (profile.id) {
				await profileClient.updateProfile({ name: profile.name, description: profile.description });
				toast.success('Profile updated successfully');
			} else {
				await profileClient.createProfile({ name: profile.name, description: profile.description });
				toast.success('Profile created successfully');
			}
			open = false;
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to save profile', {
				description: e.message
			});
		}
	};

	const handleDelete = async () => {
		if (!profile.id) return;

		try {
			await profileClient.deleteProfile({ id: profile.id });
			toast.success('Profile deleted successfully');
			open = false;
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to delete profile', {
				description: e.message
			});
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>{profile.id ? 'Edit' : 'Add'} Profile</Dialog.Title>
			<Dialog.Description>Configure your Traefik instance connection details.</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="space-y-4">
			<div class="space-y-1">
				<Label for="name">Name</Label>
				<Input id="name" bind:value={profile.name} required placeholder="traefik-site" />
			</div>

			<div class="space-y-1">
				<Label for="description">Description</Label>
				<Input
					id="description"
					bind:value={profile.description}
					placeholder="My Traefik instance"
				/>
			</div>

			<Separator />

			<div class="flex justify-end space-x-2">
				{#if profile.id}
					<Button type="button" variant="destructive" onclick={handleDelete}>Delete</Button>
				{/if}
				<Button type="submit">{profile.id ? 'Update' : 'Create'}</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
