<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import { z } from 'zod';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.compress = {
		excludedContentTypes: [],
		includeContentTypes: [],
		defaultEncoding: '',
		...middleware.compress
	};

	const compressSchema = z.object({
		minResponseBodyBytes: z.coerce
			.number({ required_error: 'Min Response Body Bytes is required' })
			.int()
			.nonnegative(),
		defaultEncoding: z.string({ required_error: 'Default Encoding is required' }).trim().optional(),
		excludedContentTypes: z
			.array(z.string({ required_error: 'Excluded Content Types is required' }).trim())
			.optional(),
		includeContentTypes: z
			.array(z.string({ required_error: 'Include Content Types is required' }).trim())
			.optional()
	});

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			compressSchema.parse(middleware.compress); // Parse the compress object
			errors = {};
			return { isValid: true, errors: null };
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
				return { isValid: false, errors: err.flatten().fieldErrors };
			}
			return { isValid: false, errors: { general: ['Unexpected error'] } };
		}
	};
</script>

{#if middleware.compress}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="min-response-body-bytes" class="text-right">Min Response Body Bytes</Label>
		<div class="relative col-span-3">
			<Input
				id="min-response-body-bytes"
				name="min-response-body-bytes"
				type="number"
				bind:value={middleware.compress.minResponseBodyBytes}
				on:input={validate}
				placeholder="1024"
				{disabled}
			/>
			{#if errors.minResponseBodyBytes}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.minResponseBodyBytes}</div>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="default-encoding" class="text-right">Default Encoding</Label>
		<div class="relative col-span-3">
			<Input
				id="default-encoding"
				name="default-encoding"
				type="text"
				bind:value={middleware.compress.defaultEncoding}
				on:input={validate}
				placeholder="gzip"
				{disabled}
			/>
			{#if errors.defaultEncoding}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.defaultEncoding}</div>
			{/if}
		</div>
	</div>
	<ArrayInput
		bind:items={middleware.compress.excludedContentTypes}
		label="Excluded Content Types"
		placeholder="text/event-stream"
		{disabled}
	/>
	<ArrayInput
		bind:items={middleware.compress.includeContentTypes}
		label="Include Content Types"
		placeholder="application/json"
		{disabled}
	/>
{/if}
