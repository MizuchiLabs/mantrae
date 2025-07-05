<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { toast } from 'svelte-sonner';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import { UpdateUserRequestSchema, type User } from '$lib/gen/mantrae/v1/user_pb';
	import { userClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { user } from '$lib/stores/user';
	import { create } from '@bufbuild/protobuf';

	interface Props {
		data?: User[];
		item: User;
		open?: boolean;
	}

	let { data = $bindable(), item = $bindable(), open = $bindable(false) }: Props = $props();

	let password = $state('');
	let isSelf = $derived(item?.id === user.value?.id);

	const handleSubmit = async () => {
		try {
			if (item.id) {
				let payload = create(UpdateUserRequestSchema);
				payload.id = item.id;
				payload.username = item.username;
				payload.email = item.email;
				payload.isAdmin = item.isAdmin;
				if (item.password && item.password.length > 0) {
					payload.password = item.password;
				}
				await userClient.updateUser(payload);
				toast.success(`User ${item.username} updated successfully.`);
			} else {
				await userClient.createUser({
					username: item.username,
					password: item.password,
					email: item.email,
					isAdmin: item.isAdmin
				});
				toast.success(`User ${item.username} created successfully.`);
			}

			// Refresh data
			if (!data) return;
			let response = await userClient.listUsers({
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.users;
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
			toast.success('EntryPoint deleted successfully');

			// Refresh data
			if (!data) return;
			let response = await userClient.listUsers({
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.users;
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
			{#if isSelf}
				<Dialog.Title>Update Profile</Dialog.Title>
				<Dialog.Description>Configure your profile settings</Dialog.Description>
			{:else}
				<div class="flex flex-row justify-between gap-2">
					<div>
						<Dialog.Title>{item?.id ? 'Update' : 'Add'} User</Dialog.Title>
						<Dialog.Description>Configure your user settings</Dialog.Description>
					</div>
					<div class="mr-4 flex items-center gap-2">
						<Label for="admin">Admin</Label>
						<Switch
							id="admin"
							checked={item.isAdmin || false}
							onCheckedChange={(e) => (item.isAdmin = e)}
						/>
					</div>
				</div>
			{/if}
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="space-y-4">
			<!-- Username -->
			<div class="space-y-1">
				<Label for="name">Name</Label>
				<Input name="name" type="text" bind:value={item.username} placeholder="Name" required />
			</div>

			<!-- Email -->
			<div class="space-y-1">
				<Label for="Email">Email</Label>
				<Input name="Email" type="email" bind:value={item.email} placeholder="Email" />
			</div>

			<!-- Password -->
			<div class="space-y-1">
				{#if item.id}
					<Label for="Password">Password (leave empty to keep current)</Label>
					<PasswordInput bind:value={password} />
				{:else}
					<Label for="Password">Password</Label>
					<PasswordInput bind:value={item.password} required />
				{/if}
			</div>

			<Separator />

			<div class="flex w-full flex-row gap-2">
				{#if item.id}
					<Button type="button" variant="destructive" onclick={handleDelete} class="flex-1">
						Delete
					</Button>
				{/if}
				<Button type="submit" class="flex-1">{item.id ? 'Update' : 'Create'}</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
