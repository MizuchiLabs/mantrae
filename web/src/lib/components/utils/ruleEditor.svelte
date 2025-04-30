<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import { ValidateRule } from './ruleString';
	import { CircleCheck, CircleX } from '@lucide/svelte';
	import { ruleTab } from '$lib/stores/common';

	interface Props {
		rule: string | undefined;
		type: 'http' | 'tcp';
		disabled?: boolean;
	}

	let { rule = $bindable(), type = $bindable(), disabled = false }: Props = $props();

	// Simplified rule templates
	const ruleTemplates = {
		http: [
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
		],
		tcp: ['HostSNI(`domain`)', 'HostSNIRegexp(`regexp`)', 'ClientIP(`ip`)', 'ALPN(`protocol`)']
	};

	// Editor state
	let valid = $state(true);
	let showDropdown = $state(false);
	let filteredRules = $state<string[]>([]);
	let selectedRuleIndex = $state(0);
	let simpleDisabled = $state(false);
	let cursorPosition = $state(0);
	let placeholderPositions: { start: number; end: number }[] = [];
	let currentPlaceholderIndex = 0;

	// Simple mode state
	let host = $state('');
	let path = $state('');

	// Parse existing rule on mount
	$effect(() => {
		if (!rule) return;

		const hostPattern = type === 'http' ? /Host\(`(.*?)`\)/ : /HostSNI\(`(.*?)`\)/;
		const pathPattern = /Path\(`(.*?)`\)/;

		host = rule.match(hostPattern)?.[1] ?? '';
		path = type === 'http' ? (rule.match(pathPattern)?.[1] ?? '') : '';

		checkSimpleMode();
	});

	// Check if rule can be displayed in simple mode
	function checkSimpleMode() {
		if (!rule) {
			simpleDisabled = false;
			return;
		}

		const conditions = rule
			.split(/(&&|\|\|)/)
			.filter((part) => part.trim() && !['&&', '||'].includes(part));
		simpleDisabled =
			conditions.length > 2 ||
			conditions.filter((c) => c.includes('Host')).length > 1 ||
			conditions.filter((c) => c.includes('Path')).length > 1;

		if (simpleDisabled && ruleTab.value === 'simple') {
			ruleTab.value = 'advanced';
		}
	}

	// Simple mode handler
	function handleSimpleInput() {
		if (!host && !path) {
			rule = undefined;
			return;
		}

		if (type === 'http') {
			rule = [host ? `Host(\`${host}\`)` : null, path ? `Path(\`${path}\`)` : null]
				.filter(Boolean)
				.join(' && ');
		} else {
			rule = host ? `HostSNI(\`${host}\`)` : undefined;
		}
	}

	function updatePlaceholderPositions() {
		placeholderPositions = [];
		if (!rule) return;

		const regex = /`([^`]*)`/g;
		let match;
		while ((match = regex.exec(rule)) !== null) {
			placeholderPositions.push({
				start: match.index + 1,
				end: match.index + match[0].length - 1
			});
		}
	}

	// Advanced mode handlers
	function handleRuleInput(event: Event) {
		const textarea = event.target as HTMLTextAreaElement;
		cursorPosition = textarea.selectionStart;
		const lastWord = rule?.slice(0, cursorPosition).split(/\s+/).pop();

		if (!lastWord) {
			showDropdown = false;
			return;
		}

		filteredRules = ruleTemplates[type].filter((template) =>
			template.toLowerCase().startsWith(lastWord.toLowerCase())
		);

		showDropdown = filteredRules.length > 0;
		selectedRuleIndex = 0;
		valid = rule ? ValidateRule(rule) : true;

		updatePlaceholderPositions();
	}

	function insertRule(template: string) {
		const textarea = document.getElementById('rulesTextarea') as HTMLTextAreaElement;
		const beforeCursor = rule?.slice(0, cursorPosition) ?? '';
		const afterCursor = rule?.slice(cursorPosition) ?? '';
		const lastWord = beforeCursor.split(/\s+/).pop() ?? '';

		// const newCursorPos = cursorPosition - lastWord.length;
		rule = beforeCursor.slice(0, -lastWord.length) + template + afterCursor;

		showDropdown = false;
		currentPlaceholderIndex = 0;
		updatePlaceholderPositions();

		setTimeout(() => {
			textarea.focus();
			if (placeholderPositions.length > 0) {
				const { start, end } = placeholderPositions[0];
				textarea.setSelectionRange(start, end);
			}
		}, 0);
	}

	function moveToPlaceholder(textarea: HTMLTextAreaElement) {
		if (currentPlaceholderIndex < placeholderPositions.length) {
			const { start, end } = placeholderPositions[currentPlaceholderIndex];
			textarea.focus();
			textarea.setSelectionRange(start, end);
		}
	}

	function appendAndOperator() {
		if (!rule) return;
		const textarea = document.getElementById('rulesTextarea') as HTMLTextAreaElement;
		const lastNonSpaceIndex = rule.search(/\S\s*$/);

		rule = rule.slice(0, lastNonSpaceIndex + 1) + ' && ' + rule.slice(lastNonSpaceIndex + 1);

		textarea.focus();
		textarea.setSelectionRange(rule.length, rule.length);
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (showDropdown) {
			switch (event.key) {
				case 'ArrowDown':
					event.preventDefault();
					selectedRuleIndex = (selectedRuleIndex + 1) % filteredRules.length;
					break;
				case 'ArrowUp':
					event.preventDefault();
					selectedRuleIndex = (selectedRuleIndex - 1 + filteredRules.length) % filteredRules.length;
					break;
				case 'Enter':
					event.preventDefault();
					insertRule(filteredRules[selectedRuleIndex]);
					break;
				case 'Escape':
					showDropdown = false;
					break;
			}
		} else {
			if (event.key === 'Tab' && !event.shiftKey) {
				event.preventDefault();
				const textarea = event.target as HTMLTextAreaElement;
				currentPlaceholderIndex = (currentPlaceholderIndex + 1) % placeholderPositions.length;
				moveToPlaceholder(textarea);
			} else if (event.key === 'Tab' && event.shiftKey) {
				event.preventDefault();
				appendAndOperator();
			}
		}
	}
</script>

<Tabs.Root
	value={ruleTab.value}
	onValueChange={(value) => (ruleTab.value = value)}
	class="flex flex-col gap-2"
>
	<div class="flex justify-end">
		<Tabs.List class="h-8">
			<Tabs.Trigger value="simple" class="px-2 py-0.5 font-bold" disabled={simpleDisabled}
				>Simple</Tabs.Trigger
			>
			<Tabs.Trigger value="advanced" class="px-2 py-0.5 font-bold">Advanced</Tabs.Trigger>
		</Tabs.List>
	</div>

	<Tabs.Content value="simple">
		<div class="grid grid-cols-8 items-center gap-2">
			<Label for="host" class="col-span-1 text-right">Domain</Label>
			<Input
				id="host"
				bind:value={host}
				oninput={handleSimpleInput}
				placeholder="example.com"
				class={type === 'http' ? 'col-span-5' : 'col-span-7'}
				{disabled}
			/>
			{#if type === 'http'}
				<Input
					id="path"
					bind:value={path}
					oninput={handleSimpleInput}
					placeholder="/path"
					class="col-span-2"
					{disabled}
				/>
			{/if}
		</div>
	</Tabs.Content>

	<Tabs.Content value="advanced">
		<Label for="rule">Rules</Label>
		<div class="relative mb-4 rounded-lg border">
			<Textarea
				id="rulesTextarea"
				placeholder="Add rules here"
				rows={3}
				bind:value={rule}
				class="w-full border-0 font-mono text-sm focus-visible:ring-0"
				oninput={handleRuleInput}
				onkeydown={handleKeyDown}
				{disabled}
			/>

			{#if showDropdown}
				<ul class="bg-card absolute mt-1 max-h-48 w-80 overflow-y-auto rounded-lg border p-2">
					{#each filteredRules as template, i (template)}
						<button
							role="option"
							aria-selected={i === selectedRuleIndex}
							class="hover:bg-muted w-full cursor-pointer rounded p-1 text-left font-mono text-sm"
							class:bg-muted={i === selectedRuleIndex}
							onclick={() => insertRule(template)}
						>
							{template}
						</button>
					{/each}
				</ul>
			{/if}

			{#if !disabled}
				<div
					class="text-muted-foreground flex items-center justify-end gap-1 border-t px-3 py-2 text-sm"
				>
					{#if valid}
						<span>Valid</span>
						<CircleCheck size="1rem" />
					{:else}
						<span>Invalid</span>
						<CircleX size="1rem" />
					{/if}
				</div>
			{/if}
		</div>

		{#if !disabled}
			<div class="text-muted-foreground text-xs">
				<span class="font-bold">Examples:</span>
				<ul class="list-inside list-disc">
					{#each ruleTemplates[type] as template (template)}
						<li>{template}</li>
					{/each}
				</ul>
			</div>
		{/if}
	</Tabs.Content>
</Tabs.Root>
