<script lang="ts">
	import { api, profile, routers, services } from '$lib/api';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import RouterModal from '$lib/components/modals/router.svelte';
	import { renderComponent } from '$lib/components/ui/data-table';
	import type { Router, Service, TLS } from '$lib/types/router';
	import type { ColumnDef } from '@tanstack/table-core';
	import { Edit, Trash } from 'lucide-svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';

	export const routerModalState = writable({
		isOpen: false,
		router: {} as Router | undefined,
		service: {} as Service | undefined
	});

	const deleteRouter = async (name: string | undefined, type: string) => {
		if (!name) return;
		try {
			let routerProvider = name.split('@')[1];
			if (routerProvider !== 'http') {
				toast.error('Router not managed by Mantrae!');
				return;
			}

			await api.deleteRouter($profile.id, name, type);
			await api.getTraefikConfig($profile.id);
			toast.success('Router deleted');
		} catch (err) {
			toast.error(err.message);
		}
	};

	type RouterWithService = { router: Router; service: Service };
	let mergedData: RouterWithService[] = $state([]);

	const columns: ColumnDef<RouterWithService>[] = [
		{
			header: 'Name',
			accessorFn: (row) => row.router.name,
			id: 'name',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('name') as string;
				return name.split('@')[0];
			}
		},
		{
			header: 'Type',
			accessorFn: (row) => row.router.type,
			id: 'type',
			enableSorting: true,
			cell: ({ row }) => {
				const type = row.getValue('type') as string;
				return renderComponent(ColumnBadge, { label: type });
			}
		},
		{
			header: 'Provider',
			accessorFn: (row) => row.router.name,
			id: 'provider',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('provider') as string;
				return renderComponent(ColumnBadge, {
					label: name.split('@')[1].toLowerCase(),
					variant: 'secondary'
				});
			}
		},
		{
			header: 'Entrypoints',
			accessorFn: (row) => row.router.entryPoints,
			id: 'entryPoints',
			enableSorting: true,
			cell: ({ row }) => {
				const entryPoints = row.getValue('entryPoints') as string[];
				return renderComponent(ColumnBadge, {
					label: entryPoints.join(', '),
					variant: 'secondary'
				});
			}
		},
		{
			header: 'Cert Resolver',
			accessorFn: (row) => row.router.tls,
			id: 'tls',
			enableSorting: true,
			cell: ({ row }) => {
				const tls = row.getValue('tls') as TLS;
				if (!tls) {
					return renderComponent(ColumnBadge, { label: 'None', variant: 'secondary' });
				}
				return renderComponent(ColumnBadge, {
					label: tls.certResolver as string,
					variant: 'secondary',
					class: 'bg-slate-300 dark:bg-slate-700'
				});
			}
		},
		{
			id: 'actions',
			cell: ({ row }) => {
				return renderComponent(TableActions, {
					actions: [
						{
							label: 'Edit Router',
							icon: Edit,
							onClick: () => {
								routerModalState.set({
									isOpen: true,
									router: row.original.router,
									service: row.original.service
								});
							}
						},
						{
							label: 'Delete Router',
							icon: Trash,
							variant: 'destructive',
							onClick: () => {
								deleteRouter(row.original.router.name, row.original.router.type);
							}
						}
					]
				});
			}
		}
	];

	profile.subscribe(async (value) => {
		if (!value?.id) return;
		await api.getTraefikConfig($profile.id);
		mergedData = $routers.map((router) => {
			const service = $services.find((service) => service.name === router.name);
			if (!service) {
				return { router, service: {} as Service };
			}
			return { router, service };
		});
	});
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<div class="container flex flex-col gap-4">
	<div class="flex flex-col justify-start">
		<h1 class="text-2xl font-bold">Router Management</h1>
		<span class="text-sm text-muted-foreground">Total Routers: {$routers.length}</span>
	</div>
	<DataTable
		{columns}
		data={mergedData || []}
		createButton={{
			label: 'Add Router',
			modal: RouterModal
		}}
	/>
</div>

<RouterModal
	bind:router={$routerModalState.router}
	bind:service={$routerModalState.service}
	bind:open={$routerModalState.isOpen}
/>
