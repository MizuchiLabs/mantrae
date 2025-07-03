<script lang="ts" generics="TData">
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { AlertTriangle, Eye, Globe, Key, Shield, Lock, Unlock } from '@lucide/svelte';
	import type { Column } from '@tanstack/table-core';
	import type { RouterTCPTLSConfig, RouterTLSConfig } from '$lib/gen/zen/traefik-schemas';
	import { truncateText } from '$lib/utils';

	type Props = {
		tls?: RouterTLSConfig | RouterTCPTLSConfig;
		column?: Column<TData, unknown>;
		compact?: boolean;
	};

	let { tls, column, compact = false }: Props = $props();

	// Type guard to check if TLS config has passthrough (TCP only)
	function isTCPTLSConfig(
		config: RouterTLSConfig | RouterTCPTLSConfig
	): config is RouterTCPTLSConfig {
		return 'passthrough' in config;
	}

	// Determine TLS status and configuration
	const tlsStatus = $derived(() => {
		if (!tls) return { enabled: false, type: 'disabled' };

		const hasResolver = !!tls.certResolver;
		const hasPassthrough = isTCPTLSConfig(tls) && !!tls.passthrough;
		const hasOptions = !!tls.options;
		const hasDomains = !!tls.domains?.length;
		const hasConfig = hasResolver || hasPassthrough || hasOptions || hasDomains;

		return {
			enabled: true,
			type: hasConfig ? 'configured' : 'basic',
			hasResolver,
			hasPassthrough,
			hasOptions,
			hasDomains,
			configCount: [hasResolver, hasPassthrough, hasOptions, hasDomains].filter(Boolean).length
		};
	});

	function getBadgeVariant() {
		if (!tlsStatus().enabled) return 'outline';
		return tlsStatus().type === 'configured' ? 'default' : 'secondary';
	}

	function getStatusIcon() {
		const status = tlsStatus();
		if (!status.enabled) return Unlock;
		if (status.type === 'configured') return Shield;
		return Lock;
	}

	function getStatusText() {
		const status = tlsStatus();
		if (!status.enabled) return 'Disabled';
		if (status.hasResolver) return truncateText(tls?.certResolver || 'Auto SSL', 20);
		if (status.hasPassthrough) return 'Passthrough';
		if (status.hasOptions) return 'Custom';
		return 'Enabled';
	}
</script>

{#if tls}
	<HoverCard.Root openDelay={200}>
		<HoverCard.Trigger>
			{@const StatusIcon = getStatusIcon()}
			{@const variant = getBadgeVariant()}
			{@const status = tlsStatus()}

			<Badge
				{variant}
				onclick={() => column?.setFilterValue?.(tls?.certResolver ?? 'enabled')}
				class="cursor-pointer transition-colors duration-200 hover:shadow-sm
					   {compact ? 'px-2 text-xs' : 'text-sm'}"
			>
				<StatusIcon class="mr-1 h-3 w-3 shrink-0" />
				<span class="max-w-[100px] truncate sm:max-w-none">
					{getStatusText()}
				</span>
				{#if !compact && status.configCount !== undefined && status.configCount > 1}
					<span class="ml-1 text-xs opacity-70">+{status.configCount - 1}</span>
				{/if}
			</Badge>
		</HoverCard.Trigger>

		<HoverCard.Content class="w-auto max-w-sm" side="bottom" sideOffset={8}>
			<div class="space-y-3">
				<!-- Header -->
				<div class="border-border flex items-center gap-2 border-b pb-2">
					<Shield class="h-4 w-4 text-green-500" />
					<span class="text-sm font-medium">TLS Configuration</span>
				</div>

				<!-- Configuration Details -->
				<div class="space-y-2">
					{#if tls.certResolver}
						<div class="flex items-start gap-2 rounded-md bg-blue-50 p-2 dark:bg-blue-950/30">
							<Key class="mt-0.5 h-4 w-4 shrink-0 text-blue-500" />
							<div class="min-w-0 flex-1">
								<div class="text-sm font-medium text-blue-900 dark:text-blue-100">
									Certificate Resolver
								</div>
								<div class="truncate text-xs text-blue-700 dark:text-blue-300">
									{truncateText(tls.certResolver, 20)}
								</div>
								<div class="mt-1 text-xs text-blue-600 dark:text-blue-400">
									Automatic SSL certificate management
								</div>
							</div>
						</div>
					{/if}

					{#if isTCPTLSConfig(tls) && tls.passthrough}
						<div class="flex items-start gap-2 rounded-md bg-green-50 p-2 dark:bg-green-950/30">
							<Globe class="mt-0.5 h-4 w-4 shrink-0 text-green-500" />
							<div class="min-w-0 flex-1">
								<div class="text-sm font-medium text-green-900 dark:text-green-100">
									TLS Passthrough
								</div>
								<div class="text-xs text-green-700 dark:text-green-300">
									TLS termination handled by backend
								</div>
							</div>
						</div>
					{/if}

					{#if tls.options}
						<div class="flex items-start gap-2 rounded-md bg-purple-50 p-2 dark:bg-purple-950/30">
							<Eye class="mt-0.5 h-4 w-4 shrink-0 text-purple-500" />
							<div class="min-w-0 flex-1">
								<div class="text-sm font-medium text-purple-900 dark:text-purple-100">
									TLS Options
								</div>
								<div class="truncate text-xs text-purple-700 dark:text-purple-300">
									{tls.options}
								</div>
								<div class="mt-1 text-xs text-purple-600 dark:text-purple-400">
									Custom TLS settings applied
								</div>
							</div>
						</div>
					{/if}

					{#if tls.domains?.length}
						<div class="flex items-start gap-2 rounded-md bg-orange-50 p-2 dark:bg-orange-950/30">
							<Globe class="mt-0.5 h-4 w-4 shrink-0 text-orange-500" />
							<div class="min-w-0 flex-1">
								<div class="text-sm font-medium text-orange-900 dark:text-orange-100">
									Certificate Domains ({tls.domains.length})
								</div>
								<div class="mt-1 space-y-1">
									{#each tls.domains.slice(0, 3) as domain (domain.main)}
										<div class="text-xs">
											<span class="font-medium text-orange-800 dark:text-orange-200">
												{domain.main}
											</span>
											{#if domain.sans?.length}
												<span class="text-orange-600 dark:text-orange-400">
													(+{domain.sans.length} SAN{domain.sans.length > 1 ? 's' : ''})
												</span>
											{/if}
										</div>
									{/each}
									{#if tls.domains.length > 3}
										<div class="text-xs text-orange-600 dark:text-orange-400">
											+{tls.domains.length - 3} more domains...
										</div>
									{/if}
								</div>
							</div>
						</div>
					{/if}

					<!-- Basic TLS warning -->
					{#if !tls.certResolver && !(isTCPTLSConfig(tls) && tls.passthrough) && !tls.options && !tls.domains?.length}
						<div class="flex items-start gap-2 rounded-md bg-yellow-50 p-2 dark:bg-yellow-950/30">
							<AlertTriangle class="mt-0.5 h-4 w-4 shrink-0 text-yellow-500" />
							<div class="min-w-0 flex-1">
								<div class="text-sm font-medium text-yellow-900 dark:text-yellow-100">
									Basic TLS Enabled
								</div>
								<div class="text-xs text-yellow-700 dark:text-yellow-300">
									TLS is enabled but has minimal configuration. Consider setting up a certificate
									resolver or defining specific domains.
								</div>
							</div>
						</div>
					{/if}
				</div>

				<!-- Quick Actions -->
				<!-- {#if !compact} -->
				<!-- 	<div class="border-border flex gap-2 border-t pt-2"> -->
				<!-- 		<Tooltip.Provider> -->
				<!-- 			<Tooltip.Root delayDuration={300}> -->
				<!-- 				<Tooltip.Trigger> -->
				<!-- 					<Badge -->
				<!-- 						variant="outline" -->
				<!-- 						class="hover:bg-muted cursor-pointer text-xs" -->
				<!-- 						onclick={() => column?.setFilterValue?.('resolver')} -->
				<!-- 					> -->
				<!-- 						Filter by resolver -->
				<!-- 					</Badge> -->
				<!-- 				</Tooltip.Trigger> -->
				<!-- 				<Tooltip.Content side="bottom"> -->
				<!-- 					Show only routers using certificate resolvers -->
				<!-- 				</Tooltip.Content> -->
				<!-- 			</Tooltip.Root> -->
				<!-- 		</Tooltip.Provider> -->
				<!-- 	</div> -->
				<!-- {/if} -->
			</div>
		</HoverCard.Content>
	</HoverCard.Root>
{:else}
	<!-- Disabled TLS -->
	<Tooltip.Provider>
		<Tooltip.Root delayDuration={300}>
			<Tooltip.Trigger>
				<Badge
					variant="outline"
					class="cursor-pointer transition-colors duration-200
						   {compact ? 'px-2 text-xs' : 'text-sm'}"
					onclick={() => column?.setFilterValue?.('disabled')}
				>
					<Unlock class="text-muted-foreground mr-1 h-3 w-3" />
					<span class="text-muted-foreground">Disabled</span>
				</Badge>
			</Tooltip.Trigger>
			<Tooltip.Content side="top">
				<div class="text-center">
					<div class="font-medium">TLS Disabled</div>
					<div class="text-muted-foreground text-xs">This router does not use TLS encryption</div>
				</div>
			</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
{/if}
