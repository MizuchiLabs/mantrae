<script lang="ts">
	import { dnsClient } from '$lib/api';
	import DNSModal from '$lib/components/modals/DNSModal.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import ColumnCheck from '$lib/components/tables/ColumnCheck.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { DnsProviderType, type DnsProvider } from '$lib/gen/mantrae/v1/dns_provider_pb';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { ConnectError } from '@connectrpc/connect';
	import { CircleCheck, CircleSlash, Globe, Pencil, Trash } from '@lucide/svelte';
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
			enableSorting: true,
			enableHiding: false
		},
		{
			header: 'Provider',
			accessorKey: 'type',
			enableSorting: true,
			enableGlobalFilter: false,
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
			header: 'IP Address',
			accessorKey: 'config.ip',
			id: 'ip',
			enableSorting: true,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				if (row.original.config?.autoUpdate) {
					// utilClient.getPublicIP({}).then((res) => {
					return renderComponent(ColumnBadge, {
						label: 'auto',
						variant: 'secondary',
						class: 'hover:cursor-pointer'
					});
					// });
				} else {
					return renderComponent(ColumnBadge, {
						label: row.getValue('ip') as string,
						class: 'hover:cursor-pointer'
					});
				}
			}
		},
		{
			header: 'Default',
			accessorKey: 'isDefault',
			enableSorting: true,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				return renderComponent(TableActions, {
					actions: [
						{
							type: 'button',
							label: row.original.isDefault ? 'Disable' : 'Enable',
							icon: row.original.isDefault ? CircleCheck : CircleSlash,
							iconProps: {
								class: row.original.isDefault ? 'text-green-500 size-5' : 'text-red-500 size-5',
								size: 20
							},
							onClick: () => toggleItem(row.original, !row.original.isDefault)
						}
					]
				});
			}
		},
		{
			header: 'Proxied',
			accessorKey: 'config.proxied',
			id: 'proxied',
			enableSorting: true,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				let checked = row.getValue('proxied') as boolean;
				return renderComponent(ColumnCheck, { checked: checked });
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
							onClick: () => deleteItem(row.original)
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

	const deleteItem = async (item: DnsProvider) => {
		try {
			await dnsClient.deleteDnsProvider({ id: item.id });
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`DNS Provider ${item.name} deleted`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete DNS Provider', { description: e.message });
		}
	};

	async function toggleItem(item: DnsProvider, isDefault: boolean) {
		try {
			await dnsClient.updateDnsProvider({
				id: item.id,
				name: item.name,
				type: item.type,
				config: item.config,
				isDefault: isDefault
			});
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(
				`DNS Provider ${item.name} ${isDefault ? 'set as default' : 'removed as default'}`
			);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to update DNS Provider', { description: e.message });
		}
	}

	async function bulkDelete(rows: DnsProvider[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} DNS Providers?`);
			if (!confirmed) return;

			for (const row of rows) {
				await dnsClient.deleteDnsProvider({ id: row.id });
			}
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`Successfully deleted ${rows.length} DNS Providers`);
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

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="bg-primary/10 rounded-lg p-2">
					<Globe class="text-primary h-6 w-6" />
				</div>
				DNS Management
			</h1>
			<p class="text-muted-foreground mt-1">Manage your DNS providers</p>
		</div>
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
