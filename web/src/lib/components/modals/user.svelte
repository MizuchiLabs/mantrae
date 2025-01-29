<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import type { User } from '$lib/types';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api, loading, user as currentUser } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import Separator from '../ui/separator/separator.svelte';

	interface Props {
		user: User | undefined;
		open?: boolean;
	}

	let { user = $bindable({} as User), open = $bindable(false) }: Props = $props();

	let password = $state('');
	let isSelf = $derived(user.id === $currentUser?.id);

	const handleSubmit = async () => {
		if (!user.username) return;
		if (user.password) user.password = password;
		if (user.id) {
			await api.updateUser(user);
			toast.success(`User ${user.username} updated successfully`);
		} else {
			await api.createUser(user);
			toast.success(`User ${user.username} created successfully`);
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			{#if isSelf}
				<Dialog.Title>Update Profile</Dialog.Title>
			{:else}
				<Dialog.Title>{user.id ? 'Update' : 'Add'} User</Dialog.Title>
			{/if}
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="space-y-4">
			<!-- Username -->
			<div class="space-y-1">
				<Label for="name">Name</Label>
				<Input name="name" type="text" bind:value={user.username} placeholder="Name" required />
			</div>

			<!-- Email -->
			<div class="space-y-1">
				<Label for="Email">Email</Label>
				<Input name="Email" type="email" bind:value={user.email} placeholder="Email" />
			</div>

			<!-- Password -->
			<div class="space-y-1">
				{#if user.id}
					<Label for="Password">Password (leave empty to keep current)</Label>
					<PasswordInput bind:password />
				{:else}
					<Label for="Password">Password</Label>
					<PasswordInput bind:password required />
				{/if}
			</div>

			<!-- Admin -->
			{#if !isSelf}
				<div class="flex items-center gap-2 space-y-1">
					<Label for="admin">Set Admin</Label>
					<Switch id="admin" checked={user.isAdmin || false} class="col-span-3" />
				</div>
			{/if}

			<Separator />

			<Button type="submit" disabled={$loading} class="w-full">
				{user.id ? 'Update' : 'Save'}
			</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
