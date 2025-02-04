<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import {
		type Middleware,
		type SupportedMiddleware,
		type SupportedMiddlewareHTTP,
		type SupportedMiddlewareTCP,
		type UpsertMiddlewareParams,
		getDefaultValuesForType,
		getMetadataForMiddleware,
		getTCPDefaultValuesForType,
		HTTPMiddlewareKeys,
		TCPMiddlewareKeys
	} from '$lib/types/middlewares';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import DynamicForm from '../forms/DynamicForm.svelte';
	import PluginForm from '../forms/PluginForm.svelte';

	interface Props {
		middleware: Middleware;
		open?: boolean;
		mode?: 'create' | 'edit';
		disabled?: boolean;
	}
	let {
		middleware = $bindable(),
		open = $bindable(false),
		mode = 'create',
		disabled
	}: Props = $props();

	type FormData = Record<string, unknown>;
	let currentFormData = $state<FormData>({});

	// Helper to extract middleware config
	function extractMiddlewareConfig(mw: Middleware): FormData {
		if (!mw.type) return {};

		// Get the default structure for the middleware type
		const defaultConfig =
			mw.protocol === 'http'
				? getDefaultValuesForType(mw.type as SupportedMiddlewareHTTP)
				: getTCPDefaultValuesForType(mw.type as SupportedMiddlewareTCP);

		// Get the current config values
		const currentConfig: FormData = {};
		for (const [key, value] of Object.entries(mw)) {
			if (!['name', 'protocol', 'type'].includes(key)) {
				currentConfig[key] = value;
			}
		}

		// Deep merge the default structure with the current config
		const mergeObjects = (target: Record<string, unknown>, source: Record<string, unknown>) => {
			for (const key of Object.keys(source)) {
				if (source[key] instanceof Object && !Array.isArray(source[key])) {
					target[key] = target[key] || {};
					mergeObjects(
						target[key] as Record<string, unknown>,
						source[key] as Record<string, unknown>
					);
				} else {
					// Only use default if there's no current value
					if (!(key in target)) {
						target[key] = source[key];
					}
				}
			}
		};

		// Start with current config and fill in missing values from defaults
		const result = { ...currentConfig } as FormData;
		mergeObjects(result, defaultConfig as Record<string, unknown>);

		return result;
	}

	// Initialize or update form data when middleware changes or modal opens
	$effect(() => {
		if (open && middleware.type) {
			const existingConfig = extractMiddlewareConfig(middleware);

			if (Object.keys(existingConfig).length > 0) {
				currentFormData = existingConfig as FormData;
			} else {
				if (middleware.protocol === 'http') {
					currentFormData = getDefaultValuesForType(
						middleware.type as SupportedMiddlewareHTTP
					) as FormData;
				} else {
					currentFormData = getTCPDefaultValuesForType(
						middleware.type as SupportedMiddlewareTCP
					) as FormData;
				}
			}
		}
	});

	const handleTypeChange = (value: string) => {
		middleware.type = value as SupportedMiddleware;
	};
	const handleProtocolChange = () => {
		middleware.protocol = middleware.protocol === 'http' ? 'tcp' : 'http';
		middleware.type = undefined;
	};

	let isValid = $derived(!!middleware.name && !!middleware.type);
	const onSubmit = async (formData: Record<string, unknown>) => {
		if (!isValid) {
			toast.error('Please fill in all required fields');
			return;
		}

		try {
			const params: UpsertMiddlewareParams = {
				name: middleware.name,
				protocol: middleware.protocol,
				type: middleware.type,
				middleware: {},
				tcpMiddleware: {}
			};

			if (middleware.protocol === 'http' && middleware.type) {
				params.middleware = {
					[middleware.type]: structuredClone(formData)
				};
			} else if (middleware.protocol === 'tcp' && middleware.type) {
				params.tcpMiddleware = {
					[middleware.type]: structuredClone(formData)
				};
			}

			await api.upsertMiddleware(params);
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
			<Dialog.Title>{mode === 'edit' ? 'Update' : 'Add'} Middleware</Dialog.Title>
			<Dialog.Description>Configure your Traefik middleware</Dialog.Description>
		</Dialog.Header>

		<div class="flex flex-col gap-4 py-2">
			{#if mode === 'create'}
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
								{middleware.type ? middleware.type : 'Select type'}
							</Select.Trigger>
							<Select.Content>
								{#if middleware.protocol === 'http'}
									{#each HTTPMiddlewareKeys as type}
										{#if type.value !== 'plugin'}
											<Select.Item value={type.value}>{type.label}</Select.Item>
										{/if}
									{/each}
								{:else if middleware.protocol === 'tcp'}
									{#each TCPMiddlewareKeys as type}
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
			{/if}
			<div class="flex flex-col gap-2">
				<Label for="name">Name</Label>
				<Input type="text" name="name" bind:value={middleware.name} required {disabled} />
			</div>
		</div>

		<Separator />

		{#if middleware.type}
			{#if middleware.type === 'plugin'}
				<PluginForm bind:data={currentFormData} {onSubmit} {disabled} />
			{:else}
				<DynamicForm
					data={currentFormData}
					metadata={getMetadataForMiddleware(middleware.type)}
					{onSubmit}
					{disabled}
				/>
			{/if}
		{/if}
	</Dialog.Content>
</Dialog.Root>
