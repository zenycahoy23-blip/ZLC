<script>
  import { createEventDispatcher } from 'svelte'
  import Register from '../components/Register.svelte'
  import Login from '../components/Login.svelte'
  import VerifyEmail from '../components/VerifyEmail.svelte'
  import VerifyPhone from '../components/VerifyPhone.svelte'
  import Verify2FA from '../components/Verify2FA.svelte'

  const dispatch = createEventDispatcher()

  let authStep = 'login'
  let tempUserData = null

  function handleRegisterClick() {
    authStep = 'register'
  }

  function handleLoginClick() {
    authStep = 'login'
  }

  function handleRegisterSuccess(event) {
    tempUserData = event.detail
    authStep = 'verify-email'
  }

  function handleEmailVerified(event) {
    tempUserData = { ...tempUserData, ...event.detail }
    authStep = 'verify-phone'
  }

  function handlePhoneVerified(event) {
    tempUserData = { ...tempUserData, ...event.detail }
    authStep = 'login'
  }

  function handleLoginSuccess(event) {
    tempUserData = event.detail
    authStep = 'verify-2fa'
  }

  function handle2FASuccess(event) {
    dispatch('login', event.detail)
  }
</script>

<div class="auth-container">
  <div class="auth-header">
    <div class="logo-section">
      <span class="logo-icon">🔄</span>
      <h1>Uno Reverse</h1>
      <p class="tagline">Secure Authentication & RBAC System</p>
    </div>
  </div>

  <div class="auth-card">
    {#if authStep === 'register'}
      <Register on:success={handleRegisterSuccess} on:switchToLogin={handleLoginClick} />
    {:else if authStep === 'login'}
      <Login on:success={handleLoginSuccess} on:switchToRegister={handleRegisterClick} />
    {:else if authStep === 'verify-email'}
      <VerifyEmail email={tempUserData?.email} on:success={handleEmailVerified} />
    {:else if authStep === 'verify-phone'}
      <VerifyPhone userId={tempUserData?.user_id} on:success={handlePhoneVerified} />
    {:else if authStep === 'verify-2fa'}
      <Verify2FA userId={tempUserData?.user_id} email={tempUserData?.email} on:success={handle2FASuccess} />
    {/if}
  </div>
</div>

<style>
  .auth-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    padding: 20px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    animation: fadeIn 0.6s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .auth-header {
    text-align: center;
    margin-bottom: 40px;
    color: white;
  }

  .logo-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }

  .logo-icon {
    font-size: 56px;
    display: block;
    animation: spin 3s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  h1 {
    margin: 0;
    font-size: 48px;
    font-weight: 700;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  }

  .tagline {
    margin: 0;
    font-size: 16px;
    opacity: 0.9;
  }

  .auth-card {
    background: white;
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    width: 100%;
    max-width: 420px;
    padding: 48px 40px;
    animation: slideUp 0.6s ease-out;
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
</style>
