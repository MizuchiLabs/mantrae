<script lang="ts" generics="TData, TValue">
	import {
		type ColumnDef,
		type PaginationState,
		type SortingState,
		type ColumnFiltersState,
		getCoreRowModel,
		getPaginationRowModel,
		getSortedRowModel,
		getFilteredRowModel,
		type VisibilityState,
		type RowSelectionState
	} from '@tanstack/table-core';
	import { rankItem } from '@tanstack/match-sorter-utils';
	import {
		createSvelteTable,
		FlexRender,
		renderComponent
	} from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { TraefikSource } from '$lib/types';
	import { source } from '$lib/stores/source';
	import { api, rdps } from '$lib/api';
	import {
		ArrowDown,
		ArrowUp,
		ChevronLeft,
		ChevronRight,
		ChevronsLeft,
		ChevronsRight,
		Delete,
		Plus,
		Search,
		X
	} from '@lucide/svelte';
	import { limit } from '$lib/stores/common';
	import BulkActions from './BulkActions.svelte';
	import type { BulkAction } from './types';

	type DataTableProps<TData, TValue> = {
		columns: ColumnDef<TData, TValue>[];
		data: TData[];
		createButton?: {
			label: string;
			onClick: () => void;
		};
		showSourceTabs?: boolean;
		onRowSelection?: (selectedRows: TData[]) => void;
		getRowClassName?: (row: TData) => string;
		bulkActions?: BulkAction<TData>[] | undefined;
	};

	let {
		data,
		columns,
		createButton,
		showSourceTabs,
		getRowClassName,
		bulkActions
	}: DataTableProps<TData, TValue> = $props();

	// Pagination
	const pageSizeOptions = [10, 20, 30, 40, 50];
	let pagination = $state<PaginationState>({
		pageIndex: 0,
		pageSize: parseInt(limit.value ?? pageSizeOptions[0].toString())
	});
	let sorting = $state<SortingState>([]);
	let columnFilters = $derived<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let rowSelection = $state<RowSelectionState>({});
	let globalFilter = $state<string>('');

	// Table
	const table = createSvelteTable({
		get data() {
			return data;
		},
		columns: [
			{
				id: 'select',
				header: ({ table }) =>
					renderComponent(Checkbox, {
						checked: table.getIsAllPageRowsSelected(),
						indeterminate: table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected(),
						onCheckedChange: (value) => table.toggleAllPageRowsSelected(!!value),
						'aria-label': 'Select all'
					}),
				cell: ({ row }) =>
					renderComponent(Checkbox, {
						checked: row.getIsSelected(),
						onCheckedChange: (value) => row.toggleSelected(!!value),
						'aria-label': 'Select row'
					}),
				enableSorting: false,
				enableHiding: false
			},
			...columns
		],
		autoResetAll: true,
		filterFns: {
			fuzzy: (row, columnId, value, addMeta) => {
				const itemRank = rankItem(row.getValue(columnId), value);
				addMeta({ itemRank });
				return itemRank.passed;
			}
		},
		globalFilterFn: (row, columnId, value, addMeta) => {
			const itemRank = rankItem(row.getValue(columnId), value);
			addMeta({ itemRank });
			return itemRank.passed;
		},
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				pagination = updater(pagination);
			} else {
				pagination = updater;
			}
		},
		onSortingChange: (updater) => {
			if (typeof updater === 'function') {
				sorting = updater(sorting);
			} else {
				sorting = updater;
			}
		},
		onColumnFiltersChange: (updater) => {
			if (typeof updater === 'function') {
				columnFilters = updater(columnFilters);
			} else {
				columnFilters = updater;
			}
		},
		onColumnVisibilityChange: (updater) => {
			if (typeof updater === 'function') {
				columnVisibility = updater(columnVisibility);
			} else {
				columnVisibility = updater;
			}
		},
		onRowSelectionChange: (updater) => {
			if (typeof updater === 'function') {
				rowSelection = updater(rowSelection);
			} else {
				rowSelection = updater;
			}
		},
		onGlobalFilterChange: (updater) => {
			if (typeof updater === 'function') {
				globalFilter = updater(globalFilter);
			} else {
				globalFilter = updater;
			}
		},
		state: {
			get pagination() {
				return pagination;
			},
			get sorting() {
				return sorting;
			},
			get columnFilters() {
				return columnFilters;
			},
			get columnVisibility() {
				return columnVisibility;
			},
			get rowSelection() {
				return rowSelection;
			},
			get globalFilter() {
				return globalFilter;
			}
		}
	});

	// Update localStorage and fetch config when tab changes
	async function handleTabChange(value: string) {
		if (!source.isValid(value)) return;
		source.value = value;
		// Reset table state
		table.resetRowSelection();
		table.resetColumnFilters();
		table.resetGlobalFilter();
		table.resetColumnOrder();
		table.resetPagination();
		await Promise.all([api.getTraefikConfig(source.value), api.listDNSProviders()]);
	}
	function handleLimitChange(value: string) {
		if (!value) return;
		table.setPageSize(Number(value));
		pagination.pageSize = Number(value);
		limit.value = value;
	}
	function clearFilter(columnId: string) {
		const column = table.getColumn(columnId);
		if (column) column.setFilterValue(undefined);
	}
</script>

<div>
	<div class="flex items-center justify-between gap-2 py-4">
		<div class="relative flex items-center">
			<Search class="text-muted-foreground absolute left-3" size={16} />
			<Input
				placeholder="Search..."
				bind:value={globalFilter}
				oninput={() => table.setGlobalFilter(String(globalFilter))}
				class="w-[200px] pl-9 lg:w-[350px]"
			/>
			<Delete
				class="text-muted-foreground absolute right-4"
				size={16}
				onclick={() => table.setGlobalFilter('')}
			/>
		</div>

		<!-- Tabs -->
		{#if showSourceTabs}
			<Tabs.Root bind:value={source.value} onValueChange={handleTabChange}>
				<Tabs.List class="grid w-[400px] grid-cols-3">
					<Tabs.Trigger value={TraefikSource.LOCAL}>Local</Tabs.Trigger>
					<Tabs.Trigger value={TraefikSource.API}>API</Tabs.Trigger>
					<Tabs.Trigger value={TraefikSource.AGENT}>Agent</Tabs.Trigger>
				</Tabs.List>
			</Tabs.Root>
		{/if}

		{#if table.getState().columnFilters.length > 0}
			<Button onclick={() => table.setColumnFilters([])}>Clear Filters</Button>
			{#each table.getState().columnFilters as filter (filter.id)}
				<Badge
					variant="secondary"
					class="hover:bg-muted-foreground/20 flex items-center gap-1 hover:cursor-pointer"
					onclick={() => clearFilter(filter.id)}
				>
					<X size={12} />
					{filter.id.toLowerCase()}: {String(filter.value)}
				</Badge>
			{/each}
		{/if}

		<!-- Column Visibility -->
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button {...props} variant="outline" class="ml-auto">Columns</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end">
				{#each table.getAllColumns().filter((col) => col.getCanHide()) as column (column.id)}
					<DropdownMenu.CheckboxItem
						class="capitalize"
						closeOnSelect={false}
						bind:checked={() => column.getIsVisible(), (v) => column.toggleVisibility(!!v)}
					>
						{column.columnDef.header}
					</DropdownMenu.CheckboxItem>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>

		{#if createButton}
			<Button variant="default" onclick={createButton.onClick}>
				<Plus />
				{createButton.label}
			</Button>
		{/if}
	</div>

	<!-- Table -->
	<div class="rounded-md border">
		{#key source.value + $rdps + JSON.stringify(data)}
			<Table.Root>
				<Table.Header>
					{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
						<Table.Row>
							{#each headerGroup.headers as header (header.id)}
								<Table.Head>
									{#if !header.isPlaceholder}
										<div class="flex items-center">
											<Button
												variant="ghost"
												size="sm"
												class="-ml-3 h-8 data-[sortable=false]:cursor-default"
												data-sortable={header.column.getCanSort()}
												onclick={() => header.column.toggleSorting()}
											>
												<FlexRender
													content={header.column.columnDef.header}
													context={header.getContext()}
												/>
												{#if header.column.getCanSort()}
													{#if header.column.getIsSorted() === 'asc'}
														<ArrowDown />
													{:else if header.column.getIsSorted() === 'desc'}
														<ArrowUp />
													{/if}
												{/if}
											</Button>
										</div>
									{/if}
								</Table.Head>
							{/each}
						</Table.Row>
					{/each}
				</Table.Header>
				<Table.Body>
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row
							data-state={row.getIsSelected() && 'selected'}
							class={getRowClassName ? getRowClassName(row.original) : ''}
						>
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell>
									<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
								</Table.Cell>
							{/each}
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={columns.length} class="h-24 text-center">No results.</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
				<Table.Footer>
					<Table.Row class="border-t">
						<Table.Cell colspan={columns.length}>Total</Table.Cell>
						<Table.Cell class="mr-4 text-right"
							>{table.getPrePaginationRowModel().rows.length}</Table.Cell
						>
					</Table.Row>
				</Table.Footer>
			</Table.Root>
		{/key}
	</div>
	{#if table.getSelectedRowModel().rows.length > 0 && bulkActions && bulkActions.length > 0}
		<BulkActions
			selectedCount={table.getFilteredSelectedRowModel().rows.length}
			totalCount={table.getFilteredRowModel().rows.length}
			actions={bulkActions}
			selectedItems={table.getSelectedRowModel().rows.map((row) => row.original)}
		/>
	{/if}

	<!-- Pagination -->
	<div class="flex items-center justify-between py-4">
		<div>
			<Select.Root
				type="single"
				allowDeselect={false}
				value={limit.value ?? pagination.pageSize.toString()}
				onValueChange={handleLimitChange}
			>
				<Select.Trigger class="w-[180px]">
					{pagination.pageSize}
				</Select.Trigger>
				<Select.Content>
					{#each pageSizeOptions as size (size)}
						<Select.Item value={size.toString()}>{size}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
		<div class="flex items-center justify-end gap-2">
			<Button
				variant="outline"
				size="icon"
				onclick={() => table.firstPage()}
				disabled={!table.getCanPreviousPage()}
			>
				<ChevronsLeft />
			</Button>
			<Button
				variant="outline"
				size="icon"
				onclick={() => table.previousPage()}
				disabled={!table.getCanPreviousPage()}
			>
				<ChevronLeft />
			</Button>
			<span class="text-muted-foreground text-sm">
				Page {pagination.pageIndex + 1} / {table.getPageCount()}
			</span>
			<Button
				variant="outline"
				size="icon"
				onclick={() => table.nextPage()}
				disabled={!table.getCanNextPage()}
			>
				<ChevronRight />
			</Button>
			<Button
				variant="outline"
				size="icon"
				onclick={() => table.lastPage()}
				disabled={!table.getCanNextPage()}
			>
				<ChevronsRight />
			</Button>
		</div>
	</div>
</div>
