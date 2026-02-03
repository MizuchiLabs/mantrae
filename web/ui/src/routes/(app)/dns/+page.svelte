<script lang="ts">
	import { dns } from '$lib/api/dns.svelte';
	import DNSModal from '$lib/components/modals/DNSModal.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import ColumnCheck from '$lib/components/tables/ColumnCheck.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { DNSProviderType, type DNSProvider } from '$lib/gen/mantrae/v1/dns_provider_pb';
	import { ConnectError } from '@connectrpc/connect';
	import { CircleCheck, CircleSlash, Globe, Pencil, Trash } from '@lucide/svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { toast } from 'svelte-sonner';

	let data = $state({} as DNSProvider);
	let open = $state(false);

	const dnsList = dns.list();
	const updateDNS = dns.update();
	const deleteDNS = dns.delete();

	const columns: ColumnDef<DNSProvider>[] = [
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
					| DNSProviderType.DNS_PROVIDER_TYPE_CLOUDFLARE
					| DNSProviderType.DNS_PROVIDER_TYPE_POWERDNS
					| DNSProviderType.DNS_PROVIDER_TYPE_TECHNITIUM
					| DNSProviderType.DNS_PROVIDER_TYPE_PIHOLE;
				let label = 'Unspecified';
				switch (type) {
					case DNSProviderType.DNS_PROVIDER_TYPE_CLOUDFLARE:
						label = 'Cloudflare';
						break;
					case DNSProviderType.DNS_PROVIDER_TYPE_POWERDNS:
						label = 'PowerDNS';
						break;
					case DNSProviderType.DNS_PROVIDER_TYPE_TECHNITIUM:
						label = 'Technitium';
						break;
					case DNSProviderType.DNS_PROVIDER_TYPE_PIHOLE:
						label = 'PiHole';
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
					return renderComponent(ColumnBadge, {
						label: 'auto',
						variant: 'secondary',
						class: 'hover:cursor-pointer'
					});
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
							onClick: () =>
								updateDNS.mutate({ ...row.original, isDefault: !row.original.isDefault })
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
								data = row.original;
								open = true;
							}
						},
						{
							type: 'popover',
							label: 'Delete Provider',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteDNS.mutate({ id: row.original.id }),
							popover: {
								title: 'Delete Provider?',
								description: 'This DNS provider will be permanently deleted.',
								confirmLabel: 'Delete',
								cancelLabel: 'Cancel'
							}
						}
					]
				});
			}
		}
	];

	const bulkActions: BulkAction<DNSProvider>[] = [
		{
			type: 'button',
			label: 'Delete',
			icon: Trash,
			variant: 'destructive',
			onClick: bulkDelete
		}
	];

	async function bulkDelete(rows: DNSProvider[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} DNS Providers?`);
			if (!confirmed) return;

			for (const row of rows) {
				deleteDNS.mutate({ id: row.id });
			}
			toast.success(`Successfully deleted ${rows.length} DNS Providers`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete DNS Providers', { description: e.message });
		}
	}
</script>

<svelte:head>
	<title>DNS Providers - Mantrae</title>
	<meta
		name="description"
		content="Manage your DNS providers for automatic DNS challenge resolution with Let's Encrypt"
	/>
</svelte:head>

<DNSModal bind:open {data} />

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="rounded-lg bg-primary/10 p-2">
					<Globe class="h-6 w-6 text-primary" />
				</div>
				DNS Management
			</h1>
			<p class="mt-1 text-muted-foreground">Manage your DNS providers</p>
		</div>
	</div>

	<DataTable
		data={dnsList.data}
		{columns}
		{bulkActions}
		createButton={{
			label: 'Add Provider',
			onClick: () => (open = true)
		}}
	/>
</div>
