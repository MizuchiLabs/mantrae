<script lang="ts">
	import { agentClient } from '$lib/api';
	import AgentModal from '$lib/components/modals/AgentModal.svelte';
	import ColumnBadge from '$lib/components/tables/ColumnBadge.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import type { Agent } from '$lib/gen/mantrae/v1/agent_pb';
	import { DateFormat } from '$lib/stores/common';
	import { profile } from '$lib/stores/profile';
	import { agents } from '$lib/stores/realtime';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { ConnectError } from '@connectrpc/connect';
	import { Bot, KeyRound, Pencil, Trash } from '@lucide/svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { toast } from 'svelte-sonner';
	import { readable } from 'svelte/store';

	let item = $state({} as Agent);
	let open = $state(false);

	const columns: ColumnDef<Agent>[] = [
		{
			header: 'Hostname',
			accessorKey: 'hostname',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('hostname') as string;
				if (!name) {
					return 'Connect your agent!';
				}
				return name;
			}
		},
		{
			header: 'Endpoint',
			accessorKey: 'activeIp',
			enableSorting: true,
			cell: ({ row }) => {
				return renderComponent(ColumnBadge, {
					label: row.original.activeIp || 'Unknown',
					class: 'hover:cursor-pointer'
				});
			}
		},
		{
			header: 'Last Seen',
			accessorKey: 'updatedAt',
			enableSorting: true,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				if (row.original.updatedAt === undefined || !row.original.hostname) {
					return renderComponent(ColumnBadge, { label: 'Never', class: 'text-xs' });
				}
				return DateFormat.format(timestampDate(row.original.updatedAt));
			}
		},
		{
			id: 'actions',
			enableHiding: false,
			enableGlobalFilter: false,
			cell: ({ row }) => {
				let editText = row.original.hostname ? 'Edit Agent' : 'Connect Agent';
				let editIcon = row.original.hostname ? Pencil : KeyRound;
				return renderComponent(TableActions, {
					actions: [
						{
							type: 'button',
							label: editText,
							icon: editIcon,
							onClick: () => {
								item = row.original;
								open = true;
							}
						},
						{
							type: 'popover',
							label: 'Delete Agent',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original),
							popover: {
								title: 'Delete Agent?',
								description:
									'This agent will will be permanently deleted. This will also delete all associated routers.',
								confirmLabel: 'Delete',
								cancelLabel: 'Cancel'
							}
						}
					]
				});
			}
		}
	];

	const now = readable(new Date(), (set) => {
		const interval = setInterval(() => {
			set(new Date());
		}, 1000);
		return () => clearInterval(interval);
	});
	function getAgentStatus(agent: Agent) {
		if (!agent.updatedAt) return false;
		const lastSeen = new Date(timestampDate(agent.updatedAt));
		const lastSeenInSeconds = ($now.getTime() - lastSeen.getTime()) / 1000;
		return lastSeenInSeconds <= 20 ? true : false;
	}

	const bulkActions: BulkAction<Agent>[] = [
		{
			type: 'button',
			label: 'Delete',
			icon: Trash,
			variant: 'destructive',
			onClick: bulkDelete
		}
	];

	const deleteItem = async (item: Agent) => {
		try {
			await agentClient.deleteAgent({ id: item.id });
			toast.success(`Agent ${item.hostname ?? item.id} deleted`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete agent', { description: e.message });
		}
	};

	async function bulkDelete(rows: Agent[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${rows.length} agents?`);
			if (!confirmed) return;

			for (const row of rows) {
				await agentClient.deleteAgent({ id: row.id });
			}
			toast.success(`Successfully deleted ${rows.length} agents`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete agents', { description: e.message });
		}
	}
	async function createAgent() {
		try {
			const response = await agentClient.createAgent({ profileId: profile.id });
			if (!response.agent) return;
			item = response.agent;

			toast.success('Agent created');
			open = true;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to create agent', { description: e.message });
		}
	}

	$effect(() => {
		if (profile.isValid()) {
			agentClient.listAgents({ profileId: profile.id }).then((response) => {
				agents.set(response.agents);
			});
		}
	});
</script>

<svelte:head>
	<title>Agents</title>
</svelte:head>

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="bg-primary/10 rounded-lg p-2">
					<Bot class="text-primary h-6 w-6" />
				</div>
				Agents
			</h1>
			<p class="text-muted-foreground mt-1">Connect your agents</p>
		</div>
	</div>

	<DataTable
		data={$agents}
		{columns}
		{bulkActions}
		rowClassModifiers={{
			'bg-red-300/25 dark:bg-red-700/25': (r) => !getAgentStatus(r)
		}}
		createButton={{
			label: 'Add Agent',
			onClick: createAgent
		}}
	/>
</div>

<AgentModal bind:open bind:item />
