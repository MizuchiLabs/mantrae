<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { createEventDispatcher } from 'svelte';

	export let label: string;
	export let items: string[] = [''];
	const dispatch = createEventDispatcher();

	const addItem = () => {
		items = [...items, ''];
		dispatch('update', items);
	};

	const removeItem = (index: number) => {
		items = items.filter((_, i) => i !== index);
		dispatch('update', items);
	};
</script>

<div class="space-y-2 rounded-lg py-2">
	<div class="flex flex-row items-center justify-between">
		<Label for="item" class="text-right">{label}</Label>
		<Button class="h-8 w-4 bg-red-400 text-black" on:click={() => addItem()}>
			<iconify-icon icon="fa6-solid:plus" />
		</Button>
	</div>
	{#each items as item, index (index)}
		<div class="flex items-center justify-end gap-1">
			<Input
				id="item"
				type="text"
				bind:value={items[index]}
				placeholder={`${label} ${index + 1}`}
				class="focus-visible:ring-0 focus-visible:ring-offset-0"
			/>
			{#if items.length > 1 && index >= 1}
				<Button on:click={() => removeItem(index)} class="h-8 w-4">
					<iconify-icon icon="fa6-solid:minus" />
				</Button>
			{/if}
		</div>
	{/each}
</div>
