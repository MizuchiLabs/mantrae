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
		Check,
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
	import { pageIndex, pageSize } from '$lib/stores/common';

	type DataTableProps<TData, TValue> = {
		data: TData[];
		columns: ColumnDef<TData, TValue>[];
		rowCount: number;
		viewMode?: 'table' | 'grid';
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
		cardConfig?: {
			titleKey?: string;
			subtitleKey?: string;
			excludeColumns?: string[];
		};
	};

	let {
		data,
		columns,
		rowCount,
		viewMode = 'table',
		onPaginationChange,
		onSortingChange,
		onRowSelection,
		getRowClassName,
		rowClassModifiers,
		bulkActions,
		createButton,
		cardConfig = {}
	}: DataTableProps<TData, TValue> = $props();

	// Pagination
	const pageSizeOptions = [5, 10, 25, 50, 100];
	let pagination = $state<PaginationState>({
		pageIndex: pageIndex.value ?? 0,
		pageSize: pageSize.value ?? 10
	});
	let pageCount = $derived(Math.ceil(rowCount / pagination.pageSize));
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
				enableHiding: false,
				enableGlobalFilter: false
			},
			...columns
		],
		manualPagination: true,
		get rowCount() {
			return rowCount;
		},
		get pageCount() {
			return pageCount;
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
			pageIndex.value = pagination.pageIndex;
			pageSize.value = pagination.pageSize;
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
	function getVisibleColumns() {
		return table
			.getAllColumns()
			.filter(
				(col) =>
					col.getIsVisible() &&
					col.id !== 'select' &&
					col.id !== 'actions' &&
					!cardConfig.excludeColumns?.includes(col.id)
			);
	}
</script>

<div>
	<div class="flex flex-col gap-2 py-4 sm:flex-row sm:items-center sm:justify-between">
		<div class="relative flex items-center">
			<Search class="text-muted-foreground absolute left-3" size={16} />
			<Input
				placeholder="Search..."
				bind:value={globalFilter}
				oninput={() => table.setGlobalFilter(String(globalFilter))}
				class="w-full pl-9 sm:w-[180px] lg:w-[350px]"
			/>
			<Delete
				class="text-muted-foreground absolute right-4"
				size={16}
				onclick={() => table.setGlobalFilter('')}
			/>
		</div>

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
		{#if viewMode === 'table'}
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
		{/if}

		{#if createButton}
			<Button variant="default" onclick={createButton.onClick}>
				<Plus />
				{createButton.label}
			</Button>
		{/if}
	</div>

	{#if viewMode === 'table'}
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
								<Table.Cell colspan={columns.length} class="h-24 text-center"
									>No results.</Table.Cell
								>
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
	{:else}
		<!-- Grid View -->
		<div
			class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
		>
			{#each table.getRowModel().rows as row (row.id)}
				{@const isSelected = row.getIsSelected()}
				{@const titleCell = cardConfig.titleKey
					? row.getVisibleCells().find((c) => c.column.id === cardConfig.titleKey)
					: null}
				{@const subtitleCell = cardConfig.subtitleKey
					? row.getVisibleCells().find((c) => c.column.id === cardConfig.subtitleKey)
					: null}
				{@const actionsCell = row.getVisibleCells().find((c) => c.column.id === 'actions')}
				{@const visibleColumns = getVisibleColumns()}

				<!-- svelte-ignore a11y_click_events_have_key_events -->
				<!-- svelte-ignore a11y_no_static_element_interactions -->
				<div
					class="group bg-card text-card-foreground hover:shadow-primary/5 relative overflow-hidden rounded-xl border shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md {computeRowClasses(
						row.original
					)} {isSelected
						? 'ring-primary bg-primary/5 border-primary/20 ring-2'
						: 'hover:border-primary/20'}"
					onclick={() => row.toggleSelected()}
					ontouchstart={(e) => {
						let touchTimer;
						touchTimer = setTimeout(() => {
							// Long press action - could show context menu
							if (actionsCell) {
								e.preventDefault();
								// You could dispatch a custom event here for showing actions menu
							}
						}, 500);

						// Clear timer on touch end
						const clearTimer = () => {
							clearTimeout(touchTimer);
							document.removeEventListener('touchend', clearTimer);
							document.removeEventListener('touchmove', clearTimer);
						};

						document.addEventListener('touchend', clearTimer);
						document.addEventListener('touchmove', clearTimer);
					}}
				>
					{#if isSelected}
						<Check class="absolute bottom-3 left-3 z-10 text-green-500 " />
					{/if}

					<div class="space-y-3 p-4">
						<!-- Header Section -->
						<div class="flex items-center justify-between space-y-2">
							{#if titleCell && cardConfig.titleKey}
								{@const titleColumn = table.getColumn(cardConfig.titleKey)}
								<h3 class="line-clamp-2 pr-8 text-base leading-tight font-semibold">
									<FlexRender
										content={titleColumn?.columnDef.cell}
										context={titleCell.getContext()}
									/>
								</h3>
							{/if}

							{#if subtitleCell && cardConfig.subtitleKey}
								{@const subtitleColumn = table.getColumn(cardConfig.subtitleKey)}
								<div class="text-muted-foreground text-sm">
									<FlexRender
										content={subtitleColumn?.columnDef.cell}
										context={subtitleCell.getContext()}
									/>
								</div>
							{/if}
						</div>

						<!-- Content Section - Show only most important fields -->
						{#if visibleColumns.length > 0}
							<div class="space-y-2.5 border-t pt-3">
								{#each visibleColumns as column (column.id)}
									{@const cell = row.getVisibleCells().find((c) => c.column.id === column.id)}
									{#if cell && column.id !== cardConfig.titleKey && column.id !== cardConfig.subtitleKey}
										<div class="flex items-center justify-between gap-2 text-xs">
											<span
												class="text-muted-foreground min-w-0 truncate font-medium tracking-wide uppercase"
											>
												{column.columnDef.header}
											</span>
											<div class="min-w-0 flex-1 text-right">
												<FlexRender
													content={cell.column.columnDef.cell}
													context={cell.getContext()}
												/>
											</div>
										</div>
									{/if}
								{/each}
							</div>
						{/if}

						<!-- Actions Section -->
						{#if actionsCell}
							<div
								class="flex justify-end border-t pt-3 opacity-0 transition-opacity duration-200 group-hover:opacity-100 sm:opacity-100"
							>
								<div class="flex gap-1">
									<FlexRender
										content={actionsCell.column.columnDef.cell}
										context={actionsCell.getContext()}
									/>
								</div>
							</div>
						{/if}
					</div>

					<!-- Hover Overlay -->
					<div
						class="from-primary/5 pointer-events-none absolute inset-0 bg-gradient-to-t to-transparent opacity-0 transition-opacity duration-200 group-hover:opacity-100"
					></div>
				</div>
			{:else}
				<div class="col-span-full flex h-32 items-center justify-center">
					<div class="text-center space-y-2">
						<div class="text-muted-foreground text-4xl">📋</div>
						<p class="text-muted-foreground font-medium">No results found</p>
						<p class="text-muted-foreground text-sm">Try adjusting your search or filters</p>
					</div>
				</div>
			{/each}
		</div>
	{/if}

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
