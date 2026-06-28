<script>
  import { createEventDispatcher } from 'svelte'
  import AdminLanding from '../components/AdminLanding.svelte'
  import UserDashboard from '../components/UserDashboard.svelte'

  const dispatch = createEventDispatcher()

  export let user
  export let token

  function handleLogout() {
    dispatch('logout')
  }

  let currentTab = 'profile'
</script>

<div class="dashboard-container">
  <nav class="navbar">
    <div class="navbar-brand">
      <div class="logo">
        <span class="logo-icon">�</span>
        <h1>ZLC</h1>
      </div>
      <p class="user-email">{user?.email}</p>
    </div>
    <button class="logout-btn" on:click={handleLogout}>
      <span>🚪</span> Logout
    </button>
  </nav>

  <div class="dashboard-content">
    {#if user?.role_id === 1}
      <AdminLanding {user} {token} />
    {:else}
      <UserDashboard {user} {token} />
    {/if}
  </div>
</div>

<style>
  .dashboard-container {
    min-height: 100vh;
    background: radial-gradient(circle at top, rgba(15, 23, 42, 0.95), #020617 80%);
    color: #e2e8f0;
  }

  .navbar {
    background: rgba(15, 23, 42, 0.95);
    padding: 18px 36px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid rgba(148, 163, 184, 0.15);
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .navbar-brand {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .logo-icon {
    font-size: 32px;
    display: block;
  }

  .logo h1 {
    margin: 0;
    font-size: 28px;
    color: #60a5fa;
    font-weight: 700;
  }

  .user-email {
    color: #94a3b8;
    font-size: 13px;
    margin: 0;
    margin-top: -4px;
  }

  .logout-btn {
    background: linear-gradient(135deg, #7c3aed 0%, #22d3ee 100%);
    color: white;
    border: none;
    padding: 14px 26px;
    border-radius: 999px;
    cursor: pointer;
    font-weight: 700;
    transition: transform 0.25s ease, box-shadow 0.25s ease;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 15px;
    box-shadow: 0 10px 30px rgba(124, 58, 237, 0.25);
  }

  .logout-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 14px 32px rgba(34, 211, 238, 0.3);
  }

  .dashboard-content {
    padding: 48px 24px 64px;
    max-width: 1320px;
    margin: 0 auto;
    animation: slideIn 0.45s ease-out;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(24px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .dashboard-content > :global(*) {
    border-radius: 24px;
    background: rgba(15, 23, 42, 0.85);
    box-shadow: 0 24px 80px rgba(15, 23, 42, 0.35);
    border: 1px solid rgba(148, 163, 184, 0.12);
  }
</style>
