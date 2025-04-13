<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { Pencil, Trash, Users } from 'lucide-svelte';
	import { type User } from '$lib/types';
	import UserModal from '$lib/components/modals/user.svelte';
	import { api, users } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { onMount } from 'svelte';
	import { DateFormat } from '$lib/stores/common';

	interface ModalState {
		isOpen: boolean;
		user?: User;
	}

	const initialModalState: ModalState = { isOpen: false };
	let modalState = $state(initialModalState);

	const deleteUser = async (user: User) => {
		try {
			await api.deleteUser(user.id);
			toast.success('User deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

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
				const admin = row.getValue('isAdmin') as boolean;
				if (admin) {
					return renderComponent(ColumnBadge, {
						label: 'Yes',
						variant: 'default'
					});
				} else {
					return renderComponent(ColumnBadge, {
						label: 'No',
						variant: 'secondary'
					});
				}
			}
		},
		{
			header: 'Last Login',
			accessorKey: 'lastLogin',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('lastLogin') as string;
				return DateFormat.format(new Date(date));
			}
		},
		{
			header: 'Created',
			accessorKey: 'createdAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('createdAt') as string;
				return DateFormat.format(new Date(date));
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
								modalState = {
									isOpen: true,
									user: row.original
								};
							}
						},
						{
							type: 'button',
							label: 'Delete User',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => {
								deleteUser(row.original);
							},
							disabled: row.original.id === 1
						}
					]
				});
			}
		}
	];

	onMount(async () => {
		await api.listUsers();
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
		{columns}
		data={$users || []}
		createButton={{
			label: 'Add User',
			onClick: () => {
				modalState = { isOpen: true };
			}
		}}
	/>
</div>

<UserModal bind:open={modalState.isOpen} user={modalState.user} />
