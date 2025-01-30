export class LocalStorage<T> {
	#key: string;
	#value = $state<T | undefined>(undefined);

	constructor(key: string, initial?: T) {
		this.#key = key;

		if (typeof localStorage !== 'undefined') {
			const stored = localStorage.getItem(key);
			this.#value = stored ? JSON.parse(stored) : initial;

			if (stored === null && initial !== undefined) {
				localStorage.setItem(key, JSON.stringify(initial));
			}
		} else {
			this.#value = initial;
		}
	}

	get value(): T | undefined {
		return this.#value;
	}

	set value(newValue: T | undefined) {
		this.#value = newValue;
		if (typeof localStorage !== 'undefined') {
			localStorage.setItem(this.#key, JSON.stringify(newValue));
		}
	}
}

// Helper function to create a new localStorage store
export function createLocalStorage<T>(key: string, defaultValue: T) {
	return new LocalStorage<T>(key, defaultValue);
}
