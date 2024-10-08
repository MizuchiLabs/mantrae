<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { getTraefikConfig, config, dynamic } from '$lib/api';
	import { onMount } from 'svelte';
	import { darkMode } from '$lib/utils';
	import Highlight, { LineNumbers } from 'svelte-highlight';
	import yaml from 'svelte-highlight/languages/yaml';
	import github from 'svelte-highlight/styles/github';
	import githubDark from 'svelte-highlight/styles/github-dark';
	import { Copy, CopyCheck } from 'lucide-svelte';

	let copyText = 'Copy';
	const copy = () => {
		navigator.clipboard.writeText($dynamic);
		copyText = 'Copied!';
		setTimeout(() => {
			copyText = 'Copy';
		}, 2000);
	};

	onMount(async () => {
		await getTraefikConfig();
	});
</script>

<svelte:head>
	{#if $darkMode}
		{@html githubDark}
	{:else}
		{@html github}
	{/if}
</svelte:head>

<Dialog.Root>
	<Dialog.Trigger>
		<Button variant="ghost" on:click={getTraefikConfig}>
			<iconify-icon icon="devicon:traefikproxy" width="24" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<Tabs.Root value="overview" class="mt-4 max-w-2xl">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
				<Tabs.Trigger value="config">Config</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="overview">
				<Card.Root>
					<Card.Content class="flex flex-col gap-2">
						<span class="mt-4 border-b border-gray-200 pb-2 font-bold">Traefik Information</span>

						<!-- Version -->
						<div class="mt-2 grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1">Version</span>
							<div class="col-span-3 space-x-2">
								{#if $config?.version}
									<Badge variant="secondary" class="bg-blue-300">
										v{$config?.version}
									</Badge>
								{/if}
							</div>
						</div>

						<!-- Entrypoints -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1">Entrypoints</span>
							<div class="col-span-3 space-x-2">
								{#each $config.entrypoints ?? [] as entrypoint}
									{#if entrypoint.asDefault}
										<Badge variant="secondary" class="bg-green-300"
											>{entrypoint.name}{entrypoint.address}</Badge
										>
									{:else}
										<Badge variant="secondary">{entrypoint.name}{entrypoint.address}</Badge>
									{/if}
								{/each}
							</div>
						</div>

						<!-- Features -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1">Features</span>
							<div class="col-span-3 space-x-2">
								{#if $config.overview?.features.tracing}
									<Badge variant="secondary">Tracing</Badge>
								{/if}
								{#if $config.overview?.features.metrics}
									<Badge variant="secondary">Metrics</Badge>
								{/if}
								{#if $config.overview?.features.accessLog}
									<Badge variant="secondary">Access Log</Badge>
								{/if}
							</div>
						</div>

						<!-- Providers -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1">Providers</span>
							<div class="col-span-3 space-x-2">
								{#each $config.overview?.providers ?? [] as provider}
									{#if provider === 'http'}
										<Badge variant="secondary" class="bg-yellow-300">{provider}</Badge>
									{:else}
										<Badge variant="secondary">{provider}</Badge>
									{/if}
								{/each}
							</div>
						</div>

						<span class="mt-2 border-b border-gray-200 pb-2 font-bold"> Router Overview </span>

						<!-- HTTP Overview -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1 font-mono">HTTP</span>
							<div class="col-span-3 space-x-2">
								<Badge variant="secondary">
									Routers: {$config?.overview?.http.routers.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Services: {$config?.overview?.http.services.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Middlewares: {$config?.overview?.http.middlewares.total ?? 0}
								</Badge>
							</div>
						</div>

						<!-- TCP Overview -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1 font-mono">TCP</span>
							<div class="col-span-3 space-x-2">
								<Badge variant="secondary">
									Routers: {$config?.overview?.tcp.routers.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Services: {$config?.overview?.tcp.services.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Middlewares: {$config?.overview?.tcp.middlewares.total ?? 0}
								</Badge>
							</div>
						</div>

						<!-- UDP Overview -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1 font-mono">UDP</span>
							<div class="col-span-3 space-x-2">
								<Badge variant="secondary">
									Routers: {$config?.overview?.udp.routers.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Services: {$config?.overview?.udp.services.total ?? 0}
								</Badge>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
			<Tabs.Content value="config">
				<Card.Root>
					<Card.Header>
						<Card.Title class="flex items-center justify-between gap-2">
							Dynamic Config
							<button
								on:click={copy}
								class="flex flex-row items-center gap-2 rounded p-2 text-sm font-medium hover:bg-gray-100"
							>
								{copyText}
								{#if copyText === 'Copied!'}
									<CopyCheck size="1rem" />
								{:else}
									<Copy size="1rem" />
								{/if}
							</button>
						</Card.Title>
						<Card.Description>
							This is the current dynamic configuration your Traefik instance is using.
						</Card.Description>
					</Card.Header>
					<Card.Content class="text-sm">
						<Highlight code={$dynamic} language={yaml} let:highlighted>
							<LineNumbers {highlighted} hideBorder wrapLines />
						</Highlight>
						<!-- <Textarea -->
						<!-- 	value={$dynamic} -->
						<!-- 	rows={30} -->
						<!-- 	class="bg-gray-50 focus-visible:ring-0 focus-visible:ring-offset-0" -->
						<!-- 	on:click={(e) => e.target?.select()} -->
						<!-- 	readonly -->
						<!-- /> -->
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
		</Tabs.Root>
		<Dialog.Close class="w-full">
			<Button class="w-full">Close</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
