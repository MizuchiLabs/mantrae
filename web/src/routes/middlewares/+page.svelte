<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import MiddlewareModal from '$lib/components/modals/middleware.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import type { Middleware } from '$lib/types/middlewares';
	import { Eye, Layers, Pencil, Trash } from 'lucide-svelte';
	import { TraefikSource } from '$lib/types';
	import { api, profile, middlewares, source } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import type { SupportedMiddleware } from '$lib/components/forms/mw_registry';

	interface ModalState {
		isOpen: boolean;
		middleware: Middleware;
	}

	const initialModalState: ModalState = {
		isOpen: false,
		middleware: {
			name: '',
			protocol: 'http',
			type: undefined
		}
	};

	let modalState = $state(initialModalState);

	function openCreateModal() {
		modalState = {
			isOpen: true,
			middleware: initialModalState.middleware
		};
	}

	const deleteMiddleware = async (middleware: Middleware) => {
		try {
			let provider = middleware.name.split('@')[1];
			if (provider !== 'http') {
				toast.error('Middleware not managed by Mantrae!');
				return;
			}

			await api.deleteMiddleware($profile.id, middleware);
			toast.success('Middleware deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

	const columns: ColumnDef<Middleware>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('name') as string;
				if (!name) return;
				return name.split('@')[0];
			}
		},
		{
			header: 'Protocol',
			accessorKey: 'protocol',
			enableSorting: true,
			cell: ({ row }) => {
				const protocol = row.getValue('protocol') as string;
				if (!protocol) return;
				return renderComponent(ColumnBadge, { label: protocol });
			}
		},
		{
			header: 'Type',
			accessorKey: 'type',
			enableSorting: true,
			cell: ({ row }) => {
				const type = row.getValue('type') as SupportedMiddleware;
				if (!type) return;
				return renderComponent(ColumnBadge, { label: type });
			}
		},
		{
			header: 'Provider',
			accessorFn: (row) => row.name,
			id: 'provider',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('provider') as string;
				if (!name) return;
				return renderComponent(ColumnBadge, {
					label: name.split('@')[1].toLowerCase(),
					variant: 'secondary'
				});
			}
		},
		{
			id: 'actions',
			enableHiding: false,
			cell: ({ row }) => {
				if ($source === TraefikSource.LOCAL) {
					return renderComponent(TableActions, {
						actions: [
							{
								label: 'Edit Middleware',
								icon: Pencil,
								onClick: () => {
									console.log(row.original);
									modalState = {
										isOpen: true,
										middleware: row.original
									};
								}
							},
							{
								label: 'Delete Middleware',
								icon: Trash,
								classProps: 'text-destructive',
								onClick: () => {
									deleteMiddleware(row.original);
								}
							}
						]
					});
				} else {
					return renderComponent(TableActions, {
						actions: [
							{
								label: 'View Middleware',
								icon: Eye,
								onClick: () => {
									modalState = {
										isOpen: true,
										middleware: row.original
									};
								}
							}
						]
					});
				}
			}
		}
	];

	profile.subscribe((value) => {
		if (value.id) {
			api.getTraefikConfig(value.id, $source);
		}
	});
</script>

<svelte:head>
	<title>Middlewares</title>
</svelte:head>

<Tabs.Root value={$source}>
	<Tabs.Content value={TraefikSource.LOCAL}>
		<div class="flex flex-col gap-4">
			<div class="flex items-center justify-start gap-2">
				<Layers />
				<h1 class="text-2xl font-bold">Middleware Management</h1>
			</div>
			<DataTable
				{columns}
				data={$middlewares || []}
				showSourceTabs={true}
				createButton={{
					label: 'Add Middleware',
					onClick: openCreateModal
				}}
			/>
		</div>
	</Tabs.Content>
	<Tabs.Content value={TraefikSource.API}>
		<div class="flex flex-col gap-4">
			<div class="flex items-center justify-start gap-2">
				<Layers />
				<h1 class="text-2xl font-bold">Middleware Management</h1>
			</div>

			<DataTable {columns} data={$middlewares || []} showSourceTabs={true} />
		</div>
	</Tabs.Content>
	<Tabs.Content value={TraefikSource.AGENT}>
		<div class="flex flex-col gap-4">
			<div class="flex items-center justify-start gap-2">
				<Layers />
				<h1 class="text-2xl font-bold">Middleware Management</h1>
			</div>

			<DataTable {columns} data={$middlewares || []} showSourceTabs={true} />
		</div>
	</Tabs.Content>
</Tabs.Root>

<MiddlewareModal bind:open={modalState.isOpen} bind:middleware={modalState.middleware} />
