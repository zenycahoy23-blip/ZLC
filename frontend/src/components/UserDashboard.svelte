<script>
  import { createEventDispatcher } from 'svelte'
  import { onMount } from 'svelte'

  const dispatch = createEventDispatcher()

  export let user
  export let token

  let userData = null
  let loading = false
  let error = ''
  let successMessage = ''
  let currentTab = 'overview'
  let showDeleteConfirm = false
  let deleteConfirmText = ''

  onMount(() => {
    fetchDashboard()
  })

  async function fetchDashboard() {
    loading = true
    error = ''
    try {
      const response = await fetch('http://localhost:8083/api/user/dashboard', {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      const data = await response.json()
      userData = data
    } catch (err) {
      error = 'Failed to load dashboard'
    } finally {
      loading = false
    }
  }

  function getRoleName(roleId) {
    const roles = { 1: 'Admin', 2: 'Moderator', 3: 'User' }
    return roles[roleId] || 'Unknown'
  }

  function handleTabChange(tab) {
    currentTab = tab
  }

  function handleDeleteClick() {
    showDeleteConfirm = true
    deleteConfirmText = ''
  }

  async function handleDeleteAccount() {
    if (deleteConfirmText !== 'DELETE') {
      error = 'Please type "DELETE" to confirm'
      return
    }

    try {
      loading = true
      const response = await fetch('http://localhost:8083/api/delete-account', {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ user_id: userData.user_id })
      })

      if (response.ok) {
        successMessage = 'Account deleted successfully. Redirecting...'
        setTimeout(() => {
          dispatch('logout')
        }, 2000)
      } else {
        error = 'Failed to delete account'
      }
    } catch (err) {
      error = 'Error deleting account: ' + err.message
    } finally {
      loading = false
    }
  }

  function handleCancelDelete() {
    showDeleteConfirm = false
    deleteConfirmText = ''
  }

  function copyToClipboard(text) {
    navigator.clipboard.writeText(text)
    successMessage = 'Copied to clipboard!'
    setTimeout(() => {
      successMessage = ''
    }, 2000)
  }
</script>

<div class="user-dashboard-container">
  {#if successMessage}
    <div class="success-banner">
      <span class="icon">✅</span>
      <p>{successMessage}</p>
    </div>
  {/if}

  {#if error}
    <div class="error-banner">
      <span class="icon">❌</span>
      <p>{error}</p>
      <button class="close-btn" on:click={() => (error = '')}>×</button>
    </div>
  {/if}

  <!-- Header Section -->
  <div class="dashboard-header">
    <div class="header-content">
      <div class="header-icon">👤</div>
      <div class="header-text">
        <h1>Account Dashboard</h1>
        <p>Manage your account settings and security</p>
      </div>
    </div>
    <div class="header-status">
      <div class="status-badge active">
        <span class="status-dot"></span>
        Active Account
      </div>
    </div>
  </div>

  <!-- Tab Navigation -->
  <div class="tab-navigation">
    <button
      class="tab-button {currentTab === 'overview' ? 'active' : ''}"
      on:click={() => handleTabChange('overview')}
    >
      <span class="tab-icon">📊</span>
      <span>Overview</span>
    </button>
    <button
      class="tab-button {currentTab === 'security' ? 'active' : ''}"
      on:click={() => handleTabChange('security')}
    >
      <span class="tab-icon">🔒</span>
      <span>Security</span>
    </button>
    <button
      class="tab-button {currentTab === 'settings' ? 'active' : ''}"
      on:click={() => handleTabChange('settings')}
    >
      <span class="tab-icon">⚙️</span>
      <span>Settings</span>
    </button>
  </div>

  <!-- Tab Content -->
  <div class="tab-content">
    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
        <p>Loading your dashboard...</p>
      </div>
    {:else if currentTab === 'overview'}
      <!-- Overview Tab -->
      <div class="content-section">
        <div class="cards-grid">
          <!-- Profile Card -->
          <div class="content-card profile-card">
            <div class="card-icon">👤</div>
            <div class="card-header">
              <h3>Account Information</h3>
              <span class="card-badge">Primary</span>
            </div>

            <div class="card-content">
              <div class="info-row">
                <span class="info-label">Email Address</span>
                <div class="info-value-with-action">
                  <span class="info-value">{userData?.email || 'N/A'}</span>
                  <button
                    class="copy-btn"
                    on:click={() => copyToClipboard(userData?.email)}
                    title="Copy email"
                  >
                    📋
                  </button>
                </div>
              </div>

              <div class="info-row">
                <span class="info-label">User ID</span>
                <div class="info-value-with-action">
                  <span class="info-value">#{userData?.user_id || 'N/A'}</span>
                  <button
                    class="copy-btn"
                    on:click={() => copyToClipboard(userData?.user_id)}
                    title="Copy ID"
                  >
                    📋
                  </button>
                </div>
              </div>

              <div class="info-row">
                <span class="info-label">Account Role</span>
                <span class="role-badge">{getRoleName(userData?.role_id)}</span>
              </div>

              <div class="info-row">
                <span class="info-label">Account Status</span>
                <span class="status-indicator">🟢 Active</span>
              </div>
            </div>
          </div>

          <!-- Activity Card -->
          <div class="content-card activity-card">
            <div class="card-icon">📈</div>
            <div class="card-header">
              <h3>Account Activity</h3>
              <span class="card-badge">Recent</span>
            </div>

            <div class="card-content">
              <div class="activity-row">
                <span class="activity-icon">🔓</span>
                <div>
                  <p class="activity-label">Last Login</p>
                  <p class="activity-value">
                    {userData?.last_login ? new Date(userData.last_login).toLocaleString() : 'This is your first login'}
                  </p>
                </div>
              </div>

              <div class="activity-row">
                <span class="activity-icon">⏰</span>
                <div>
                  <p class="activity-label">Account Created</p>
                  <p class="activity-value">Recently</p>
                </div>
              </div>

              <div class="activity-row">
                <span class="activity-icon">📊</span>
                <div>
                  <p class="activity-label">Account Standing</p>
                  <p class="activity-value">In Good Standing</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    {:else if currentTab === 'security'}
      <!-- Security Tab -->
      <div class="content-section">
        <div class="cards-grid">
          <!-- Security Features Card -->
          <div class="content-card security-card">
            <div class="card-icon">🔒</div>
            <div class="card-header">
              <h3>Security Features</h3>
              <span class="card-badge success">Protected</span>
            </div>

            <div class="card-content">
              <div class="security-feature">
                <div class="feature-check">✅</div>
                <div class="feature-info">
                  <h4>Password Hashing</h4>
                  <p>Your password is securely hashed using bcrypt</p>
                </div>
              </div>

              <div class="security-feature">
                <div class="feature-check">✅</div>
                <div class="feature-info">
                  <h4>Email Verification</h4>
                  <p>Your email address has been verified</p>
                </div>
              </div>

              <div class="security-feature">
                <div class="feature-check">✅</div>
                <div class="feature-info">
                  <h4>Phone Verification</h4>
                  <p>Your phone number has been verified via OTP</p>
                </div>
              </div>

              <div class="security-feature">
                <div class="feature-check">✅</div>
                <div class="feature-info">
                  <h4>Two-Factor Authentication</h4>
                  <p>2FA is enabled for all login attempts</p>
                </div>
              </div>

              <div class="security-feature">
                <div class="feature-check">✅</div>
                <div class="feature-info">
                  <h4>Login History Tracking</h4>
                  <p>All login attempts are logged for security</p>
                </div>
              </div>

              <div class="security-feature">
                <div class="feature-check">✅</div>
                <div class="feature-info">
                  <h4>IP Address Logging</h4>
                  <p>Login locations are tracked and monitored</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    {:else if currentTab === 'settings'}
      <!-- Settings Tab -->
      <div class="content-section">
        <div class="settings-container">
          <!-- Danger Zone -->
          <div class="danger-zone-card">
            <div class="danger-header">
              <div class="danger-icon">⚠️</div>
              <div>
                <h3>Danger Zone</h3>
                <p>Irreversible and destructive actions</p>
              </div>
            </div>

            <div class="danger-content">
              <div class="danger-action">
                <div class="action-description">
                  <h4>Delete Account</h4>
                  <p>Permanently delete your account and all associated data. This action cannot be undone.</p>
                </div>
                <button class="btn-danger" on:click={handleDeleteClick}>
                  Delete Account
                </button>
              </div>
            </div>
          </div>

          <!-- Delete Confirmation Modal -->
          {#if showDeleteConfirm}
            <div class="modal-overlay" on:click={handleCancelDelete}>
              <div class="modal-content" on:click|stopPropagation>
                <div class="modal-header">
                  <h2>Delete Account</h2>
                  <button class="modal-close" on:click={handleCancelDelete}>×</button>
                </div>

                <div class="modal-body">
                  <div class="warning-icon">⚠️</div>
                  <p class="warning-title">This action cannot be undone</p>
                  <p class="warning-text">
                    Deleting your account will permanently remove all your data, including:
                  </p>
                  <ul class="warning-list">
                    <li>Account profile information</li>
                    <li>Login history and audit logs</li>
                    <li>All authentication records</li>
                    <li>Activity tracking data</li>
                  </ul>

                  <div class="confirmation-input">
                    <label>To confirm deletion, type <strong>DELETE</strong></label>
                    <input
                      type="text"
                      placeholder="Type DELETE"
                      bind:value={deleteConfirmText}
                      on:keydown={(e) => {
                        if (e.key === 'Enter' && deleteConfirmText === 'DELETE') {
                          handleDeleteAccount()
                        }
                      }}
                    />
                  </div>
                </div>

                <div class="modal-footer">
                  <button class="btn-secondary" on:click={handleCancelDelete}>Cancel</button>
                  <button
                    class="btn-danger"
                    disabled={deleteConfirmText !== 'DELETE' || loading}
                    on:click={handleDeleteAccount}
                  >
                    {loading ? 'Deleting...' : 'Delete My Account'}
                  </button>
                </div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .user-dashboard-container {
    animation: fadeIn 0.4s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* Banners */
  .success-banner,
  .error-banner {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 20px;
    border-radius: 12px;
    margin-bottom: 24px;
    animation: slideIn 0.3s ease-out;
  }

  .success-banner {
    background: #d4edda;
    border: 1px solid #c3e6cb;
    color: #155724;
  }

  .error-banner {
    background: #f8d7da;
    border: 1px solid #f5c6cb;
    color: #721c24;
    justify-content: space-between;
  }

  .success-banner .icon,
  .error-banner .icon {
    font-size: 20px;
  }

  .close-btn {
    background: none;
    border: none;
    color: inherit;
    font-size: 24px;
    cursor: pointer;
    padding: 0;
  }

  /* Header */
  .dashboard-header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border-radius: 16px;
    padding: 32px;
    margin-bottom: 32px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 8px 32px rgba(102, 126, 234, 0.25);
  }

  .header-content {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .header-icon {
    font-size: 48px;
    background: rgba(255, 255, 255, 0.2);
    width: 80px;
    height: 80px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .header-text h1 {
    margin: 0 0 8px 0;
    font-size: 28px;
  }

  .header-text p {
    margin: 0;
    opacity: 0.9;
    font-size: 14px;
  }

  .header-status {
    display: flex;
    gap: 12px;
  }

  .status-badge {
    display: flex;
    align-items: center;
    gap: 8px;
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(10px);
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 13px;
    font-weight: 600;
    border: 1px solid rgba(255, 255, 255, 0.3);
  }

  .status-dot {
    width: 8px;
    height: 8px;
    background: #4ade80;
    border-radius: 50%;
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.5;
    }
  }

  /* Tab Navigation */
  .tab-navigation {
    display: flex;
    gap: 12px;
    margin-bottom: 32px;
    background: rgba(255, 255, 255, 0.05);
    padding: 12px;
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    flex-wrap: wrap;
  }

  .tab-button {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.7);
    border: 1px solid rgba(255, 255, 255, 0.2);
    padding: 12px 20px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.3s;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
  }

  .tab-button:hover {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.3);
    transform: translateY(-2px);
  }

  .tab-button.active {
    background: white;
    color: #667eea;
    border-color: white;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.15);
  }

  .tab-icon {
    font-size: 16px;
  }

  /* Content */
  .tab-content {
    animation: fadeIn 0.3s ease-out;
  }

  .content-section {
    margin-bottom: 32px;
  }

  .cards-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
    gap: 24px;
  }

  .content-card {
    background: white;
    border-radius: 12px;
    padding: 28px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    border: 1px solid rgba(0, 0, 0, 0.05);
    transition: all 0.3s;
    position: relative;
  }

  .content-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  }

  .card-icon {
    font-size: 32px;
    margin-bottom: 12px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 2px solid #f5f5f5;
  }

  .card-header h3 {
    margin: 0;
    font-size: 18px;
    color: #333;
  }

  .card-badge {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 11px;
    font-weight: 700;
    text-transform: uppercase;
  }

  .card-badge.success {
    background: linear-gradient(135deg, #4ade80 0%, #22c55e 100%);
  }

  .card-content {
    display: flex;
    flex-direction: column;
    gap: 18px;
  }

  /* Info Rows */
  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid #f5f5f5;
  }

  .info-row:last-child {
    border-bottom: none;
  }

  .info-label {
    color: #999;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .info-value-with-action {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .info-value {
    color: #333;
    font-weight: 600;
    font-size: 14px;
  }

  .copy-btn {
    background: none;
    border: none;
    font-size: 16px;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
    transition: all 0.2s;
  }

  .copy-btn:hover {
    background: #f5f5f5;
  }

  .role-badge {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 12px;
    font-weight: 600;
  }

  .status-indicator {
    color: #4ade80;
    font-weight: 600;
    font-size: 14px;
  }

  /* Activity Rows */
  .activity-row {
    display: flex;
    gap: 12px;
    padding: 12px;
    background: #f8f9fa;
    border-radius: 8px;
  }

  .activity-icon {
    font-size: 20px;
    flex-shrink: 0;
  }

  .activity-label {
    margin: 0;
    color: #999;
    font-size: 12px;
    text-transform: uppercase;
    font-weight: 600;
  }

  .activity-value {
    margin: 4px 0 0 0;
    color: #333;
    font-size: 14px;
    font-weight: 500;
  }

  /* Security Features */
  .security-feature {
    display: flex;
    gap: 12px;
    padding: 16px;
    background: #f8f9fa;
    border-radius: 8px;
  }

  .feature-check {
    font-size: 24px;
    flex-shrink: 0;
  }

  .feature-info h4 {
    margin: 0 0 4px 0;
    color: #333;
    font-size: 14px;
    font-weight: 600;
  }

  .feature-info p {
    margin: 0;
    color: #999;
    font-size: 13px;
  }

  /* Loading State */
  .loading-state {
    text-align: center;
    padding: 60px 20px;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(255, 255, 255, 0.2);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 20px;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  /* Settings */
  .settings-container {
    max-width: 600px;
  }

  .danger-zone-card {
    background: white;
    border-radius: 12px;
    padding: 28px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    border: 2px solid #fee;
    border-left: 4px solid #dc3545;
  }

  .danger-header {
    display: flex;
    gap: 16px;
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 2px solid #fee;
  }

  .danger-icon {
    font-size: 32px;
  }

  .danger-header h3 {
    margin: 0 0 4px 0;
    color: #dc3545;
    font-size: 18px;
  }

  .danger-header p {
    margin: 0;
    color: #999;
    font-size: 13px;
  }

  .danger-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .danger-action {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 20px;
  }

  .action-description h4 {
    margin: 0 0 4px 0;
    color: #333;
    font-size: 16px;
    font-weight: 600;
  }

  .action-description p {
    margin: 0;
    color: #999;
    font-size: 13px;
  }

  /* Buttons */
  .btn-danger {
    background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
    white-space: nowrap;
  }

  .btn-danger:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(220, 53, 69, 0.3);
  }

  .btn-danger:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Modal */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
  }

  .modal-content {
    background: white;
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    max-width: 500px;
    width: 90%;
    animation: slideUp 0.3s ease-out;
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 24px;
    border-bottom: 1px solid #f5f5f5;
  }

  .modal-header h2 {
    margin: 0;
    font-size: 20px;
    color: #333;
  }

  .modal-close {
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: #999;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .modal-close:hover {
    color: #333;
  }

  .modal-body {
    padding: 24px;
    text-align: center;
  }

  .warning-icon {
    font-size: 48px;
    margin-bottom: 16px;
    display: block;
  }

  .warning-title {
    margin: 0 0 12px 0;
    color: #333;
    font-size: 18px;
    font-weight: 600;
  }

  .warning-text {
    margin: 0 0 16px 0;
    color: #999;
    font-size: 14px;
  }

  .warning-list {
    text-align: left;
    background: #fef5f5;
    border-left: 3px solid #dc3545;
    padding: 16px;
    border-radius: 6px;
    margin: 16px 0;
    list-style: none;
    color: #333;
    font-size: 13px;
  }

  .warning-list li {
    padding: 6px 0;
  }

  .warning-list li:before {
    content: '• ';
    color: #dc3545;
    font-weight: 700;
    margin-right: 8px;
  }

  .confirmation-input {
    text-align: left;
    margin-top: 20px;
  }

  .confirmation-input label {
    display: block;
    margin-bottom: 8px;
    color: #333;
    font-size: 13px;
    font-weight: 600;
  }

  .confirmation-input input {
    width: 100%;
    padding: 10px 12px;
    border: 2px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
    font-family: monospace;
    transition: all 0.2s;
    box-sizing: border-box;
  }

  .confirmation-input input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .modal-footer {
    display: flex;
    gap: 12px;
    padding: 24px;
    border-top: 1px solid #f5f5f5;
  }

  .btn-secondary {
    flex: 1;
    background: #f5f5f5;
    color: #333;
    border: none;
    padding: 10px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-secondary:hover {
    background: #eee;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .dashboard-header {
      flex-direction: column;
      text-align: center;
      gap: 20px;
    }

    .header-content {
      flex-direction: column;
    }

    .header-icon {
      width: 60px;
      height: 60px;
      font-size: 32px;
    }

    .header-text h1 {
      font-size: 24px;
    }

    .cards-grid {
      grid-template-columns: 1fr;
    }

    .danger-action {
      flex-direction: column;
    }

    .tab-navigation {
      overflow-x: auto;
    }

    .info-row {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }
  }
</style>
