<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { SiteRoutes } from './routes';

	type Crumb = {
		title: string;
		url: string;
		isLast: boolean;
	};

	const breadcrumbs = $derived.by(() => {
		const path = page.url.pathname;
		const segments = path.split('/').filter(Boolean); // chop into parts
		let currentPath = '';
		const crumbs: Crumb[] = [];

		for (let i = 0; i < segments.length; i++) {
			currentPath += `/${segments[i]}`;
			const match = SiteRoutes.find((r) => r.url.replace(/\/$/, '') === currentPath);
			if (match) {
				crumbs.push({
					title: match.title,
					url: match.url,
					isLast: i === segments.length - 1
				});
			} else {
				// fallback to raw segment if no match found
				crumbs.push({
					title: segments[i].charAt(0).toUpperCase() + segments[i].slice(1),
					url: currentPath,
					isLast: i === segments.length - 1
				});
			}
		}

		// Special case for root
		if (path === '/') {
			const root = SiteRoutes.find((r) => r.url === '/');
			return [
				{
					title: root?.title || 'Home',
					url: '/',
					isLast: true
				}
			];
		}

		return [{ title: 'Home', url: '/', isLast: false }, ...crumbs];
	});
</script>

<header
	class="flex h-16 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-16"
>
	<div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
		<Sidebar.Trigger class="-ml-1" />
		<Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />

		<!-- Breadcrumb Navigation -->
		<Breadcrumb.Root>
			<Breadcrumb.List>
				{#each breadcrumbs as crumb, i (crumb.title)}
					<Breadcrumb.Item>
						{#if !crumb.isLast}
							<Breadcrumb.Link href={crumb.url}>{crumb.title}</Breadcrumb.Link>
						{:else}
							<Breadcrumb.Page>{crumb.title}</Breadcrumb.Page>
						{/if}
					</Breadcrumb.Item>
					{#if i < breadcrumbs.length - 1}
						<Breadcrumb.Separator />
					{/if}
				{/each}
			</Breadcrumb.List>
		</Breadcrumb.Root>

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
						<Tabs.Trigger value="middlewares" class="px-2 py-0.5 font-bold">
							Middlewares
						</Tabs.Trigger>
						<Tabs.Trigger value="plugins" class="px-2 py-0.5 font-bold">Plugins</Tabs.Trigger>
					</Tabs.List>
				</Tabs.Root>
			</div>
		{/if}
	</div>
</header>
