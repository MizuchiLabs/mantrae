<script lang="ts">
	import '../app.css';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppSidebar from '$lib/components/nav/AppSidebar.svelte';
	import AppHeader from '$lib/components/nav/AppHeader.svelte';
	import AppFooter from '$lib/components/nav/AppFooter.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import { onMount } from 'svelte';
	import { api, BASE_URL, user } from '$lib/api';
	import autoAnimate from '@formkit/auto-animate';
	// import CommandCenter from '$lib/components/utils/commandCenter.svelte';

	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();

	// Realtime updates
	const eventSource = new EventSource(`${BASE_URL}/events`);
	eventSource.onmessage = (event) => {
		if (!$user) return;
		let data = JSON.parse(event.data);
		switch (data.type) {
			case 'profile_updated':
				api.listProfiles();
				break;
			case 'user_updated':
				api.listUsers();
				break;
			case 'provider_updated':
				api.listDNSProviders();
				break;
			case 'agent_updated':
				api.listAgentsByProfile();
				break;
		}
	};

	onMount(async () => {
		if (!$user) return;
		await api.load();
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
