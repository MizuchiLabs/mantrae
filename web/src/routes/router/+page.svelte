<script lang="ts">
	import { routerClient } from '$lib/api';
	import RouterModal from '$lib/components/modals/router.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import ColumnCheck from '$lib/components/tables/ColumnCheck.svelte';
	import ColumnRule from '$lib/components/tables/ColumnRule.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { RouterType, type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { RouterTLSConfig } from '$lib/gen/zen/traefik-schemas';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { profile } from '$lib/stores/profile';
	import { ConnectError } from '@connectrpc/connect';
	import { Bot, CircleCheck, CircleSlash, Pencil, Route, Trash } from '@lucide/svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	let item = $state({} as Router);
	let open = $state(false);

	// Data state
	let data = $state<Router[]>([]);
	let rowCount = $state<number>(0);

	const columns: ColumnDef<Router>[] = [
		{
			header: 'Name',
			accessorKey: 'name',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('name') as string;
				if (row.original.agentId) {
					return renderComponent(ColumnBadge<Router>, {
						label: name?.split('@')[0],
						variant: 'outline',
						icon: Bot
					});
				}
				return renderComponent(ColumnBadge<Router>, {
					label: name?.split('@')[0],
					variant: 'outline'
				});
			}
		},
		{
			header: 'Protocol',
			accessorKey: 'type',
			enableSorting: true,
			cell: ({ row, column }) => {
				let protocol = row.getValue('type') as RouterType.HTTP | RouterType.TCP | RouterType.UDP;

				let label = 'Unspecified';
				if (protocol === RouterType.HTTP) {
					label = 'HTTP';
				} else if (protocol === RouterType.TCP) {
					label = 'TCP';
				} else if (protocol === RouterType.UDP) {
					label = 'UDP';
				}
				return renderComponent(ColumnBadge<Router>, {
					label: label,
					class: 'hover:cursor-pointer italic',
					column: column
				});
			}
		},
		{
			header: 'EntryPoints',
			accessorKey: 'config.entryPoints',
			id: 'entrypoints',
			enableSorting: true,
			filterFn: 'arrIncludes',
			cell: ({ row, column }) => {
				let entrypoints: string[] = [];
				if (row.original.config?.entryPoints !== undefined) {
					entrypoints = row.getValue('entrypoints') as string[];
				}
				return renderComponent(ColumnBadge<Router>, {
					label: entrypoints?.length ? entrypoints : 'None',
					variant: entrypoints?.length ? 'secondary' : 'outline',
					class: 'hover:cursor-pointer',
					column: entrypoints?.length ? column : undefined
				});
			}
		},
		{
			header: 'Middlewares',
			accessorKey: 'config.middlewares',
			id: 'middlewares',
			enableSorting: true,
			filterFn: 'arrIncludes',
			cell: ({ row, column }) => {
				let middlewares: string[] = [];
				if (row.original.config?.middlewares !== undefined) {
					middlewares = row.getValue('middlewares') as string[];
				}
				return renderComponent(ColumnBadge<Router>, {
					label: middlewares?.length ? middlewares : 'None',
					variant: middlewares?.length ? 'secondary' : 'outline',
					class: 'hover:cursor-pointer',
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
			filterFn: (row, columnId, filterValue) => {
				const tls = row.getValue(columnId) as RouterTLSConfig;
				return tls?.certResolver === filterValue;
			},
			cell: ({ row, column }) => {
				let tls = undefined;
				if (row.original.config?.tls !== undefined) {
					tls = row.getValue('tls') as RouterTLSConfig;
				}

				return renderComponent(ColumnBadge<Router>, {
					label: tls?.certResolver ? tls.certResolver : 'Disabled',
					variant: tls?.certResolver ? 'secondary' : 'outline',
					class: tls?.certResolver ? 'bg-slate-300 dark:bg-slate-700' : '',
					tls,
					column: tls?.certResolver ? column : undefined
				});
			}
		},
		{
			header: 'Enabled',
			accessorKey: 'enabled',
			enableSorting: true,
			cell: ({ row }) => {
				let checked = row.getValue('enabled') as boolean;
				return renderComponent(ColumnCheck, { checked: checked });
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
							onClick: () => deleteItem(row.original.id, row.original.type)
						}
					]
				});
			}
		}
	];

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

	const deleteItem = async (id: bigint, type: RouterType) => {
		try {
			await routerClient.deleteRouter({ id: id, type: type });
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success('Router deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
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

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<Route />
		<h1 class="text-2xl font-bold">Routers</h1>
	</div>
	<DataTable
		{data}
		{columns}
		{rowCount}
		{onPaginationChange}
		{bulkActions}
		rowClassModifiers={{
			'bg-red-300/25 dark:bg-red-700/25': (r) => !r.enabled,
			'bg-green-300/25 dark:bg-green-700/25': (r) => r.agentId !== ''
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
