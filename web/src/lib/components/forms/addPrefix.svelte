<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';

	export let middleware: Middleware;
	export let disabled = false;

	const addPrefixSchema = z.object({
		prefix: z
			.string({ required_error: 'Prefix is required' })
			.trim()
			.min(1, 'Prefix is required')
			.default('/foo')
	});
	middleware.addPrefix = addPrefixSchema.parse({ ...middleware.addPrefix });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.addPrefix = addPrefixSchema.parse(middleware.addPrefix); // Parse the addPrefix object
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};
</script>

{#if middleware.addPrefix}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="prefix" class="text-right">Prefix</Label>
		<div class="relative col-span-3">
			<Input
				type="text"
				id="prefix"
				placeholder="/foo"
				bind:value={middleware.addPrefix.prefix}
				on:input={validate}
				{disabled}
			/>
			{#if errors.prefix}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.prefix}</div>
			{/if}
		</div>
	</div>
{/if}
