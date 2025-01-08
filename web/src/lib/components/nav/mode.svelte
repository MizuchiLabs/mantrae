<script lang="ts">
	import { browser } from '$app/environment';
	import { darkMode } from '$lib/utils';
	import { Moon, Sun } from 'lucide-svelte';
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

<Button
	variant="ghost"
	on:click={handleSwitchDarkMode}
	class="mb-2 rounded-full text-gray-600 dark:text-white"
	size="icon"
>
	{#if $darkMode}
		<Sun size="1.25rem" />
	{:else}
		<Moon size="1.25rem" />
	{/if}
</Button>
