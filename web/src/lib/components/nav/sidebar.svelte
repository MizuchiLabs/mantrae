<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import logo from '$lib/images/logo.svg';
	import Mode from './mode.svelte';
	import { page } from '$app/stores';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Settings, Users, Route, Layers, Globe, Blocks, Bot, Home } from 'lucide-svelte';

	$: active = $page.url.pathname;

	export const routes = [
		{
			name: 'Home',
			path: '/',
			icon: Home
		},
		{
			name: 'Router',
			path: '/router/',
			icon: Route
		},
		{
			name: 'Middlewares',
			path: '/middlewares/',
			icon: Layers
		},
		{
			name: 'Plugins',
			path: '/plugins/',
			icon: Blocks
		},
		{
			name: 'Users',
			path: '/users/',
			icon: Users
		},
		{
			name: 'Agents',
			path: '/agents/',
			icon: Bot
		},
		{
			name: 'DNS',
			path: '/dns/',
			icon: Globe
		},
		{
			name: 'Settings',
			path: '/settings/',
			icon: Settings
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
							class:dark:text-white={active !== route.path}
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
		{/each}
	</div>

	<Mode />
</nav>
