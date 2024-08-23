<script lang="ts">
	import { profile, deleteMiddleware, middlewares } from '$lib/api';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Select from '$lib/components/ui/select';
	import CreateMiddleware from '$lib/components/modals/createMiddleware.svelte';
	import Pagination from '$lib/components/tables/pagination.svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import UpdateMiddleware from '$lib/components/modals/updateMiddleware.svelte';
	import type { Selected } from 'bits-ui';
	import type { Middleware } from '$lib/types/middlewares';
	import { Input } from '$lib/components/ui/input';
	import { onMount } from 'svelte';

	let search = '';
	let count = 0;
	let currentPage = 1;
	let fMiddlewares: Middleware[] = [];
	let perPage: Selected<number> | undefined = JSON.parse(
		localStorage.getItem('limit') ?? '{"value": 10, "label": "10"}'
	);
	$: search, $middlewares, currentPage, searchMiddleware();

	// Reset the page to 1 when the search input changes
	$: {
		if (search) {
			currentPage = 1;
		}
	}

	function searchMiddleware() {
		let items = $middlewares.filter((middleware) => {
			if (localProvider && middleware.provider !== 'http') return false;
			const searchParts = search.toLowerCase().split(' ');
			return searchParts.every((part) =>
				part.startsWith('@provider:')
					? middleware.provider?.toLowerCase() === part.split(':')[1]
					: part.startsWith('@type:')
						? middleware.type?.toLowerCase() === part.split(':')[1]
						: middleware.name.toLowerCase().includes(part)
			);
		});

		count = items.length || 1;
		fMiddlewares = paginate(items);
	}

	const paginate = (middlewares: Middleware[]) => {
		const itemsPerPage = perPage?.value ?? 10;
		return middlewares.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);
	};

	let columns: Selected<string>[] | undefined = [
		{ value: 'name', label: 'Name' },
		{ value: 'provider', label: 'Provider' },
		{ value: 'type', label: 'Type' }
	];
	let selectedColumns: string[] = JSON.parse(
		localStorage.getItem('middleware-columns') ?? '["name", "provider", "type"]'
	);
	$: showColumn = (column: string): boolean => {
		return selectedColumns.includes(column);
	};
	const changeColumns = (columns: Selected<string>[] | undefined) => {
		if (columns === undefined) return;
		selectedColumns = columns.map((c) => c.value);
		localStorage.setItem('middleware-columns', JSON.stringify(selectedColumns));
	};

	// Only show local middlewares not external ones
	let localProvider = localStorage.getItem('local-provider') === 'true';
	const toggleProvider = () => {
		localProvider = !localProvider;
		search = localProvider ? '@provider:http' : '';
		localStorage.setItem('local-provider', localProvider.toString());
	};

	onMount(() => {
		search = localProvider ? '@provider:http' : '';
		searchMiddleware();
	});
</script>

<svelte:head>
	<title>Middlewares | {$profile}</title>
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
					{#if showColumn('name')}
						<Table.Head>Name</Table.Head>
					{/if}
					{#if showColumn('provider')}
						<Table.Head>Provider</Table.Head>
					{/if}
					{#if showColumn('type')}
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
							class={showColumn('name')
								? 'max-w-[180px] overflow-hidden text-ellipsis whitespace-nowrap'
								: 'hidden'}
						>
							{middleware.name.split('@')[0]}
						</Table.Cell>
						<Table.Cell class={showColumn('provider') ? 'font-medium' : 'hidden'}>
							<span
								class="inline-flex cursor-pointer select-none items-center rounded-full bg-slate-300 px-2.5 py-0.5 text-xs font-semibold text-slate-800 hover:bg-red-300 focus:outline-none"
								on:click={() => (search = `@provider:${middleware.provider}`)}
								aria-hidden
							>
								{middleware.provider}
							</span>
						</Table.Cell>
						<Table.Cell class={showColumn('type') ? 'font-medium' : 'hidden'}>
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
									on:click={() => deleteMiddleware($profile, middleware.name)}
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
			<strong>{$middlewares.length}</strong> middlewares
		</div>
	</Card.Footer>
</Card.Root>

<Pagination {count} bind:perPage bind:currentPage />
