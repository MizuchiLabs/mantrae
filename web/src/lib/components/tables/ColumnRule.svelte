<script lang="ts">
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { RouterType } from '$lib/gen/mantrae/v1/router_pb';
	import { Link, List } from '@lucide/svelte';

	interface Props {
		rule: string;
		routerType: RouterType.HTTP | RouterType.TCP;
	}
	let { rule, routerType }: Props = $props();

	interface ParsedRule {
		value: string;
		isClickable: boolean;
		operator?: '&&' | '||';
	}

	function parseRule(rule: string): ParsedRule[] {
		if (!rule) return [];

		const patterns = {
			Host: /Host\(`(.*?)`\)/g,
			Path: /Path\(`(.*?)`\)/g,
			PathPrefix: /PathPrefix\(`(.*?)`\)/g,
			HostSNI: /HostSNI\(`(.*?)`\)/g
		};

		// If rule contains negation, return as complex
		if (rule.includes('!')) {
			return [
				{
					value: rule,
					isClickable: false
				}
			];
		}

		const parsedRules: ParsedRule[] = [];

		// Split by && and || while preserving the operators
		const parts = rule
			.split(/(\|\||&&)/)
			.map((part) => part.trim())
			.filter(Boolean);

		for (let i = 0; i < parts.length; i++) {
			const part = parts[i];

			// If it's an operator, add it to the previous rule
			if (part === '&&' || part === '||') {
				if (parsedRules.length > 0) {
					parsedRules[parsedRules.length - 1].operator = part;
				}
				continue;
			}

			// Try to match each pattern
			let matched = false;
			for (const [type, pattern] of Object.entries(patterns)) {
				const matches = [...part.matchAll(pattern)];
				if (matches.length > 0) {
					matches.forEach((match) => {
						parsedRules.push({
							value: match[1],
							isClickable: ['Host', 'HostSNI'].includes(type)
						});
					});
					matched = true;
					break;
				}
			}

			// If no patterns matched for this part, treat it as raw
			if (!matched) {
				parsedRules.push({
					value: part,
					isClickable: false
				});
			}
		}

		return parsedRules;
	}

	let parsedRules: ParsedRule[] = $derived(parseRule(rule));

	function getUrl(domain: string): string {
		const prefix = routerType === RouterType.HTTP ? 'http://' : 'https://';
		return `${prefix}${domain}`;
	}
</script>

{#if parsedRules.length === 0}
	<Badge variant="outline" class="text-xs">None</Badge>
{:else if parsedRules.length === 1}
	{#if parsedRules[0].isClickable}
		<a
			href={getUrl(parsedRules[0].value)}
			target="_blank"
			rel="noopener noreferrer"
			class="inline-flex items-center gap-1.5 text-sm text-blue-500 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
		>
			<Link class="h-4 w-4" />
			{parsedRules[0].value}
		</a>
	{/if}
{:else}
	<HoverCard.Root openDelay={100}>
		<HoverCard.Trigger>
			<button class="inline-flex items-center gap-1.5 text-sm">
				<List class="h-4 w-4" />
				<Badge variant="outline" class="text-xs">Multiple</Badge>
			</button>
		</HoverCard.Trigger>
		<HoverCard.Content class="w-auto">
			<div class="flex flex-col gap-2">
				{#each parsedRules as { value, isClickable } (value)}
					<div class="flex flex-col">
						<div class="flex items-center gap-2">
							{#if isClickable}
								<div class="bg-muted-foreground/20 h-4 w-0.5"></div>
								<a
									href={getUrl(value)}
									target="_blank"
									rel="noopener noreferrer"
									class="inline-flex items-center gap-1 text-sm text-blue-500 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
								>
									{value}
									<Link class="h-3 w-3" />
								</a>
							{:else}
								<div class="bg-muted-foreground/20 h-4 w-0.5"></div>
								<span class="text-sm">{value}</span>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		</HoverCard.Content>
	</HoverCard.Root>
{/if}
