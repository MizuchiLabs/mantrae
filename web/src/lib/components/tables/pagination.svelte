<script lang="ts">
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import ChevronRight from 'lucide-svelte/icons/chevron-right';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import * as Select from '$lib/components/ui/select';
	import type { Selected } from 'bits-ui';

	export let count: number;
	export let perPage: Selected<number> | undefined;
	export let currentPage: number;
	const limits = [5, 10, 25, 50, 100];
	const changeLimit = (limit: Selected<number> | undefined) => {
		if (limit === undefined) return;
		perPage = limit;
		localStorage.setItem('limit', JSON.stringify(limit));
	};
</script>

<div class="flex flex-row justify-between">
	<Select.Root selected={perPage} onSelectedChange={changeLimit}>
		<Select.Trigger class="w-[180px]">
			<Select.Value placeholder="Select a limit" />
		</Select.Trigger>
		<Select.Content>
			{#each limits as limit}
				<Select.Item value={limit} label={limit.toString()}>{limit}</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
	<div class="flex">
		<Pagination.Root
			{count}
			page={currentPage}
			perPage={perPage?.value ?? 10}
			let:pages
			let:currentPage
			onPageChange={(page) => (currentPage = page)}
		>
			<Pagination.Content>
				<Pagination.Item>
					<Pagination.PrevButton>
						<ChevronLeft class="h-4 w-4" />
						<span class="hidden sm:block">Previous</span>
					</Pagination.PrevButton>
				</Pagination.Item>
				{#each pages as page (page.key)}
					{#if page.type === 'ellipsis'}
						<Pagination.Item>
							<Pagination.Ellipsis />
						</Pagination.Item>
					{:else}
						<Pagination.Item>
							<Pagination.Link {page} isActive={currentPage === page.value}>
								{page.value}
							</Pagination.Link>
						</Pagination.Item>
					{/if}
				{/each}
				<Pagination.Item>
					<Pagination.NextButton>
						<span class="hidden sm:block">Next</span>
						<ChevronRight class="h-4 w-4" />
					</Pagination.NextButton>
				</Pagination.Item>
			</Pagination.Content>
		</Pagination.Root>
	</div>
</div>
