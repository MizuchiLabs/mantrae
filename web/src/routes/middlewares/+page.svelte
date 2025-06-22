<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import MiddlewareModal from '$lib/components/modals/middleware.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { Layers, Pencil, Trash } from '@lucide/svelte';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { profile } from '$lib/stores/profile';
	import type { BulkAction } from '$lib/components/tables/types';
	import { MiddlewareType, type Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
	import { onMount } from 'svelte';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { middlewareClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import type { TCPMiddleware } from '$lib/gen/tygo/dynamic';
	import type { Middleware as HTTPMiddleware } from '$lib/gen/tygo/dynamic';
	import type { JsonObject } from '@bufbuild/protobuf';

	let item = $state({} as Middleware);
	let open = $state(false);

	// Data state
	let data = $state<Middleware[]>([]);
	let rowCount = $state<number>(0);

	const columns: ColumnDef<Middleware>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('name') as string;
				return name?.split('@')[0];
			}
		},
		{
			header: 'Protocol',
			accessorKey: 'type',
			enableSorting: true,
			cell: ({ row, column }) => {
				let protocol = row.getValue('type') as MiddlewareType.HTTP | MiddlewareType.TCP;

				let label = 'Unspecified';
				if (protocol === MiddlewareType.HTTP) {
					label = 'HTTP';
				} else if (protocol === MiddlewareType.TCP) {
					label = 'TCP';
				}
				return renderComponent(ColumnBadge<Middleware>, {
					label: label,
					class: 'hover:cursor-pointer italic',
					column: column
				});
			}
		},
		{
			header: 'Type',
			accessorKey: 'config',
			enableSorting: true,
			cell: ({ row, column }) => {
				let config = row.getValue('config') as JsonObject;
				let label = config ? Object.keys(config)[0] : 'unknown';

				return renderComponent(ColumnBadge<Middleware>, {
					label: label,
					class: 'hover:cursor-pointer',
					column: column
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
							label: 'Edit Middleware',
							icon: Pencil,
							onClick: () => {
								item = row.original;
								open = true;
							}
						},
						{
							type: 'button',
							label: 'Delete Middleware',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original.id, row.original.type)
						}
					]
				});
			}
		}
	];

	const bulkActions: BulkAction<Middleware>[] = [
		{
			type: 'button',
			label: 'Delete',
			icon: Trash,
			variant: 'destructive',
			onClick: bulkDelete
		}
	];

	async function onPaginationChange(p: PaginationState) {
		await refreshData(p.pageSize, p.pageIndex);
	}

	const deleteItem = async (id: bigint, type: MiddlewareType) => {
		try {
			await middlewareClient.deleteMiddleware({ id: id, type: type });
			await refreshData(pageSize.value ?? 10, 0);
			toast.success('Router deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
	};

	async function bulkDelete(selectedRows: Middleware[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${selectedRows.length} routers?`);
			if (!confirmed) return;

			const rows = selectedRows.map((row) => ({ id: row.id, type: row.type }));
			for (const row of rows) {
				await middlewareClient.deleteMiddleware({ id: row.id, type: row.type });
			}
			await refreshData(pageSize.value ?? 10, 0);
			toast.success(`Successfully deleted ${selectedRows.length} routers`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete routers', { description: e.message });
		}
	}
	async function refreshData(pageSize: number, pageIndex: number) {
		const response = await middlewareClient.listMiddlewares({
			profileId: profile.id,
			limit: BigInt(pageSize),
			offset: BigInt(pageIndex * pageSize)
		});
		data = response.middlewares;
		rowCount = Number(response.totalCount);
	}

	onMount(async () => {
		await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
	});
</script>

<svelte:head>
	<title>Middlewares</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<Layers />
		<h1 class="text-2xl font-bold">Middleware Management</h1>
	</div>
	<DataTable
		{data}
		{columns}
		{rowCount}
		{onPaginationChange}
		{bulkActions}
		createButton={{
			label: 'Create Middleware',
			onClick: () => {
				item = { type: MiddlewareType.HTTP } as Middleware;
				open = true;
			}
		}}
	/>
</div>

<MiddlewareModal bind:open bind:item bind:data />
