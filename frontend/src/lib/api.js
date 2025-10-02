const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://172.22.246.143:8080/api/v1';

export const api = {
  // Budget APIs
  async getBudgets() {
    const response = await fetch(`${API_BASE_URL}/budgets`);
    if (!response.ok) throw new Error('Fehler beim Laden der Budgets');
    return response.json();
  },

  async getBudget(id) {
    const response = await fetch(`${API_BASE_URL}/budgets/${id}`);
    if (!response.ok) throw new Error('Fehler beim Laden des Budgets');
    return response.json();
  },

  async createBudget(budget) {
    const response = await fetch(`${API_BASE_URL}/budgets`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(budget)
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || 'Fehler beim Erstellen des Budgets');
    }
    return response.json();
  },

  async updateBudget(id, budget) {
    const response = await fetch(`${API_BASE_URL}/budgets/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(budget)
    });
    if (!response.ok) throw new Error('Fehler beim Aktualisieren des Budgets');
    return response.json();
  },

  async deleteBudget(id) {
    const response = await fetch(`${API_BASE_URL}/budgets/${id}`, {
      method: 'DELETE'
    });
    if (!response.ok) throw new Error('Fehler beim Löschen des Budgets');
    return response.json();
  },

  // Budget Lines APIs
  async getBudgetLines(budgetId) {
    const response = await fetch(`${API_BASE_URL}/budgets/${budgetId}/lines`);
    if (!response.ok) throw new Error('Fehler beim Laden der Budget-Zeilen');
    return response.json();
  },

  async createBudgetLine(line) {
    const response = await fetch(`${API_BASE_URL}/budget-lines`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(line)
    });
    if (!response.ok) throw new Error('Fehler beim Erstellen der Budget-Zeile');
    return response.json();
  },

  // Actuals APIs
  async getActuals(budgetId) {
    const response = await fetch(`${API_BASE_URL}/budgets/${budgetId}/actuals`);
    if (!response.ok) throw new Error('Fehler beim Laden der Ist-Daten');
    return response.json();
  },

  async importActuals(budgetId, file) {
    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(`${API_BASE_URL}/budgets/${budgetId}/import/actuals`, {
      method: 'POST',
      body: formData
    });
    if (!response.ok) throw new Error('Fehler beim Importieren der Ist-Daten');
    return response.json();
  },

  // Reports APIs
  async getPlanActualReport(budgetId) {
    const response = await fetch(`${API_BASE_URL}/budgets/${budgetId}/reports/plan-actual`);
    if (!response.ok) throw new Error('Fehler beim Laden des Plan-Ist-Berichts');
    return response.json();
  },

  async getCategorySummary(budgetId) {
    const response = await fetch(`${API_BASE_URL}/budgets/${budgetId}/reports/category-summary`);
    if (!response.ok) throw new Error('Fehler beim Laden der Kategorien-Zusammenfassung');
    return response.json();
  }
};
