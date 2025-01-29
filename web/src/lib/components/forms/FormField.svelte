<script lang="ts">
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Plus, Trash } from 'lucide-svelte';
	import { mwNames } from '$lib/api';

	interface Props {
		key: string;
		path: string;
		type: string;
		data: Record<string, unknown>;
		disabled?: boolean;
	}

	let { key, path, type, data = $bindable(), disabled }: Props = $props();

	$effect(() => {
		console.log(key, path, type, data);
	});

	type FormValue = string | number | boolean | string[] | Record<string, unknown>;

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

	function handleObjectChange(e: Event, key: string) {
		const target = e.target as HTMLInputElement;
		const currentValue = getNestedValue(data, path) as Record<string, unknown>;
		const newValue = { ...currentValue, [key]: target.value };
		setNestedValue(data, path, newValue);
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
</script>

<div class="grid gap-2">
	<Label for={path}>{formatLabel(key)}</Label>

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
						{fieldValue?.length > 0 ? fieldValue?.join(', ') : 'Select Middlewares'}
					</Select.Trigger>
					<Select.Content>
						{#each $mwNames as name}
							<Select.Item value={name}>{name}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			{/if}
		</div>
	{:else if type === 'object'}
		<div class="grid gap-2">
			{#each Object.entries(fieldValue || {}) as [key, value]}
				<div class="flex flex-col gap-1">
					<Label>{formatLabel(key)}</Label>
					<Input
						type="text"
						value={value as string}
						onchange={(e) => handleObjectChange(e, key)}
						{disabled}
					/>
				</div>
			{/each}
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
			{#each (fieldValue as string[]) || [] as value, i}
				<div class="flex gap-2">
					<Input
						type="text"
						{value}
						onchange={(e) => handleArrayChange(i, (e.target as HTMLInputElement).value)}
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
			{#if !disabled}
				<Button type="button" variant="outline" onclick={addArrayItem} class="w-full">
					<Plus />
					Add {key.charAt(0).toUpperCase() + key.slice(1)}
				</Button>
			{/if}
		</div>
	{:else if type === 'number'}
		<Input
			type="number"
			id={path}
			value={fieldValue as number}
			onchange={handleChange}
			{disabled}
		/>
	{:else}
		<Input type="text" id={path} value={fieldValue as string} onchange={handleChange} {disabled} />
	{/if}
</div>
