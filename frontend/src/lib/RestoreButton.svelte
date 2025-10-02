<script>
  import { deletedBudgets, removeDeletedBudget, clearAllDeleted } from './deletedBudgetsStore.js';
  import { api } from './api';

  let showPopup = false;

  function togglePopup() {
    showPopup = !showPopup;
  }

  function closePopup() {
    showPopup = false;
  }

  async function restoreBudget(deleted) {
    try {
      // Recreate the budget
      await api.createBudget({
        name: deleted.name,
        description: deleted.description,
        period_start: deleted.period_start,
        period_end: deleted.period_end,
        status: deleted.status
      });

      // Remove from deleted list
      removeDeletedBudget(deleted.id);

      // Reload page
      window.location.reload();
    } catch (error) {
      console.error('Error restoring budget:', error);
      alert('Fehler beim Wiederherstellen: ' + error.message);
    }
  }

  function formatDeletedTime(isoString) {
    const date = new Date(isoString);
    const now = new Date();
    const diff = now - date;

    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return 'Gerade eben';
    if (minutes < 60) return `Vor ${minutes} Min.`;
    if (hours < 24) return `Vor ${hours} Std.`;
    return `Vor ${days} Tag${days > 1 ? 'en' : ''}`;
  }

  function formatFullDate(isoString) {
    return new Date(isoString).toLocaleString('de-DE', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function handleClearAll() {
    if (confirm('Möchtest du alle gelöschten Budgets dauerhaft entfernen?')) {
      clearAllDeleted();
    }
  }
</script>

<!-- Floating Restore Button -->
<button
  class="floating-restore-btn"
  class:has-items={$deletedBudgets.length > 0}
  on:click={togglePopup}
  title="Gelöschte Budgets wiederherstellen"
>
  <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
    <path d="M13 3a9 9 0 00-9 9H1l3.89 3.89.07.14L9 12H6c0-3.87 3.13-7 7-7s7 3.13 7 7-3.13 7-7 7c-1.93 0-3.68-.79-4.94-2.06l-1.42 1.42A8.954 8.954 0 0013 21a9 9 0 000-18zm-1 5v5l4.28 2.54.72-1.21-3.5-2.08V8H12z"/>
  </svg>
  {#if $deletedBudgets.length > 0}
    <span class="badge-count">{$deletedBudgets.length}</span>
  {/if}
</button>

<!-- Popup -->
{#if showPopup}
  <div class="popup-overlay" on:click={closePopup}>
    <div class="popup-content" on:click|stopPropagation>
      <div class="popup-header">
        <h3>
          <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor" style="display: inline-block; vertical-align: middle; margin-right: 0.5rem;">
            <path d="M10 2a8 8 0 00-8 8h-1l3 3 .07.14L7 10H5c0-3.31 2.69-6 6-6s6 2.69 6 6-2.69 6-6 6c-1.66 0-3.14-.69-4.22-1.78l-1.42 1.42A7.95 7.95 0 0010 18a8 8 0 000-16zm-1 4v5l3.62 2.16.72-1.21-2.84-1.7V6H9z"/>
          </svg>
          Gelöschte Budgets
        </h3>
        <div class="header-actions">
          {#if $deletedBudgets.length > 0}
            <button class="clear-all-btn" on:click={handleClearAll}>
              Alle löschen
            </button>
          {/if}
          <button class="close-popup-btn" on:click={closePopup}>
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <path d="M5 5l10 10M15 5L5 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </button>
        </div>
      </div>

      <div class="popup-body">
        {#if $deletedBudgets.length === 0}
          <div class="empty-state">
            <svg width="64" height="64" viewBox="0 0 64 64" fill="none">
              <circle cx="32" cy="32" r="30" fill="#f3f4f6"/>
              <path d="M32 20v16M32 42h.01" stroke="#9ca3af" stroke-width="3" stroke-linecap="round"/>
            </svg>
            <p>Keine gelöschten Budgets</p>
          </div>
        {:else}
          <div class="deleted-list">
            {#each $deletedBudgets as deleted}
              <div class="deleted-item">
                <div class="item-icon">
                  <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M8 2a2 2 0 00-2 2v1H4a1 1 0 000 2h1v9a2 2 0 002 2h6a2 2 0 002-2V7h1a1 1 0 100-2h-2V4a2 2 0 00-2-2H8zm0 2h4v1H8V4z" opacity="0.5"/>
                  </svg>
                </div>
                <div class="item-content">
                  <div class="item-name">{deleted.name}</div>
                  <div class="item-meta">
                    <span class="item-time" title={formatFullDate(deleted.deletedAt)}>
                      {formatDeletedTime(deleted.deletedAt)}
                    </span>
                    {#if deleted.description}
                      <span class="item-separator">•</span>
                      <span class="item-desc">{deleted.description.substring(0, 30)}{deleted.description.length > 30 ? '...' : ''}</span>
                    {/if}
                  </div>
                </div>
                <button class="restore-btn" on:click={() => restoreBudget(deleted)}>
                  <svg width="18" height="18" viewBox="0 0 18 18" fill="currentColor">
                    <path d="M9 2a7 7 0 00-7 7H1l2.5 2.5L6 9H4c0-2.76 2.24-5 5-5s5 2.24 5 5-2.24 5-5 5c-1.38 0-2.63-.56-3.54-1.46l-1.41 1.41A6.95 6.95 0 009 16a7 7 0 000-14z"/>
                  </svg>
                  Wiederherstellen
                </button>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  /* Floating Restore Button */
  .floating-restore-btn {
    position: fixed;
    bottom: 2rem;
    left: 2rem;
    width: 64px;
    height: 64px;
    background: #64748b;
    color: white;
    border: none;
    border-radius: 50%;
    box-shadow: 0 4px 12px rgba(100, 116, 139, 0.4);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    z-index: 1000;
    position: relative;
  }

  .floating-restore-btn:hover {
    background: #475569;
    transform: scale(1.05);
    box-shadow: 0 6px 16px rgba(100, 116, 139, 0.5);
  }

  .floating-restore-btn.has-items {
    background: #7c3aed;
  }

  .floating-restore-btn.has-items:hover {
    background: #6d28d9;
  }

  :global(.dark) .floating-restore-btn {
    background: #64748b;
  }

  :global(.dark) .floating-restore-btn.has-items {
    background: #7c3aed;
  }

  .badge-count {
    position: absolute;
    top: -4px;
    right: -4px;
    min-width: 24px;
    height: 24px;
    background: #ef4444;
    color: white;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 600;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 6px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  /* Popup */
  .popup-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    backdrop-filter: blur(4px);
  }

  .popup-content {
    background: white;
    border-radius: 0.75rem;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
    width: 90%;
    max-width: 600px;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  :global(.dark) .popup-content {
    background: #1e293b;
  }

  .popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e5e7eb;
  }

  :global(.dark) .popup-header {
    border-bottom-color: #334155;
  }

  .popup-header h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1a202c;
    display: flex;
    align-items: center;
  }

  :global(.dark) .popup-header h3 {
    color: #e2e8f0;
  }

  .header-actions {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  .clear-all-btn {
    padding: 0.5rem 1rem;
    background: #fee2e2;
    color: #991b1b;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .clear-all-btn:hover {
    background: #fecaca;
  }

  :global(.dark) .clear-all-btn {
    background: rgba(239, 68, 68, 0.2);
    color: #fca5a5;
  }

  .close-popup-btn {
    background: none;
    border: none;
    color: #94a3b8;
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 0.375rem;
    transition: all 0.2s ease;
  }

  .close-popup-btn:hover {
    background: #f1f5f9;
    color: #64748b;
  }

  :global(.dark) .close-popup-btn:hover {
    background: #334155;
    color: #cbd5e1;
  }

  .popup-body {
    flex: 1;
    overflow-y: auto;
    padding: 1rem;
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
    margin-bottom: 1rem;
  }

  .deleted-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .deleted-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    transition: all 0.2s ease;
  }

  .deleted-item:hover {
    background: #f3f4f6;
  }

  :global(.dark) .deleted-item {
    background: #0f172a;
    border-color: #334155;
  }

  :global(.dark) .deleted-item:hover {
    background: #1e293b;
  }

  .item-icon {
    flex-shrink: 0;
    width: 40px;
    height: 40px;
    background: #e5e7eb;
    border-radius: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #64748b;
  }

  :global(.dark) .item-icon {
    background: #334155;
    color: #94a3b8;
  }

  .item-content {
    flex: 1;
    min-width: 0;
  }

  .item-name {
    font-weight: 600;
    color: #1a202c;
    margin-bottom: 0.25rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  :global(.dark) .item-name {
    color: #e2e8f0;
  }

  .item-meta {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.813rem;
    color: #64748b;
  }

  :global(.dark) .item-meta {
    color: #94a3b8;
  }

  .item-time {
    font-weight: 500;
  }

  .item-separator {
    color: #cbd5e1;
  }

  .item-desc {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .restore-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: #7c3aed;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    flex-shrink: 0;
  }

  .restore-btn:hover {
    background: #6d28d9;
  }
</style>
