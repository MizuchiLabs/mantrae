<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';
	import { onDestroy } from 'svelte';

	export let middleware: Middleware;
	export let disabled = false;

	const schema = z.object({
		amount: z
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
