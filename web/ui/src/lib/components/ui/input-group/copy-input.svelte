<script lang="ts">
	import * as InputGroup from '$lib/components/ui/input-group';
	import { CheckIcon, CopyIcon } from '@lucide/svelte';
	import { UseClipboard } from '$lib/hooks/use-clipboard.svelte.js';
	import type { ComponentProps } from 'svelte';
	import type { Input } from '../input';

	let { value = $bindable(''), ...restProps }: ComponentProps<typeof Input> = $props();
	const clipboard = new UseClipboard();
</script>

<InputGroup.Root>
	<InputGroup.Input bind:value class="truncate" {...restProps} />
	<InputGroup.Addon align="inline-end">
		<InputGroup.Button
			aria-label="Copy"
			title="Copy"
			size="icon-xs"
			onclick={() => clipboard.copy(value)}
		>
			{#if clipboard.copied}
				<CheckIcon />
			{:else}
				<CopyIcon />
			{/if}
		</InputGroup.Button>
	</InputGroup.Addon>
</InputGroup.Root>
