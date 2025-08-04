<script lang="ts">
	import { Check, Copy, X } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { UseClipboard } from '$lib/hooks/use-clipboard.svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import { scale } from 'svelte/transition';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type' | 'files'>>;
	let {
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		...restProps
	}: Props = $props();

	const clipboard = new UseClipboard();
	let animationDuration = 500;
</script>

<div class="relative w-full">
	<Input
		bind:ref
		bind:value
		data-slot="input"
		class={className + ' pr-10'}
		type="text"
		{...restProps}
	/>
	<Button
		variant="ghost"
		size="icon"
		class="absolute inset-y-0 right-1 flex items-center justify-center p-2 text-muted-foreground hover:bg-transparent hover:text-red-400"
		onclick={async () => await clipboard.copy(value)}
	>
		{#if clipboard.status === 'success'}
			<div in:scale={{ duration: animationDuration, start: 0.85 }}>
				<Check />
				<span class="sr-only">Copied</span>
			</div>
		{:else if clipboard.status === 'failure'}
			<div in:scale={{ duration: animationDuration, start: 0.85 }}>
				<X />
				<span class="sr-only">Failed to copy</span>
			</div>
		{:else}
			<div in:scale={{ duration: animationDuration, start: 0.85 }}>
				<Copy />
				<span class="sr-only">Copy</span>
			</div>
		{/if}
	</Button>
</div>
