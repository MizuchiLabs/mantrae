<script lang="ts" generics="TData">
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
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
		truncateAt?: number;
	};

	let {
		label,
		variant = 'default',
		icon,
		iconProps,
		column,
		limit = 3,
		truncateAt = 15,
		...restProps
	}: Props = $props();

	const items = Array.isArray(label) ? label : [label];
	const visible = items.slice(0, limit);
	const hidden = items.slice(limit);

	function truncateText(text: string, maxLength: number): string {
		return text.length > maxLength ? `${text.slice(0, maxLength)}...` : text;
	}
</script>

<div class="flex max-w-full flex-wrap items-center gap-1">
	{#each visible as item, index (item)}
		{@const truncated = truncateText(item, truncateAt)}
		{@const shouldShowTooltip = item.length > truncateAt}

		{#if shouldShowTooltip}
			<Tooltip.Provider>
				<Tooltip.Root delayDuration={300}>
					<Tooltip.Trigger>
						<Badge
							{variant}
							onclick={() => column?.setFilterValue?.(item)}
							class="text-sm transition-colors duration-200 hover:cursor-pointer 
                            {index > 0 ? 'hidden sm:flex' : ''}"
							{...restProps}
						>
							{#if icon && index === 0}
								{@const Icon = icon}
								<Icon size={12} class="shrink-0" {...iconProps} />
							{/if}
							<span class="max-w-[8rem] truncate sm:max-w-none">{truncated}</span>
						</Badge>
					</Tooltip.Trigger>
					<Tooltip.Content side="top" class="max-w-xs break-words">
						{item}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{:else}
			<Badge
				{variant}
				onclick={() => column?.setFilterValue?.(item)}
				class="text-sm transition-colors duration-200 hover:cursor-pointer
					   {index > 0 ? 'hidden sm:flex' : ''}"
				{...restProps}
			>
				{#if icon && index === 0}
					{@const Icon = icon}
					<Icon size={12} class="shrink-0" {...iconProps} />
				{/if}
				<span class="max-w-[8rem] truncate sm:max-w-none">{item}</span>
			</Badge>
		{/if}
	{/each}

	{#if hidden.length > 0}
		<HoverCard.Root openDelay={200}>
			<HoverCard.Trigger>
				<Badge
					variant="outline"
					class="cursor-pointer text-xs transition-colors duration-200"
					{...restProps}
				>
					{#if visible.length > 1}
						<span class="sm:hidden">+{visible.length - 1 + hidden.length}</span>
						<span class="hidden sm:block">+{hidden.length}</span>
					{:else}
						+{hidden.length}
					{/if}
				</Badge>
			</HoverCard.Trigger>
			<HoverCard.Content class="w-auto max-w-sm">
				<div class="flex flex-wrap gap-1">
					{#each visible.slice(1) as item (item)}
						<Badge
							{variant}
							onclick={() => column?.setFilterValue?.(item)}
							class="text-xs hover:cursor-pointer"
							{...restProps}
						>
							{#if icon}
								{@const Icon = icon}
								<Icon size={12} {...iconProps} />
							{/if}
							{item}
						</Badge>
					{/each}
					{#each hidden as item (item)}
						<Badge
							{variant}
							onclick={() => column?.setFilterValue?.(item)}
							class="text-xs hover:cursor-pointer"
							{...restProps}
						>
							{#if icon}
								{@const Icon = icon}
								<Icon size={12} {...iconProps} />
							{/if}
							{item}
						</Badge>
					{/each}
				</div>
			</HoverCard.Content>
		</HoverCard.Root>
	{/if}
</div>
