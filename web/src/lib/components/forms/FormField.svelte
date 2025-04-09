<script lang="ts">
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Plus, Trash } from 'lucide-svelte';
	import { mwNames } from '$lib/api';
	import type { FieldMetadata } from '$lib/types/middlewares';
	import Separator from '../ui/separator/separator.svelte';

	interface Props {
		key: string;
		path: string;
		type: string;
		data: Record<string, unknown>;
		metadata?: FieldMetadata;
		disabled?: boolean;
	}

	let { key, path, type, data = $bindable(), metadata = {}, disabled }: Props = $props();

	type FormValue =
		| string
		| number
		| boolean
		| string[]
		| Record<string, unknown>
		| Record<string, unknown>[];

	function getNestedValue(obj: Record<string, unknown>, path: string): FormValue {
		return path.split('.').reduce<unknown>((acc, part) => {
			if (acc && typeof acc === 'object') {
				return (acc as Record<string, unknown>)[part];
			}
			return undefined;
		}, obj) as FormValue;
	}

	function setNestedValue(obj: Record<string, unknown>, path: string, value: FormValue) {
		const parts = path.split('.');
		const last = parts.pop()!;
		const target = parts.reduce((acc, part) => {
			if (!acc[part] || typeof acc[part] !== 'object') {
				acc[part] = {};
			}
			return acc[part] as Record<string, unknown>;
		}, obj);
		target[last] = value;
	}

	function handleChange(e: Event) {
		const target = e.target as HTMLInputElement;
		const newValue = type === 'number' ? Number(target.value) : target.value;
		setNestedValue(data, path, newValue);
	}

	function handleSwitchChange(checked: boolean) {
		setNestedValue(data, path, checked);
	}

	function handleArrayChange(index: number, value: string) {
		const array = (getNestedValue(data, path) as string[]) || [];
		array[index] = value;
		setNestedValue(data, path, array);
	}

	function addArrayItem() {
		const array = (getNestedValue(data, path) as string[]) || [];
		array.push('');
		setNestedValue(data, path, array);
	}

	function removeArrayItem(index: number) {
		const array = (getNestedValue(data, path) as string[]) || [];
		array.splice(index, 1);
		setNestedValue(data, path, array);
	}

	const fieldValue = $derived(getNestedValue(data, path));

	function formatLabel(str: string): string {
		return str
			.split(/(?=[A-Z])/)
			.join(' ')
			.replace(/^\w/, (c) => c.toUpperCase());
	}

	// Extra special cases
	const isChainMiddleware = $derived(path === 'middlewares');
	function handleMiddlewareChange(values: string[]) {
		setNestedValue(data, path, values);
	}

	function handleObjectArrayChange(index: number, field: string, value: string) {
		const array = (getNestedValue(data, path) as Record<string, unknown>[]) || [];
		if (!isObjectArray(array)) return;
		if (!array[index]) {
			array[index] = {};
		}
		array[index] = { ...array[index], [field]: value };
		setNestedValue(data, path, array);
	}

	function removeObjectArrayItem(index: number) {
		if (index < 1) return;
		const array = (getNestedValue(data, path) as string[]) || [];
		array.splice(index, 1);
		setNestedValue(data, path, array);
	}

	function addObjectArrayItem() {
		const array = (getNestedValue(data, path) as Record<string, unknown>[]) || [];
		// Use the first item as a template, creating an object with the same keys but empty values
		const template = array[0] || {};
		const newItem = Object.keys(template).reduce(
			(obj, key) => {
				obj[key] = '';
				return obj;
			},
			{} as Record<string, unknown>
		);
		array.push(newItem);
		setNestedValue(data, path, array);
	}

	// Helper to detect if array contains objects
	function isObjectArray(value: unknown[]): value is Record<string, unknown>[] {
		return (
			Array.isArray(value) && value.length > 0 && typeof value[0] === 'object' && value[0] !== null
		);
	}
</script>

<div class="grid gap-2">
	<Label for={path} class="flex flex-row items-center justify-between">
		{formatLabel(key)}
		{#if metadata.description}
			<span class="text-muted-foreground ml-1 text-sm">
				{metadata.description}
			</span>
		{/if}
	</Label>

	{#if isChainMiddleware}
		<div class="flex flex-col gap-2">
			{#if Array.isArray(fieldValue) && fieldValue.length > 0 && disabled}
				{#each fieldValue as middleware}
					<div class="flex items-center gap-2">
						<Input type="text" value={middleware} readonly {disabled} />
						{#if !disabled}
							<Button variant="destructive" size="icon" type="button">
								<Trash class="h-4 w-4" />
							</Button>
						{/if}
					</div>
				{/each}
			{/if}
			{#if !disabled}
				<Select.Root
					type="multiple"
					value={fieldValue as string[]}
					onValueChange={handleMiddlewareChange}
					{disabled}
				>
					<Select.Trigger>
						{Array.isArray(fieldValue) && fieldValue.length > 0
							? fieldValue.join(', ')
							: 'Select Middlewares'}
					</Select.Trigger>
					<Select.Content>
						{#each $mwNames as name}
							<Select.Item value={name}>{name}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			{/if}
		</div>
	{:else if type === 'boolean'}
		<Switch
			id={path}
			checked={fieldValue as boolean}
			onCheckedChange={handleSwitchChange}
			{disabled}
		/>
	{:else if type === 'array'}
		<div class="flex flex-col gap-2">
			{#if Array.isArray(fieldValue) && isObjectArray(fieldValue)}
				<div class="ml-4 flex flex-col gap-2 rounded border-l p-4">
					{#each fieldValue as item, i}
						<div class="flex flex-col gap-2 rounded">
							{#each Object.entries(item) as [field, value]}
								<div class="grid grid-cols-4 items-center gap-2">
									<Label class="col-span-1">{formatLabel(field)}</Label>
									<Input
										type="text"
										value={value as string}
										class="col-span-3"
										onchange={(e) =>
											handleObjectArrayChange(i, field, (e.target as HTMLInputElement).value)}
										placeholder={metadata.placeholder}
										{disabled}
									/>
								</div>
							{/each}
							<Separator class="my-2" />
							{#if !disabled}
								<Button
									variant="secondary"
									size="icon"
									type="button"
									class="w-full text-red-500"
									onclick={() => removeObjectArrayItem(i)}
								>
									<Trash />
								</Button>
							{/if}
						</div>
					{/each}
				</div>
			{:else}
				{#each (fieldValue as string[]) || [] as value, i}
					<div class="flex gap-2">
						<Input
							type="text"
							{value}
							onchange={(e) => handleArrayChange(i, (e.target as HTMLInputElement).value)}
							placeholder={metadata.placeholder}
							{disabled}
						/>
						{#if !disabled}
							<Button
								variant="ghost"
								size="icon"
								type="button"
								class="text-red-500"
								onclick={() => removeArrayItem(i)}
							>
								<Trash />
							</Button>
						{/if}
					</div>
				{/each}
			{/if}

			{#if !disabled}
				<Button
					type="button"
					variant="outline"
					onclick={Array.isArray(fieldValue) && isObjectArray(fieldValue)
						? addObjectArrayItem
						: addArrayItem}
					class="w-full"
				>
					<Plus />
					Add {key.charAt(0).toUpperCase() + key.slice(1)}
				</Button>
			{/if}
		</div>
	{:else if type === 'number'}
		<Input
			type="number"
			id={path}
			value={fieldValue !== undefined ? (fieldValue as number) : ''}
			onchange={handleChange}
			placeholder={metadata.placeholder}
			{disabled}
		/>
		{#if metadata.examples?.length}
			<div class="text-muted-foreground text-sm">
				Examples: {metadata.examples.join(', ')}
			</div>
		{/if}
	{:else}
		<Input
			type="text"
			id={path}
			value={fieldValue as string}
			placeholder={metadata.placeholder}
			onchange={handleChange}
			{disabled}
		/>
		{#if metadata.examples?.length}
			<div class="text-muted-foreground text-sm">
				Examples: {metadata.examples.join(', ')}
			</div>
		{/if}
	{/if}
</div>
