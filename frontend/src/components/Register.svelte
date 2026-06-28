<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  let email = ''
  let password = ''
  let confirmPassword = ''
  let phoneNumber = ''
  let error = ''
  let loading = false

  async function handleRegister() {
    error = ''

    if (!email || !password || !confirmPassword || !phoneNumber) {
      error = 'All fields are required'
      return
    }

    if (password !== confirmPassword) {
      error = 'Passwords do not match'
      return
    }

    if (password.length < 6) {
      error = 'Password must be at least 6 characters'
      return
    }

    if (phoneNumber.length < 10) {
      error = 'Please enter a valid phone number (at least 10 digits)'
      return
    }

    loading = true

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8084'
      const response = await fetch(`${apiUrl}/api/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, phone_number: phoneNumber })
      })

      const data = await response.json()

      if (!response.ok) {
        error = data.error || 'Registration failed'
        return
      }

      dispatch('success', { email })
    } catch (err) {
      error = 'Network error. Please try again.'
    } finally {
      loading = false
    }
  }
</script>

<div class="register-form">
  <h2>Create Account</h2>

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
    type="tel"
    placeholder="Phone Number (e.g., +1234567890)"
    bind:value={phoneNumber}
    disabled={loading}
  />

  <input
    type="password"
    placeholder="Password"
    bind:value={password}
    disabled={loading}
  />

  <input
    type="password"
    placeholder="Confirm Password"
    bind:value={confirmPassword}
    disabled={loading}
  />

  <button on:click={handleRegister} disabled={loading}>
    {loading ? 'Creating...' : 'Register'}
  </button>

  <p class="switch-auth">
    Already have an account?
    <button type="button" on:click={() => dispatch('switchToLogin')} class="link-btn">
      Login here
    </button>
  </p>
</div>

<style>
  .register-form {
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
