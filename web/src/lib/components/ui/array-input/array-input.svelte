<script lang="ts">
	import { run } from 'svelte/legacy';

	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import HoverInfo from '$lib/components/utils/hoverInfo.svelte';
	import { cn } from '$lib/utils';
	import autoAnimate from '@formkit/auto-animate';
	import { Minus, Plus } from 'lucide-svelte';
	import { createEventDispatcher, onMount } from 'svelte';

	interface Props {
		class?: string | undefined | null;
		label: string;
		placeholder: string;
		items: string[] | undefined;
		disabled?: boolean;
		helpText?: string | undefined;
	}

	let {
		class: className = undefined,
		label,
		placeholder,
		items = $bindable(),
		disabled = false,
		helpText = undefined
	}: Props = $props();
	const dispatch = createEventDispatcher();

	const verifyArray = () => {
		if (!items || items.length === 0) {
			items = [''];
		}
	};

	const addItem = () => {
		items = [...(items ?? []), ''];
		dispatch('update', items);
	};

	const removeItem = (index: number) => {
		items = items?.filter((_, i) => i !== index);
		dispatch('update', items);
	};

	const update = (index: number, target: EventTarget | null) => {
		const value = (target as HTMLInputElement).value;
		items = items?.map((_, i) => (i === index ? value : _));
		dispatch('update', items);
	};

	onMount(() => {
		verifyArray();
	});
	run(() => {
		items, verifyArray();
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
		{#each items || [] as item, index}
			<li class="flex flex-row items-center justify-end gap-1">
				{#if !disabled}
					<div class="absolute mr-2 flex flex-row items-center justify-between gap-1">
						{#if index === 0}
							<Button
								onclick={() => addItem()}
								class="h-8 w-8 rounded-full bg-red-400 text-black"
								size="icon"
							>
								<Plus size="1rem" />
							</Button>
						{/if}
						{#if (items?.length ?? 0) > 1 && index >= 1}
							<Button onclick={() => removeItem(index)} class="h-8 w-8 rounded-full" size="icon">
								<Minus size="1rem" />
							</Button>
						{/if}
					</div>
				{/if}
				{#if items}
					<Input
						id="item"
						type="text"
						bind:value={items[index]}
						placeholder={disabled ? '' : placeholder}
						oninput={(e) => update(index, e.target)}
						{disabled}
					/>
				{/if}
			</li>
		{/each}
	</ul>
</div>
