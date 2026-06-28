<script>
  import { onMount, createEventDispatcher } from 'svelte'
  import AdminAuditDashboard from '../components/AdminAuditDashboard.svelte'

  export let user
  export let token

  const dispatch = createEventDispatcher()
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8082'

  let currentTab = 'dashboard'
  let stats = { totalUsers: 0, verifiedUsers: 0, failedLogins: 0 }
  let loading = true

  const tabs = [
    { id: 'dashboard', label: '📊 Dashboard', icon: '📊' },
    { id: 'audit', label: '📋 Audit Logs', icon: '📋' }
  ]

  onMount(async () => {
    await loadStats()
  })

  async function loadStats() {
    loading = true
    try {
      const [usersRes, statsRes] = await Promise.all([
        fetch(`${apiUrl}/api/admin/all-users`, {
          headers: { 'Authorization': `Bearer ${token}` }
        }),
        fetch(`${apiUrl}/api/admin/dashboard-stats`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
      ])

      const users = await usersRes.json()
      const statsData = await statsRes.json()

      stats = {
        totalUsers: statsData?.total_users || users?.length || 0,
        verifiedUsers: users?.filter(u => u.email_verified && u.phone_verified)?.length || 0,
        failedLogins: statsData?.failed_login_attempts || 0
      }
    } catch (err) {
      console.error('Failed to load stats', err)
    } finally {
      loading = false
    }
  }

  function handleLogout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    dispatch('logout')
  }

  function handleTabChange(tabId) {
    currentTab = tabId
  }
</script>

<div class="admin-container">
  <!-- Header with Logout -->
  <div class="admin-header">
    <div class="header-left">
      <h2>Admin Dashboard</h2>
      <p class="admin-subtitle">System Administration & Monitoring</p>
    </div>
    <div class="header-right">
      <span class="admin-email">{user?.email}</span>
      <span class="admin-badge">ADMIN</span>
      <button class="logout-btn" on:click={handleLogout}>
        🚪 Logout
      </button>
    </div>
  </div>

  <!-- Tabs Navigation -->
  <div class="tabs-navigation">
    {#each tabs as tab}
      <button
        class="tab-btn {currentTab === tab.id ? 'active' : ''}"
        on:click={() => handleTabChange(tab.id)}
      >
        <span class="tab-icon">{tab.icon}</span>
        <span class="tab-label">{tab.label}</span>
      </button>
    {/each}
  </div>

  <!-- Tab Content -->
  <div class="tab-content">
    {#if currentTab === 'dashboard'}
      <div class="dashboard-grid">
        <div class="stat-card">
          <div class="stat-icon">👥</div>
          <div class="stat-details">
            <p class="stat-label">Total Users</p>
            <p class="stat-value">{loading ? '...' : stats.totalUsers}</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">✅</div>
          <div class="stat-details">
            <p class="stat-label">Verified Users</p>
            <p class="stat-value">{loading ? '...' : stats.verifiedUsers}</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">📊</div>
          <div class="stat-details">
            <p class="stat-label">System Status</p>
            <p class="stat-value" style="color: #90ee90;">Online</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">⚠️</div>
          <div class="stat-details">
            <p class="stat-label">Failed Logins (24h)</p>
            <p class="stat-value">{loading ? '...' : stats.failedLogins}</p>
          </div>
        </div>
      </div>
    {:else if currentTab === 'audit'}
      <AdminAuditDashboard {token} />
    {/if}
  </div>
</div>

<style>
  .admin-container {
    padding: 40px;
    max-width: 1400px;
    margin: 0 auto;
  }

  .admin-header {
    margin-bottom: 40px;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 12px;
    padding: 30px;
    color: white;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header-left h2 {
    margin: 0 0 10px 0;
    font-size: 32px;
  }

  .admin-subtitle {
    margin: 0;
    opacity: 0.8;
    font-size: 16px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .admin-email {
    font-size: 14px;
    opacity: 0.9;
  }

  .admin-badge {
    background: rgba(255, 255, 255, 0.2);
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 700;
  }

  .logout-btn {
    background: rgba(220, 53, 69, 0.3);
    color: white;
    border: 2px solid rgba(220, 53, 69, 0.5);
    padding: 10px 20px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 700;
    transition: all 0.2s;
  }

  .logout-btn:hover {
    background: rgba(220, 53, 69, 0.7);
    border-color: #dc3545;
    transform: translateY(-2px);
  }

  .tabs-navigation {
    display: flex;
    gap: 12px;
    margin-bottom: 30px;
    flex-wrap: wrap;
    background: rgba(255, 255, 255, 0.05);
    padding: 12px;
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .tab-btn {
    background: rgba(255, 255, 255, 0.1);
    color: white;
    border: 1px solid rgba(255, 255, 255, 0.2);
    padding: 12px 24px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.3s;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
  }

  .tab-btn:hover {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.3);
    transform: translateY(-2px);
  }

  .tab-btn.active {
    background: white;
    color: #667eea;
    border-color: white;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.15);
  }

  .tab-icon { font-size: 18px; }

  .tab-label { display: none; }

  .tab-content {
    animation: fadeIn 0.3s ease-out;
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 20px;
    margin-bottom: 40px;
  }

  .stat-card {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 12px;
    padding: 24px;
    display: flex;
    align-items: center;
    gap: 20px;
    color: white;
    transition: all 0.3s;
  }

  .stat-card:hover {
    background: rgba(255, 255, 255, 0.15);
    transform: translateY(-5px);
  }

  .stat-icon { font-size: 40px; min-width: 60px; }

  .stat-label {
    margin: 0;
    opacity: 0.8;
    font-size: 14px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .stat-value {
    margin: 8px 0 0 0;
    font-size: 28px;
    font-weight: 700;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @media (min-width: 768px) {
    .tab-label { display: inline; }
  }

  @media (max-width: 768px) {
    .admin-container { padding: 20px; }
    .admin-header { flex-direction: column; gap: 16px; align-items: flex-start; }
    .header-right { justify-content: flex-start; }
    .admin-header h2 { font-size: 24px; }
    .tabs-navigation { gap: 8px; padding: 8px; }
    .tab-btn { padding: 10px 16px; font-size: 12px; }
    .dashboard-grid { grid-template-columns: 1fr; }
  }
</style>