<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import MiddlewareModal from '$lib/components/modals/MiddlewareModal.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import {
		Bot,
		CircleCheck,
		CircleSlash,
		Globe,
		Layers,
		Network,
		Pencil,
		Power,
		PowerOff,
		Trash
	} from '@lucide/svelte';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import type { BulkAction } from '$lib/components/tables/types';
	import { type Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
	import { ConnectError } from '@connectrpc/connect';
	import ColumnText from '$lib/components/tables/ColumnText.svelte';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { middleware } from '$lib/api/middleware.svelte';

	let data = $state({} as Middleware);
	let open = $state(false);

	const middlewareList = middleware.list();
	const updateMiddleware = middleware.update();
	const deleteMiddleware = middleware.delete();

	const columns: ColumnDef<Middleware>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true,
			enableHiding: false,
			cell: ({ row }) => {
				return renderComponent(ColumnText, {
					label: row.getValue('name') as string,
					icon: row.original.agentId ? Bot : undefined,
					iconProps: { class: 'text-green-500', size: 20 },
					class: 'text-sm'
				});
			}
		},
		{
			header: 'Protocol',
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
				let protocol = row.getValue('type') as ProtocolType.HTTP | ProtocolType.TCP;

				let label = 'Unspecified';
				let icon = undefined;
				if (protocol === ProtocolType.HTTP) {
					label = 'HTTP';
					icon = Globe;
				} else if (protocol === ProtocolType.TCP) {
					label = 'TCP';
					icon = Network;
				}
				return renderComponent(ColumnBadge<Middleware>, {
					label,
					icon,
					variant: 'outline',
					class: 'hover:cursor-pointer',
					column: column
				});
			}
		},
		{
			header: 'Type',
			accessorKey: 'config',
			enableSorting: true,
			enableGlobalFilter: false,
			filterFn: (row, filterValue) => {
				const label = Object.keys(row.original.config ?? {})[0] ?? 'unknown';
				return label.toLowerCase().includes(filterValue.toLowerCase());
			},
			cell: ({ row, column }) => {
				let label = Object.keys(row.original.config ?? {})[0] ?? 'unknown';
				return renderComponent(ColumnBadge<Middleware>, {
					label: label,
					class: 'hover:cursor-pointer',
					column: column
				});
			}
		},
		{
			header: 'Default',
			accessorKey: 'isDefault',
			enableSorting: true,
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
							onClick: () => {
								updateMiddleware.mutate({ ...row.original, isDefault: !row.original.isDefault });
							}
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
							label: row.original.enabled ? 'Disable' : 'Enable',
							icon: row.original.enabled ? Power : PowerOff,
							iconProps: {
								class: row.original.enabled ? 'text-green-500' : 'text-red-500'
							},
							onClick: () => {
								updateMiddleware.mutate({ ...row.original, enabled: !row.original.enabled });
							}
						},
						{
							type: 'button',
							label: 'Edit Middleware',
							icon: Pencil,
							onClick: () => {
								data = row.original;
								open = true;
							}
						},
						{
							type: 'popover',
							label: 'Delete Middleware',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteMiddleware.mutate({ ...row.original }),
							popover: {
								title: 'Delete Middleware?',
								description: 'This middleware and its configuration will be permanently deleted.',
								confirmLabel: 'Delete',
								cancelLabel: 'Cancel'
							}
						}
					]
				});
			}
		}
	];
	function getProtocolLabel(protocol: ProtocolType): string {
		if (protocol === ProtocolType.HTTP) return 'HTTP';
		if (protocol === ProtocolType.TCP) return 'TCP';
		return 'Unspecified';
	}
	const bulkActions: BulkAction<Middleware>[] = [
		{
			type: 'button',
			label: 'Delete',
			icon: Trash,
			variant: 'destructive',
			onClick: bulkDelete
		}
	];

	async function bulkDelete(selectedRows: Middleware[]) {
		try {
			const confirmed = confirm(
				`Are you sure you want to delete ${selectedRows.length} middlewares?`
			);
			if (!confirmed) return;

			const rows = selectedRows.map((row) => ({ id: row.id, type: row.type }));
			for (const row of rows) {
				deleteMiddleware.mutate({ ...row });
			}
			toast.success(`Successfully deleted ${selectedRows.length} middlewares`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete middlewares', { description: e.message });
		}
	}
</script>

<svelte:head>
	<title>Middlewares - Mantrae</title>
	<meta
		name="description"
		content="Manage HTTP and TCP middlewares to customize your reverse proxy behavior"
	/>
</svelte:head>

<MiddlewareModal bind:open {data} />

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="rounded-lg bg-primary/10 p-2">
					<Layers class="h-6 w-6 text-primary" />
				</div>
				Middlewares
			</h1>
			<p class="mt-1 text-muted-foreground">Configure your middlewares</p>
		</div>
	</div>

	<DataTable
		data={middlewareList.data}
		{columns}
		{bulkActions}
		createButton={{
			label: 'Create Middleware',
			onClick: () => {
				data = { type: ProtocolType.HTTP } as Middleware;
				open = true;
			}
		}}
	/>
</div>
