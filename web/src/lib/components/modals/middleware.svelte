<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import {
		type Middleware,
		type SupportedMiddlewareHTTP,
		type SupportedMiddlewareTCP,
		type UpsertMiddlewareParams,
		getDefaultValuesForType,
		getTCPDefaultValuesForType,
		HTTPMiddlewareKeys,
		TCPMiddlewareKeys
	} from '$lib/types/middlewares';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api, profile } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import DynamicForm from '../forms/DynamicForm.svelte';
	import { safeClone } from '$lib/utils';

	interface Props {
		middleware: Middleware;
		open?: boolean;
		disabled?: boolean;
	}
	let { middleware = $bindable(), open = $bindable(false), disabled }: Props = $props();

	let currentFormData = $state<Record<string, unknown>>({});

	// Helper to extract middleware config
	function extractMiddlewareConfig(mw: Middleware): Record<string, unknown> {
		if (!mw.type) return {};

		// Filter out the base properties to get the config
		const config: Record<string, unknown> = {};
		for (const [key, value] of Object.entries(mw)) {
			if (!['name', 'protocol', 'type'].includes(key)) {
				config[key] = value;
			}
		}
		return config;
	}

	// Initialize or update form data when middleware changes or modal opens
	$effect(() => {
		if (open && middleware.type) {
			const existingConfig = extractMiddlewareConfig(middleware);

			if (Object.keys(existingConfig).length > 0) {
				currentFormData = safeClone(existingConfig);
			} else {
				// Fall back to default values if no existing data
				currentFormData =
					middleware.protocol === 'http'
						? getDefaultValuesForType(middleware.type as SupportedMiddlewareHTTP)
						: getTCPDefaultValuesForType(middleware.type as SupportedMiddlewareTCP);
			}
		}
	});

	const handleTypeChange = (value: string) => {
		middleware.type = value as SupportedMiddlewareHTTP | SupportedMiddlewareTCP;
	};
	const handleProtocolChange = () => {
		middleware.protocol = middleware.protocol === 'http' ? 'tcp' : 'http';
		middleware.type = undefined;
	};

	const onSubmit = async (formData: Record<string, unknown>) => {
		if (!middleware.type) return;

		const data = {
			...middleware,
			[middleware.type]: structuredClone(formData)
		};

		try {
			// Ensure proper name formatting and synchronization
			if (!data.name) {
				toast.error('Middleware name is required');
				return;
			}
			if (!data.name.includes('@')) {
				middleware.name = `${middleware.name}@http`;
			}

			const params: UpsertMiddlewareParams = {
				name: data.name,
				protocol: data.protocol,
				type: data.type,
				...(data.protocol === 'http' ? { middleware: data } : { tcpMiddleware: data })
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
			<Dialog.Description>Configure your Traefik middleware</Dialog.Description>
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
							{middleware.type ? middleware.type : 'Select type'}
						</Select.Trigger>
						<Select.Content>
							{#if middleware.protocol === 'http'}
								{#each HTTPMiddlewareKeys as type}
									<Select.Item value={type.value}>{type.label}</Select.Item>
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
			<div class="flex flex-col gap-2">
				<Label for="name">Name</Label>
				<Input type="text" name="name" bind:value={middleware.name} required {disabled} />
			</div>
		</div>

		<Separator />

		{#if middleware.type}
			<DynamicForm data={currentFormData} {onSubmit} {disabled} />
		{/if}
	</Dialog.Content>
</Dialog.Root>
