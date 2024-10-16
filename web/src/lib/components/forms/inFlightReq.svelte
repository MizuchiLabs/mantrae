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

	const inFlightReqSchema = z.object({
		amount: z.coerce
			.number({ required_error: 'Amount is required' })
			.int()
			.nonnegative()
			.default(50),
		sourceCriterion: z
			.object({
				ipStrategy: z
					.object({
						depth: z.coerce.number().int().nonnegative().optional(),
						excludedIPs: z.array(CustomIPSchemaOptional).optional()
					})
					.default({}),
				requestHeaderName: z.string().trim().optional(),
				requestHost: z.boolean().optional()
			})
			.default({})
	});
	middleware.inFlightReq = inFlightReqSchema.parse({ ...middleware.inFlightReq });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.inFlightReq = inFlightReqSchema.parse(middleware.inFlightReq); // Parse the inFlightReq object
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};
</script>

{#if middleware.inFlightReq}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="amount" class="text-right">Amount</Label>
		<div class="relative col-span-3">
			<Input
				id="amount"
				name="amount"
				type="number"
				bind:value={middleware.inFlightReq.amount}
				on:input={validate}
				placeholder="50"
				{disabled}
			/>
			{#if errors.amount}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.amount}</div>
			{/if}
		</div>
	</div>

	<header class="border-b border-gray-200 py-2 font-bold">Source Criterion</header>
	{#if middleware.inFlightReq.sourceCriterion && middleware.inFlightReq.sourceCriterion.ipStrategy}
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="depth" class="text-right">Depth</Label>
			<div class="relative col-span-3">
				<Input
					id="depth"
					name="depth"
					type="number"
					bind:value={middleware.inFlightReq.sourceCriterion.ipStrategy.depth}
					on:input={validate}
					placeholder="0"
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
					bind:value={middleware.inFlightReq.sourceCriterion.requestHeaderName}
					on:input={validate}
					placeholder="X-CustomHeader"
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
				bind:checked={middleware.inFlightReq.sourceCriterion.requestHost}
				class="col-span-3"
				{disabled}
			/>
		</div>
		<ArrayInput
			bind:items={middleware.inFlightReq.sourceCriterion.ipStrategy.excludedIPs}
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
