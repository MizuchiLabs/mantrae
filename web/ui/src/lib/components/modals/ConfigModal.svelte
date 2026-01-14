<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import CopyButton from '$lib/components/ui/copy-button/copy-button.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { createHighlighter } from 'shiki';
	import { getConfig } from '$lib/api';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';

	type Props = {
		open?: boolean;
		item: Profile;
	};

	let { open = $bindable(false), item = $bindable() }: Props = $props();

	let lang: 'json' | 'yaml' = $state('yaml');
	let config = $derived(getConfig(lang, item));

	let code = $derived.by(async () => {
		const currentLang = lang;
		const currentConfig = await config;

		let highlighter = await createHighlighter({
			themes: ['catppuccin-latte', 'catppuccin-mocha'],
			langs: ['json', 'yaml']
		});

		return highlighter.codeToHtml(currentConfig, {
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

		<div class="relative overflow-x-auto rounded-xl">
			{#await config then value}
				<CopyButton text={value} class="absolute top-1 right-4" />
			{/await}

			{#await code then value}
				<div class="w-full overflow-x-auto">
					<div class="min-w-fit whitespace-pre">
						{@html value}
					</div>
				</div>
			{/await}
		</div>
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
