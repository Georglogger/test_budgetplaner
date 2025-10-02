<script>
  import { api } from './api';
  import { addDeletedBudget } from './deletedBudgetsStore.js';

  export let budgetId = null;
  export let onDelete = null; // Callback after successful deletion

  let showConfirmDialog = false;
  let deleting = false;

  function openConfirmDialog() {
    if (!budgetId) {
      console.warn('No budget ID provided for deletion');
      return;
    }
    showConfirmDialog = true;
  }

  function closeConfirmDialog() {
    showConfirmDialog = false;
  }

  async function confirmDelete() {
    if (!budgetId) return;

    deleting = true;
    try {
      // Get budget details before deleting
      const budget = await api.getBudget(budgetId);

      // Delete the budget
      await api.deleteBudget(budgetId);

      // Add to deleted budgets list (persistent)
      addDeletedBudget({ ...budget, id: budgetId });

      closeConfirmDialog();

      // Call callback if provided
      if (onDelete) {
        onDelete();
      } else {
        // Default: navigate to budget list
        window.location.hash = '#/';
      }
    } catch (error) {
      console.error('Error deleting budget:', error);
      alert('Fehler beim Löschen: ' + error.message);
    } finally {
      deleting = false;
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter') {
      confirmDelete();
    } else if (e.key === 'Escape') {
      closeConfirmDialog();
    }
  }
</script>

<!-- Floating Delete Button -->
{#if budgetId}
  <button class="floating-delete-btn" on:click={openConfirmDialog} title="Budget löschen">
    <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
      <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
    </svg>
  </button>
{/if}

<!-- Confirmation Dialog -->
{#if showConfirmDialog}
  <div class="modal-overlay" on:click={closeConfirmDialog} on:keydown={handleKeyDown}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <div class="warning-icon">
          <svg width="48" height="48" viewBox="0 0 48 48" fill="currentColor">
            <circle cx="24" cy="24" r="22" fill="#fef3c7"/>
            <path d="M24 16v12M24 32h.01" stroke="#f59e0b" stroke-width="3" stroke-linecap="round"/>
          </svg>
        </div>
        <h3>Budget löschen?</h3>
        <p class="modal-description">
          Bist du sicher, dass du dieses Budget löschen möchtest?<br/>
          Du kannst es später über den Wiederherstellen-Button wiederherstellen.
        </p>
      </div>

      <div class="modal-footer">
        <button class="btn-cancel" on:click={closeConfirmDialog} disabled={deleting}>
          Abbrechen
        </button>
        <button class="btn-delete" on:click={confirmDelete} disabled={deleting}>
          {#if deleting}
            <span class="spinner"></span>
            Wird gelöscht...
          {:else}
            <svg width="18" height="18" viewBox="0 0 18 18" fill="currentColor">
              <path d="M5 15c0 1.1.9 2 2 2h6c1.1 0 2-.9 2-2V5H5v10zM15 2h-3l-1-1H7L6 2H3v2h12V2z"/>
            </svg>
            Löschen
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  /* Floating Delete Button */
  .floating-delete-btn {
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    width: 64px;
    height: 64px;
    background: #ef4444;
    color: white;
    border: none;
    border-radius: 50%;
    box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    z-index: 1000;
  }

  .floating-delete-btn:hover {
    background: #dc2626;
    transform: scale(1.05);
    box-shadow: 0 6px 16px rgba(239, 68, 68, 0.5);
  }

  .floating-delete-btn:active {
    transform: scale(0.95);
  }

  :global(.dark) .floating-delete-btn {
    background: #ef4444;
  }

  :global(.dark) .floating-delete-btn:hover {
    background: #dc2626;
  }

  /* Modal */
  .modal-overlay {
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

  .modal-content {
    background: white;
    border-radius: 0.75rem;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
    width: 90%;
    max-width: 450px;
    overflow: hidden;
  }

  :global(.dark) .modal-content {
    background: #1e293b;
  }

  .modal-header {
    padding: 2rem;
    text-align: center;
  }

  .warning-icon {
    margin: 0 auto 1.5rem;
    display: flex;
    justify-content: center;
  }

  .modal-header h3 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a202c;
    margin-bottom: 0.75rem;
  }

  :global(.dark) .modal-header h3 {
    color: #e2e8f0;
  }

  .modal-description {
    font-size: 0.875rem;
    color: #64748b;
    line-height: 1.6;
  }

  :global(.dark) .modal-description {
    color: #94a3b8;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 1.5rem;
    border-top: 1px solid #e5e7eb;
  }

  :global(.dark) .modal-footer {
    border-top-color: #334155;
  }

  .btn-cancel,
  .btn-delete {
    padding: 0.625rem 1.5rem;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .btn-cancel {
    background: white;
    border: 1px solid #e5e7eb;
    color: #64748b;
  }

  .btn-cancel:hover:not(:disabled) {
    background: #f9fafb;
  }

  .btn-cancel:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  :global(.dark) .btn-cancel {
    background: #0f172a;
    border-color: #334155;
    color: #94a3b8;
  }

  :global(.dark) .btn-cancel:hover:not(:disabled) {
    background: #1e293b;
  }

  .btn-delete {
    background: #ef4444;
    border: none;
    color: white;
  }

  .btn-delete:hover:not(:disabled) {
    background: #dc2626;
  }

  .btn-delete:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>
