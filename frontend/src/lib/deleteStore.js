import { writable } from 'svelte/store';

// Store for the last deleted budget
export const lastDeleted = writable(null);
export const showUndoToast = writable(false);

let undoTimeout = null;

export function setLastDeleted(budget) {
  lastDeleted.set(budget);
  showUndoToast.set(true);

  // Clear previous timeout
  if (undoTimeout) {
    clearTimeout(undoTimeout);
  }

  // Auto-hide undo toast after 10 seconds
  undoTimeout = setTimeout(() => {
    showUndoToast.set(false);
    lastDeleted.set(null);
  }, 10000);
}

export function clearLastDeleted() {
  lastDeleted.set(null);
  showUndoToast.set(false);
  if (undoTimeout) {
    clearTimeout(undoTimeout);
  }
}
