<script lang="ts">
	import { serversTransportClient } from '$lib/api';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { type ServersTransport } from '$lib/gen/mantrae/v1/servers_transport_pb';
	import { profile } from '$lib/stores/profile';
	import type { IconComponent } from '$lib/types';
	import { ConnectError } from '@connectrpc/connect';
	import ServerTransportModal from '$lib/components/modals/ServerTransportModal.svelte';
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
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { serversTransports } from '$lib/stores/realtime';

	let item = $state({} as ServersTransport);
	let open = $state(false);

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
								row.original.enabled = !row.original.enabled;
								updateItem(row.original);
							}
						},
						{
							type: 'button',
							label: 'Edit Transport',
							icon: Pencil,
							onClick: () => {
								item = row.original;
								open = true;
							}
						},
						{
							type: 'popover',
							label: 'Delete Transport',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original.id),
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

	async function deleteItem(id: bigint) {
		try {
			await serversTransportClient.deleteServersTransport({ id: id });
			toast.success('Transport deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete transport', { description: e.message });
		}
	}

	async function updateItem(item: ServersTransport) {
		try {
			await serversTransportClient.updateServersTransport({
				id: item.id,
				name: item.name,
				config: item.config,
				type: item.type,
				enabled: item.enabled
			});
			toast.success(`Transport ${item.name} updated`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to update transport', { description: e.message });
		}
	}

	async function bulkDelete(rows: ServersTransport[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} transports?`);
			if (!confirmed) return;

			for (const s of rows) {
				await serversTransportClient.deleteServersTransport({ id: s.id });
			}
			toast.success(`Successfully deleted ${rows.length} transports`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete transports', { description: e.message });
		}
	}

	$effect(() => {
		if (profile.isValid()) {
			serversTransportClient.listServersTransports({ profileId: profile.id }).then((response) => {
				serversTransports.set(response.serversTransports);
			});
		}
	});
</script>

<svelte:head>
	<title>Server Transports</title>
</svelte:head>

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="bg-primary/10 rounded-lg p-2">
					<Truck class="text-primary h-6 w-6" />
				</div>
				Server Transports
			</h1>
			<p class="text-muted-foreground mt-1">Manage your server transports</p>
		</div>
	</div>

	<DataTable
		data={$serversTransports}
		{columns}
		{bulkActions}
		createButton={{
			label: 'Create Transport',
			onClick: () => {
				item = {} as ServersTransport;
				open = true;
			}
		}}
	/>
</div>

<ServerTransportModal bind:open bind:item />
