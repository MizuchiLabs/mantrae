export class LocalStorage<T> {
	#key: string;
	#value = $state<T | undefined>(undefined);

	constructor(key: string, initial?: T) {
		this.#key = key;

		if (typeof localStorage !== 'undefined') {
			const stored = localStorage.getItem(key);
			this.#value = stored ? this.#parseWithBigInt(stored) : initial;

			if (stored === null && initial !== undefined) {
				localStorage.setItem(key, this.#stringifyWithBigInt(initial));
			}
		} else {
			this.#value = initial;
		}
	}

	#parseWithBigInt(text: string): T {
		try {
			return JSON.parse(text, (_, value) => {
				if (typeof value === 'string' && value.match(/^\d+n$/)) {
					return BigInt(value.slice(0, -1));
				}
				return value;
			});
		} catch {
			return text as unknown as T;
		}
	}

	#stringifyWithBigInt(value: T): string {
		return JSON.stringify(value, (_, value) => {
			// Convert BigInt to string representation with 'n' suffix
			if (typeof value === 'bigint') {
				return value.toString() + 'n';
			}
			return value;
		});
	}

	get value(): T | undefined {
		return this.#value;
	}

	set value(newValue: T | undefined) {
		this.#value = newValue;
		if (typeof localStorage !== 'undefined') {
			localStorage.setItem(this.#key, this.#stringifyWithBigInt(newValue!));
		}
	}
}

// Helper function to create a new localStorage store
export function createLocalStorage<T>(key: string, defaultValue: T) {
	return new LocalStorage<T>(key, defaultValue);
}
