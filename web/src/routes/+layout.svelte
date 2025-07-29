<script lang="ts">
	import AppFooter from '$lib/components/nav/AppFooter.svelte';
	import AppHeader from '$lib/components/nav/AppHeader.svelte';
	import AppSidebar from '$lib/components/nav/AppSidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher } from 'mode-watcher';
	import { fade } from 'svelte/transition';
	import { user } from '$lib/stores/user';
	import { page } from '$app/state';
	import '../app.css';

	interface Props {
		children?: import('svelte').Snippet;
	}
	let { children }: Props = $props();
</script>

<ModeWatcher />
<Toaster />

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
