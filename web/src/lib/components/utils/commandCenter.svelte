<script lang="ts">
	import { goto } from '$app/navigation';
	import { getService, routers } from '$lib/api';
	import * as Command from '$lib/components/ui/command';
	import { newRouter, newService, type Router, type Service } from '$lib/types/config';
	import { Layers, Route, Settings } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import RouterModal from '../modals/routerModal.svelte';

	let open = false;
	let searchQuery = '';
	$: filteredRouters = $routers.filter((router: Router) =>
		router.name.toLowerCase().includes(searchQuery.toLowerCase())
	);

	let router: Router;
	let service: Service;
	let disabled = false;
	let openModal = false;

	const createModal = async () => {
		router = newRouter();
		service = newService();
		disabled = false;
		openModal = true;
	};
	const updateModal = async (r: Router) => {
		open = false;
		if (r.provider === 'http') {
			disabled = false;
		} else {
			disabled = true;
		}
		router = r;
		service = getService(router);
		openModal = true;
	};

	onMount(() => {
		function handleKeydown(e: KeyboardEvent) {
			// Check if the focused element is an input, textarea, or contenteditable element
			const focusedElement = document.activeElement;
			const isEditableElement =
				focusedElement?.tagName === 'INPUT' || focusedElement?.tagName === 'TEXTAREA';

			// If focused element is editable, do not run global shortcuts
			if (isEditableElement) {
				return;
			}

			//const currentTime = new Date().getTime(); // Current timestamp
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
					case 'n':
						createModal();
						open = false;
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
				}
			}
		}

		document.addEventListener('keydown', handleKeydown);
		return () => {
			document.removeEventListener('keydown', handleKeydown);
		};
	});
</script>

<div class="hidden">
	<RouterModal {router} {service} {disabled} bind:open={openModal} />
</div>

<Command.Dialog bind:open>
	<Command.Input placeholder="Type a command or search..." bind:value={searchQuery} />
	<Command.List class="max-h-[800px] overflow-y-auto overflow-x-hidden">
		<Command.Empty>No results found.</Command.Empty>
		{#if searchQuery !== ''}
			<Command.Group heading="Routers">
				{#each filteredRouters as router}
					<Command.Item onSelect={() => updateModal(router)}>
						<Route class="mr-2 h-4 w-4" />
						<span>{router.name}</span>
					</Command.Item>
				{/each}
			</Command.Group>
		{/if}
		<Command.Separator />
		<Command.Group heading="Jump to">
			<Command.Item
				onSelect={() => {
					open = false;
					goto('/');
				}}
			>
				<Route class="mr-2 h-4 w-4" />
				<span>Router</span>
				<Command.Shortcut>⌘R</Command.Shortcut>
			</Command.Item>
			<Command.Item
				onSelect={() => {
					open = false;
					goto('/middlewares/');
				}}
			>
				<Layers class="mr-2 h-4 w-4" />
				<span>Middlewares</span>
				<Command.Shortcut>⌘M</Command.Shortcut>
			</Command.Item>
			<Command.Item
				onSelect={() => {
					open = false;
					goto('/settings/');
				}}
			>
				<Settings class="mr-2 h-4 w-4" />
				<span>Settings</span>
				<Command.Shortcut>⌘S</Command.Shortcut>
			</Command.Item>
		</Command.Group>
	</Command.List>
</Command.Dialog>
