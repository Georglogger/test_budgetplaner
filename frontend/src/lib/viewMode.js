import { writable } from 'svelte/store';

// Load from localStorage or default to 'employee'
const stored = typeof window !== 'undefined' ? localStorage.getItem('viewMode') : null;
export const viewMode = writable(stored || 'employee'); // 'employee' or 'customer'

// Save to localStorage when changed
viewMode.subscribe(value => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('viewMode', value);
  }
});

// Password for switching (you should change this!)
const SWITCH_PASSWORD = 'TEST';

export function switchViewMode(password) {
  if (password === SWITCH_PASSWORD) {
    viewMode.update(mode => mode === 'employee' ? 'customer' : 'employee');
    return true;
  }
  return false;
}

export function isCustomerView() {
  let current;
  viewMode.subscribe(value => current = value)();
  return current === 'customer';
}

export function isEmployeeView() {
  let current;
  viewMode.subscribe(value => current = value)();
  return current === 'employee';
}
