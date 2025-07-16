<script lang="ts">
	import { routerClient, serviceClient } from '$lib/api';
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
		TriangleAlert,
		Waves
	} from '@lucide/svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { type IconComponent } from '$lib/types';

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
					truncate: true,
					maxLength: 20
				});
			}
		},
		{
			header: 'Type',
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
				const protocol = row.getValue('type') as RouterType;
				const label = getProtocolLabel(protocol);
				const iconMap: Record<RouterType, IconComponent> = {
					[RouterType.HTTP]: Globe,
					[RouterType.TCP]: Network,
					[RouterType.UDP]: Waves,
					[RouterType.UNSPECIFIED]: TriangleAlert
				};
				return renderComponent(ColumnBadge<Router>, {
					label,
					icon: iconMap[protocol],
					variant: 'outline',
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
				let entrypoints = row.original.config?.entryPoints as string[];
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
				let middlewares = row.original.config?.middlewares as string[];
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
			id: 'rules',
			enableSorting: true,
			cell: ({ row }) => {
				return renderComponent(ColumnRule, {
					rule: (row.original.config?.rule as string) ?? '',
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
				const tls = row.original.config?.tls as RouterTLSConfig | RouterTCPTLSConfig;
				return renderComponent(ColumnTls<Router>, { tls, column });
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
							label: row.original.enabled ? 'Disable' : 'Enable',
							icon: row.original.enabled ? Power : PowerOff,
							iconProps: {
								class: row.original.enabled ? 'text-green-500' : 'text-red-500'
							},
							onClick: () => {
								row.original.enabled = !row.original.enabled;
								updateItem(row.original);
							}
						},
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
							type: 'popover',
							label: 'Delete Router',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original),
							popover: {
								title: 'Delete Router?',
								description: 'This router and its configuration will be permanently deleted.',
								confirmLabel: 'Delete',
								cancelLabel: 'Cancel'
							}
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

	const updateItem = async (item: Router) => {
		try {
			await routerClient.updateRouter({
				id: item.id,
				name: item.name,
				type: item.type,
				config: item.config,
				enabled: item.enabled,
				dnsProviders: item.dnsProviders
			});
			const service = await serviceClient.getServiceByRouter({
				name: item.name,
				type: item.type
			});
			if (service.service) {
				service.service.enabled = item.enabled;
				await serviceClient.updateService({
					id: service.service.id,
					name: service.service.name,
					config: service.service.config,
					type: service.service.type,
					enabled: service.service.enabled
				});
			}
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`Router ${item.name} updated`);
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

	onMount(() => {
		refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);

		// Set up resize handler
		const handleResize = () => {
			const isMobile = window.matchMedia('(max-width: 768px)').matches;
			if (isMobile) viewMode = 'grid';
		};

		handleResize(); // Check initial state
		window.addEventListener('resize', handleResize);
		return () => window.removeEventListener('resize', handleResize);
	});
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<div class="flex flex-col gap-4 sm:gap-6">
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div class="space-y-2">
			<h1 class="flex items-center gap-3 text-2xl font-bold tracking-tight sm:text-3xl">
				<div class="bg-primary/10 flex h-10 w-10 items-center justify-center rounded-lg">
					<Route class="text-primary h-5 w-5 sm:h-6 sm:w-6" />
				</div>
				<span class="truncate">Routers</span>
			</h1>
			<p class="text-muted-foreground text-sm sm:text-base">Manage your routers and services</p>
		</div>

		<!-- View Toggle (Don't show on mobile) -->
		<div class="hidden items-center gap-2 self-start sm:self-auto md:flex">
			<Button
				variant={viewMode === 'table' ? 'default' : 'outline'}
				size="sm"
				onclick={() => (viewMode = 'table')}
				class="flex-1 sm:flex-none"
			>
				<Table class="h-4 w-4 sm:mr-2" />
				<span class="hidden sm:block">Table</span>
			</Button>
			<Button
				variant={viewMode === 'grid' ? 'default' : 'outline'}
				size="sm"
				onclick={() => (viewMode = 'grid')}
				class="flex-1 sm:flex-none"
			>
				<LayoutGrid class="h-4 w-4 sm:mr-2" />
				<span class="hidden sm:block">Grid</span>
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
