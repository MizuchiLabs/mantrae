<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Eye, EyeOff } from '@lucide/svelte';
	import { Input } from '$lib/components/ui/input/index.js';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type' | 'files'>> & {
		showPassword?: boolean;
	};
	let {
		ref = $bindable(null),
		value = $bindable(),
		showPassword = $bindable(false),
		class: className,
		...restProps
	}: Props = $props();
</script>

<div class="relative w-full">
	<Input
		bind:ref
		bind:value
		data-slot="input"
		class={className + ' pr-10'}
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
