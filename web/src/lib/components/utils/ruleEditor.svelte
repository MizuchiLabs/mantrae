<script lang="ts">
	import { Label } from '$lib/components/ui/label/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { ValidateRule } from './ruleString';

	export let rule: string;

	const rules = [
		'Header(`key`, `value`)',
		'HeaderRegexp(`key`, `regexp`)',
		'Host(`domain`)',
		'HostRegexp(`regexp`)',
		'Method(`method`)',
		'Path(`path`)',
		'PathPrefix(`prefix`)',
		'PathRegexp(`regexp`)',
		'Query(`key`, `value`)',
		'QueryRegexp(`key`, `regexp`)',
		'ClientIP(`ip`)'
	];

	let valid = true;
	let cursorPosition = 0;
	let showDropdown = false;
	let filteredRules: string[] = [];
	let selectedRuleIndex = 0;
	let placeholderPositions: { start: number; end: number }[] = [];
	let currentPlaceholderIndex = 0;
	function handleRuleInput(event: InputEvent) {
		if (event.target === undefined) return;
		cursorPosition = (event.target as HTMLTextAreaElement).selectionStart;
		const lastWord = rule?.slice(0, cursorPosition).split(/\s+/).pop();

		if (lastWord === undefined) return;
		if (lastWord.length > 0) {
			filteredRules = rules.filter((rule) => rule.toLowerCase().startsWith(lastWord.toLowerCase()));
			showDropdown = filteredRules.length > 0;
			selectedRuleIndex = 0;
		} else {
			showDropdown = false;
		}

		// Recalculate the placeholder positions
		placeholderPositions = [];
		const regex = /`([^`]*)`/g;
		let match;
		while ((match = regex.exec(rule)) !== null) {
			placeholderPositions.push({
				start: match.index + 1,
				end: match.index + match[0].length - 1
			});
		}
		valid = ValidateRule(rule);
	}

	function insertRule(newRule: string) {
		const textarea = document.getElementById('rulesTextarea') as HTMLTextAreaElement;
		const beforeCursor = textarea.value.slice(0, cursorPosition);
		const afterCursor = textarea.value.slice(cursorPosition);

		// Replace the last word with the selected rule
		const lastWord = beforeCursor.split(/\s+/).pop();
		if (lastWord === undefined) return;

		const newCursorPos = cursorPosition - (lastWord?.length || 0);
		textarea.value = beforeCursor.slice(0, -lastWord.length) + newRule + afterCursor;

		// Update the Svelte reactive variable
		rule = textarea.value;

		// Calculate positions of placeholders (text inside backticks)
		placeholderPositions = [];
		const regex = /`([^`]*)`/g;
		let match;
		while ((match = regex.exec(newRule)) !== null) {
			placeholderPositions.push({
				start: newCursorPos + match.index + 1,
				end: newCursorPos + match.index + match[0].length - 1
			});
		}

		currentPlaceholderIndex = 0;
		moveToPlaceholder(textarea);

		// Hide the dropdown
		showDropdown = false;
	}

	function moveToPlaceholder(textarea: HTMLTextAreaElement) {
		if (currentPlaceholderIndex < placeholderPositions.length) {
			const { start, end } = placeholderPositions[currentPlaceholderIndex];
			textarea.focus();
			textarea.setSelectionRange(start, end);
		} else {
			textarea.focus();
			const newCursorPos = placeholderPositions[placeholderPositions.length - 1].end + 1;
			textarea.setSelectionRange(newCursorPos, newCursorPos);
		}
	}

	function appendAndOperator() {
		const textarea = document.getElementById('rulesTextarea') as HTMLTextAreaElement;
		const textValue = textarea.value;
		const lastNonSpaceIndex = textValue.search(/\S\s*$/); // Find the last non-space character

		// Insert ' && ' after the last non-space character
		textarea.value =
			textValue.slice(0, lastNonSpaceIndex + 1) + ' && ' + textValue.slice(lastNonSpaceIndex + 1);

		// Move cursor to the end
		textarea.focus();
		textarea.setSelectionRange(textarea.value.length, textarea.value.length);
	}

	function handleRuleKeys(event: KeyboardEvent) {
		if (showDropdown) {
			if (event.key === 'ArrowDown' || event.key === 'Tab') {
				event.preventDefault();
				selectedRuleIndex = (selectedRuleIndex + 1) % filteredRules.length;
			} else if (event.key === 'ArrowUp' || (event.shiftKey && event.key === 'Tab')) {
				event.preventDefault();
				selectedRuleIndex = (selectedRuleIndex - 1 + filteredRules.length) % filteredRules.length;
			} else if (event.key === 'Enter') {
				event.preventDefault();
				insertRule(filteredRules[selectedRuleIndex]);
			} else if (event.key === 'Escape') {
				showDropdown = false;
			}
		} else {
			if (!event.shiftKey && event.key === 'Tab') {
				event.preventDefault();
				// Move to the next placeholder
				currentPlaceholderIndex = (currentPlaceholderIndex + 1) % placeholderPositions.length;
				moveToPlaceholder(document.getElementById('rulesTextarea') as HTMLTextAreaElement);
			} else if (event.shiftKey && event.key === 'Tab') {
				event.preventDefault();
				appendAndOperator();
			} else if (event.shiftKey && event.key === 'Enter') {
				// unfocus the textarea
				document.getElementById('rulesTextarea')?.blur();
			}
		}
	}
</script>

<div class="space-y-1">
	<Label for="rule">Rules</Label>
	<div class="mb-4 rounded-lg border border-gray-200">
		<div class="rounded-t-lg">
			<Textarea
				placeholder="Add rules here"
				rows={3}
				id="rulesTextarea"
				bind:value={rule}
				class="w-full border-0 font-mono text-sm focus-visible:ring-0 focus-visible:ring-offset-0"
				on:input={handleRuleInput}
				on:keydown={handleRuleKeys}
			/>
			{#if showDropdown}
				<ul
					class="absolute mt-1 flex max-h-48 w-80 flex-col gap-2 overflow-y-auto rounded-lg border bg-white p-2 dark:bg-gray-800"
				>
					{#each filteredRules as rule, i}
						<li
							class="cursor-pointer font-mono text-sm hover:bg-gray-200"
							class:bg-gray-200={i === selectedRuleIndex}
							on:click={() => insertRule(rule)}
							aria-hidden
						>
							{rule}
						</li>
					{/each}
				</ul>
			{/if}
		</div>
		<div
			class="flex items-center justify-end gap-1 border-t px-3 py-2 text-sm text-muted-foreground dark:border-gray-600"
		>
			{#if valid}
				<p>Valid</p>
				<iconify-icon icon="fa6-solid:circle-check" />
			{:else}
				<p>Invalid</p>
				<iconify-icon icon="fa6-solid:circle-xmark" />
			{/if}
		</div>
	</div>
	<div class="ml-2 flex items-center justify-between">
		<div class="text-xs text-muted-foreground">
			<span class="font-bold">Rule Examples:</span>
			<ul class="list-inside list-disc">
				<li>Host(`example.com`)</li>
				<li>Path(`/hello`)</li>
				<li>PathPrefix(`/hello`)</li>
				<li>PathRegexp(`/hello/[0-9]+`)</li>
				<li>Method(`GET`)</li>
				<li>Header(`X-Forwarded-For`, `.*`)</li>
				<li>Query(`page`, `[0-9]+`)</li>
				<li>QueryRegexp(`page`, `[0-9]+`)</li>
			</ul>
		</div>
	</div>
</div>
