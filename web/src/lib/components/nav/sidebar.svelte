<script lang="ts">
	import { page } from '$app/stores';
	import Mode from './mode.svelte';

	$: active = $page.url.pathname;

	export const routes = [
		{
			name: 'Routers',
			path: '/',
			icon: 'fa6-solid:route'
		},
		{
			name: 'Middlewares',
			path: '/middlewares',
			icon: 'fa6-solid:layer-group'
		}
	];
</script>

<nav
	class="fixed hidden h-full w-16 flex-col items-center justify-between border-r bg-background sm:flex"
>
	<div class="flex flex-col">
		{#each routes as route}
			<a
				href={route.path}
				class="p-4 hover:text-red-300"
				class:text-gray-600={active !== `${route.path}`}
				class:text-red-400={active === `${route.path}`}
			>
				<iconify-icon icon={route.icon} class="text-xl" />
				<span class="sr-only">{route.name}</span>
			</a>
		{/each}
	</div>
	<div class="mb-2">
		<Mode />
	</div>
</nav>
