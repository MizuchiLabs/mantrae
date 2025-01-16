<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { z } from 'zod';
	import { onDestroy } from 'svelte';

	interface Props {
		middleware: Middleware;
		disabled?: boolean;
	}

	let { middleware = $bindable(), disabled = false }: Props = $props();

	const schema = z.object({
		maxRequestBodyBytes: z
			.union([z.string(), z.number()])
			.transform((value) => (value === '' ? null : Number(value)))
			.nullish(),
		memRequestBodyBytes: z
			.union([z.string(), z.number()])
			.transform((value) => (value === '' ? null : Number(value)))
			.nullish(),
		maxResponseBodyBytes: z
			.union([z.string(), z.number()])
			.transform((value) => (value === '' ? null : Number(value)))
			.nullish(),
		memResponseBodyBytes: z
			.union([z.string(), z.number()])
			.transform((value) => (value === '' ? null : Number(value)))
			.nullish(),
		retryExpression: z.string().trim().nullish()
	});
	middleware.content = schema.parse({ ...middleware.content });

	let errors: Record<any, string[] | undefined> = $state({});
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
	<Label for="max-request-body-bytes" class="text-right">Max Request Body Bytes</Label>
	<div class="relative col-span-3">
		<Input
			id="max-request-body-bytes"
			name="max-request-body-bytes"
			type="number"
			bind:value={middleware.content.maxRequestBodyBytes}
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
			bind:value={middleware.content.memRequestBodyBytes}
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
			bind:value={middleware.content.maxResponseBodyBytes}
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
			bind:value={middleware.content.memResponseBodyBytes}
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
			bind:value={middleware.content.retryExpression}
			on:input={validate}
			placeholder="IsNetworkError() && Attempts() < 2"
			{disabled}
		/>
		{#if errors.retryExpression}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.retryExpression}</div>
		{/if}
	</div>
</div>
