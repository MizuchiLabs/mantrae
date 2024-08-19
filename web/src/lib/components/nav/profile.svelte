<script lang="ts">
	import ChevronsUpDown from 'lucide-svelte/icons/chevrons-up-down';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import CreateProfile from '../modals/createProfile.svelte';
	import UpdateProfile from '../modals/updateProfile.svelte';
	import { profiles, profile } from '$lib/api';

	let open = false;
	function handleProfileClick(name: string) {
		profile.set(name);
		localStorage.setItem('profile', name);
		open = false;
	}
</script>

<Popover.Root bind:open>
	<Popover.Trigger asChild let:builder>
		<Button
			builders={[builder]}
			variant="outline"
			role="combobox"
			aria-expanded={open}
			class="w-[200px] justify-between"
		>
			{$profile || 'Select a profile'}
			<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
		</Button>
	</Popover.Trigger>
	<Popover.Content class="w-[200px] p-0">
		<Command.Root>
			<Command.Input placeholder="Search profile..." />
			<Command.Empty>No profile found.</Command.Empty>
			<Command.Group>
				{#each Object.keys($profiles) ?? [] as name}
					<Command.Item class="flex w-full flex-row items-center justify-between">
						<span class="w-full py-2" on:click={() => handleProfileClick(name)} aria-hidden
							>{name}</span
						>
						<UpdateProfile {name} />
					</Command.Item>
				{/each}
				<Command.Item>
					<CreateProfile />
				</Command.Item>
			</Command.Group>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
