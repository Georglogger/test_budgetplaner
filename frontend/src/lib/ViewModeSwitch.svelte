<script>
  import { viewMode, switchViewMode } from './viewMode.js';

  let showPasswordDialog = false;
  let password = '';
  let errorMessage = '';

  function openDialog() {
    showPasswordDialog = true;
    password = '';
    errorMessage = '';
  }

  function closeDialog() {
    showPasswordDialog = false;
    password = '';
    errorMessage = '';
  }

  function handleSwitch() {
    const success = switchViewMode(password);
    if (success) {
      closeDialog();
    } else {
      errorMessage = 'Falsches Passwort!';
      password = '';
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter') {
      handleSwitch();
    } else if (e.key === 'Escape') {
      closeDialog();
    }
  }
</script>

<div class="view-mode-switch">
  <button class="switch-btn" on:click={openDialog} title="Ansicht wechseln">
    {#if $viewMode === 'employee'}
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path d="M10 2a4 4 0 100 8 4 4 0 000-8zm0 10c-4.42 0-8 1.79-8 4v2h16v-2c0-2.21-3.58-4-8-4z"/>
      </svg>
      <span>Mitarbeiter</span>
    {:else}
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path d="M10 2a4 4 0 100 8 4 4 0 000-8zM4 12c0-2 2.67-3 6-3s6 1 6 3v2H4v-2z"/>
        <path d="M17 12h2v4h-2z" opacity="0.5"/>
      </svg>
      <span>Kunde</span>
    {/if}
    <svg class="switch-icon" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
      <path d="M8 3v10M5 6l3-3 3 3M5 10l3 3 3-3" stroke="currentColor" stroke-width="1.5" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  </button>
</div>

{#if showPasswordDialog}
  <!-- Modal Overlay -->
  <div class="modal-overlay" on:click={closeDialog}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h3>Ansicht wechseln</h3>
        <button class="close-btn" on:click={closeDialog}>
          <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
            <path d="M5 5l10 10M15 5L5 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </button>
      </div>

      <div class="modal-body">
        <p class="modal-description">
          Wechsel zu: <strong>{$viewMode === 'employee' ? 'Kunden-Ansicht' : 'Mitarbeiter-Ansicht'}</strong>
        </p>

        <div class="input-group">
          <label for="password">Passwort</label>
          <input
            id="password"
            type="password"
            bind:value={password}
            placeholder="Passwort eingeben..."
            class:error={errorMessage}
            on:keydown={handleKeyDown}
            autofocus
          />
          {#if errorMessage}
            <span class="error-message">{errorMessage}</span>
          {/if}
        </div>

        <div class="info-box">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <circle cx="8" cy="8" r="7" fill="#7c3aed"/>
            <path d="M8 7v4M8 5h.01" stroke="white" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <span>Die Kunden-Ansicht ist vereinfacht und zeigt weniger Details.</span>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn-cancel" on:click={closeDialog}>Abbrechen</button>
        <button class="btn-confirm" on:click={handleSwitch}>Wechseln</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .view-mode-switch {
    position: relative;
  }

  .switch-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
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

  .switch-btn:hover {
    background: #f9fafb;
    border-color: #7c3aed;
    color: #7c3aed;
  }

  :global(.dark) .switch-btn {
    background: #1e293b;
    border-color: #334155;
    color: #94a3b8;
  }

  :global(.dark) .switch-btn:hover {
    background: #334155;
    border-color: #7c3aed;
    color: #a78bfa;
  }

  .switch-icon {
    margin-left: 0.25rem;
    opacity: 0.6;
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
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e5e7eb;
  }

  :global(.dark) .modal-header {
    border-bottom-color: #334155;
  }

  .modal-header h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1a202c;
  }

  :global(.dark) .modal-header h3 {
    color: #e2e8f0;
  }

  .close-btn {
    background: none;
    border: none;
    color: #94a3b8;
    cursor: pointer;
    padding: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 0.375rem;
    transition: all 0.2s ease;
  }

  .close-btn:hover {
    background: #f1f5f9;
    color: #64748b;
  }

  :global(.dark) .close-btn:hover {
    background: #334155;
    color: #cbd5e1;
  }

  .modal-body {
    padding: 1.5rem;
  }

  .modal-description {
    font-size: 0.875rem;
    color: #64748b;
    margin-bottom: 1.5rem;
  }

  :global(.dark) .modal-description {
    color: #94a3b8;
  }

  .modal-description strong {
    color: #7c3aed;
    font-weight: 600;
  }

  :global(.dark) .modal-description strong {
    color: #a78bfa;
  }

  .input-group {
    margin-bottom: 1rem;
  }

  .input-group label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: #64748b;
    margin-bottom: 0.5rem;
  }

  :global(.dark) .input-group label {
    color: #94a3b8;
  }

  .input-group input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    transition: all 0.2s ease;
  }

  .input-group input:focus {
    outline: none;
    border-color: #7c3aed;
    box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.1);
  }

  .input-group input.error {
    border-color: #ef4444;
  }

  :global(.dark) .input-group input {
    background: #0f172a;
    border-color: #334155;
    color: #e2e8f0;
  }

  :global(.dark) .input-group input:focus {
    border-color: #7c3aed;
    box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.2);
  }

  .error-message {
    display: block;
    font-size: 0.813rem;
    color: #ef4444;
    margin-top: 0.5rem;
  }

  .info-box {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 1rem;
    background: #ede9fe;
    border-radius: 0.5rem;
    font-size: 0.813rem;
    color: #5b21b6;
    line-height: 1.5;
  }

  :global(.dark) .info-box {
    background: rgba(124, 58, 237, 0.1);
    color: #c4b5fd;
  }

  .info-box svg {
    flex-shrink: 0;
    margin-top: 0.125rem;
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
  .btn-confirm {
    padding: 0.625rem 1.5rem;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-cancel {
    background: white;
    border: 1px solid #e5e7eb;
    color: #64748b;
  }

  .btn-cancel:hover {
    background: #f9fafb;
  }

  :global(.dark) .btn-cancel {
    background: #0f172a;
    border-color: #334155;
    color: #94a3b8;
  }

  :global(.dark) .btn-cancel:hover {
    background: #1e293b;
  }

  .btn-confirm {
    background: #7c3aed;
    border: none;
    color: white;
  }

  .btn-confirm:hover {
    background: #6d28d9;
  }
</style>
