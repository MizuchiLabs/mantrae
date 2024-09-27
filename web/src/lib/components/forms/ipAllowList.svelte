<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import { z } from 'zod';
	import { CustomIPSchema, CustomIPSchemaOptional } from '../utils/validation';

	export let middleware: Middleware;
	export let disabled = false;

	const emptyIpAllowList = {
		sourceRange: [],
		ipStrategy: { excludedIPs: [] }
	};
	const defaultTemplate = {
		sourceRange: ['192.168.0.0/16', '172.16.0.0/12', '127.0.0.1/32', '10.0.0.0/8'],
		ipStrategy: { excludedIPs: [] }
	};

	const ipAllowListSchema = z.object({
		sourceRange: z.array(CustomIPSchema).default([]),
		ipStrategy: z
			.object({
				depth: z.coerce.number().int().nonnegative().optional(),
				excludedIPs: z.array(CustomIPSchemaOptional).optional()
			})
			.default({})
	});
	middleware.ipAllowList = ipAllowListSchema.parse({ ...middleware.ipAllowList });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.ipAllowList = ipAllowListSchema.parse(middleware.ipAllowList);
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
			middleware.ipAllowList = defaultTemplate;
		} else {
			middleware.ipAllowList = emptyIpAllowList;
		}
		validate();
	};
</script>

<div class="flex items-center justify-end gap-2">
	<Button on:click={toggleIpAllowList}>
		{isTemplate ? 'Clear IPs' : 'Add Private IP range'}
	</Button>
</div>

{#if middleware.ipAllowList}
	<ArrayInput
		bind:items={middleware.ipAllowList.sourceRange}
		label="Source Range"
		placeholder="192.168.1.1/32"
		on:update={validate}
		{disabled}
	/>
	{#if errors.sourceRange}
		<div class="col-span-4 text-right text-sm text-red-500">{errors.sourceRange}</div>
	{/if}

	{#if middleware.ipAllowList.ipStrategy && middleware.middlewareType === 'http'}
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="depth" class="text-right">Depth</Label>
			<div class="relative col-span-3">
				<Input
					id="depth"
					name="depth"
					type="text"
					bind:value={middleware.ipAllowList.ipStrategy.depth}
					on:input={validate}
					placeholder="0"
					{disabled}
				/>
				{#if errors.ipStrategy}
					<div class="col-span-4 text-right text-sm text-red-500">{errors.ipStrategy}</div>
				{/if}
			</div>
		</div>

		<ArrayInput
			bind:items={middleware.ipAllowList.ipStrategy.excludedIPs}
			label="Excluded IPs"
			placeholder="192.168.1.1"
			on:update={validate}
			{disabled}
		/>
		{#if errors.ipStrategy}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.ipStrategy}</div>
		{/if}
	{/if}
{/if}
