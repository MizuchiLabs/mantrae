<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { entrypoints, middlewares, getService } from '$lib/api';
	import { type Router } from '$lib/types/config';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import type { Selected } from 'bits-ui';
	import Service from '../forms/service.svelte';

	export let router: Router;
	let service = getService(router.service + '@' + router.provider);

	const getSelectedEntrypoints = (router: Router): Selected<unknown>[] => {
		let list = router?.entrypoints?.map((entrypoint) => {
			return { value: entrypoint, label: entrypoint };
		});
		return list ?? [];
	};
	const getSelectedMiddlewares = (router: Router): Selected<unknown>[] => {
		let list = router?.middlewares?.map((middleware) => {
			return { value: middleware, label: middleware };
		});
		return list ?? [];
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
					<Card.Content class="space-y-2">
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
						<div class="grid grid-cols-4 items-center gap-4">
							<Label for="entrypoints" class="text-right">Entrypoints</Label>
							<Select.Root multiple={true} selected={getSelectedEntrypoints(router)}>
								<Select.Trigger class="col-span-3">
									<Select.Value placeholder="Select an entrypoint" />
								</Select.Trigger>
								<Select.Content>
									{#each $entrypoints || [] as entrypoint}
										<Select.Item value={entrypoint.name}>
											<div class="flex flex-row items-center gap-2">
												{entrypoint.name}
												{#if entrypoint.http}
													{#if 'tls' in entrypoint.http}
														<iconify-icon icon="fa6-solid:lock" class=" text-green-400" />
													{/if}
												{/if}
											</div>
										</Select.Item>
									{/each}
								</Select.Content>
							</Select.Root>
						</div>
						<div
							class="grid grid-cols-4 items-center gap-4"
							class:hidden={router.routerType === 'udp'}
						>
							<Label for="middlewares" class="text-right">Middlewares</Label>
							<Select.Root multiple={true} selected={getSelectedMiddlewares(router)}>
								<Select.Trigger class="col-span-3">
									<Select.Value placeholder="Select a middleware" />
								</Select.Trigger>
								<Select.Content>
									{#each $middlewares as middleware}
										{#if router.routerType === middleware.middlewareType}
											<Select.Item value={middleware.name}>
												{middleware.name}
											</Select.Item>
										{/if}
									{/each}
								</Select.Content>
							</Select.Root>
						</div>
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
