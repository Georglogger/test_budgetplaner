import { writable } from 'svelte/store';

// Load from localStorage
const loadDeletedBudgets = () => {
  if (typeof window === 'undefined') return [];
  const stored = localStorage.getItem('deletedBudgets');
  return stored ? JSON.parse(stored) : [];
};

// Store for all deleted budgets (persists forever)
export const deletedBudgets = writable(loadDeletedBudgets());

// Save to localStorage when changed
deletedBudgets.subscribe(value => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('deletedBudgets', JSON.stringify(value));
  }
});

export function addDeletedBudget(budget) {
  deletedBudgets.update(budgets => {
    return [{
      ...budget,
      deletedAt: new Date().toISOString()
    }, ...budgets];
  });
}

export function removeDeletedBudget(budgetId) {
  deletedBudgets.update(budgets => {
    return budgets.filter(b => b.id !== budgetId);
  });
}

export function clearAllDeleted() {
  deletedBudgets.set([]);
}
