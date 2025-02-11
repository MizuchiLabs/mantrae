import { createLocalStorage } from '$lib/storage.svelte';

const getInitialTheme = () => {
	if (typeof window === 'undefined') return 'light';

	const savedTheme = window.localStorage.getItem('theme');
	const theme =
		savedTheme === 'light' || savedTheme === 'dark'
			? savedTheme
			: window.matchMedia('(prefers-color-scheme: dark)').matches
				? 'dark'
				: 'light';

	document.documentElement.classList.remove('light', 'dark');
	document.documentElement.classList.add(theme);

	return theme;
};

class ThemeStore {
	private store = createLocalStorage<string>('theme', getInitialTheme());

	get value(): string | undefined {
		return this.store.value;
	}

	set value(value: string) {
		this.store.value = value;

		if (typeof window !== 'undefined') {
			window.localStorage.setItem('theme', value);
			document.documentElement.classList.remove('light', 'dark');
			document.documentElement.classList.add(value);
		}
	}

	toggle(): void {
		this.value = this.value === 'light' ? 'dark' : 'light';
	}
}
export const theme = new ThemeStore();
