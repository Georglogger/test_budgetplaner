<script>
  import { onMount } from 'svelte';

  let activeTab = 'aktiv'; // 'aktiv', 'nicht-aktiv', 'archiviert', 'alle'
  let searchQuery = '';
  let showFilters = false;
  let employees = [];
  let loading = false;

  const tabs = [
    { id: 'aktiv', label: 'Aktiv', count: 0 },
    { id: 'nicht-aktiv', label: 'Nicht-Aktiv', count: 0 },
    { id: 'archiviert', label: 'Archiviert', count: 0 },
    { id: 'alle', label: 'Alle', count: 0 }
  ];

  onMount(() => {
    loadEmployees();
  });

  function loadEmployees() {
    loading = true;
    // TODO: Load from API
    setTimeout(() => {
      employees = [];
      loading = false;
    }, 500);
  }

  function toggleFilters() {
    showFilters = !showFilters;
  }
</script>

<div class="personal-page">
  <!-- Tab Navigation with Counts -->
  <div class="tabs-header">
    <div class="tabs-container">
      {#each tabs as tab}
        <button
          class="tab"
          class:active={activeTab === tab.id}
          on:click={() => activeTab = tab.id}
        >
          {tab.label}
          <span class="tab-count">{tab.count}</span>
        </button>
      {/each}
    </div>

    <div class="header-actions">
      <div class="employee-count">0/50 Employees</div>
      <button class="add-employee-btn">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 3v10M3 8h10" stroke="white" stroke-width="2" stroke-linecap="round"/>
        </svg>
        Mitarbeiter:In hinzufügen
        <span class="badge-n">N</span>
      </button>
      <button class="instruction-btn">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 12v-1M8 4v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5" fill="none"/>
        </svg>
        Anleitung
      </button>
    </div>
  </div>

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

  <!-- Submenu Card -->
  <div class="submenu-card">
    <div class="submenu-icon">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
        <path d="M16 7c0 1.657-1.343 3-3 3s-3-1.343-3-3 1.343-3 3-3 3 1.343 3 3zM12 14c-3.866 0-7 1.343-7 3v2h14v-2c0-1.657-3.134-3-7-3z" fill="#7c3aed"/>
      </svg>
    </div>
    <div class="submenu-content">
      <h3 class="submenu-title">Personal</h3>
      <p class="submenu-subtitle">Plane Personal und Freelancer</p>
    </div>
  </div>

  <!-- Empty State -->
  {#if !loading && employees.length === 0}
    <div class="empty-state">
      <div class="empty-illustration">
        <svg width="220" height="220" viewBox="0 0 220 220">
          <!-- Person with magnifying glass illustration -->
          <ellipse cx="110" cy="180" rx="80" ry="15" fill="#f3e8ff" opacity="0.5"/>

          <!-- Person head -->
          <circle cx="90" cy="80" r="25" fill="#e5e7eb"/>
          <path d="M90 75c-3 0-5-2-5-5s2-5 5-5 5 2 5 5-2 5-5 5z" fill="#9ca3af"/>

          <!-- Person body -->
          <path d="M65 105c0-13 11-24 25-24s25 11 25 24v40H65v-40z" fill="#e5e7eb"/>

          <!-- Magnifying glass -->
          <circle cx="140" cy="110" r="30" fill="none" stroke="#cbd5e1" stroke-width="4"/>
          <circle cx="140" cy="110" r="20" fill="none" stroke="#cbd5e1" stroke-width="3"/>
          <line x1="160" y1="130" x2="180" y2="150" stroke="#cbd5e1" stroke-width="6" stroke-linecap="round"/>

          <!-- Papers/documents -->
          <rect x="120" y="155" width="30" height="35" rx="2" fill="#e5e7eb"/>
          <line x1="125" y1="165" x2="145" y2="165" stroke="#cbd5e1" stroke-width="2"/>
          <line x1="125" y1="172" x2="140" y2="172" stroke="#cbd5e1" stroke-width="2"/>
          <line x1="125" y1="179" x2="143" y2="179" stroke="#cbd5e1" stroke-width="2"/>

          <!-- Money symbols -->
          <text x="165" y="175" font-size="20" fill="#cbd5e1">$</text>
          <text x="55" y="165" font-size="20" fill="#cbd5e1">$</text>
        </svg>
      </div>

      <h2 class="empty-title">Kein geplantes Personal gefunden</h2>

      <p class="empty-description">
        Prüfe deine Filtereinstellung oder klicke [+ Mitarbeiter:In Hinzufügen], um neues Personal zu planen.<br/>
        Nur aktive Mitarbeiter:Innen werden in der Liquiditätsplanung berücksichtigt.
      </p>
    </div>
  {/if}

  {#if loading}
    <div class="loading">Lade Personal...</div>
  {/if}
</div>

<style>
  .personal-page {
    max-width: 1400px;
    margin: 0 auto;
  }

  /* Tabs Header */
  .tabs-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    gap: 1rem;
  }

  .tabs-container {
    display: flex;
    gap: 0.5rem;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.625rem 1.25rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: #64748b;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .tab:hover {
    background: #f9fafb;
    color: #1a202c;
  }

  .tab.active {
    background: white;
    border-color: #7c3aed;
    color: #7c3aed;
  }

  :global(.dark) .tab {
    background: #1e293b;
    border-color: #334155;
    color: #94a3b8;
  }

  :global(.dark) .tab:hover {
    background: #334155;
    color: #e2e8f0;
  }

  :global(.dark) .tab.active {
    background: #1e293b;
    border-color: #7c3aed;
    color: #a78bfa;
  }

  .tab-count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 1.5rem;
    height: 1.5rem;
    padding: 0 0.375rem;
    background: #f1f5f9;
    border-radius: 0.375rem;
    font-size: 0.75rem;
    font-weight: 600;
    color: #64748b;
  }

  .tab.active .tab-count {
    background: #ede9fe;
    color: #7c3aed;
  }

  :global(.dark) .tab-count {
    background: #334155;
    color: #94a3b8;
  }

  :global(.dark) .tab.active .tab-count {
    background: rgba(124, 58, 237, 0.2);
    color: #a78bfa;
  }

  /* Header Actions */
  .header-actions {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .employee-count {
    font-size: 0.875rem;
    color: #64748b;
    font-weight: 500;
  }

  :global(.dark) .employee-count {
    color: #94a3b8;
  }

  .add-employee-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.625rem 1.25rem;
    background: #7c3aed;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .add-employee-btn:hover {
    background: #6d28d9;
  }

  .badge-n {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.25rem;
    height: 1.25rem;
    background: rgba(255, 255, 255, 0.25);
    border-radius: 0.25rem;
    font-size: 0.75rem;
    font-weight: 600;
  }

  .instruction-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.625rem 1rem;
    background: #fef3c7;
    color: #92400e;
    border: 1px solid #fcd34d;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .instruction-btn:hover {
    background: #fde68a;
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

  /* Submenu Card */
  .submenu-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.25rem;
    background: #f5f3ff;
    border: 2px solid #e9d5ff;
    border-radius: 0.75rem;
    margin-bottom: 1.5rem;
    transition: all 0.2s ease;
  }

  :global(.dark) .submenu-card {
    background: rgba(124, 58, 237, 0.1);
    border-color: rgba(124, 58, 237, 0.3);
  }

  .submenu-icon {
    width: 3rem;
    height: 3rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    border-radius: 0.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  :global(.dark) .submenu-icon {
    background: #1e293b;
  }

  .submenu-content {
    flex: 1;
  }

  .submenu-title {
    font-size: 1rem;
    font-weight: 600;
    color: #7c3aed;
    margin-bottom: 0.25rem;
  }

  :global(.dark) .submenu-title {
    color: #a78bfa;
  }

  .submenu-subtitle {
    font-size: 0.875rem;
    color: #6b21a8;
  }

  :global(.dark) .submenu-subtitle {
    color: #c4b5fd;
  }

  /* Empty State */
  .empty-state {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
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
    line-height: 1.6;
    max-width: 700px;
    margin: 0 auto;
  }

  :global(.dark) .empty-description {
    color: #94a3b8;
  }

  .loading {
    padding: 4rem 2rem;
    text-align: center;
    color: #64748b;
  }
</style>
