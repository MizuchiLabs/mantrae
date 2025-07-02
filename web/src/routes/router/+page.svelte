<script lang="ts">
	import { routerClient } from '$lib/api';
	import RouterModal from '$lib/components/modals/RouterModal.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import ColumnRule from '$lib/components/tables/ColumnRule.svelte';
	import ColumnText from '$lib/components/tables/ColumnText.svelte';
	import ColumnTls from '$lib/components/tables/ColumnTLS.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { Button } from '$lib/components/ui/button/index.js';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { RouterTCPTLSConfig, RouterTLSConfig } from '$lib/gen/zen/traefik-schemas';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { profile } from '$lib/stores/profile';
	import { ConnectError } from '@connectrpc/connect';
	import {
		Bot,
		CircleCheck,
		CircleSlash,
		Globe,
		LayoutGrid,
		Network,
		Pencil,
		Power,
		PowerOff,
		Route,
		Table,
		Trash,
		Waves
	} from '@lucide/svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	let item = $state({} as Router);
	let open = $state(false);

	// Data state
	let data = $state<Router[]>([]);
	let rowCount = $state<number>(0);
	let viewMode = $state<'table' | 'grid'>('table');

	const columns: ColumnDef<Router>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true,
			cell: ({ row }) => {
				return renderComponent(ColumnText, {
					label: row.getValue('name') as string,
					icon: row.original.agentId ? Bot : undefined,
					iconProps: { class: 'text-green-500', size: 20 },
					class: 'text-sm'
				});
			}
		},
		{
			header: 'Protocol',
			accessorKey: 'type',
			enableSorting: true,
			enableGlobalFilter: false,
			filterFn: (row, columnId, filterValue) => {
				const protocol = row.getValue(columnId) as RouterType;

				// Handle both enum value and display label filtering
				if (typeof filterValue === 'string') {
					const displayLabel = getProtocolLabel(protocol);
					return (
						displayLabel.toLowerCase().includes(filterValue.toLowerCase()) ||
						protocol.toString().toLowerCase().includes(filterValue.toLowerCase())
					);
				}

				// Direct enum comparison for badge clicking
				return protocol === filterValue;
			},
			cell: ({ row, column }) => {
				let protocol = row.getValue('type') as RouterType.HTTP | RouterType.TCP | RouterType.UDP;

				let label = 'Unspecified';
				let icon = undefined;
				if (protocol === RouterType.HTTP) {
					label = 'HTTP';
					icon = Globe;
				} else if (protocol === RouterType.TCP) {
					label = 'TCP';
					icon = Network;
				} else if (protocol === RouterType.UDP) {
					label = 'UDP';
					icon = Waves;
				}
				return renderComponent(ColumnBadge<Router>, {
					label,
					icon,
					variant: 'outline',
					class: 'hover:cursor-pointer',
					column: column
				});
			}
		},
		{
			header: 'EntryPoints',
			accessorKey: 'config.entryPoints',
			id: 'entrypoints',
			enableSorting: true,
			enableGlobalFilter: false,
			filterFn: 'arrIncludes',
			cell: ({ row, column }) => {
				let entrypoints: string[] = [];
				if (row.original.config?.entryPoints !== undefined) {
					entrypoints = row.getValue('entrypoints') as string[];
				}
				return renderComponent(ColumnBadge<Router>, {
					label: entrypoints?.length ? entrypoints : 'None',
					variant: entrypoints?.length ? 'secondary' : 'outline',
					column: entrypoints?.length ? column : undefined
				});
			}
		},
		{
			header: 'Middlewares',
			accessorKey: 'config.middlewares',
			id: 'middlewares',
			enableSorting: true,
			enableGlobalFilter: false,
			filterFn: 'arrIncludes',
			cell: ({ row, column }) => {
				let middlewares: string[] = [];
				if (row.original.config?.middlewares !== undefined) {
					middlewares = row.getValue('middlewares') as string[];
				}
				return renderComponent(ColumnBadge<Router>, {
					label: middlewares?.length ? middlewares : 'None',
					variant: middlewares?.length ? 'secondary' : 'outline',
					column: middlewares?.length ? column : undefined
				});
			}
		},
		{
			header: 'Rules',
			accessorKey: 'config.rule',
			id: 'rule',
			enableSorting: true,
			cell: ({ row }) => {
				let rule = '';
				if (row.original.config?.rule !== undefined) {
					rule = row.getValue('rule') as string;
				}
				return renderComponent(ColumnRule, {
					rule: rule,
					routerType: row.original.type as RouterType.HTTP | RouterType.TCP
				});
			}
		},
		{
			header: 'TLS',
			accessorKey: 'config.tls',
			id: 'tls',
			enableSorting: true,
			enableGlobalFilter: false,
			filterFn: (row, columnId, filterValue) => {
				const tls = row.getValue(columnId) as RouterTLSConfig;
				return tls?.certResolver === filterValue;
			},
			cell: ({ row, column }) => {
				let tls: RouterTLSConfig | RouterTCPTLSConfig = {};
				if (row.original.config?.tls !== undefined) {
					tls = row.getValue('tls') as RouterTLSConfig | RouterTCPTLSConfig;
				}
				return renderComponent(ColumnTls<Router>, { tls, column });
			}
		},
		{
			header: 'Enabled',
			accessorKey: 'enabled',
			enableSorting: true,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				return renderComponent(TableActions, {
					actions: [
						{
							type: 'button',
							label: row.original.enabled ? 'Disable' : 'Enable',
							icon: row.original.enabled ? Power : PowerOff,
							iconProps: {
								class: row.original.enabled ? 'text-green-500 size-5' : 'text-red-500 size-5'
							},
							onClick: () => toggleItem(row.original, !row.original.enabled)
						}
					]
				});
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
							label: 'Edit Router',
							icon: Pencil,
							onClick: () => {
								item = row.original;
								open = true;
							}
						},
						{
							type: 'button',
							label: 'Delete Router',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original)
						}
					]
				});
			}
		}
	];

	// Helper functions to avoid repetition
	function getProtocolLabel(protocol: RouterType): string {
		if (protocol === RouterType.HTTP) return 'HTTP';
		if (protocol === RouterType.TCP) return 'TCP';
		if (protocol === RouterType.UDP) return 'UDP';
		return 'Unspecified';
	}

	const bulkActions: BulkAction<Router>[] = [
		{
			type: 'button',
			label: 'Enable',
			icon: CircleCheck,
			variant: 'outline',
			onClick: (e) => bulk(e, 'enable')
		},
		{
			type: 'button',
			label: 'Disable',
			icon: CircleSlash,
			variant: 'outline',
			onClick: (e) => bulk(e, 'disable')
		},
		{
			type: 'button',
			label: 'Delete',
			icon: Trash,
			variant: 'destructive',
			onClick: (e) => bulk(e, 'delete')
		}
	];

	async function onPaginationChange(p: PaginationState) {
		await refreshData(p.pageSize, p.pageIndex);
	}

	const deleteItem = async (item: Router) => {
		try {
			await routerClient.deleteRouter({ id: item.id, type: item.type });
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`Router ${item.name} deleted`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
	};

	const toggleItem = async (item: Router, enabled: boolean) => {
		try {
			await routerClient.updateRouter({
				id: item.id,
				name: item.name,
				type: item.type,
				config: item.config,
				enabled: enabled,
				dnsProviders: item.dnsProviders
			});
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`Router ${item.name} ${enabled ? 'enabled' : 'disabled'}`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to update router', { description: e.message });
		}
	};

	async function bulk(rows: Router[], action: string) {
		try {
			const confirmed = confirm(`Are you sure you want to ${action} ${rows.length} routers?`);
			if (!confirmed) return;

			switch (action) {
				case 'delete':
					for (const row of rows) {
						await routerClient.deleteRouter({ id: row.id, type: row.type });
					}
					break;
				case 'disable':
					for (const row of rows) {
						await routerClient.updateRouter({
							id: row.id,
							name: row.name,
							type: row.type,
							config: row.config,
							dnsProviders: row.dnsProviders,
							enabled: false
						});
					}
					break;
				case 'enable':
					for (const row of rows) {
						await routerClient.updateRouter({
							id: row.id,
							name: row.name,
							type: row.type,
							config: row.config,
							dnsProviders: row.dnsProviders,
							enabled: true
						});
					}
					break;
			}

			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`Successfully ${action}d ${rows.length} routers`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error(`Failed to ${action}d routers`, { description: e.message });
		}
	}
	async function refreshData(pageSize: number, pageIndex: number) {
		const response = await routerClient.listRouters({
			profileId: profile.id,
			limit: BigInt(pageSize),
			offset: BigInt(pageIndex * pageSize)
		});
		data = response.routers;
		rowCount = Number(response.totalCount);
	}

	onMount(async () => {
		await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
	});
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="bg-primary/10 rounded-lg p-2">
					<Route class="text-primary h-6 w-6" />
				</div>
				Routers
			</h1>
			<p class="text-muted-foreground mt-1">Manage your routers and services</p>
		</div>

		<!-- View Toggle -->
		<div class="flex items-center gap-2">
			<Button
				variant={viewMode === 'table' ? 'default' : 'outline'}
				size="sm"
				onclick={() => (viewMode = 'table')}
			>
				<Table class="h-4 w-4" />
				Table
			</Button>
			<Button
				variant={viewMode === 'grid' ? 'default' : 'outline'}
				size="sm"
				onclick={() => (viewMode = 'grid')}
			>
				<LayoutGrid class="h-4 w-4" />
				Grid
			</Button>
		</div>
	</div>

	<DataTable
		{data}
		{columns}
		{rowCount}
		{viewMode}
		{onPaginationChange}
		{bulkActions}
		rowClassModifiers={{}}
		cardConfig={{
			titleKey: 'name',
			subtitleKey: 'type',
			excludeColumns: []
		}}
		createButton={{
			label: 'Create Router',
			onClick: () => {
				item = { type: RouterType.HTTP } as Router;
				open = true;
			}
		}}
	/>
</div>

<RouterModal bind:open bind:item bind:data />
