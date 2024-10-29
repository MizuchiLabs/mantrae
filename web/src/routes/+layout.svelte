<script>
	import '../app.css';
	import Sidebar from '$lib/components/nav/sidebar.svelte';
	import Header from '$lib/components/nav/header.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import { API_URL, getProfiles, loggedIn, getProviders, getUsers, getVersion } from '$lib/api';
	import Footer from '$lib/components/nav/footer.svelte';
	import autoAnimate from '@formkit/auto-animate';
	import { onMount } from 'svelte';
	//import CommandCenter from '$lib/components/utils/commandCenter.svelte';

	// Realtime updates
	const eventSource = new EventSource(`${API_URL}/events`);
	eventSource.onmessage = (event) => {
		if (!$loggedIn) return;
		switch (event.data) {
			case 'profiles':
				getProfiles();
				break;
			case 'users':
				getUsers();
				break;
			case 'providers':
				getProviders();
				break;
		}
	};

	onMount(async () => {
		if (!$loggedIn) return;
		await getProfiles();
		await getProviders();
		await getVersion();
	});
</script>

<Toaster />
<!--<CommandCenter />-->

<div class="app flex min-h-screen flex-col">
	{#if $loggedIn}
		<Sidebar />
		<div class="flex flex-1 flex-col sm:pl-14">
			<Header />
			<main class="flex flex-grow flex-col gap-4 sm:px-2" use:autoAnimate={{ duration: 100 }}>
				<slot />
			</main>
			<Footer />
		</div>
	{:else}
		<div class="flex h-screen flex-col items-center justify-center">
			<slot />
		</div>
	{/if}
</div>
