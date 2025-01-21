<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import { Badge } from '$lib/components/ui/badge';

	type Props = ComponentProps<typeof Badge> & {
		label: string | string[];
		limit?: number;
	};

	let { variant = 'default', label, limit = 3, ...restProps }: Props = $props();

	const isArray = Array.isArray(label);
	const items = isArray ? label : [label];
	const visibleItems = items.slice(0, limit);
	const remaining = items.length - limit;
</script>

<div class="flex flex-col items-start gap-1">
	{#each visibleItems as item}
		<Badge {variant} {...restProps}>{item}</Badge>
	{/each}
	{#if remaining > 0}
		<Badge variant="outline" {...restProps}>+{remaining} more</Badge>
	{/if}
</div>
