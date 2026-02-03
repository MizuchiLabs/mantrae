<script lang="ts">
	import { fade } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { user } from '$lib/api/users.svelte';
	import { page } from '$app/state';

	let { children } = $props();

	const redirectRoutes = ['/login'];
	const currentUser = user.self();

	$effect(() => {
		// Only redirect if user is authenticated AND on a public-only route
		if (currentUser.isSuccess && currentUser.data && redirectRoutes.includes(page.url.pathname)) {
			goto('/');
		}
	});
</script>

<main class="flex h-screen w-full items-center justify-center" in:fade={{ duration: 200 }}>
	{@render children?.()}
</main>
