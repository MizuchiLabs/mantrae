<script lang="ts">
	import * as Select from '$lib/components/ui/select';
	import Input from '$lib/components/ui/input/input.svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import { Delete } from 'lucide-svelte';
	import type { Selected } from 'bits-ui';
	import { onMount } from 'svelte';
	import { LOCAL_PROVIDER_SK } from '$lib/store';

	export let search: string;
	export let columns: Selected<string>[];
	export let columnName: string;
	export let fColumns: string[];

	const clearSearch = () => {
		search = '';
		localProvider = false;
		localStorage.setItem(LOCAL_PROVIDER_SK, localProvider.toString());
	};

	// Only show local routers not external ones
	let localProvider = localStorage.getItem(LOCAL_PROVIDER_SK) === 'true';
	const toggleProvider = () => {
		localProvider = !localProvider;
		search = localProvider ? '@provider:http' : '';
		localStorage.setItem(LOCAL_PROVIDER_SK, localProvider.toString());
	};

	const changeColumns = (columns: Selected<string>[] | undefined) => {
		if (columns === undefined) return;
		fColumns = columns.map((c) => c.value);
		localStorage.setItem(columnName, JSON.stringify(fColumns));
	};

	onMount(() => {
		search = localProvider ? '@provider:http' : '';
	});
</script>

<div class="flex flex-row items-center justify-between">
	<div class="flex flex-row items-center gap-1">
		<div class="flex flex-row items-center justify-end gap-1">
			<div class="absolute flex flex-row items-center justify-between gap-1">
				<Button
					variant="ghost"
					class="mr-1 rounded-full hover:bg-transparent"
					on:click={clearSearch}
					size="icon"
					aria-hidden
				>
					<Delete size="1.25rem" class="text-muted-foreground hover:text-red-400" />
				</Button>
			</div>
			<Input
				type="text"
				placeholder="Search..."
				class="w-80 focus-visible:ring-0 focus-visible:ring-offset-0"
				bind:value={search}
			/>
		</div>
		<button
			class={buttonVariants({ variant: 'outline' })}
			class:bg-primary={localProvider}
			class:text-primary-foreground={localProvider}
			on:click={toggleProvider}
		>
			Local Only
		</button>
	</div>

	<Select.Root
		multiple
		selected={fColumns.map((c) => ({ value: c, label: c }))}
		onSelectedChange={changeColumns}
	>
		<Select.Trigger class="w-[180px]">
			<Select.Value placeholder="Columns" />
		</Select.Trigger>
		<Select.Content>
			{#each columns as column}
				<Select.Item value={column.value} label={column.label}>
					{column.label}
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
</div>
