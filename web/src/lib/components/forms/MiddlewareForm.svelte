<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Middleware } from '$lib/types/middlewares';
	import { MIDDLEWARE_REGISTRY } from './mw-registry';

	export let selectedType: string | undefined;
	export let value: Partial<Middleware> = {};

	function updateField(fieldName: string, fieldValue: any) {
		value = { ...value, [fieldName]: fieldValue };
	}

	const middlewareTypes = Object.entries(MIDDLEWARE_REGISTRY).map(([key, config]) => ({
		value: key,
		label: config.label
	}));
</script>

<div class="flex items-center gap-2">
	<Select.Root>
		<Select.Trigger class="w-[180px]">Select a middleware type</Select.Trigger>
		<Select.Content>
			{#each middlewareTypes as type}
				<Select.Item value={type.value}>{type.label}</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>

	{#if selectedType && MIDDLEWARE_REGISTRY[selectedType]}
		<div class="space-y-4">
			{#each Object.entries(MIDDLEWARE_REGISTRY[selectedType].fields) as [fieldName, field]}
				<div class="space-y-2">
					<label for={fieldName} class="text-sm font-medium">
						{field.label}
					</label>

					{#if field.type === 'text'}
						<Input
							type="text"
							id={fieldName}
							value={value[fieldName] ?? field.defaultValue ?? ''}
							oninput={(e) => updateField(fieldName, e.currentTarget.value)}
							required={field.required}
						/>
					{:else if field.type === 'boolean'}
						<Checkbox
							id={fieldName}
							checked={value[fieldName] ?? field.defaultValue ?? false}
							onchange={(e) => updateField(fieldName, e.currentTarget.checked)}
						/>
					{:else if field.type === 'array'}
						<!-- Implement array input component -->
					{/if}

					{#if field.description}
						<p class="text-sm text-muted-foreground">{field.description}</p>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
