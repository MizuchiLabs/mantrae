<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	middleware.chain = middleware.chain ?? { middlewares: [] };
</script>

<ArrayInput bind:items={middleware.chain.middlewares} label="Middlewares" />
