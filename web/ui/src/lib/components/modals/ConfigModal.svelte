<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import CopyButton from '$lib/components/ui/copy-button/copy-button.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { createHighlighter } from 'shiki';
	import { util } from '$lib/api/util.svelte';
	import YAML from 'yaml';

	type Props = {
		open?: boolean;
	};
	let { open = $bindable(false) }: Props = $props();

	const config = $derived(util.config());
	let lang: 'json' | 'yaml' = $state('json');

	const formatted = $derived.by(() => {
		if (!config.data) return '';
		return lang === 'json'
			? JSON.stringify(JSON.parse(config.data), null, 2)
			: YAML.stringify(JSON.parse(config.data), {
					indent: 2,
					lineWidth: 0, // no forced wrapping
					collectionStyle: 'block' // no inline `{}` or `[]`
				});
	});

	let code = $derived.by(async () => {
		const currentLang = lang;

		let highlighter = await createHighlighter({
			themes: ['catppuccin-latte', 'catppuccin-mocha'],
			langs: ['json', 'yaml']
		});

		return highlighter.codeToHtml(formatted, {
			lang: currentLang,
			themes: {
				light: 'catppuccin-latte',
				dark: 'catppuccin-mocha'
			}
		});
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="max-h-[90vh] w-fit max-w-[90vw] overflow-y-auto px-4 py-2 sm:min-w-160">
		<Dialog.Header class="flex justify-between gap-2 py-4">
			<Dialog.Title>Dynamic Config</Dialog.Title>
			<Dialog.Description>
				This is the current dynamic configuration of this profile.
			</Dialog.Description>
		</Dialog.Header>
		<Button
			variant="outline"
			size="sm"
			class="w-full"
			onclick={() => (lang = lang === 'json' ? 'yaml' : 'json')}
		>
			{lang === 'json' ? 'Show YAML' : 'Show JSON'}
		</Button>

		{#if config.isSuccess && config.data !== ''}
			<div class="relative overflow-x-auto rounded-xl">
				<CopyButton text={formatted} class="absolute top-1 right-4" />

				{#await code then value}
					<div class="w-full overflow-x-auto">
						<div class="min-w-fit whitespace-pre">
							{@html value}
						</div>
					</div>
				{/await}
			</div>
		{/if}
	</Dialog.Content>
</Dialog.Root>

<style>
	:global(.dark .shiki),
	:global(.dark .shiki span) {
		color: var(--shiki-dark) !important;
		background-color: var(--shiki-dark-bg) !important;
		font-style: var(--shiki-dark-font-style) !important;
		font-weight: var(--shiki-dark-font-weight) !important;
		text-decoration: var(--shiki-dark-text-decoration) !important;
	}
	:global(.shiki) {
		padding: 1rem !important;
		border-radius: 0.5rem;
		margin: 0 !important;
	}

	:global(.shiki code) {
		padding: 0 !important;
	}
</style>
