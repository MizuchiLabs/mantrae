<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import RouterModal from '$lib/components/modals/router.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import type { Router, Service, TLS } from '$lib/types/router';
	import { Pencil, Route, Trash } from 'lucide-svelte';
	import { TraefikSource } from '$lib/types';
	import { api, routerServiceMerge, type RouterWithService } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { source } from '$lib/stores/source';
	import { onMount } from 'svelte';

	interface ModalState {
		isOpen: boolean;
		mode: 'create' | 'edit';
		router: Router;
		service: Service;
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
	const initialModalState: ModalState = {
		isOpen: false,
		mode: 'create',
		router: defaultRouter,
		service: defaultService
	};
	let modalState = $state(initialModalState);

	function openCreateModal() {
		modalState = {
			isOpen: true,
			mode: 'create',
			router: defaultRouter,
			service: defaultService
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

	const deleteRouter = async (router: Router) => {
		try {
			let routerProvider = router.name.split('@')[1];
			if (routerProvider !== 'http') {
				toast.error('Router not managed by Mantrae!');
				return;
			}

			await api.deleteRouter(router);
			toast.success('Router deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

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
			header: 'Protocol',
			accessorFn: (row) => row.router.protocol,
			id: 'protocol',
			enableSorting: true,
			cell: ({ row }) => {
				const protocol = row.getValue('protocol') as string;
				return renderComponent(ColumnBadge, { label: protocol });
			}
		},
		{
			header: 'Provider',
			accessorFn: (row) => row.router.name,
			id: 'provider',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('provider') as string;
				const provider = name.split('@')[1];
				if (!provider && source.value === TraefikSource.AGENT) {
					return renderComponent(ColumnBadge, {
						label: 'agent',
						variant: 'secondary'
					});
				} else if (!provider) {
					return renderComponent(ColumnBadge, {
						label: 'local',
						variant: 'secondary'
					});
				} else {
					return renderComponent(ColumnBadge, {
						label: provider.toLowerCase(),
						variant: 'secondary'
					});
				}
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
				if (!resolver?.certResolver) {
					return renderComponent(ColumnBadge, {
						label: 'None',
						variant: 'secondary',
						class: 'bg-slate-300 dark:bg-slate-700'
					});
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
				if (!status) return renderComponent(ColumnBadge, { label: 'N/A', variant: 'secondary' });
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
			enableHiding: false,
			cell: ({ row }) => {
				if (source.value === TraefikSource.LOCAL) {
					return renderComponent(TableActions, {
						actions: [
							{
								label: 'Edit Router',
								icon: Pencil,
								onClick: () => {
									openEditModal(row.original.router, row.original.service);
								}
							},
							{
								label: 'Delete Router',
								icon: Trash,
								classProps: 'text-destructive',
								onClick: () => {
									deleteRouter(row.original.router);
								}
							}
						]
					});
				} else {
					return renderComponent(TableActions, {
						actions: [
							{
								label: 'Edit Router',
								icon: Pencil,
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

	onMount(() => {
		api.getTraefikConfig(source.value);
	});
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<Tabs.Root value={source.value}>
	<Tabs.Content value={TraefikSource.LOCAL}>
		<div class="flex flex-col gap-4">
			<div class="flex items-center justify-start gap-2">
				<Route />
				<h1 class="text-2xl font-bold">Router Management</h1>
			</div>
			<DataTable
				{columns}
				data={$routerServiceMerge || []}
				showSourceTabs={true}
				createButton={{
					label: 'Add Router',
					onClick: openCreateModal
				}}
			/>
		</div>
	</Tabs.Content>
	<Tabs.Content value={TraefikSource.API}>
		<div class="flex flex-col gap-4">
			<div class="flex items-center justify-start gap-2">
				<Route />
				<h1 class="text-2xl font-bold">Router Management</h1>
			</div>
			<DataTable {columns} data={$routerServiceMerge || []} showSourceTabs={true} />
		</div>
	</Tabs.Content>
	<Tabs.Content value={TraefikSource.AGENT}>
		<div class="flex flex-col gap-4">
			<div class="flex items-center justify-start gap-2">
				<Route />
				<h1 class="text-2xl font-bold">Router Management</h1>
			</div>
			<DataTable {columns} data={$routerServiceMerge || []} showSourceTabs={true} />
		</div>
	</Tabs.Content>
</Tabs.Root>

<RouterModal
	bind:open={modalState.isOpen}
	mode={modalState.mode}
	bind:router={modalState.router}
	bind:service={modalState.service}
/>
