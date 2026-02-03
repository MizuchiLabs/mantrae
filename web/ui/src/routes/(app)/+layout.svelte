<script lang="ts">
	import { page } from '$app/state';
	import logo from '$lib/assets/logo.svg';
	import AppFooter from '$lib/components/nav/AppFooter.svelte';
	import AppHeader from '$lib/components/nav/AppHeader.svelte';
	import AppSidebar from '$lib/components/nav/AppSidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher } from 'mode-watcher';
	import { fade } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { profile } from '$lib/api/profiles.svelte';
	import { user } from '$lib/api/users.svelte';

	let { children } = $props();

	const profiles = profile.list();
	const currentUser = user.self();

	$effect(() => {
		if (currentUser.isError) goto('/login');
	});
	$effect(() => {
		if (currentUser.isSuccess && currentUser.data && !profiles.data) {
			profiles.refetch();
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={logo} />
</svelte:head>

<ModeWatcher />
<Toaster />

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
