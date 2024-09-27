<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';
	import { CustomTimeUnitSchemaOptional } from '../utils/validation';

	export let middleware: Middleware;
	export let disabled = false;

	const circuitBreakerSchema = z.object({
		expression: z
			.string({ required_error: 'Expression is required' })
			.trim()
			.min(1, 'Expression is required')
			.default('LatencyAtQuantileMS(50.0) > 100'),
		checkPeriod: CustomTimeUnitSchemaOptional,
		fallbackDuration: CustomTimeUnitSchemaOptional,
		recoveryDuration: CustomTimeUnitSchemaOptional,
		responseCode: z.coerce.number().int().nonnegative().optional()
	});
	middleware.circuitBreaker = circuitBreakerSchema.parse({ ...middleware.circuitBreaker });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.circuitBreaker = circuitBreakerSchema.parse(middleware.circuitBreaker); // Parse the circuitBreaker object
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};
</script>

{#if middleware.circuitBreaker}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="expression" class="text-right">Expression</Label>
		<div class="relative col-span-3">
			<Input
				id="expression"
				name="expression"
				type="text"
				bind:value={middleware.circuitBreaker.expression}
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
				bind:value={middleware.circuitBreaker.checkPeriod}
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
				bind:value={middleware.circuitBreaker.fallbackDuration}
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
				bind:value={middleware.circuitBreaker.recoveryDuration}
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
				bind:value={middleware.circuitBreaker.responseCode}
				on:input={validate}
				placeholder="503"
				{disabled}
			/>
			{#if errors.responseCode}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.responseCode}</div>
			{/if}
		</div>
	</div>
{/if}
