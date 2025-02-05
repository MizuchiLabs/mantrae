<script lang="ts">
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
	import ColumnRule from '$lib/components/tables/ColumnRule.svelte';

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
			await api.deleteRouter(router);
			toast.success('Router deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

	const defaultColumns: ColumnDef<RouterWithService>[] = [
		{
			header: 'Name',
			accessorKey: 'router.name',
			id: 'name',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('name') as string;
				return name?.split('@')[0];
			}
		},
		{
			header: 'Protocol',
			accessorKey: 'router.protocol',
			id: 'protocol',
			enableSorting: true,
			cell: ({ row }) => {
				return renderComponent(ColumnBadge, {
					label: row.getValue('protocol') as string
				});
			}
		},
		{
			header: 'Provider',
			accessorKey: 'router.name',
			id: 'provider',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('provider') as string;
				const provider = name?.split('@')[1];
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
			accessorKey: 'router.entryPoints',
			id: 'entryPoints',
			enableSorting: true,
			cell: ({ row }) => {
				return renderComponent(ColumnBadge, {
					label: row.getValue('entryPoints') as string[],
					variant: 'secondary'
				});
			}
		},
		{
			header: 'Middlewares',
			accessorKey: 'router.middlewares',
			id: 'middlewares',
			enableSorting: true,
			cell: ({ row }) => {
				if (row.original.router.protocol === 'udp') return;
				return renderComponent(ColumnBadge, {
					label: row.getValue('middlewares') as string[],
					variant: 'secondary'
				});
			}
		},
		{
			header: 'Domains',
			accessorKey: 'router.rule',
			id: 'rule',
			enableSorting: true,
			cell: ({ row }) => {
				if (row.original.router.protocol === 'udp') return;
				return renderComponent(ColumnRule, {
					rule: row.getValue('rule') as string,
					protocol: row.original.router.protocol as 'http' | 'tcp'
				});
			}
		},
		{
			header: 'Cert Resolver',
			accessorKey: 'router.tls',
			id: 'tls',
			enableSorting: true,
			cell: ({ row }) => {
				const tls = row.getValue('tls') as TLS;
				if (!tls?.certResolver) {
					return renderComponent(ColumnBadge, {
						label: 'None',
						variant: 'secondary',
						class: 'bg-slate-300 dark:bg-slate-700'
					});
				}
				return renderComponent(ColumnBadge, {
					label: tls.certResolver as string,
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
				if (status === undefined) {
					return renderComponent(ColumnBadge, {
						label: 'N/A',
						variant: 'secondary'
					});
				}
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

	let columns: ColumnDef<RouterWithService>[] = $derived(
		source.value === TraefikSource.LOCAL
			? defaultColumns.filter((c) => c.id !== 'provider' && c.id !== 'serverStatus')
			: defaultColumns
	);

	onMount(async () => {
		await api.getTraefikConfig(source.value);
		await api.listDNSProviders();
	});
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<Route />
		<h1 class="text-2xl font-bold">Router Management</h1>
	</div>
	{#if source.value === TraefikSource.LOCAL}
		<DataTable
			{columns}
			data={$routerServiceMerge || []}
			showSourceTabs={true}
			createButton={{
				label: 'Add Router',
				onClick: openCreateModal
			}}
		/>
	{:else}
		<DataTable {columns} data={$routerServiceMerge || []} showSourceTabs={true} />
	{/if}
</div>

<RouterModal
	mode={modalState.mode}
	bind:open={modalState.isOpen}
	bind:router={modalState.router}
	bind:service={modalState.service}
/>
