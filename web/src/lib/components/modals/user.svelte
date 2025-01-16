<script lang="ts">
	import { upsertUser } from '$lib/api';
	import UserForm from '$lib/components/forms/user.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import type { User } from '$lib/types/base';

	interface Props {
		user: User;
		disabled?: boolean;
		open?: boolean;
	}

	let { user = $bindable(), disabled = false, open = $bindable(false) }: Props = $props();

	let userForm: UserForm = $state();
	const update = async () => {
		const valid = userForm.validate();
		if (!valid) return;

		await upsertUser(user);
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<UserForm bind:user bind:this={userForm} {disabled} />
		<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
