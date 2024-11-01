<script lang="ts">
	import * as Command from '$lib/components/ui/command';
	import { goto } from '$app/navigation';
	import { getMiddlewares, getRouters, middlewares, profile, routers, services } from '$lib/api';
	import { newRouter, newService, type Router, type Service } from '$lib/types/config';
	import { onMount } from 'svelte';
	import RouterModal from '../modals/router.svelte';
	import MiddlewareModal from '../modals/middleware.svelte';
	import { newMiddleware, type Middleware } from '$lib/types/middlewares';
	import { Earth, Layers, Route, Settings, Users } from 'lucide-svelte';

	let open = false;
	let searchQuery = '';
	let timeout: NodeJS.Timeout;
	let fRouters: Router[] = [];
	let fMiddlewares: Middleware[] = [];

	// Debounced search function
	function debounceSearch() {
		clearTimeout(timeout);
		timeout = setTimeout(() => {
			// Ensure stores are available before filtering
			fRouters = $routers?.filter((router) =>
				router.name.toLowerCase().includes(searchQuery.toLowerCase())
			);
			fMiddlewares = $middlewares?.filter((middleware) =>
				middleware.name.toLowerCase().includes(searchQuery.toLowerCase())
			);
		}, 300);
	}

	$: debounceSearch(); // Trigger debounce on every searchQuery change

	let router: Router;
	let service: Service;
	let middleware: Middleware;
	let disabled = false;
	let openRouterModal = false;
	let openMiddlewareModal = false;

	const createRouter = async () => {
		open = false;
		router = newRouter();
		service = newService();
		disabled = false;
		openRouterModal = true;
	};
	const updateRouter = async (r: Router) => {
		open = false;
		if (r.provider === 'http') {
			disabled = false;
		} else {
			disabled = true;
		}
		router = r;
		service = $services?.find((s: Service) => s.name === r.name) ?? newService();
		openRouterModal = true;
	};

	const createMiddleware = async () => {
		open = false;
		middleware = newMiddleware();
		disabled = false;
		openMiddlewareModal = true;
	};
	const updateMiddleware = async (m: Middleware) => {
		open = false;
		middleware = m;
		disabled = false;
		openMiddlewareModal = true;
	};

	const routes = [
		{ name: 'Routers', path: '/', icon: Route },
		{ name: 'Middlewares', path: '/middlewares/', icon: Layers },
		{ name: 'Settings', path: '/settings/', icon: Settings },
		{ name: 'Users', path: '/settings/users/', icon: Users },
		{ name: 'DNS', path: '/settings/dns/', icon: Earth }
	];

	onMount(() => {
		function handleKeydown(e: KeyboardEvent) {
			// Check if the focused element is an input or textarea
			const focusedElement = document.activeElement;
			const isEditableElement =
				focusedElement?.tagName === 'INPUT' || focusedElement?.tagName === 'TEXTAREA';

			// If focused element is editable, do not run global shortcuts
			if (isEditableElement) {
				return;
			}

			if (e.key === '/') {
				open = !open;
				e.preventDefault();
			}

			if (e.metaKey || e.ctrlKey) {
				e.preventDefault();

				switch (e.key) {
					case 'k':
						open = !open;
						break;
					case 'r':
						goto('/');
						open = false;
						break;
					case 'm':
						goto('/middlewares/');
						open = false;
						break;
					case 's':
						goto('/settings/');
						open = false;
						break;
					case 'u':
						goto('/settings/users/');
						open = false;
						break;
					case 'd':
						goto('/settings/dns/');
						open = false;
						break;
				}
			}
		}

		document.addEventListener('keydown', handleKeydown);
		return () => {
			document.removeEventListener('keydown', handleKeydown);
		};
	});

	// Get routers when the profile changes
	profile.subscribe((value) => {
		if (!value?.id) return;
		getRouters(value.id);
		getMiddlewares(value.id);
	});
</script>

<div class="hidden">
	<RouterModal {router} {service} {disabled} bind:open={openRouterModal} />
	<MiddlewareModal {middleware} {disabled} bind:open={openMiddlewareModal} />
</div>

<Command.Dialog bind:open>
	<Command.Root loop>
		<Command.Input placeholder="Type a command or search..." bind:value={searchQuery} />
		<Command.List class="max-h-[800px] overflow-y-auto overflow-x-hidden">
			{#if searchQuery !== ''}
				<Command.Group heading="Routers">
					<Command.Empty>No results found.</Command.Empty>
					{#each fRouters || [] as router}
						<Command.Item onSelect={() => updateRouter(router)} value={router.name}>
							<Route class="mr-2 h-4 w-4" />
							<span>{router.name}</span>
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Group heading="Middlewares">
					<Command.Empty>No results found.</Command.Empty>
					{#each fMiddlewares || [] as middleware}
						<Command.Item onSelect={() => updateMiddleware(middleware)} value={middleware.name}>
							<Layers class="mr-2 h-4 w-4" />
							<span>{middleware.name}</span>
						</Command.Item>
					{/each}
				</Command.Group>
			{/if}
			<Command.Separator />
			<Command.Group heading="Create">
				<Command.Item onSelect={() => createRouter()}>
					<Route class="mr-2 h-4 w-4" />
					<span>Create Router</span>
				</Command.Item>
				<Command.Item onSelect={() => createMiddleware()}>
					<Layers class="mr-2 h-4 w-4" />
					<span>Create Middleware</span>
				</Command.Item>
			</Command.Group>
			<Command.Separator />
			<Command.Group heading="Jump to">
				<Command.Empty>No results found.</Command.Empty>
				{#each routes as route}
					<Command.Item
						onSelect={() => {
							open = false;
							goto(route.path);
						}}
					>
						<svelte:component this={route.icon} class="mr-2 h-4 w-4" />
						<span>{route.name}</span>
						<Command.Shortcut>⌘{route.name[0]}</Command.Shortcut>
					</Command.Item>
				{/each}
			</Command.Group>
		</Command.List>
	</Command.Root>
</Command.Dialog>
