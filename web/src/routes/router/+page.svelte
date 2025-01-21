<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import RouterModal from '$lib/components/modals/router.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import type { Router, Service, TLS } from '$lib/types/router';
	import { Edit, Trash } from 'lucide-svelte';
	import { TraefikSource } from '$lib/types';
	import { api, profile, routers, services } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { SOURCE_TAB_SK } from '$lib/store';

	interface RouterModalState {
		isOpen: boolean;
		mode: 'create' | 'edit';
		router?: Router;
		service?: Service;
	}

	const initialModalState: RouterModalState = {
		isOpen: false,
		mode: 'create'
	};

	let modalState = $state(initialModalState);

	function openCreateModal() {
		modalState = {
			isOpen: true,
			mode: 'create'
		};
	}

	function openEditModal(router: Router, service: Service) {
		modalState = {
			isOpen: true,
			mode: 'edit',
			router,
			service
		};
	}

	const deleteRouter = async (name: string | undefined, type: string) => {
		if (!name) return;
		try {
			let routerProvider = name.split('@')[1];
			if (routerProvider !== 'http') {
				toast.error('Router not managed by Mantrae!');
				return;
			}

			await api.deleteRouter($profile.id, name, type);
			toast.success('Router deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

	type RouterWithService = { router: Router; service: Service };
	let mergedData: RouterWithService[] = $state([]);

	const columns: ColumnDef<RouterWithService>[] = [
		{
			header: 'Name',
			accessorFn: (row) => row.router.name,
			id: 'name',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('name') as string;
				return name.split('@')[0];
			}
		},
		{
			header: 'Type',
			accessorFn: (row) => row.router.type,
			id: 'type',
			enableSorting: true,
			cell: ({ row }) => {
				const type = row.getValue('type') as string;
				return renderComponent(ColumnBadge, { label: type });
			}
		},
		{
			header: 'Provider',
			accessorFn: (row) => row.router.name,
			id: 'provider',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('provider') as string;
				return renderComponent(ColumnBadge, {
					label: name.split('@')[1].toLowerCase(),
					variant: 'secondary'
				});
			}
		},
		{
			header: 'Entrypoints',
			accessorFn: (row) => row.router.entryPoints,
			id: 'entryPoints',
			enableSorting: true,
			cell: ({ row }) => {
				const ep = row.getValue('entryPoints') as string[];
				return renderComponent(ColumnBadge, {
					label: ep,
					variant: 'secondary'
				});
			}
		},
		{
			header: 'Middlewares',
			accessorFn: (row) => row.router.middlewares,
			id: 'middlewares',
			enableSorting: true,
			cell: ({ row }) => {
				const middlewares = row.getValue('middlewares') as string[];
				return renderComponent(ColumnBadge, {
					label: middlewares,
					variant: 'secondary'
				});
			}
		},
		{
			header: 'Cert Resolver',
			accessorFn: (row) => row.router.tls,
			id: 'resolver',
			enableSorting: true,
			cell: ({ row }) => {
				const resolver = row.getValue('resolver') as TLS;
				if (!resolver) {
					return renderComponent(ColumnBadge, { label: 'None', variant: 'secondary' });
				}
				return renderComponent(ColumnBadge, {
					label: resolver.certResolver as string,
					variant: 'secondary',
					class: 'bg-slate-300 dark:bg-slate-700'
				});
			}
		},
		{
			header: 'Server Status',
			accessorFn: (row) => row.service.serverStatus,
			id: 'serverStatus',
			enableSorting: true,
			cell: ({ row }) => {
				const status = row.getValue('serverStatus') as Record<string, string>;
				if (!status) return;
				const upCount = Object.values(status).filter((status) => status === 'UP').length;
				const totalCount = Object.values(status).length;
				const greenBadge = 'bg-green-300 dark:bg-green-600';
				const redBadge = 'bg-red-300 dark:bg-red-600';
				return renderComponent(ColumnBadge, {
					label: `${upCount}/${totalCount}`,
					variant: 'secondary',
					class: upCount === totalCount ? greenBadge : redBadge
				});
			}
		},
		{
			id: 'actions',
			cell: ({ row }) => {
				const provider = row.getValue('provider') as string;
				if (provider === 'http') {
					return renderComponent(TableActions, {
						actions: [
							{
								label: 'Edit Router',
								icon: Edit,
								onClick: () => {
									openEditModal(row.original.router, row.original.service);
								}
							},
							{
								label: 'Delete Router',
								icon: Trash,
								variant: 'destructive',
								onClick: () => {
									deleteRouter(row.original.router.name, row.original.router.type);
								}
							}
						]
					});
				} else {
					return renderComponent(TableActions, {
						actions: [
							{
								label: 'Edit Router',
								icon: Edit,
								onClick: () => {
									openEditModal(row.original.router, row.original.service);
								}
							}
						]
					});
				}
			}
		}
	];

	let selectedTab = $state(localStorage.getItem(SOURCE_TAB_SK) || TraefikSource.LOCAL);

	// Update localStorage and fetch config when tab changes
	function handleTabChange(value: string) {
		localStorage.setItem(SOURCE_TAB_SK, value);
		selectedTab = value;
		fetchTraefikConfig();
	}

	// Fetch config based on selected source
	async function fetchTraefikConfig() {
		if (!$profile?.id) return;
		const source = selectedTab === TraefikSource.API ? TraefikSource.API : TraefikSource.LOCAL;
		await api.getTraefikConfig($profile.id, source);
	}

	profile.subscribe(async (value) => {
		if (!value?.id) return;
		await fetchTraefikConfig();
	});
	$effect(() => {
		if ($routers?.length) {
			mergedData = $routers.map((router) => {
				const service = $services.find((service) => service.name === router.name);
				return {
					router,
					service: service || ({} as Service)
				};
			});
		} else {
			mergedData = [];
		}
	});
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<Tabs.Root value={selectedTab} onValueChange={handleTabChange}>
	<Tabs.Content value={TraefikSource.LOCAL}>
		<div class="flex flex-col gap-4">
			<div class="flex flex-col justify-start">
				<h1 class="text-2xl font-bold">Router Management</h1>
				<span class="text-sm text-muted-foreground">Total Routers: {$routers.length}</span>
			</div>
			<DataTable
				{columns}
				data={mergedData || []}
				createButton={{
					label: 'Add Router',
					onClick: openCreateModal
				}}
			/>
		</div>
	</Tabs.Content>
	<Tabs.Content value={TraefikSource.API}>
		<div class="flex flex-col gap-4">
			<div class="flex flex-col justify-start">
				<h1 class="text-2xl font-bold">Router Management</h1>
				<span class="text-sm text-muted-foreground">Total Routers: {$routers.length}</span>
			</div>
			<DataTable
				{columns}
				data={mergedData || []}
				createButton={{
					label: 'Add Router',
					onClick: openCreateModal
				}}
			/>
		</div>
	</Tabs.Content>
</Tabs.Root>

<RouterModal
	bind:open={modalState.isOpen}
	mode={modalState.mode}
	router={modalState.router}
	service={modalState.service}
/>
