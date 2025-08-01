<script lang="ts">
	import { entryPointClient } from '$lib/api';
	import EntryPointModal from '$lib/components/modals/EntryPointModal.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import type { EntryPoint } from '$lib/gen/mantrae/v1/entry_point_pb';
	import { profile } from '$lib/stores/profile';
	import { entryPoints } from '$lib/stores/realtime';
	import { ConnectError } from '@connectrpc/connect';
	import { CircleCheck, CircleSlash, EthernetPort, Pencil, Trash } from '@lucide/svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { toast } from 'svelte-sonner';

	let item = $state({} as EntryPoint);
	let open = $state(false);

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
							onClick: () => toggleItem(row.original, !row.original.isDefault)
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
								item = row.original;
								open = true;
							}
						},
						{
							type: 'popover',
							label: 'Delete EntryPoint',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original.id),
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

	async function deleteItem(id: bigint) {
		try {
			await entryPointClient.deleteEntryPoint({ id: id });
			toast.success('EntryPoint deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete entrypoint', { description: e.message });
		}
	}

	async function toggleItem(item: EntryPoint, isDefault: boolean) {
		try {
			await entryPointClient.updateEntryPoint({
				id: item.id,
				name: item.name,
				address: item.address,
				isDefault: isDefault
			});
			toast.success(
				`Entry point ${item.name} ${isDefault ? 'set as default' : 'removed as default'}`
			);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to update entry point', { description: e.message });
		}
	}

	async function bulkDelete(rows: EntryPoint[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} entrypoints?`);
			if (!confirmed) return;

			for (const e of rows) {
				await entryPointClient.deleteEntryPoint({ id: e.id });
			}
			toast.success(`Successfully deleted ${rows.length} entrypoints`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete entry points', { description: e.message });
		}
	}

	$effect(() => {
		if (profile.isValid()) {
			entryPointClient.listEntryPoints({ profileId: profile.id }).then((response) => {
				entryPoints.set(response.entryPoints);
			});
		}
	});
</script>

<svelte:head>
	<title>EntryPoints</title>
</svelte:head>

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="bg-primary/10 rounded-lg p-2">
					<EthernetPort class="text-primary h-6 w-6" />
				</div>
				Entry Points
			</h1>
			<p class="text-muted-foreground mt-1">Manage your entry points</p>
		</div>
	</div>

	<DataTable
		data={$entryPoints}
		{columns}
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

<EntryPointModal bind:open bind:item />
