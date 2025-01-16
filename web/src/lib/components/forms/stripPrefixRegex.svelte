<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	interface Props {
		middleware: Middleware;
		disabled?: boolean;
	}

	let { middleware = $bindable(), disabled = false }: Props = $props();
</script>

<ArrayInput
	bind:items={middleware.content.regex}
	label="Regex"
	placeholder="/foo/[a-z0-9]+/[0-9]+/"
	{disabled}
/>
