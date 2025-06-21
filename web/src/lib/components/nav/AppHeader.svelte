<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
</script>

<header class="bg-primary-foreground flex h-16 items-center justify-between border-b px-4">
	<Sidebar.Trigger />
	{#if page.url.pathname.split('/')[1] === 'middlewares'}
		<Tabs.Root
			class="flex flex-col gap-2"
			value={page.url.pathname}
			onValueChange={(value) => {
				if (value === 'middlewares') {
					goto('/middlewares/');
				} else {
					goto('/middlewares/plugins');
				}
			}}
		>
			<Tabs.List>
				<Tabs.Trigger value="middlewares" class="px-2 py-0.5 font-bold">Middlewares</Tabs.Trigger>
				<Tabs.Trigger value="plugins" class="px-2 py-0.5 font-bold">Plugins</Tabs.Trigger>
			</Tabs.List>
		</Tabs.Root>
	{/if}
</header>
