<script lang="ts">
	import ChevronsUpDown from 'lucide-svelte/icons/chevrons-up-down';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import ProfileModal from '../modals/profile.svelte';
	import { profiles, getProfile, profile as currentProfile } from '$lib/api';
	import { newProfile, type Profile } from '$lib/types/base';
	import { Pencil, Plus } from 'lucide-svelte';

	let profile: Profile;
	let open = false;
	let openModal = false;

	const createModal = () => {
		profile = newProfile();
		openModal = true;
		open = false;
	};
	const updateModal = (p: Profile) => {
		profile = p;
		openModal = true;
		open = false;
	};

	function handleProfileClick(id: number) {
		getProfile(id);
		open = false;
	}
</script>

<ProfileModal bind:profile bind:open={openModal} />

<Popover.Root bind:open>
	<Popover.Trigger asChild let:builder>
		<Button
			builders={[builder]}
			variant="outline"
			role="combobox"
			aria-expanded={open}
			class="w-[200px] justify-between"
		>
			{$currentProfile?.name ?? 'Select a profile'}
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
						onSelect={() => handleProfileClick(p.id)}
					>
						<span>{p.name}</span>
						<Button
							variant="ghost"
							class="h-8 w-8 rounded-full bg-orange-400"
							size="icon"
							on:click={(event) => {
								event.stopPropagation(); // Prevents the click from bubbling to Command.Item
								updateModal(p);
							}}
						>
							<Pencil size="1rem" />
						</Button>
					</Command.Item>
				{/each}
				<Command.Item
					class="flex w-full flex-row items-center justify-between"
					onSelect={createModal}
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
