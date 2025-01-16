<script lang="ts">
	import { browser } from '$app/environment';
	import { Moon, Sun } from 'lucide-svelte';
	import Button from '../ui/button/button.svelte';

	let darkMode = $state(false);
	function handleSwitchDarkMode() {
		darkMode = !darkMode;
		localStorage.setItem('mode', darkMode ? 'dark' : 'light');

		darkMode
			? document.documentElement.classList.add('dark')
			: document.documentElement.classList.remove('dark');
	}

	if (browser) {
		if (
			localStorage.mode === 'dark' ||
			(!('mode' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)
		) {
			document.documentElement.classList.add('dark');
			darkMode = true;
		} else {
			document.documentElement.classList.remove('dark');
			darkMode = false;
		}
	}
</script>

<Button
	variant="ghost"
	onclick={handleSwitchDarkMode}
	class="mb-2 rounded-full text-gray-600 dark:text-white"
	size="icon"
>
	{#if darkMode}
		<Sun size="1.25rem" />
	{:else}
		<Moon size="1.25rem" />
	{/if}
</Button>
