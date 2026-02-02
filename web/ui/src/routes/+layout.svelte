<script lang="ts">
	import { page } from '$app/state';
	import { queryClient, transport } from '$lib/api/client';
	import logo from '$lib/assets/logo.svg';
	import AppFooter from '$lib/components/nav/AppFooter.svelte';
	import AppHeader from '$lib/components/nav/AppHeader.svelte';
	import AppSidebar from '$lib/components/nav/AppSidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Toaster } from '$lib/components/ui/sonner';
	import { setTransport } from '$lib/query';
	import { profile } from '$lib/stores/profile';
	import { subscribe, unsubscribe } from '$lib/stores/realtime';
	import { user } from '$lib/stores/user';
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import { ModeWatcher } from 'mode-watcher';
	import { onDestroy } from 'svelte';
	import { fade } from 'svelte/transition';
	import './layout.css';

	let { children } = $props();
	setTransport(transport);

	$effect(() => {
		if (user.isLoggedIn() && profile.isValid()) {
			subscribe();
		}
	});
	onDestroy(() => {
		unsubscribe();
	});
</script>

<svelte:head>
	<link rel="icon" href={logo} />
</svelte:head>

<ModeWatcher />
<Toaster />

<QueryClientProvider client={queryClient}>
	{#if user.isLoggedIn()}
		<Sidebar.Provider>
			<AppSidebar variant="inset" />
			<Sidebar.Inset>
				<AppHeader />
				<div class="flex flex-1 flex-col">
					<div class="@container/main flex flex-1 flex-col gap-2">
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
			</Sidebar.Inset>
		</Sidebar.Provider>
	{:else}
		<main class="flex h-screen w-full items-center justify-center" in:fade={{ duration: 200 }}>
			{@render children?.()}
		</main>
	{/if}
</QueryClientProvider>
