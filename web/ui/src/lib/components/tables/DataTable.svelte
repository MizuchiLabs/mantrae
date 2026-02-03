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
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
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
	import BulkActions from './BulkActions.svelte';
	import type { BulkAction } from './types';
	import { pageIndex, pageSize } from '$lib/store.svelte';

	type DataTableProps<TData, TValue> = {
		data?: TData[];
		columns: ColumnDef<TData, TValue>[];
		onPaginationChange?: (pagination: PaginationState) => void;
		onSortingChange?: (sorting: SortingState) => void;
		onRowSelection?: (rowSelection: RowSelectionState) => void;
		getRowClassName?: (row: TData) => string;
		rowClassModifiers?: Record<string, (row: TData) => boolean>;
		bulkActions?: BulkAction<TData>[] | undefined;
		createButton?: {
			label: string;
			onClick: () => void;
		};
	};

	let {
		data = [],
		columns,
		onPaginationChange,
		onSortingChange,
		onRowSelection,
		getRowClassName,
		rowClassModifiers,
		bulkActions,
		createButton
	}: DataTableProps<TData, TValue> = $props();

	// Pagination
	const pageSizeOptions = [5, 10, 25, 50, 100];
	let pagination = $state<PaginationState>({
		pageIndex: 0,
		pageSize: pageSize.current ?? 10
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
		get columns() {
			return [
				{
					id: 'select',
					header: ({ table }: { table: any }) =>
						renderComponent(Checkbox, {
							checked: table.getIsAllPageRowsSelected(),
							indeterminate: table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected(),
							onCheckedChange: (value) => table.toggleAllPageRowsSelected(!!value),
							'aria-label': 'Select all'
						}),
					cell: ({ row }: { row: any }) =>
						renderComponent(Checkbox, {
							checked: row.getIsSelected(),
							onCheckedChange: (value) => row.toggleSelected(!!value),
							'aria-label': 'Select row'
						}),
					enableSorting: false,
					enableHiding: false,
					enableGlobalFilter: false
				},
				...columns
			];
		},
		filterFns: {
			fuzzy: (row, columnId, value, addMeta) => {
				const itemRank = rankItem(row.getValue(columnId), value);
				addMeta({ itemRank });
				return itemRank.passed;
			},
			arrIncludes: (row, columnId, value) => {
				const cellValue = row.getValue(columnId) as string[];
				if (!Array.isArray(cellValue)) return false;
				return cellValue.some((item) => item.toLowerCase().includes(value.toLowerCase()));
			}
		},
		globalFilterFn: 'includesString',
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
			if (onPaginationChange) onPaginationChange(pagination);
			pageIndex.current = pagination.pageIndex;
			pageSize.current = pagination.pageSize;
		},
		onSortingChange: (updater) => {
			if (typeof updater === 'function') {
				sorting = updater(sorting);
			} else {
				sorting = updater;
			}
			if (onSortingChange) onSortingChange(sorting);
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
			if (onRowSelection) onRowSelection(rowSelection);
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

	// helper to merge all classes into one string
	function computeRowClasses(row: TData) {
		const classes: string[] = [];

		if (getRowClassName) {
			const c = getRowClassName(row);
			if (c) classes.push(c);
		}

		if (rowClassModifiers) {
			for (const [cls, fn] of Object.entries(rowClassModifiers)) {
				if (fn(row)) classes.push(cls);
			}
		}
		return classes.join(' ');
	}

	function clearFilter(columnId: string) {
		const column = table.getColumn(columnId);
		if (column) column.setFilterValue(undefined);
	}
</script>

<div>
	<div class="flex flex-col gap-2 py-4 sm:flex-row sm:items-center sm:justify-between">
		<div class="relative flex items-center">
			<Search class="absolute left-3 text-muted-foreground" size={16} />
			<Input
				placeholder="Search..."
				bind:value={globalFilter}
				oninput={() => table.setGlobalFilter(String(globalFilter))}
				class="w-full pl-9 sm:w-[180px] lg:w-[350px]"
			/>
			<Delete
				class="absolute right-4 text-muted-foreground"
				size={16}
				onclick={() => table.setGlobalFilter('')}
			/>
		</div>

		{#if table.getState().columnFilters.length > 0}
			<Button onclick={() => table.setColumnFilters([])}>Clear Filters</Button>
			{#each table.getState().columnFilters as filter (filter.id)}
				<Badge
					variant="secondary"
					class="flex items-center gap-1 hover:cursor-pointer hover:bg-muted-foreground/20"
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
		{#key table.getRowModel().rows}
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
							class={computeRowClasses(row.original)}
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
						<Table.Cell class="mr-4 text-right">
							{table.getPaginationRowModel().rows.length}
						</Table.Cell>
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
	<div class="flex flex-col gap-4 py-4 sm:flex-row sm:items-center sm:justify-between">
		<!-- Page size selector -->
		<div class="flex justify-center sm:justify-start">
			<Select.Root
				type="single"
				allowDeselect={false}
				value={pagination.pageSize.toString()}
				onValueChange={(value) => table.setPageSize(Number(value))}
			>
				<Select.Trigger class="w-full sm:w-[180px]">
					{pagination.pageSize}
				</Select.Trigger>
				<Select.Content>
					{#each pageSizeOptions as size (size)}
						<Select.Item value={size.toString()}>{size}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>

		<!-- Pagination controls -->
		<div class="flex flex-wrap items-center justify-center gap-2 text-sm sm:justify-end">
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
			<span class="text-sm text-muted-foreground">
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
