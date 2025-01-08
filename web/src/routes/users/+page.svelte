<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { onMount } from 'svelte';
	import { deleteUser, getUsers, users } from '$lib/api';
	import UserModal from '$lib/components/modals/user.svelte';
	import { newUser, type User } from '$lib/types/base';
	import { UserIcon, Plus, Crown } from 'lucide-svelte';

	let user: User;
	let openModal = false;
	let disabled = false;

	const createModal = () => {
		user = newUser();
		disabled = false;
		openModal = true;
	};

	const updateModal = (u: User) => {
		user = u;
		disabled = true;
		openModal = true;
	};

	onMount(async () => {
		if ($users === undefined) {
			await getUsers();
		}
	});
</script>

<UserModal bind:user bind:open={openModal} {disabled} />

<div class="mt-4 flex flex-col gap-4 px-4 md:flex-row">
	<Button class="flex items-center gap-2 bg-red-400 text-black" on:click={createModal}>
		<span>Add User</span>
		<Plus size="1rem" />
	</Button>
</div>

<div class="flex flex-col gap-4 px-4 md:flex-row">
	{#if $users}
		{#each $users as u}
			<Card.Root class="w-full md:w-[400px]">
				<Card.Header>
					<Card.Title class="flex items-center justify-between gap-2">
						<span>{u.username}</span>
						<div class="flex items-center gap-2">
							{#if u.isAdmin}
								<Crown />
							{:else}
								<UserIcon />
							{/if}
						</div>
					</Card.Title>
				</Card.Header>
				<Card.Content class="space-y-2"></Card.Content>
				<Card.Footer class="flex items-center gap-2">
					{#if u.username !== 'admin'}
						<Button
							variant="ghost"
							class="w-full bg-red-400 text-black"
							on:click={() => deleteUser(u.id)}>Delete</Button
						>
					{/if}
					<Button
						variant="ghost"
						class="w-full bg-orange-400 text-black"
						on:click={() => updateModal(u)}>Edit</Button
					>
				</Card.Footer>
			</Card.Root>
		{/each}
	{/if}
</div>
