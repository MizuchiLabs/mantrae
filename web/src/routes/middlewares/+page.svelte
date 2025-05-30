<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import MiddlewareModal from '$lib/components/modals/middleware.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import type { DeleteMiddlewareParams, Middleware } from '$lib/types/middlewares';
	import { Layers, Pencil, Trash } from '@lucide/svelte';
	import { TraefikSource } from '$lib/types';
	import { api, middlewares } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { source } from '$lib/stores/source';
	import { profile } from '$lib/stores/profile';
	import type { BulkAction } from '$lib/components/tables/types';

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
			if (!profile.hasValidId() || !profile.id) {
				toast.error('Invalid profile ID');
				return;
			}
			const params: DeleteMiddlewareParams = {
				profileId: profile?.id,
				name: middleware.name,
				protocol: middleware.protocol
			};
			await api.deleteMiddleware(params);
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

			const items = selectedRows.map((row) => ({
				name: row.name,
				protocol: row.protocol
			}));
			await api.bulkDeleteMiddleware(items);
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
				return name?.split('@')[0];
			}
		},
		{
			header: 'Protocol',
			accessorKey: 'protocol',
			enableSorting: true,
			cell: ({ row, column }) => {
				return renderComponent(ColumnBadge<Middleware>, {
					label: row.getValue('protocol') as string,
					class: 'hover:cursor-pointer',
					column: column
				});
			}
		},
		{
			header: 'Type',
			accessorKey: 'type',
			enableSorting: true,
			cell: ({ row, column }) => {
				return renderComponent(ColumnBadge<Middleware>, {
					label: row.getValue('type') as string,
					class: 'hover:cursor-pointer',
					column: column
				});
			}
		},
		{
			header: 'Provider',
			accessorFn: (row) => row.name,
			id: 'provider',
			enableSorting: true,
			cell: ({ row, column }) => {
				const name = row.getValue('provider') as string;
				const provider = name?.split('@')[1];
				return renderComponent(ColumnBadge<Middleware>, {
					label: provider ? provider.toLowerCase() : 'http',
					variant: 'secondary',
					class: 'hover:cursor-pointer',
					column: column
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
			type: 'button',
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
