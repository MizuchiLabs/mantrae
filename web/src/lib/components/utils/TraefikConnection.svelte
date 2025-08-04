<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Link, Terminal, FileText } from '@lucide/svelte';
	import { buildConnectionString } from '$lib/api';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import CopyButton from '../ui/copy-button/copy-button.svelte';

	interface Props {
		profile: Profile;
		variant?: 'compact' | 'full';
	}

	let { profile, variant = 'compact' }: Props = $props();

	function getYamlConfig(endpoint: string): string {
		return `providers:
  http:
    endpoint: "${endpoint}"`;
	}

	function getCliConfig(endpoint: string): string {
		return `--providers.http.endpoint=${endpoint}`;
	}
</script>

{#await buildConnectionString(profile) then connectionString}
	{#if variant === 'compact'}
		<Popover.Root>
			<Popover.Trigger>
				<Button size="sm" variant="outline" class="w-full gap-2">
					<Link class="h-3 w-3" />
					Connection
				</Button>
			</Popover.Trigger>
			<Popover.Content class="w-80">
				<div class="space-y-4">
					<div class="space-y-2">
						<h4 class="leading-none font-medium">Traefik HTTP Provider</h4>
						<p class="text-sm text-muted-foreground">
							Copy the configuration for your Traefik setup
						</p>
					</div>

					<div class="space-y-3">
						<!-- YAML Config -->
						<div class="space-y-2">
							<div class="flex items-center justify-between">
								<Badge variant="outline" class="gap-1">
									<FileText class="h-3 w-3" />
									YAML
								</Badge>
								<CopyButton text={getYamlConfig(connectionString)} />
							</div>
							<pre class="block overflow-x-auto rounded bg-muted p-2 text-xs"><code
									>{getYamlConfig(connectionString)}</code
								></pre>
						</div>

						<!-- CLI Config -->
						<div class="space-y-2">
							<div class="flex items-center justify-between">
								<Badge variant="outline" class="gap-1">
									<Terminal class="h-3 w-3" />
									CLI
								</Badge>
								<CopyButton text={getCliConfig(connectionString)} />
							</div>
							<code class="block rounded bg-muted p-2 text-xs break-all">
								{getCliConfig(connectionString)}
							</code>
						</div>
					</div>
				</div>
			</Popover.Content>
		</Popover.Root>
	{:else}
		<!-- Full variant for dedicated sections -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Link class="h-4 w-4" />
					Traefik Connection
				</Card.Title>
				<Card.Description>Copy the configuration for your Traefik setup</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-3">
				<!-- YAML Config -->
				<div class="space-y-1">
					<div class="flex items-center justify-between">
						<Badge variant="outline" class="gap-1">
							<FileText class="h-3 w-3" />
							Static Configuration (YAML)
						</Badge>
						<CopyButton text={getYamlConfig(connectionString)} />
					</div>
					<pre class="block overflow-x-auto rounded bg-muted p-3 text-sm"><code
							>{getYamlConfig(connectionString)}</code
						></pre>
				</div>

				<!-- CLI Config -->
				<div class="space-y-1">
					<div class="flex items-center justify-between">
						<Badge variant="outline" class="gap-1">
							<Terminal class="h-3 w-3" />
							Command Line
						</Badge>
						<CopyButton text={getCliConfig(connectionString)} />
					</div>
					<code class="block rounded bg-muted p-3 text-sm break-all">
						{getCliConfig(connectionString)}
					</code>
				</div>
			</Card.Content>
		</Card.Root>
	{/if}
{/await}
