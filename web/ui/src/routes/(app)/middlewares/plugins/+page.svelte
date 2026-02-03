<script lang="ts">
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { Search, Package, Terminal, Info, X } from '@lucide/svelte';
	import { Textarea } from '$lib/components/ui/textarea';
	import YAML from 'yaml';
	import { fade } from 'svelte/transition';
	import CopyButton from '$lib/components/ui/copy-button/copy-button.svelte';
	import { type Plugin } from '$lib/gen/mantrae/v1/middleware_pb';
	import { marshalConfig } from '$lib/types';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { middleware } from '$lib/api/middleware.svelte';
	import { toast } from 'svelte-sonner';

	// State
	let open = $state(false);
	let search = $state('');
	let currentPage = $state(1);
	let perPage = $state(12);
	let selectedPlugin = $state<Plugin | undefined>(undefined);
	let yamlSnippet = $state('');
	let cliSnippet = $state('');

	const plugins = middleware.plugins();

	// Derived values
	let filteredPlugins = $derived(
		plugins.data?.filter(
			(plugin) =>
				!search ||
				plugin.name.toLowerCase().includes(search.toLowerCase()) ||
				plugin.displayName.toLowerCase().includes(search.toLowerCase())
		) || []
	);

	let totalPages = $derived(Math.ceil(filteredPlugins.length / perPage));

	let paginatedPlugins = $derived(() => {
		const start = (currentPage - 1) * perPage;
		return filteredPlugins.slice(start, start + perPage);
	});

	// When search changes, reset to first page
	$effect(() => {
		if (search) currentPage = 1;
	});

	const createMutation = middleware.create();

	async function installPlugin(plugin: Plugin) {
		if (!plugin.snippet) return;

		try {
			const data = YAML.parse(plugin.snippet.yaml);
			const middlewares = data.http?.middlewares;

			if (!middlewares) {
				toast.error('Invalid plugin snippet format');
				return;
			}

			const pluginContent = Object.values(middlewares)[0] as { plugin?: Record<string, any> };
			const name = Object.keys(pluginContent?.plugin || {})[0];

			if (!name) {
				toast.error('Could not determine plugin name');
				return;
			}

			createMutation.mutate({
				type: ProtocolType.HTTP,
				name: name,
				config: marshalConfig(pluginContent)
			});

			selectedPlugin = plugin;
			yamlSnippet = generateYamlSnippet(plugin);
			cliSnippet = generateCmdSnippet(plugin);
			open = true;
		} catch (e) {
			console.error(e);
			toast.error('Failed to prepare plugin installation');
		}
	}

	function generateYamlSnippet(plugin: Plugin) {
		const pluginName = plugin.name.split('/').slice(-1)[0];
		return `experimental:
  plugins:
    ${pluginName}:
      moduleName: ${plugin.name}
      version: ${plugin.latestVersion}`;
	}

	function generateCmdSnippet(plugin: Plugin) {
		const pluginName = plugin.name.split('/').slice(-1)[0];
		return `--experimental.plugins.${pluginName}.moduleName=${plugin.name}
--experimental.plugins.${pluginName}.version=${plugin.latestVersion}`;
	}
</script>

<svelte:head>
	<title>Plugin Marketplace - Mantrae</title>
	<meta
		name="description"
		content="Browse and manage Traefik plugins for extended reverse proxy functionality"
	/>
</svelte:head>

<div class="mx-auto flex w-full max-w-[1600px] flex-col gap-6 p-6 md:p-8">
	<!-- Header Section -->
	<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
		<div class="space-y-1">
			<h1 class="text-2xl font-bold tracking-tight">Plugin Marketplace</h1>
			<p class="text-sm text-muted-foreground">
				Discover and install plugins to extend Traefik's capabilities.
			</p>
		</div>

		<div class="relative w-full md:w-80">
			<Search class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
			<Input type="text" placeholder="Search plugins..." class="pr-9 pl-9" bind:value={search} />
			{#if search}
				<Button
					variant="ghost"
					size="icon"
					class="absolute top-0 right-0 h-full px-3 hover:bg-transparent"
					onclick={() => (search = '')}
				>
					<span class="sr-only">Clear search</span>
					<X class="h-4 w-4" />
				</Button>
			{/if}
		</div>
	</div>

	{#if plugins.isLoading}
		<!-- Loading State -->
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each Array(8) as _}
				<div class="flex flex-col gap-3 rounded-xl border p-4">
					<div class="flex items-center gap-3">
						<Skeleton class="h-12 w-12 rounded-full" />
						<div class="space-y-2">
							<Skeleton class="h-4 w-24" />
							<Skeleton class="h-3 w-16" />
						</div>
					</div>
					<Skeleton class="h-16 w-full" />
					<div class="mt-auto flex justify-between pt-2">
						<Skeleton class="h-5 w-16" />
						<Skeleton class="h-9 w-20" />
					</div>
				</div>
			{/each}
		</div>
	{:else if plugins.error}
		<!-- Error State -->
		<div
			class="flex h-100 flex-col items-center justify-center gap-4 rounded-xl border border-dashed bg-muted/30 text-center"
		>
			<div class="rounded-full bg-red-100 p-3 dark:bg-red-900/20">
				<Info class="h-6 w-6 text-red-600 dark:text-red-400" />
			</div>
			<div class="space-y-1">
				<h3 class="font-semibold">Failed to load plugins</h3>
				<p class="text-sm text-muted-foreground">
					Something went wrong while fetching the plugin list.
				</p>
			</div>
			<Button variant="outline" onclick={() => plugins.refetch()}>Try Again</Button>
		</div>
	{:else if filteredPlugins.length === 0}
		<!-- Empty State -->
		<div
			class="flex h-100 flex-col items-center justify-center gap-4 rounded-xl border border-dashed bg-muted/30 text-center"
		>
			<div class="rounded-full bg-muted p-3">
				<Package class="h-6 w-6 text-muted-foreground" />
			</div>
			<div class="space-y-1">
				<h3 class="font-semibold">No plugins found</h3>
				<p class="text-sm text-muted-foreground">
					Try adjusting your search terms or browse all plugins.
				</p>
			</div>
			{#if search}
				<Button variant="secondary" onclick={() => (search = '')}>Clear Search</Button>
			{/if}
		</div>
	{:else}
		<!-- Plugin Grid -->
		<div
			class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
			transition:fade
		>
			{#each paginatedPlugins() as plugin (plugin.name)}
				<Card.Root
					class="group flex h-full flex-col overflow-hidden transition-all hover:border-primary/50 hover:shadow-md"
				>
					<Card.Header class="space-y-4">
						<div class="flex items-start justify-between gap-4">
							<div class="flex items-center gap-3">
								<Avatar.Root class="h-10 w-10 border bg-background">
									<Avatar.Image src={plugin.iconUrl} alt={plugin.displayName} />
									<Avatar.Fallback class="font-semibold text-primary">
										{plugin.displayName.slice(0, 2).toUpperCase()}
									</Avatar.Fallback>
								</Avatar.Root>
								<div>
									<Card.Title class="text-base leading-tight font-semibold">
										{plugin.displayName}
									</Card.Title>
									<p class="text-xs text-muted-foreground" title={plugin.name}>
										by {plugin.author || 'Unknown'}
									</p>
								</div>
							</div>
						</div>

						<Card.Description class="line-clamp-3 h-16 text-sm">
							{plugin.summary || 'No description available.'}
						</Card.Description>
					</Card.Header>

					<Card.Content class="mt-auto flex items-center justify-between border-t px-6 pt-5">
						<Badge variant="secondary" class="font-mono text-xs font-normal">
							{plugin.latestVersion}
						</Badge>
						<Button size="sm" onclick={() => installPlugin(plugin)}>Install</Button>
					</Card.Content>
				</Card.Root>
			{/each}
		</div>

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="flex justify-center py-4">
				<Pagination.Root
					count={filteredPlugins.length}
					{perPage}
					onPageChange={(page) => (currentPage = page)}
				>
					{#snippet children({ pages, currentPage })}
						<Pagination.Content>
							<Pagination.Item>
								<Pagination.PrevButton />
							</Pagination.Item>
							{#each pages as page (page.key)}
								{#if page.type === 'ellipsis'}
									<Pagination.Item>
										<Pagination.Ellipsis />
									</Pagination.Item>
								{:else}
									<Pagination.Item>
										<Pagination.Link {page} isActive={currentPage === page.value}>
											{page.value}
										</Pagination.Link>
									</Pagination.Item>
								{/if}
							{/each}
							<Pagination.Item>
								<Pagination.NextButton />
							</Pagination.Item>
						</Pagination.Content>
					{/snippet}
				</Pagination.Root>
			</div>
		{/if}
	{/if}
</div>

<!-- Installation Dialog -->
<Dialog.Root bind:open>
	<Dialog.Content class="max-w-2xl gap-0 overflow-hidden p-0 outline-none">
		<Dialog.Header class="border-b bg-muted/30 px-6 py-4">
			<div class="flex items-center gap-3">
				<div class="rounded-lg border bg-background p-2">
					<Avatar.Root class="h-8 w-8">
						<Avatar.Image src={selectedPlugin?.iconUrl} alt={selectedPlugin?.displayName} />
						<Avatar.Fallback>{selectedPlugin?.displayName?.slice(0, 2)}</Avatar.Fallback>
					</Avatar.Root>
				</div>
				<div class="space-y-1">
					<Dialog.Title>Configure {selectedPlugin?.displayName}</Dialog.Title>
					<Dialog.Description>
						Middleware created in Mantrae. Now update your Traefik static configuration.
					</Dialog.Description>
				</div>
			</div>
		</Dialog.Header>

		<div class="space-y-4 p-6">
			<div
				class="rounded-md bg-blue-50 p-3 text-sm text-blue-900 dark:bg-blue-900/20 dark:text-blue-200"
			>
				<div class="flex gap-2">
					<Info class="mt-0.5 h-4 w-4 shrink-0" />
					<p>
						To enable this plugin, you must add one of the following snippets to your Traefik static
						configuration file (traefik.yml or command line flags) and restart Traefik.
					</p>
				</div>
			</div>

			<Tabs.Root value="yaml" class="w-full">
				<Tabs.List class="mb-4 grid w-full grid-cols-2">
					<Tabs.Trigger value="yaml" class="gap-2">
						<span class="font-mono text-xs">traefik.yml</span>
					</Tabs.Trigger>
					<Tabs.Trigger value="cli" class="gap-2">
						<Terminal />
						CLI Flags
					</Tabs.Trigger>
				</Tabs.List>

				<Tabs.Content value="yaml" class="group relative mt-0">
					<div class="relative rounded-md border bg-muted/50">
						<Textarea
							value={yamlSnippet}
							rows={yamlSnippet?.split('\n').length + 2}
							readonly
							class="min-h-[120px] resize-none border-0 bg-transparent font-mono text-sm focus-visible:ring-0"
						/>
						<div class="absolute top-2 right-2 rounded-lg bg-muted">
							<CopyButton text={yamlSnippet} />
						</div>
					</div>
				</Tabs.Content>

				<Tabs.Content value="cli" class="group relative mt-0">
					<div class="relative rounded-md border bg-muted/50">
						<Textarea
							bind:value={cliSnippet}
							class="min-h-[120px] resize-none border-0 bg-transparent font-mono text-sm break-all focus-visible:ring-0"
							rows={cliSnippet?.split('\n').length + 2}
							readonly
						/>
						<div class="absolute top-2 right-2 rounded-lg bg-muted">
							<CopyButton text={cliSnippet} />
						</div>
					</div>
				</Tabs.Content>
			</Tabs.Root>
		</div>
	</Dialog.Content>
</Dialog.Root>
