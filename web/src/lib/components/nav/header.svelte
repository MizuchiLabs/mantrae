<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Button } from '$lib/components/ui/button/index';
	import { api, profiles, profile, source } from '$lib/api';
	import Profile from './profile.svelte';
	import InfoModal from '../modals/info.svelte';
	// import Warning from '../modals/warning.svelte';
	import { ArrowLeft, LogOut } from 'lucide-svelte';
	import { SOURCE_TAB_SK } from '$lib/store';
	import { TraefikSource } from '$lib/types';
	import { onMount } from 'svelte';

	// Update localStorage and fetch config when tab changes
	async function handleTabChange(value: TraefikSource) {
		localStorage.setItem(SOURCE_TAB_SK, value);
		source.set(value);
		if (!$profile?.id) return;
		await api.getTraefikConfig($profile.id, $source);
	}

	onMount(async () => {
		let savedSource = localStorage.getItem(SOURCE_TAB_SK) as TraefikSource;
		if (savedSource) {
			source.set(savedSource);
		}
	});
</script>

<nav class="flex h-16 items-center justify-between border-b bg-primary-foreground pl-4">
	<div class="flex flex-row items-center gap-2">
		<Profile />
		{#if !$profiles}
			<span class="ml-2 flex items-center gap-2 text-sm text-muted-foreground">
				<ArrowLeft size="1rem" />
				No profiles configured, create one here
			</span>
		{:else}
			<InfoModal />
		{/if}
		<!-- <Warning /> -->
	</div>

	<Tabs.Root value={$source} onValueChange={(value) => handleTabChange(value as TraefikSource)}>
		<Tabs.List class="grid w-[400px] grid-cols-2">
			<Tabs.Trigger value={TraefikSource.LOCAL}>Local</Tabs.Trigger>
			<Tabs.Trigger value={TraefikSource.API}>API</Tabs.Trigger>
		</Tabs.List>
	</Tabs.Root>

	<div class="mr-2 flex flex-row items-center gap-2">
		<Button variant="ghost" onclick={api.logout} size="icon">
			<LogOut size="1rem" />
		</Button>
	</div>
</nav>
