<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { entrypoints, overview, version, dynamicJSON, dynamicYAML, api } from '$lib/api';
	import Highlight, { LineNumbers } from 'svelte-highlight';
	import { json, yaml } from 'svelte-highlight/languages';
	import CopyButton from '../ui/copy-button/copy-button.svelte';
	import { profile } from '$lib/stores/profile';

	let isYaml = $state(false);
	let { open = $bindable(false) } = $props();

	$effect(() => {
		if (profile.isValid()) {
			api.getDynamicConfig();
		}
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-[850px] overflow-y-auto">
		<Tabs.Root value="overview" class="mt-4">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
				<Tabs.Trigger value="config">Config</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="overview">
				<Card.Root>
					<Card.Content class="flex flex-col gap-4">
						<span class="border-b border-gray-200 pb-2 font-bold">Traefik Information</span>

						<!-- Version -->
						<div class="mt-2 grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1">Version</span>
							<div class="col-span-3 space-x-2">
								{#if $version}
									<Badge variant="secondary" class="bg-blue-300 dark:bg-blue-600">
										v{$version}
									</Badge>
								{/if}
							</div>
						</div>

						<!-- Entrypoints -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1">Entrypoints</span>
							<div class="col-span-3 space-x-2">
								{#each $entrypoints ?? [] as entrypoint (entrypoint.name)}
									{#if entrypoint.asDefault}
										<Tooltip.Provider>
											<Tooltip.Root delayDuration={100}>
												<Tooltip.Trigger>
													<Badge variant="secondary" class="bg-green-300 dark:bg-green-600">
														{entrypoint.name}
													</Badge>
												</Tooltip.Trigger>
												<Tooltip.Content>
													{entrypoint.address}
												</Tooltip.Content>
											</Tooltip.Root>
										</Tooltip.Provider>
									{:else}
										<Tooltip.Provider>
											<Tooltip.Root delayDuration={100}>
												<Tooltip.Trigger>
													<Badge variant="secondary">
														{entrypoint.name}
													</Badge>
												</Tooltip.Trigger>
												<Tooltip.Content>
													{entrypoint.address}
												</Tooltip.Content>
											</Tooltip.Root>
										</Tooltip.Provider>
									{/if}
								{/each}
							</div>
						</div>

						<!-- Features -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1">Features</span>
							<div class="col-span-3 space-x-2">
								{#if $overview?.features?.tracing}
									<Badge variant="secondary">Tracing</Badge>
								{/if}
								{#if $overview?.features?.metrics}
									<Badge variant="secondary">Metrics</Badge>
								{/if}
								{#if $overview?.features?.accessLog}
									<Badge variant="secondary">Access Log</Badge>
								{/if}
							</div>
						</div>

						<!-- Providers -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1">Providers</span>
							<div class="col-span-3 space-x-2">
								{#each $overview?.providers ?? [] as provider (provider)}
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
									Routers: {$overview?.http?.routers.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Services: {$overview?.http?.services.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Middlewares: {$overview?.http?.middlewares.total ?? 0}
								</Badge>
							</div>
						</div>

						<!-- TCP Overview -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1 font-mono">TCP</span>
							<div class="col-span-3 space-x-2">
								<Badge variant="secondary">
									Routers: {$overview?.tcp?.routers.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Services: {$overview?.tcp?.services.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Middlewares: {$overview?.tcp?.middlewares.total ?? 0}
								</Badge>
							</div>
						</div>

						<!-- UDP Overview -->
						<div class="grid grid-cols-4 items-center gap-2 text-sm">
							<span class="col-span-1 font-mono">UDP</span>
							<div class="col-span-3 space-x-2">
								<Badge variant="secondary">
									Routers: {$overview?.udp?.routers.total ?? 0}
								</Badge>
								<Badge variant="secondary">
									Services: {$overview?.udp?.services.total ?? 0}
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
							<div class="flex items-center gap-2">
								<Button variant="outline" size="sm" onclick={() => (isYaml = !isYaml)}>
									{isYaml ? 'Show JSON' : 'Show YAML'}
								</Button>
								<CopyButton text={isYaml ? $dynamicYAML : $dynamicJSON} />
							</div>
						</Card.Title>
						<Card.Description>
							This is the current dynamic configuration your Traefik instance is using.
						</Card.Description>
					</Card.Header>
					<Card.Content class="text-sm">
						{#if $dynamicJSON && $dynamicJSON !== 'null'}
							<Highlight
								language={isYaml ? yaml : json}
								code={isYaml ? $dynamicYAML : $dynamicJSON}
								let:highlighted
							>
								<LineNumbers {highlighted} />
							</Highlight>
						{:else}
							<p class="flex items-center justify-center">
								No dynamic configuration, add some routers.
							</p>
						{/if}
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
		</Tabs.Root>
		<Dialog.Close class="w-full">
			<Button class="w-full">Close</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
