<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.stripPrefixRegex = { regex: [], ...middleware.stripPrefixRegex };
</script>

{#if middleware.stripPrefixRegex}
	<ArrayInput
		bind:items={middleware.stripPrefixRegex.regex}
		label="Regex"
		placeholder="/foo/[a-z0-9]+/[0-9]+/"
			{disabled}
	/>
{/if}
