<!-- @migration-task Error while migrating Svelte code: `$routers` is an illegal variable name. To reference a global variable called `$routers`, use `globalThis.$routers`
https://svelte.dev/e/global_reference_invalid -->
<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Select from '$lib/components/ui/select';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { api, profile } from '$lib/api';
	import Pagination from '$lib/components/tables/pagination.svelte';
	import { newRouter, newService, type Router, type Service } from '$lib/types/router';
	import type { Selected } from 'bits-ui';
	import RouterModal from '$lib/components/modals/router.svelte';
	import Search from '$lib/components/tables/search.svelte';
	import { page } from '$app/stores';
	import {
		Lock,
		Eye,
		Pencil,
		Bot,
		X,
		SquareArrowOutUpRight,
		ShieldAlert,
		TriangleAlert,
		BotOff,
		Plus
	} from 'lucide-svelte';
	import { LIMIT_SK, ROUTER_COLUMN_SK } from '$lib/store';

	let search = '';
	let count = 0;
	let currentPage = 1;
	let fRouters: Router[] = [];
	let perPage: Selected<number> | undefined = JSON.parse(
		localStorage.getItem(LIMIT_SK) ?? '{"value": 10, "label": "10"}'
	);
	$: search, $routers, currentPage, perPage, searchRouter();

	page.subscribe((p) => {
		if (p.url.searchParams.get('search')) {
			search = p.url.searchParams.get('search') ?? '';
		}
	});

	// Reset the page to 1 when the search input changes
	$: {
		if (search) {
			currentPage = 1;
		}
	}

	const searchRouter = () => {
		if ($routers === undefined) return;
		let items = $routers?.filter((router) => {
			const searchParts = search.toLowerCase().split(' ');
			return searchParts.every((part) =>
				part.startsWith('@provider:')
					? router.provider?.toLowerCase() === part.split(':')[1]
					: part.startsWith('@protocol:')
						? router.protocol.toLowerCase() === part.split(':')[1]
						: part.startsWith('@dns:')
							? getDNSProviderName(router)?.toLowerCase() === part.split(':')[1]
							: router.name.toLowerCase().includes(part)
			);
		});

		count = items?.length || 1;
		fRouters = paginate(items);
	};

	const paginate = (routers: Router[]) => {
		const itemsPerPage = perPage?.value ?? 10;
		return routers?.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage) ?? [];
	};

	// Open the modal for the selected router
	let router: Router;
	let service: Service;
	let openModal = false;
	let disabled = false;
	const createModal = () => {
		if (!$profile?.id) return;
		router = newRouter();
		service = newService();
		router.profileId = $profile.id;
		service.profileId = $profile.id;
		disabled = false;
		openModal = true;
	};
	const updateModal = (r: Router) => {
		if (!r) return;
		if (r.provider === 'http') {
			disabled = false;
		} else {
			disabled = true;
		}
		router = r;
		service = $services?.find((s: Service) => s.name === r.name) ?? newService();
		openModal = true;
	};

	let columns: Selected<string>[] | undefined = [
		{ value: 'name', label: 'Name' },
		{ value: 'provider', label: 'Provider' },
		{ value: 'dns', label: 'DNS' },
		{ value: 'protocol', label: 'Protocol' },
		{ value: 'rule', label: 'Rule' },
		{ value: 'entrypoints', label: 'Entrypoints' },
		{ value: 'middlewares', label: 'Middlewares' },
		{ value: 'serviceStatus', label: 'Service Status' }
	];
	let fColumns: string[] = JSON.parse(
		localStorage.getItem(ROUTER_COLUMN_SK) ?? JSON.stringify(columns.map((c) => c.value))
	);

	const getSelectedEntrypoints = (router: Router): Selected<unknown>[] => {
		let list = router?.entryPoints?.map((ep) => {
			return { value: ep, label: ep };
		});
		return list ?? [];
	};
	const getSelectedMiddlewares = (router: Router): Selected<unknown>[] => {
		if (router?.middlewares === undefined) return [];
		return router.middlewares?.map((middleware) => {
			return { value: middleware, label: middleware };
		});
	};

	// Show how many services are up and total services
	const getServiceStatus = (router: Router) => {
		if (router === undefined) return { status: '0/0', ok: false };
		let service = $services?.find((s: Service) => s.name === router.service);
		let total = service?.loadBalancer?.servers?.length || 0;
		let up = Object.values(service?.serverStatus || {}).filter((status) => status === 'UP').length;
		return { status: `${up}/${total}`, ok: up === total };
	};

	const hasTLS = (router: Router): boolean => {
		if (router.entryPoints === undefined) return false;

		return router?.entryPoints?.some((e) => {
			let entrypoint = $entrypoints?.find((ep) => ep.name === e);
			return entrypoint?.http?.tls !== undefined;
		});
	};

	function getHost(router: Router): string {
		const hostRegex = /Host\(`([^`]+)`\)/;
		const pathPrefixRegex = /PathPrefix\(`([^`]+)`\)/;
		const schema = router.tls?.certResolver ? 'https' : 'http';

		const match = router.rule?.match(hostRegex);
		const matchPath = router.rule?.match(pathPrefixRegex);
		let link = '';
		if (match && match[1]) {
			link = `${schema}://${match[1]}`;
			if (matchPath && matchPath[1]) {
				link = `${link}${matchPath[1]}`;
			}
		}
		return link;
	}

	// Add reactive variables for bulk actions
	let allChecked = false;
	let selectedRouters: Router[] = [];
	let bulkEntrypoints: Selected<string>[] | undefined = [];
	let bulkMiddlewares: Selected<string>[] | undefined = [];
	let bulkDnsProvider: Selected<number> | undefined = undefined;
	let lastSelectedIndex: number | null = null;
	let shiftKeyPressed = false;

	const toggleRouterSelection = (router: Router) => {
		const currentIndex = fRouters.findIndex((r) => r.id === router.id);

		// Check if shift key is held
		if (shiftKeyPressed && lastSelectedIndex !== null) {
			const start = Math.min(lastSelectedIndex, currentIndex);
			const end = Math.max(lastSelectedIndex, currentIndex);

			// Select all routers between the last selected and the current one
			const rangeToSelect = fRouters.slice(start, end + 1);
			const allSelected = rangeToSelect.every((r) =>
				selectedRouters.some((selected) => selected.id === r.id)
			);

			// If all routers in range are selected, deselect them, otherwise select them
			rangeToSelect.forEach((r) => {
				if (allSelected) {
					// Deselect if already selected
					selectedRouters = selectedRouters.filter((sr) => sr.id !== r.id);
				} else {
					// Select if not selected
					if (!selectedRouters.some((selected) => selected.id === r.id)) {
						selectedRouters = [...selectedRouters, r];
					}
				}
			});
		} else {
			// Toggle individual selection
			if (selectedRouters.some((selected) => selected.id === router.id)) {
				// Deselect if already selected
				selectedRouters = selectedRouters.filter((r) => r.id !== router.id);
			} else {
				// Select if not selected
				selectedRouters = [...selectedRouters, router];
			}
		}

		// Update the last selected index
		lastSelectedIndex = currentIndex;
	};

	const getDNSProviderName = (r: Router) => {
		return $provider?.find((p) => p.id === r.dnsProvider)?.name || undefined;
	};

	const applyBulkChanges = () => {
		selectedRouters.forEach((router) => {
			if (bulkEntrypoints && router.provider === 'http') {
				if (bulkEntrypoints?.length > 0) {
					toggleEntrypoint(router, bulkEntrypoints, true);
				}
			}
			if (bulkMiddlewares && router.provider === 'http') {
				if (bulkMiddlewares?.length > 0) {
					toggleMiddleware(router, bulkMiddlewares, true);
				}
			}
			if (bulkDnsProvider) {
				router.dnsProvider = bulkDnsProvider.value;
				upsertRouter(router);
			}
		});
		// Reset after applying changes
		selectedRouters = [];
		bulkEntrypoints = [];
		bulkMiddlewares = [];
		bulkDnsProvider = undefined;
		allChecked = false;
	};

	// Get routers when the profile changes
	profile.subscribe((value) => {
		if (!value?.id) return;
		getServices();
		getRouters();
	});
</script>

<svelte:head>
	<title>Routers</title>
</svelte:head>

<RouterModal bind:router bind:service bind:open={openModal} bind:disabled />

<div class="mt-4 flex flex-col gap-4 p-4">
	<Search bind:search {columns} columnName="router-columns" bind:fColumns />

	<Card.Root>
		<Card.Header class="grid grid-cols-2 items-center justify-between">
			<div>
				<Card.Title>Routers</Card.Title>
				<Card.Description>
					Total routers managed by Mantrae {$routers?.filter((r) => r.provider === 'http').length}
				</Card.Description>
			</div>
			<div class="justify-self-end">
				<Button
					variant="secondary"
					class="flex items-center gap-2 bg-red-400 text-black"
					on:click={createModal}
				>
					<span>Create Router</span>
					<Plus size="1rem" />
				</Button>
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
						{#if fColumns.includes('name')}
							<Table.Head>Name</Table.Head>
						{/if}
						{#if fColumns.includes('provider')}
							<Table.Head>Provider</Table.Head>
						{/if}
						{#if fColumns.includes('dns')}
							<Table.Head>DNS</Table.Head>
						{/if}
						{#if fColumns.includes('protocol')}
							<Table.Head>Protocol</Table.Head>
						{/if}
						{#if fColumns.includes('rule')}
							<Table.Head class="hidden lg:table-cell">Rule</Table.Head>
						{/if}
						{#if fColumns.includes('entrypoints')}
							<Table.Head class="hidden lg:table-cell">Entrypoints</Table.Head>
						{/if}
						{#if fColumns.includes('middlewares')}
							<Table.Head class="hidden lg:table-cell">Middlewares</Table.Head>
						{/if}
						{#if fColumns.includes('serviceStatus')}
							<Table.Head>Service Status</Table.Head>
						{/if}
						<Table.Head>
							<span class="sr-only">Edit</span>
						</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each fRouters as router}
						<Table.Row class={hasTLS(router) ? 'bg-green-100/40 dark:bg-green-800/40' : ''}>
							<Table.Cell class="min-w-[2rem]">
								<div
									on:keydown={(e) => (e.key === 'Shift' ? (shiftKeyPressed = true) : null)}
									on:keyup={(e) => (e.key === 'Shift' ? (shiftKeyPressed = false) : null)}
									aria-hidden
								>
									<Checkbox
										id={router.name}
										checked={selectedRouters.some((sr) => sr.id === router.id)}
										onCheckedChange={() => toggleRouterSelection(router)}
									/>
								</div>
							</Table.Cell>
							<Table.Cell class={fColumns.includes('name') ? 'font-medium' : 'hidden'}>
								<div class="flex flex-row items-center gap-1">
									{#if getHost(router)}
										<a
											href={getHost(router)}
											target="_blank"
											rel="noreferrer"
											class="flex flex-row items-center gap-1 text-blue-600"
										>
											{router.name.split('@')[0]}
											<SquareArrowOutUpRight size="1rem" />
										</a>
									{:else}
										{router.name.split('@')[0]}
									{/if}

									{#if router.agentId}
										{#if router.errors && router.errors.agent}
											<HoverCard.Root openDelay={400}>
												<HoverCard.Trigger>
													<BotOff size="1.25rem" class="ml-1 text-red-600" />
												</HoverCard.Trigger>
												<HoverCard.Content>
													<div class="text-sm">Agent Error: {router.errors.agent}</div>
												</HoverCard.Content>
											</HoverCard.Root>
										{:else}
											<Bot size="1.25rem" class="ml-1 text-green-600" />
										{/if}
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell class={fColumns.includes('provider') ? 'font-medium' : 'hidden'}>
								<span
									class="inline-flex cursor-pointer select-none items-center rounded-full bg-slate-300 px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
									on:click={() => (search = `@provider:${router.provider}`)}
									aria-hidden
								>
									{router.provider}
								</span>
							</Table.Cell>
							<Table.Cell class={fColumns.includes('dns') ? 'font-medium' : 'hidden'}>
								<div class="flex flex-row items-center gap-1">
									<span
										class="inline-flex cursor-pointer select-none items-center rounded-full px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
										class:bg-green-300={getDNSProviderName(router)}
										class:bg-blue-300={!getDNSProviderName(router)}
										on:click={() => (search = `@dns:${getDNSProviderName(router)}`)}
										aria-hidden
									>
										{getDNSProviderName(router) ? getDNSProviderName(router) : 'None'}
									</span>
									{#if router.errors && router.errors.dns}
										<HoverCard.Root openDelay={400}>
											<HoverCard.Trigger>
												<Badge variant="secondary" class="bg-orange-300">
													<TriangleAlert size="1rem" />
												</Badge>
											</HoverCard.Trigger>
											<HoverCard.Content class="text-sm text-slate-800">
												DNS Error: {router.errors.dns}
											</HoverCard.Content>
										</HoverCard.Root>
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell class={fColumns.includes('protocol') ? 'font-medium' : 'hidden'}>
								<span
									class="inline-flex cursor-pointer select-none items-center rounded-full px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
									class:bg-green-300={router.protocol === 'http'}
									class:bg-blue-300={router.protocol === 'tcp'}
									class:bg-purple-300={router.protocol === 'udp'}
									on:click={() => (search = `@protocol:${router.protocol}`)}
									aria-hidden
								>
									{router.protocol}
								</span>
							</Table.Cell>
							<Table.Cell
								class={fColumns.includes('rule')
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
							<Table.Cell
								class={fColumns.includes('entrypoints') ? 'hidden lg:table-cell' : 'hidden'}
							>
								<Select.Root
									multiple={true}
									selected={getSelectedEntrypoints(router)}
									onSelectedChange={(value) => toggleEntrypoint(router, value, true)}
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
															<Lock size="1rem" class=" text-green-400" />
														{/if}
													{/if}
												</div>
											</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>
							</Table.Cell>
							<Table.Cell
								class={fColumns.includes('middlewares') ? 'hidden lg:table-cell' : 'hidden'}
							>
								<div class:hidden={router.protocol === 'udp'}>
									<Select.Root
										multiple={true}
										selected={getSelectedMiddlewares(router)}
										onSelectedChange={(value) => toggleMiddleware(router, value, true)}
										disabled={router.provider !== 'http'}
									>
										<Select.Trigger class="w-[180px]">
											<Select.Value placeholder="Select a middleware" />
										</Select.Trigger>
										<Select.Content class="text-sm">
											{#each $middlewares as middleware}
												{#if router.protocol === middleware.protocol}
													<Select.Item value={middleware.name}>
														{middleware.name}
													</Select.Item>
												{/if}
											{/each}
										</Select.Content>
									</Select.Root>
								</div>
							</Table.Cell>
							<Table.Cell class={fColumns.includes('serviceStatus') ? 'font-medium' : 'hidden'}>
								<div class="flex flex-row items-center gap-2">
									<span
										class="rounded-full px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-opacity-80 focus:outline-none"
										class:bg-green-300={getServiceStatus(router).ok}
										class:bg-red-300={!getServiceStatus(router).ok}
									>
										{getServiceStatus(router).status}
									</span>
									{#if router.errors && router.errors.ssl}
										<HoverCard.Root openDelay={400}>
											<HoverCard.Trigger>
												<Badge variant="secondary" class="bg-red-300">
													<ShieldAlert size="1rem" />
												</Badge>
											</HoverCard.Trigger>
											<HoverCard.Content>
												<div class="text-sm">SSL Error: {router.errors.ssl}</div>
											</HoverCard.Content>
										</HoverCard.Root>
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell class="min-w-[100px]">
								{#if router.provider === 'http'}
									<Button
										variant="ghost"
										class="h-8 w-8 rounded-full bg-orange-400"
										size="icon"
										on:click={() => updateModal(router)}
									>
										<Pencil size="1rem" />
									</Button>
									{#if !router.agentId}
										<Button
											variant="ghost"
											class="h-8 w-8 rounded-full bg-red-400"
											size="icon"
											on:click={() => deleteRouter(router)}
										>
											<X size="1rem" />
										</Button>
									{/if}
								{:else}
									<Button
										variant="ghost"
										class="h-8 w-8 rounded-full bg-green-400"
										size="icon"
										on:click={() => updateModal(router)}
									>
										<Eye size="1rem" />
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
				Showing <strong>{fRouters.length > 0 ? 1 : 0}-{fRouters?.length ?? 0}</strong> of
				<strong>{$routers?.length ?? 0}</strong> routers
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
						<!--Bulk update entrypoints-->
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
						{#if $provider}
							<Select.Root
								selected={bulkDnsProvider}
								onSelectedChange={(value) => (bulkDnsProvider = value)}
							>
								<Select.Trigger class="w-[200px]">
									<Select.Value placeholder="DNS Provider" />
								</Select.Trigger>
								<Select.Content>
									{#each $provider as p}
										<Select.Item value={p.id}>{p.name}</Select.Item>
									{/each}
								</Select.Content>
							</Select.Root>
						{/if}
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
