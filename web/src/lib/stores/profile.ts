import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
import { createLocalStorage } from '$lib/storage.svelte';

class ProfileStore {
	private store = createLocalStorage<Profile | null>('selected_profile', null);

	get value(): Profile | null {
		return this.store.value ?? null;
	}

	set value(profile: Profile | null) {
		this.store.value = profile;
	}

	// Helper methods for safe access
	get id(): bigint | undefined {
		return this.value?.id;
	}

	get name(): string | undefined {
		return this.value?.name;
	}

	get description(): string | undefined {
		return this.value?.description;
	}

	get token(): string | undefined {
		return this.value?.token;
	}

	// Validation methods
	hasValidId(): boolean {
		return typeof this.id === 'bigint' && this.id > 0n;
	}

	hasValidName(): boolean {
		return typeof this.name === 'string' && this.name.length > 0;
	}

	isValid(): boolean {
		return this.value !== null && this.hasValidId() && this.hasValidName();
	}

	clear(): void {
		this.store.value = null;
	}
}

export const profile = new ProfileStore();
