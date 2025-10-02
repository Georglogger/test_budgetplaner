<script>
  import { api } from '../lib/api';

  export let budgetId = null;

  let formData = {
    name: '',
    description: '',
    period_start: '',
    period_end: '',
    status: 'draft',
    created_by: 'user@example.com'
  };

  let lines = [];
  let newLine = {
    category: '',
    subcategory: '',
    amount: 0,
    driver: '',
    driver_value: 0
  };

  let loading = false;
  let error = null;
  let success = false;

  async function handleSubmit() {
    loading = true;
    error = null;

    try {
      // Konvertiere Datumswerte in ISO 8601 Format
      const budgetData = {
        ...formData,
        period_start: new Date(formData.period_start).toISOString(),
        period_end: new Date(formData.period_end).toISOString()
      };

      console.log('Sending budget data:', JSON.stringify(budgetData, null, 2));
      const budget = await api.createBudget(budgetData);

      // Budget-Zeilen speichern
      for (const line of lines) {
        await api.createBudgetLine({
          ...line,
          budget_id: budget.id
        });
      }

      success = true;
      setTimeout(() => {
        window.location.hash = `#/budgets/${budget.id}`;
      }, 1500);
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function addLine() {
    if (newLine.category && newLine.amount) {
      lines = [...lines, { ...newLine }];
      newLine = {
        category: '',
        subcategory: '',
        amount: 0,
        driver: '',
        driver_value: 0
      };
    }
  }

  function removeLine(index) {
    lines = lines.filter((_, i) => i !== index);
  }

  function calculateTotal() {
    return lines.reduce((sum, line) => sum + Number(line.amount), 0);
  }
</script>

<div class="container">
  <div class="header">
    <h1>Neues Budget erstellen</h1>
    <button class="btn-secondary" on:click={() => window.location.hash = '#/'}>
      Zurück
    </button>
  </div>

  {#if success}
    <div class="success">Budget erfolgreich erstellt!</div>
  {/if}

  {#if error}
    <div class="error">{error}</div>
  {/if}

  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-section">
      <h2>Budget-Details</h2>

      <div class="form-group">
        <label for="name">Budget-Name *</label>
        <input
          id="name"
          type="text"
          bind:value={formData.name}
          required
          placeholder="z.B. Jahresbudget 2025"
        />
      </div>

      <div class="form-group">
        <label for="description">Beschreibung</label>
        <textarea
          id="description"
          bind:value={formData.description}
          rows="3"
          placeholder="Optionale Beschreibung des Budgets"
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="period_start">Start-Datum *</label>
          <input
            id="period_start"
            type="date"
            bind:value={formData.period_start}
            required
          />
        </div>

        <div class="form-group">
          <label for="period_end">End-Datum *</label>
          <input
            id="period_end"
            type="date"
            bind:value={formData.period_end}
            required
          />
        </div>
      </div>

      <div class="form-group">
        <label for="status">Status</label>
        <select id="status" bind:value={formData.status}>
          <option value="draft">Entwurf</option>
          <option value="approved">Genehmigt</option>
          <option value="closed">Abgeschlossen</option>
        </select>
      </div>
    </div>

    <div class="form-section">
      <h2>Budget-Zeilen</h2>

      <div class="line-form">
        <div class="form-row">
          <div class="form-group">
            <label for="category">Kategorie *</label>
            <input
              id="category"
              type="text"
              bind:value={newLine.category}
              placeholder="z.B. Personal"
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
            <button type="button" class="btn-secondary" on:click={addLine}>
              Zeile hinzufügen
            </button>
          </div>
        </div>
      </div>

      {#if lines.length > 0}
        <table class="lines-table">
          <thead>
            <tr>
              <th>Kategorie</th>
              <th>Unterkategorie</th>
              <th>Betrag</th>
              <th>Treiber</th>
              <th>Treiber-Wert</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {#each lines as line, i}
              <tr>
                <td>{line.category}</td>
                <td>{line.subcategory || '-'}</td>
                <td class="amount">{Number(line.amount).toFixed(2)} €</td>
                <td>{line.driver || '-'}</td>
                <td>{line.driver_value || '-'}</td>
                <td>
                  <button type="button" class="btn-delete" on:click={() => removeLine(i)}>
                    Löschen
                  </button>
                </td>
              </tr>
            {/each}
            <tr class="total-row">
              <td colspan="2"><strong>Gesamt</strong></td>
              <td class="amount"><strong>{calculateTotal().toFixed(2)} €</strong></td>
              <td colspan="3"></td>
            </tr>
          </tbody>
        </table>
      {/if}
    </div>

    <div class="form-actions">
      <button type="submit" class="btn-primary" disabled={loading}>
        {loading ? 'Wird gespeichert...' : 'Budget erstellen'}
      </button>
      <button type="button" class="btn-secondary" on:click={() => window.location.hash = '#/'}>
        Abbrechen
      </button>
    </div>
  </form>
</div>

<style>
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  h1 {
    font-size: 2rem;
    font-weight: bold;
    color: #1a202c;
  }

  h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a202c;
    margin-bottom: 1.5rem;
  }

  .form-section {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    padding: 2rem;
    margin-bottom: 2rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
    flex: 1;
  }

  .form-row {
    display: flex;
    gap: 1rem;
  }

  label {
    display: block;
    font-weight: 500;
    color: #374151;
    margin-bottom: 0.5rem;
  }

  input, textarea, select {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-family: inherit;
  }

  input:focus, textarea:focus, select:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .btn-primary, .btn-secondary, .btn-delete {
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

  .btn-primary:hover:not(:disabled) {
    background-color: #2563eb;
  }

  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-secondary {
    background-color: #f3f4f6;
    color: #374151;
  }

  .btn-secondary:hover {
    background-color: #e5e7eb;
  }

  .btn-delete {
    background-color: #fee2e2;
    color: #991b1b;
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
  }

  .btn-delete:hover {
    background-color: #fecaca;
  }

  .line-form {
    background-color: #f9fafb;
    padding: 1.5rem;
    border-radius: 0.5rem;
    margin-bottom: 1.5rem;
  }

  .lines-table {
    width: 100%;
    border-collapse: collapse;
  }

  .lines-table th {
    background-color: #f9fafb;
    padding: 0.75rem;
    text-align: left;
    font-weight: 600;
    color: #374151;
    border-bottom: 2px solid #e5e7eb;
  }

  .lines-table td {
    padding: 0.75rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .lines-table .amount {
    text-align: right;
    font-family: 'Courier New', monospace;
  }

  .total-row {
    background-color: #f9fafb;
    font-weight: 600;
  }

  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }

  .error {
    background-color: #fee2e2;
    color: #991b1b;
    padding: 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
  }

  .success {
    background-color: #d1fae5;
    color: #065f46;
    padding: 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
  }
</style>
