<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { createEventDispatcher, onMount } from 'svelte';
	import autoAnimate from '@formkit/auto-animate';

	export let label: string;
	export let placeholder: string;
	export let items: string[] | undefined;
	export let disabled = false;
	const dispatch = createEventDispatcher();

	const addItem = () => {
		items = [...(items ?? []), ''];
		dispatch('update', items);
	};

	const removeItem = (index: number) => {
		items = items?.filter((_, i) => i !== index);
		dispatch('update', items);
	};

	onMount(() => {
		if (!items) items = [''];
		items = items.length > 0 ? items : [''];
	});
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="item" class="text-right">{label}</Label>
	<ul class="col-span-3 space-y-2" use:autoAnimate={{ duration: 100 }}>
		{#each items ?? [] as _, index}
			<li class="flex flex-row items-center justify-end gap-1">
				{#if !disabled}
					<div class="absolute mr-2 flex flex-row items-center justify-between gap-1">
						<Button class="h-8 w-4 rounded-full bg-red-400 text-black" on:click={() => addItem()}>
							<iconify-icon icon="fa6-solid:plus" />
						</Button>
						{#if (items?.length ?? 0) > 1 && index >= 1}
							<Button on:click={() => removeItem(index)} class="h-8 w-4 rounded-full ">
								<iconify-icon icon="fa6-solid:minus" />
							</Button>
						{/if}
					</div>
				{/if}
				{#if items}
					<Input
						id="item"
						type="text"
						bind:value={items[index]}
						placeholder={`${placeholder}`}
						class="focus-visible:ring-0 focus-visible:ring-offset-0"
						{disabled}
					/>
				{/if}
			</li>
		{/each}
	</ul>
</div>
