<script lang="ts" generics="T">
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Plus, XIcon } from '@lucide/svelte';

	interface Props {
		values?: T[];
		placeholder?: string;
		maxItems?: number;
		/** Extract display string from item (default: String(item)) */
		getLabel?: (item: T) => string;
		/** Create new item from input string (default: input as T) */
		createItem?: (input: string) => T;
		/** Extract value for duplicate checking (default: getLabel) */
		getValue?: (item: T) => string;
		/** Called when values change */
		onchange?: (values: T[]) => void;
	}

	let {
		values = $bindable([]),
		placeholder = 'Add item...',
		maxItems,
		getLabel = (item) => String(item),
		createItem = (input) => input as T,
		getValue = getLabel,
		onchange
	}: Props = $props();

	let inputValue = $state('');

	function addItem() {
		const trimmed = inputValue.trim();
		if (!trimmed) return;
		if (values.some((v) => getValue(v) === trimmed)) return;
		if (maxItems && values.length >= maxItems) return;

		values = [...values, createItem(trimmed)];
		inputValue = '';
		onchange?.(values);
	}

	function removeItem(index: number) {
		values = values.filter((_, i) => i !== index);
		onchange?.(values);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			e.preventDefault();
			addItem();
		}
		if (e.key === 'Backspace' && inputValue === '' && values.length > 0) {
			removeItem(values.length - 1);
		}
	}
</script>

<div class="flex w-full flex-col gap-2">
	<InputGroup.Root>
		<InputGroup.Input
			bind:value={inputValue}
			{placeholder}
			onkeydown={handleKeydown}
			disabled={maxItems ? values.length >= maxItems : false}
		/>
		<InputGroup.Addon align="inline-end">
			{#if maxItems}
				<InputGroup.Text class="text-xs text-muted-foreground">
					{values.length}/{maxItems}
				</InputGroup.Text>
			{/if}
			<InputGroup.Button
				variant="default"
				class="rounded-full"
				size="icon-xs"
				disabled={!inputValue.trim() || (maxItems ? values.length >= maxItems : false)}
				onclick={addItem}
			>
				<Plus />
				<span class="sr-only">Add</span>
			</InputGroup.Button>
		</InputGroup.Addon>
	</InputGroup.Root>

	{#if values.length > 0}
		<div class="flex flex-wrap gap-1.5">
			{#each values as item, index (getValue(item))}
				<Badge variant="secondary" class="bg-card" onclick={() => (inputValue = getLabel(item))}>
					{getLabel(item)}
					<button
						type="button"
						class="rounded-full p-0.5 transition-colors hover:bg-muted-foreground/20"
						onclick={() => removeItem(index)}
					>
						<XIcon class="size-3" />
						<span class="sr-only">Remove {getLabel(item)}</span>
					</button>
				</Badge>
			{/each}
		</div>
	{/if}
</div>
