<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { api, loading } from '$lib/api';
	import type { Profile } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import Separator from '../ui/separator/separator.svelte';

	interface Props {
		profile?: Profile;
		open?: boolean;
	}

	let { profile = $bindable({} as Profile), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			// Strip trailing slashes from URL
			if (profile.url?.endsWith('/')) {
				profile.url = profile.url.slice(0, -1);
			}

			if (profile.id) {
				await api.updateProfile(profile as Profile);
				toast.success('Profile updated successfully');
			} else {
				await api.createProfile(profile as Profile);
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
			await api.deleteProfile(profile.id);
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
				<Label for="url">URL</Label>
				<Input id="url" bind:value={profile.url} required placeholder="http://localhost:8080" />
			</div>

			<div class="space-y-1">
				<Label for="username">Username (optional)</Label>
				<Input id="username" bind:value={profile.username} placeholder="admin" autocomplete="off" />
			</div>

			<div class="space-y-1">
				<Label for="password">Password (optional)</Label>
				<PasswordInput bind:value={profile.password} autocomplete="new-password" />
			</div>

			<div class="flex items-center space-x-2">
				<Checkbox id="tls" checked={profile.tls} />
				<Label for="tls">Enable TLS</Label>
			</div>

			<Separator />

			<div class="flex justify-end space-x-2">
				{#if profile.id}
					<Button type="button" variant="destructive" onclick={handleDelete} disabled={$loading}>
						Delete
					</Button>
				{/if}
				<Button type="submit" disabled={$loading}>{profile.id ? 'Update' : 'Create'}</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
