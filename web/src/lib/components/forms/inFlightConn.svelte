<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';

	export let middleware: Middleware;
	export let disabled = false;

	const inFlightConnSchema = z.object({
		amount: z.coerce
			.number({ required_error: 'Amount is required' })
			.int()
			.nonnegative()
			.default(50)
	});
	middleware.content = inFlightConnSchema.parse({
		...middleware.content
	});

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.content = inFlightConnSchema.parse(middleware.content);
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="amount" class="text-right">Amount</Label>
	<div class="relative col-span-3">
		<Input
			id="amount"
			name="amount"
			type="number"
			bind:value={middleware.content.amount}
			on:input={validate}
			placeholder="50"
			{disabled}
		/>
		{#if errors.amount}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.amount}</div>
		{/if}
	</div>
</div>
