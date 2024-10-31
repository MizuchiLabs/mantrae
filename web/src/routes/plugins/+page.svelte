<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Delete } from 'lucide-svelte';
	import type { Plugin } from '$lib/types/plugins';
	import { profile, getPlugins, plugins, upsertMiddleware } from '$lib/api';
	import { onMount } from 'svelte';
	import type { Selected } from 'bits-ui';
	import { LIMIT_SK } from '$lib/store';
	import Pagination from '$lib/components/tables/pagination.svelte';
	import { Textarea } from '$lib/components/ui/textarea';
	import { toast } from 'svelte-sonner';
	import YAML from 'yaml';
	import { newMiddleware } from '$lib/types/middlewares';

	let open = false;
	let search: string = '';
	let fPlugins: Plugin[] = [];
	let count = 0;
	let currentPage = 1;
	let perPage: Selected<number> | undefined = JSON.parse(
		localStorage.getItem(LIMIT_SK) ?? '{"value": 10, "label": "10"}'
	);

	// Reset the page to 1 when the search input changes
	$: search, (currentPage = 1);

	// Watch for changes in search or currentPage
	$: fPlugins = getFilteredPlugins($plugins, search);
	$: paginatedPlugins = paginate(fPlugins, currentPage, perPage?.value ?? 10);
	$: count = fPlugins?.length || 1;

	const getFilteredPlugins = (plugins: Plugin[], search: string) => {
		if (!search) return plugins; // Return all if no search
		return plugins.filter(
			(plugin) =>
				plugin.name.toLowerCase().includes(search.toLowerCase()) ||
				plugin.displayName.toLowerCase().includes(search.toLowerCase())
		);
	};

	const paginate = (plugins: Plugin[], page: number, itemsPerPage: number) => {
		const start = (page - 1) * itemsPerPage;
		return plugins?.slice(start, start + itemsPerPage);
	};

	let selectedPlugin: Plugin | undefined;
	let yamlSnippet: string = '';
	const installPlugin = async (plugin: Plugin) => {
		if (!$profile.id) return;
		const data = YAML.parse(plugin.snippet.yaml);
		const content = getPluginContent(data);
		let name = Object.keys(content)[0];
		if (name.startsWith('my-')) {
			name = name.slice(3);
		}
		let middleware = newMiddleware();
		middleware.name = name;
		middleware.content = Object.values(content)[0];
		middleware.profileId = $profile.id;
		middleware.provider = 'http';
		middleware.protocol = 'http';
		middleware.type = 'plugin';
		await upsertMiddleware(middleware);
		selectedPlugin = plugin;
		yamlSnippet = `experimental:
  plugins:
    ${plugin.name.split('/').slice(-1)[0]}:
      moduleName: ${plugin.name}
      version: ${plugin.latestVersion}`;
		open = true;
	};

	function getPluginContent(data: any) {
		let content: Record<string, Record<string, any>> = {};

		// Iterate through HTTP middlewares
		if (data.http && data.http.middlewares) {
			Object.entries(data.http.middlewares).forEach(([name, middleware]: any) => {
				if (middleware.plugin) {
					// There should only be one plugin key, which is the last part we need
					const pluginNames = Object.keys(middleware.plugin);
					if (pluginNames.length > 0) {
						const lastPluginName = pluginNames[pluginNames.length - 1];
						content[name] = middleware.plugin[lastPluginName]; // Get the content of the last plugin
					}
				}
			});
		}

		return content;
	}

	const copyToClipboard = () => {
		navigator.clipboard.writeText(yamlSnippet);
		toast.success('Copied!');
	};

	onMount(async () => {
		await getPlugins();
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
				on:click={() => (search = '')}
				size="icon"
				aria-hidden
			>
				<Delete size="1.25rem" class="text-muted-foreground hover:text-red-400" />
			</Button>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
		{#each paginatedPlugins || [] as plugin}
			<Card.Root class="flex h-[300px] w-full flex-col">
				<Card.Header class="flex-grow">
					<Card.Title class="mb-2 flex flex-row items-center gap-4">
						<Avatar.Root class="h-12 w-12">
							<Avatar.Image src={plugin.iconUrl} alt="@shadcn" />
							<Avatar.Fallback>{plugin.displayName.slice(0, 2)}</Avatar.Fallback>
						</Avatar.Root>
						{plugin.displayName}
					</Card.Title>
					<Card.Description class="line-clamp-3 overflow-hidden text-ellipsis"
						>{plugin.summary}</Card.Description
					>
				</Card.Header>
				<Card.Content class="mt-auto flex flex-row items-center justify-between">
					<Badge variant="secondary">{plugin.latestVersion}</Badge>
					<Button variant="default" on:click={() => installPlugin(plugin)}>Install</Button>
				</Card.Content>
			</Card.Root>
		{/each}
	</div>

	<Pagination {count} bind:perPage bind:currentPage />
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
				on:click={copyToClipboard}
				class="p-2"
				readonly
			/>
		</div>
	</Dialog.Content>
</Dialog.Root>
