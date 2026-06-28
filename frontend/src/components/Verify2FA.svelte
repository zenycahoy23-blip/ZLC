<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let userId
  export let email

  let token = ''
  let error = ''
  let loading = false

  async function handleVerify2FA() {
    error = ''

    if (!token) {
      error = '2FA code is required'
      return
    }

    loading = true

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8083'
      const response = await fetch(`${apiUrl}/api/verify-2fa`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, token })
      })

      const data = await response.json()

      if (!response.ok) {
        error = data.error || '2FA verification failed'
        return
      }

      dispatch('success', {
        token: data.token,
        user: {
          user_id: userId,
          email,
          role_id: data.role_id
        }
      })
    } catch (err) {
      error = 'Network error. Please try again.'
    } finally {
      loading = false
    }
  }
</script>

<div class="verify-2fa-form">
  <h2>Two-Factor Authentication</h2>
  <p class="info">Enter the 6-digit code sent to {email}</p>

  {#if error}
    <div class="error-message">{error}</div>
  {/if}

  <input
    type="text"
    placeholder="000000"
    bind:value={token}
    maxlength="6"
    disabled={loading}
  />

  <button on:click={handleVerify2FA} disabled={loading}>
    {loading ? 'Verifying...' : 'Verify Code'}
  </button>
</div>

<style>
  .verify-2fa-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  h2 {
    color: #333;
    margin-bottom: 10px;
  }

  .info {
    color: #666;
    font-size: 14px;
  }

  input {
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
    transition: border 0.3s;
    text-align: center;
    letter-spacing: 4px;
    font-weight: 500;
  }

  input:focus {
    outline: none;
    border-color: #667eea;
  }

  input:disabled {
    background: #f5f5f5;
    cursor: not-allowed;
  }

  button {
    padding: 12px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.3s;
  }

  button:hover:not(:disabled) {
    background: #5568d3;
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .error-message {
    background: #f8d7da;
    color: #721c24;
    padding: 12px;
    border-radius: 6px;
    font-size: 14px;
  }
</style>
