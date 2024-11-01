<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
</script>

<ArrayInput
	bind:items={middleware.content.users}
	label="Users"
	placeholder="user:password"
	helpText="Username and password are separated by a colon. Password will be hashed automatically. You will not be able to see the password again!"
	{disabled}
/>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="users-file" class="text-right">Users File</Label>
	<Input
		id="users-file"
		name="users-file"
		type="text"
		bind:value={middleware.content.usersFile}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="/path/to/my/usersfile"
		{disabled}
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="realm" class="text-right">Realm</Label>
	<Input
		id="realm"
		name="realm"
		type="text"
		bind:value={middleware.content.realm}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="traefik"
		{disabled}
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="remove-header" class="text-right">Remove Header</Label>
	<Switch
		id="remove-header"
		bind:checked={middleware.content.removeHeader}
		class="col-span-3"
		{disabled}
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="header-field" class="text-right">Header Field</Label>
	<Input
		id="header-field"
		name="header-field"
		type="text"
		bind:value={middleware.content.headerField}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="X-WebAuth-User"
		{disabled}
	/>
</div>
