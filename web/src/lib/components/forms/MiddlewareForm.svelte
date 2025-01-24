<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Middleware } from '$lib/types/middlewares';
	import { MIDDLEWARE_REGISTRY, type SupportedMiddleware } from './mw-registry';

	type Props = {
		middleware: Partial<Middleware>;
	};

	let { middleware }: Props = $props();

	let selectedType: SupportedMiddleware | undefined = $state();

	function updateField<K extends keyof Middleware>(fieldName: K, fieldValue: Middleware[K]) {
		middleware = { ...middleware, [fieldName]: fieldValue };
	}

	const middlewareTypes = Object.entries(MIDDLEWARE_REGISTRY).map(([key, config]) => ({
		value: key,
		label: config.label
	}));

	function handleArrayInput(fieldName: string, value: string) {
		const arrayValue = value.split(',').map((item) => item.trim());
		updateField(fieldName as keyof Middleware, arrayValue as never);
	}

	$effect(() => {
		if (!selectedType) return;
		let registry = MIDDLEWARE_REGISTRY[selectedType];
		console.log(registry);
		console.log(Object.entries(registry.fields));
	});
</script>

<div class="flex items-center gap-2">
	<Select.Root type="single" bind:value={selectedType}>
		<Select.Trigger class="w-[380px]">
			{selectedType ? MIDDLEWARE_REGISTRY[selectedType].label : 'Select a middleware type'}
		</Select.Trigger>
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
							value={middleware[fieldName as keyof Middleware] ?? ''}
							oninput={(e) =>
								updateField(fieldName as keyof Middleware, e.currentTarget.value as never)}
						/>
					{:else if field.type === 'number'}
						<Input
							type="number"
							id={fieldName}
							value={middleware[fieldName as keyof Middleware] ?? 0}
							oninput={(e) =>
								updateField(fieldName as keyof Middleware, Number(e.currentTarget.value) as never)}
						/>
					{:else if field.type === 'boolean'}
						<Checkbox
							id={fieldName}
							checked={middleware[fieldName as keyof Middleware] ?? false}
							onCheckedChange={(checked) =>
								updateField(fieldName as keyof Middleware, checked as never)}
						/>
					{:else if field.type === 'array'}
						<Input
							type="text"
							id={fieldName}
							value={Array.isArray(middleware[fieldName as keyof Middleware])
								? (middleware[fieldName as keyof Middleware] as string[]).join(', ')
								: ''}
							oninput={(e) => handleArrayInput(fieldName, e.currentTarget.value)}
							placeholder="Enter comma-separated values"
						/>
					{/if}

					{#if field.description}
						<p class="text-sm text-muted-foreground">{field.description}</p>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
