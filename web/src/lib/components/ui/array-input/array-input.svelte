<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { createEventDispatcher, onMount } from 'svelte';

	export let label: string;
	export let items: string[];
	const dispatch = createEventDispatcher();

	const addItem = () => {
		items = [...items, ''];
		dispatch('update', items);
	};

	const removeItem = (index: number) => {
		items = items.filter((_, i) => i !== index);
		dispatch('update', items);
	};

	onMount(() => {
		items = items.length > 0 ? items : [''];
	});
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="item" class="text-right">{label}</Label>
	<div class="col-span-3 space-y-2">
		{#each items as _, index}
			<div class="flex flex-row items-center justify-end gap-1">
				<div class="absolute mr-2 flex flex-row items-center justify-between gap-1">
					<Button class="h-8 w-4 rounded-full bg-red-400 text-black" on:click={() => addItem()}>
						<iconify-icon icon="fa6-solid:plus" />
					</Button>
					{#if items.length > 1 && index >= 1}
						<Button on:click={() => removeItem(index)} class="h-8 w-4 rounded-full ">
							<iconify-icon icon="fa6-solid:minus" />
						</Button>
					{/if}
				</div>
				<Input
					id="item"
					type="text"
					bind:value={items[index]}
					placeholder={`${label} ${index + 1}`}
					class="focus-visible:ring-0 focus-visible:ring-offset-0"
				/>
			</div>
		{/each}
	</div>
</div>
