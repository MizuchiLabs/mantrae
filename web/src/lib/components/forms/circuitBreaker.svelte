<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';
	import { CustomTimeUnitSchemaOptional } from '../utils/validation';
	import { onDestroy } from 'svelte';

	export let middleware: Middleware;
	export let disabled = false;

	const schema = z.object({
		expression: z
			.string()
			.trim()
			.min(1, 'Expression is required')
			.default('LatencyAtQuantileMS(50.0) > 100'),
		checkPeriod: CustomTimeUnitSchemaOptional,
		fallbackDuration: CustomTimeUnitSchemaOptional,
		recoveryDuration: CustomTimeUnitSchemaOptional,
		responseCode: z
			.union([z.string(), z.number()])
			.transform((value) => (value === '' ? null : Number(value)))
			.nullish()
	});
	middleware.content = schema.parse({ ...middleware.content });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.content = schema.parse(middleware.content);
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};

	onDestroy(() => {
		validate();
	});
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="expression" class="text-right">Expression</Label>
	<div class="relative col-span-3">
		<Input
			id="expression"
			name="expression"
			type="text"
			bind:value={middleware.content.expression}
			on:input={validate}
			placeholder="LatencyAtQuantileMS(50.0) > 100"
			{disabled}
		/>
		{#if errors.expression}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.expression}</div>
		{/if}
	</div>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="check-period" class="text-right">Check Period</Label>
	<div class="relative col-span-3">
		<Input
			id="check-period"
			name="check-period"
			type="text"
			bind:value={middleware.content.checkPeriod}
			on:input={validate}
			placeholder="100ms"
			{disabled}
		/>
		{#if errors.checkPeriod}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.checkPeriod}</div>
		{/if}
	</div>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="fallback-duration" class="text-right">Fallback Duration</Label>
	<div class="relative col-span-3">
		<Input
			id="fallback-duration"
			name="fallback-duration"
			type="text"
			bind:value={middleware.content.fallbackDuration}
			on:input={validate}
			placeholder="10s"
			{disabled}
		/>
		{#if errors.fallbackDuration}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.fallbackDuration}</div>
		{/if}
	</div>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="recovery-duration" class="text-right">Recovery Duration</Label>
	<div class="relative col-span-3">
		<Input
			id="recovery-duration"
			name="recovery-duration"
			type="text"
			bind:value={middleware.content.recoveryDuration}
			on:input={validate}
			placeholder="10s"
			{disabled}
		/>
		{#if errors.recoveryDuration}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.recoveryDuration}</div>
		{/if}
	</div>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="response-code" class="text-right">Response Code</Label>
	<div class="relative col-span-3">
		<Input
			id="response-code"
			name="response-code"
			type="number"
			bind:value={middleware.content.responseCode}
			on:input={validate}
			placeholder="503"
			{disabled}
		/>
		{#if errors.responseCode}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.responseCode}</div>
		{/if}
	</div>
</div>
