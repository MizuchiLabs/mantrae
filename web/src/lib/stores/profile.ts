import { createLocalStorage } from '$lib/storage.svelte';
import type { Profile } from '$lib/types';

class ProfileStore {
	private store = createLocalStorage<Profile | null>('selected_profile', null);

	get value(): Profile | null {
		return this.store.value ?? null;
	}

	set value(profile: Profile | null) {
		this.store.value = profile;
	}

	// Helper methods for safe access
	get id(): number | undefined {
		return this.value?.id;
	}

	get name(): string | undefined {
		return this.value?.name;
	}

	// Validation methods
	hasValidId(): boolean {
		return typeof this.id === 'number' && !isNaN(this.id);
	}

	hasValidName(): boolean {
		return typeof this.name === 'string' && this.name.length > 0;
	}

	isValid(): boolean {
		return this.value !== null && this.hasValidId() && this.hasValidName();
	}
}

export const profile = new ProfileStore();
