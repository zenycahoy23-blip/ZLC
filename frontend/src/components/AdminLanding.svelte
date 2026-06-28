<script>
  import { onMount, createEventDispatcher } from 'svelte'
  import AdminAuditDashboard from './AdminAuditDashboard.svelte'

  export let user
  export let token

  const dispatch = createEventDispatcher()

  let activeTab = 'audit-dashboard'
  let stats = { totalUsers: 0, verifiedUsers: 0, auditLogs: 0 }
  let loading = true
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8082'

  onMount(async () => {
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
        auditLogs: statsData?.audit_logs || 0
      }
    } catch (err) {
      console.error('Failed to load admin stats', err)
    } finally {
      loading = false
    }
  })

  function handleLogout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    dispatch('logout')
  }
</script>

<div class="admin-dashboard">
  <div class="admin-header">
    <div class="header-content">
      <h1>👑 Admin Control Panel</h1>
      <p>Manage users, view audit logs, and monitor system activity</p>
    </div>
    <div class="admin-info">
      <span class="admin-badge">ADMIN</span>
      <p>{user?.email}</p>
      <button class="logout-btn" on:click={handleLogout}>
        🚪 Logout
      </button>
    </div>
  </div>

  <div class="stats-grid">
    <div class="stat-card">
      <span class="stat-icon">👥</span>
      <div>
        <p class="stat-label">Total Users</p>
        <p class="stat-value">{loading ? '...' : stats.totalUsers}</p>
      </div>
    </div>
    <div class="stat-card">
      <span class="stat-icon">✅</span>
      <div>
        <p class="stat-label">Verified Users</p>
        <p class="stat-value">{loading ? '...' : stats.verifiedUsers}</p>
      </div>
    </div>
    <div class="stat-card">
      <span class="stat-icon">📝</span>
      <div>
        <p class="stat-label">Audit Logs</p>
        <p class="stat-value">{loading ? '...' : stats.auditLogs}</p>
      </div>
    </div>
    <div class="stat-card">
      <span class="stat-icon">⚙️</span>
      <div>
        <p class="stat-label">System Status</p>
        <p class="stat-value" style="color: #28a745;">Online</p>
      </div>
    </div>
  </div>

  <div class="tabs">
    <button class={activeTab === 'audit-dashboard' ? 'active' : ''} on:click={() => activeTab = 'audit-dashboard'}>
      📊 Audit Dashboard
    </button>
    <button class={activeTab === 'user-activity' ? 'active' : ''} on:click={() => activeTab = 'user-activity'}>
      👤 User Activity
    </button>
  </div>

  {#if activeTab === 'audit-dashboard'}
    <AdminAuditDashboard {token} />
  {:else if activeTab === 'user-activity'}
    <p class="no-selection">Select a user to view their activity</p>
  {/if}
</div>

<style>
  .admin-dashboard {
    animation: slideIn 0.5s ease-out;
  }

  @keyframes slideIn {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .admin-header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 40px;
    border-radius: 16px;
    margin-bottom: 40px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 10px 40px rgba(102, 126, 234, 0.3);
  }

  .header-content h1 {
    margin: 0 0 10px 0;
    font-size: 36px;
    font-weight: 700;
  }

  .header-content p {
    margin: 0;
    opacity: 0.9;
  }

  .admin-info {
    text-align: right;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 8px;
  }

  .admin-badge {
    background: rgba(255, 255, 255, 0.2);
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 700;
  }

  .admin-info p {
    margin: 0;
    font-size: 14px;
  }

  .logout-btn {
    background: rgba(255, 255, 255, 0.15);
    color: white;
    border: 2px solid rgba(255, 255, 255, 0.4);
    padding: 8px 18px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 600;
    transition: all 0.2s;
  }

  .logout-btn:hover {
    background: rgba(255, 255, 255, 0.3);
    border-color: white;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
    margin-bottom: 40px;
  }

  .stat-card {
    background: white;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
    display: flex;
    gap: 16px;
    align-items: center;
    transition: all 0.3s;
  }

  .stat-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  }

  .stat-icon {
    font-size: 32px;
  }

  .stat-label {
    margin: 0;
    color: #999;
    font-size: 12px;
    text-transform: uppercase;
    font-weight: 600;
  }

  .stat-value {
    margin: 4px 0 0 0;
    font-size: 28px;
    font-weight: 700;
    color: #333;
  }

  .tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 30px;
    border-bottom: 2px solid #eee;
  }

  .tabs button {
    background: none;
    border: none;
    padding: 16px 20px;
    color: #666;
    font-weight: 600;
    cursor: pointer;
    border-bottom: 3px solid transparent;
    transition: all 0.3s;
    margin-bottom: -2px;
  }

  .tabs button.active {
    color: #667eea;
    border-bottom-color: #667eea;
  }

  .no-selection {
    text-align: center;
    color: #999;
    padding: 40px 20px;
    background: white;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  }
</style>