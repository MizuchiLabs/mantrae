<script lang="ts" generics="TData">
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { AlertTriangle, Eye, Globe, Key } from '@lucide/svelte';
	import type { Column } from '@tanstack/table-core';
	import type { RouterTCPTLSConfig, RouterTLSConfig } from '$lib/gen/zen/traefik-schemas';

	type Props = {
		tls?: RouterTLSConfig | RouterTCPTLSConfig;
		column?: Column<TData, unknown>;
	};

	let { tls, column }: Props = $props();
</script>

{#if tls}
	<HoverCard.Root openDelay={100}>
		<HoverCard.Trigger>
			<Badge variant="secondary" onclick={() => column?.setFilterValue?.(tls?.certResolver ?? '')}>
				{tls.certResolver ? tls.certResolver : 'Enabled'}
			</Badge>
		</HoverCard.Trigger>
		<HoverCard.Content class="max-w-[300px] space-y-2 text-sm">
			{#if tls.certResolver}
				<div class="text-muted-foreground flex items-center gap-2 italic">
					<Key size={14} class="text-blue-500" />
					<span><strong>Resolver:</strong> {tls.certResolver}</span>
				</div>
			{/if}
			{#if tls?.passthrough}
				<div class="text-muted-foreground flex items-center gap-2 italic">
					<Globe size={14} class="text-green-500" />
					<span><strong>Passthrough:</strong> Enabled</span>
				</div>
			{/if}
			{#if tls.options}
				<div class="text-muted-foreground flex items-center gap-2 italic">
					<Eye size={14} class="text-purple-500" />
					<span><strong>Options:</strong> {tls.options}</span>
				</div>
			{/if}
			{#if tls.domains?.length}
				<div class="text-muted-foreground flex items-center gap-2 italic">
					<Globe size={14} class="mt-1 text-orange-400" />
					<div class="space-y-1">
						<strong>Domains:</strong>
						<ul class="ml-4 list-disc">
							{#each tls.domains as d (d)}
								<li>{d.main}{d.sans ? ` (${d.sans.join(', ')})` : ''}</li>
							{/each}
						</ul>
					</div>
				</div>
			{/if}
			{#if !tls.certResolver && !tls.passthrough && !tls.options && !tls.domains?.length}
				<div class="text-muted-foreground flex items-center gap-2 italic">
					<AlertTriangle size={14} class="text-yellow-400" />
					<span>TLS is enabled but has no configuration</span>
				</div>
			{/if}
		</HoverCard.Content>
	</HoverCard.Root>
{:else}
	<Badge variant="outline">Disabled</Badge>
{/if}
