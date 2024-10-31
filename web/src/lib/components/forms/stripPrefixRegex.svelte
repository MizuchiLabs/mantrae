<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
</script>

<ArrayInput
	bind:items={middleware.content.regex}
	label="Regex"
	placeholder="/foo/[a-z0-9]+/[0-9]+/"
	{disabled}
/>
