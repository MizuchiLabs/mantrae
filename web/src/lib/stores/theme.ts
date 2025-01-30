import { writable } from 'svelte/store';

// Light/Dark Mode
const getInitialTheme = () => {
	if (typeof window !== 'undefined') {
		const savedTheme = window.localStorage.getItem('theme') as string;
		if (savedTheme === 'light' || savedTheme === 'dark') return savedTheme;

		return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
	}
	return 'light';
};

export const theme = writable<'light' | 'dark'>(getInitialTheme());

// Subscribe to changes and update localStorage and document class
if (typeof window !== 'undefined') {
	theme.subscribe((value) => {
		window.localStorage.setItem('theme', value);
		document.documentElement.classList.remove('light', 'dark');
		document.documentElement.classList.add(value);
	});
}
