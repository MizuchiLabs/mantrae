<script lang="ts">
	import { entrypoint } from '$lib/api/entrypoints.svelte';
	import EntryPointModal from '$lib/components/modals/EntryPointModal.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import type { EntryPoint } from '$lib/gen/mantrae/v1/entry_point_pb';
	import { ConnectError } from '@connectrpc/connect';
	import { CircleCheck, CircleSlash, EthernetPort, Pencil, Trash } from '@lucide/svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { toast } from 'svelte-sonner';

	let data = $state({} as EntryPoint);
	let open = $state(false);

	const entryPointList = entrypoint.list();
	const updateEntryPoint = entrypoint.update();
	const deleteEntryPoint = entrypoint.delete();

	const columns: ColumnDef<EntryPoint>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true,
			enableHiding: false
		},
		{
			header: 'Address',
			accessorKey: 'address',
			enableSorting: true,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				return renderComponent(ColumnBadge, {
					label: row.original.address || 'None',
					class: 'hover:cursor-pointer'
				});
			}
		},
		{
			header: 'Default',
			accessorKey: 'isDefault',
			enableGlobalFilter: false,
			cell: ({ row }) => {
				return renderComponent(TableActions, {
					actions: [
						{
							type: 'button',
							label: row.original.isDefault ? 'Disable' : 'Enable',
							icon: row.original.isDefault ? CircleCheck : CircleSlash,
							iconProps: {
								class: row.original.isDefault ? 'text-green-500 size-5' : 'text-red-500 size-5',
								size: 20
							},
							onClick: () =>
								updateEntryPoint.mutate({ ...row.original, isDefault: !row.original.isDefault })
						}
					]
				});
			}
		},
		{
			id: 'actions',
			enableHiding: false,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				return renderComponent(TableActions, {
					actions: [
						{
							type: 'button',
							label: 'Edit EntryPoint',
							icon: Pencil,
							onClick: () => {
								data = row.original;
								open = true;
							}
						},
						{
							type: 'popover',
							label: 'Delete EntryPoint',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteEntryPoint.mutate({ id: row.original.id }),
							popover: {
								title: 'Delete EntryPoint?',
								description: 'This entry point will be permanently deleted.',
								confirmLabel: 'Delete',
								cancelLabel: 'Cancel'
							}
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

	async function bulkDelete(rows: EntryPoint[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} entrypoints?`);
			if (!confirmed) return;

			for (const e of rows) {
				deleteEntryPoint.mutate({ id: e.id });
			}
			toast.success(`Successfully deleted ${rows.length} entrypoints`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete entry points', { description: e.message });
		}
	}
</script>

<svelte:head>
	<title>EntryPoints - Mantrae</title>
	<meta
		name="description"
		content="Configure entrypoints for your reverse proxy to listen for incoming connections"
	/>
</svelte:head>

<EntryPointModal bind:open {data} />

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="rounded-lg bg-primary/10 p-2">
					<EthernetPort class="h-6 w-6 text-primary" />
				</div>
				Entry Points
			</h1>
			<p class="mt-1 text-muted-foreground">Manage your entry points</p>
		</div>
	</div>

	<DataTable
		data={entryPointList.data}
		{columns}
		{bulkActions}
		createButton={{
			label: 'Create EntryPoint',
			onClick: () => (open = true)
		}}
	/>
</div>
