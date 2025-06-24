<script lang="ts">
	import { entryPointClient, routerClient } from '$lib/api';
	import EntryPointModal from '$lib/components/modals/entrypoint.svelte';
	import ColumnCheck from '$lib/components/tables/ColumnCheck.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import type { EntryPoint } from '$lib/gen/mantrae/v1/entry_point_pb';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { profile } from '$lib/stores/profile';
	import { ConnectError } from '@connectrpc/connect';
	import { EthernetPort, Pencil, Trash } from '@lucide/svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	let item = $state({} as EntryPoint);
	let open = $state(false);

	// Data state
	let data = $state<EntryPoint[]>([]);
	let rowCount = $state<number>(0);

	const columns: ColumnDef<EntryPoint>[] = [
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
			header: 'Address',
			accessorKey: 'address',
			id: 'address',
			enableSorting: true,
			cell: ({ row }) => {
				const address = row.getValue('address') as string;
				return address;
			}
		},
		{
			header: 'Default',
			accessorKey: 'isDefault',
			id: 'isDefault',
			enableHiding: false,
			cell: ({ row }) => {
				let isDefault = row.getValue('isDefault') as boolean;
				return renderComponent(ColumnCheck, { checked: isDefault });
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
							label: 'Edit EntryPoint',
							icon: Pencil,
							onClick: () => {
								item = row.original;
								open = true;
							}
						},
						{
							type: 'button',
							label: 'Delete EntryPoint',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original.id)
						}
					]
				});
			}
		}
	];

	const bulkActions: BulkAction<EntryPoint>[] = [
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

	async function deleteItem(id: bigint) {
		try {
			await entryPointClient.deleteEntryPoint({ id: id });
			await refreshData(pageSize.value ?? 10, 0);
			toast.success('EntryPoint deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
	}

	async function bulkDelete(selectedRows: EntryPoint[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${selectedRows.length} routers?`);
			if (!confirmed) return;

			const routerIds = selectedRows.map((row) => ({ id: row.id }));
			for (const router of routerIds) {
				await routerClient.deleteRouter(router);
			}
			await refreshData(pageSize.value ?? 10, 0);
			toast.success(`Successfully deleted ${selectedRows.length} routers`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete routers', { description: e.message });
		}
	}

	async function refreshData(pageSize: number, pageIndex: number) {
		const response = await entryPointClient.listEntryPoints({
			profileId: profile.id,
			limit: BigInt(pageSize),
			offset: BigInt(pageIndex * pageSize)
		});
		data = response.entryPoints;
		rowCount = Number(response.totalCount);
	}

	onMount(async () => {
		await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
	});
</script>

<svelte:head>
	<title>EntryPoints</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<EthernetPort />
		<h1 class="text-2xl font-bold">Entry Points</h1>
	</div>
	<DataTable
		{data}
		{columns}
		{rowCount}
		{onPaginationChange}
		{bulkActions}
		createButton={{
			label: 'Create EntryPoint',
			onClick: () => {
				item = {} as EntryPoint;
				open = true;
			}
		}}
	/>
</div>

<EntryPointModal bind:open bind:item bind:data />
