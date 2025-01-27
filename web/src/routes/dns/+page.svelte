<script lang="ts">
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import DNSModal from '$lib/components/modals/dns.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { Globe, Pencil, Trash } from 'lucide-svelte';
	import { type DNSProvider } from '$lib/types';
	import { api, dnsProviders } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { onMount } from 'svelte';
	import { DateFormat } from '$lib/store';

	interface ModalState {
		isOpen: boolean;
		dnsProvider?: DNSProvider;
	}

	const initialModalState: ModalState = { isOpen: false };
	let modalState = $state(initialModalState);

	const deleteUser = async (d: DNSProvider) => {
		try {
			await api.deleteDNSProvider(d.id);
			toast.success('DNSProvider deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

	const columns: ColumnDef<DNSProvider>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true
		},
		{
			header: 'Type',
			accessorKey: 'type',
			enableSorting: true
		},
		{
			header: 'Default Provider',
			accessorKey: 'isActive',
			enableSorting: true,
			cell: ({ row }) => {
				const active = row.getValue('isActive') as boolean;
				if (active) {
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
			header: 'Created',
			accessorKey: 'createdAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('createdAt') as string;
				return DateFormat.format(new Date(date));
			}
		},
		{
			header: 'Updated',
			accessorKey: 'updatedAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('updatedAt') as string;
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
							label: 'Edit DNSProvider',
							icon: Pencil,
							onClick: () => {
								modalState = {
									isOpen: true,
									dnsProvider: row.original
								};
							}
						},
						{
							label: 'Delete DNSProvider',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => {
								deleteUser(row.original);
							}
						}
					]
				});
			}
		}
	];

	onMount(async () => {
		await api.listDNSProviders();
	});
</script>

<svelte:head>
	<title>Users</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<Globe />
		<h1 class="text-2xl font-bold">DNS Management</h1>
	</div>
	<DataTable
		{columns}
		data={$dnsProviders || []}
		createButton={{
			label: 'Add Provider',
			onClick: () => {
				modalState = { isOpen: true };
			}
		}}
	/>
</div>

<DNSModal bind:open={modalState.isOpen} dns={modalState.dnsProvider} />
