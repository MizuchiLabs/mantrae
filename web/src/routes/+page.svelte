<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Select from '$lib/components/ui/select';
	import {
		profile,
		deleteRouter,
		entrypoints,
		middlewares,
		routers,
		services,
		updateRouter
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
						: router.service.toLowerCase().includes(part)
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
		{ value: 'type', label: 'Type' },
		{ value: 'rule', label: 'Rule' },
		{ value: 'entrypoints', label: 'Entrypoints' },
		{ value: 'middlewares', label: 'Middlewares' }
	];
	let selectedColumns: string[] = JSON.parse(
		localStorage.getItem('router-columns') ??
			'["name", "provider", "type", "rule", "entrypoints", "middlewares"]'
	);
	$: showColumn = (column: string): boolean => {
		return selectedColumns.includes(column);
	};
	const changeColumns = (columns: Selected<string>[] | undefined) => {
		if (columns === undefined) return;
		selectedColumns = columns.map((c) => c.value);
		localStorage.setItem('router-columns', JSON.stringify(selectedColumns));
	};

	const toggleEntrypoint = (router: Router, item: Selected<unknown>[] | undefined) => {
		if (item === undefined) return;
		router.entrypoints = item.map((i) => i.value) as string[];
		updateRouter($profile, router, router.name);
	};
	const toggleMiddleware = (router: Router, item: Selected<unknown>[] | undefined) => {
		if (item === undefined) return;
		router.middlewares = item.map((i) => i.value) as string[];
		updateRouter($profile, router, router.name);
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
		let service = $services.find((s) => s.name === `${router.service}@${router.provider}`);
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

	onMount(() => {
		search = localProvider ? '@provider:http' : '';
		searchRouter();
	});
</script>

<svelte:head>
	<title>Routers {$profile ? `| ${$profile}` : ''}</title>
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
					{#if showColumn('name')}
						<Table.Head>Name</Table.Head>
					{/if}
					{#if showColumn('provider')}
						<Table.Head>Provider</Table.Head>
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
						<Table.Cell class={showColumn('name') ? 'font-medium' : 'hidden'}>
							{router.service}
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
									on:click={() => deleteRouter($profile, router.name)}
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

<Pagination {count} bind:perPage bind:currentPage />
