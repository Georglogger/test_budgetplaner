<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { darkMode, toggleDarkMode } from '../lib/darkMode';
  import DeleteButton from '../lib/DeleteButton.svelte';

  export let budgetId;

  let budget = null;
  let budgetLines = [];
  let report = [];
  let categorySummary = [];
  let loading = true;
  let error = null;
  let showAddLine = false;
  let newLine = {
    category: '',
    subcategory: '',
    amount: 0,
    driver: '',
    driver_value: 0
  };

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    try {
      loading = true;
      [budget, budgetLines, report, categorySummary] = await Promise.all([
        api.getBudget(budgetId),
        api.getBudgetLines(budgetId),
        api.getPlanActualReport(budgetId),
        api.getCategorySummary(budgetId)
      ]);
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  async function addBudgetLine() {
    try {
      await api.createBudgetLine({
        ...newLine,
        budget_id: budgetId
      });

      // Reset form
      newLine = {
        category: '',
        subcategory: '',
        amount: 0,
        driver: '',
        driver_value: 0
      };
      showAddLine = false;

      // Reload data
      await loadData();
    } catch (e) {
      error = e.message;
    }
  }

  function calculateBudgetTotal() {
    return budgetLines.reduce((sum, line) => sum + Number(line.amount), 0);
  }

  function formatCurrency(amount) {
    return new Intl.NumberFormat('de-DE', {
      style: 'currency',
      currency: 'EUR'
    }).format(amount);
  }

  function formatPercent(value) {
    return `${value >= 0 ? '+' : ''}${value.toFixed(1)}%`;
  }

  function getVarianceClass(variance) {
    if (Math.abs(variance) < 0.05) return 'neutral';
    return variance > 0 ? 'negative' : 'positive';
  }

  function calculateTotals() {
    return categorySummary.reduce((acc, cat) => ({
      planned: acc.planned + cat.planned,
      actual: acc.actual + cat.actual,
      variance: acc.variance + cat.variance
    }), { planned: 0, actual: 0, variance: 0 });
  }
</script>

<div class="container">
  {#if loading}
    <p>Lade Dashboard...</p>
  {:else if error}
    <div class="error">{error}</div>
  {:else}
    <div class="header">
      <div>
        <h1>{budget?.name}</h1>
        <p class="subtitle">{budget?.description || 'Keine Beschreibung'}</p>
      </div>
      <div class="actions">
        <button class="btn-dark-toggle" on:click={toggleDarkMode} title="Dark Mode umschalten">
          {#if $darkMode}
            ☀️
          {:else}
            🌙
          {/if}
        </button>
        <button class="btn-secondary" on:click={() => window.location.hash = '#/'}>
          Zurück zur Übersicht
        </button>
        <button class="btn-secondary" on:click={() => window.location.hash = `#/budgets/${budgetId}/table`}>
          📊 Tabellenansicht
        </button>
        <button class="btn-primary" on:click={() => window.location.hash = `#/budgets/${budgetId}/import`}>
          Ist-Daten importieren
        </button>
      </div>
    </div>

    <!-- Budget-Zeilen -->
    <div class="budget-lines-section">
      <div class="section-header">
        <h2>Budget-Zeilen</h2>
        <button class="btn-primary" on:click={() => showAddLine = !showAddLine}>
          {showAddLine ? 'Abbrechen' : '+ Neue Zeile'}
        </button>
      </div>

      {#if showAddLine}
        <div class="add-line-form">
          <div class="form-row">
            <div class="form-group">
              <label for="category">Kategorie *</label>
              <input
                id="category"
                type="text"
                bind:value={newLine.category}
                placeholder="z.B. Personal"
                required
              />
            </div>

            <div class="form-group">
              <label for="subcategory">Unterkategorie</label>
              <input
                id="subcategory"
                type="text"
                bind:value={newLine.subcategory}
                placeholder="z.B. Gehälter"
              />
            </div>

            <div class="form-group">
              <label for="amount">Betrag *</label>
              <input
                id="amount"
                type="number"
                step="0.01"
                bind:value={newLine.amount}
                placeholder="0.00"
                required
              />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="driver">Treiber</label>
              <input
                id="driver"
                type="text"
                bind:value={newLine.driver}
                placeholder="z.B. anzahl_mitarbeiter"
              />
            </div>

            <div class="form-group">
              <label for="driver_value">Treiber-Wert</label>
              <input
                id="driver_value"
                type="number"
                step="0.01"
                bind:value={newLine.driver_value}
                placeholder="0"
              />
            </div>

            <div class="form-group">
              <label>&nbsp;</label>
              <button type="button" class="btn-primary" on:click={addBudgetLine}>
                Speichern
              </button>
            </div>
          </div>
        </div>
      {/if}

      {#if budgetLines.length === 0}
        <p class="empty-state">Noch keine Budget-Zeilen vorhanden. Klicke auf "+ Neue Zeile" um eine hinzuzufügen.</p>
      {:else}
        <div class="table-container">
          <table class="budget-lines-table">
            <thead>
              <tr>
                <th>Kategorie</th>
                <th>Unterkategorie</th>
                <th class="number">Betrag</th>
                <th>Treiber</th>
                <th class="number">Treiber-Wert</th>
              </tr>
            </thead>
            <tbody>
              {#each budgetLines as line}
                <tr>
                  <td>{line.category}</td>
                  <td>{line.subcategory || '-'}</td>
                  <td class="number">{formatCurrency(line.amount)}</td>
                  <td>{line.driver || '-'}</td>
                  <td class="number">{line.driver_value || '-'}</td>
                </tr>
              {/each}
              <tr class="total-row">
                <td colspan="2"><strong>Gesamt</strong></td>
                <td class="number"><strong>{formatCurrency(calculateBudgetTotal())}</strong></td>
                <td colspan="2"></td>
              </tr>
            </tbody>
          </table>
        </div>
      {/if}
    </div>

    <!-- Kategorien-Zusammenfassung -->
    <div class="summary-section">
      <h2>Kategorien-Übersicht</h2>
      <div class="cards-grid">
        {#each categorySummary as category}
          <div class="summary-card">
            <h3>{category.category}</h3>
            <div class="summary-values">
              <div class="value-row">
                <span class="label">Plan:</span>
                <span class="value">{formatCurrency(category.planned)}</span>
              </div>
              <div class="value-row">
                <span class="label">Ist:</span>
                <span class="value">{formatCurrency(category.actual)}</span>
              </div>
              <div class="value-row variance {getVarianceClass(category.variance / category.planned)}">
                <span class="label">Abweichung:</span>
                <span class="value">{formatCurrency(category.variance)}</span>
              </div>
            </div>
          </div>
        {/each}

        <!-- Gesamt-Karte -->
        {#if categorySummary.length > 0}
          {@const totals = calculateTotals()}
          <div class="summary-card total-card">
            <h3>Gesamt</h3>
            <div class="summary-values">
              <div class="value-row">
                <span class="label">Plan:</span>
                <span class="value">{formatCurrency(totals.planned)}</span>
              </div>
              <div class="value-row">
                <span class="label">Ist:</span>
                <span class="value">{formatCurrency(totals.actual)}</span>
              </div>
              <div class="value-row variance {getVarianceClass(totals.variance / totals.planned)}">
                <span class="label">Abweichung:</span>
                <span class="value">{formatCurrency(totals.variance)}</span>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>

    <!-- Detaillierter Plan-Ist-Vergleich -->
    <div class="report-section">
      <h2>Detaillierter Plan-Ist-Vergleich</h2>
      {#if report.length === 0}
        <p class="empty-state">Noch keine Daten vorhanden.</p>
      {:else}
        <div class="table-container">
          <table class="report-table">
            <thead>
              <tr>
                <th>Kategorie</th>
                <th>Unterkategorie</th>
                <th class="number">Plan</th>
                <th class="number">Ist</th>
                <th class="number">Abweichung</th>
                <th class="number">Abw. %</th>
              </tr>
            </thead>
            <tbody>
              {#each report as row}
                <tr>
                  <td>{row.category}</td>
                  <td>{row.subcategory || '-'}</td>
                  <td class="number">{formatCurrency(row.planned)}</td>
                  <td class="number">{formatCurrency(row.actual)}</td>
                  <td class="number variance {getVarianceClass(row.variance_pct / 100)}">
                    {formatCurrency(row.variance)}
                  </td>
                  <td class="number variance {getVarianceClass(row.variance_pct / 100)}">
                    {formatPercent(row.variance_pct)}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 2rem;
  }

  h1 {
    font-size: 2rem;
    font-weight: bold;
    color: #1a202c;
    margin: 0;
  }

  .subtitle {
    color: #6b7280;
    margin-top: 0.5rem;
  }

  .actions {
    display: flex;
    gap: 1rem;
  }

  h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a202c;
    margin-bottom: 1.5rem;
  }

  .summary-section, .report-section {
    margin-bottom: 3rem;
  }

  .cards-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1.5rem;
  }

  .summary-card {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .total-card {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
  }

  .total-card h3,
  .total-card .label,
  .total-card .value {
    color: white;
  }

  .summary-card h3 {
    font-size: 1.125rem;
    font-weight: 600;
    margin-bottom: 1rem;
    color: #1a202c;
  }

  .summary-values {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .value-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .label {
    font-size: 0.875rem;
    color: #6b7280;
  }

  .value {
    font-size: 1.125rem;
    font-weight: 600;
    color: #1a202c;
    font-family: 'Courier New', monospace;
  }

  .variance {
    padding: 0.5rem;
    border-radius: 0.5rem;
  }

  .variance.positive {
    background-color: #d1fae5;
  }

  .variance.positive .value {
    color: #065f46;
  }

  .variance.negative {
    background-color: #fee2e2;
  }

  .variance.negative .value {
    color: #991b1b;
  }

  .variance.neutral {
    background-color: #f3f4f6;
  }

  .table-container {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    overflow: hidden;
  }

  .report-table {
    width: 100%;
    border-collapse: collapse;
  }

  .report-table th {
    background-color: #f9fafb;
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: #374151;
    border-bottom: 2px solid #e5e7eb;
  }

  .report-table td {
    padding: 1rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .report-table .number {
    text-align: right;
    font-family: 'Courier New', monospace;
  }

  .btn-primary, .btn-secondary {
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    font-weight: 500;
    transition: all 0.2s;
  }

  .btn-primary {
    background-color: #3b82f6;
    color: white;
  }

  .btn-primary:hover {
    background-color: #2563eb;
  }

  .btn-secondary {
    background-color: #f3f4f6;
    color: #374151;
  }

  .btn-secondary:hover {
    background-color: #e5e7eb;
  }

  .error {
    background-color: #fee2e2;
    color: #991b1b;
    padding: 1rem;
    border-radius: 0.5rem;
    margin: 1rem 0;
  }

  .empty-state {
    text-align: center;
    color: #6b7280;
    padding: 3rem;
    background: white;
    border-radius: 0.75rem;
    border: 1px dashed #d1d5db;
  }

  .budget-lines-section {
    margin-bottom: 3rem;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .add-line-form {
    background-color: #f9fafb;
    padding: 1.5rem;
    border-radius: 0.75rem;
    margin-bottom: 1.5rem;
    border: 1px solid #e5e7eb;
  }

  .form-row {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .form-row:last-child {
    margin-bottom: 0;
  }

  .form-group {
    flex: 1;
  }

  .form-group label {
    display: block;
    font-weight: 500;
    color: #374151;
    margin-bottom: 0.5rem;
    font-size: 0.875rem;
  }

  .form-group input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-family: inherit;
  }

  .form-group input:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .budget-lines-table {
    width: 100%;
    border-collapse: collapse;
  }

  .budget-lines-table th {
    background-color: #f9fafb;
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: #374151;
    border-bottom: 2px solid #e5e7eb;
  }

  .budget-lines-table td {
    padding: 1rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .budget-lines-table .number {
    text-align: right;
    font-family: 'Courier New', monospace;
  }

  .total-row {
    background-color: #f9fafb;
    font-weight: 600;
  }

  /* Dark Mode Styles */
  :global(.dark) .container {
    color: #e2e8f0;
  }

  :global(.dark) h1,
  :global(.dark) h2 {
    color: #e2e8f0;
  }

  :global(.dark) .subtitle {
    color: #94a3b8;
  }

  .btn-dark-toggle {
    padding: 0.5rem 0.75rem;
    font-size: 1.25rem;
    background-color: #334155;
    border: 1px solid #475569;
  }

  .btn-dark-toggle:hover {
    background-color: #475569;
  }

  :global(.dark) .btn-secondary {
    background-color: #334155;
    border-color: #475569;
    color: #e2e8f0;
  }

  :global(.dark) .btn-secondary:hover {
    background-color: #475569;
  }

  :global(.dark) .summary-card {
    background: #1e293b;
    border-color: #334155;
  }

  :global(.dark) .summary-card h3 {
    color: #e2e8f0;
  }

  :global(.dark) .label {
    color: #94a3b8;
  }

  :global(.dark) .value {
    color: #e2e8f0;
  }

  :global(.dark) .table-container {
    background: #1e293b;
    border-color: #334155;
  }

  :global(.dark) .report-table th {
    background-color: #0f172a;
    color: #cbd5e1;
    border-bottom-color: #475569;
  }

  :global(.dark) .report-table td {
    border-bottom-color: #334155;
    color: #e2e8f0;
  }

  :global(.dark) .empty-state {
    color: #94a3b8;
    border-color: #475569;
  }

  :global(.dark) .add-line-form {
    background-color: #0f172a;
    border-color: #334155;
  }

  :global(.dark) .form-group input {
    background: #1e293b;
    border-color: #475569;
    color: #e2e8f0;
  }

  :global(.dark) .budget-lines-table {
    background: #1e293b;
  }

  :global(.dark) .budget-lines-table th {
    background-color: #0f172a;
    color: #cbd5e1;
    border-bottom-color: #475569;
  }

  :global(.dark) .budget-lines-table td {
    border-bottom-color: #334155;
    color: #e2e8f0;
  }

  :global(.dark) .total-row {
    background-color: #0f172a;
  }
</style>

<DeleteButton {budgetId} />
