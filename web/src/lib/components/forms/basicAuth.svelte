<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: HttpMiddleware;
	middleware.basicAuth = middleware.basicAuth ?? {
		users: [''],
		realm: '',
		removeHeader: false,
		headerField: ''
	};
</script>

<ArrayInput bind:items={middleware.basicAuth.users} label="Users" />
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="users-file" class="text-right">Users File</Label>
	<Input
		id="users-file"
		name="users-file"
		type="text"
		bind:value={middleware.basicAuth.usersFile}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Users File"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="realm" class="text-right">Realm</Label>
	<Input
		id="realm"
		name="realm"
		type="text"
		bind:value={middleware.basicAuth.realm}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Realm"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="remove-header" class="text-right">Remove Header</Label>
	<Switch
		id="remove-header"
		bind:checked={middleware.basicAuth.removeHeader}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="header-field" class="text-right">Header Field</Label>
	<Input
		id="header-field"
		name="header-field"
		type="text"
		bind:value={middleware.basicAuth.headerField}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Header Field"
	/>
</div>
