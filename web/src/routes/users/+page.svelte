<script lang="ts">
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { Pencil, Trash, Users } from '@lucide/svelte';
	import UserModal from '$lib/components/modals/user.svelte';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { DateFormat, pageIndex, pageSize } from '$lib/stores/common';
	import type { User } from '$lib/gen/mantrae/v1/user_pb';
	import { userClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { onMount } from 'svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import ColumnCheck from '$lib/components/tables/ColumnCheck.svelte';
	import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';

	let item = $state({} as User);
	let open = $state(false);

	// Data state
	let data = $state<User[]>([]);
	let rowCount = $state<number>(0);

	const columns: ColumnDef<User>[] = [
		{
			header: 'Username',
			accessorKey: 'username',
			enableSorting: true
		},
		{
			header: 'Email',
			accessorKey: 'email',
			enableSorting: true
		},
		{
			header: 'Admin',
			accessorKey: 'isAdmin',
			enableSorting: true,
			cell: ({ row }) => {
				const checked = row.getValue('isAdmin') as boolean;
				return renderComponent(ColumnCheck, { checked });
			}
		},
		{
			header: 'Last Login',
			accessorKey: 'lastLogin',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('lastLogin') as Timestamp;
				return DateFormat.format(timestampDate(date));
			}
		},
		{
			header: 'Created',
			accessorKey: 'createdAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('createdAt') as Timestamp;
				return DateFormat.format(timestampDate(date));
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
							label: 'Edit User',
							icon: Pencil,
							onClick: () => {
								item = row.original;
								open = true;
							}
						},
						{
							type: 'button',
							label: 'Delete User',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original.id)
						}
					]
				});
			}
		}
	];

	const bulkActions: BulkAction<User>[] = [
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

	const deleteItem = async (id: string) => {
		try {
			await userClient.deleteUser({ id: id });
			await refreshData(pageSize.value ?? 10, 0);
			toast.success('Router deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
	};

	async function bulkDelete(selectedRows: User[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${selectedRows.length} Users?`);
			if (!confirmed) return;

			const ids = selectedRows.map((row) => ({ id: row.id }));
			for (const row of ids) {
				await userClient.deleteUser({ id: row.id });
			}
			await refreshData(pageSize.value ?? 10, 0);
			toast.success(`Successfully deleted ${selectedRows.length} Users`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete DNS Providers', { description: e.message });
		}
	}

	async function refreshData(pageSize: number, pageIndex: number) {
		const response = await userClient.listUsers({
			limit: BigInt(pageSize),
			offset: BigInt(pageIndex * pageSize)
		});
		data = response.users;
		rowCount = Number(response.totalCount);
	}

	onMount(async () => {
		await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
	});
</script>

<svelte:head>
	<title>Users</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<Users />
		<h1 class="text-2xl font-bold">User Management</h1>
	</div>
	<DataTable
		{data}
		{columns}
		{rowCount}
		{onPaginationChange}
		{bulkActions}
		createButton={{
			label: 'Add User',
			onClick: () => {
				item = {} as User;
				open = true;
			}
		}}
	/>
</div>

<UserModal bind:open bind:item bind:data />
