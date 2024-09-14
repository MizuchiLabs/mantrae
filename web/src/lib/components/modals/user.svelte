<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import type { User } from '$lib/types/base';
	import UserForm from '$lib/components/forms/user.svelte';
	import { createUser, updateUser } from '$lib/api';
	import { toast } from 'svelte-sonner';

	export let user: User;
	export let disabled = false;
	export let open = false;
	let pw = '';
	let pwconfirm = '';

	const update = async () => {
		if (user.username === '') {
			toast.error('Username cannot be empty');
			return;
		}

		if (pw !== '') {
			if (pw === pwconfirm) {
				user.password = pw;
			} else {
				toast.error('Passwords do not match');
				return;
			}
		}

		if (user.id) {
			await updateUser(user);
		} else {
			if (user.password === '') {
				toast.error('Password cannot be empty');
				return;
			}
			await createUser(user);
		}
		open = false;
		pw = '';
		pwconfirm = '';
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content>
		<UserForm bind:user bind:pw bind:pwconfirm {disabled} />
		<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
