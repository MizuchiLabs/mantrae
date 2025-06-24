<script lang="ts">
	import { agentClient } from '$lib/api';
	import AgentModal from '$lib/components/modals/agent.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { BulkAction } from '$lib/components/tables/types';
	import { renderComponent } from '$lib/components/ui/data-table';
	import type { Agent } from '$lib/gen/mantrae/v1/agent_pb';
	import { DateFormat, pageIndex, pageSize } from '$lib/stores/common';
	import { profile } from '$lib/stores/profile';
	import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
	import { ConnectError } from '@connectrpc/connect';
	import { Bot, KeyRound, Pencil, Trash } from '@lucide/svelte';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { readable } from 'svelte/store';

	let item = $state({} as Agent);
	let open = $state(false);

	// Data state
	let data = $state<Agent[]>([]);
	let rowCount = $state<number>(0);

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
			enableSorting: true
		},
		{
			header: 'Last Seen',
			accessorKey: 'updatedAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('updatedAt') as Timestamp;
				return DateFormat.format(timestampDate(date));
			}
		},
		{
			header: 'Created',
			accessorKey: 'createdAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('createdAt') as Timestamp;
				return DateFormat.format(timestampDate(date));
			}
		},
		{
			id: 'actions',
			enableHiding: false,
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
							type: 'button',
							label: 'Delete Agent',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteItem(row.original.id)
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
		// return lastSeenInSeconds <= 30 ? 'bg-green-500/10' : 'bg-red-500/10';
		return lastSeenInSeconds <= 30 ? true : false;
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

	async function onPaginationChange(p: PaginationState) {
		await refreshData(p.pageSize, p.pageIndex);
	}

	const deleteItem = async (id: string) => {
		try {
			await agentClient.deleteAgent({ id: id });
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success('Agent deleted');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete router', { description: e.message });
		}
	};

	async function bulkDelete(selectedRows: Agent[]) {
		try {
			const confirmed = confirm(`Are you sure you want to delete ${selectedRows.length} routers?`);
			if (!confirmed) return;

			const rows = selectedRows.map((row) => ({ id: row.id }));
			for (const row of rows) {
				await agentClient.deleteAgent({ id: row.id });
			}
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			toast.success(`Successfully deleted ${selectedRows.length} agents`);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete agents', { description: e.message });
		}
	}

	async function refreshData(pageSize: number, pageIndex: number) {
		const response = await agentClient.listAgents({
			profileId: profile.id,
			limit: BigInt(pageSize),
			offset: BigInt(pageIndex * pageSize)
		});
		data = response.agents;
		rowCount = Number(response.totalCount);
	}

	async function createAgent() {
		try {
			const response = await agentClient.createAgent({ profileId: profile.id });
			if (!response.agent) return;

			toast.success('Agent created');
			await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
			item = response.agent;
			open = true;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to create agent', { description: e.message });
		}
	}

	onMount(async () => {
		await refreshData(pageSize.value ?? 10, pageIndex.value ?? 0);
	});
</script>

<svelte:head>
	<title>Agents</title>
</svelte:head>

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-start gap-2">
		<Bot />
		<h1 class="text-2xl font-bold">Agents</h1>
	</div>
	<DataTable
		{data}
		{columns}
		{rowCount}
		{onPaginationChange}
		{bulkActions}
		rowClassModifiers={{
			'bg-red-50': (r) => !getAgentStatus(r),
			'bg-green-50': (r) => getAgentStatus(r)
		}}
		createButton={{
			label: 'Add Agent',
			onClick: createAgent
		}}
	/>
</div>

<AgentModal bind:open bind:item bind:data />
