<script lang="ts">
	import { run } from 'svelte/legacy';

	import { Button } from '$lib/components/ui/button/index.js';
	import { Input, type FormInputEvent } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import HoverInfo from '$lib/components/utils/hoverInfo.svelte';
	import { cn } from '$lib/utils';
	import autoAnimate from '@formkit/auto-animate';
	import { Minus, Plus } from 'lucide-svelte';
	import { createEventDispatcher, onMount } from 'svelte';

	interface Props {
		class?: string | undefined | null;
		label: string;
		keyPlaceholder: string;
		valuePlaceholder: string;
		items: Record<string, string> | undefined;
		disabled?: boolean;
		helpText?: string | undefined;
	}

	let {
		class: className = undefined,
		label,
		keyPlaceholder,
		valuePlaceholder,
		items = $bindable(),
		disabled = false,
		helpText = undefined
	}: Props = $props();
	let internalItems: { id: string; key: string; value: string }[] = $state([]);
	const dispatch = createEventDispatcher();

	const convertToInternal = () => {
		if (items && typeof items === 'object' && Object.keys(items).length > 0) {
			internalItems = Object.entries(items).map(([key, value], index) => ({
				id: internalItems[index]?.id || generateId(),
				key,
				value
			}));
		} else {
			internalItems = [{ id: generateId(), key: '', value: '' }];
		}
	};

	const generateId = () => Math.random().toString(36).slice(2, 11);

	const addItem = () => {
		internalItems = [...internalItems, { id: generateId(), key: '', value: '' }]; // Add an item with an empty key-value pair
		dispatchConvert();
	};

	const removeItem = (id: string) => {
		internalItems = internalItems.filter((item) => item.id !== id);
		dispatchConvert();
	};

	const updateKey = (id: string, e: FormInputEvent) => {
		if (!e.target) return;
		const newKey = (e.target as HTMLInputElement).value;
		internalItems = internalItems.map((item) => (item.id === id ? { ...item, key: newKey } : item));
		dispatchConvert();
	};

	const updateValue = (id: string, e: FormInputEvent) => {
		if (!e.target) return;
		const newValue = (e.target as HTMLInputElement).value;
		internalItems = internalItems.map((item) =>
			item.id === id ? { ...item, value: newValue } : item
		);
		dispatchConvert();
	};

	const dispatchConvert = () => {
		items = internalItems.reduce(
			(obj: Record<string, string>, item) => ((obj[item.key] = item.value), obj),
			{}
		);

		dispatch('update', items);
	};

	onMount(() => {
		convertToInternal();
	});
	// Watch `items` prop and convert to array when it changes
	run(() => {
		items, convertToInternal();
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
		{#each internalItems as { id, key, value }, index (id)}
			<li class="flex flex-row items-center justify-end gap-2">
				{#if !disabled}
					<div class="absolute mr-2 flex flex-row items-center justify-between gap-1">
						{#if index === 0}
							<Button
								onclick={addItem}
								class="h-8 w-8 rounded-full bg-red-400 text-black"
								size="icon"
								tabindex={-1}
							>
								<Plus size="1rem" />
							</Button>
						{/if}
						{#if internalItems.length > 1 && index >= 1}
							<Button onclick={() => removeItem(id)} class="h-8 w-8 rounded-full" tabindex={-1}>
								<Minus size="1rem" />
							</Button>
						{/if}
					</div>
				{/if}
				<Input
					id="key"
					type="text"
					bind:value={key}
					placeholder={disabled ? '' : keyPlaceholder}
					oninput={(e) => updateKey(id, e)}
					{disabled}
				/>
				<Input
					id="value"
					type="text"
					bind:value
					placeholder={disabled ? '' : valuePlaceholder}
					on:input={(e) => updateValue(id, e)}
					{disabled}
				/>
			</li>
		{/each}
	</ul>
</div>
