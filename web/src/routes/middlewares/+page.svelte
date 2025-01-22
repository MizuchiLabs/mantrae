<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import MiddlewareModal from '$lib/components/modals/middleware.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import type { Middleware } from '$lib/types/middlewares';
	import { Edit, Layers, Trash } from 'lucide-svelte';
	import { TraefikSource } from '$lib/types';
	import { api, profile, middlewares, source } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';

	interface ModalState {
		isOpen: boolean;
		mode: 'create' | 'edit';
		middleware?: Middleware;
	}

	const initialModalState: ModalState = {
		isOpen: false,
		mode: 'create'
	};

	let modalState = $state(initialModalState);

	function openCreateModal() {
		modalState = {
			isOpen: true,
			mode: 'create'
		};
	}

	function openEditModal(middleware: Middleware) {
		modalState = {
			isOpen: true,
			mode: 'edit',
			middleware
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
				const type = row.getValue('type') as string;
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
			cell: ({ row }) => {
				if ($source === TraefikSource.LOCAL) {
					return renderComponent(TableActions, {
						actions: [
							{
								label: 'Edit Middleware',
								icon: Edit,
								onClick: () => {
									openEditModal(row.original);
								}
							},
							{
								label: 'Delete Middleware',
								icon: Trash,
								variant: 'destructive',
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
								label: 'Edit Middleware',
								icon: Edit,
								onClick: () => {
									openEditModal(row.original);
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
	$effect(() => {
		console.log($middlewares);
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

			<DataTable
				{columns}
				data={$middlewares || []}
				createButton={{
					label: 'Add Middleware',
					onClick: openCreateModal
				}}
			/>
		</div>
	</Tabs.Content>
</Tabs.Root>

<MiddlewareModal
	bind:open={modalState.isOpen}
	mode={modalState.mode}
	middleware={modalState.middleware}
/>
