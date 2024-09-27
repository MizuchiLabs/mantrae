<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';
	import { CustomTimeUnitSchemaOptional } from '../utils/validation';

	export let middleware: Middleware;
	export let disabled = false;

	const retrySchema = z.object({
		attempts: z.coerce
			.number({ required_error: 'Attempts is required' })
			.int()
			.nonnegative()
			.default(3),
		initialInterval: CustomTimeUnitSchemaOptional
	});
	middleware.retry = retrySchema.parse({ ...middleware.retry });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.retry = retrySchema.parse(middleware.retry);
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};
</script>

{#if middleware.retry}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="attempts" class="text-right">Attempts</Label>
		<div class="relative col-span-3">
			<Input
				id="attempts"
				name="attempts"
				type="number"
				bind:value={middleware.retry.attempts}
				on:input={validate}
				placeholder="3"
				min="0"
				{disabled}
			/>
			{#if errors.attempts}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.attempts}</div>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="initial-interval" class="text-right">Initial Interval</Label>
		<div class="relative col-span-3">
			<Input
				id="initial-interval"
				name="initial-interval"
				type="text"
				bind:value={middleware.retry.initialInterval}
				on:input={validate}
				placeholder="100ms"
				{disabled}
			/>
			{#if errors.initialInterval}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.initialInterval}</div>
			{/if}
		</div>
	</div>
{/if}
