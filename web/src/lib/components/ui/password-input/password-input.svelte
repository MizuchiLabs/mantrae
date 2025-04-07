<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Eye, EyeOff } from 'lucide-svelte';

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
		class={cn(
			'border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border px-3 py-2 pr-10 text-base file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
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
