<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Button } from '$lib/components/ui/button';
	import { Textarea } from '$lib/components/ui/textarea';
	import { extractSchemaFields, type FormField } from '$lib/formGenerator';
	import type { ZodSchema } from 'zod';
	import DynamicForm from './DynamicForm.svelte';
	import YAML from 'yaml';
	import { Plus, Trash } from '@lucide/svelte';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';

	interface Props {
		schema: ZodSchema;
		data: Record<string, unknown>;
		onUpdate: (data: Record<string, unknown>) => void;
	}

	let { schema, data, onUpdate }: Props = $props();

	let fields = $derived(extractSchemaFields(schema));

	// Make data reactive
	let formData = $derived(data);

	function updateField(key: string, value: unknown) {
		formData[key] = value;
		onUpdate({ ...formData });
	}

	function addArrayItem(key: string) {
		if (!formData[key]) formData[key] = [];
		(formData[key] as unknown[]).push(getDefaultValue(fields[key].arrayItemType!));
		onUpdate({ ...formData });
	}

	function removeArrayItem(key: string, index: number) {
		(formData[key] as unknown[]).splice(index, 1);
		onUpdate({ ...formData });
	}

	function addRecordItem(key: string) {
		if (!formData[key]) formData[key] = {};
		const newKey = `key${Object.keys(formData[key] as object).length + 1}`;
		(formData[key] as Record<string, unknown>)[newKey] = getDefaultValue(
			fields[key].recordValueType!
		);
		onUpdate({ ...formData });
	}

	function removeRecordItem(key: string, recordKey: string) {
		delete (formData[key] as Record<string, unknown>)[recordKey];
		onUpdate({ ...formData });
	}

	function updateRecordKey(key: string, oldKey: string, newKey: string) {
		const record = formData[key] as Record<string, unknown>;
		if (oldKey !== newKey && !record[newKey]) {
			record[newKey] = record[oldKey];
			delete record[oldKey];
			onUpdate({ ...formData });
		}
	}

	function updateRecordValue(key: string, recordKey: string, value: unknown) {
		(formData[key] as Record<string, unknown>)[recordKey] = value;
		onUpdate({ ...formData });
	}

	let yamlError = $state<string | null>(null);
	function handlePluginChange(value: string) {
		try {
			const parsed = YAML.parse(value);
			if (parsed) {
				yamlError = null;
				formData = parsed;
				onUpdate({ ...formData });
			}
		} catch (e) {
			yamlError = e instanceof Error ? e.message : 'Invalid YAML';
		}
	}

	function getPluginValue(): string {
		const value = formData;
		if (typeof value === 'string') return value;
		if (typeof value === 'object') return YAML.stringify(value, { indent: 2 });
		return '';
	}

	function getDefaultValue(field: FormField): unknown {
		switch (field.type) {
			case 'string':
				return '';
			case 'number':
				return 0;
			case 'boolean':
				return false;
			case 'array':
				return [];
			case 'object':
				return {};
			case 'record':
				return {};
			case 'plugin':
				return {};
			default:
				return '';
		}
	}
	function showDescription(field: FormField): boolean {
		switch (field.type) {
			case 'string':
				return false;
			case 'number':
				return false;
			case 'boolean':
				return true;
			case 'array':
				return true;
			case 'object':
				return true;
			case 'record':
				return true;
			case 'plugin':
				return true;
			default:
				return false;
		}
	}
	function handleYamlIndents(e: KeyboardEvent) {
		const textarea = e.target as HTMLTextAreaElement;
		if (e.key === 'Tab') {
			e.preventDefault();
			const start = textarea.selectionStart;
			const end = textarea.selectionEnd;

			textarea.setRangeText('  ', start, end, 'end');
		}
		if (e.key === 'Backspace') {
			const start = textarea.selectionStart;
			const end = textarea.selectionEnd;

			// Only act if nothing selected and cursor is after two spaces
			if (start === end && textarea.value.slice(start - 2, start) === '  ') {
				e.preventDefault();
				textarea.setRangeText('', start - 2, start, 'end');
			}
		}
	}
</script>

<div class="flex flex-col gap-4">
	{#each Object.entries(fields) as [key, field] (key)}
		<div
			class={`flex gap-2 ${field.type === 'boolean' ? 'flex-row items-center  justify-between rounded-lg border p-3' : 'flex-col'}`}
		>
			<!-- <div class="flex items-center justify-between rounded-lg border p-3"> -->
			<div class="space-y-1">
				<Label class="text-sm font-medium">
					{field.label}
					{#if !field.optional}
						<span class="text-red-500">*</span>
					{/if}
				</Label>
				{#if field.description && showDescription(field)}
					<p class="text-muted-foreground text-xs">{field.description}</p>
				{/if}
			</div>

			{#if field.type === 'string'}
				<Input
					bind:value={data[key]}
					placeholder={field.description}
					oninput={() => updateField(key, data[key])}
				/>
			{:else if field.type === 'number'}
				<Input
					type="number"
					bind:value={data[key]}
					placeholder={field.description}
					oninput={() => updateField(key, data[key])}
				/>
			{:else if field.type === 'boolean'}
				<CustomSwitch
					checked={data[key] as boolean}
					onCheckedChange={(checked) => updateField(key, checked)}
				/>
			{:else if field.type === 'plugin'}
				<Textarea
					value={getPluginValue()}
					placeholder="Edit plugin configuration as YAML"
					class="min-h-[100px] font-mono text-sm"
					oninput={(e) => handlePluginChange(e.currentTarget.value)}
					onkeydown={handleYamlIndents}
				/>
				{#if yamlError}
					<p class="text-xs text-red-400 dark:text-red-700">{yamlError}</p>
				{/if}
			{:else if field.type === 'array'}
				<div class="flex flex-col gap-2 rounded-md border p-3">
					{#if data[key] && Array.isArray(data[key])}
						{#each data[key] as _, index (index)}
							<div class="flex items-center gap-2">
								{#if field.arrayItemType?.type === 'string'}
									<Input
										bind:value={data[key][index]}
										oninput={() => updateField(key, data[key])}
									/>
								{:else if field.arrayItemType?.type === 'number'}
									<Input
										type="number"
										bind:value={data[key][index]}
										oninput={() => updateField(key, data[key])}
									/>
								{:else if field.arrayItemType?.type === 'object' && field.arrayItemType.nestedSchema}
									<div class="flex-1 rounded border p-2">
										<DynamicForm
											schema={field.arrayItemType.nestedSchema}
											data={(data[key][index] as Record<string, unknown>) || {}}
											onUpdate={(nestedData) => {
												(data[key] as unknown[])[index] = nestedData;
												updateField(key, data[key]);
											}}
										/>
									</div>
								{/if}
								<Button
									variant="ghost"
									size="icon"
									class="text-red-500"
									onclick={() => removeArrayItem(key, index)}
								>
									<Trash />
								</Button>
							</div>
						{/each}
					{/if}
					<Button variant="outline" size="sm" onclick={() => addArrayItem(key)}>
						<Plus />
						Add {field.label}
					</Button>
				</div>
			{:else if field.type === 'record'}
				<div class="flex flex-col gap-2 rounded-md border p-3">
					{#if formData[key] && typeof formData[key] === 'object'}
						{#each Object.entries(formData[key] as Record<string, unknown>) as [recordKey, recordValue] (recordKey)}
							<div class="flex items-center gap-2">
								<Input
									value={recordKey}
									placeholder="Key"
									oninput={(e) => updateRecordKey(key, recordKey, e.currentTarget.value)}
								/>
								{#if field.recordValueType?.type === 'string'}
									<Input
										value={recordValue as string}
										placeholder="Value"
										oninput={(e) => updateRecordValue(key, recordKey, e.currentTarget.value)}
									/>
								{:else if field.recordValueType?.type === 'number'}
									<Input
										type="number"
										value={recordValue as number}
										placeholder="Value"
										oninput={(e) =>
											updateRecordValue(key, recordKey, parseFloat(e.currentTarget.value) || 0)}
									/>
								{/if}
								<Button
									variant="outline"
									size="sm"
									onclick={() => removeRecordItem(key, recordKey)}
								>
									Remove
								</Button>
							</div>
						{/each}
					{/if}
					<Button variant="outline" size="sm" onclick={() => addRecordItem(key)}>
						Add {field.label} Entry
					</Button>
				</div>
			{:else if field.type === 'object' && field.nestedSchema}
				<div class="rounded-md border p-3">
					<DynamicForm
						schema={field.nestedSchema}
						data={(formData[key] as Record<string, unknown>) || {}}
						onUpdate={(nestedData) => updateField(key, nestedData)}
					/>
				</div>
			{/if}
		</div>
	{/each}
</div>
