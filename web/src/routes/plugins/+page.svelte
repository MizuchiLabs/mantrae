<script lang="ts">
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Delete } from 'lucide-svelte';
	import type { Plugin } from '$lib/types';
	import { profile, api, plugins } from '$lib/api';
	import { onMount } from 'svelte';
	import { Textarea } from '$lib/components/ui/textarea';
	import { toast } from 'svelte-sonner';
	import YAML from 'yaml';
	import type { HTTPMiddleware, UpsertMiddlewareParams } from '$lib/types/middlewares';

	// State
	let open = $state(false);
	let search = $state('');
	let currentPage = $state(1);
	let perPage = $state(10);
	let selectedPlugin = $state<Plugin | undefined>(undefined);
	let yamlSnippet = $state('');

	// Derived values
	let filteredPlugins = $derived(
		$plugins.filter(
			(plugin) =>
				!search ||
				plugin.name.toLowerCase().includes(search.toLowerCase()) ||
				plugin.displayName.toLowerCase().includes(search.toLowerCase())
		)
	);

	let paginatedPlugins = $derived(() => {
		const start = (currentPage - 1) * perPage;
		return [...filteredPlugins].slice(start, start + perPage);
	});

	// When search changes, reset to first page
	$effect(() => {
		if (search) currentPage = 1;
	});

	async function installPlugin(plugin: Plugin) {
		if (!$profile.id) return;

		const data = YAML.parse(plugin.snippet.yaml);
		const pluginContent = extractPluginContent(data);
		const name = Object.keys(pluginContent)[0];

		const middleware: UpsertMiddlewareParams = {
			name: `${name}@http`,
			protocol: 'http',
			type: 'plugin',
			middleware: {
				name: name,
				protocol: 'http',
				plugin: {
					[name]: pluginContent[name]
				}
			}
		};
		console.log(middleware);
		await api.upsertMiddleware($profile.id, middleware);

		selectedPlugin = plugin;
		yamlSnippet = generateYamlSnippet(plugin);
		open = true;
	}

	function extractPluginContent(data: Record<string, any>) {
		const middlewares = data.http?.middlewares;
		if (!middlewares) return null;

		const firstMiddleware = Object.values(middlewares)[0];
		return firstMiddleware?.plugin || null;
	}

	function generateYamlSnippet(plugin: Plugin) {
		return `experimental:
  plugins:
    ${plugin.name.split('/').slice(-1)[0]}:
      moduleName: ${plugin.name}
      version: ${plugin.latestVersion}`;
	}

	function copyToClipboard() {
		navigator.clipboard.writeText(yamlSnippet);
		toast.success('Copied!');
	}

	onMount(async () => {
		await api.getMiddlewarePlugins();
	});
</script>

<svelte:head>
	<title>Plugins</title>
</svelte:head>

<div class="mt-4 flex flex-col gap-4 p-4">
	<div class="flex flex-row items-center gap-1">
		<div class="relative flex flex-row items-center">
			<Input
				type="text"
				placeholder="Search..."
				class="w-80 focus-visible:ring-0 focus-visible:ring-offset-0"
				bind:value={search}
			/>

			<Button
				variant="ghost"
				class="absolute right-0 mr-1 rounded-full hover:bg-transparent"
				onclick={() => (search = '')}
				size="icon"
				aria-hidden
			>
				<Delete size="1.25rem" class="text-muted-foreground hover:text-red-400" />
			</Button>
		</div>
	</div>

	<!-- Plugin Grid -->
	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4">
		{#each paginatedPlugins() as plugin (plugin.name)}
			<Card.Root class="flex h-[300px] w-full flex-col">
				<Card.Header class="flex-grow">
					<Card.Title class="mb-2 flex flex-row items-center gap-4">
						<Avatar.Root class="h-12 w-12">
							<Avatar.Image src={plugin.iconUrl} alt={plugin.displayName} />
							<Avatar.Fallback>{plugin.displayName.slice(0, 2)}</Avatar.Fallback>
						</Avatar.Root>
						{plugin.displayName}
					</Card.Title>
					<Card.Description class="line-clamp-3">{plugin.summary}</Card.Description>
				</Card.Header>
				<Card.Content class="mt-auto flex flex-row items-center justify-between">
					<Badge variant="secondary">{plugin.latestVersion}</Badge>
					<Button onclick={() => installPlugin(plugin)}>Install</Button>
				</Card.Content>
			</Card.Root>
		{/each}
	</div>

	<!-- Pagination -->
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

<Dialog.Root bind:open>
	<Dialog.Content class="max-w-xl">
		<Dialog.Header>
			<Dialog.Title>Install {selectedPlugin?.displayName}</Dialog.Title>
			<Dialog.Description>
				Add this snippet to your Traefik Static Config:
				<span class="text-xs text-slate-500">(Click to copy)</span>
			</Dialog.Description>
		</Dialog.Header>
		<div class="flex flex-col gap-4">
			<Textarea
				bind:value={yamlSnippet}
				rows={yamlSnippet?.split('\n').length || 5}
				onclick={copyToClipboard}
				class="p-2"
				readonly
			/>
		</div>
	</Dialog.Content>
</Dialog.Root>
