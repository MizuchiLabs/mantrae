<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import logo from '$lib/images/logo.svg';
	import Mode from './mode.svelte';
	import { page } from '$app/stores';
	import { slide } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Settings, Users, Route, Layers, Earth } from 'lucide-svelte';

	$: active = $page.url.pathname;

	export const routes = [
		{
			name: 'Router',
			path: '/',
			icon: Route,
			subRoutes: []
		},
		{
			name: 'Middlewares',
			path: '/middlewares/',
			icon: Layers,
			subRoutes: []
		},
		{
			name: 'Settings',
			path: '/settings/',
			icon: Settings,
			subRoutes: [
				{
					name: 'General',
					path: '/settings/',
					icon: Settings
				},
				{
					name: 'Users',
					path: '/settings/users/',
					icon: Users
				},
				{
					name: 'DNS',
					path: '/settings/dns/',
					icon: Earth
				}
			]
		}
	];
</script>

<nav
	class="fixed hidden h-screen w-16 flex-col items-center justify-between border-r bg-primary-foreground sm:flex"
>
	<div class="flex flex-col items-center gap-2 p-4">
		<img src={logo} alt="Mantrae Logo" class="mb-4 w-8" />

		<!-- Base Routes -->
		{#each routes as route}
			<Tooltip.Root openDelay={500}>
				<Tooltip.Trigger>
					<Button
						variant="ghost"
						class="h-12 w-12 rounded-full hover:bg-foreground/5"
						href={route.path}
					>
						<div
							class:text-gray-600={active !== route.path}
							class:text-red-400={active === route.path}
						>
							<svelte:component this={route.icon} />
						</div>
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content side="right" align="center">
					{route.name}
				</Tooltip.Content>
			</Tooltip.Root>

			<!-- Sub Routes -->
			{#if route.subRoutes.length > 0 && active.includes(route.path)}
				<nav
					class="fixed left-16 top-0 hidden h-screen w-48 border-r bg-background text-left sm:flex"
					transition:slide={{ delay: 100, duration: 200, easing: quintOut, axis: 'x' }}
				>
					<div class="flex w-full flex-col items-center justify-start">
						<span class="mb-2 p-4 text-lg font-bold">Settings</span>
						<div class="flex w-full flex-col gap-2 px-2">
							{#each route.subRoutes as subRoute}
								<Button variant="ghost" class="flex h-12 w-full justify-start" href={subRoute.path}>
									<div
										class="flex gap-2 hover:bg-accent/20 hover:text-red-300"
										class:text-gray-600={active !== `${subRoute.path}`}
										class:text-red-400={active === `${subRoute.path}`}
									>
										<svelte:component this={subRoute.icon} />
										{subRoute.name}
									</div>
								</Button>
							{/each}
						</div>
					</div>
				</nav>
			{/if}
		{/each}
	</div>

	<Mode />
</nav>
