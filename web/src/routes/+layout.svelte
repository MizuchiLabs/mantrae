<script>
	import '../app.css';
	import Profile from '$lib/components/nav/profile.svelte';
	import Sidebar from '$lib/components/nav/sidebar.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import { API_URL, getProfiles, loggedIn, getProviders, getUsers } from '$lib/api';
	import Footer from '$lib/components/nav/footer.svelte';
	import autoAnimate from '@formkit/auto-animate';
	import { onMount } from 'svelte';

	// Realtime updates
	const eventSource = new EventSource(`${API_URL}/events`);
	eventSource.onmessage = (event) => {
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
	});
</script>

<Toaster />

<div class="app flex min-h-screen flex-col">
	{#if $loggedIn}
		<Sidebar />
		<div class="flex flex-1 flex-col sm:py-4 sm:pl-14">
			<main
				class="flex flex-1 flex-col gap-4 p-4 sm:px-6 sm:py-0"
				use:autoAnimate={{ duration: 100 }}
			>
				<div class="mb-6 flex flex-row items-center justify-between">
					<Profile />
				</div>
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
