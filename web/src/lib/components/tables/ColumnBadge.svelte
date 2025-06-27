<script lang="ts" generics="TData">
	import type { ComponentProps } from 'svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { AlertTriangle, Eye, Globe, Key, type IconProps } from '@lucide/svelte';
	import { type Component } from 'svelte';
	import type { Column } from '@tanstack/table-core';
	import type { RouterTCPTLSConfig } from '$lib/gen/zen/traefik-schemas';

	type IconComponent = Component<IconProps, Record<string, never>, ''>;
	type Props = ComponentProps<typeof Badge> & {
		label: string | string[];
		icon?: IconComponent;
		column?: Column<TData, unknown>;
		limit?: number;
		tls?: RouterTCPTLSConfig;
	};

	let { variant = 'default', label, column, icon, limit = 3, tls, ...restProps }: Props = $props();

	const isArray = Array.isArray(label);
	const items = isArray ? label : [label];
	const visibleItems = items.slice(0, limit);
	const remaining = items.length - limit;
</script>

{#if tls}
	<HoverCard.Root>
		<HoverCard.Trigger>
			<div class="flex flex-col items-start gap-1">
				{#each visibleItems as item (item)}
					<Badge {variant} onclick={() => column?.setFilterValue(item)} {...restProps}>
						{item}
					</Badge>
				{/each}
				{#if remaining > 0}
					<Badge variant="outline" {...restProps}>+{remaining} more</Badge>
				{/if}
			</div>
		</HoverCard.Trigger>
		<HoverCard.Content class="max-w-[300px] space-y-2 text-sm">
			{#if tls.certResolver}
				<div class="text-muted-foreground flex items-center gap-2 italic">
					<Key size={14} class="text-blue-500" />
					<span><strong>Resolver:</strong> {tls.certResolver}</span>
				</div>
			{/if}
			{#if tls.passthrough}
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
	<div class="flex flex-col items-start gap-1">
		{#each visibleItems as item (item)}
			<Badge {variant} onclick={() => column?.setFilterValue(item)} {...restProps}>
				{#if icon}
					{@const Icon = icon}
					<Icon size={16} />
				{/if}
				{item}
			</Badge>
		{/each}
		{#if remaining > 0}
			<Badge variant="outline" {...restProps}>+{remaining} more</Badge>
		{/if}
	</div>
{/if}
