<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import type { Middleware } from '$lib/types/middlewares';
	import { CustomIPSchema } from '../utils/validation';
	import { z } from 'zod';

	export let middleware: Middleware;
	export let disabled = false;
	const defaultTemplate = {
		sourceRange: ['192.168.0.0/16', '172.16.0.0/12', '127.0.0.1/32', '10.0.0.0/8']
	};

	const tcpIpAllowListSchema = z.object({
		sourceRange: z.array(CustomIPSchema).default(['127.0.0.1/32'])
	});
	middleware.tcpIpAllowList = tcpIpAllowListSchema.parse({ ...middleware.tcpIpAllowList });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.tcpIpAllowList = tcpIpAllowListSchema.parse(middleware.tcpIpAllowList); // Parse the tcpIpAllowList object
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};

	let isTemplate = false;
	const toggleIpAllowList = () => {
		isTemplate = !isTemplate;
		if (isTemplate) {
			middleware.tcpIpAllowList = defaultTemplate;
		} else {
			middleware.tcpIpAllowList = { sourceRange: [] };
		}
		validate();
	};
</script>

<div class="flex items-center justify-end gap-2">
	<Button on:click={toggleIpAllowList}>
		{isTemplate ? 'Clear IPs' : 'Add Private IP range'}
	</Button>
</div>

{#if middleware.tcpIpAllowList}
	<ArrayInput
		bind:items={middleware.tcpIpAllowList.sourceRange}
		label="Source Range"
		placeholder="192.168.1.1/32"
		on:update={validate}
		{disabled}
	/>
	{#if errors.sourceRange}
		<div class="col-span-4 text-right text-sm text-red-500">{errors.sourceRange}</div>
	{/if}
{/if}
