<script lang="ts">
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Delete } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { Textarea } from '$lib/components/ui/textarea';
	import { toast } from 'svelte-sonner';
	import YAML from 'yaml';
	import { slide } from 'svelte/transition';
	import CopyButton from '$lib/components/ui/copy-button/copy-button.svelte';
	import { MiddlewareType, type Plugin } from '$lib/gen/mantrae/v1/middleware_pb';
	import { middlewareClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { profile } from '$lib/stores/profile';
	import { marshalConfig } from '$lib/types';

	// State
	let open = $state(false);
	let search = $state('');
	let currentPage = $state(1);
	let perPage = $state(10);
	let plugins = $state([] as Plugin[]);
	let selectedPlugin = $state<Plugin | undefined>(undefined);
	let yamlSnippet = $state('');
	let cliSnippet = $state('');

	// Derived values
	let filteredPlugins = $derived(
		plugins.filter(
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
		if (!plugin.snippet) return;
		const data = YAML.parse(plugin.snippet.yaml);

		const middlewares = data.http?.middlewares;
		if (!middlewares) return;
		const pluginContent = Object.values(middlewares)[0];
		const name = Object.keys(pluginContent?.plugin || {})[0];

		try {
			await middlewareClient.createMiddleware({
				profileId: profile.value?.id,
				type: MiddlewareType.HTTP,
				name: name,
				config: marshalConfig(pluginContent)
			});
			toast.success('Installed plugin!', {
				description: `The plugin ${name} has been added to the middlewares.`
			});
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to install plugin', { description: e.message });
		}

		selectedPlugin = plugin;
		yamlSnippet = generateYamlSnippet(plugin);
		cliSnippet = generateCmdSnippet(plugin);
		open = true;
	}

	function generateYamlSnippet(plugin: Plugin) {
		return `experimental:
  plugins:
    ${plugin.name.split('/').slice(-1)[0]}:
      moduleName: ${plugin.name}
      version: ${plugin.latestVersion}`;
	}

	function generateCmdSnippet(plugin: Plugin) {
		return `--experimental.plugins.${plugin.name.split('/').slice(-1)[0]}.moduleName=${plugin.name}
--experimental.plugins.${plugin.name.split('/').slice(-1)[0]}.version=${plugin.latestVersion}`;
	}

	onMount(async () => {
		const response = await middlewareClient.getMiddlewarePlugins({});
		plugins = response.plugins;
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
				<Card.Header class="grow">
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
	<Dialog.Content class="no-scrollbar max-h-[95vh] max-w-2xl overflow-auto">
		<Dialog.Header>
			<Dialog.Title>Install {selectedPlugin?.displayName}</Dialog.Title>
			<Dialog.Description>
				Add this snippet to your Traefik Static Config:
				<span class="text-xs text-slate-500">(Click to copy)</span>
			</Dialog.Description>
		</Dialog.Header>
		<Tabs.Root class="flex flex-col gap-2">
			<div class="flex justify-end" transition:slide={{ duration: 200 }}>
				<Tabs.List class="h-8">
					<Tabs.Trigger value="yaml" class="px-2 py-0.5 font-bold">YAML</Tabs.Trigger>
					<Tabs.Trigger value="cli" class="px-2 py-0.5 font-bold">CLI</Tabs.Trigger>
				</Tabs.List>
			</div>
			<Tabs.Content value="yaml" class="relative">
				<Textarea value={yamlSnippet} rows={yamlSnippet?.split('\n').length || 5} readonly />
				<CopyButton text={yamlSnippet} class="absolute top-1 right-1" />
			</Tabs.Content>
			<Tabs.Content value="cli" class="relative overflow-x-auto">
				<Textarea
					bind:value={cliSnippet}
					class="break-all whitespace-pre-wrap"
					rows={cliSnippet?.split('\n').length || 2}
					readonly
				/>
				<CopyButton text={cliSnippet} class="bg-background absolute top-1 right-1" />
			</Tabs.Content>
		</Tabs.Root>
	</Dialog.Content>
</Dialog.Root>
