<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import type { User } from '$lib/types/base';
	import UserForm from '$lib/components/forms/user.svelte';
	import { upsertUser } from '$lib/api';

	export let user: User;
	export let disabled = false;
	export let open = false;

	let userForm: UserForm;
	const update = async () => {
		const valid = userForm.validate();
		if (!valid) return;

		await upsertUser(user);
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content>
		<UserForm bind:user bind:this={userForm} {disabled} />
		<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
