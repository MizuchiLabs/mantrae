<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import RouterModal from '$lib/components/modals/router.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { Pencil, Route, Trash } from '@lucide/svelte';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import ColumnRule from '$lib/components/tables/ColumnRule.svelte';
	import { routerClient } from '$lib/api';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import { ConnectError } from '@connectrpc/connect';
	import { profile } from '$lib/stores/profile';
	import type { RouterTLSConfig } from '$lib/gen/tygo/dynamic';
	import type { BulkAction } from '$lib/components/tables/types';
	import { onMount } from 'svelte';
	import { pageIndex, pageSize } from '$lib/stores/common';

	let modalRouter = $state({} as Router);
	let modalRouterOpen = $state(false);

	// Data state
	let data = $state<Router[]>([]);
	let rowCount = $state<number>(0);

	const deleteRouter = async (id: bigint) => {
		try {
			await routerClient.deleteRouter({ id: id });
			await refreshRouters(pageSize.value ?? 10, 0);
			toast.success('Router deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
	};

	async function handleBulkDelete(selectedRows: Router[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${selectedRows.length} routers?`);
			if (!confirmed) return;

			const routerIds = selectedRows.map((row) => ({ id: row.id }));
			for (const router of routerIds) {
				await routerClient.deleteRouter(router);
			}
			await refreshRouters(pageSize.value ?? 10, 0);
			toast.success(`Successfully deleted ${selectedRows.length} routers`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete routers', { description: e.message });
		}
	}

	const columns: ColumnDef<Router>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			id: 'name',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('name') as string;
				return name?.split('@')[0];
			}
		},
		{
			header: 'Protocol',
			accessorKey: 'type',
			id: 'protocol',
			enableSorting: true,
			cell: ({ row, column }) => {
				let protocol = row.getValue('protocol') as
					| RouterType.HTTP
					| RouterType.TCP
					| RouterType.UDP;

				let label = 'Unspecified';
				if (protocol === RouterType.HTTP) {
					label = 'HTTP';
				} else if (protocol === RouterType.TCP) {
					label = 'TCP';
				} else if (protocol === RouterType.UDP) {
					label = 'UDP';
				}
				return renderComponent(ColumnBadge<Router>, {
					label: label,
					class: 'hover:cursor-pointer italic',
					column: column
				});
			}
		},
		{
			header: 'Entrypoints',
			accessorKey: 'entryPoints',
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
			accessorKey: 'middlewares',
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
			accessorKey: 'config.rule',
			id: 'rule',
			enableSorting: true,
			cell: ({ row }) => {
				let rule = '';
				if (row.original.config?.rule !== undefined) {
					rule = row.getValue('rule') as string;
				}
				return renderComponent(ColumnRule, {
					rule: rule,
					routerType: row.original.type as RouterType.HTTP | RouterType.TCP
				});
			}
		},
		{
			header: 'TLS',
			accessorKey: 'config.tls',
			id: 'tls',
			enableSorting: true,
			filterFn: (row, columnId, filterValue) => {
				const tls = row.getValue(columnId) as RouterTLSConfig;
				return tls?.certResolver === filterValue;
			},
			cell: ({ row, column }) => {
				let tls = undefined;
				if (row.original.config?.tls !== undefined) {
					tls = row.getValue('tls') as RouterTLSConfig;
				}

				let label = 'Disabled';
				if (tls) {
					label = tls.certResolver ? tls.certResolver : 'Enabled';
				}
				return renderComponent(ColumnBadge<Router>, {
					label,
					variant: tls ? 'secondary' : 'outline',
					class: tls ? 'bg-slate-300 dark:bg-slate-700' : '',
					tls,
					column: tls?.certResolver ? column : undefined
				});
			}
		},
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
								modalRouter = row.original;
								modalRouterOpen = true;
							}
						},
						{
							type: 'button',
							label: 'Delete Router',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteRouter(row.original.id)
						}
					]
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

	async function handlePaginationChange(p: PaginationState) {
		await refreshRouters(p.pageSize, p.pageIndex);
	}
	async function refreshRouters(pageSize: number, pageIndex: number) {
		const response = await routerClient.listRouters({
			profileId: profile.id,
			limit: BigInt(pageSize),
			offset: BigInt(pageIndex * pageSize)
		});
		data = response.routers;
		rowCount = Number(response.totalCount);
	}

	onMount(async () => {
		await refreshRouters(pageSize.value ?? 10, pageIndex.value ?? 0);
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
	<DataTable
		{data}
		{columns}
		{rowCount}
		createButton={{
			label: 'Create Router',
			onClick: () => {
				modalRouter = { type: RouterType.HTTP } as Router;
				modalRouterOpen = true;
			}
		}}
		onRowSelection={(selections) => console.log(selections)}
		onPaginationChange={handlePaginationChange}
		bulkActions={routerBulkActions}
	/>
</div>

<RouterModal bind:open={modalRouterOpen} bind:router={modalRouter} bind:data />
