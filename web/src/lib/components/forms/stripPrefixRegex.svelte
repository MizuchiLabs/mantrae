<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: HttpMiddleware;
	middleware.stripPrefixRegex = middleware.stripPrefixRegex ?? {
		regex: ['']
	};
</script>

<ArrayInput bind:items={middleware.stripPrefixRegex.regex} label="Regex" />
