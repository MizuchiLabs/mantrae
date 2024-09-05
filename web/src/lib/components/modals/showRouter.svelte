<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { config, provider, toggleDNSProvider } from '$lib/api';
	import { newService, type Router } from '$lib/types/config';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import type { Selected } from 'bits-ui';
	import Service from '../forms/service.svelte';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let router: Router;
	let service = $config?.services?.[router.service + '@' + router.provider] ?? newService();

	const getSelectedDNSProvider = (router: Router): Selected<unknown> | undefined => {
		return router?.dnsProvider
			? { value: router.dnsProvider, label: router.dnsProvider }
			: undefined;
	};
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<Button variant="ghost" class="h-8 w-4 rounded-full bg-green-400">
			<iconify-icon icon="fa6-solid:eye" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[520px]">
		<Tabs.Root value="router" class="mt-4 w-[470px]">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="router">Router</Tabs.Trigger>
				<Tabs.Trigger value="service">Service</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="router">
				<Card.Root>
					<Card.Header>
						<Card.Title class="flex items-center justify-between gap-1">
							<span>Router</span>
							<div>
								<Badge variant="secondary" class="bg-blue-400">
									Type: {router.routerType}
								</Badge>
								<Badge variant="secondary" class="bg-green-400">
									Provider: {router.provider}
								</Badge>
							</div>
						</Card.Title>
						<Card.Description>
							Make changes to your Router here. Click save when you're done.
						</Card.Description>
					</Card.Header>
					<Card.Content class="space-y-4">
						<div class="grid grid-cols-4 items-center gap-4">
							<Label for="name" class="text-right">Name</Label>
							<Input
								id="name"
								name="name"
								type="text"
								class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
								value={router.service}
								placeholder="Name of the router"
								disabled
							/>
						</div>
						<ArrayInput
							items={router.entrypoints}
							label="Entrypoints"
							placeholder=""
							disabled={true}
						/>
						<ArrayInput
							items={router.middlewares}
							label="Middlewares"
							placeholder=""
							disabled={true}
						/>
						{#if $provider}
							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="provider" class="text-right">DNS Provider</Label>
								<Select.Root
									selected={getSelectedDNSProvider(router)}
									onSelectedChange={(value) => toggleDNSProvider(router, value)}
								>
									<Select.Trigger class="col-span-3">
										<Select.Value placeholder="Select a dns provider" />
									</Select.Trigger>
									<Select.Content>
										<Select.Item value="" label="">None</Select.Item>

										{#each $provider as provider}
											<Select.Item value={provider.name} class="flex items-center gap-2">
												{provider.name} ({provider.type})
												{#if provider.is_active}
													<iconify-icon icon="fa6-solid:star" class="text-yellow-400" />
												{/if}
											</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>
							</div>
						{/if}
						<div class:hidden={router.routerType === 'udp'}>
							<RuleEditor bind:rule={router.rule} disabled={true} />
						</div>
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
			<Tabs.Content value="service">
				<Service bind:service disabled={true} />
			</Tabs.Content>
		</Tabs.Root>
		<Dialog.Close class="w-full">
			<Button class="w-full">Close</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
