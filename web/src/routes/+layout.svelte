<script lang="ts">
	import AppFooter from '$lib/components/nav/AppFooter.svelte';
	import AppHeader from '$lib/components/nav/AppHeader.svelte';
	import AppSidebar from '$lib/components/nav/AppSidebar.svelte';
	// import AppCenter from '$lib/components/nav/AppCenter.svelte';
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher } from 'mode-watcher';
	import { fade } from 'svelte/transition';
	import '../app.css';
	import { user } from '$lib/stores/user';
	import { onMount } from 'svelte';
	import { profile } from '$lib/stores/profile';
	import { profileClient } from '$lib/api';

	interface Props {
		children?: import('svelte').Snippet;
	}
	let { children }: Props = $props();

	onMount(async () => {
		if (user.isLoggedIn() && !profile.id) {
			const response = await profileClient.listProfiles({});
			profile.value = response.profiles[0];
		}
	});
</script>

<ModeWatcher />
<Toaster />
<!-- <AppCenter /> -->

<Sidebar.Provider>
	{#if user.isLoggedIn()}
		<div class="bg-background flex h-screen w-full">
			<AppSidebar />
			<div class="flex w-full flex-1 flex-col">
				<AppHeader />

				<main class="flex-1 overflow-auto p-8 px-8">
					{#key page.url.pathname}
						<div in:fade={{ duration: 200 }}>
							{@render children?.()}
						</div>
					{/key}
				</main>

				<AppFooter />
			</div>
		</div>
	{:else}
		<main class="flex h-screen w-full items-center justify-center" in:fade={{ duration: 200 }}>
			{@render children?.()}
		</main>
	{/if}
</Sidebar.Provider>
