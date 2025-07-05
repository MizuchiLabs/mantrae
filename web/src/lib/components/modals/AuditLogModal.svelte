<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { auditLogClient } from '$lib/api';
	import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
	import Separator from '../ui/separator/separator.svelte';

	interface Props {
		open?: boolean;
	}
	let { open = $bindable(false) }: Props = $props();

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
</script>

<Dialog.Root bind:open>
	<Dialog.Content
		class="max-h-[90vh] w-fit max-w-[90vw] overflow-y-auto px-4 py-2 sm:min-w-[40rem]"
	>
		<Dialog.Header class="flex justify-between gap-2 py-4">
			<Dialog.Title>Audit Logs</Dialog.Title>
		</Dialog.Header>
		<div class="space-y-3 text-sm">
			{#await auditLogClient.listAuditLogs({ limit: 1000n, offset: 0n }) then result}
				{#each result.auditLogs || [] as log (log.id)}
					<div class="flex items-center justify-start gap-3">
						{#if log.agentId}
							<div class="mt-2 h-2 w-2 rounded-full bg-blue-500"></div>
						{:else if log.userId}
							<div class="mt-2 h-2 w-2 rounded-full bg-green-500"></div>
						{:else}
							<div class="mt-2 h-2 w-2 rounded-full bg-orange-500"></div>
						{/if}
						<div>
							<p class="text-sm">{log.details}</p>

							<div class="text-muted-foreground flex items-center gap-2 text-xs">
								{#if log.createdAt}
									<span class="text-muted-foreground text-xs">
										{timeAgo(log.createdAt)}
									</span>
								{/if}
								{#if log.agentId}
									<span class="rounded bg-blue-100 px-1.5 py-0.5 text-blue-700" title={log.agentId}>
										Agent: {log.agentName || `...${log.agentId.slice(-8)}`}
									</span>
								{:else if log.userId}
									<span
										class="rounded bg-green-100 px-1.5 py-0.5 text-green-700"
										title={log.userId}
									>
										User: {log.userName || `...${log.userId.slice(-8)}`}
									</span>
								{/if}
							</div>
						</div>
					</div>
					<Separator />
				{/each}
			{/await}
		</div>
	</Dialog.Content>
</Dialog.Root>
