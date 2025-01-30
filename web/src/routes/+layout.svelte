<script lang="ts">
	import '../app.css';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppSidebar from '$lib/components/nav/AppSidebar.svelte';
	import AppHeader from '$lib/components/nav/AppHeader.svelte';
	import AppFooter from '$lib/components/nav/AppFooter.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import { onMount } from 'svelte';
	import { api, BASE_URL } from '$lib/api';
	import autoAnimate from '@formkit/auto-animate';
	import { source } from '$lib/stores/source';
	import { user } from '$lib/stores/user';
	// import CommandCenter from '$lib/components/utils/commandCenter.svelte';

	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();

	// Realtime updates
	interface Event {
		type: string;
		message: string;
	}

	const eventSource = new EventSource(`${BASE_URL}/events`);
	eventSource.onmessage = async (event) => {
		if (!user.isLoggedIn()) return;
		let data: Event = JSON.parse(event.data);
		switch (data.message) {
			case 'profile':
				api.listProfiles();
				break;
			case 'traefik':
				if (!source.isValid(source.value)) return;
				await api.getTraefikConfig(source.value);
				break;
			case 'user':
				api.listUsers();
				break;
			case 'dns':
				api.listDNSProviders();
				break;
			case 'agent':
				api.listAgentsByProfile();
				break;
			default:
				break;
		}
	};

	onMount(async () => {
		await api.load();
	});
</script>

<Toaster />
<!-- <CommandCenter /> -->

<Sidebar.Provider>
	{#if user.isLoggedIn()}
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
