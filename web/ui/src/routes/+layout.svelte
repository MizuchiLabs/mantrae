<script lang="ts">
	import { queryClient, transport } from '$lib/api/client';
	import { Toaster } from '$lib/components/ui/sonner';
	import { setTransport } from '$lib/query';
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import { ModeWatcher } from 'mode-watcher';
	import './layout.css';

	let { children } = $props();
	setTransport(transport);
</script>

<ModeWatcher />
<Toaster />

<QueryClientProvider client={queryClient}>
	{@render children?.()}
</QueryClientProvider>
