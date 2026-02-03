<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import CopyButton from '$lib/components/ui/copy-button/copy-button.svelte';
	import { createHighlighter } from 'shiki';
	import { util } from '$lib/api/util.svelte';
	import YAML from 'yaml';
	import { onMount } from 'svelte';

	type Props = {
		open?: boolean;
	};
	let { open = $bindable(false) }: Props = $props();

	const config = $derived(util.config());
	let lang: 'json' | 'yaml' = $state('json');
	let highlighter: Awaited<ReturnType<typeof createHighlighter>> | null = $state(null);

	onMount(async () => {
		highlighter = await createHighlighter({
			themes: ['catppuccin-latte', 'catppuccin-mocha'],
			langs: ['json', 'yaml']
		});
	});

	const formatted = $derived.by(() => {
		if (!config.data) return '';
		try {
			const obj = JSON.parse(config.data);
			return lang === 'json'
				? JSON.stringify(obj, null, 2)
				: YAML.stringify(obj, {
						indent: 2,
						lineWidth: 0,
						collectionStyle: 'block'
					});
		} catch (e) {
			return config.data;
		}
	});

	const codeHtml = $derived.by(() => {
		if (!highlighter || !formatted) return '';
		return highlighter.codeToHtml(formatted, {
			lang,
			themes: {
				light: 'catppuccin-latte',
				dark: 'catppuccin-mocha'
			}
		});
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="max-h-[85vh] w-full max-w-4xl p-0 gap-0 overflow-hidden flex flex-col">
		<Dialog.Header class="px-6 py-4 border-b">
			<Dialog.Title>Configuration</Dialog.Title>
			<Dialog.Description>
				View the current dynamic configuration.
			</Dialog.Description>
		</Dialog.Header>

		<div class="flex flex-col flex-1 min-h-0">
			<Tabs.Root value={lang} onValueChange={(v) => (lang = v as 'json' | 'yaml')} class="flex flex-col flex-1 min-h-0">
				<div class="flex items-center justify-between px-4 py-2 border-b bg-muted/40">
					<Tabs.List class="grid w-[200px] grid-cols-2">
						<Tabs.Trigger value="json">JSON</Tabs.Trigger>
						<Tabs.Trigger value="yaml">YAML</Tabs.Trigger>
					</Tabs.List>
					<CopyButton text={formatted} />
				</div>

				<div class="flex-1 overflow-auto bg-card relative min-h-0">
					{#if codeHtml}
						<div class="p-4 text-sm font-mono leading-relaxed tab-size-2">
							{@html codeHtml}
						</div>
					{:else}
						<div class="p-6 space-y-3">
							<div class="h-4 w-3/4 bg-muted/50 rounded animate-pulse"></div>
							<div class="h-4 w-1/2 bg-muted/50 rounded animate-pulse"></div>
							<div class="h-4 w-2/3 bg-muted/50 rounded animate-pulse"></div>
							<div class="h-4 w-1/3 bg-muted/50 rounded animate-pulse"></div>
						</div>
					{/if}
				</div>
			</Tabs.Root>
		</div>
	</Dialog.Content>
</Dialog.Root>

<style>
	:global(.dark .shiki),
	:global(.dark .shiki span) {
		color: var(--shiki-dark) !important;
		background-color: transparent !important;
		font-style: var(--shiki-dark-font-style) !important;
		font-weight: var(--shiki-dark-font-weight) !important;
		text-decoration: var(--shiki-dark-text-decoration) !important;
	}
	
	:global(.shiki) {
		background-color: transparent !important;
		margin: 0 !important;
	}

	:global(.shiki code) {
		background-color: transparent !important;
		padding: 0 !important;
	}
	
	.tab-size-2 {
		-moz-tab-size: 2;
		-o-tab-size: 2;
		tab-size: 2;
	}
</style>
