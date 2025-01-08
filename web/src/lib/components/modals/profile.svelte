<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { deleteProfile, upsertProfile } from '$lib/api';
	import type { Profile } from '$lib/types/base';
	import ProfileForm from '../forms/profile.svelte';

	export let profile: Profile;
	export let open = false;

	const update = async () => {
		// Strip trailing slashes
		if (profile.url.endsWith('/')) {
			profile.url = profile.url.slice(0, -1);
		}

		await upsertProfile(profile);
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger />
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Profile</Dialog.Title>
			<Dialog.Description>Update the profile to manage your Traefik instance.</Dialog.Description>
		</Dialog.Header>
		<ProfileForm bind:profile />
		<div class="flex flex-row gap-2">
			{#if profile.id}
				<Button
					type="submit"
					class="w-full bg-red-400"
					on:click={() => {
						deleteProfile(profile);
						open = false;
					}}
				>
					Delete
				</Button>
			{/if}
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		</div>
	</Dialog.Content>
</Dialog.Root>
