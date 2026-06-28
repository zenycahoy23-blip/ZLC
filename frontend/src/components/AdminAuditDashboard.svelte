<script>
  import { onMount } from 'svelte'

  export let token = ''

  let stats = null
  let loginTrends = null
  let users = null
  let loading = true
  let error = ''

  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8082'

  onMount(async () => {
    if (!token) {
      token = localStorage.getItem('token')
    }
    if (!token) {
      error = 'No authentication token found'
      loading = false
      return
    }
    await loadData()
  })

  async function loadData() {
    loading = true
    error = ''
    try {
      const [statsRes, trendsRes, usersRes] = await Promise.all([
        fetch(`${apiUrl}/api/admin/dashboard-stats`, {
          headers: { 'Authorization': `Bearer ${token}` }
        }),
        fetch(`${apiUrl}/api/admin/login-trends`, {
          headers: { 'Authorization': `Bearer ${token}` }
        }),
        fetch(`${apiUrl}/api/admin/all-users`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
      ])

      stats = await statsRes.json()
      loginTrends = await trendsRes.json()
      users = await usersRes.json()
    } catch (err) {
      error = 'Failed to load dashboard data'
    } finally {
      loading = false
    }
  }

  async function unlockAccount(userId) {
    try {
      const response = await fetch(`${apiUrl}/api/admin/unlock-account?user_id=${userId}`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      })
      if (response.ok) {
        await loadData()
      }
    } catch (err) {
      console.error('Error unlocking account', err)
    }
  }

  function getRoleName(roleId) {
    const roles = { 1: 'Admin', 2: 'Moderator', 3: 'User' }
    return roles[roleId] || 'Unknown'
  }
</script>

<div class="admin-dashboard">
  <div class="dashboard-header">
    <h1>Admin Audit Dashboard</h1>
    <button class="refresh-btn" on:click={loadData}>🔄 Refresh</button>
  </div>

  {#if loading}
    <div class="loading">Loading dashboard...</div>
  {:else if error}
    <div class="error">{error}</div>
  {:else}
    <!-- KPI Cards -->
    <div class="kpi-grid">
      <div class="kpi-card">
        <div class="kpi-label">Total Users</div>
        <div class="kpi-value">{stats?.total_users || 0}</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">New Users Today</div>
        <div class="kpi-value">{stats?.new_users_today || 0}</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">Active Users</div>
        <div class="kpi-value">{stats?.active_users || 0}</div>
      </div>
      <div class="kpi-card warning">
        <div class="kpi-label">Failed Logins</div>
        <div class="kpi-value">{stats?.failed_login_attempts || 0}</div>
      </div>
      <div class="kpi-card danger">
        <div class="kpi-label">Locked Accounts</div>
        <div class="kpi-value">{stats?.locked_accounts || 0}</div>
      </div>
      <div class="kpi-card info">
        <div class="kpi-label">Verification Pending</div>
        <div class="kpi-value">{stats?.verification_pending || 0}</div>
      </div>
    </div>

    <!-- Login Trends -->
    {#if loginTrends && loginTrends.length > 0}
      <div class="chart-section">
        <h2>Login Trends (Last 7 Days)</h2>
        <div class="chart-table">
          <table>
            <thead>
              <tr>
                <th>Date</th>
                <th>Total Logins</th>
                <th>Successful</th>
                <th>Failed</th>
              </tr>
            </thead>
            <tbody>
              {#each loginTrends as trend}
                <tr>
                  <td>{trend.date}</td>
                  <td>{trend.total}</td>
                  <td class="success">{trend.successful}</td>
                  <td class="danger">{trend.failed}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    {/if}

    <!-- Users Table -->
    {#if users && users.length > 0}
      <div class="users-section">
        <h2>All Users</h2>
        <div class="users-table">
          <table>
            <thead>
              <tr>
                <th>Email</th>
                <th>Phone</th>
                <th>Role</th>
                <th>Email Verified</th>
                <th>Phone Verified</th>
                <th>Status</th>
                <th>Failed Attempts</th>
                <th>Last Login</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each users as user}
                <tr>
                  <td>{user.email}</td>
                  <td>{user.phone_number || '-'}</td>
                  <td>{getRoleName(user.role_id)}</td>
                  <td>
                    <span class="badge {user.email_verified ? 'success' : 'danger'}">
                      {user.email_verified ? '✓' : '✗'}
                    </span>
                  </td>
                  <td>
                    <span class="badge {user.phone_verified ? 'success' : 'danger'}">
                      {user.phone_verified ? '✓' : '✗'}
                    </span>
                  </td>
                  <td>
                    <span class="badge {user.account_locked ? 'danger' : 'success'}">
                      {user.account_locked ? 'Locked' : 'Active'}
                    </span>
                  </td>
                  <td>{user.failed_login_attempts || 0}</td>
                  <td>{user.last_login ? new Date(user.last_login).toLocaleString() : '-'}</td>
                  <td>
                    {#if user.account_locked}
                      <button class="action-btn" on:click={() => unlockAccount(user.id)}>
                        Unlock
                      </button>
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    {:else}
      <div class="empty">No users found</div>
    {/if}
  {/if}
</div>

<style>
  .admin-dashboard {
    padding: 20px;
    max-width: 1400px;
    margin: 0 auto;
  }

  .dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
  }

  h1 {
    color: #333;
    font-size: 32px;
    margin: 0;
  }

  h2 {
    color: #555;
    margin-top: 0;
    margin-bottom: 15px;
    font-size: 20px;
  }

  .loading, .error, .empty {
    padding: 20px;
    text-align: center;
    border-radius: 8px;
    margin-bottom: 20px;
  }

  .loading { background: #f0f0f0; color: #666; }
  .error { background: #f8d7da; color: #721c24; }
  .empty { background: #f0f0f0; color: #999; padding: 40px; }

  .kpi-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 15px;
    margin-bottom: 30px;
  }

  .kpi-card {
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 20px;
    text-align: center;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    transition: transform 0.2s;
  }

  .kpi-card:hover { transform: translateY(-2px); }
  .kpi-card.warning { border-color: #ffc107; background: #fffbf0; }
  .kpi-card.danger { border-color: #dc3545; background: #fff5f5; }
  .kpi-card.info { border-color: #17a2b8; background: #f0f9fb; }

  .kpi-label { font-size: 14px; color: #999; margin-bottom: 10px; text-transform: uppercase; }
  .kpi-value { font-size: 32px; font-weight: bold; color: #333; }

  .chart-section, .users-section {
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 30px;
    overflow-x: auto;
  }

  table { width: 100%; border-collapse: collapse; }

  th, td { padding: 12px; text-align: left; border-bottom: 1px solid #eee; font-size: 14px; }
  th { background: #f5f5f5; font-weight: 600; color: #333; }
  tr:hover { background: #f9f9f9; }

  td.success { color: #28a745; font-weight: 600; }
  td.danger { color: #dc3545; font-weight: 600; }

  .badge {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
  }

  .badge.success { background: #d4edda; color: #155724; }
  .badge.danger { background: #f8d7da; color: #721c24; }

  .action-btn {
    padding: 6px 12px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
    font-weight: 600;
    transition: background 0.2s;
  }

  .action-btn:hover { background: #5568d3; }

  .refresh-btn {
    background: #667eea;
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    transition: background 0.2s;
  }

  .refresh-btn:hover { background: #5568d3; }
</style>