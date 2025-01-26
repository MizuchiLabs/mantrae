<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import BaseForm from '$lib/components/forms/BaseForm.svelte';
	import type { Middleware, UpsertMiddlewareParams } from '$lib/types/middlewares';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api, profile, mwNames } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import {
		type SupportedMiddleware,
		GetSchema,
		type ZodObjectOrRecord,
		HTTPMiddlewareTypes,
		TCPMiddlewareTypes,
		MiddlewareTypes
	} from '$lib/components/forms/mw_registry';
	import Separator from '../ui/separator/separator.svelte';
	import { z } from 'zod';

	interface Props {
		middleware: Middleware;
		open?: boolean;
	}

	let { middleware = $bindable(), open = $bindable(false) }: Props = $props();

	let mwProvider = $derived(middleware.name ? middleware.name.split('@')[1] : 'http');
	let disabled = $derived(mwProvider !== 'http' && mwProvider !== undefined);
	let schema: ZodObjectOrRecord = $derived(
		middleware.type ? GetSchema(middleware.type) : z.object({})
	);

	const handleTypeChange = (value: string) => {
		middleware.type = value as SupportedMiddleware;
	};
	const handleProtocolChange = () => {
		middleware.protocol = middleware.protocol === 'http' ? 'tcp' : 'http';
		middleware.type = undefined; // Reset type when protocol changes
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

		// TODO: Handle plugin name change
		if (middleware.type === 'plugin') {
			const pluginName = middleware.name.split('@')[0];
			middleware = {
				name: pluginName,
				protocol: 'http',
				type: 'plugin',
				plugin: {
					[pluginName]: JSON.parse((formData[pluginName] as string) ?? '{}')
				}
			};
		} else {
			middleware = {
				...middleware,
				[middleware.type as string]: formData
			};
		}
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
			toast.success(`Middleware updated successfully`);
			resetForm();
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(`Failed to update middleware`, {
				description: e.message
			});
		}
	};

	function resetForm() {
		open = false;
		middleware = {
			name: '',
			protocol: 'http',
			type: undefined
		};
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{middleware.name ? 'Update' : 'Add'} Middleware</Dialog.Title>
			<Dialog.Description>Configure your Traefik instance connection details.</Dialog.Description>
		</Dialog.Header>

		<div class="flex flex-col gap-4 py-2">
			<div class="flex flex-col gap-2">
				<span class="text-sm text-primary">Middleware Type & Protocol</span>
				<div class="flex items-center gap-4">
					<Select.Root
						type="single"
						bind:value={middleware.type}
						onValueChange={handleTypeChange}
						{disabled}
					>
						<Select.Trigger>
							{middleware.type
								? MiddlewareTypes.find((t) => t.value === middleware.type)?.label
								: 'Select a middleware type'}
						</Select.Trigger>
						<Select.Content>
							{#if middleware.protocol === 'http'}
								{#each HTTPMiddlewareTypes as type}
									<Select.Item value={type.value}>{type.label}</Select.Item>
								{/each}
							{:else if middleware.protocol === 'tcp'}
								{#each TCPMiddlewareTypes as type}
									<Select.Item value={type.value}>{type.label}</Select.Item>
								{/each}
							{/if}
						</Select.Content>
					</Select.Root>
					<div class="flex items-center gap-2">
						<Label for="protocol">{middleware.protocol.toUpperCase()}</Label>
						<Switch
							name="protocol"
							checked={middleware.protocol === 'http'}
							onCheckedChange={handleProtocolChange}
							{disabled}
						/>
					</div>
				</div>
			</div>
			<div class="flex flex-col gap-2">
				<Label for="name">Name</Label>
				<Input type="text" name="name" bind:value={middleware.name} required {disabled} />
			</div>
		</div>

		<Separator />

		{#if schema}
			{#if middleware.protocol === 'http' && middleware.type === 'chain'}
				<BaseForm {schema} data={middleware} subData={$mwNames} subMultiple={true} {onSubmit} />
			{:else}
				<BaseForm {schema} data={middleware} {onSubmit} />
			{/if}
		{/if}
	</Dialog.Content>
</Dialog.Root>
