<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import {
		profile,
		deleteRouter,
		entrypoints,
		middlewares,
		routers,
		toggleEntrypoint,
		toggleMiddleware,
		provider,
		getService
	} from '$lib/api';
	import CreateRouter from '$lib/components/modals/createRouter.svelte';
	import UpdateRouter from '$lib/components/modals/updateRouter.svelte';
	import Pagination from '$lib/components/tables/pagination.svelte';
	import type { Router } from '$lib/types/config';
	import type { Selected } from 'bits-ui';
	import Input from '$lib/components/ui/input/input.svelte';
	import ShowRouter from '$lib/components/modals/showRouter.svelte';
	import { onMount } from 'svelte';

	let search = '';
	let count = 0;
	let currentPage = 1;
	let fRouters: Router[] = [];
	let perPage: Selected<number> | undefined = JSON.parse(
		localStorage.getItem('limit') ?? '{"value": 10, "label": "10"}'
	);
	$: search, $routers, currentPage, perPage, searchRouter();

	// Reset the page to 1 when the search input changes
	$: {
		if (search) {
			currentPage = 1;
		}
	}

	const searchRouter = () => {
		let items = $routers.filter((router) => {
			if (localProvider && router.provider !== 'http') return false;
			const searchParts = search.toLowerCase().split(' ');
			return searchParts.every((part) =>
				part.startsWith('@provider:')
					? router.provider?.toLowerCase() === part.split(':')[1]
					: part.startsWith('@type:')
						? router.routerType.toLowerCase() === part.split(':')[1]
						: part.startsWith('@dns:')
							? router.dnsProvider?.toLowerCase() === part.split(':')[1]
							: router.name.toLowerCase().includes(part)
			);
		});

		count = items.length || 1;
		fRouters = paginate(items);
	};

	const paginate = (routers: Router[]) => {
		const itemsPerPage = perPage?.value ?? 10;
		return routers.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);
	};

	let columns: Selected<string>[] | undefined = [
		{ value: 'name', label: 'Name' },
		{ value: 'provider', label: 'Provider' },
		{ value: 'dns', label: 'DNS' },
		{ value: 'type', label: 'Type' },
		{ value: 'rule', label: 'Rule' },
		{ value: 'entrypoints', label: 'Entrypoints' },
		{ value: 'middlewares', label: 'Middlewares' },
		{ value: 'serviceStatus', label: 'Service Status' }
	];
	let selectedColumns: string[] = JSON.parse(
		localStorage.getItem('router-columns') ?? JSON.stringify(columns.map((c) => c.value))
	);
	$: showColumn = (column: string): boolean => {
		return selectedColumns.includes(column);
	};
	const changeColumns = (columns: Selected<string>[] | undefined) => {
		if (columns === undefined) return;
		selectedColumns = columns.map((c) => c.value);
		localStorage.setItem('router-columns', JSON.stringify(selectedColumns));
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

	// Show how many services are up and total services
	const getServiceStatus = (router: Router) => {
		let service = getService(router);
		let total = service?.loadBalancer?.servers?.length || 0;
		let up = Object.values(service?.serverStatus || {}).filter((status) => status === 'UP').length;
		return { status: `${up}/${total}`, ok: up === total };
	};

	// Only show local routers not external ones
	let localProvider = localStorage.getItem('local-provider') === 'true';
	const toggleProvider = () => {
		localProvider = !localProvider;
		search = localProvider ? '@provider:http' : '';
		localStorage.setItem('local-provider', localProvider.toString());
	};

	// Add reactive variables for bulk actions
	let allChecked = false;
	let selectedRouters: Router[] = [];
	let bulkEntrypoints: Selected<unknown>[] | undefined = [];
	let bulkMiddlewares: Selected<unknown>[] | undefined = [];
	let bulkDnsProvider: Selected<string> | undefined = undefined;
	let lastSelectedIndex: number | null = null;
	let shiftKeyPressed = false;

	const toggleRouterSelection = (router: Router) => {
		const currentIndex = fRouters.findIndex((r) => r.name === router.name);

		// Check if shift key is held
		if (shiftKeyPressed && lastSelectedIndex !== null) {
			const start = Math.min(lastSelectedIndex, currentIndex);
			const end = Math.max(lastSelectedIndex, currentIndex);

			// Select all routers between the last selected and the current one
			const rangeToSelect = fRouters.slice(start, end + 1);
			const allSelected = rangeToSelect.every((r) => selectedRouters.includes(r));

			// If all routers in range are selected, deselect them, otherwise select them
			rangeToSelect.forEach((r) => {
				if (allSelected) {
					selectedRouters = selectedRouters.filter((sr) => sr !== r); // Deselect if already selected
				} else {
					if (!selectedRouters.includes(r)) {
						selectedRouters = [...selectedRouters, r]; // Select if not selected
					}
				}
			});
		} else {
			if (selectedRouters.includes(router)) {
				selectedRouters = selectedRouters.filter((r) => r !== router);
			} else {
				selectedRouters = [...selectedRouters, router];
			}
		}
		lastSelectedIndex = currentIndex; // Update the last selected index
	};

	const applyBulkChanges = () => {
		selectedRouters.forEach((router) => {
			if (bulkEntrypoints && router.provider === 'http') {
				if (bulkEntrypoints?.length > 0) {
					toggleEntrypoint(router, bulkEntrypoints);
				}
			}
			if (bulkMiddlewares && router.provider === 'http') {
				if (bulkMiddlewares?.length > 0) {
					toggleMiddleware(router, bulkMiddlewares);
				}
			}
			if (bulkDnsProvider) {
				router.dnsProvider = bulkDnsProvider.value;
			}
		});
		// Reset after applying changes
		selectedRouters = [];
		bulkEntrypoints = [];
		bulkMiddlewares = [];
		bulkDnsProvider = undefined;
		allChecked = false;
	};

	onMount(() => {
		search = localProvider ? '@provider:http' : '';
		searchRouter();
	});
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<div class="mt-4 flex flex-col gap-4 p-4">
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
				class:bg-primary={localProvider}
				class:text-primary-foreground={localProvider}
				on:click={toggleProvider}
			>
				Local Only
			</button>
		</div>
		<Select.Root
			multiple
			selected={selectedColumns.map((c) => ({ value: c, label: c }))}
			onSelectedChange={changeColumns}
		>
			<Select.Trigger class="w-[180px]">
				<Select.Value placeholder="Columns" />
			</Select.Trigger>
			<Select.Content>
				{#each columns as column}
					<Select.Item value={column.value} label={column.label}>{column.label}</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
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
						<Table.Head>
							<Checkbox
								id="routers"
								checked={allChecked}
								onCheckedChange={() => {
									allChecked = !allChecked;
									selectedRouters = allChecked ? [...fRouters] : [];
									lastSelectedIndex = null;
								}}
							/>
						</Table.Head>
						{#if showColumn('name')}
							<Table.Head>Name</Table.Head>
						{/if}
						{#if showColumn('provider')}
							<Table.Head>Provider</Table.Head>
						{/if}
						{#if showColumn('dns')}
							<Table.Head>DNS</Table.Head>
						{/if}
						{#if showColumn('type')}
							<Table.Head>Type</Table.Head>
						{/if}
						{#if showColumn('rule')}
							<Table.Head class="hidden lg:table-cell">Rule</Table.Head>
						{/if}
						{#if showColumn('entrypoints')}
							<Table.Head class="hidden lg:table-cell">Entrypoints</Table.Head>
						{/if}
						{#if showColumn('middlewares')}
							<Table.Head class="hidden lg:table-cell">Middlewares</Table.Head>
						{/if}
						{#if showColumn('serviceStatus')}
							<Table.Head>Service Status</Table.Head>
						{/if}
						<Table.Head>
							<span class="sr-only">Edit</span>
						</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each fRouters as router}
						<Table.Row>
							<Table.Cell class="min-w-[2rem]">
								<div
									on:keydown={(e) => (e.key === 'Shift' ? (shiftKeyPressed = true) : null)}
									on:keyup={(e) => (e.key === 'Shift' ? (shiftKeyPressed = false) : null)}
									aria-hidden
								>
									<Checkbox
										id={router.name}
										checked={selectedRouters.includes(router)}
										onCheckedChange={() => toggleRouterSelection(router)}
									/>
								</div>
							</Table.Cell>
							<Table.Cell class={showColumn('name') ? 'font-medium' : 'hidden'}>
								{router.name.split('@')[0]}
							</Table.Cell>
							<Table.Cell class={showColumn('provider') ? 'font-medium' : 'hidden'}>
								<span
									class="inline-flex cursor-pointer select-none items-center rounded-full bg-slate-300 px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
									on:click={() => (search = `@provider:${router.provider}`)}
									aria-hidden
								>
									{router.provider}
								</span>
							</Table.Cell>
							<Table.Cell class={showColumn('dns') ? 'font-medium' : 'hidden'}>
								<span
									class="inline-flex cursor-pointer select-none items-center rounded-full px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
									class:bg-green-300={router.dnsProvider}
									class:bg-blue-300={!router.dnsProvider}
									on:click={() => (search = `@dns:${router.dnsProvider}`)}
									aria-hidden
								>
									{router.dnsProvider ? router.dnsProvider : 'None'}
								</span>
							</Table.Cell>
							<Table.Cell class={showColumn('type') ? 'font-medium' : 'hidden'}>
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
								class={showColumn('rule')
									? 'hidden max-w-[180px] overflow-hidden text-ellipsis whitespace-nowrap lg:table-cell'
									: 'hidden'}
							>
								{#if 'rule' in router}
									{#if router?.rule !== ''}
										{router.rule}
									{:else}
										<span class="text-muted-foreground">N/A</span>
									{/if}
								{/if}
							</Table.Cell>
							<Table.Cell class={showColumn('entrypoints') ? 'hidden lg:table-cell' : 'hidden'}>
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
							<Table.Cell class={showColumn('middlewares') ? 'hidden lg:table-cell' : 'hidden'}>
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
							<Table.Cell class={showColumn('serviceStatus') ? 'font-medium' : 'hidden'}>
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
										on:click={() => deleteRouter(router.name)}
									>
										<iconify-icon icon="fa6-solid:xmark" />
									</Button>
								{:else}
									<ShowRouter {router} />
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

	<Pagination {count} bind:perPage bind:currentPage />

	<!-- Bulk Edit Footer -->
	{#if selectedRouters.length > 0}
		<div class="fixed bottom-4 flex flex-row items-center justify-between">
			<Card.Root>
				<Card.Content class="flex flex-row items-center justify-between gap-4 p-4 shadow-md">
					<div class="flex flex-col items-center justify-start gap-4 md:flex-row">
						<!-- Bulk update entrypoints -->
						<Select.Root
							multiple={true}
							selected={bulkEntrypoints}
							onSelectedChange={(value) => (bulkEntrypoints = value)}
						>
							<Select.Trigger class="w-[200px]">
								<Select.Value placeholder="EntryPoints" />
							</Select.Trigger>
							<Select.Content>
								{#each $entrypoints as entrypoint}
									<Select.Item value={entrypoint.name}>{entrypoint.name}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>

						<!-- Bulk update middlewares -->
						<Select.Root
							multiple={true}
							selected={bulkMiddlewares}
							onSelectedChange={(value) => (bulkMiddlewares = value)}
						>
							<Select.Trigger class="w-[200px]">
								<Select.Value placeholder="Middlewares" />
							</Select.Trigger>
							<Select.Content>
								{#each $middlewares as middleware}
									<Select.Item value={middleware.name}>{middleware.name}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>

						<!-- Bulk update DNS Provider -->
						<Select.Root
							selected={bulkDnsProvider}
							onSelectedChange={(value) => (bulkDnsProvider = value)}
						>
							<Select.Trigger class="w-[200px]">
								<Select.Value placeholder="DNS Provider" />
							</Select.Trigger>
							<Select.Content>
								{#each $provider as p}
									<Select.Item value={p.name}>{p.name}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>

					<div class="flex flex-row items-center gap-2">
						<Button variant="secondary" on:click={() => (selectedRouters = [])}>Clear</Button>
						<Button on:click={applyBulkChanges}>Apply Changes</Button>
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	{/if}
</div>
