import { writable } from 'svelte/store';
import { API_URL } from '../config';

export const targets = writable<string[]>([]);

export async function fetchTargets() {
    const response = await fetch(`${API_URL}/metrics/data`);
    const data = await response.json();
    targets.set(Object.keys(data));
}

export async function fetchMetrics() {
    const response = await fetch(`${API_URL}/metrics/data`);
    return await response.json();
}
