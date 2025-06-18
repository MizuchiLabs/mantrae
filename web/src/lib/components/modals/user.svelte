<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { toast } from 'svelte-sonner';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import type { User } from '$lib/gen/mantrae/v1/user_pb';
	import { userClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { userId } from '$lib/stores/common';

	interface Props {
		user?: User;
		open?: boolean;
	}

	let { user = $bindable({} as User), open = $bindable(false) }: Props = $props();

	let password = $state('');
	let isSelf = $derived(user?.id === userId.value);

	const handleSubmit = async () => {
		try {
			if (user.id) {
				await userClient.updateUser({
					id: user.id,
					username: user.username,
					email: user.email,
					isAdmin: user.isAdmin,
					password: user.password
				});
				toast.success(`User ${user.username} updated successfully.`);
			} else {
				await userClient.createUser({
					username: user.username,
					password: user.password,
					email: user.email,
					isAdmin: user.isAdmin
				});
				toast.success(`User ${user.username} created successfully.`);
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error(`Failed to update user`, { description: e.message });
		}
		password = '';
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			{#if isSelf}
				<Dialog.Title>Update Profile</Dialog.Title>
			{:else}
				<Dialog.Title>{user?.id ? 'Update' : 'Add'} User</Dialog.Title>
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
					<PasswordInput bind:value={password} />
				{:else}
					<Label for="Password">Password</Label>
					<PasswordInput bind:value={password} required />
				{/if}
			</div>

			<!-- Admin -->
			{#if !isSelf}
				<div class="flex items-center gap-2 space-y-1">
					<Label for="admin">Set Admin</Label>
					<Switch
						id="admin"
						checked={user.isAdmin || false}
						onCheckedChange={(e) => (user.isAdmin = e)}
						class="col-span-3"
					/>
				</div>
			{/if}

			<Separator />

			<Button type="submit" class="w-full">
				{user.id ? 'Update' : 'Save'}
			</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
