<script>
  import { onMount } from 'svelte'

  export let userId
  let activity = null
  let loading = true
  let error = ''
  let token = ''
  let activeTab = 'login_history'

  onMount(async () => {
    token = localStorage.getItem('token')
    if (!token) {
      error = 'No authentication token found'
      loading = false
      return
    }

    try {
      const response = await fetch(`http://localhost:8082/api/admin/user-activity?user_id=${userId}`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      activity = await response.json()
      loading = false
    } catch (err) {
      error = 'Failed to load activity data'
      loading = false
    }
  })

  function formatDate(dateString) {
    return new Date(dateString).toLocaleString()
  }
</script>

<div class="user-activity">
  <h2>User Activity Audit</h2>

  {#if loading}
    <div class="loading">Loading activity data...</div>
  {:else if error}
    <div class="error">{error}</div>
  {:else if activity}
    <!-- Tabs -->
    <div class="tabs">
      <button 
        class="tab {activeTab === 'login_history' ? 'active' : ''}"
        on:click={() => activeTab = 'login_history'}
      >
        Login History ({activity.login_history?.length || 0})
      </button>
      <button 
        class="tab {activeTab === 'failed_logins' ? 'active' : ''}"
        on:click={() => activeTab = 'failed_logins'}
      >
        Failed Logins ({activity.failed_logins?.length || 0})
      </button>
      <button 
        class="tab {activeTab === 'password_changes' ? 'active' : ''}"
        on:click={() => activeTab = 'password_changes'}
      >
        Password Changes ({activity.password_changes?.length || 0})
      </button>
      <button 
        class="tab {activeTab === 'contact_changes' ? 'active' : ''}"
        on:click={() => activeTab = 'contact_changes'}
      >
        Contact Changes ({activity.contact_changes?.length || 0})
      </button>
      <button 
        class="tab {activeTab === 'profile_updates' ? 'active' : ''}"
        on:click={() => activeTab = 'profile_updates'}
      >
        Profile Updates ({activity.profile_updates?.length || 0})
      </button>
    </div>

    <!-- Tab Content -->
    <div class="tab-content">
      {#if activeTab === 'login_history' && activity.login_history}
        <div class="table-section">
          {#if activity.login_history.length === 0}
            <p class="no-data">No login history</p>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>Login Time</th>
                  <th>Logout Time</th>
                  <th>IP Address</th>
                  <th>Device</th>
                  <th>Browser</th>
                  <th>Location</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {#each activity.login_history as log}
                  <tr>
                    <td>{formatDate(log.login_time)}</td>
                    <td>{log.logout_time ? formatDate(log.logout_time) : '-'}</td>
                    <td><code>{log.ip_address}</code></td>
                    <td>{log.device_info}</td>
                    <td>{log.browser_info}</td>
                    <td>{log.location || '-'}</td>
                    <td>
                      <span class="status {log.status}">
                        {log.status}
                      </span>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
      {/if}

      {#if activeTab === 'failed_logins' && activity.failed_logins}
        <div class="table-section">
          {#if activity.failed_logins.length === 0}
            <p class="no-data">No failed login attempts</p>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>Time</th>
                  <th>IP Address</th>
                  <th>Device</th>
                  <th>Browser</th>
                  <th>Failure Reason</th>
                </tr>
              </thead>
              <tbody>
                {#each activity.failed_logins as log}
                  <tr>
                    <td>{formatDate(log.login_time)}</td>
                    <td><code>{log.ip_address}</code></td>
                    <td>{log.device_info}</td>
                    <td>{log.browser_info}</td>
                    <td><span class="danger">{log.failure_reason}</span></td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
      {/if}

      {#if activeTab === 'password_changes' && activity.password_changes}
        <div class="table-section">
          {#if activity.password_changes.length === 0}
            <p class="no-data">No password changes</p>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>Changed At</th>
                  <th>IP Address</th>
                  <th>Device</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {#each activity.password_changes as log}
                  <tr>
                    <td>{formatDate(log.changed_at)}</td>
                    <td><code>{log.ip_address}</code></td>
                    <td>{log.device_info}</td>
                    <td>
                      <span class="status {log.status}">
                        {log.status}
                      </span>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
      {/if}

      {#if activeTab === 'contact_changes' && activity.contact_changes}
        <div class="table-section">
          {#if activity.contact_changes.length === 0}
            <p class="no-data">No contact changes</p>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>Changed At</th>
                  <th>Type</th>
                  <th>Old Value</th>
                  <th>New Value</th>
                  <th>IP Address</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {#each activity.contact_changes as log}
                  <tr>
                    <td>{formatDate(log.changed_at)}</td>
                    <td><strong>{log.type}</strong></td>
                    <td><code>{log.old_value}</code></td>
                    <td><code>{log.new_value}</code></td>
                    <td><code>{log.ip_address}</code></td>
                    <td>
                      <span class="status {log.status}">
                        {log.status}
                      </span>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
      {/if}

      {#if activeTab === 'profile_updates' && activity.profile_updates}
        <div class="table-section">
          {#if activity.profile_updates.length === 0}
            <p class="no-data">No profile updates</p>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>Updated At</th>
                  <th>Field</th>
                  <th>Old Value</th>
                  <th>New Value</th>
                  <th>IP Address</th>
                </tr>
              </thead>
              <tbody>
                {#each activity.profile_updates as log}
                  <tr>
                    <td>{formatDate(log.updated_at)}</td>
                    <td><strong>{log.field}</strong></td>
                    <td><code>{log.old_value}</code></td>
                    <td><code>{log.new_value}</code></td>
                    <td><code>{log.ip_address}</code></td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .user-activity {
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 20px;
  }

  h2 {
    margin-top: 0;
    color: #333;
    margin-bottom: 20px;
  }

  .loading,
  .error {
    padding: 20px;
    text-align: center;
    border-radius: 8px;
  }

  .loading {
    background: #f0f0f0;
    color: #666;
  }

  .error {
    background: #f8d7da;
    color: #721c24;
  }

  /* Tabs */
  .tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
    flex-wrap: wrap;
    border-bottom: 2px solid #eee;
  }

  .tab {
    padding: 10px 15px;
    background: none;
    border: none;
    border-bottom: 3px solid transparent;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    color: #666;
    transition: all 0.2s;
  }

  .tab:hover {
    color: #333;
  }

  .tab.active {
    color: #667eea;
    border-bottom-color: #667eea;
  }

  /* Table */
  .table-section {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid #eee;
    font-size: 14px;
  }

  th {
    background: #f5f5f5;
    font-weight: 600;
    color: #333;
  }

  tr:hover {
    background: #f9f9f9;
  }

  code {
    background: #f5f5f5;
    padding: 2px 6px;
    border-radius: 3px;
    font-family: 'Courier New', monospace;
    font-size: 12px;
  }

  .status {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
  }

  .status.success {
    background: #d4edda;
    color: #155724;
  }

  .status.failed {
    background: #f8d7da;
    color: #721c24;
  }

  .status.pending {
    background: #fff3cd;
    color: #856404;
  }

  .danger {
    color: #dc3545;
  }

  .no-data {
    text-align: center;
    color: #999;
    padding: 30px;
  }

  @media (max-width: 768px) {
    .tabs {
      flex-direction: column;
    }

    table {
      font-size: 12px;
    }

    th,
    td {
      padding: 8px;
    }
  }
</style>
