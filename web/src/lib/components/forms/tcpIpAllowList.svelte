<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.tcpIpAllowList = { sourceRange: [], ...middleware.tcpIpAllowList };
</script>

{#if middleware.tcpIpAllowList}
	<ArrayInput
		bind:items={middleware.tcpIpAllowList.sourceRange}
		label="Source Range"
		placeholder="192.168.1.1/32"
		{disabled}
	/>
{/if}
