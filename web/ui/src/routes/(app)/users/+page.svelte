<script lang="ts">
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { Pencil, Trash, Users } from '@lucide/svelte';
	import UserModal from '$lib/components/modals/UserModal.svelte';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import type { User } from '$lib/gen/mantrae/v1/user_pb';
	import { ConnectError } from '@connectrpc/connect';
	import type { BulkAction } from '$lib/components/tables/types';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import { formatTs } from '$lib/utils';
	import { user } from '$lib/api/users.svelte';

	let data = $state({} as User);
	let open = $state(false);

	const userList = user.list();
	const deleteUser = user.delete();

	const columns: ColumnDef<User>[] = [
		{
			header: 'Username',
			accessorKey: 'username',
			enableSorting: true
		},
		{
			header: 'Email',
			accessorKey: 'email',
			enableSorting: true,
			cell: ({ row }) => {
				return renderComponent(ColumnBadge, { label: row.original.email || 'None' });
			}
		},
		{
			header: 'Last Login',
			accessorKey: 'lastLogin',
			enableSorting: true,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				if (row.original.lastLogin === undefined) {
					return renderComponent(ColumnBadge, { label: 'Never' });
				}
				return formatTs(row.original.lastLogin, 'relative');
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
								data = row.original;
								open = true;
							}
						},
						{
							type: 'popover',
							label: 'Delete User',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteUser.mutate({ id: row.original.id }),
							popover: {
								title: 'Delete User?',
								description: 'This user will be permanently deleted.',
								confirmLabel: 'Delete',
								cancelLabel: 'Cancel'
							}
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

	async function bulkDelete(rows: User[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} Users?`);
			if (!confirmed) return;

			for (const row of rows) {
				deleteUser.mutate({ id: row.id });
			}
			toast.success(`Successfully deleted ${rows.length} Users`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete DNS Providers', { description: e.message });
		}
	}
</script>

<svelte:head>
	<title>User Management</title>
	<meta
		name="description"
		content="Manage your Mantrae users and access permissions for your reverse proxy system"
	/>
</svelte:head>

<UserModal bind:open {data} />

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="rounded-lg bg-primary/10 p-2">
					<Users class="h-6 w-6 text-primary" />
				</div>
				User Management
			</h1>
			<p class="mt-1 text-muted-foreground">Manage your users and access management</p>
		</div>
	</div>

	<DataTable
		data={userList.data}
		{columns}
		{bulkActions}
		createButton={{
			label: 'Add User',
			onClick: () => (open = true)
		}}
	/>
</div>
