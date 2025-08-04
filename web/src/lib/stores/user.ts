import type { User } from '$lib/gen/mantrae/v1/user_pb';
import { createLocalStorage } from '$lib/storage.svelte';
import type { Timestamp } from '@bufbuild/protobuf/wkt';

class UserStore {
	private store = createLocalStorage<User | null>('current_user', null);

	get value(): User | null {
		return this.store.value ?? null;
	}

	set value(user: User | null) {
		this.store.value = user;
	}

	get id(): string | undefined {
		return this.value?.id;
	}

	get username(): string | undefined {
		return this.value?.username;
	}

	get email(): string | undefined {
		return this.value?.email;
	}

	get lastLogin(): Timestamp | undefined {
		return this.value?.lastLogin;
	}

	get createdAt(): Timestamp | undefined {
		return this.value?.createdAt;
	}

	get updatedAt(): Timestamp | undefined {
		return this.value?.updatedAt;
	}

	// Validation methods
	hasValidId(): boolean {
		return typeof this.id === 'string' && this.id.length > 0;
	}

	hasValidUsername(): boolean {
		return typeof this.username === 'string' && this.username.length > 0;
	}

	isLoggedIn(): boolean {
		return this.value !== null && this.hasValidId();
	}

	hasEmail(): boolean {
		return typeof this.email === 'string' && this.email.length > 0;
	}

	// Clear user data (logout)
	clear(): void {
		this.store.value = null;
	}
}

export const user = new UserStore();
