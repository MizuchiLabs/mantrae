<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable, { type BulkAction } from '$lib/components/tables/DataTable.svelte';
	import MiddlewareModal from '$lib/components/modals/middleware.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import type { Middleware, SupportedMiddleware } from '$lib/types/middlewares';
	import { Layers, Pencil, Trash } from 'lucide-svelte';
	import { TraefikSource } from '$lib/types';
	import { api, middlewares } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { source } from '$lib/stores/source';
	import { profile } from '$lib/stores/profile';

	interface ModalState {
		isOpen: boolean;
		middleware: Middleware;
		mode?: 'create' | 'edit';
	}
	const initialModalState: ModalState = {
		isOpen: false,
		mode: 'create',
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
			mode: 'create',
			middleware: initialModalState.middleware
		};
	}

	const deleteMiddleware = async (middleware: Middleware) => {
		if (!source.isLocal()) return;
		try {
			await api.deleteMiddleware(middleware);
			toast.success('Middleware deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

	async function handleBulkDelete(selectedRows: Middleware[]) {
		try {
			const confirmed = confirm(
				`Are you sure you want to delete ${selectedRows.length} middlewares?`
			);
			if (!confirmed) return;

			await Promise.all(selectedRows.map((row) => api.deleteMiddleware(row)));
			toast.success(`Successfully deleted ${selectedRows.length} middlewares`);
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(`Failed to delete middlewares: ${e.message}`);
		}
	}

	const defaultColumns: ColumnDef<Middleware>[] = [
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
					label: name.split('@')[1]?.toLowerCase(),
					variant: 'secondary'
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
								modalState = {
									isOpen: true,
									mode: 'edit',
									middleware: row.original
								};
							}
						},
						{
							type: 'button',
							label: 'Delete Middleware',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => {
								deleteMiddleware(row.original);
							},
							disabled: source.value !== TraefikSource.LOCAL
						}
					],
					shareObject: source.value === TraefikSource.LOCAL ? row.original : undefined
				});
			}
		}
	];

	const mwBulkActions: BulkAction<Middleware>[] = [
		{
			label: 'Delete',
			icon: Trash,
			variant: 'destructive',
			onClick: handleBulkDelete
		}
	];

	let columns: ColumnDef<Middleware>[] = $derived(
		source.value === TraefikSource.LOCAL
			? defaultColumns.filter((c) => c.id !== 'provider')
			: defaultColumns
	);

	$effect(() => {
		if (profile.isValid() && source.value) {
			api.getTraefikConfig(source.value);
		}
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
	{#if source.value === TraefikSource.LOCAL}
		<DataTable
			{columns}
			data={$middlewares || []}
			showSourceTabs={true}
			createButton={{
				label: 'Add Middleware',
				onClick: openCreateModal
			}}
			bulkActions={mwBulkActions}
		/>
	{:else}
		<DataTable {columns} data={$middlewares || []} showSourceTabs={true} />
	{/if}
</div>

<MiddlewareModal
	bind:open={modalState.isOpen}
	bind:middleware={modalState.middleware}
	mode={modalState.mode}
/>
