<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let email

  let token = ''
  let error = ''
  let loading = false
  let resending = false

  async function handleVerify() {
    error = ''

    if (!token) {
      error = 'Token is required'
      return
    }

    loading = true

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8083'
      const response = await fetch(`${apiUrl}/api/verify-email`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, token })
      })

      const data = await response.json()

      if (!response.ok) {
        error = data.error || 'Verification failed'
        return
      }

      dispatch('success', { email, user_id: data.user_id })
    } catch (err) {
      error = 'Network error. Please try again.'
    } finally {
      loading = false
    }
  }

  async function handleResendToken() {
    error = ''
    resending = true

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8083'
      const response = await fetch(`${apiUrl}/api/resend-verification`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
      })

      const data = await response.json()

      if (!response.ok) {
        error = data.error || 'Failed to resend token'
        return
      }

      error = ''
      token = ''
    } catch (err) {
      error = 'Network error. Please try again.'
    } finally {
      resending = false
    }
  }
</script>

<div class="verify-form">
  <h2>Verify Email</h2>
  <p class="info">Enter the token sent to {email}</p>

  {#if error}
    <div class="error-message">{error}</div>
  {/if}

  <input
    type="text"
    placeholder="Verification Token"
    bind:value={token}
    disabled={loading}
  />

  <button on:click={handleVerify} disabled={loading}>
    {loading ? 'Verifying...' : 'Verify'}
  </button>

  <button on:click={handleResendToken} disabled={resending} class="resend-btn">
    {resending ? 'Sending...' : 'Send Again'}
  </button>
</div>

<style>
  .verify-form {
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

  .resend-btn {
    background: #6c757d;
  }

  .resend-btn:hover:not(:disabled) {
    background: #5a6268;
  }

  .error-message {
    background: #f8d7da;
    color: #721c24;
    padding: 12px;
    border-radius: 6px;
    font-size: 14px;
  }
</style>
