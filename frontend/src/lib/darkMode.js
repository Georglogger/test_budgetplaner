import { writable } from 'svelte/store';

// Check if user has a preference saved
const storedTheme = localStorage.getItem('theme');
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

// Initialize with stored preference or system preference
const initialTheme = storedTheme || (prefersDark ? 'dark' : 'light');

export const darkMode = writable(initialTheme === 'dark');

// Update localStorage when theme changes
darkMode.subscribe(isDark => {
  localStorage.setItem('theme', isDark ? 'dark' : 'light');

  if (isDark) {
    document.documentElement.classList.add('dark');
    document.body.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
    document.body.classList.remove('dark');
  }
});

export function toggleDarkMode() {
  darkMode.update(value => !value);
}
