<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	middleware.tcpIpAllowList = middleware.tcpIpAllowList ?? { sourceRange: [] };
</script>

<ArrayInput bind:items={middleware.tcpIpAllowList.sourceRange} label="Source Range" />
