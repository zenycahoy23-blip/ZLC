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
        <span class="logo-icon">🔄</span>
        <h1>Uno Reverse</h1>
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
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .navbar {
    background: rgba(255, 255, 255, 0.95);
    padding: 20px 40px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
    backdrop-filter: blur(10px);
  }

  .navbar-brand {
    display: flex;
    align-items: center;
    gap: 15px;
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
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    font-weight: 700;
  }

  .user-email {
    color: #666;
    font-size: 13px;
    margin: 0;
    margin-top: -5px;
  }

  .logout-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 12px 24px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.3s;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 15px;
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
  }

  .logout-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
  }

  .logout-btn:active {
    transform: translateY(0);
  }

  .dashboard-content {
    padding: 40px;
    max-width: 1400px;
    margin: 0 auto;
    animation: slideIn 0.5s ease-out;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
