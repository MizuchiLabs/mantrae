<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
	import { Search, User, Bot, TriangleAlert } from '@lucide/svelte';
	import type { AuditLog } from '$lib/gen/mantrae/v1/auditlog_pb';
	import { audit } from '$lib/api/util.svelte';

	interface Props {
		open?: boolean;
	}
	let { open = $bindable(false) }: Props = $props();

	const logs = audit.logs(100n);
	let searchQuery = $state('');

	function timeAgo(date: Timestamp) {
		const dateTime = new Date(timestampDate(date));
		const seconds = Math.floor((new Date().getTime() - dateTime.getTime()) / 1000);

		if (seconds < 60) return `${seconds} second${seconds !== 1 ? 's' : ''} ago`;

		const intervals = [
			{ label: 'year', seconds: 31536000 },
			{ label: 'month', seconds: 2592000 },
			{ label: 'day', seconds: 86400 },
			{ label: 'hour', seconds: 3600 },
			{ label: 'minute', seconds: 60 }
		];

		for (const interval of intervals) {
			const count = Math.floor(seconds / interval.seconds);
			if (count >= 1) {
				return `${count} ${interval.label}${count !== 1 ? 's' : ''} ago`;
			}
		}
	}
	function getActivityColor(log: AuditLog) {
		if (log.agentId) return 'bg-blue-500';
		if (log.userId) return 'bg-green-500';
		return 'bg-orange-500';
	}

	function getActivityBadgeColor(log: AuditLog) {
		if (log.agentId) return 'bg-blue-100 text-blue-700 border-blue-200';
		if (log.userId) return 'bg-green-100 text-green-700 border-green-200';
		return 'bg-orange-100 text-orange-700 border-orange-200';
	}

	let filteredLogs = $derived(
		logs.data?.filter(
			(log) =>
				log.details?.toLowerCase().includes(searchQuery.toLowerCase()) ||
				log.agentName?.toLowerCase().includes(searchQuery.toLowerCase()) ||
				log.userName?.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);
</script>

<Dialog.Root bind:open>
	<Dialog.Content
		class="no-scrollbar max-h-[90vh] w-fit max-w-[90vw] overflow-y-auto px-4 py-2 sm:min-w-200"
	>
		<Dialog.Header class="space-y-3 py-2">
			<div class="flex items-center justify-between">
				<div class="space-y-1">
					<Dialog.Title class="text-lg font-semibold">Audit Logs</Dialog.Title>
					<Dialog.Description>
						System activity and security events across all users and agents
					</Dialog.Description>
				</div>
			</div>

			<!-- Search -->
			<div class="space-y-2">
				<Label for="search" class="text-sm font-medium">Search Logs</Label>
				<div class="relative">
					<Search class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
					<Input
						id="search"
						bind:value={searchQuery}
						placeholder="Search by activity, user, or agent..."
						class="pl-10"
					/>
				</div>
			</div>
		</Dialog.Header>

		<Separator />

		<div class="space-y-1 py-2">
			{#if filteredLogs === undefined}
				<div class="flex items-center justify-center py-12">
					<div class="space-y-2 text-center">
						<TriangleAlert class="mx-auto h-8 w-8 text-muted-foreground" />
						<p class="text-sm text-muted-foreground">
							{searchQuery ? 'No logs match your search criteria' : 'No audit logs found'}
						</p>
					</div>
				</div>
			{:else}
				{#each filteredLogs as log (log.id)}
					<div class="group rounded-lg border p-4 transition-colors hover:bg-muted/30">
						<div class="flex items-start gap-3">
							<!-- Activity Indicator -->
							<div class="flex flex-col items-center gap-1 pt-1">
								<div class="h-2 w-2 rounded-full {getActivityColor(log)}"></div>
								<div class="h-8 w-px bg-border"></div>
							</div>

							<div class="flex-1 space-y-2">
								<!-- Activity Details -->
								<div class="space-y-1">
									<p class="text-sm leading-relaxed font-medium">{log.details}</p>

									<div class="flex items-center gap-2 text-xs text-muted-foreground">
										{#if log.createdAt}
											<span>{timeAgo(log.createdAt)}</span>
										{/if}

										{#if log.agentId || log.userId}
											<span>•</span>
										{/if}

										{#if log.agentId}
											<div class="flex items-center gap-1">
												<Bot class="h-3 w-3" />
												<Badge
													variant="outline"
													class="h-5 px-2 text-xs font-medium {getActivityBadgeColor(log)}"
													title={log.agentId}
												>
													{log.agentName || `...${log.agentId.slice(-8)}`}
												</Badge>
											</div>
										{:else if log.userId}
											<div class="flex items-center gap-1">
												<User class="h-3 w-3" />
												<Badge
													variant="outline"
													class="h-5 px-2 text-xs font-medium {getActivityBadgeColor(log)}"
													title={log.userId}
												>
													{log.userName || `...${log.userId.slice(-8)}`}
												</Badge>
											</div>
										{/if}
									</div>
								</div>
							</div>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
