<script lang="ts">
	import { Button, type ButtonProps } from '$lib/components/ui/button';
	import { UseClipboard } from '$lib/hooks/use-clipboard.svelte';
	import { cn } from '$lib/utils';
	import { Check, Copy, X, type IconProps } from '@lucide/svelte';
	import type { Component } from 'svelte';
	import { scale } from 'svelte/transition';

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	// omit href so you can't create a link
	interface Props extends Omit<ButtonProps, 'href'> {
		label?: string;
		text: string;
		icon?: IconComponent;
		animationDuration?: number;
		onCopy?: (status: UseClipboard['status']) => void;
	}

	let {
		label,
		text,
		icon,
		animationDuration = 500,
		variant = 'ghost',
		size = 'icon',
		onCopy,
		class: className,
		...restProps
	}: Props = $props();

	const clipboard = new UseClipboard();
</script>

<Button
	{...restProps}
	{variant}
	{size}
	class={cn(className) + ' right-0.5 '}
	type="button"
	name="copy"
	tabindex={-1}
	onclick={async () => {
		const status = await clipboard.copy(text);

		onCopy?.(status);
	}}
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
		<div in:scale={{ duration: animationDuration, start: 0.85 }} class="flex items-center gap-1">
			{label ?? ''}
			{#if icon}
				{@const Icon = icon}
				<Icon />
			{:else}
				<Copy />
			{/if}
			<span class="sr-only">Copy</span>
		</div>
	{/if}
</Button>
