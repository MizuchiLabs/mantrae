<script lang="ts" generics="TData">
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { type IconProps } from '@lucide/svelte';
	import { type Component } from 'svelte';
	import type { ComponentProps } from 'svelte';
	import type { Column } from '@tanstack/table-core';

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	type Props = ComponentProps<typeof Badge> & {
		label: string | string[];
		icon?: IconComponent;
		iconProps?: IconProps;
		column?: Column<TData, unknown>;
		limit?: number;
	};

	let {
		label,
		variant = 'default',
		icon,
		iconProps,
		column,
		limit = 3,
		...restProps
	}: Props = $props();

	const items = Array.isArray(label) ? label : [label];
	const visible = items.slice(0, limit);
	const hidden = items.slice(limit);
</script>

<div class="flex flex-col items-start gap-1">
	{#each visible as item (item)}
		<Badge
			{variant}
			onclick={() => column?.setFilterValue?.(item)}
			class="flex items-center gap-1 hover:cursor-pointer"
			{...restProps}
		>
			{#if icon}
				{@const Icon = icon}
				<Icon size={14} {...iconProps} />
			{/if}
			{item}
		</Badge>
	{/each}

	{#if hidden.length}
		<HoverCard.Root openDelay={100}>
			<HoverCard.Trigger>
				<Badge variant="outline" class="cursor-pointer" {...restProps}>
					+{hidden.length} more
				</Badge>
			</HoverCard.Trigger>
			<HoverCard.Content class="w-auto">
				<div class="flex flex-col gap-1">
					{#each hidden as item (item)}
						<Badge
							{variant}
							onclick={() => column?.setFilterValue?.(item)}
							class="flex items-center gap-1 hover:cursor-pointer"
							{...restProps}
						>
							{#if icon}
								{@const Icon = icon}
								<Icon size={14} {...iconProps} />
							{/if}
							{item}
						</Badge>
					{/each}
				</div>
			</HoverCard.Content>
		</HoverCard.Root>
	{/if}
</div>
