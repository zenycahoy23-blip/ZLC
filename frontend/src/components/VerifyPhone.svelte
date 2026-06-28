<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let userId

  let otp = ''
  let error = ''
  let loading = false
  let showResend = false
  let resendLoading = false
  let cooldownSeconds = 0

  async function handleVerifyPhoneOTP() {
    error = ''

    if (!otp || otp.length !== 6) {
      error = 'Please enter a valid 6-digit OTP'
      return
    }

    loading = true

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8084'
      const response = await fetch(`${apiUrl}/api/verify-phone-otp`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, otp })
      })

      const data = await response.json()

      if (!response.ok) {
        error = data.error || 'OTP verification failed'
        return
      }

      dispatch('success', { userId })
    } catch (err) {
      error = 'Network error. Please try again.'
    } finally {
      loading = false
    }
  }

  async function handleResendOTP() {
    error = ''
    resendLoading = true

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8084'
      const response = await fetch(`${apiUrl}/api/resend-phone-otp`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId })
      })

      const data = await response.json()

      if (!response.ok) {
        error = data.error || 'Failed to resend OTP'
        return
      }

      showResend = false
      otp = ''
      
      // Show cooldown timer
      cooldownSeconds = 30
      const countdown = setInterval(() => {
        cooldownSeconds--
        if (cooldownSeconds <= 0) {
          clearInterval(countdown)
          showResend = true
        }
      }, 1000)

      error = 'OTP resent successfully!'
    } catch (err) {
      error = 'Network error. Please try again.'
    } finally {
      resendLoading = false
    }
  }

  function handleKeyPress(e) {
    if (e.key === 'Enter') {
      handleVerifyPhoneOTP()
    }
  }
</script>

<div class="verify-phone-form">
  <h2>Verify Phone Number</h2>

  <p class="subtitle">Enter the OTP sent to your email</p>

  {#if error}
    <div class="error-message">{error}</div>
  {/if}

  <input
    type="text"
    placeholder="Enter 6-digit OTP"
    bind:value={otp}
    maxlength="6"
    disabled={loading}
    on:keypress={handleKeyPress}
  />

  <button on:click={handleVerifyPhoneOTP} disabled={loading || otp.length !== 6}>
    {loading ? 'Verifying...' : 'Verify OTP'}
  </button>

  <div class="resend-section">
    {#if showResend}
      <button
        type="button"
        on:click={handleResendOTP}
        disabled={resendLoading}
        class="resend-btn"
      >
        {resendLoading ? 'Sending...' : 'Resend OTP'}
      </button>
    {:else if cooldownSeconds > 0}
      <p class="cooldown">Resend available in {cooldownSeconds}s</p>
    {:else}
      <button
        type="button"
        on:click={handleResendOTP}
        disabled={resendLoading}
        class="resend-btn"
      >
        {resendLoading ? 'Sending...' : 'Resend OTP'}
      </button>
    {/if}
  </div>
</div>

<style>
  .verify-phone-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  h2 {
    color: #333;
    margin-bottom: 10px;
  }

  .subtitle {
    font-size: 14px;
    color: #666;
    margin-bottom: 10px;
  }

  input {
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 24px;
    text-align: center;
    letter-spacing: 5px;
    font-weight: 600;
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

  button:not(.resend-btn) {
    padding: 12px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.3s;
  }

  button:not(.resend-btn):hover:not(:disabled) {
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

  .resend-section {
    text-align: center;
    margin-top: 10px;
  }

  .resend-btn {
    background: none;
    border: 1px solid #667eea;
    color: #667eea;
    padding: 10px 16px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.3s;
  }

  .resend-btn:hover:not(:disabled) {
    background: #f0f0ff;
  }

  .cooldown {
    color: #999;
    font-size: 14px;
    margin: 0;
  }
</style>
