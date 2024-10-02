<script lang="ts">
	import { browser } from '$app/environment';
	import { darkMode } from '$lib/utils';
	import Button from '../ui/button/button.svelte';

	function handleSwitchDarkMode() {
		darkMode.update((d) => !d);
		localStorage.setItem('mode', $darkMode ? 'dark' : 'light');

		$darkMode
			? document.documentElement.classList.add('dark')
			: document.documentElement.classList.remove('dark');
	}

	if (browser) {
		if (
			localStorage.mode === 'dark' ||
			(!('mode' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)
		) {
			document.documentElement.classList.add('dark');
			darkMode.set(true);
		} else {
			document.documentElement.classList.remove('dark');
			darkMode.set(false);
		}
	}
</script>

<Button variant="ghost" on:click={handleSwitchDarkMode} class="mb-2 h-12 w-12 rounded-full">
	{#if darkMode}
		<iconify-icon icon="line-md:sunny-outline-to-moon-loop-transition" width="20" height="20" />
	{:else}
		<iconify-icon icon="line-md:moon-alt-to-sunny-outline-loop-transition" width="20" height="20" />
	{/if}
</Button>
