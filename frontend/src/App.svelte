<script>
  import { onMount } from 'svelte'
  import Home from './views/Home.svelte'
  import Auth from './views/Auth.svelte'
  import Dashboard from './views/Dashboard.svelte'
  import AdminPanel from './views/AdminPanel.svelte'

  let currentView = 'home'
  let user = null
  let token = null

  onMount(() => {
    const savedToken = localStorage.getItem('token')
    const savedUser = localStorage.getItem('user')
    if (savedToken && savedUser) {
      token = savedToken
      user = JSON.parse(savedUser)
      currentView = 'dashboard'
    }
  })

  function handleLogin(event) {
    const { token: newToken, user: newUser } = event.detail
    token = newToken
    user = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
    currentView = 'dashboard'
  }

  function handleLogout() {
    token = null
    user = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    currentView = 'home'
  }

  function handleNavigateToAuth() {
    currentView = 'auth'
  }
</script>

{#if currentView === 'home'}
  <Home on:navigateToAuth={handleNavigateToAuth} />
{:else if currentView === 'auth'}
  <Auth on:login={handleLogin} />
{:else if currentView === 'dashboard'}
  {#if user?.role_id === 1}
    <AdminPanel {user} {token} on:logout={handleLogout} />
  {:else}
    <Dashboard {user} {token} on:logout={handleLogout} />
  {/if}
{/if}

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background: linear-gradient(135deg, #1e3c72 0%, #2a5298 50%, #7e22ce 100%);
    min-height: 100vh;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    color: #fff;
  }
  :global(#app) {
    width: 100%;
    min-height: 100vh;
  }
</style>
