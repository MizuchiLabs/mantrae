<script lang="ts">
	import { transport } from '$lib/api/transport.svelte';
	import ServerTransportModal from '$lib/components/modals/ServerTransportModal.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { type ServersTransport } from '$lib/gen/mantrae/v1/servers_transport_pb';
	import type { IconComponent } from '$lib/types';
	import { ConnectError } from '@connectrpc/connect';
	import {
		Globe,
		Network,
		Pencil,
		Power,
		PowerOff,
		Trash,
		TriangleAlert,
		Truck
	} from '@lucide/svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { toast } from 'svelte-sonner';

	let data = $state({} as ServersTransport);
	let open = $state(false);

	const transportList = transport.list();
	const deleteTransport = transport.delete();
	const updateTransport = transport.update();

	const columns: ColumnDef<ServersTransport>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true,
			enableHiding: false
		},
		{
			header: 'Type',
			accessorKey: 'type',
			enableSorting: true,
			enableGlobalFilter: false,
			filterFn: (row, columnId, filterValue) => {
				const protocol = row.getValue(columnId) as ProtocolType;

				// Handle both enum value and display label filtering
				if (typeof filterValue === 'string') {
					const displayLabel = getProtocolLabel(protocol);
					return (
						displayLabel.toLowerCase().includes(filterValue.toLowerCase()) ||
						protocol.toString().toLowerCase().includes(filterValue.toLowerCase())
					);
				}

				// Direct enum comparison for badge clicking
				return protocol === filterValue;
			},
			cell: ({ row, column }) => {
				const protocol = row.getValue('type') as ProtocolType;
				const label = getProtocolLabel(protocol);
				const iconMap: Partial<Record<ProtocolType, IconComponent>> = {
					[ProtocolType.HTTP]: Globe,
					[ProtocolType.TCP]: Network,
					[ProtocolType.UNSPECIFIED]: TriangleAlert
				};
				return renderComponent(ColumnBadge<ServersTransport>, {
					label,
					icon: iconMap[protocol],
					variant: 'outline',
					column: column
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
							label: row.original.enabled ? 'Disable' : 'Enable',
							icon: row.original.enabled ? Power : PowerOff,
							iconProps: {
								class: row.original.enabled ? 'text-green-500' : 'text-red-500'
							},
							onClick: () => {
								updateTransport.mutate({ ...row.original, enabled: !row.original.enabled });
							}
						},
						{
							type: 'button',
							label: 'Edit Transport',
							icon: Pencil,
							onClick: () => {
								data = row.original;
								open = true;
							}
						},
						{
							type: 'popover',
							label: 'Delete Transport',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteTransport.mutate({ ...row.original }),
							popover: {
								title: 'Delete Transport?',
								description: 'This transport will be permanently deleted.',
								confirmLabel: 'Delete',
								cancelLabel: 'Cancel'
							}
						}
					]
				});
			}
		}
	];

	// Helper functions to avoid repetition
	function getProtocolLabel(protocol: ProtocolType): string {
		if (protocol === ProtocolType.HTTP) return 'HTTP';
		if (protocol === ProtocolType.TCP) return 'TCP';
		return 'Unspecified';
	}

	const bulkActions: BulkAction<ServersTransport>[] = [
		{
			type: 'button',
			label: 'Delete',
			icon: Trash,
			variant: 'destructive',
			onClick: bulkDelete
		}
	];

	async function bulkDelete(rows: ServersTransport[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} transports?`);
			if (!confirmed) return;

			for (const s of rows) {
				deleteTransport.mutate({ id: s.id });
			}
			toast.success(`Successfully deleted ${rows.length} transports`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete transports', { description: e.message });
		}
	}
</script>

<svelte:head>
	<title>Server Transports - Mantrae</title>
	<meta
		name="description"
		content="Configure HTTP and TCP server transports for your reverse proxy services"
	/>
</svelte:head>

<ServerTransportModal bind:open {data} />

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="rounded-lg bg-primary/10 p-2">
					<Truck class="h-6 w-6 text-primary" />
				</div>
				Server Transports
			</h1>
			<p class="mt-1 text-muted-foreground">Manage your server transports</p>
		</div>
	</div>

	<DataTable
		data={transportList.data}
		{columns}
		{bulkActions}
		createButton={{
			label: 'Create Transport',
			onClick: () => (open = true)
		}}
	/>
</div>
