<script>
  import { onMount, createEventDispatcher } from 'svelte'

  export let user
  export let token

  const dispatch = createEventDispatcher()

  let logs = []
  let users = []
  let activeTab = 'logs'
  let loading = false

  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8082'

  onMount(() => {
    fetchLogs()
  })

  async function fetchLogs() {
    loading = true
    try {
      const response = await fetch(`${apiUrl}/api/admin/logs`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      const data = await response.json()
      logs = data || []
    } catch (err) {
      console.error('Error fetching logs:', err)
    } finally {
      loading = false
    }
  }

  async function fetchUsers() {
    loading = true
    try {
      const response = await fetch(`${apiUrl}/api/admin/users`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      const data = await response.json()
      users = data || []
    } catch (err) {
      console.error('Error fetching users:', err)
    } finally {
      loading = false
    }
  }

  function switchTab(tab) {
    activeTab = tab
    if (tab === 'logs') fetchLogs()
    else if (tab === 'users') fetchUsers()
  }

  function getRoleName(roleId) {
    const roles = { 1: 'Admin', 2: 'Moderator', 3: 'User' }
    return roles[roleId] || 'Unknown'
  }

  function handleLogout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    dispatch('logout')
  }
</script>

<div class="admin-dashboard">
  <!-- Header with logout -->
  <div class="dashboard-header">
    <div class="header-left">
      <h1>👑 Admin Dashboard</h1>
      <p>Welcome, <strong>{user?.email}</strong></p>
    </div>
    <div class="header-right">
      <span class="admin-badge">ADMIN</span>
      <button class="logout-btn" on:click={handleLogout}>
        🚪 Logout
      </button>
    </div>
  </div>

  <!-- Tabs -->
  <div class="tabs">
    <button class={activeTab === 'logs' ? 'active' : ''} on:click={() => switchTab('logs')}>
      📋 Audit Logs
    </button>
    <button class={activeTab === 'users' ? 'active' : ''} on:click={() => switchTab('users')}>
      👥 Users
    </button>
  </div>

  <!-- Audit Logs Tab -->
  {#if activeTab === 'logs'}
    <div class="content">
      <div class="content-header">
        <h2>Audit Logs</h2>
        <button class="refresh-btn" on:click={fetchLogs}>🔄 Refresh</button>
      </div>
      {#if loading}
        <div class="loading">Loading...</div>
      {:else if logs.length === 0}
        <div class="empty">No logs found</div>
      {:else}
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>User ID</th>
                <th>Action</th>
                <th>Timestamp</th>
                <th>IP Address</th>
                <th>Details</th>
              </tr>
            </thead>
            <tbody>
              {#each logs as log}
                <tr>
                  <td>{log.user_id}</td>
                  <td><span class="action-badge">{log.action}</span></td>
                  <td>{new Date(log.timestamp).toLocaleString()}</td>
                  <td>{log.ip_address}</td>
                  <td>{log.details || '-'}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Users Tab -->
  {#if activeTab === 'users'}
    <div class="content">
      <div class="content-header">
        <h2>User Management</h2>
        <button class="refresh-btn" on:click={fetchUsers}>🔄 Refresh</button>
      </div>
      {#if loading}
        <div class="loading">Loading...</div>
      {:else if users.length === 0}
        <div class="empty">No users found</div>
      {:else}
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Email</th>
                <th>Role</th>
                <th>Email Verified</th>
                <th>Phone Verified</th>
                <th>Created</th>
              </tr>
            </thead>
            <tbody>
              {#each users as u}
                <tr>
                  <td>{u.id}</td>
                  <td>{u.email}</td>
                  <td>{getRoleName(u.role_id)}</td>
                  <td>
                    <span class="badge {u.email_verified ? 'success' : 'danger'}">
                      {u.email_verified ? '✓ Yes' : '✗ No'}
                    </span>
                  </td>
                  <td>
                    <span class="badge {u.phone_verified ? 'success' : 'danger'}">
                      {u.phone_verified ? '✓ Yes' : '✗ No'}
                    </span>
                  </td>
                  <td>{new Date(u.created_at).toLocaleString()}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .admin-dashboard {
    background: white;
    border-radius: 12px;
    padding: 30px;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  }

  .dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
    padding-bottom: 20px;
    border-bottom: 2px solid #eee;
  }

  .header-left h1 {
    margin: 0 0 5px 0;
    font-size: 28px;
    color: #333;
  }

  .header-left p {
    margin: 0;
    color: #666;
    font-size: 14px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .admin-badge {
    background: #667eea;
    color: white;
    padding: 6px 14px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 700;
  }

  .logout-btn {
    background: #f8d7da;
    color: #721c24;
    border: 2px solid #f5c6cb;
    padding: 8px 18px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 600;
    transition: all 0.2s;
  }

  .logout-btn:hover {
    background: #dc3545;
    color: white;
    border-color: #dc3545;
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
    padding: 12px 20px;
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

  .content-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .content-header h2 {
    margin: 0;
    color: #333;
  }

  .refresh-btn {
    background: #667eea;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 600;
    transition: background 0.2s;
  }

  .refresh-btn:hover { background: #5568d3; }

  .loading, .empty {
    text-align: center;
    padding: 40px;
    color: #999;
    background: #f9f9f9;
    border-radius: 8px;
  }

  .table-container { overflow-x: auto; }

  table { width: 100%; border-collapse: collapse; }

  th {
    background: #f5f5f5;
    padding: 12px;
    text-align: left;
    font-weight: 600;
    color: #333;
    border-bottom: 2px solid #ddd;
  }

  td {
    padding: 12px;
    border-bottom: 1px solid #eee;
    color: #666;
    font-size: 14px;
  }

  tr:hover { background: #f9f9f9; }

  .badge {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
  }

  .badge.success { background: #d4edda; color: #155724; }
  .badge.danger { background: #f8d7da; color: #721c24; }

  .action-badge {
    background: #e3f2fd;
    color: #1976d2;
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 600;
  }
</style>