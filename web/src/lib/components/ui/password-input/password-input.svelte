<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Eye, EyeOff } from 'lucide-svelte';

	let showPassword = $state(false);
	interface Props {
		password: string | undefined;
		required?: boolean;
		[props: string]: any;
	}

	let { password = $bindable(), required, ...restProps }: Props = $props();
</script>

<div class="flex flex-row items-center justify-end gap-1">
	{#if showPassword}
		<Input
			id="password"
			type="text"
			bind:value={password}
			placeholder="••••••••"
			{required}
			{...restProps}
		/>
	{:else}
		<Input
			id="password"
			type="password"
			bind:value={password}
			placeholder="••••••••"
			{required}
			{...restProps}
		/>
	{/if}
	<Button
		variant="ghost"
		size="icon"
		class="absolute hover:bg-transparent hover:text-red-400"
		onclick={() => (showPassword = !showPassword)}
	>
		{#if showPassword}
			<Eye size="1rem" />
		{:else}
			<EyeOff size="1rem" />
		{/if}
	</Button>
</div>
