<script lang="ts">
	import { userClient } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { UpdateUserRequestSchema, type User } from '$lib/gen/mantrae/v1/user_pb';
	import { user } from '$lib/stores/user';
	import { create } from '@bufbuild/protobuf';
	import { ConnectError } from '@connectrpc/connect';
	import { toast } from 'svelte-sonner';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import ConfirmButton from '../ui/confirm-button/confirm-button.svelte';
	import { Trash2 } from '@lucide/svelte';

	interface Props {
		item: User;
		open?: boolean;
	}

	let { item = $bindable(), open = $bindable(false) }: Props = $props();

	let password = $state('');
	let isSelf = $derived(item?.id === user.value?.id);

	const handleSubmit = async () => {
		try {
			if (item.id) {
				let payload = create(UpdateUserRequestSchema);
				payload.id = item.id;
				payload.username = item.username;
				payload.email = item.email;
				if (item.password && item.password.length > 0) {
					payload.password = item.password;
				}
				await userClient.updateUser(payload);
				toast.success(`User ${item.username} updated successfully.`);
			} else {
				await userClient.createUser({
					username: item.username,
					password: item.password,
					email: item.email
				});
				toast.success(`User ${item.username} created successfully.`);
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error(`Failed to update user`, { description: e.message });
		}
		password = '';
		open = false;
	};

	const handleDelete = async () => {
		if (!item.id) return;

		try {
			await userClient.deleteUser({ id: item.id });
			toast.success('User deleted successfully');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete entry point', { description: e.message });
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[425px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>
				{isSelf ? 'Update Profile' : `${item?.id ? 'Edit' : 'Create'} User`}
			</Dialog.Title>
			<Dialog.Description>
				{isSelf
					? 'Update your account information and preferences'
					: 'Configure user account details and permissions'}
			</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="space-y-6">
			<!-- Basic Information -->
			<div class="space-y-4">
				<div class="space-y-2">
					<Label for="username" class="flex items-center gap-2 text-sm font-medium">Username</Label>
					<Input
						id="username"
						bind:value={item.username}
						placeholder="Enter username"
						required
						class="transition-colors"
					/>
					<p class="text-xs text-muted-foreground">Display name for the user account</p>
				</div>

				<div class="space-y-2">
					<Label for="email" class="flex items-center gap-2 text-sm font-medium">Email</Label>
					<Input
						id="email"
						type="email"
						bind:value={item.email}
						placeholder="user@example.com"
						class="transition-colors"
					/>
					<p class="text-xs text-muted-foreground">
						Email address for notifications and account recovery
					</p>
				</div>

				<div class="space-y-2">
					<!-- <Label class="text-sm font-medium">Security</Label> -->
					{#if item.id}
						<Label for="password" class="text-sm font-normal text-muted-foreground">Password</Label>
						<PasswordInput id="password" bind:value={password} />
						<p class="text-xs text-muted-foreground">
							Only enter a new password if you want to change it
						</p>
					{:else}
						<Label for="password" class="text-sm font-normal text-muted-foreground">Password</Label>
						<PasswordInput id="password" bind:value={item.password} required />
						<p class="text-xs text-muted-foreground">Secure password for the user account</p>
					{/if}
				</div>
			</div>

			<!-- {#if !isSelf} -->
			<!-- 	<Separator /> -->

			<!-- Permissions -->
			<!-- 	<div class="space-y-4"> -->
			<!-- 		<div class="space-y-2"> -->
			<!-- 			<Label class="text-sm font-medium">Permissions</Label> -->
			<!-- 			<div class="flex items-center justify-between rounded-lg border p-3"> -->
			<!-- 				<div class="space-y-1"> -->
			<!-- 					<Label for="admin" class="text-sm font-normal">Administrator Access</Label> -->
			<!-- 					<p class="text-muted-foreground text-xs"> -->
			<!-- 						Grant full system access and user management privileges -->
			<!-- 					</p> -->
			<!-- 				</div> -->
			<!-- 				<CustomSwitch bind:checked={item.isAdmin} size="md" /> -->
			<!-- 			</div> -->

			<!-- 			{#if item.isAdmin} -->
			<!-- 				<div class="rounded-lg border border-amber-200 bg-amber-50 p-3"> -->
			<!-- 					<p class="text-xs text-amber-800"> -->
			<!-- 						<strong>Note:</strong> Admin users have full access to all system features and can -->
			<!-- 						manage other users. -->
			<!-- 					</p> -->
			<!-- 				</div> -->
			<!-- 			{/if} -->
			<!-- 		</div> -->
			<!-- 	</div> -->
			<!-- {/if} -->

			<Separator />

			<!-- Actions -->
			<div class="flex gap-2">
				{#if item.id && !isSelf}
					<ConfirmButton
						title="Delete User"
						description="This user will be permanently deleted."
						confirmLabel="Delete"
						cancelLabel="Cancel"
						icon={Trash2}
						class="text-destructive"
						onclick={handleDelete}
					/>
				{/if}
				<Button type="submit" class="flex-1">
					{item.id ? 'Update' : 'Create'}
					{isSelf ? 'Profile' : 'User'}
				</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
