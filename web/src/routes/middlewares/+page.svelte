<script lang="ts">
	import { deleteMiddleware, middlewares } from '$lib/api';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import Pagination from '$lib/components/tables/pagination.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { Selected } from 'bits-ui';
	import { newMiddleware, type Middleware } from '$lib/types/middlewares';
	import Search from '$lib/components/tables/search.svelte';
	import MiddlewareModal from '$lib/components/modals/middleware.svelte';
	import { Eye, Pencil, Plus, X } from 'lucide-svelte';
	import { LIMIT_SK, MIDDLEWARE_COLUMN_SK } from '$lib/store';
	import { profile, getMiddlewares } from '$lib/api';

	let search = '';
	let count = 0;
	let currentPage = 1;
	let fMiddlewares: Middleware[] = [];
	let perPage: Selected<number> | undefined = JSON.parse(
		localStorage.getItem(LIMIT_SK) ?? '{"value": 10, "label": "10"}'
	);
	$: search, $middlewares, currentPage, perPage, searchMiddleware();

	// Reset the page to 1 when the search input changes
	$: {
		if (search) {
			currentPage = 1;
		}
	}

	function searchMiddleware() {
		if ($middlewares === undefined) return;
		let items = $middlewares?.filter((middleware) => {
			const searchParts = search.toLowerCase().split(' ');
			return searchParts.every((part) =>
				part.startsWith('@provider:')
					? middleware.provider?.toLowerCase() === part.split(':')[1]
					: part.startsWith('@type:')
						? middleware.type?.toLowerCase() === part.split(':')[1]
						: middleware.name.toLowerCase().includes(part)
			);
		});

		count = items?.length || 1;
		fMiddlewares = paginate(items);
	}

	const paginate = (middlewares: Middleware[]) => {
		const itemsPerPage = perPage?.value ?? 10;
		return middlewares?.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage) ?? [];
	};

	let middleware = newMiddleware();
	let openModal = false;
	let disabled = false;
	const createModal = async () => {
		middleware = newMiddleware();
		disabled = false;
		openModal = true;
	};
	const updateModal = async (m: Middleware) => {
		if (m.provider === 'http') {
			disabled = false;
		} else {
			disabled = true;
		}
		middleware = m;
		openModal = true;
	};

	let columns: Selected<string>[] | undefined = [
		{ value: 'name', label: 'Name' },
		{ value: 'provider', label: 'Provider' },
		{ value: 'type', label: 'Type' }
	];
	let fColumns: string[] = JSON.parse(
		localStorage.getItem(MIDDLEWARE_COLUMN_SK) ?? JSON.stringify(columns.map((c) => c.value))
	);

	profile.subscribe((value) => {
		if (!value?.id) return;
		getMiddlewares();
	});
</script>

<svelte:head>
	<title>Middlewares</title>
</svelte:head>

<MiddlewareModal {middleware} bind:open={openModal} bind:disabled />

<div class="mt-4 flex flex-col gap-4 p-4">
	<Search bind:search {columns} columnName="middleware-columns" bind:fColumns />

	<Card.Root>
		<Card.Header class="grid grid-cols-2 items-center justify-between">
			<div>
				<Card.Title>Middlewares</Card.Title>
				<Card.Description
					>Total middlewares managed by Mantrae {$middlewares?.filter((m) => m.provider === 'http')
						.length}</Card.Description
				>
			</div>
			<div class="justify-self-end">
				<Button
					variant="secondary"
					class="flex items-center gap-2 bg-red-400 text-black"
					on:click={createModal}
				>
					<span>Add Middleware</span>
					<Plus size="1rem" />
				</Button>
			</div>
		</Card.Header>
		<Card.Content>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						{#if fColumns.includes('name')}
							<Table.Head>Name</Table.Head>
						{/if}
						{#if fColumns.includes('provider')}
							<Table.Head>Provider</Table.Head>
						{/if}
						{#if fColumns.includes('type')}
							<Table.Head class="hidden md:table-cell">Type</Table.Head>
						{/if}
						<Table.Head>
							<span class="sr-only">Delete</span>
						</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each fMiddlewares as middleware}
						<Table.Row>
							<Table.Cell
								class={fColumns.includes('name')
									? 'max-w-[180px] overflow-hidden text-ellipsis whitespace-nowrap'
									: 'hidden'}
							>
								{middleware.name.split('@')[0]}
							</Table.Cell>
							<Table.Cell class={fColumns.includes('provider') ? 'font-medium' : 'hidden'}>
								<span
									class="inline-flex cursor-pointer select-none items-center rounded-full bg-slate-300 px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
									on:click={() => (search = `@provider:${middleware.provider}`)}
									aria-hidden
								>
									{middleware.provider}
								</span>
							</Table.Cell>
							<Table.Cell class={fColumns.includes('type') ? 'font-medium' : 'hidden'}>
								<span
									class="inline-flex cursor-pointer select-none items-center rounded-full bg-slate-300 px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
									on:click={() => (search = `@type:${middleware.type}`)}
									aria-hidden
								>
									{middleware.type}
								</span>
							</Table.Cell>
							<Table.Cell class="min-w-[100px]">
								{#if middleware.provider === 'http'}
									<Button
										variant="ghost"
										class="h-8 w-8 rounded-full bg-orange-400"
										size="icon"
										on:click={() => updateModal(middleware)}
									>
										<Pencil size="1rem" />
									</Button>
									<Button
										variant="ghost"
										class="h-8 w-8 rounded-full bg-red-400"
										size="icon"
										on:click={() => deleteMiddleware(middleware)}
									>
										<X size="1rem" />
									</Button>
								{:else}
									<Button
										variant="ghost"
										class="h-8 w-8 rounded-full bg-green-400"
										size="icon"
										on:click={() => updateModal(middleware)}
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
				Showing <strong>{fMiddlewares?.length > 0 ? 1 : 0}-{fMiddlewares?.length ?? 0}</strong>
				of
				<strong>{$middlewares?.length ?? 0}</strong> middlewares
			</div>
		</Card.Footer>
	</Card.Root>

	<Pagination {count} bind:perPage bind:currentPage />
</div>
