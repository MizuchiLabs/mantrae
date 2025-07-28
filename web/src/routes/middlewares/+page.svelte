<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import MiddlewareModal from '$lib/components/modals/MiddlewareModal.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
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
	import { profile } from '$lib/stores/profile';
	import type { BulkAction } from '$lib/components/tables/types';
	import { MiddlewareType, type Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
	import { onMount } from 'svelte';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { middlewareClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import ColumnText from '$lib/components/tables/ColumnText.svelte';

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
				const protocol = row.getValue(columnId) as MiddlewareType;

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
				let protocol = row.getValue('type') as MiddlewareType.HTTP | MiddlewareType.TCP;

				let label = 'Unspecified';
				let icon = undefined;
				if (protocol === MiddlewareType.HTTP) {
					label = 'HTTP';
					icon = Globe;
				} else if (protocol === MiddlewareType.TCP) {
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
								row.original.isDefault = !row.original.isDefault;
								updateItem(row.original);
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
								row.original.enabled = !row.original.enabled;
								updateItem(row.original);
							}
						},
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
							type: 'popover',
							label: 'Delete Middleware',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original.id, row.original.type),
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
	function getProtocolLabel(protocol: MiddlewareType): string {
		if (protocol === MiddlewareType.HTTP) return 'HTTP';
		if (protocol === MiddlewareType.TCP) return 'TCP';
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

	async function onPaginationChange(p: PaginationState) {
		await refreshData(p.pageSize, p.pageIndex);
	}

	const deleteItem = async (id: bigint, type: MiddlewareType) => {
		try {
			await middlewareClient.deleteMiddleware({ id: id, type: type });
			await refreshData(pageSize.value ?? 10, 0);
			toast.success('Middleware deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete middleware', { description: e.message });
		}
	};

	const updateItem = async (item: Middleware) => {
		try {
			await middlewareClient.updateMiddleware({
				id: item.id,
				name: item.name,
				type: item.type,
				config: item.config,
				enabled: item.enabled,
				isDefault: item.isDefault
			});
			await refreshData(pageSize.value ?? 10, 0);
			toast.success(`Middleware ${item.name} updated`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to update middleware', { description: e.message });
		}
	};

	async function bulkDelete(selectedRows: Middleware[]) {
		try {
			const confirmed = confirm(
				`Are you sure you want to delete ${selectedRows.length} middlewares?`
			);
			if (!confirmed) return;

			const rows = selectedRows.map((row) => ({ id: row.id, type: row.type }));
			for (const row of rows) {
				await middlewareClient.deleteMiddleware({ id: row.id, type: row.type });
			}
			await refreshData(pageSize.value ?? 10, 0);
			toast.success(`Successfully deleted ${selectedRows.length} middlewares`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete middlewares', { description: e.message });
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

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="bg-primary/10 rounded-lg p-2">
					<Layers class="text-primary h-6 w-6" />
				</div>
				Middlewares
			</h1>
			<p class="text-muted-foreground mt-1">Configure your middlewares</p>
		</div>
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
