<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import RouterModal from '$lib/components/modals/router.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { Pencil, Route, Trash } from '@lucide/svelte';
	import { TraefikSource } from '$lib/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import ColumnRule from '$lib/components/tables/ColumnRule.svelte';
	import { profile } from '$lib/stores/profile';
	import type { BulkAction } from '$lib/components/tables/types';
	import { routerClient } from '$lib/api';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { Service } from '$lib/gen/mantrae/v1/service_pb';

	interface ModalState {
		isOpen: boolean;
		mode: 'create' | 'edit';
		router: Router;
		service: Service;
	}

	const initialModalState: ModalState = {
		isOpen: false,
		mode: 'create',
		router: {} as Router,
		service: {} as Service
	};
	let modalState = $state(initialModalState);

	function openCreateModal() {
		modalState = {
			isOpen: true,
			mode: 'create',
			router: {} as Router,
			service: {} as Service
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
			if (!profile.hasValidId() || !profile.id) {
				toast.error('Invalid profile ID');
				return;
			}
			await routerClient.deleteRouter({ id: router.id });
			toast.success('Router deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

	// async function handleBulkDelete(selectedRows: RouterWithService[]) {
	// 	try {
	// 		const confirmed = confirm(`Are you sure you want to delete ${selectedRows.length} routers?`);
	// 		if (!confirmed) return;
	//
	// 		const items = selectedRows.map((row) => ({
	// 			name: row.router.name,
	// 			protocol: row.router.protocol
	// 		}));
	// 		await api.bulkDeleteRouter(items);
	// 		toast.success(`Successfully deleted ${selectedRows.length} routers`);
	// 	} catch (err: unknown) {
	// 		const e = err as Error;
	// 		toast.error(`Failed to delete routers: ${e.message}`);
	// 	}
	// }

	const columns: ColumnDef<Router>[] = [
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
			accessorKey: 'router.type',
			id: 'protocol',
			enableSorting: true,
			cell: ({ row, column }) => {
				return renderComponent(ColumnBadge<Router>, {
					label: row.getValue('protocol') as string,
					class: 'hover:cursor-pointer',
					column: column
				});
			}
		},
		{
			header: 'Provider',
			accessorKey: 'router.name',
			id: 'provider',
			enableSorting: true,
			cell: ({ row, column }) => {
				const name = row.getValue('provider') as string;
				const provider = name?.split('@')[1];
				return renderComponent(ColumnBadge<Router>, {
					label: provider ? provider.toLowerCase() : 'http',
					variant: 'secondary',
					class: 'hover:cursor-pointer',
					column: column
				});
			}
		},
		{
			header: 'Entrypoints',
			accessorKey: 'router.entryPoints',
			id: 'entrypoints',
			enableSorting: true,
			filterFn: 'arrIncludes',
			cell: ({ row, column }) => {
				const entrypoints = row.getValue('entrypoints') as string[];
				return renderComponent(ColumnBadge<Router>, {
					label: entrypoints?.length ? entrypoints : 'None',
					variant: entrypoints?.length ? 'secondary' : 'outline',
					class: 'hover:cursor-pointer',
					column: entrypoints?.length ? column : undefined
				});
			}
		},
		{
			header: 'Middlewares',
			accessorKey: 'router.middlewares',
			id: 'middlewares',
			enableSorting: true,
			filterFn: 'arrIncludes',
			cell: ({ row, column }) => {
				const middlewares = row.getValue('middlewares') as string[];
				return renderComponent(ColumnBadge<Router>, {
					label: middlewares?.length ? middlewares : 'None',
					variant: middlewares?.length ? 'secondary' : 'outline',
					class: 'hover:cursor-pointer',
					column: middlewares?.length ? column : undefined
				});
			}
		},
		{
			header: 'Rules',
			accessorKey: 'router.rule',
			id: 'rule',
			enableSorting: true,
			cell: ({ row }) => {
				if (row.original.type === RouterType.UDP) return;
				return renderComponent(ColumnRule, {
					rule: row.getValue('rule') as string,
					type: row.original.type as RouterType.HTTP | RouterType.TCP
				});
			}
		},
		// {
		// 	header: 'DNS',
		// 	accessorKey: 'router.name',
		// 	id: 'dns',
		// 	enableSorting: true,
		// 	filterFn: (row, columnId, filterValue) => {
		// 		const routerName = row.getValue(columnId) as string;
		// 		const matches = $rdps?.some(
		// 			(rdp) => rdp.routerName === routerName && rdp.providerName === filterValue
		// 		);
		// 		return !!matches;
		// 	},
		// 	cell: ({ row, column }) => {
		// 		// Return early if no rdps data
		// 		if (!$rdps) {
		// 			return renderComponent(ColumnBadge, {
		// 				label: ['Disabled'],
		// 				variant: 'outline'
		// 			});
		// 		}
		// 		const name = row.getValue('dns') as string;
		// 		const dns = $rdps?.filter((item) => item.routerName === name);
		// 		const rdpNames = dns.length ? [...new Set(dns.map((item) => item.providerName))] : [];
		//
		// 		return renderComponent(ColumnBadge<RouterWithService>, {
		// 			label: rdpNames.length ? rdpNames : ['Disabled'],
		// 			variant: rdpNames.length ? 'secondary' : 'outline',
		// 			class: rdpNames.length ? 'bg-blue-300 dark:bg-blue-700' : undefined,
		// 			column: rdpNames.length ? column : undefined
		// 		});
		// 	}
		// },
		// {
		// 	header: 'TLS',
		// 	accessorKey: 'router.tls',
		// 	id: 'tls',
		// 	enableSorting: true,
		// 	filterFn: (row, columnId, filterValue) => {
		// 		const tls = row.getValue(columnId) as TLS;
		// 		return tls?.certResolver === filterValue;
		// 	},
		// 	cell: ({ row, column }) => {
		// 		const tls = row.getValue('tls') as TLS;
		//
		// 		let label = 'Disabled';
		// 		if (tls) {
		// 			label = tls.certResolver ? tls.certResolver : 'Enabled';
		// 		}
		// 		return renderComponent(ColumnBadge<RouterWithService>, {
		// 			label,
		// 			variant: tls ? 'secondary' : 'outline',
		// 			class: tls ? 'bg-slate-300 dark:bg-slate-700' : '',
		// 			tls,
		// 			column: tls?.certResolver ? column : undefined
		// 		});
		// 	}
		// },
		// {
		// 	header: 'Server Status',
		// 	accessorFn: (row) => row.service.serverStatus,
		// 	id: 'serverStatus',
		// 	enableSorting: true,
		// 	cell: ({ row }) => {
		// 		const status = row.getValue('serverStatus') as Record<string, string>;
		// 		if (status === undefined) {
		// 			return renderComponent(ColumnBadge, {
		// 				label: 'N/A',
		// 				variant: 'secondary'
		// 			});
		// 		}
		// 		const upCount = Object.values(status).filter((status) => status === 'UP').length;
		// 		const totalCount = Object.values(status).length;
		// 		const greenBadge = 'bg-green-300 dark:bg-green-600';
		// 		const redBadge = 'bg-red-300 dark:bg-red-600';
		// 		return renderComponent(ColumnBadge, {
		// 			label: `${upCount}/${totalCount}`,
		// 			variant: 'secondary',
		// 			class: upCount === totalCount ? greenBadge : redBadge
		// 		});
		// 	}
		// },
		{
			id: 'actions',
			enableHiding: false,
			cell: ({ row }) => {
				return renderComponent(TableActions, {
					actions: [
						{
							type: 'button',
							label: 'Edit Router',
							icon: Pencil,
							onClick: () => {
								openEditModal(row.original.router, row.original.service);
							}
						},
						{
							type: 'button',
							label: 'Delete Router',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => {
								deleteRouter(row.original.router);
							},
							disabled: source.value !== TraefikSource.LOCAL
						}
					],
					shareObject: source.value === TraefikSource.LOCAL ? row.original : undefined
				});
			}
		}
	];

	const routerBulkActions: BulkAction<Router>[] = [
		{
			type: 'button',
			label: 'Delete',
			icon: Trash,
			variant: 'destructive',
			onClick: handleBulkDelete
		}
	];
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<Route />
		<h1 class="text-2xl font-bold">Router Management</h1>
	</div>
	{#await routerClient.listRouters({}) then value}
		<DataTable
			{columns}
			data={value.routers}
			bulkActions={routerBulkActions}
			createButton={{
				label: 'Add Router',
				onClick: openCreateModal
			}}
		/>
	{/await}
</div>

<RouterModal
	mode={modalState.mode}
	bind:open={modalState.isOpen}
	bind:router={modalState.router}
	bind:service={modalState.service}
/>
