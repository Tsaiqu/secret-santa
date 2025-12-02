import { writable } from 'svelte/store';

// Pusty magazyn, który wypełnimy danymi pobranymi z Go
export const participantsStore = writable([]);
