<script>
	import '../app.css';
	import Profile from '$lib/components/nav/profile.svelte';
	import Sidebar from '$lib/components/nav/sidebar.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Toaster } from '$lib/components/ui/sonner';
	import { onMount } from 'svelte';
	import { getProfiles, profiles, activeProfile, API_URL } from '$lib/api';
	import Footer from '$lib/components/nav/footer.svelte';

	onMount(async () => {
		await getProfiles();
		activeProfile.set($profiles[0] ?? {});
	});
</script>

<Toaster />

<div class="app flex min-h-screen flex-col">
	<Sidebar />
	<div class="flex flex-1 flex-col sm:py-4 sm:pl-14">
		<main class="flex flex-1 flex-col gap-4 p-4 sm:px-6 sm:py-0">
			<div class="mb-6 flex flex-row items-center justify-between">
				<Profile />
				<Button variant="default" href={`${API_URL}/${$activeProfile.name}`}>
					Download Config
					<iconify-icon icon="fa6-solid:download" class="ml-2" />
				</Button>
			</div>
			<slot />
		</main>
		<Footer />
	</div>
</div>
