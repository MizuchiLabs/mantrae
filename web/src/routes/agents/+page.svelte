<script lang="ts">
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import TableActions from '$lib/components/tables/TableActions.svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { Bot, KeyRound, Pencil, Trash } from 'lucide-svelte';
	import { type Agent } from '$lib/types';
	import AgentModal from '$lib/components/modals/agent.svelte';
	import { api, agents, profile } from '$lib/api';
	import { renderComponent } from '$lib/components/ui/data-table';
	import { toast } from 'svelte-sonner';
	import { DateFormat } from '$lib/store';

	interface ModalState {
		isOpen: boolean;
		agent?: Agent;
	}

	const initialModalState: ModalState = { isOpen: false };
	let modalState = $state(initialModalState);

	const deleteAgent = async (agent: Agent) => {
		try {
			await api.deleteAgent(agent.id);
			toast.success('Agent deleted');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

	const copyToken = async (agent: Agent) => {
		try {
			await navigator.clipboard.writeText(agent.token);
			toast.success('Token copied to clipboard');
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(e.message);
		}
	};

	const columns: ColumnDef<Agent>[] = [
		{
			header: 'Hostname',
			accessorKey: 'hostname',
			enableSorting: true,
			cell: ({ row }) => {
				const name = row.getValue('hostname') as string;
				if (!name) {
					return 'Connect your agent!';
				} else {
					return name;
				}
			}
		},
		{
			header: 'Endpoint',
			accessorKey: 'activeIp',
			enableSorting: true
		},
		// {
		// 	header: 'Containers',
		// 	accessorKey: 'containers',
		// 	enableSorting: true,
		// 	cell: ({ row }) => {
		// 		const admin = row.getValue('isAdmin') as boolean;
		// 		if (admin) {
		// 			return renderComponent(ColumnBadge, {
		// 				label: 'Yes',
		// 				variant: 'default'
		// 			});
		// 		} else {
		// 			return renderComponent(ColumnBadge, {
		// 				label: 'No',
		// 				variant: 'secondary'
		// 			});
		// 		}
		// 	}
		// },
		{
			header: 'Last Seen',
			accessorKey: 'updatedAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('updatedAt') as string;
				return DateFormat.format(new Date(date));
			}
		},
		{
			header: 'Created',
			accessorKey: 'createdAt',
			enableSorting: true,
			cell: ({ row }) => {
				const date = row.getValue('createdAt') as string;
				return DateFormat.format(new Date(date));
			}
		},
		{
			id: 'actions',
			enableHiding: false,
			cell: ({ row }) => {
				return renderComponent(TableActions, {
					actions: [
						{
							label: 'Edit Agent',
							icon: Pencil,
							onClick: () => {
								modalState = {
									isOpen: true,
									agent: row.original
								};
							}
						},
						{
							label: 'Copy Token',
							icon: KeyRound,
							classProps: 'text-green-500',
							onClick: () => copyToken(row.original)
						},
						{
							label: 'Delete Agent',
							icon: Trash,
							classProps: 'text-destructive',
							onClick: () => deleteAgent(row.original)
						}
					]
				});
			}
		}
	];

	profile.subscribe((value) => {
		if (value.id) {
			api.listAgentsByProfile();
		}
	});
</script>

<svelte:head>
	<title>Agents</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-start gap-2">
		<Bot />
		<h1 class="text-2xl font-bold">Agent Management</h1>
	</div>
	<DataTable
		{columns}
		data={$agents || []}
		createButton={{
			label: 'Add Agent',
			onClick: () => api.createAgent()
		}}
	/>
</div>

<AgentModal bind:open={modalState.isOpen} agent={modalState.agent} />
