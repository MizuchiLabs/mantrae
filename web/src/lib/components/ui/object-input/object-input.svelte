<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input, type FormInputEvent } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { createEventDispatcher, onMount } from 'svelte';
	import autoAnimate from '@formkit/auto-animate';
	import { cn } from '$lib/utils';
	import HoverInfo from '$lib/components/utils/hoverInfo.svelte';

	let className: string | undefined | null = undefined;
	export { className as class };
	export let label: string;
	export let keyPlaceholder: string;
	export let valuePlaceholder: string;
	export let items: Record<string, string> = {};
	export let disabled = false;
	export let helpText: string | undefined = undefined;
	const dispatch = createEventDispatcher();

	const addItem = () => {
		items = { ...items, '': '' }; // Add an empty key-value pair
		dispatch('update', items);
	};

	const removeItem = (key: string) => {
		const { [key]: _, ...rest } = items;
		items = rest;
		dispatch('update', items);
	};

	const updateKey = (oldKey: string, e: FormInputEvent) => {
		if (!e.target) return;
		const newKey = (e.target as HTMLInputElement).value;
		if (newKey !== oldKey) {
			const { [oldKey]: value, ...rest } = items;
			items = { ...rest, [newKey]: value };
			dispatch('update', items);
		}
	};

	const updateValue = (key: string, e: FormInputEvent) => {
		if (!e.target) return;
		const value = (e.target as HTMLInputElement).value;
		items = { ...items, [key]: value };
		dispatch('update', items);
	};

	onMount(() => {
		if (!items || typeof items !== 'object' || Object.keys(items).length === 0) {
			items = { '': '' };
		}
	});
</script>

<div class={cn('grid grid-cols-4 items-center gap-4', className)}>
	<Label for="item" class="col-span-1 flex items-center justify-end gap-0.5">
		{label}
		{#if helpText}
			<HoverInfo text={helpText} />
		{/if}
	</Label>
	<ul class="col-span-3 space-y-2" use:autoAnimate={{ duration: 100 }}>
		{#each Object.entries(items || {}) as [key, value], index}
			<li class="flex flex-row items-center justify-end gap-2">
				{#if !disabled}
					<div class="absolute mr-2 flex flex-row items-center justify-between gap-1">
						{#if index === 0}
							<Button class="h-8 w-4 rounded-full bg-red-400 text-black" on:click={addItem}>
								<iconify-icon icon="fa6-solid:plus" />
							</Button>
						{/if}
						{#if Object.keys(items).length > 1 && index >= 1}
							<Button on:click={() => removeItem(key)} class="h-8 w-4 rounded-full">
								<iconify-icon icon="fa6-solid:minus" />
							</Button>
						{/if}
					</div>
				{/if}
				<Input
					id="key"
					type="text"
					bind:value={key}
					placeholder={disabled ? '' : keyPlaceholder}
					on:input={(e) => updateKey(key, e)}
					class="focus-visible:ring-0 focus-visible:ring-offset-0"
					{disabled}
				/>
				<Input
					id="value"
					type="text"
					bind:value
					placeholder={disabled ? '' : valuePlaceholder}
					on:input={(e) => updateValue(key, e)}
					class="focus-visible:ring-0 focus-visible:ring-offset-0"
					{disabled}
				/>
			</li>
		{/each}
	</ul>
</div>
