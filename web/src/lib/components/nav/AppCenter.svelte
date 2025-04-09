<script lang="ts">
	import * as Command from '$lib/components/ui/command/index.js';
	import RouterModal from '../modals/router.svelte';
	import MiddlewareModal from '../modals/middleware.svelte';
	import { Layers, Route } from 'lucide-svelte';
	import type { Router, Service } from '$lib/types/router';
	import type { Middleware } from '$lib/types/middlewares';
	import { middlewares, routerServiceMerge } from '$lib/api';
	import { source } from '$lib/stores/source';
	import { TraefikSource } from '$lib/types';

	let open = $state(false);
	let searchQuery = $state('');

	interface ModalState {
		isOpenRouter: boolean;
		isOpenMiddleware: boolean;
		mode: 'create' | 'edit';
		router: Router;
		service: Service;
		middleware: Middleware;
	}

	const defaultRouter: Router = {
		name: '',
		protocol: 'http',
		tls: {},
		entryPoints: [],
		middlewares: [],
		rule: '',
		service: ''
	};
	const defaultService: Service = {
		name: defaultRouter.name,
		protocol: defaultRouter.protocol,
		loadBalancer: {
			servers: [],
			passHostHeader: true
		}
	};
	const defaultMiddleware: Middleware = {
		name: '',
		protocol: 'http',
		type: undefined
	};
	const initialModalState: ModalState = {
		isOpenRouter: false,
		isOpenMiddleware: false,
		mode: 'create',
		router: defaultRouter,
		service: defaultService,
		middleware: defaultMiddleware
	};
	let modalState = $state(initialModalState);

	const baseModal = () => {
		return (modalState = {
			isOpenRouter: false,
			isOpenMiddleware: false,
			mode: 'create',
			router: defaultRouter,
			service: defaultService,
			middleware: defaultMiddleware
		});
	};
	const createRouter = () => {
		source.value = TraefikSource.LOCAL;
		modalState = baseModal();
		modalState.mode = 'create';
		modalState.isOpenRouter = true;
		open = false;
	};
	const createMiddleware = () => {
		source.value = TraefikSource.LOCAL;
		modalState = baseModal();
		modalState.mode = 'create';
		modalState.isOpenMiddleware = true;
		open = false;
	};
	const updateRouter = (router: Router, service: Service) => {
		modalState = baseModal();
		modalState.mode = 'edit';
		modalState.isOpenRouter = true;
		modalState.router = router;
		modalState.service = service;
		open = false;
	};
	const updateMiddleware = (middleware: Middleware) => {
		source.value = TraefikSource.LOCAL;
		modalState = baseModal();
		modalState.mode = 'edit';
		modalState.isOpenMiddleware = true;
		modalState.middleware = middleware;
		open = false;
	};

	// Keyboard shortcuts
	const handleKeydown = (e: KeyboardEvent) => {
		const isEditableElement =
			document.activeElement?.tagName === 'INPUT' || document.activeElement?.tagName === 'TEXTAREA';

		if (isEditableElement) return;

		// Command palette toggle
		if (e.key === '/' || ((e.ctrlKey || e.metaKey) && e.key === 'k')) {
			e.preventDefault();
			open = !open;
			return;
		}
	};

	$effect(() => {
		document.addEventListener('keydown', handleKeydown);
		return () => document.removeEventListener('keydown', handleKeydown);
	});
</script>

<RouterModal
	mode={modalState.mode}
	bind:router={modalState.router}
	bind:service={modalState.service}
	bind:open={modalState.isOpenRouter}
/>
<MiddlewareModal
	mode={modalState.mode}
	bind:middleware={modalState.middleware}
	bind:open={modalState.isOpenMiddleware}
/>

<Command.Dialog bind:open>
	<Command.Root loop>
		<Command.Input placeholder="Search..." bind:value={searchQuery} />
		<Command.List class="max-h-[800px] overflow-x-hidden overflow-y-auto">
			{#if searchQuery !== ''}
				<Command.Group heading="Routers">
					<Command.Empty>No results found.</Command.Empty>
					{#each $routerServiceMerge || [] as m}
						<Command.Item onSelect={() => updateRouter(m.router, m.service)} value={m.router.name}>
							<Route class="mr-2 h-4 w-4" />
							<span>{m.router.name}</span>
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Group heading="Middlewares">
					<Command.Empty>No results found.</Command.Empty>
					{#each $middlewares || [] as m}
						<Command.Item onSelect={() => updateMiddleware(m)} value={m.name}>
							<Layers class="mr-2 h-4 w-4" />
							<span>{m.name}</span>
						</Command.Item>
					{/each}
				</Command.Group>
			{/if}
			<Command.Separator />
			<Command.Group heading="Create">
				<Command.Item onSelect={createRouter}>
					<Route class="mr-2 h-4 w-4" />
					<span>Create Router</span>
				</Command.Item>
				<Command.Item onSelect={createMiddleware}>
					<Layers class="mr-2 h-4 w-4" />
					<span>Create Middleware</span>
				</Command.Item>
			</Command.Group>
		</Command.List>
	</Command.Root>
</Command.Dialog>
