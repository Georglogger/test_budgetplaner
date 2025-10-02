<script>
  import { api } from '../lib/api';

  export let budgetId;

  let file = null;
  let uploading = false;
  let result = null;
  let error = null;

  function handleFileChange(event) {
    file = event.target.files[0];
    result = null;
    error = null;
  }

  async function handleUpload() {
    if (!file) {
      error = 'Bitte wählen Sie eine Datei aus';
      return;
    }

    uploading = true;
    error = null;

    try {
      result = await api.importActuals(budgetId, file);
    } catch (e) {
      error = e.message;
    } finally {
      uploading = false;
    }
  }

  function downloadTemplate() {
    const csv = 'category,subcategory,amount,date\nPersonal,Gehälter,50000.00,2025-01-01\nMarketing,Online-Werbung,10000.00,2025-01-15\n';
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'actuals_template.csv';
    a.click();
    URL.revokeObjectURL(url);
  }
</script>

<div class="container">
  <div class="header">
    <h1>Ist-Daten importieren</h1>
    <button class="btn-secondary" on:click={() => window.location.hash = `#/budgets/${budgetId}`}>
      Zurück zum Dashboard
    </button>
  </div>

  <div class="content">
    <div class="info-box">
      <h3>📋 CSV-Format</h3>
      <p>Die CSV-Datei muss folgende Spalten enthalten:</p>
      <ul>
        <li><strong>category</strong>: Kategorie (z.B. "Personal")</li>
        <li><strong>subcategory</strong>: Unterkategorie (z.B. "Gehälter")</li>
        <li><strong>amount</strong>: Betrag als Dezimalzahl (z.B. 50000.00)</li>
        <li><strong>date</strong>: Datum im Format YYYY-MM-DD (z.B. 2025-01-01)</li>
      </ul>
      <button class="btn-link" on:click={downloadTemplate}>
        📥 Beispiel-Vorlage herunterladen
      </button>
    </div>

    <div class="upload-section">
      <div class="upload-box">
        <input
          type="file"
          accept=".csv"
          on:change={handleFileChange}
          id="file-input"
        />
        <label for="file-input" class="file-label">
          {#if file}
            <span class="file-name">📄 {file.name}</span>
          {:else}
            <span class="placeholder">Klicken Sie hier, um eine CSV-Datei auszuwählen</span>
          {/if}
        </label>
      </div>

      <button
        class="btn-primary"
        on:click={handleUpload}
        disabled={!file || uploading}
      >
        {uploading ? 'Importiere...' : 'Hochladen & Importieren'}
      </button>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    {#if result}
      <div class="result-box">
        <h3>Import-Ergebnis</h3>
        <div class="result-stats">
          <div class="stat success">
            <span class="stat-label">Erfolgreich importiert:</span>
            <span class="stat-value">{result.imported}</span>
          </div>
          <div class="stat total">
            <span class="stat-label">Gesamt verarbeitet:</span>
            <span class="stat-value">{result.total}</span>
          </div>
        </div>

        {#if result.errors && result.errors.length > 0}
          <div class="errors-section">
            <h4>Fehler beim Import:</h4>
            <ul class="error-list">
              {#each result.errors as errorMsg}
                <li>{errorMsg}</li>
              {/each}
            </ul>
          </div>
        {:else}
          <div class="success-message">
            ✅ Alle Daten wurden erfolgreich importiert!
          </div>
        {/if}

        <button class="btn-primary" on:click={() => window.location.hash = `#/budgets/${budgetId}`}>
          Zum Dashboard
        </button>
      </div>
    {/if}
  </div>
</div>

<style>
  .container {
    max-width: 900px;
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

  .content {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .info-box {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    padding: 2rem;
  }

  .info-box h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1a202c;
    margin-bottom: 1rem;
  }

  .info-box p {
    color: #6b7280;
    margin-bottom: 1rem;
  }

  .info-box ul {
    list-style: none;
    padding: 0;
    margin: 1rem 0;
  }

  .info-box li {
    padding: 0.5rem 0;
    color: #374151;
  }

  .btn-link {
    background: none;
    border: none;
    color: #3b82f6;
    cursor: pointer;
    font-size: 1rem;
    text-decoration: underline;
    padding: 0;
    margin-top: 1rem;
  }

  .btn-link:hover {
    color: #2563eb;
  }

  .upload-section {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .upload-box {
    position: relative;
  }

  #file-input {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
  }

  .file-label {
    display: block;
    padding: 3rem;
    border: 2px dashed #d1d5db;
    border-radius: 0.75rem;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s;
  }

  .file-label:hover {
    border-color: #3b82f6;
    background-color: #f9fafb;
  }

  .placeholder {
    color: #6b7280;
    font-size: 1.125rem;
  }

  .file-name {
    color: #1a202c;
    font-size: 1.125rem;
    font-weight: 500;
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
    align-self: flex-end;
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

  .error {
    background-color: #fee2e2;
    color: #991b1b;
    padding: 1rem;
    border-radius: 0.5rem;
  }

  .result-box {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    padding: 2rem;
  }

  .result-box h3 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a202c;
    margin-bottom: 1.5rem;
  }

  .result-stats {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .stat {
    padding: 1.5rem;
    border-radius: 0.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .stat.success {
    background-color: #d1fae5;
  }

  .stat.total {
    background-color: #dbeafe;
  }

  .stat-label {
    color: #6b7280;
    font-size: 0.875rem;
  }

  .stat-value {
    font-size: 2rem;
    font-weight: bold;
    color: #1a202c;
  }

  .success-message {
    background-color: #d1fae5;
    color: #065f46;
    padding: 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1.5rem;
    font-weight: 500;
  }

  .errors-section {
    margin-bottom: 1.5rem;
  }

  .errors-section h4 {
    color: #991b1b;
    margin-bottom: 0.75rem;
  }

  .error-list {
    background-color: #fee2e2;
    padding: 1rem;
    border-radius: 0.5rem;
    list-style: none;
  }

  .error-list li {
    color: #991b1b;
    padding: 0.25rem 0;
    font-size: 0.875rem;
  }
</style>
