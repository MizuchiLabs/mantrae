<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type User } from '$lib/gen/mantrae/v1/user_pb';
	import PasswordInput from '../ui/password-input/password-input.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import { user } from '$lib/api/users.svelte';

	interface Props {
		data?: User;
		open?: boolean;
	}
	let { data, open = $bindable(false) }: Props = $props();

	let password = $state('');
	let userData = $state({} as User);
	$effect(() => {
		if (data) userData = { ...data };
	});
	$effect(() => {
		if (!open) {
			userData = {} as User;
			password = '';
		}
	});

	const createMutation = user.create();
	const updateMutation = user.update();
	function onsubmit() {
		if (password !== '') {
			userData.password = password;
		}
		if (userData.id) {
			updateMutation.mutate({ ...userData });
		} else {
			createMutation.mutate({ ...userData });
		}
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[425px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>
				{userData?.id ? 'Edit' : 'Create'} User
			</Dialog.Title>
			<Dialog.Description>Configure user account details and permissions</Dialog.Description>
		</Dialog.Header>

		<form {onsubmit} class="space-y-6">
			<div class="space-y-4">
				<div class="space-y-2">
					<Label for="username" class="flex items-center gap-2 text-sm font-medium">Username</Label>
					<Input
						id="username"
						bind:value={userData.username}
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
						bind:value={userData.email}
						placeholder="user@example.com"
						class="transition-colors"
					/>
					<p class="text-xs text-muted-foreground">
						Email address for notifications and account recovery
					</p>
				</div>

				<div class="space-y-2">
					{#if userData.id}
						<Label for="password" class="text-sm font-normal text-muted-foreground">Password</Label>
						<PasswordInput id="password" bind:value={password} />
						<p class="text-xs text-muted-foreground">
							Only enter a new password if you want to change it
						</p>
					{:else}
						<Label for="password" class="text-sm font-normal text-muted-foreground">Password</Label>
						<PasswordInput id="password" bind:value={password} required />
						<p class="text-xs text-muted-foreground">Secure password for the user account</p>
					{/if}
				</div>
			</div>

			<Separator />

			<Button type="submit" class="w-full">{userData.id ? 'Update' : 'Create'}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
