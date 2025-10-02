<script>
  import { darkMode, toggleDarkMode } from './darkMode';
  import ViewModeSwitch from './ViewModeSwitch.svelte';
</script>

<header class="top-header">
  <div class="header-left">
    <button class="back-btn">
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path d="M12 4L6 10L12 16" stroke="currentColor" stroke-width="2" fill="none"/>
      </svg>
    </button>

    <div class="logo">
      <div class="logo-icon">T</div>
    </div>
  </div>

  <div class="header-center">
    <button class="search-btn">
      <svg width="18" height="18" viewBox="0 0 18 18" fill="none">
        <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5"/>
        <path d="M12.5 12.5L16 16" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
      <span>Suche & Quick-Actions</span>
      <kbd class="kbd">
        <span>STRG</span>
        <span>K</span>
      </kbd>
    </button>
  </div>

  <div class="header-right">
    <ViewModeSwitch />

    <button class="icon-btn">
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <rect x="3" y="3" width="6" height="6" rx="1"/>
        <rect x="11" y="3" width="6" height="6" rx="1"/>
        <rect x="3" y="11" width="6" height="6" rx="1"/>
        <rect x="11" y="11" width="6" height="6" rx="1"/>
      </svg>
    </button>

    <button class="icon-btn" on:click={toggleDarkMode}>
      {#if $darkMode}
        ☀️
      {:else}
        🌙
      {/if}
    </button>

    <button class="icon-btn">
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path d="M10 10c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4m0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4"/>
      </svg>
    </button>
  </div>
</header>

<style>
  .top-header {
    height: 64px;
    background: white;
    border-bottom: 1px solid #e5e7eb;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 1.5rem;
    gap: 2rem;
    position: fixed;
    top: 0;
    left: 200px;
    right: 0;
    z-index: 10;
    transition: all 0.3s ease;
  }

  :global(.dark) .top-header {
    background: #0f172a;
    border-bottom-color: #334155;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .back-btn {
    background: none;
    border: none;
    color: #64748b;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 0.375rem;
    transition: all 0.15s ease;
  }

  .back-btn:hover {
    background: #f1f5f9;
    color: #334155;
  }

  :global(.dark) .back-btn {
    color: #94a3b8;
  }

  :global(.dark) .back-btn:hover {
    background: #1e293b;
    color: #e2e8f0;
  }

  .logo {
    display: flex;
    align-items: center;
  }

  .logo-icon {
    width: 32px;
    height: 32px;
    background: linear-gradient(135deg, #7c3aed 0%, #a78bfa 100%);
    color: white;
    border-radius: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 1.125rem;
  }

  .header-center {
    flex: 1;
    max-width: 600px;
  }

  .search-btn {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.625rem 1rem;
    background: #f8f9fa;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    color: #64748b;
    cursor: pointer;
    transition: all 0.15s ease;
    font-size: 0.875rem;
  }

  .search-btn:hover {
    background: #f1f5f9;
    border-color: #cbd5e1;
  }

  :global(.dark) .search-btn {
    background: #1e293b;
    border-color: #334155;
    color: #94a3b8;
  }

  :global(.dark) .search-btn:hover {
    background: #334155;
    border-color: #475569;
  }

  .search-btn span {
    flex: 1;
    text-align: left;
  }

  .kbd {
    display: flex;
    gap: 0.25rem;
  }

  .kbd span {
    padding: 0.125rem 0.375rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.25rem;
    font-size: 0.75rem;
    font-weight: 500;
    color: #64748b;
  }

  :global(.dark) .kbd span {
    background: #0f172a;
    border-color: #475569;
    color: #94a3b8;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .icon-btn {
    background: none;
    border: none;
    color: #64748b;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 0.375rem;
    transition: all 0.15s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .icon-btn:hover {
    background: #f1f5f9;
    color: #334155;
  }

  :global(.dark) .icon-btn {
    color: #94a3b8;
  }

  :global(.dark) .icon-btn:hover {
    background: #1e293b;
    color: #e2e8f0;
  }
</style>
