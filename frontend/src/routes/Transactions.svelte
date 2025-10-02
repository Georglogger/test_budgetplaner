<script>
  import { onMount } from 'svelte';

  let searchQuery = '';
  let activeTab = 'uncategorized'; // 'uncategorized', 'bank', 'open-items', 'budgets', 'reconciled', 'late', 'all'
  let showFilters = false;
  let selectedFilters = {
    similarTitles: false,
    differentAmounts: false
  };

  let transactions = [];
  let loading = false;

  const tabs = [
    { id: 'uncategorized', label: 'Zu kategorisieren', icon: '🏷️' },
    { id: 'bank', label: 'Bank', icon: '🏦' },
    { id: 'open-items', label: 'Offene Posten', icon: '📦' },
    { id: 'budgets', label: 'Budgets', icon: '💼' },
    { id: 'reconciled', label: 'Abgleichen', icon: '✓' },
    { id: 'late', label: 'Spät', icon: '⏰' },
    { id: 'all', label: 'Alle', icon: '📋' }
  ];

  onMount(() => {
    loadTransactions();
  });

  function loadTransactions() {
    loading = true;
    // TODO: Load from API
    setTimeout(() => {
      transactions = [];
      loading = false;
    }, 500);
  }

  function toggleFilters() {
    showFilters = !showFilters;
  }

  function resetFilters() {
    selectedFilters = {
      similarTitles: false,
      differentAmounts: false
    };
    searchQuery = '';
  }
</script>

<div class="transactions-page">
  <!-- Search and Filter Bar -->
  <div class="search-bar">
    <div class="search-input-wrapper">
      <svg class="search-icon" width="18" height="18" viewBox="0 0 18 18" fill="none">
        <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5"/>
        <path d="M12.5 12.5L16 16" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
      <input
        type="text"
        placeholder="Suchen"
        bind:value={searchQuery}
        class="search-input"
      />
    </div>

    <button class="filter-btn" class:active={showFilters} on:click={toggleFilters}>
      <svg width="18" height="18" viewBox="0 0 18 18" fill="currentColor">
        <path d="M2 4h14M4 9h10M6 14h6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
      Filtern
    </button>
  </div>

  <!-- Tab Navigation -->
  <div class="tabs-container">
    {#each tabs as tab}
      <button
        class="tab"
        class:active={activeTab === tab.id}
        on:click={() => activeTab = tab.id}
      >
        <span class="tab-icon">{tab.icon}</span>
        {tab.label}
      </button>
    {/each}
  </div>

  <!-- Filter Options -->
  {#if showFilters}
    <div class="filter-options">
      <label class="filter-checkbox">
        <input type="checkbox" bind:checked={selectedFilters.similarTitles} />
        <span class="checkbox-label">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M13 4L6 11L3 8" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </span>
        <span class="filter-text">Transaktionen mit ähnlichen Titeln zusammenfassen</span>
      </label>

      <label class="filter-checkbox">
        <input type="checkbox" bind:checked={selectedFilters.differentAmounts} />
        <span class="checkbox-label">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M13 4L6 11L3 8" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </span>
        <span class="filter-text">Transaktionen mit unterschiedlichen Beträgen zusammenfassen</span>
      </label>
    </div>
  {/if}

  <!-- Action Bar -->
  <div class="action-bar">
    <button class="action-btn">
      <svg width="18" height="18" viewBox="0 0 18 18" fill="currentColor">
        <path d="M9 2v14M2 9h14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
      Aktion
    </button>

    <div class="action-checkboxes">
      <label class="action-checkbox">
        <input type="checkbox" bind:checked={selectedFilters.similarTitles} />
        Transaktionen mit ähnlichen Titeln zusammenfassen
      </label>

      <label class="action-checkbox">
        <input type="checkbox" bind:checked={selectedFilters.differentAmounts} />
        Transaktionen mit unterschiedlichen Beträgen zusammenfassen
      </label>
    </div>

    <div class="info-pill">
      <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <circle cx="8" cy="8" r="7" fill="#7c3aed"/>
        <path d="M8 7v4M8 5h.01" stroke="white" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <span>Wähle mehrere Transaktionen, um Aktionen auszuführen oder nutze Kategorie-Zuweisen-Button</span>
    </div>

    <button class="category-btn">
      Regeln verwalten
    </button>
  </div>

  <!-- Table Header -->
  <div class="table-header">
    <div class="table-columns">
      <div class="column-title">Titel</div>
      <div class="column-value">Wert</div>
      <div class="column-source">Datenquelle</div>
    </div>
  </div>

  <!-- Empty State -->
  {#if !loading && transactions.length === 0}
    <div class="empty-state">
      <div class="empty-illustration">
        <svg width="200" height="200" viewBox="0 0 200 200">
          <!-- Simple illustration placeholder -->
          <circle cx="100" cy="80" r="30" fill="#e5e7eb"/>
          <rect x="70" y="120" width="60" height="60" rx="4" fill="#e5e7eb"/>
          <line x1="60" y1="140" x2="140" y2="140" stroke="#cbd5e1" stroke-width="2"/>
          <line x1="70" y1="155" x2="130" y2="155" stroke="#cbd5e1" stroke-width="2"/>
          <line x1="75" y1="170" x2="125" y2="170" stroke="#cbd5e1" stroke-width="2"/>
        </svg>
      </div>

      <h2 class="empty-title">Keine Vorschläge gefunden</h2>

      <p class="empty-description">
        Keine Transaktionen gefunden, für die es Vorschläge zur Kategorisierung gibt.
      </p>

      <div class="empty-reasons">
        <p>Das kann folgende Gründe haben:</p>
        <ul>
          <li>Es gibt keine Transaktionen mit Vorschlägen</li>
          <li>Alle Datenquellen sind ausgeschaltet (siehe Overlay im Header)</li>
          <li>Es existieren keine Transaktionen, weil keine Datenquellen vorhanden sind</li>
        </ul>
      </div>

      <button class="reset-btn" on:click={resetFilters}>Reset Filters</button>
    </div>
  {/if}

  {#if loading}
    <div class="loading">Lade Transaktionen...</div>
  {/if}
</div>

<style>
  .transactions-page {
    max-width: 1400px;
    margin: 0 auto;
  }

  /* Search Bar */
  .search-bar {
    display: flex;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .search-input-wrapper {
    flex: 1;
    position: relative;
  }

  .search-icon {
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    color: #94a3b8;
  }

  .search-input {
    width: 100%;
    padding: 0.75rem 1rem 0.75rem 3rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    background: white;
    transition: all 0.2s ease;
  }

  .search-input:focus {
    outline: none;
    border-color: #7c3aed;
    box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.1);
  }

  :global(.dark) .search-input {
    background: #1e293b;
    border-color: #334155;
    color: #e2e8f0;
  }

  .filter-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.5rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: #64748b;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .filter-btn:hover {
    background: #f9fafb;
  }

  .filter-btn.active {
    background: #7c3aed;
    border-color: #7c3aed;
    color: white;
  }

  :global(.dark) .filter-btn {
    background: #1e293b;
    border-color: #334155;
    color: #94a3b8;
  }

  :global(.dark) .filter-btn:hover {
    background: #334155;
  }

  /* Tabs */
  .tabs-container {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
    padding: 0.5rem;
    background: white;
    border-radius: 0.75rem;
    border: 1px solid #e5e7eb;
    overflow-x: auto;
  }

  :global(.dark) .tabs-container {
    background: #1e293b;
    border-color: #334155;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.625rem 1rem;
    background: transparent;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: #64748b;
    cursor: pointer;
    transition: all 0.2s ease;
    white-space: nowrap;
  }

  .tab:hover {
    background: #f9fafb;
    color: #1a202c;
  }

  .tab.active {
    background: #7c3aed;
    color: white;
  }

  :global(.dark) .tab {
    color: #94a3b8;
  }

  :global(.dark) .tab:hover {
    background: #334155;
    color: #e2e8f0;
  }

  :global(.dark) .tab.active {
    background: #7c3aed;
    color: white;
  }

  .tab-icon {
    font-size: 1rem;
  }

  /* Filter Options */
  .filter-options {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    margin-bottom: 1.5rem;
  }

  :global(.dark) .filter-options {
    background: #1e293b;
    border-color: #334155;
  }

  .filter-checkbox {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
  }

  .filter-checkbox input[type="checkbox"] {
    display: none;
  }

  .checkbox-label {
    width: 1.25rem;
    height: 1.25rem;
    border: 2px solid #e5e7eb;
    border-radius: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .checkbox-label svg {
    opacity: 0;
    transition: opacity 0.2s ease;
  }

  .filter-checkbox input[type="checkbox"]:checked + .checkbox-label {
    background: #7c3aed;
    border-color: #7c3aed;
  }

  .filter-checkbox input[type="checkbox"]:checked + .checkbox-label svg {
    opacity: 1;
    stroke: white;
  }

  .filter-text {
    font-size: 0.875rem;
    color: #64748b;
  }

  :global(.dark) .filter-text {
    color: #94a3b8;
  }

  /* Action Bar */
  .action-bar {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    margin-bottom: 1rem;
  }

  :global(.dark) .action-bar {
    background: #1e293b;
    border-color: #334155;
  }

  .action-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: transparent;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    color: #64748b;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .action-btn:hover {
    background: #f9fafb;
  }

  :global(.dark) .action-btn {
    border-color: #334155;
    color: #94a3b8;
  }

  .action-checkboxes {
    display: flex;
    gap: 1rem;
    flex: 1;
  }

  .action-checkbox {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.813rem;
    color: #64748b;
    cursor: pointer;
  }

  :global(.dark) .action-checkbox {
    color: #94a3b8;
  }

  .info-pill {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: #ede9fe;
    border-radius: 0.5rem;
    font-size: 0.813rem;
    color: #5b21b6;
  }

  :global(.dark) .info-pill {
    background: rgba(124, 58, 237, 0.2);
    color: #c4b5fd;
  }

  .category-btn {
    padding: 0.5rem 1rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: #64748b;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .category-btn:hover {
    background: #f9fafb;
  }

  :global(.dark) .category-btn {
    background: #334155;
    border-color: #475569;
    color: #94a3b8;
  }

  /* Table Header */
  .table-header {
    background: white;
    border: 1px solid #e5e7eb;
    border-bottom: none;
    border-radius: 0.75rem 0.75rem 0 0;
    padding: 1rem;
  }

  :global(.dark) .table-header {
    background: #1e293b;
    border-color: #334155;
  }

  .table-columns {
    display: grid;
    grid-template-columns: 1fr 200px 200px;
    gap: 1rem;
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.025em;
    color: #64748b;
  }

  /* Empty State */
  .empty-state {
    background: white;
    border: 1px solid #e5e7eb;
    border-top: none;
    border-radius: 0 0 0.75rem 0.75rem;
    padding: 4rem 2rem;
    text-align: center;
  }

  :global(.dark) .empty-state {
    background: #1e293b;
    border-color: #334155;
  }

  .empty-illustration {
    margin: 0 auto 2rem;
    display: flex;
    justify-content: center;
  }

  .empty-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a202c;
    margin-bottom: 1rem;
  }

  :global(.dark) .empty-title {
    color: #e2e8f0;
  }

  .empty-description {
    font-size: 1rem;
    color: #64748b;
    margin-bottom: 2rem;
  }

  :global(.dark) .empty-description {
    color: #94a3b8;
  }

  .empty-reasons {
    max-width: 600px;
    margin: 0 auto 2rem;
    text-align: left;
  }

  .empty-reasons p {
    font-size: 0.875rem;
    color: #64748b;
    margin-bottom: 0.5rem;
  }

  .empty-reasons ul {
    list-style: disc;
    padding-left: 1.5rem;
    color: #64748b;
    font-size: 0.875rem;
  }

  .empty-reasons li {
    margin-bottom: 0.25rem;
  }

  :global(.dark) .empty-reasons p,
  :global(.dark) .empty-reasons ul {
    color: #94a3b8;
  }

  .reset-btn {
    padding: 0.75rem 2rem;
    background: #7c3aed;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .reset-btn:hover {
    background: #6d28d9;
  }

  .loading {
    padding: 4rem 2rem;
    text-align: center;
    color: #64748b;
  }
</style>
