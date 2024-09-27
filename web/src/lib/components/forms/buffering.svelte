<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.buffering = {
		retryExpression: '',
		...middleware.buffering
	};

	const bufferingSchema = z.object({
		maxRequestBodyBytes: z.coerce
			.number({ required_error: 'Max Request Body Bytes is required' })
			.int()
			.nonnegative()
			.optional(),
		memRequestBodyBytes: z.coerce
			.number({ required_error: 'Mem Request Body Bytes is required' })
			.int()
			.nonnegative()
			.optional(),
		maxResponseBodyBytes: z.coerce
			.number({ required_error: 'Max Response Body Bytes is required' })
			.int()
			.nonnegative()
			.optional(),
		memResponseBodyBytes: z.coerce
			.number({ required_error: 'Mem Response Body Bytes is required' })
			.int()
			.nonnegative()
			.optional(),
		retryExpression: z.string({ required_error: 'Retry Expression is required' }).optional()
	});

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			bufferingSchema.parse(middleware.buffering); // Parse the buffering object
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

{#if middleware.buffering}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="max-request-body-bytes" class="text-right">Max Request Body Bytes</Label>
		<div class="relative col-span-3">
			<Input
				id="max-request-body-bytes"
				name="max-request-body-bytes"
				type="number"
				bind:value={middleware.buffering.maxRequestBodyBytes}
				on:input={validate}
				placeholder="0"
				{disabled}
			/>
			{#if errors.maxRequestBodyBytes}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.maxRequestBodyBytes}</div>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="mem-request-body-bytes" class="text-right">Mem Request Body Bytes</Label>
		<div class="relative col-span-3">
			<Input
				id="mem-request-body-bytes"
				name="mem-request-body-bytes"
				type="number"
				bind:value={middleware.buffering.memRequestBodyBytes}
				on:input={validate}
				placeholder="1048576"
				{disabled}
			/>
			{#if errors.memRequestBodyBytes}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.memRequestBodyBytes}</div>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="max-response-body-bytes" class="text-right">Max Response Body Bytes</Label>
		<div class="relative col-span-3">
			<Input
				id="max-response-body-bytes"
				name="max-response-body-bytes"
				type="number"
				bind:value={middleware.buffering.maxResponseBodyBytes}
				on:input={validate}
				placeholder="0"
				{disabled}
			/>
			{#if errors.maxResponseBodyBytes}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.maxResponseBodyBytes}</div>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="mem-response-body-bytes" class="text-right">Mem Response Body Bytes</Label>
		<div class="relative col-span-3">
			<Input
				id="mem-response-body-bytes"
				name="mem-response-body-bytes"
				type="number"
				bind:value={middleware.buffering.memResponseBodyBytes}
				on:input={validate}
				placeholder="1048576"
				{disabled}
			/>
			{#if errors.memResponseBodyBytes}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.memResponseBodyBytes}</div>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="retry-expression" class="text-right">Retry Expression</Label>
		<div class="relative col-span-3">
			<Input
				id="retry-expression"
				name="retry-expression"
				type="text"
				bind:value={middleware.buffering.retryExpression}
				on:input={validate}
				placeholder="IsNetworkError() && Attempts() < 2"
				{disabled}
			/>
			{#if errors.retryExpression}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.retryExpression}</div>
			{/if}
		</div>
	</div>
{/if}
