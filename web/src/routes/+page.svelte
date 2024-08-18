<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Select from '$lib/components/ui/select';
	import {
		activeProfile,
		deleteRouter,
		entrypoints,
		middlewares,
		routers,
		services,
		updateProfile
	} from '$lib/api';
	import CreateRouter from '$lib/components/modals/createRouter.svelte';
	import UpdateRouter from '$lib/components/modals/updateRouter.svelte';
	import Pagination from '$lib/components/tables/pagination.svelte';
	import type { Router } from '$lib/types/config';
	import type { Selected } from 'bits-ui';
	import Input from '$lib/components/ui/input/input.svelte';
	import { onMount } from 'svelte';

	let search = '';
	let count = 0;
	let currentPage = 1;
	let fRouters: Router[] = [];
	let perPage: Selected<number> | undefined = { value: 10, label: '10' }; // Items per page
	$: search, $routers, currentPage, searchRouter();

	// Reset the page to 1 when the search input changes
	$: {
		if (search) {
			currentPage = 1;
		}
	}

	const searchRouter = () => {
		let items: Router[] = [...$routers];

		if (search) {
			const searchParts = search.split(' ').map((part) => part.toLowerCase());
			let providerSearch = '';
			let typeSearch = '';
			let generalSearch = [];

			for (const part of searchParts) {
				if (part.startsWith('@provider:')) {
					providerSearch = part.split(':')[1];
				} else if (part.startsWith('@type:')) {
					typeSearch = part.split(':')[1];
				} else {
					generalSearch.push(part);
				}
			}

			items = items.filter((router) => {
				const matchesProvider = providerSearch
					? router.provider?.toLowerCase() === providerSearch
					: true;
				const matchesType = typeSearch ? router.routerType.toLowerCase() === typeSearch : true;
				const matchesGeneral = generalSearch.every((term) =>
					router.service.toLowerCase().includes(term)
				);
				return matchesProvider && matchesType && matchesGeneral;
			});
		}

		fRouters = paginate(items);
		count = items.length || 1;
	};

	const paginate = (routers: Router[]) => {
		const itemsPerPage = perPage?.value ?? 10;
		return routers.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);
	};

	// let columns: Selected<string>[] | undefined = [
	// 	{ value: 'name', label: 'Name' },
	// 	{ value: 'provider', label: 'Provider' },
	// 	{ value: 'type', label: 'Type' },
	// 	{ value: 'rule', label: 'Rule' },
	// 	{ value: 'entrypoints', label: 'Entrypoints' },
	// 	{ value: 'middlewares', label: 'Middlewares' }
	// ];
	// $: selectedColumns = [...columns];
	//
	// $: showColumn = (column: string): boolean => {
	// 	return selectedColumns.find((c) => c.value === column) !== undefined;
	// };

	const toggleEntrypoint = (router: Router, item: Selected<unknown>[] | undefined) => {
		if (item === undefined) return;
		router.entrypoints = item.map((i) => i.value) as string[];
		updateProfile($activeProfile.name, $activeProfile);
	};
	const toggleMiddleware = (router: Router, item: Selected<unknown>[] | undefined) => {
		if (item === undefined) return;
		router.middlewares = item.map((i) => i.value) as string[];
		updateProfile($activeProfile.name, $activeProfile);
	};
	const getSelectedEntrypoints = (router: Router): Selected<unknown>[] => {
		let list = router?.entrypoints?.map((entrypoint) => {
			return { value: entrypoint, label: entrypoint };
		});
		return list ?? [];
	};
	const getSelectedMiddlewares = (router: Router): Selected<unknown>[] => {
		if (router?.middlewares === undefined) return [];
		return router.middlewares.map((middleware) => {
			return { value: middleware, label: middleware };
		});
	};

	const getServiceStatus = (router: Router): Record<string, string | boolean> => {
		let service = $services.find((s) => s.name === router.service + '@' + router.provider);
		let totalServices = service?.loadBalancer?.servers?.length || 0;

		let upServices = 0;
		if (service?.serverStatus !== undefined) {
			upServices = Object.keys(service.serverStatus).filter(
				(key) => service.serverStatus !== undefined && service.serverStatus[key] === 'UP'
			).length;
		}
		let status = `${upServices}/${totalServices}`;
		let ok = upServices === totalServices ? true : false;
		return { status: status, ok: ok };
	};

	// Only show local routers not external ones
	let onlyLocal = localStorage.getItem('local-only') === 'true';
	const toggleLocalOnly = () => {
		onlyLocal = !onlyLocal;
		localStorage.setItem('localOnly', onlyLocal.toString());
		search = onlyLocal ? '@provider:http' : '';
	};

	onMount(() => {
		onlyLocal = localStorage.getItem('localOnly') === 'true';
		search = onlyLocal ? '@provider:http' : '';
	});
</script>

<svelte:head>
	<title>Routers | {$activeProfile?.name}</title>
	<meta name="description" content="Traefik Web UI" />
</svelte:head>

<div class="flex flex-row items-center justify-between">
	<div class="flex flex-row items-center gap-1">
		<Input
			type="text"
			placeholder="Search..."
			class="w-80 focus-visible:ring-0 focus-visible:ring-offset-0"
			bind:value={search}
		/>
		<Button variant="outline" on:click={() => (search = '')} aria-hidden>
			<iconify-icon icon="fa6-solid:xmark" />
		</Button>
		<button
			class={buttonVariants({ variant: 'outline' })}
			class:bg-primary={onlyLocal}
			class:text-primary-foreground={onlyLocal}
			on:click={toggleLocalOnly}
		>
			Local Only
		</button>
	</div>
	<!-- <Select.Root -->
	<!-- 	multiple={true} -->
	<!-- 	selected={selectedColumns} -->
	<!-- 	onSelectedChange={(value) => (selectedColumns = value ?? [])} -->
	<!-- > -->
	<!-- 	<Select.Trigger class="w-[180px]"> -->
	<!-- 		<Select.Value placeholder="Columns" /> -->
	<!-- 	</Select.Trigger> -->
	<!-- 	<Select.Content> -->
	<!-- 		{#each columns as column} -->
	<!-- 			<Select.Item value={column.value}>{column.label}</Select.Item> -->
	<!-- 		{/each} -->
	<!-- 	</Select.Content> -->
	<!-- </Select.Root> -->
</div>

<Card.Root>
	<Card.Header class="grid grid-cols-2 items-center justify-between">
		<div>
			<Card.Title>Routers</Card.Title>
			<Card.Description>Manage your Routers and view their status.</Card.Description>
		</div>
		<div class="justify-self-end">
			<CreateRouter />
		</div>
	</Card.Header>
	<Card.Content>
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Name</Table.Head>
					<Table.Head>Provider</Table.Head>
					<Table.Head>Type</Table.Head>
					<Table.Head class="hidden lg:table-cell">Rule</Table.Head>
					<Table.Head class="hidden lg:table-cell">Entrypoints</Table.Head>
					<Table.Head class="hidden lg:table-cell">Middlewares</Table.Head>
					<Table.Head>Service Status</Table.Head>
					<Table.Head>
						<span class="sr-only">Edit</span>
					</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each fRouters as router}
					<Table.Row>
						<Table.Cell class="font-medium">
							{router.service}
						</Table.Cell>
						<Table.Cell class="font-medium">
							<span
								class="inline-flex cursor-pointer select-none items-center rounded-full bg-slate-300 px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
								on:click={() => (search = `@provider:${router.provider}`)}
								aria-hidden
							>
								{router.provider}
							</span>
						</Table.Cell>
						<Table.Cell class="font-medium">
							<span
								class="inline-flex cursor-pointer select-none items-center rounded-full px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
								class:bg-green-300={router.routerType === 'http'}
								class:bg-blue-300={router.routerType === 'tcp'}
								class:bg-purple-300={router.routerType === 'udp'}
								on:click={() => (search = `@type:${router.routerType}`)}
								aria-hidden
							>
								{router.routerType}
							</span>
						</Table.Cell>
						<Table.Cell
							class="hidden max-w-[180px] overflow-hidden text-ellipsis whitespace-nowrap lg:table-cell"
						>
							{#if 'rule' in router}
								{#if router?.rule !== ''}
									{router.rule}
								{:else}
									<span class="text-muted-foreground">N/A</span>
								{/if}
							{/if}
						</Table.Cell>
						<Table.Cell class="hidden lg:table-cell">
							<Select.Root
								multiple={true}
								selected={getSelectedEntrypoints(router)}
								onSelectedChange={(value) => toggleEntrypoint(router, value)}
								disabled={router.provider !== 'http'}
							>
								<Select.Trigger class="w-[150px]">
									<Select.Value placeholder="Select an entrypoint" />
								</Select.Trigger>
								<Select.Content class="text-sm">
									{#each $entrypoints as entrypoint}
										<Select.Item value={entrypoint.name}>
											<div class="flex flex-row items-center gap-2">
												{entrypoint.name}
												{#if entrypoint.http}
													{#if 'tls' in entrypoint.http}
														<iconify-icon icon="fa6-solid:lock" class=" text-green-400" />
													{/if}
												{/if}
											</div>
										</Select.Item>
									{/each}
								</Select.Content>
							</Select.Root>
						</Table.Cell>
						<Table.Cell class="hidden lg:table-cell">
							<div class:hidden={router.routerType === 'udp'}>
								<Select.Root
									multiple={true}
									selected={getSelectedMiddlewares(router)}
									onSelectedChange={(value) => toggleMiddleware(router, value)}
									disabled={router.provider !== 'http'}
								>
									<Select.Trigger class="w-[180px]">
										<Select.Value placeholder="Select a middleware" />
									</Select.Trigger>
									<Select.Content class="text-sm">
										{#each $middlewares as middleware}
											{#if router.routerType === middleware.middlewareType}
												<Select.Item value={middleware.name}>
													{middleware.name}
												</Select.Item>
											{/if}
										{/each}
									</Select.Content>
								</Select.Root>
							</div>
						</Table.Cell>
						<Table.Cell>
							<span
								class="rounded-full px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-opacity-80 focus:outline-none"
								class:bg-green-300={getServiceStatus(router).ok}
								class:bg-red-300={!getServiceStatus(router).ok}
							>
								{getServiceStatus(router).status}
							</span>
						</Table.Cell>
						<Table.Cell class="min-w-[100px]">
							{#if router.provider === 'http' || router.provider === undefined}
								<UpdateRouter {router} />
								<Button
									variant="ghost"
									class="h-8 w-4 rounded-full bg-red-400"
									on:click={() => deleteRouter($activeProfile.name, router.name)}
								>
									<iconify-icon icon="fa6-solid:xmark" />
								</Button>
							{/if}
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card.Content>
	<Card.Footer>
		<div class="text-xs text-muted-foreground">
			Showing <strong>{fRouters.length > 0 ? 1 : 0}-{fRouters.length}</strong> of
			<strong>{$routers.length}</strong> routers
		</div>
	</Card.Footer>
</Card.Root>

<Pagination {count} {perPage} bind:currentPage />
