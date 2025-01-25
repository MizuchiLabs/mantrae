<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import BaseForm from '$lib/components/forms/BaseForm.svelte';
	import type { Middleware, UpsertMiddlewareParams } from '$lib/types/middlewares';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { api, profile, mwNames } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import {
		type SupportedMiddleware,
		GetSchema,
		MiddlewareTypes,
		type ZodObjectOrRecord
	} from '$lib/components/forms/mw_registry';
	import Separator from '../ui/separator/separator.svelte';
	import { z } from 'zod';

	interface Props {
		middleware?: Middleware;
		open?: boolean;
		mode: 'create' | 'edit' | 'view';
	}

	const defaultMiddleware: Middleware = {
		name: '',
		protocol: 'http'
	};

	let {
		middleware = $bindable(defaultMiddleware),
		open = $bindable(false),
		mode = 'view'
	}: Props = $props();

	let disabled = $state(mode === 'view');

	let schema: ZodObjectOrRecord = $state(GetSchema(middleware.type));

	const handleSelect = (value: string) => {
		middleware.type = value as SupportedMiddleware;
		schema = GetSchema(middleware.type);
	};

	function getBaseType(fieldSchema: z.ZodTypeAny | unknown) {
		if (fieldSchema instanceof z.ZodOptional || fieldSchema instanceof z.ZodDefault) {
			return fieldSchema._def.innerType;
		}
		return fieldSchema;
	}

	const onSubmit = async (data: FormData) => {
		let formData: Record<string, unknown> = {};

		// Special handling for arrays from FormData
		const entries = Array.from(data.entries());
		entries.forEach(([key, value]) => {
			// Check if the key is for an array (ends with [])
			if (key.endsWith('[]')) {
				const cleanKey = key.replace('[]', '');
				if (!formData[cleanKey]) {
					formData[cleanKey] = [];
				}
				(formData[cleanKey] as unknown[]).push(value);
			} else {
				formData[key] = value;
			}
		});

		// Coerce number fields based on schema
		if (schema instanceof z.ZodObject) {
			Object.entries(schema.shape).forEach(([key, field]) => {
				const baseType = getBaseType(field);
				if (baseType instanceof z.ZodNumber && typeof formData[key] === 'string') {
					formData[key] = Number(formData[key]);
				}
				// Handle array fields specifically for chain middleware
				if (baseType instanceof z.ZodArray && Array.isArray(formData[key])) {
					// Ensure the array is preserved
					formData[key] = [...formData[key]];
				}
			});
		}

		middleware = {
			...middleware,
			[middleware.type as string]: formData
		};
		try {
			// Ensure proper name formatting and synchronization
			if (!middleware.name) return;
			if (!middleware.name.includes('@')) {
				middleware.name = `${middleware.name}@http`;
			}

			let params: UpsertMiddlewareParams = {
				name: middleware.name,
				protocol: middleware.protocol,
				type: middleware.type,
				...(middleware.protocol === 'http' ? { middleware } : { tcpMiddleware: middleware })
			};

			await api.upsertMiddleware($profile.id, params);
			open = false;
			toast.success(`Middleware ${mode === 'create' ? 'created' : 'updated'} successfully`);
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(`Failed to ${mode} router`, {
				description: e.message
			});
		}
	};

	$effect(() => {
		if (open && middleware.type) {
			schema = GetSchema(middleware.type);
		}
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{middleware.name ? 'Update' : 'Add'} Middleware</Dialog.Title>
			<Dialog.Description>Configure your Traefik instance connection details.</Dialog.Description>
		</Dialog.Header>

		<div class="flex flex-col gap-4 py-2">
			<div class="flex flex-col gap-2">
				<span class="text-sm text-primary">Middleware Type</span>
				<Select.Root
					type="single"
					bind:value={middleware.type}
					onValueChange={handleSelect}
					{disabled}
				>
					<Select.Trigger>
						{middleware.type
							? MiddlewareTypes.find((t) => t.value === middleware.type)?.label
							: 'Select a middleware type'}
					</Select.Trigger>
					<Select.Content>
						{#each MiddlewareTypes as type}
							<Select.Item value={type.value}>{type.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
			<div class="flex flex-col gap-2">
				<Label for="name">Name</Label>
				<Input type="text" name="name" bind:value={middleware.name} required {disabled} />
			</div>
		</div>

		<Separator />

		{#if schema}
			{#if middleware.type === 'chain'}
				<BaseForm {schema} data={middleware} subData={$mwNames} subMultiple={true} {onSubmit} />
			{:else}
				<BaseForm {schema} data={middleware} {onSubmit} />
			{/if}
		{/if}
	</Dialog.Content>
</Dialog.Root>
