<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Eye, EyeOff } from '@lucide/svelte';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type'>> & {
		showPassword?: boolean;
	};
	let {
		showPassword = $bindable(false),
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		...restProps
	}: Props = $props();
</script>

<div class="relative w-full">
	<input
		bind:this={ref}
		bind:value
		data-slot="input"
		class={cn(
			'border-input bg-background selection:bg-primary dark:bg-input/30 selection:text-primary-foreground ring-offset-background placeholder:text-muted-foreground flex h-9 w-full min-w-0 rounded-md border px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
			'focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]',
			'aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive',
			className
		)}
		type={showPassword ? 'text' : 'password'}
		placeholder="••••••••"
		{...restProps}
	/>
	<Button
		variant="ghost"
		size="icon"
		class="text-muted-foreground absolute inset-y-0 right-1 flex items-center justify-center p-2 hover:bg-transparent hover:text-red-400"
		onclick={() => (showPassword = !showPassword)}
	>
		{#if showPassword}
			<Eye size="1rem" />
		{:else}
			<EyeOff size="1rem" />
		{/if}
	</Button>
</div>
