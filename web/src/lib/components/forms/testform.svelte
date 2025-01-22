<script lang="ts">
	import { Field, Label, FieldErrors, Control, Description, Fieldset, Legend } from 'formsnap';
	import { superForm } from 'sveltekit-superforms';
	import type { Writable } from 'svelte/store';

	export let schema;
	export let data;
	export let formData: Writable<unknown>;

	const { form, enhance } = superForm(data, { validators: schema });
</script>

<form method="POST" use:enhance>
	{#each Object.keys(schema.shape) as fieldName}
		{#if schema.shape[fieldName]._def.typeName === 'ZodString'}
			<Field {form} name={fieldName}>
				<Control>
					<Label>{fieldName}</Label>
					<input type="text" bind:value={$formData[fieldName]} />
				</Control>
				<FieldErrors />
			</Field>
		{/if}
		{#if schema.shape[fieldName]._def.typeName === 'ZodEnum'}
			<Field {form} name={fieldName}>
				<Control>
					<Label>{fieldName}</Label>
					<select bind:value={$formData[fieldName]}>
						{#each schema.shape[fieldName]._def.values as value}
							<option {value}>{value}</option>
						{/each}
					</select>
				</Control>
				<FieldErrors />
			</Field>
		{/if}
		{#if schema.shape[fieldName]._def.typeName === 'ZodBoolean'}
			<Field {form} name={fieldName}>
				<Control>
					<Label>{fieldName}</Label>
					<input type="checkbox" bind:checked={$formData[fieldName]} />
				</Control>
				<FieldErrors />
			</Field>
		{/if}
	{/each}

	<button type="submit">Submit</button>
</form>
