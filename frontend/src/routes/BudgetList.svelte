<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { darkMode, toggleDarkMode } from '../lib/darkMode';
  import { viewMode } from '../lib/viewMode.js';
  import RestoreButton from '../lib/RestoreButton.svelte';

  let budgets = [];
  let loading = true;
  let error = null;

  onMount(async () => {
    try {
      budgets = await api.getBudgets();
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  });

  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('de-DE');
  }

  function getStatusBadge(status) {
    const badges = {
      draft: { color: 'bg-gray-200 text-gray-800', text: 'Entwurf' },
      approved: { color: 'bg-green-200 text-green-800', text: 'Genehmigt' },
      closed: { color: 'bg-red-200 text-red-800', text: 'Abgeschlossen' }
    };
    return badges[status] || badges.draft;
  }
</script>

<div class="page-container">
  <div class="page-header">
    <h1 class="page-title">
      Budgetplanung
      {#if $viewMode === 'customer'}<span class="view-badge">Kunden-Ansicht</span>{/if}
    </h1>

    {#if $viewMode === 'employee'}
      <button class="btn-new-budget" on:click={() => window.location.hash = '#/budgets/new'}>
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <path d="M10 3v14M3 10h14" stroke="white" stroke-width="2" stroke-linecap="round"/>
        </svg>
        Neues Budget
      </button>
    {/if}
  </div>

  <div class="page-content">

    {#if loading}
      <div class="loading-state">Lade Budgets...</div>
    {:else if error}
      <div class="error-state">{error}</div>
    {:else if budgets.length === 0}
      <div class="empty-state">
        <svg width="64" height="64" viewBox="0 0 64 64" fill="none">
          <circle cx="32" cy="32" r="30" fill="#f3f4f6"/>
          <path d="M32 20v16M32 42h.01" stroke="#9ca3af" stroke-width="3" stroke-linecap="round"/>
        </svg>
        <p>Noch keine Budgets vorhanden.</p>
        {#if $viewMode === 'employee'}
          <button class="btn-create-first" on:click={() => window.location.hash = '#/budgets/new'}>
            Erstes Budget erstellen
          </button>
        {/if}
      </div>
    {:else}
      <div class="budget-grid">
        {#each budgets as budget}
          <div class="budget-card" on:click={() => window.location.hash = `#/budgets/${budget.id}`}>
            <div class="card-header">
              <h3 class="card-title">{budget.name}</h3>
              <span class="status-badge {getStatusBadge(budget.status).color}">
                {getStatusBadge(budget.status).text}
              </span>
            </div>
            <p class="card-description">{budget.description || 'Keine Beschreibung'}</p>
            <div class="card-footer">
              <span class="card-date">📅 {formatDate(budget.period_start)} - {formatDate(budget.period_end)}</span>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<RestoreButton />

<style>
  .page-container {
    height: 100vh;
    display: flex;
    flex-direction: column;
    padding: 0;
    margin: 0;
    margin-top: -64px;
    padding-top: 64px;
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem 2rem;
    background: white;
    border-bottom: 1px solid #e5e7eb;
    flex-shrink: 0;
  }

  :global(.dark) .page-header {
    background: #1e293b;
    border-bottom-color: #334155;
  }

  .page-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a202c;
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  :global(.dark) .page-title {
    color: #e2e8f0;
  }

  .view-badge {
    display: inline-block;
    padding: 0.375rem 0.875rem;
    background: #ede9fe;
    color: #7c3aed;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  :global(.dark) .view-badge {
    background: rgba(124, 58, 237, 0.2);
    color: #a78bfa;
  }

  .btn-new-budget {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.5rem;
    background: #7c3aed;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-new-budget:hover {
    background: #6d28d9;
  }

  .page-content {
    flex: 1;
    overflow-y: auto;
    padding: 2rem;
    background: #f9fafb;
  }

  :global(.dark) .page-content {
    background: #0f172a;
  }

  .loading-state,
  .error-state {
    padding: 4rem 2rem;
    text-align: center;
    color: #64748b;
  }

  .error-state {
    color: #ef4444;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    text-align: center;
    color: #94a3b8;
  }

  .empty-state svg {
    margin-bottom: 1.5rem;
  }

  .empty-state p {
    font-size: 1rem;
    margin-bottom: 1.5rem;
  }

  .btn-create-first {
    padding: 0.75rem 1.5rem;
    background: #7c3aed;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-create-first:hover {
    background: #6d28d9;
  }

  .budget-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 1.5rem;
  }

  .budget-card {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    padding: 1.5rem;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  :global(.dark) .budget-card {
    background: #1e293b;
    border-color: #334155;
  }

  .budget-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
  }

  :global(.dark) .budget-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
  }

  .card-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1a202c;
    margin: 0;
  }

  :global(.dark) .card-title {
    color: #e2e8f0;
  }

  .status-badge {
    padding: 0.25rem 0.75rem;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .bg-gray-200 { background-color: #e5e7eb; }
  .text-gray-800 { color: #1f2937; }
  .bg-green-200 { background-color: #bbf7d0; }
  .text-green-800 { color: #166534; }
  .bg-red-200 { background-color: #fecaca; }
  .text-red-800 { color: #991b1b; }

  .card-description {
    color: #6b7280;
    margin: 0.5rem 0;
    font-size: 0.875rem;
  }

  :global(.dark) .card-description {
    color: #94a3b8;
  }

  .card-footer {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid #e5e7eb;
  }

  :global(.dark) .card-footer {
    border-top-color: #475569;
  }

  .card-date {
    color: #6b7280;
    font-size: 0.875rem;
  }

  :global(.dark) .card-date {
    color: #94a3b8;
  }
</style>
