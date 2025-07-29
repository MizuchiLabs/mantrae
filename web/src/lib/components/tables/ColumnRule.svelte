<script lang="ts">
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Badge, type BadgeVariant } from '$lib/components/ui/badge/index.js';
	import { List, Globe, Route, Shield, Eye } from '@lucide/svelte';
	import { truncateText } from '$lib/utils';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';

	interface Props {
		rule: string;
		protocol: ProtocolType.HTTP | ProtocolType.TCP;
		maxDisplayLength?: number;
		showIcons?: boolean;
	}

	let { rule, protocol, maxDisplayLength = 40, showIcons = true }: Props = $props();

	interface ParsedRule {
		type: 'Host' | 'Path' | 'PathPrefix' | 'HostSNI' | 'Method' | 'Headers' | 'Query' | 'Complex';
		value: string;
		isClickable: boolean;
		operator?: '&&' | '||';
		priority?: number;
	}

	function parseRule(rule: string): ParsedRule[] {
		if (!rule) return [];

		const patterns = {
			Host: { regex: /Host\(`(.*?)`\)/g, clickable: true, priority: 1 },
			HostSNI: { regex: /HostSNI\(`(.*?)`\)/g, clickable: true, priority: 1 },
			Path: { regex: /Path\(`(.*?)`\)/g, clickable: false, priority: 2 },
			PathPrefix: { regex: /PathPrefix\(`(.*?)`\)/g, clickable: false, priority: 2 },
			Method: { regex: /Method\(`(.*?)`\)/g, clickable: false, priority: 3 },
			Headers: { regex: /Headers\(`(.*?)`\)/g, clickable: false, priority: 4 },
			Query: { regex: /Query\(`(.*?)`\)/g, clickable: false, priority: 5 }
		};

		// Handle complex rules (with negation or complex logic)
		if (rule.includes('!')) {
			return [
				{
					type: 'Complex',
					value: rule,
					isClickable: false,
					priority: 10
				}
			];
		}

		const parsedRules: ParsedRule[] = [];

		// Split by logical operators while preserving them
		const parts = rule
			.split(/(\|\||&&)/)
			.map((part) => part.trim())
			.filter(Boolean);

		for (let i = 0; i < parts.length; i++) {
			const part = parts[i];

			// Handle operators
			if (part === '&&' || part === '||') {
				if (parsedRules.length > 0) {
					parsedRules[parsedRules.length - 1].operator = part;
				}
				continue;
			}

			// Try to match patterns
			let matched = false;
			for (const [type, config] of Object.entries(patterns)) {
				const matches = [...part.matchAll(config.regex)];
				if (matches.length > 0) {
					matches.forEach((match) => {
						parsedRules.push({
							type: type as ParsedRule['type'],
							value: match[1],
							isClickable: config.clickable,
							priority: config.priority
						});
					});
					matched = true;
					break;
				}
			}

			// If no pattern matched, treat as complex
			if (!matched) {
				parsedRules.push({
					type: 'Complex',
					value: part,
					isClickable: false,
					priority: 10
				});
			}
		}

		// Sort by priority for better display
		return parsedRules.sort((a, b) => (a.priority || 10) - (b.priority || 10));
	}

	let parsedRules: ParsedRule[] = $derived(parseRule(rule));
	const shouldShowMultiple = $derived(parsedRules.length > 1);
	const primaryRule = $derived(parsedRules[0]);

	function getUrl(domain: string): string {
		// For HTTP routers, try HTTPS first, fallback to HTTP
		if (protocol === ProtocolType.HTTP) {
			return `https://${domain}`;
		}
		// For TCP/SNI, just show the domain
		return domain;
	}

	function getRuleIcon(type: ParsedRule['type']) {
		const iconMap = {
			Host: Globe,
			HostSNI: Shield,
			Path: Route,
			PathPrefix: Route,
			Method: Eye,
			Headers: Eye,
			Query: Eye,
			Complex: List
		};
		return iconMap[type] || List;
	}

	function getRuleVariant(type: ParsedRule['type']): BadgeVariant {
		const variantMap = {
			Host: 'default',
			HostSNI: 'secondary',
			Path: 'outline',
			PathPrefix: 'outline',
			Method: 'outline',
			Headers: 'outline',
			Query: 'outline',
			Complex: 'destructive'
		};
		if (type in variantMap) return variantMap[type] as BadgeVariant;
		return 'outline';
		// return variantMap[type] || 'outline';
	}

	function truncateValue(value: string, maxLen: number): string {
		return value.length > maxLen ? `${value.slice(0, maxLen)}...` : value;
	}
</script>

{#if parsedRules.length === 0}
	<Badge variant="outline" class="text-sm transition-colors duration-200">None</Badge>
{:else if !shouldShowMultiple && primaryRule}
	<!-- Single rule display -->
	{@const Icon = showIcons ? getRuleIcon(primaryRule.type) : null}
	{@const variant = getRuleVariant(primaryRule.type)}
	{@const truncated = truncateValue(primaryRule.value, maxDisplayLength)}
	{@const shouldShowTooltip = primaryRule.value.length > maxDisplayLength}

	{#if primaryRule.isClickable}
		{#if shouldShowTooltip}
			<Tooltip.Provider>
				<Tooltip.Root delayDuration={300}>
					<Tooltip.Trigger>
						<a
							href={getUrl(primaryRule.value)}
							target="_blank"
							rel="noopener noreferrer"
							class="inline-flex max-w-full items-center gap-1.5 text-sm text-blue-600
								   transition-colors duration-200 hover:text-blue-800 dark:text-blue-400
								   dark:hover:text-blue-300"
						>
							<span class="truncate">{truncated}</span>
						</a>
					</Tooltip.Trigger>
					<Tooltip.Content side="top" class="max-w-xs break-words">
						<div class="space-y-1">
							<div class="font-medium">{primaryRule.type}</div>
							<div class="text-xs">{primaryRule.value}</div>
						</div>
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{:else}
			<a
				href={getUrl(primaryRule.value)}
				target="_blank"
				rel="noopener noreferrer"
				class="inline-flex max-w-full items-center gap-1.5 text-sm text-blue-600
					   transition-colors duration-200 hover:text-blue-800 dark:text-blue-400
					   dark:hover:text-blue-300"
			>
				<span class="truncate">{truncated}</span>
			</a>
		{/if}
	{:else if shouldShowTooltip}
		<Tooltip.Provider>
			<Tooltip.Root delayDuration={300}>
				<Tooltip.Trigger>
					<Badge {variant} class="max-w-full">
						<span class="truncate">{truncated}</span>
					</Badge>
				</Tooltip.Trigger>
				<Tooltip.Content side="top" class="max-w-xs break-words">
					<div class="space-y-1">
						<div class="font-medium">{primaryRule.type}</div>
						<div class="text-xs">{primaryRule.value}</div>
					</div>
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{:else}
		<Badge {variant} class="max-w-full">
			{#if Icon}
				<Icon class="mr-1 h-3 w-3 shrink-0" />
			{/if}
			<span class="truncate">{truncated}</span>
		</Badge>
	{/if}
{:else}
	<!-- Multiple rules display -->
	<HoverCard.Root openDelay={200}>
		<HoverCard.Trigger>
			<div class="inline-flex items-center gap-2">
				<!-- Show primary rule on desktop -->
				<div class="hidden max-w-[200px] items-center gap-1.5 sm:flex">
					{#if primaryRule.isClickable}
						<a
							href={getUrl(primaryRule.value)}
							target="_blank"
							rel="noopener noreferrer"
							class="inline-flex items-center gap-1 truncate text-sm text-blue-600
								   hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
						>
							<span class="truncate">{truncateText(primaryRule.value, 20)}</span>
						</a>
					{:else}
						<Badge variant={getRuleVariant(primaryRule.type)} class="h-6 text-xs">
							<span class="truncate">{truncateText(primaryRule.value, 20)}</span>
						</Badge>
					{/if}
				</div>

				<!-- Multiple indicator -->
				<Badge
					variant="outline"
					class="hover:bg-muted cursor-pointer text-xs transition-colors duration-200"
				>
					<span class="hidden sm:inline">+{parsedRules.length - 1}</span>
					<span class="sm:hidden">{parsedRules.length} rules</span>
				</Badge>
			</div>
		</HoverCard.Trigger>

		<HoverCard.Content class="w-auto max-w-md" side="bottom" sideOffset={8}>
			<div class="space-y-3">
				<div class="border-border flex items-center gap-2 border-b pb-2">
					<List class="text-muted-foreground h-4 w-4" />
					<span class="text-sm font-medium">Traefik Rules ({parsedRules.length})</span>
				</div>

				<div class="max-h-60 space-y-1 overflow-y-auto">
					{#each parsedRules as rule, index (rule.value + index)}
						<div class="bg-muted/30 flex items-center gap-2 rounded-md p-2">
							{#if showIcons}
								{@const Icon = getRuleIcon(rule.type)}

								<div class="flex min-w-0 flex-1 items-center gap-1">
									<Icon class="text-muted-foreground h-3 w-3 shrink-0" />
									<span class="text-muted-foreground min-w-0 text-xs font-medium">
										{rule.type}:
									</span>

									{#if rule.isClickable}
										<a
											href={getUrl(rule.value)}
											target="_blank"
											rel="noopener noreferrer"
											class="flex flex-1 items-center gap-1
											   truncate text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
										>
											<span class="truncate">{rule.value}</span>
										</a>
									{:else}
										<span class="flex-1 truncate text-sm">{rule.value}</span>
									{/if}
								</div>

								{#if rule.operator}
									<Badge variant="outline" class="h-5 px-1 py-0 text-xs">
										{rule.operator}
									</Badge>
								{/if}
							{/if}
						</div>
					{/each}
				</div>
			</div>
		</HoverCard.Content>
	</HoverCard.Root>
{/if}
