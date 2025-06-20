<script lang="ts">
	import { dnsClient } from '$lib/api';
	import DNSModal from '$lib/components/modals/dns.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import ColumnCheck from '$lib/components/tables/ColumnCheck.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { DnsProviderType, type DnsProvider } from '$lib/gen/mantrae/v1/dns_provider_pb';
	import { DateFormat, pageIndex, pageSize } from '$lib/stores/common';
	import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
	import { ConnectError } from '@connectrpc/connect';
	import { Globe, Pencil, Trash } from '@lucide/svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	let item = $state({} as DnsProvider);
	let open = $state(false);

	// Data state
	let data = $state<DnsProvider[]>([]);
	let rowCount = $state<number>(0);

	const columns: ColumnDef<DnsProvider>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true
		},
		{
			header: 'Provider',
			accessorKey: 'type',
			enableSorting: true,
			cell: ({ row }) => {
				let type = row.getValue('type') as
					| DnsProviderType.CLOUDFLARE
					| DnsProviderType.POWERDNS
					| DnsProviderType.TECHNITIUM;
				let label = 'Unspecified';
				switch (type) {
					case DnsProviderType.CLOUDFLARE:
						label = 'Cloudflare';
						break;
					case DnsProviderType.POWERDNS:
						label = 'PowerDNS';
						break;
					case DnsProviderType.TECHNITIUM:
						label = 'Technitium';
						break;
				}
				return renderComponent(ColumnBadge, {
					label: label,
					class: 'hover:cursor-pointer'
				});
			}
		},
		{
			header: 'Default',
			accessorKey: 'isActive',
			enableSorting: true,
			cell: ({ row }) => {
				let checked = row.getValue('isActive') as boolean;
				return renderComponent(ColumnCheck, { checked: checked });
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
			header: 'Updated',
			accessorKey: 'updatedAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('updatedAt') as Timestamp;
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
							label: 'Edit Provider',
							icon: Pencil,
							onClick: () => {
								item = row.original;
								open = true;
							}
						},
						{
							type: 'button',
							label: 'Delete Provider',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original.id)
						}
					]
				});
			}
		}
	];

	const bulkActions: BulkAction<DnsProvider>[] = [
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

	const deleteItem = async (id: bigint) => {
		try {
			await dnsClient.deleteDnsProvider({ id: id });
			await refreshData(pageSize.value ?? 10, 0);
			toast.success('Router deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
	};

	async function bulkDelete(selectedRows: DnsProvider[]) {
		try {
			const confirmed = confirm(
				`Are you sure you want to delete ${selectedRows.length} DNS Providers?`
			);
			if (!confirmed) return;

			const ids = selectedRows.map((row) => ({ id: row.id }));
			for (const row of ids) {
				await dnsClient.deleteDnsProvider({ id: row.id });
			}
			await refreshData(pageSize.value ?? 10, 0);
			toast.success(`Successfully deleted ${selectedRows.length} DNS Providers`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete DNS Providers', { description: e.message });
		}
	}

	async function refreshData(pageSize: number, pageIndex: number) {
		const response = await dnsClient.listDnsProviders({
			limit: BigInt(pageSize),
			offset: BigInt(pageIndex * pageSize)
		});
		data = response.dnsProviders;
		rowCount = Number(response.totalCount);
	}

	onMount(async () => {
		await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
	});
</script>

<svelte:head>
	<title>DNS Providers</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<Globe />
		<h1 class="text-2xl font-bold">DNS Management</h1>
	</div>
	<DataTable
		{data}
		{columns}
		{rowCount}
		{onPaginationChange}
		{bulkActions}
		createButton={{
			label: 'Add Provider',
			onClick: () => {
				item = {} as DnsProvider;
				open = true;
			}
		}}
	/>
</div>

<DNSModal bind:open bind:item bind:data />
