<script lang="ts">
	import { Button } from '$lib/components/ui/button/index';
	import { Textarea } from '$lib/components/ui/textarea/index.js';

	import Separator from '../ui/separator/separator.svelte';
	import { loading } from '$lib/api';
	import YAML from 'yaml';

	interface Props {
		data: Record<string, unknown>;
		onSubmit: (data: Record<string, unknown>) => void;
		disabled?: boolean;
	}

	let { data = $bindable(), onSubmit, disabled }: Props = $props();

	// Form state
	let formData = $state(YAML.stringify(data, { indent: 2 }));
	let errorMessage = $state<string | null>(null);
	let isValid = $state(true);

	// Validate YAML and return boolean
	function validateYAML(input: string): boolean {
		try {
			YAML.parse(input);
			errorMessage = null;
			return true;
		} catch (e) {
			errorMessage = e instanceof Error ? e.message : 'Invalid YAML';
			return false;
		}
	}

	// Format YAML with proper indentation
	function formatYAML(input: string): string {
		try {
			const parsed = YAML.parse(input);
			return YAML.stringify(parsed, { indent: 2 });
		} catch {
			return input;
		}
	}

	// Handle input changes
	function handleInput(e: Event) {
		const input = (e.target as HTMLTextAreaElement).value;
		formData = input;
		isValid = validateYAML(input);
	}

	// Format on blur
	function handleBlur() {
		if (isValid) {
			formData = formatYAML(formData);
		}
	}

	// Handle tab key
	function handleKeydown(e: KeyboardEvent & { currentTarget: HTMLTextAreaElement }) {
		if (e.key === 'Tab') {
			e.preventDefault();
			const target = e.currentTarget;
			const start = target.selectionStart;
			const end = target.selectionEnd;
			const newPosition = start + 2;

			// Handle selected text
			if (start !== end) {
				const lines = formData.split('\n');
				let startLine = formData.substring(0, start).split('\n').length - 1;
				let endLine = formData.substring(0, end).split('\n').length - 1;

				// Indent or unindent selected lines
				const newLines = lines.map((line, i) => {
					if (i >= startLine && i <= endLine) {
						return e.shiftKey ? line.replace(/^ {2}/, '') : '  ' + line;
					}
					return line;
				});

				formData = newLines.join('\n');
				requestAnimationFrame(() => {
					target.setSelectionRange(newPosition, newPosition);
				});
			} else {
				// Insert tab at cursor position
				formData = formData.substring(0, start) + '  ' + formData.substring(end);
				requestAnimationFrame(() => {
					target.setSelectionRange(newPosition, newPosition);
				});
			}
		}
	}

	// Handle form submission
	function handleSubmit(e: Event) {
		e.preventDefault();
		if (!isValid) return;

		try {
			const parsed = YAML.parse(formData);
			onSubmit(parsed);
		} catch (e) {
			errorMessage = e instanceof Error ? e.message : 'Failed to parse YAML';
		}
	}
</script>

<form onsubmit={handleSubmit}>
	<div class="grid gap-4">
		<Textarea
			bind:value={formData}
			rows={formData.split('\n').length}
			class={!isValid ? 'border-red-500 font-mono' : 'font-mono'}
			oninput={handleInput}
			onblur={handleBlur}
			onkeydown={handleKeydown}
			{disabled}
		/>
		{#if errorMessage}
			<p class="text-sm text-red-500">{errorMessage}</p>
		{/if}
	</div>

	<Separator class="my-4" />
	<Button type="submit" class="w-full" disabled={$loading || !isValid || disabled}>Save</Button>
</form>
