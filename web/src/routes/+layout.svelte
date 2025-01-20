<script lang="ts">
	import '../app.css';
	// import Sidebar from '$lib/components/nav/sidebar.svelte';
	// import Header from '$lib/components/nav/header.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	// import Footer from '$lib/components/nav/footer.svelte';
	import autoAnimate from '@formkit/auto-animate';
	import Header from '$lib/components/nav/header.svelte';
	import Footer from '$lib/components/nav/footer.svelte';
	import Sidebar from '$lib/components/nav/sidebar.svelte';
	import { onMount } from 'svelte';
	import { api, profiles, profile, user } from '$lib/api';
	import { PROFILE_SK } from '$lib/store';
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
		if ($profiles.length === 0) return;

		const savedProfileID = parseInt(localStorage.getItem(PROFILE_SK) ?? '');
		const switchProfile = $profiles.find((p) => p.id === savedProfileID) ?? $profiles[0];
		profile.set(switchProfile);
	});
</script>

<Toaster />
<!-- <CommandCenter /> -->

<div class="app flex min-h-screen flex-col">
	{#if $user}
		<Sidebar />
		<div class="flex w-full flex-1 flex-col sm:pl-16">
			<Header />

			<main class="flex-1 overflow-auto p-6" use:autoAnimate={{ duration: 100 }}>
				{@render children?.()}
			</main>

			<Footer />
		</div>
	{:else}
		<div class="flex h-screen flex-col items-center justify-center">
			{@render children?.()}
		</div>
	{/if}
</div>
