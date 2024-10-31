<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';

	export let middleware: Middleware;
	export let disabled = false;

	// Define validation schema for addPrefix content
	const addPrefixSchema = z.object({
		prefix: z
			.string({ required_error: 'Prefix is required' })
			.trim()
			.min(1, 'Prefix is required')
			.default('/foo')
	});

	// Parse and validate middleware.content for addPrefix
	let errors: Record<string, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.content = addPrefixSchema.parse(middleware.content);
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};

	// Initial validation
	middleware.content = addPrefixSchema.parse({ ...middleware.content });
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="prefix" class="text-right">Prefix</Label>
	<div class="relative col-span-3">
		<Input
			type="text"
			id="prefix"
			placeholder="/foo"
			bind:value={middleware.content.prefix}
			on:input={validate}
			{disabled}
		/>
		{#if errors.prefix}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.prefix}</div>
		{/if}
	</div>
</div>
