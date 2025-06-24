<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
</script>

<header
	class="flex h-16 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-16"
>
	<div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
		<Sidebar.Trigger class="-ml-1" />
		<Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />
		<h1 class="text-base font-medium">
			{page.url.pathname.split('/')[1].charAt(0).toUpperCase() +
				page.url.pathname.split('/')[1].slice(1)}
		</h1>

		{#if page.url.pathname.split('/')[1] === 'middlewares'}
			<div class="mr-4 ml-auto flex items-center gap-2">
				<Tabs.Root
					class="flex flex-col gap-2"
					value={page.url.pathname.split('/')[2] || 'middlewares'}
					onValueChange={(value) => {
						goto(`/middlewares/${value === 'middlewares' ? '' : value + '/'}`);
					}}
				>
					<Tabs.List>
						<Tabs.Trigger value="middlewares" class="px-2 py-0.5 font-bold"
							>Middlewares</Tabs.Trigger
						>
						<Tabs.Trigger value="plugins" class="px-2 py-0.5 font-bold">Plugins</Tabs.Trigger>
					</Tabs.List>
				</Tabs.Root>
			</div>
		{/if}
	</div>
</header>
