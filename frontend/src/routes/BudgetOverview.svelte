<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { viewMode } from '../lib/viewMode.js';
  import RestoreButton from '../lib/RestoreButton.svelte';
  import { Chart, registerables } from 'chart.js';
  import { darkMode } from '../lib/darkMode';

  Chart.register(...registerables);

  let budgets = [];
  let loading = true;
  let error = null;
  let chartCanvas;
  let chartInstance;

  onMount(async () => {
    try {
      budgets = await api.getBudgets();
      renderChart();
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  });

  function renderChart() {
    if (!chartCanvas || budgets.length === 0) return;

    if (chartInstance) {
      chartInstance.destroy();
    }

    const ctx = chartCanvas.getContext('2d');

    // Calculate totals for each budget
    const labels = budgets.map(b => b.name);
    const revenueData = budgets.map(() => Math.random() * 100000); // Placeholder - replace with real data
    const expenseData = budgets.map(() => Math.random() * 80000); // Placeholder - replace with real data

    chartInstance = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: labels,
        datasets: [
          {
            label: 'Einnahmen',
            data: revenueData,
            backgroundColor: '#10b981',
            borderRadius: 6,
            maxBarThickness: 60
          },
          {
            label: 'Ausgaben',
            data: expenseData,
            backgroundColor: '#ef4444',
            borderRadius: 6,
            maxBarThickness: 60
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: true,
            position: 'top',
            labels: {
              font: { size: 14 },
              padding: 20,
              color: $darkMode ? '#e2e8f0' : '#1a202c'
            }
          }
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              color: $darkMode ? '#94a3b8' : '#64748b',
              callback: function(value) {
                return '€' + value.toLocaleString('de-DE');
              }
            },
            grid: {
              color: $darkMode ? '#334155' : '#e5e7eb'
            }
          },
          x: {
            ticks: {
              color: $darkMode ? '#94a3b8' : '#64748b'
            },
            grid: {
              display: false
            }
          }
        }
      }
    });
  }

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

  function openBudget(budgetId) {
    window.location.hash = `#/budgets/${budgetId}`;
  }

  $: if (chartCanvas && !loading) {
    renderChart();
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
      <!-- Chart Section -->
      <div class="chart-container">
        <div class="chart-header">
          <h2>Übersicht aller Budgets</h2>
        </div>
        <div class="chart-wrapper">
          <canvas bind:this={chartCanvas}></canvas>
        </div>
      </div>

      <!-- Table Section -->
      <div class="table-container">
        <table class="budget-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Beschreibung</th>
              <th>Zeitraum</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {#each budgets as budget}
              <tr class="budget-row" on:click={() => openBudget(budget.id)}>
                <td class="budget-name">{budget.name}</td>
                <td class="budget-description">{budget.description || 'Keine Beschreibung'}</td>
                <td class="budget-period">
                  {formatDate(budget.period_start)} - {formatDate(budget.period_end)}
                </td>
                <td>
                  <span class="status-badge {getStatusBadge(budget.status).color}">
                    {getStatusBadge(budget.status).text}
                  </span>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</div>

<RestoreButton />

<style>
  .page-container {
    width: 100%;
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
    padding: 1.5rem 0 1.5rem 2rem;
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
    margin-right: 2rem;
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
    padding: 0;
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

  /* Chart Section */
  .chart-container {
    background: white;
    border-radius: 0;
    padding: 2rem 0 2rem 2rem;
    margin-bottom: 0;
    margin-right: 0;
    box-shadow: none;
    border-bottom: 1px solid #e5e7eb;
  }

  :global(.dark) .chart-container {
    background: #1e293b;
    box-shadow: none;
    border-bottom-color: #334155;
  }

  .chart-header {
    margin-bottom: 1.5rem;
  }

  .chart-header h2 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1a202c;
  }

  :global(.dark) .chart-header h2 {
    color: #e2e8f0;
  }

  .chart-wrapper {
    height: 300px;
    position: relative;
  }

  /* Table Section */
  .table-container {
    background: white;
    border-radius: 0;
    overflow: hidden;
    margin-right: 0;
    box-shadow: none;
  }

  :global(.dark) .table-container {
    background: #1e293b;
    box-shadow: none;
  }

  .budget-table {
    width: 100%;
    border-collapse: collapse;
  }

  .budget-table thead {
    background: #f9fafb;
    border-bottom: 2px solid #e5e7eb;
  }

  :global(.dark) .budget-table thead {
    background: #0f172a;
    border-bottom-color: #334155;
  }

  .budget-table th {
    padding: 1rem 0 1rem 2rem;
    text-align: left;
    font-size: 0.875rem;
    font-weight: 600;
    color: #64748b;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .budget-table th:first-child {
    padding-left: 2rem;
  }

  :global(.dark) .budget-table th {
    color: #94a3b8;
  }

  .budget-row {
    border-bottom: 1px solid #e5e7eb;
    cursor: pointer;
    transition: background 0.15s ease;
  }

  :global(.dark) .budget-row {
    border-bottom-color: #334155;
  }

  .budget-row:hover {
    background: #f9fafb;
  }

  :global(.dark) .budget-row:hover {
    background: #334155;
  }

  .budget-table td {
    padding: 1rem 0 1rem 2rem;
    font-size: 0.875rem;
    color: #1a202c;
  }

  .budget-table td:first-child {
    padding-left: 2rem;
  }

  :global(.dark) .budget-table td {
    color: #e2e8f0;
  }

  .budget-name {
    font-weight: 600;
    color: #7c3aed;
  }

  :global(.dark) .budget-name {
    color: #a78bfa;
  }

  .budget-description {
    color: #64748b;
  }

  :global(.dark) .budget-description {
    color: #94a3b8;
  }

  .budget-period {
    color: #64748b;
    font-size: 0.813rem;
  }

  :global(.dark) .budget-period {
    color: #94a3b8;
  }

  .status-badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 9999px;
    font-size: 0.813rem;
    font-weight: 500;
  }

  .bg-gray-200 { background-color: #e5e7eb; }
  .text-gray-800 { color: #1f2937; }
  .bg-green-200 { background-color: #bbf7d0; }
  .text-green-800 { color: #166534; }
  .bg-red-200 { background-color: #fecaca; }
  .text-red-800 { color: #991b1b; }
</style>
