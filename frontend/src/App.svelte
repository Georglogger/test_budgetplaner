<script>
  import BudgetOverview from './routes/BudgetOverview.svelte';
  import BudgetForm from './routes/BudgetForm.svelte';
  import Dashboard from './routes/Dashboard.svelte';
  import ImportActuals from './routes/ImportActuals.svelte';
  import TableView from './routes/TableView.svelte';
  import Transactions from './routes/Transactions.svelte';
  import Personal from './routes/Personal.svelte';
  import Sidebar from './lib/Sidebar.svelte';
  import TopHeader from './lib/TopHeader.svelte';

  let currentRoute = { component: BudgetOverview, props: {}, name: 'planning' };

  function parseHash() {
    const hash = window.location.hash.slice(1) || '/';

    if (hash === '/' || hash === '') {
      currentRoute = { component: BudgetOverview, props: {}, name: 'planning' };
    } else if (hash === '/transactions') {
      currentRoute = { component: Transactions, props: {}, name: 'transactions' };
    } else if (hash === '/personal' || hash === '/personal/employees') {
      currentRoute = { component: Personal, props: {}, name: 'personal' };
    } else if (hash === '/contracts') {
      currentRoute = { component: BudgetOverview, props: {}, name: 'contracts' };
    } else if (hash === '/reports') {
      currentRoute = { component: BudgetOverview, props: {}, name: 'reports' };
    } else if (hash === '/categories') {
      currentRoute = { component: BudgetOverview, props: {}, name: 'categories' };
    } else if (hash === '/consolidation') {
      currentRoute = { component: BudgetOverview, props: {}, name: 'consolidation' };
    } else if (hash === '/budgets/new') {
      currentRoute = { component: BudgetForm, props: {}, name: 'planning' };
    } else if (hash.match(/^\/budgets\/([^/]+)\/table$/)) {
      const budgetId = hash.match(/^\/budgets\/([^/]+)\/table$/)[1];
      currentRoute = { component: TableView, props: { budgetId }, name: 'planning' };
    } else if (hash.match(/^\/budgets\/([^/]+)\/import$/)) {
      const budgetId = hash.match(/^\/budgets\/([^/]+)\/import$/)[1];
      currentRoute = { component: ImportActuals, props: { budgetId }, name: 'planning' };
    } else if (hash.match(/^\/budgets\/([^/]+)$/)) {
      const budgetId = hash.match(/^\/budgets\/([^/]+)$/)[1];
      currentRoute = { component: Dashboard, props: { budgetId }, name: 'planning' };
    } else {
      currentRoute = { component: BudgetOverview, props: {}, name: 'planning' };
    }
  }

  // Initial route
  parseHash();

  // Listen for hash changes
  window.addEventListener('hashchange', parseHash);
</script>

<div class="flex min-h-screen bg-gray-50 dark:bg-slate-900">
  <Sidebar currentRoute={currentRoute.name} />
  <div class="flex-1 flex flex-col ml-[200px]">
    <TopHeader />
    <main class="flex-1 w-full mt-16 min-h-[calc(100vh-64px)]">
      <svelte:component this={currentRoute.component} {...currentRoute.props} />
    </main>
  </div>
</div>

<style>
  :global(html) {
    margin: 0;
    padding: 0;
  }

  :global(body) {
    margin: 0 !important;
    padding: 0 !important;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    background-color: #f9fafb;
    color: #1a202c;
    overflow-x: hidden;
  }

  :global(body.dark) {
    background-color: #0f172a;
    color: #e2e8f0;
  }

  :global(*) {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }

  /* Main element padding for non-page-container routes */
  main:not(:has(.page-container)) :global(.container) {
    padding: 2rem;
  }
</style>
