import { dev } from '$app/environment';
import { env } from '$env/dynamic/public';

// const fallbackURL = import.meta.env.PROD ? window?.location?.origin : 'http://localhost:3000';

export const BackendURL = dev ? 'http://localhost:3000' : env.PUBLIC_BACKEND_URL || ''; // Use proxy in dev mode
export const SupportEmail = env.PUBLIC_SUPPORT_EMAIL || 'support@mizuchi.dev';
export const DocsURL = env.PUBLIC_DOCS_URL || 'https://mizuchilabs.github.io/mantrae';

export const APP_NAME = 'Mantrae';
export const DEFAULT_LANGUAGE = 'en';
