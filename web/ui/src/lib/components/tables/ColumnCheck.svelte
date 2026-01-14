<script lang="ts">
	import { CircleCheck, CircleSlash } from '@lucide/svelte';

	let { checked, ...restProps }: { checked?: boolean } = $props();
</script>

{#if checked}
	<CircleCheck class="text-green-400" size="1.25rem" {...restProps} />
{:else}
	<CircleSlash class="text-muted-foreground" size="1.25rem" {...restProps} />
{/if}
