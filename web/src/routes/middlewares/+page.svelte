<script lang="ts">
	import { activeProfile, deleteMiddleware } from '$lib/api';
	import CreateMiddleware from '$lib/components/modals/createMiddleware.svelte';
	import Pagination from '$lib/components/tables/pagination.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import type { Selected } from 'bits-ui';
	import UpdateMiddleware from '$lib/components/modals/updateMiddleware.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	$: httpmiddlewares = $activeProfile?.instance?.dynamic?.httpmiddlewares ?? [];
	$: tcpmiddlewares = $activeProfile?.instance?.dynamic?.tcpmiddlewares ?? [];
	$: middlewares = [
		...httpmiddlewares.map((m) => ({ ...m, middlewareType: 'http' })),
		...tcpmiddlewares.map((m) => ({ ...m, middlewareType: 'tcp' }))
	];
	$: paginatedMiddlewares = middlewares.slice(
		(currentPage - 1) * perPage?.value!,
		currentPage * perPage?.value!
	);
	let perPage: Selected<number> | undefined = { value: 10, label: '10' }; // Items per page
	let currentPage = 1; // Current page

	const isHttp = (name: string) => name.includes('@http');
</script>

<svelte:head>
	<title>Middlewares | {$activeProfile?.name}</title>
</svelte:head>

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
				{#each paginatedMiddlewares as middleware}
					<Table.Row>
						<Table.Cell class="max-w-[180px] overflow-hidden text-ellipsis whitespace-nowrap">
							{middleware.name}
						</Table.Cell>
						<Table.Cell class="font-medium">
							<Badge variant="secondary" class="text-xs">
								{middleware.provider}
							</Badge>
						</Table.Cell>
						<Table.Cell class="font-medium">
							<Badge variant="secondary" class="text-xs">
								{middleware.type}
							</Badge>
						</Table.Cell>
						{#if middleware.provider === 'http'}
							<Table.Cell class="min-w-[100px]">
								{#if middleware.middlewareType === 'http'}
									<UpdateMiddleware httpMiddleware={middleware} tcpMiddleware={undefined} />
								{:else}
									<UpdateMiddleware httpMiddleware={undefined} tcpMiddleware={middleware} />
								{/if}
								<Button
									variant="ghost"
									class="h-8 w-4 rounded-full bg-red-400"
									on:click={() => deleteMiddleware(middleware.name)}
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
			Showing <strong
				>{paginatedMiddlewares.length > 0 ? 1 : 0}-{paginatedMiddlewares.length}</strong
			>
			of
			<strong>{middlewares.length}</strong> middlewares
		</div>
	</Card.Footer>
</Card.Root>

<Pagination
	count={middlewares.length > 0 ? middlewares.length : 1}
	{perPage}
	bind:currentPage
	on:changeLimit={(e) => (perPage = e.detail)}
/>
