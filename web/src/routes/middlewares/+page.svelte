<script lang="ts">
	import { activeProfile, deleteMiddleware } from '$lib/api';
	import CreateMiddleware from '$lib/components/modals/createMiddleware.svelte';
	import Pagination from '$lib/components/tables/pagination.svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import type { Selected } from 'bits-ui';
	import UpdateMiddleware from '$lib/components/modals/updateMiddleware.svelte';
	import type { Middleware } from '$lib/types/middlewares';
	import { Input } from '$lib/components/ui/input';
	import { onMount } from 'svelte';

	let search = '';
	let count = 0;
	let currentPage = 1;
	let fMiddlewares: Middleware[] = [];
	let perPage: Selected<number> | undefined = { value: 10, label: '10' }; // Items per page

	$: middlewares = Object.values($activeProfile?.instance?.dynamic?.middlewares ?? []);
	$: search, middlewares, currentPage, searchMiddleware();

	// Reset the page to 1 when the search input changes
	$: {
		if (search) {
			currentPage = 1;
		}
	}

	function searchMiddleware() {
		let items: Middleware[] = [...middlewares];

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

			items = items.filter((middleware) => {
				const matchesProvider = providerSearch
					? middleware.provider?.toLowerCase() === providerSearch
					: true;
				const matchesType = typeSearch ? middleware.type?.toLowerCase() === typeSearch : true;
				const matchesGeneral = generalSearch.every((term) =>
					middleware.name.toLowerCase().includes(term)
				);
				return matchesProvider && matchesType && matchesGeneral;
			});
		}

		fMiddlewares = paginate(items);
		count = items.length || 1;
	}

	const paginate = (middlewares: Middleware[]) => {
		const itemsPerPage = perPage?.value ?? 10;
		return middlewares.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);
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
	<title>Middlewares | {$activeProfile?.name}</title>
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
</div>

<Card.Root>
	<Card.Header class="grid grid-cols-2 items-center justify-between">
		<div>
			<Card.Title>Middlewares</Card.Title>
			<Card.Description>Manage your Middlewares.</Card.Description>
		</div>
		<div class="justify-self-end">
			<CreateMiddleware />
		</div>
	</Card.Header>
	<Card.Content>
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Name</Table.Head>
					<Table.Head>Provider</Table.Head>
					<Table.Head class="hidden md:table-cell">Type</Table.Head>
					<Table.Head>
						<span class="sr-only">Delete</span>
					</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each fMiddlewares as middleware}
					<Table.Row>
						<Table.Cell class="max-w-[180px] overflow-hidden text-ellipsis whitespace-nowrap">
							{middleware.name}
						</Table.Cell>
						<Table.Cell class="font-medium">
							<span
								class="inline-flex cursor-pointer select-none items-center rounded-full bg-slate-300 px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
								on:click={() => (search = `@provider:${middleware.provider}`)}
								aria-hidden
							>
								{middleware.provider}
							</span>
						</Table.Cell>
						<Table.Cell class="font-medium">
							<span
								class="inline-flex cursor-pointer select-none items-center rounded-full bg-slate-300 px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
								on:click={() => (search = `@type:${middleware.type}`)}
								aria-hidden
							>
								{middleware.type}
							</span>
						</Table.Cell>
						{#if middleware.provider === 'http'}
							<Table.Cell class="min-w-[100px]">
								<UpdateMiddleware {middleware} />
								<Button
									variant="ghost"
									class="h-8 w-4 rounded-full bg-red-400"
									on:click={() => deleteMiddleware($activeProfile.name, middleware.name)}
								>
									<iconify-icon icon="fa6-solid:xmark" />
								</Button>
							</Table.Cell>
						{/if}
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card.Content>
	<Card.Footer>
		<div class="text-xs text-muted-foreground">
			Showing <strong>{fMiddlewares.length > 0 ? 1 : 0}-{fMiddlewares.length}</strong>
			of
			<strong>{middlewares.length}</strong> middlewares
		</div>
	</Card.Footer>
</Card.Root>

<Pagination {count} {perPage} bind:currentPage />
