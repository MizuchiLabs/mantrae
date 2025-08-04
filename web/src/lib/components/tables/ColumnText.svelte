<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type IconProps } from '@lucide/svelte';
	import { type Component } from 'svelte';

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	type Props = {
		label: string;
		icon?: IconComponent;
		iconProps?: IconProps;
		truncate?: boolean;
		maxLength?: number;
		showTooltip?: boolean;
	};

	let {
		label,
		icon,
		iconProps,
		truncate = true,
		maxLength = 30,
		showTooltip = true,
		...restProps
	}: Props = $props();

	const shouldTruncate = truncate && label.length > maxLength;
	const displayLabel = shouldTruncate ? `${label.slice(0, maxLength)}...` : label;
</script>

<div class="flex max-w-full min-w-0 items-center gap-2">
	{#if shouldTruncate && showTooltip}
		<Tooltip.Provider>
			<Tooltip.Root delayDuration={300}>
				<Tooltip.Trigger>
					<Label for={label} class="cursor-help truncate text-sm font-medium" {...restProps}>
						{displayLabel}
					</Label>
				</Tooltip.Trigger>
				<Tooltip.Content side="top" class="max-w-xs break-words">
					{label}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{:else}
		<Label
			for={label}
			class="truncate text-sm font-medium {truncate ? 'max-w-full' : ''}"
			{...restProps}
		>
			{displayLabel}
		</Label>
	{/if}

	{#if icon}
		{@const Icon = icon}
		<Icon class="shrink-0 {iconProps?.class ?? ''}" size={iconProps?.size ?? 16} {...iconProps} />
	{/if}
</div>
