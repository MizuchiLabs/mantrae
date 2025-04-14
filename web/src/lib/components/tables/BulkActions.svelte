<script lang="ts" generics="T">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button } from '$lib/components/ui/button';
	import type { BulkAction } from './types';

	export let selectedCount: number;
	export let totalCount: number;
	export let actions: BulkAction<T>[] = [];
	export let selectedItems: T[] = [];
</script>

<div class="bg-muted/50 my-2 flex items-center justify-between gap-2 rounded-lg border p-2 pr-6">
	<div class="flex items-center gap-2">
		{#each actions as action (action.label)}
			{#if action.type === 'button'}
				<Button
					variant={action.variant ?? 'secondary'}
					size="sm"
					class={action.class}
					onclick={() => action.onClick?.(selectedItems)}
					disabled={action.disabled}
				>
					{#if action.icon}
						{@const Icon = action.icon}
						<Icon size={16} class="mr-1 h-4 w-4" />
					{/if}
					{action.label}
				</Button>
			{:else if action.type === 'select' && action.options}
				<Select.Root
					type="single"
					onValueChange={(value) => {
						const option = action.options?.find((o) => o.value === value);
						option?.onClick(selectedItems, value);
					}}
				>
					<Select.Trigger>
						{action.label}
					</Select.Trigger>
					<Select.Content>
						{#each action.options as option (option.value)}
							<Select.Item value={option.value}>{option.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			{/if}
		{/each}
	</div>
	<span class="text-muted-foreground text-sm">
		{selectedCount} of {totalCount} item(s) selected.
	</span>
</div>
