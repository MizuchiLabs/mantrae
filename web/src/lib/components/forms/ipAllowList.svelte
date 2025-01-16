<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import { z } from 'zod';
	import { CustomIPSchema, CustomIPSchemaOptional } from '../utils/validation';
	import { onDestroy } from 'svelte';

	interface Props {
		middleware: Middleware;
		disabled?: boolean;
	}

	let { middleware = $bindable(), disabled = false }: Props = $props();

	const templateRange = ['192.168.0.0/16', '172.16.0.0/12', '127.0.0.1/32', '10.0.0.0/8'];

	const ipAllowListSchema = z.object({
		sourceRange: z.array(CustomIPSchema).default([]),
		ipStrategy: z
			.object({
				depth: z
					.union([z.string(), z.number()])
					.transform((value) => (value === '' ? null : Number(value)))
					.nullish(),
				excludedIPs: z.array(CustomIPSchemaOptional).nullish()
			})
			.default({})
	});
	middleware.content = ipAllowListSchema.parse({ ...middleware.content });

	let errors: Record<any, string[] | undefined> = $state({});
	const validate = () => {
		try {
			middleware.content = ipAllowListSchema.parse(middleware.content);
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};

	let isTemplate = $state(middleware.content?.sourceRange?.length > 0);
	const toggleIpAllowList = () => {
		isTemplate = !isTemplate;
		if (isTemplate) {
			middleware.content.sourceRange = templateRange;
		} else {
			middleware.content.sourceRange = [];
		}
	};

	onDestroy(() => {
		validate();
	});
</script>

<div class="flex items-center justify-end gap-2">
	<Button on:click={toggleIpAllowList}>
		{isTemplate ? 'Clear IPs' : 'Add Private IP range'}
	</Button>
</div>

<ArrayInput
	bind:items={middleware.content.sourceRange}
	label="Source Range"
	placeholder="192.168.1.1/32"
	on:update={validate}
	{disabled}
/>
{#if errors.sourceRange}
	<div class="col-span-4 text-right text-sm text-red-500">{errors.sourceRange}</div>
{/if}

{#if middleware.content.ipStrategy && middleware.protocol === 'http'}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="depth" class="text-right">Depth</Label>
		<div class="relative col-span-3">
			<Input
				id="depth"
				name="depth"
				type="text"
				bind:value={middleware.content.ipStrategy.depth}
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
		bind:items={middleware.content.ipStrategy.excludedIPs}
		label="Excluded IPs"
		placeholder="192.168.1.1"
		on:update={validate}
		{disabled}
	/>
	{#if errors.ipStrategy}
		<div class="col-span-4 text-right text-sm text-red-500">{errors.ipStrategy}</div>
	{/if}
{/if}
