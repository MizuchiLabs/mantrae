<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import logo from '$lib/images/logo.svg';
	import Mode from './mode.svelte';
	import { page } from '$app/stores';
	import { slide } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';

	$: active = $page.url.pathname;

	export const routes = [
		{
			name: 'Routers',
			path: '/',
			icon: 'fa6-solid:route'
		},
		{
			name: 'Middlewares',
			path: '/middlewares/',
			icon: 'fa6-solid:layer-group'
		},
		{
			name: 'Settings',
			path: '/settings/',
			icon: 'fa6-solid:gear'
		}
	];

	export const settingsRoutes = [
		{
			name: 'General',
			path: '/settings/',
			icon: 'fa6-solid:gear'
		},
		{
			name: 'Users',
			path: '/settings/users/',
			icon: 'fa6-solid:user-group'
		},
		{
			name: 'DNS',
			path: '/settings/dns/',
			icon: 'fa6-solid:earth-americas'
		}
	];
</script>

<nav
	class="fixed hidden h-screen w-16 flex-col items-center justify-between border-r bg-primary-foreground sm:flex"
>
	<div class="flex flex-col items-center gap-6 p-4">
		<img src={logo} alt="Mantrae Logo" class="mb-4 w-8" />

		{#each routes as route}
			<a
				href={route.path}
				class="hover:text-red-300"
				class:text-gray-600={active !== `${route.path}`}
				class:text-red-400={active === `${route.path}`}
			>
				<iconify-icon icon={route.icon} class="text-xl" />
				<span class="sr-only">{route.name}</span>
			</a>
		{/each}
	</div>

	<Mode />
</nav>

<!-- Secondary Sidebar for Settings -->
{#if active.includes('/settings/')}
	<nav
		class="fixed left-16 hidden h-screen w-48 border-r bg-background text-left sm:flex"
		transition:slide={{ delay: 100, duration: 200, easing: quintOut, axis: 'x' }}
	>
		<div class="flex w-full flex-col items-center justify-start">
			<span class="mb-4 p-4 text-lg font-bold">Settings</span>
			<div class="flex w-full flex-col gap-2 px-2">
				{#each settingsRoutes as route}
					<Button variant="ghost" class="flex w-full justify-start" href={route.path}>
						<div
							class="flex gap-2 hover:bg-accent/20 hover:text-red-300"
							class:text-gray-600={active !== `${route.path}`}
							class:text-red-400={active === `${route.path}`}
						>
							<iconify-icon icon={route.icon} width="18" />
							{route.name}
						</div>
					</Button>
				{/each}
			</div>
		</div>
	</nav>
{/if}
