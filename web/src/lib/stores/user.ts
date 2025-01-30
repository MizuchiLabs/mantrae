import type { User } from '$lib/types';
import { createLocalStorage } from '$lib/storage.svelte';

class UserStore {
	private store = createLocalStorage<User | null>('current_user', null);

	get value(): User | null {
		return this.store.value ?? null;
	}

	set value(user: User | null) {
		this.store.value = user;
	}

	// Safe getters for required fields
	get id(): number | undefined {
		return this.value?.id;
	}

	get username(): string | undefined {
		return this.value?.username;
	}

	get isAdmin(): boolean {
		return this.value?.isAdmin ?? false;
	}

	// Optional field getters
	get email(): string | undefined {
		return this.value?.email;
	}

	get lastLogin(): string | undefined {
		return this.value?.lastLogin;
	}

	get createdAt(): string | undefined {
		return this.value?.createdAt;
	}

	get updatedAt(): string | undefined {
		return this.value?.updatedAt;
	}

	// Validation methods
	hasValidId(): boolean {
		return typeof this.id === 'number' && !isNaN(this.id);
	}

	hasValidUsername(): boolean {
		return typeof this.username === 'string' && this.username.length > 0;
	}

	isLoggedIn(): boolean {
		return this.value !== null && this.hasValidId() && this.hasValidUsername();
	}

	hasEmail(): boolean {
		return typeof this.email === 'string' && this.email.length > 0;
	}

	// Authorization methods
	canAccessAdmin(): boolean {
		return this.isLoggedIn() && this.isAdmin;
	}

	// Clear user data (logout)
	clear(): void {
		this.store.value = null;
	}
}

export const user = new UserStore();
