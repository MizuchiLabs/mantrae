<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import { z } from 'zod';
	import { CustomIPSchemaOptional } from '../utils/validation';

	export let middleware: Middleware;
	export let disabled = false;

	const rateLimitSchema = z.object({
		period: z.string({ required_error: 'Period is required' }).trim().default('1s'),
		average: z.coerce
			.number({ required_error: 'Average is required' })
			.int()
			.nonnegative()
			.default(0),
		burst: z.coerce.number({ required_error: 'Burst is required' }).int().nonnegative().default(1),
		sourceCriterion: z
			.object({
				ipStrategy: z
					.object({
						depth: z.coerce.number().int().nonnegative().optional(),
						excludedIPs: z.array(CustomIPSchemaOptional).optional()
					})
					.default({ excludedIPs: [] }),
				requestHeaderName: z.string().trim().optional(),
				requestHost: z.boolean().optional()
			})
			.default({})
	});
	middleware.rateLimit = rateLimitSchema.parse({ ...middleware.rateLimit });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.rateLimit = rateLimitSchema.parse(middleware.rateLimit); // Parse the rateLimit object
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};
</script>

{#if middleware.rateLimit}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="average" class="text-right">Average</Label>
		<div class="relative col-span-3">
			<Input
				id="average"
				name="average"
				type="number"
				bind:value={middleware.rateLimit.average}
				on:input={validate}
				placeholder="0"
				{disabled}
			/>
			{#if errors.average}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.average}</div>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="period" class="text-right">Period</Label>
		<div class="relative col-span-3">
			<Input
				id="period"
				name="period"
				type="text"
				bind:value={middleware.rateLimit.period}
				on:input={validate}
				placeholder="1s"
				{disabled}
			/>
			{#if errors.period}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.period}</div>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="burst" class="text-right">Burst</Label>
		<div class="relative col-span-3">
			<Input
				id="burst"
				name="burst"
				type="number"
				bind:value={middleware.rateLimit.burst}
				on:input={validate}
				placeholder="1"
				{disabled}
			/>
			{#if errors.burst}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.burst}</div>
			{/if}
		</div>
	</div>

	<header class="border-b border-gray-200 py-2 font-bold">Source Criterion</header>

	{#if middleware.rateLimit.sourceCriterion && middleware.rateLimit.sourceCriterion.ipStrategy}
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="depth" class="text-right">Depth</Label>
			<div class="relative col-span-3">
				<Input
					id="depth"
					name="depth"
					type="number"
					bind:value={middleware.rateLimit.sourceCriterion.ipStrategy.depth}
					placeholder="0"
					on:input={validate}
					{disabled}
				/>
				{#if errors.ipStrategy}
					<div class="col-span-4 text-right text-sm text-red-500">{errors.ipStrategy}</div>
				{/if}
			</div>
		</div>
		<div class="grid grid-cols-4 items-center gap-2">
			<Label for="request-header-name" class="text-right">Request Header Name</Label>
			<div class="relative col-span-3">
				<Input
					id="request-header-name"
					name="request-header-name"
					type="text"
					bind:value={middleware.rateLimit.sourceCriterion.requestHeaderName}
					on:input={validate}
					placeholder="username"
					{disabled}
				/>
				{#if errors.requestHeaderName}
					<div class="col-span-4 text-right text-sm text-red-500">{errors.requestHeaderName}</div>
				{/if}
			</div>
		</div>
		<div class="grid grid-cols-4 items-center gap-2">
			<Label for="request-host" class="text-right">Request Host</Label>
			<Switch
				id="request-host"
				bind:checked={middleware.rateLimit.sourceCriterion.requestHost}
				class="col-span-3"
				{disabled}
			/>
		</div>
		<ArrayInput
			bind:items={middleware.rateLimit.sourceCriterion.ipStrategy.excludedIPs}
			label="Excluded IPs"
			placeholder="192.168.1.1"
			on:update={validate}
			{disabled}
		/>
		{#if errors.sourceCriterion}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.sourceCriterion}</div>
		{/if}
	{/if}
{/if}
