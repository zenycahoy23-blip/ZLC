<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  let email = ''
  let password = ''
  let error = ''
  let loading = false

  async function handleLogin() {
    error = ''

    if (!email || !password) {
      error = 'Email and password are required'
      return
    }

    loading = true

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8083'
      const response = await fetch(`${apiUrl}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      })

      const data = await response.json()

      if (!response.ok) {
        error = data.error || 'Login failed'
        return
      }

      dispatch('success', {
        email,
        user_id: data.user_id
      })
    } catch (err) {
      error = 'Network error. Please try again.'
    } finally {
      loading = false
    }
  }
</script>

<div class="login-form">
  <h2>Login</h2>

  {#if error}
    <div class="error-message">{error}</div>
  {/if}

  <input
    type="email"
    placeholder="Email"
    bind:value={email}
    disabled={loading}
  />

  <input
    type="password"
    placeholder="Password"
    bind:value={password}
    disabled={loading}
  />

  <button on:click={handleLogin} disabled={loading}>
    {loading ? 'Logging in...' : 'Login'}
  </button>

  <p class="switch-auth">
    Don't have an account?
    <button type="button" on:click={() => dispatch('switchToRegister')} class="link-btn">
      Register here
    </button>
  </p>
</div>

<style>
  .login-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  h2 {
    color: #333;
    margin-bottom: 10px;
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

  button:not(.link-btn) {
    padding: 12px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.3s;
  }

  button:not(.link-btn):hover:not(:disabled) {
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

  .switch-auth {
    text-align: center;
    font-size: 14px;
    color: #666;
  }

  .link-btn {
    background: none;
    border: none;
    color: #667eea;
    cursor: pointer;
    text-decoration: underline;
    padding: 0;
    font-size: 14px;
  }
</style>
