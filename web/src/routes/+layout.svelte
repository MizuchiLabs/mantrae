<script lang="ts">
	import '../app.css';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppSidebar from '$lib/components/nav/AppSidebar.svelte';
	import AppHeader from '$lib/components/nav/AppHeader.svelte';
	import AppFooter from '$lib/components/nav/AppFooter.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import { onMount } from 'svelte';
	import { PROFILE_SK } from '$lib/store';
	import { api, profiles, profile, user } from '$lib/api';
	import autoAnimate from '@formkit/auto-animate';

	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
	// import type { User } from '$lib/types';
	// import CommandCenter from '$lib/components/utils/commandCenter.svelte';

	// Realtime updates
	// const eventSource = new EventSource(`${API_URL}/events`);
	// eventSource.onmessage = (event) => {
	// 	if (!$loggedIn) return;
	// 	let data = JSON.parse(event.data);
	// 	switch (data.type) {
	// 		case 'profile_updated':
	// 			getProfiles();
	// 			break;
	// 		case 'router_updated':
	// 			getRouters();
	// 			break;
	// 		case 'user_updated':
	// 			getUsers();
	// 			break;
	// 		case 'provider_updated':
	// 			getProviders();
	// 			break;
	// 		case 'agent_updated':
	// 			getAgents(data.message.profileId);
	// 			break;
	// 		case 'config_error':
	// 			configError.set(data.message);
	// 			break;
	// 		case 'config_ok':
	// 			configError.set('');
	// 			break;
	// 	}
	// };

	onMount(async () => {
		if (!$user) return;
		await api.listProfiles();
		if (!$profiles) return;

		const savedProfileID = parseInt(localStorage.getItem(PROFILE_SK) ?? '');
		const switchProfile = $profiles.find((p) => p.id === savedProfileID) ?? $profiles[0];
		profile.set(switchProfile);
	});
</script>

<Toaster />
<!-- <CommandCenter /> -->

<Sidebar.Provider>
	{#if $user}
		<div class="flex h-screen w-full bg-background">
			<AppSidebar />
			<div class="flex w-full flex-1 flex-col">
				<AppHeader />

				<main class="flex-1 overflow-auto p-8 px-8" use:autoAnimate={{ duration: 100 }}>
					{@render children?.()}
				</main>

				<AppFooter />
			</div>
		</div>
	{:else}
		<main class="flex h-screen w-full items-center justify-center">
			{@render children?.()}
		</main>
	{/if}
</Sidebar.Provider>
