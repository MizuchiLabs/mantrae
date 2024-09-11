<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { onMount } from 'svelte';
	import { deleteUser, getUsers, users } from '$lib/api';
	import CreateUser from '$lib/components/modals/createUser.svelte';
	import UpdateUser from '$lib/components/modals/updateUser.svelte';

	onMount(() => {
		if ($users === undefined) {
			getUsers();
		}
	});
</script>

<CreateUser />

<div class="flex flex-col gap-4 px-4 md:flex-row">
	{#if $users}
		{#each $users as u}
			<Card.Root class="w-full md:w-[400px]">
				<Card.Header>
					<Card.Title class="flex items-center justify-between gap-2">
						<span>{u.username}</span>
						<div class="flex items-center gap-2">
							{#if u.type === 'user'}
								<iconify-icon icon="fa6-solid:user" class="text-green-400" />
							{:else}
								<iconify-icon icon="fa6-solid:robot" class="text-red-400" />
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
					<UpdateUser {u} />
				</Card.Footer>
			</Card.Root>
		{/each}
	{/if}
</div>
