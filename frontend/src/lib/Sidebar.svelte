<script>
  import { createEventDispatcher } from 'svelte';

  export let currentRoute = 'planning';
  export let sidebarOpen = false;

  const dispatch = createEventDispatcher();

  const menuItems = [
    { id: 'planning', icon: '📊', label: 'Planung', route: '#/' },
    { id: 'transactions', icon: '📝', label: 'Transaktionen', route: '#/transactions' },
    { id: 'personal', icon: '👥', label: 'Personal', route: '#/personal', hasSubmenu: true },
    { id: 'contracts', icon: '📄', label: 'Verträge', route: '#/contracts' },
    { id: 'reports', icon: '📈', label: 'Berichte', route: '#/reports' },
    { id: 'categories', icon: '🏷️', label: 'Kategorien', route: '#/categories' },
    { id: 'consolidation', icon: '📊', label: 'Konsolidierung', route: '#/consolidation' }
  ];

  let expandedItems = new Set();

  function toggleItem(id) {
    if (expandedItems.has(id)) {
      expandedItems.delete(id);
    } else {
      expandedItems.add(id);
    }
    expandedItems = expandedItems;
  }

  function handleMenuClick() {
    // Close sidebar on mobile when menu item is clicked
    dispatch('closeSidebar');
  }
</script>

<aside class="sidebar" class:open={sidebarOpen}>
  <div class="sidebar-content">
    {#each menuItems as item}
      <a
        href={item.route}
        class="menu-item"
        class:active={currentRoute === item.id}
        on:click={(e) => {
          if (item.hasSubmenu) {
            toggleItem(item.id);
          } else {
            handleMenuClick();
          }
        }}
      >
        <span class="menu-icon">{item.icon}</span>
        <span class="menu-label">{item.label}</span>
        {#if item.hasSubmenu}
          <svg class="menu-arrow" width="12" height="12" viewBox="0 0 12 12" class:expanded={expandedItems.has(item.id)}>
            <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" fill="none"/>
          </svg>
        {/if}
      </a>

      {#if item.hasSubmenu && expandedItems.has(item.id)}
        <div class="submenu">
          <a href="#/personal/employees" class="submenu-item" on:click={handleMenuClick}>
            <span class="menu-icon">👥</span>
            <span class="menu-label">Personal</span>
          </a>
        </div>
      {/if}
    {/each}
  </div>
</aside>

<style>
  .sidebar {
    width: 200px;
    height: 100vh;
    background: #f8f9fa;
    border-right: 1px solid #e5e7eb;
    display: flex;
    flex-direction: column;
    position: fixed;
    left: 0;
    top: 0;
    transition: all 0.3s ease;
    z-index: 50;
  }

  /* Mobile: Hide sidebar by default */
  @media (max-width: 1023px) {
    .sidebar {
      transform: translateX(-100%);
    }

    .sidebar.open {
      transform: translateX(0);
    }
  }

  :global(.dark) .sidebar {
    background: #0f172a;
    border-right-color: #334155;
  }

  .sidebar-content {
    flex: 1;
    overflow-y: auto;
    padding: 0.5rem 0;
  }

  .menu-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    color: #64748b;
    text-decoration: none;
    transition: all 0.15s ease;
    cursor: pointer;
    border-left: 3px solid transparent;
  }

  .menu-item:hover {
    background: #f1f5f9;
    color: #334155;
  }

  .menu-item.active {
    background: #f1f5f9;
    color: #7c3aed;
    border-left-color: #7c3aed;
    font-weight: 500;
  }

  :global(.dark) .menu-item {
    color: #94a3b8;
  }

  :global(.dark) .menu-item:hover {
    background: #1e293b;
    color: #e2e8f0;
  }

  :global(.dark) .menu-item.active {
    background: #1e293b;
    color: #a78bfa;
    border-left-color: #a78bfa;
  }

  .menu-icon {
    font-size: 1.25rem;
    width: 1.5rem;
    height: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .menu-label {
    flex: 1;
    font-size: 0.875rem;
  }

  .menu-arrow {
    transition: transform 0.2s ease;
  }

  .menu-arrow.expanded {
    transform: rotate(180deg);
  }

  .submenu {
    padding-left: 1rem;
  }

  .submenu-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.625rem 1rem;
    color: #64748b;
    text-decoration: none;
    font-size: 0.875rem;
    transition: all 0.15s ease;
  }

  .submenu-item:hover {
    background: #f1f5f9;
    color: #334155;
  }

  :global(.dark) .submenu-item {
    color: #94a3b8;
  }

  :global(.dark) .submenu-item:hover {
    background: #1e293b;
    color: #e2e8f0;
  }
</style>
