<script lang="ts">
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { CircleCheck, CircleSlash, Pencil, Trash, Users } from '@lucide/svelte';
	import UserModal from '$lib/components/modals/UserModal.svelte';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { DateFormat, pageIndex, pageSize } from '$lib/stores/common';
	import type { User } from '$lib/gen/mantrae/v1/user_pb';
	import { userClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { onMount } from 'svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
	import { user } from '$lib/stores/user';

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
			enableGlobalFilter: false,
			cell: ({ row }) => {
				return renderComponent(TableActions, {
					actions: [
						{
							type: 'button',
							label: row.original.isAdmin ? 'Disable' : 'Enable',
							icon: row.original.isAdmin ? CircleCheck : CircleSlash,
							iconProps: {
								class: row.original.isAdmin ? 'text-green-500 size-5' : 'text-red-500 size-5',
								size: 20
							},
							onClick: () => toggleItem(row.original, !row.original.isAdmin)
						}
					]
				});
			}
		},
		{
			header: 'Last Login',
			accessorKey: 'lastLogin',
			enableSorting: true,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				const date = row.getValue('lastLogin') as Timestamp;
				return DateFormat.format(timestampDate(date));
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
							onClick: () => deleteItem(row.original)
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

	const deleteItem = async (item: User) => {
		try {
			await userClient.deleteUser({ id: item.id });
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`User ${item.username} deleted`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete user', { description: e.message });
		}
	};

	async function toggleItem(item: User, isAdmin: boolean) {
		try {
			if (user.id === item.id) {
				toast.error('You cannot change your own role!');
				return;
			}
			await userClient.updateUser({
				id: item.id,
				username: item.username,
				email: item.email,
				isAdmin: isAdmin
			});
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`User ${item.username} ${isAdmin ? 'set as admin' : 'removed as admin'}`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to update user', { description: e.message });
		}
	}

	async function bulkDelete(rows: User[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} Users?`);
			if (!confirmed) return;

			for (const row of rows) {
				await userClient.deleteUser({ id: row.id });
			}
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`Successfully deleted ${rows.length} Users`);
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

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="bg-primary/10 rounded-lg p-2">
					<Users class="text-primary h-6 w-6" />
				</div>
				User Management
			</h1>
			<p class="text-muted-foreground mt-1">Manage your users and access management</p>
		</div>
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
