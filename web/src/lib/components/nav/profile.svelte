<script lang="ts">
	import ChevronsUpDown from 'lucide-svelte/icons/chevrons-up-down';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import ProfileModal from '../modals/profile.svelte';
	import { api, profiles, profile } from '$lib/api';
	import type { Profile } from '$lib/types';
	import { Pencil, Plus } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { PROFILE_SK } from '$lib/store';

	let editProfile: Partial<Profile>;
	let openPopover = false;
	let openModal = false;

	const profileModal = (p?: Profile) => {
		editProfile = p ?? {};
		openModal = true;
		openPopover = false;
	};

	function changeProfile(p: Profile) {
		localStorage.setItem(PROFILE_SK, p.id.toString());
		profile.set(p);
		openPopover = false;
	}

	onMount(async () => {
		await api.listProfiles();
	});
</script>

<ProfileModal bind:profile={editProfile} bind:open={openModal} />

<Popover.Root bind:open={openPopover}>
	<Popover.Trigger asChild let:builder>
		<Button
			builders={[builder]}
			variant="outline"
			role="combobox"
			aria-expanded={openPopover}
			class="w-[200px] justify-between"
		>
			{$profile?.name ?? 'Select a profile'}
			<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
		</Button>
	</Popover.Trigger>
	<Popover.Content class="w-[200px] p-0">
		<Command.Root>
			<Command.Input placeholder="Search profile..." />
			<Command.Empty>No profile found.</Command.Empty>
			<Command.Group>
				{#each $profiles ?? [] as p}
					<Command.Item
						class="flex w-full flex-row items-center justify-between"
						onSelect={() => changeProfile(p)}
					>
						<span>{p.name}</span>
						<Button
							variant="ghost"
							class="h-8 w-8 rounded-full bg-orange-400"
							size="icon"
							on:click={(event) => {
								event.stopPropagation(); // Prevents the click from bubbling to Command.Item
								profileModal(p);
							}}
						>
							<Pencil size="1rem" />
						</Button>
					</Command.Item>
				{/each}
				<Command.Item
					class="flex w-full flex-row items-center justify-between"
					onSelect={() => profileModal()}
				>
					<span>New Profile</span>
					<Button variant="ghost" class="h-8 w-8 rounded-full bg-green-400" size="icon">
						<Plus size="1rem" />
					</Button>
				</Command.Item>
			</Command.Group>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
