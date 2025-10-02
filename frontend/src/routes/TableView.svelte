<script>
  import { onMount, afterUpdate } from 'svelte';
  import { api } from '../lib/api';
  import { Chart, registerables } from 'chart.js';
  import { darkMode, toggleDarkMode } from '../lib/darkMode';
  import DeleteButton from '../lib/DeleteButton.svelte';

  Chart.register(...registerables);

  export let budgetId;
  let chartCanvas;
  let chartInstance;

  let budget = null;
  let budgetLines = [];
  let actuals = [];
  let months = [];
  let loading = true;
  let error = null;
  let selectedMonth = null;
  let expandedCategories = new Set();
  let activeTab = 'revenue'; // 'revenue' or 'expense'
  let activeFilters = new Set(['all']); // Filter pills

  onMount(async () => {
    await loadData();
  });

  afterUpdate(() => {
    if (chartCanvas && !loading && months.length > 0) {
      renderChart();
    }
  });

  function renderChart() {
    if (chartInstance) {
      chartInstance.destroy();
    }

    const ctx = chartCanvas.getContext('2d');

    // Prepare chart data
    const revenueData = months.map(month => {
      return categoryData
        .filter(cat => cat.type === 'revenue')
        .reduce((sum, cat) => {
          return sum + cat.children.reduce((childSum, child) =>
            childSum + (child.actual[month.key] || child.planned[month.key] || 0), 0);
        }, 0);
    });

    const expenseData = months.map(month => {
      return categoryData
        .filter(cat => cat.type !== 'revenue')
        .reduce((sum, cat) => {
          return sum + cat.children.reduce((childSum, child) =>
            childSum + (child.actual[month.key] || child.planned[month.key] || 0), 0);
        }, 0);
    });

    chartInstance = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: months.map(m => m.label),
        datasets: [
          {
            label: 'Einnahmen',
            data: revenueData,
            backgroundColor: '#10b981',
            borderRadius: 6,
            maxBarThickness: 80
          },
          {
            label: 'Ausgaben',
            data: expenseData,
            backgroundColor: '#ef4444',
            borderRadius: 6,
            maxBarThickness: 80
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false
          },
          tooltip: {
            backgroundColor: '#7c3aed',
            padding: 12,
            titleFont: {
              size: 14,
              weight: 'bold'
            },
            bodyFont: {
              size: 13
            },
            callbacks: {
              label: function(context) {
                let label = context.dataset.label || '';
                if (label) {
                  label += ': ';
                }
                label += new Intl.NumberFormat('de-DE', {
                  style: 'currency',
                  currency: 'EUR',
                  minimumFractionDigits: 0
                }).format(context.parsed.y);
                return label;
              }
            }
          }
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              callback: function(value) {
                return '€' + value.toLocaleString('de-DE');
              }
            }
          }
        }
      }
    });
  }

  async function loadData() {
    try {
      loading = true;
      [budget, budgetLines, actuals] = await Promise.all([
        api.getBudget(budgetId),
        api.getBudgetLines(budgetId),
        api.getActuals(budgetId)
      ]);

      generateMonths();
      groupDataByCategory();
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function generateMonths() {
    if (!budget) return;

    const start = new Date(budget.period_start);
    const end = new Date(budget.period_end);
    months = [];

    let current = new Date(start);
    while (current <= end) {
      months.push({
        key: `${current.getFullYear()}-${String(current.getMonth() + 1).padStart(2, '0')}`,
        label: current.toLocaleDateString('de-DE', { month: 'short', year: '2-digit' }),
        fullDate: new Date(current)
      });
      current.setMonth(current.getMonth() + 1);
    }
  }

  let categoryData = [];

  function groupDataByCategory() {
    const categories = new Map();

    // Group budget lines by category
    budgetLines.forEach(line => {
      if (!categories.has(line.category)) {
        categories.set(line.category, {
          name: line.category,
          type: 'revenue', // Standardmäßig Einnahmen
          children: new Map(),
          total: 0
        });
      }

      const cat = categories.get(line.category);

      if (line.subcategory) {
        if (!cat.children.has(line.subcategory)) {
          cat.children.set(line.subcategory, {
            name: line.subcategory,
            planned: {},
            actual: {}
          });
        }

        const subcat = cat.children.get(line.subcategory);
        months.forEach(month => {
          subcat.planned[month.key] = line.amount / months.length;
        });
      }

      cat.total += line.amount;
    });

    // Add actuals
    actuals.forEach(actual => {
      const date = new Date(actual.date);
      const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;

      if (categories.has(actual.category)) {
        const cat = categories.get(actual.category);

        if (actual.subcategory && cat.children.has(actual.subcategory)) {
          const subcat = cat.children.get(actual.subcategory);
          if (!subcat.actual[monthKey]) {
            subcat.actual[monthKey] = 0;
          }
          subcat.actual[monthKey] += actual.amount;
        }
      }
    });

    categoryData = Array.from(categories.values()).map(cat => ({
      ...cat,
      children: Array.from(cat.children.values())
    }));
  }

  function formatCurrency(amount) {
    if (!amount && amount !== 0) return '-';
    return new Intl.NumberFormat('de-DE', {
      style: 'decimal',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  }

  function calculateVariance(planned, actual) {
    if (!planned || planned === 0) return 0;
    return ((actual - planned) / planned) * 100;
  }

  function getVarianceClass(variance) {
    if (Math.abs(variance) < 5) return 'neutral';
    return variance > 0 ? 'positive' : 'negative';
  }

  function handleMonthHover(month) {
    selectedMonth = month;
  }

  function toggleCategory(categoryName) {
    if (expandedCategories.has(categoryName)) {
      expandedCategories.delete(categoryName);
    } else {
      expandedCategories.add(categoryName);
    }
    expandedCategories = expandedCategories;
  }

  function toggleFilter(filterName) {
    if (filterName === 'all') {
      activeFilters = new Set(['all']);
    } else {
      activeFilters.delete('all');
      if (activeFilters.has(filterName)) {
        activeFilters.delete(filterName);
        if (activeFilters.size === 0) {
          activeFilters.add('all');
        }
      } else {
        activeFilters.add(filterName);
      }
      activeFilters = activeFilters;
    }
  }
</script>

<div class="table-view">
  <div class="top-bar">
    <div class="controls">
      <button class="nav-btn">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M11 2L5 8l6 6" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
      </button>
      <button class="nav-btn">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M5 2l6 6-6 6" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
      </button>

      <div class="date-range">
        <span class="icon">📅</span>
        <span>Monatsansicht</span>
        <span class="date-text">
          {budget ? new Date(budget.period_start).toLocaleDateString('de-DE') : ''}
          -
          {budget ? new Date(budget.period_end).toLocaleDateString('de-DE') : ''}
        </span>
      </div>

      <div class="scenario-select">
        <span class="icon">📊</span>
        <span>Basis-Szenario</span>
      </div>
    </div>

    <div class="actions">
      <button class="btn-dark-toggle" on:click={toggleDarkMode} title="Dark Mode umschalten">
        {#if $darkMode}
          ☀️
        {:else}
          🌙
        {/if}
      </button>
      <button class="btn-secondary" on:click={() => window.location.hash = `#/budgets/${budgetId}`}>
        Zurück
      </button>
      <button class="btn-export">
        <span class="icon">📤</span>
        Exportieren
      </button>
    </div>
  </div>

  {#if loading}
    <div class="loading">Lade Daten...</div>
  {:else if error}
    <div class="error">{error}</div>
  {:else}
    <!-- Tab Navigation -->
    <div class="tab-nav-container">
      <div class="tab-buttons">
        <button
          class="tab-btn"
          class:active={activeTab === 'revenue'}
          on:click={() => activeTab = 'revenue'}
        >
          <span class="tab-icon">📈</span>
          Einnahmen
        </button>
        <button
          class="tab-btn"
          class:active={activeTab === 'expense'}
          on:click={() => activeTab = 'expense'}
        >
          <span class="tab-icon">📉</span>
          Ausgaben
        </button>
      </div>

      <!-- Filter Pills -->
      <div class="filter-pills">
        <button
          class="pill"
          class:active={activeFilters.has('all')}
          on:click={() => toggleFilter('all')}
        >
          Alle
          <span class="badge">12</span>
        </button>
        <button
          class="pill"
          class:active={activeFilters.has('aktiv')}
          on:click={() => toggleFilter('aktiv')}
        >
          Aktiv
          <span class="badge">8</span>
        </button>
        <button
          class="pill"
          class:active={activeFilters.has('nicht-aktiv')}
          on:click={() => toggleFilter('nicht-aktiv')}
        >
          Nicht-Aktiv
          <span class="badge">4</span>
        </button>
      </div>
    </div>

    <!-- Chart Section -->
    <div class="chart-container">
      <canvas bind:this={chartCanvas} style="height: 280px;"></canvas>
    </div>

    <!-- Scrollable Table -->
    <div class="table-scroll-container">
      <table class="finance-table">
        <thead>
          <tr>
            <th class="sticky-col category-col">
              <div class="search-box">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                  <circle cx="7" cy="7" r="5" stroke="currentColor" stroke-width="2" fill="none"/>
                  <path d="M11 11l3 3" stroke="currentColor" stroke-width="2"/>
                </svg>
              </div>
            </th>
            <th class="trend-col"></th>
            {#each months as month}
              <th class="month-col" on:mouseenter={() => handleMonthHover(month)}>
                {month.label}
              </th>
            {/each}
          </tr>
        </thead>
        <tbody>
          <!-- Summary Rows -->
          <tr class="summary-row">
            <td class="sticky-col">
              <div class="row-header">
                <span class="icon">💰</span>
                Verfügbare Liquidität
              </div>
            </td>
            <td class="trend-col">
              <div class="sparkline">- - -</div>
            </td>
            {#each months as month}
              <td class="amount-col neutral">0</td>
            {/each}
          </tr>

          <tr class="summary-row">
            <td class="sticky-col">
              <div class="row-header">
                <span class="icon">📈</span>
                Liquidität Start
              </div>
            </td>
            <td class="trend-col">
              <div class="sparkline">- - -</div>
            </td>
            {#each months as month}
              <td class="amount-col neutral">0</td>
            {/each}
          </tr>

          <tr class="summary-row">
            <td class="sticky-col">
              <div class="row-header">
                <span class="icon">🔄</span>
                Netto-Veränderung
              </div>
            </td>
            <td class="trend-col">
              <div class="sparkline">∿ ∿</div>
            </td>
            {#each months as month}
              <td class="amount-col">
                <span class="neutral">0</span>
              </td>
            {/each}
          </tr>

          <!-- Categories -->
          {#each categoryData as category}
            <tr class="category-row expandable" on:click={() => toggleCategory(category.name)}>
              <td class="sticky-col">
                <div class="row-header">
                  <button class="expand-btn" class:expanded={expandedCategories.has(category.name)}>
                    {expandedCategories.has(category.name) ? '▼' : '▶'}
                  </button>
                  <span class="icon">{category.type === 'revenue' ? '📈' : '📉'}</span>
                  <strong>{category.name}</strong>
                </div>
              </td>
              <td class="trend-col">
                <svg width="60" height="20" class="mini-sparkline">
                  <polyline points="0,15 15,10 30,12 45,5 60,8" fill="none" stroke="#10b981" stroke-width="1.5"/>
                </svg>
              </td>
              {#each months as month}
                {@const planned = category.children.reduce((sum, child) =>
                  sum + (child.planned[month.key] || 0), 0)}
                {@const actual = category.children.reduce((sum, child) =>
                  sum + (child.actual[month.key] || 0), 0)}
                {@const value = actual || planned}
                {@const isNegative = value < 0}
                <td class="amount-col">
                  <div class="amount-cell">
                    <span class="value" class:negative={isNegative} class:positive={!isNegative && value > 0}>
                      {formatCurrency(Math.abs(value))}
                    </span>
                    {#if actual && planned}
                      {@const variance = calculateVariance(planned, actual)}
                      <span class="percentage-badge {getVarianceClass(variance)}">
                        {variance > 0 ? '+' : ''}{variance.toFixed(0)}%
                      </span>
                    {/if}
                  </div>
                </td>
              {/each}
            </tr>

            <!-- Subcategories -->
            {#if expandedCategories.has(category.name)}
              {#each category.children as subcategory}
                <tr class="subcategory-row">
                  <td class="sticky-col">
                    <div class="row-header indent-1">
                      {subcategory.name}
                    </div>
                  </td>
                  <td class="trend-col">
                    <svg width="60" height="20" class="mini-sparkline">
                      <polyline points="0,10 20,12 40,8 60,15" fill="none" stroke="#94a3b8" stroke-width="1"/>
                    </svg>
                  </td>
                  {#each months as month}
                    {@const planned = subcategory.planned[month.key] || 0}
                    {@const actual = subcategory.actual[month.key] || 0}
                    {@const value = actual || planned}
                    {@const isNegative = value < 0}
                    <td class="amount-col">
                      <div class="amount-cell">
                        <span class="value" class:negative={isNegative} class:positive={!isNegative && value > 0}>
                          {formatCurrency(Math.abs(value))}
                        </span>
                        {#if actual && planned}
                          {@const variance = calculateVariance(planned, actual)}
                          <span class="percentage-badge {getVarianceClass(variance)}">
                            {variance > 0 ? '+' : ''}{variance.toFixed(0)}%
                          </span>
                        {/if}
                      </div>
                    </td>
                  {/each}
                </tr>
              {/each}
            {/if}
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  :global(body.dark) {
    background: #0f172a;
    color: #e2e8f0;
  }

  .table-view {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background: #fafafa;
    transition: background 0.3s ease;
    padding-right: 3rem;
    padding-bottom: 3rem;
  }

  :global(.dark) .table-view {
    background: #0f172a;
  }

  .top-bar {
    background: white;
    border-bottom: 1px solid #e5e7eb;
    padding: 1rem 1.5rem;
    margin-right: 3rem;
    border-radius: 0.75rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    transition: all 0.3s ease;
  }

  :global(.dark) .top-bar {
    background: #1e293b;
    border-bottom-color: #334155;
  }

  .controls {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .nav-btn {
    padding: 0.5rem;
    border: 1px solid #e5e7eb;
    background: white;
    border-radius: 0.375rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .nav-btn:hover {
    background: #f9fafb;
  }

  :global(.dark) .nav-btn {
    background: #334155;
    border-color: #475569;
  }

  :global(.dark) .nav-btn:hover {
    background: #475569;
  }

  .date-range, .scenario-select {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.375rem;
    background: white;
    font-size: 0.875rem;
    color: #1a202c;
    font-weight: 500;
    transition: all 0.2s ease;
  }

  :global(.dark) .date-range,
  :global(.dark) .scenario-select {
    background: #334155;
    border-color: #475569;
    color: #e2e8f0;
  }

  .date-text {
    color: #1a202c;
  }

  :global(.dark) .date-text {
    color: #e2e8f0;
  }

  .actions {
    display: flex;
    gap: 0.75rem;
  }

  .btn-secondary, .btn-export, .btn-dark-toggle {
    padding: 0.5rem 1rem;
    border-radius: 0.375rem;
    border: 1px solid #e5e7eb;
    background: white;
    cursor: pointer;
    font-size: 0.875rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: all 0.2s ease;
  }

  .btn-dark-toggle {
    padding: 0.5rem 0.75rem;
    font-size: 1.25rem;
  }

  .btn-export:hover, .btn-secondary:hover, .btn-dark-toggle:hover {
    background: #f9fafb;
  }

  :global(.dark) .btn-secondary,
  :global(.dark) .btn-export,
  :global(.dark) .btn-dark-toggle {
    background: #334155;
    border-color: #475569;
    color: #e2e8f0;
  }

  :global(.dark) .btn-export:hover,
  :global(.dark) .btn-secondary:hover,
  :global(.dark) .btn-dark-toggle:hover {
    background: #475569;
  }

  .chart-container {
    background: white;
    padding: 2rem;
    margin: 1.5rem 3rem 1rem 0;
    border-radius: 0.75rem;
    border: 1px solid #e5e7eb;
    height: 350px;
    transition: all 0.3s ease;
  }

  :global(.dark) .chart-container {
    background: #1e293b;
    border-color: #334155;
  }

  .chart-placeholder {
    height: 300px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: #f9fafb;
    border-radius: 0.5rem;
    border: 2px dashed #d1d5db;
  }

  .hint {
    color: #6b7280;
    font-size: 0.875rem;
    margin-top: 0.5rem;
  }

  .table-scroll-container {
    overflow-x: auto;
    overflow-y: visible;
    margin: 0 3rem 1rem 0;
    background: white;
    border-radius: 0.75rem;
    border: 1px solid #e5e7eb;
    transition: all 0.3s ease;
  }

  :global(.dark) .table-scroll-container {
    background: #1e293b;
    border-color: #334155;
  }

  .finance-table {
    width: 100%;
    border-collapse: separate;
    border-spacing: 0;
    font-size: 0.813rem;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  }

  .finance-table thead {
    position: sticky;
    top: 0;
    z-index: 10;
    background: white;
  }

  :global(.dark) .finance-table thead {
    background: #1e293b;
  }

  .finance-table th {
    padding: 0.5rem 1rem;
    text-align: left;
    font-weight: 500;
    color: #64748b;
    background: #f8fafc;
    border-bottom: 1px solid #e2e8f0;
    white-space: nowrap;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.025em;
  }

  .finance-table th:last-child {
    padding-right: 3rem;
  }

  :global(.dark) .finance-table th {
    color: #cbd5e1;
    background: #0f172a;
    border-bottom-color: #475569;
  }

  .sticky-col {
    position: sticky;
    left: 0;
    z-index: 5;
    background: white;
    border-right: 1px solid #e5e7eb;
    min-width: 250px;
  }

  :global(.dark) .sticky-col {
    background: #1e293b;
    border-right-color: #475569;
  }

  :global(.dark) .finance-table thead .sticky-col {
    background: #0f172a;
  }

  .category-col {
    min-width: 250px;
  }

  .trend-col {
    width: 80px;
    text-align: center;
  }

  .month-col {
    min-width: 120px;
    text-align: center;
  }

  .search-box {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.25rem 0.5rem;
    color: #94a3b8;
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 0.375rem;
    font-size: 0.75rem;
  }

  :global(.dark) .search-box {
    background: #334155;
    border-color: #475569;
    color: #94a3b8;
  }

  .finance-table tbody tr {
    border-bottom: 1px solid #f3f4f6;
  }

  :global(.dark) .finance-table tbody tr {
    border-bottom-color: #334155;
  }

  .finance-table tbody tr:hover {
    background: #fafafa;
  }

  :global(.dark) .finance-table tbody tr:hover {
    background: #334155;
  }

  .finance-table td {
    padding: 0.625rem 1rem;
    border-bottom: 1px solid #f3f4f6;
    color: #1a202c;
    font-size: 0.813rem;
  }

  .finance-table td:last-child {
    padding-right: 3rem;
  }

  :global(.dark) .finance-table td {
    border-bottom-color: #334155;
    color: #e2e8f0;
  }

  .row-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .indent-1 {
    padding-left: 2rem;
  }

  .expand-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.25rem;
    color: #6b7280;
    font-size: 0.65rem;
    transition: transform 0.2s ease;
    width: 20px;
    text-align: center;
  }

  .expand-btn.expanded {
    transform: rotate(0deg);
  }

  .icon {
    font-size: 1rem;
  }

  .sparkline {
    font-size: 0.75rem;
    color: #6b7280;
    text-align: center;
  }

  .amount-col {
    text-align: right;
    font-family: 'Courier New', monospace;
  }

  .amount-cell {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 0.5rem;
  }

  .value {
    font-weight: 500;
    font-size: 0.813rem;
  }

  .value.negative {
    color: #ef4444;
  }

  .value.positive {
    color: #10b981;
  }

  .percentage-badge {
    font-size: 0.688rem;
    padding: 0.125rem 0.313rem;
    border-radius: 0.25rem;
    font-weight: 600;
    min-width: 36px;
    text-align: center;
  }

  .percentage-badge.positive {
    background: #d1fae5;
    color: #065f46;
  }

  .percentage-badge.negative {
    background: #fee2e2;
    color: #991b1b;
  }

  .percentage-badge.neutral {
    background: #f3f4f6;
    color: #6b7280;
  }

  .mini-sparkline {
    display: block;
  }

  .summary-row {
    font-weight: 500;
  }

  .category-row {
    background: #f9fafb;
    cursor: pointer;
    transition: background 0.15s ease;
  }

  .category-row:hover {
    background: #f3f4f6;
  }

  :global(.dark) .category-row {
    background: #0f172a;
  }

  :global(.dark) .category-row:hover {
    background: #1e293b;
  }

  .category-row strong {
    font-weight: 600;
    font-size: 0.813rem;
  }

  .subcategory-row {
    color: #6b7280;
    font-size: 0.813rem;
  }

  .subcategory-row:hover {
    background: #fafafa;
  }

  :global(.dark) .subcategory-row {
    color: #94a3b8;
  }

  :global(.dark) .subcategory-row:hover {
    background: #1e293b;
  }

  .loading, .error {
    padding: 2rem;
    text-align: center;
  }

  .error {
    color: #991b1b;
    background: #fee2e2;
    margin: 1rem;
    border-radius: 0.5rem;
  }

  /* Tab Navigation */
  .tab-nav-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 1.5rem 3rem 1rem 0;
    padding: 0.5rem;
    background: white;
    border-radius: 0.75rem;
    border: 1px solid #e5e7eb;
    transition: all 0.3s ease;
  }

  :global(.dark) .tab-nav-container {
    background: #1e293b;
    border-color: #334155;
  }

  .tab-buttons {
    display: flex;
    gap: 0.5rem;
    padding: 0.25rem;
    background: #f3f4f6;
    border-radius: 0.5rem;
  }

  :global(.dark) .tab-buttons {
    background: #0f172a;
  }

  .tab-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1.25rem;
    background: transparent;
    border: none;
    border-radius: 0.375rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: #64748b;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .tab-btn:hover {
    color: #1a202c;
  }

  .tab-btn.active {
    background: white;
    color: #7c3aed;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  :global(.dark) .tab-btn {
    color: #94a3b8;
  }

  :global(.dark) .tab-btn:hover {
    color: #e2e8f0;
  }

  :global(.dark) .tab-btn.active {
    background: #334155;
    color: #a78bfa;
  }

  .tab-icon {
    font-size: 1rem;
  }

  /* Filter Pills */
  .filter-pills {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  .pill {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.375rem 0.875rem;
    background: #f3f4f6;
    border: 1px solid #e5e7eb;
    border-radius: 9999px;
    font-size: 0.813rem;
    font-weight: 500;
    color: #64748b;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .pill:hover {
    background: #e5e7eb;
    border-color: #d1d5db;
  }

  .pill.active {
    background: #7c3aed;
    border-color: #7c3aed;
    color: white;
  }

  :global(.dark) .pill {
    background: #334155;
    border-color: #475569;
    color: #94a3b8;
  }

  :global(.dark) .pill:hover {
    background: #475569;
    border-color: #64748b;
  }

  :global(.dark) .pill.active {
    background: #7c3aed;
    border-color: #7c3aed;
    color: white;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 1.25rem;
    height: 1.25rem;
    padding: 0 0.375rem;
    background: #e5e7eb;
    border-radius: 9999px;
    font-size: 0.688rem;
    font-weight: 600;
    color: #64748b;
  }

  .pill.active .badge {
    background: rgba(255, 255, 255, 0.25);
    color: white;
  }

  :global(.dark) .badge {
    background: #475569;
    color: #cbd5e1;
  }

  :global(.dark) .pill.active .badge {
    background: rgba(255, 255, 255, 0.25);
    color: white;
  }
</style>

<DeleteButton {budgetId} />
