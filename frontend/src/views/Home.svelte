<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  function handleGetStarted() {
    dispatch('navigateToAuth')
  }
</script>

<div class="home-container">
  <!-- Hero Section -->
  <nav class="navbar">
    <div class="logo">
      <span class="logo-icon">🔐</span>
      <h2>ZLC</h2>
    </div>
    <button class="get-started-btn" on:click={handleGetStarted}>Get Started</button>
  </nav>

  <section class="hero">
    <div class="hero-content">
      <h1>Enterprise-Grade Security</h1>
      <p class="hero-subtitle">Advanced authentication with 2FA, phone verification, and comprehensive audit logs</p>
      <button class="cta-btn" on:click={handleGetStarted}>
        Start Now
        <span>→</span>
      </button>
    </div>
    <div class="hero-visual">
      <div class="floating-card card-1">
        <span class="icon">🔐</span>
        <p>2FA Security</p>
      </div>
      <div class="floating-card card-2">
        <span class="icon">📱</span>
        <p>Phone OTP</p>
      </div>
      <div class="floating-card card-3">
        <span class="icon">📊</span>
        <p>Audit Logs</p>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section class="features">
    <h2>Powerful Features</h2>
    <div class="features-grid">
      <div class="feature-card">
        <div class="feature-icon">🔐</div>
        <h3>2-Factor Authentication</h3>
        <p>Email-based OTP verification adds an extra layer of security to your account</p>
      </div>

      <div class="feature-card">
        <div class="feature-icon">📱</div>
        <h3>Phone Verification</h3>
        <p>Verify your phone number with time-limited OTP codes for enhanced security</p>
      </div>

      <div class="feature-card">
        <div class="feature-icon">📊</div>
        <h3>Audit Logs</h3>
        <p>Comprehensive logging of all user activities and login attempts</p>
      </div>

      <div class="feature-card">
        <div class="feature-icon">🛡️</div>
        <h3>Account Protection</h3>
        <p>Automatic account lockout after failed login attempts to prevent unauthorized access</p>
      </div>

      <div class="feature-card">
        <div class="feature-icon">👥</div>
        <h3>Role-Based Access</h3>
        <p>Fine-grained access control with admin and user roles</p>
      </div>

      <div class="feature-card">
        <div class="feature-icon">📈</div>
        <h3>Real-time Monitoring</h3>
        <p>Monitor login trends and user activity with Prometheus & Grafana dashboards</p>
      </div>
    </div>
  </section>

  <!-- How It Works Section -->
  <section class="how-it-works">
    <h2>How It Works</h2>
    <div class="steps-container">
      <div class="step">
        <div class="step-number">1</div>
        <h3>Register</h3>
        <p>Create your account with email and password</p>
      </div>
      <div class="step-arrow">→</div>
      <div class="step">
        <div class="step-number">2</div>
        <h3>Verify Email</h3>
        <p>Enter the code sent to your email</p>
      </div>
      <div class="step-arrow">→</div>
      <div class="step">
        <div class="step-number">3</div>
        <h3>Verify Phone</h3>
        <p>Complete phone verification with OTP</p>
      </div>
      <div class="step-arrow">→</div>
      <div class="step">
        <div class="step-number">4</div>
        <h3>Login Securely</h3>
        <p>Use 2FA to access your dashboard</p>
      </div>
    </div>
  </section>

  <!-- Security Highlights -->
  <section class="security">
    <div class="security-content">
      <h2>Security First Approach</h2>
      <ul class="security-list">
        <li>✓ Industry-standard bcrypt password hashing</li>
        <li>✓ Time-limited OTP codes (10 minutes)</li>
        <li>✓ Account lockout after 5 failed attempts (30 minutes)</li>
        <li>✓ Comprehensive audit trail of all activities</li>
        <li>✓ Device and browser tracking</li>
        <li>✓ IP address logging for security monitoring</li>
      </ul>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="cta-section">
    <h2>Ready to Secure Your Account?</h2>
    <p>Join thousands of users who trust our secure authentication system</p>
    <button class="cta-btn-large" on:click={handleGetStarted}>Get Started Now</button>
  </section>

  <!-- Footer -->
  <footer class="footer">
    <p>&copy; 2026 ZLC Secure Authentication System. All rights reserved.</p>
  </footer>
</div>

<style>
  .home-container {
    background: linear-gradient(135deg, #1e3c72 0%, #2a5298 50%, #7e22ce 100%);
    min-height: 100vh;
    color: white;
    overflow-x: hidden;
  }

  /* Navbar */
  .navbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 40px;
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 24px;
    font-weight: 700;
  }

  .logo-icon {
    font-size: 32px;
  }

  .get-started-btn {
    background: white;
    color: #1e3c72;
    border: none;
    padding: 12px 24px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
  }

  .get-started-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  }

  /* Hero Section */
  .hero {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 80px 40px;
    max-width: 1200px;
    margin: 0 auto;
    position: relative;
  }

  .hero-content {
    flex: 1;
    z-index: 10;
  }

  .hero h1 {
    font-size: 64px;
    margin: 0 0 20px 0;
    line-height: 1.2;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  }

  .hero-subtitle {
    font-size: 20px;
    opacity: 0.9;
    margin-bottom: 30px;
    line-height: 1.6;
  }

  .cta-btn {
    background: white;
    color: #1e3c72;
    border: none;
    padding: 16px 40px;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 8px;
    transition: all 0.3s;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  }

  .cta-btn:hover {
    transform: translateY(-3px);
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.3);
  }

  .cta-btn span {
    font-size: 20px;
  }

  /* Hero Visual */
  .hero-visual {
    flex: 1;
    position: relative;
    height: 400px;
  }

  .floating-card {
    position: absolute;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 12px;
    padding: 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    animation: float 3s ease-in-out infinite;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0px);
    }
    50% {
      transform: translateY(-20px);
    }
  }

  .card-1 {
    top: 0;
    right: 0;
    animation-delay: 0s;
  }

  .card-2 {
    top: 150px;
    right: 100px;
    animation-delay: 0.5s;
  }

  .card-3 {
    top: 280px;
    right: 20px;
    animation-delay: 1s;
  }

  .floating-card .icon {
    font-size: 32px;
  }

  .floating-card p {
    margin: 0;
    font-size: 14px;
    text-align: center;
  }

  /* Features Section */
  .features {
    padding: 80px 40px;
    background: rgba(0, 0, 0, 0.1);
  }

  .features h2 {
    font-size: 48px;
    text-align: center;
    margin-bottom: 60px;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  }

  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 30px;
    max-width: 1200px;
    margin: 0 auto;
  }

  .feature-card {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 12px;
    padding: 40px 30px;
    text-align: center;
    transition: all 0.3s;
  }

  .feature-card:hover {
    background: rgba(255, 255, 255, 0.15);
    transform: translateY(-10px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
  }

  .feature-icon {
    font-size: 48px;
    margin-bottom: 20px;
  }

  .feature-card h3 {
    font-size: 20px;
    margin: 0 0 12px 0;
  }

  .feature-card p {
    margin: 0;
    opacity: 0.9;
    line-height: 1.6;
  }

  /* How It Works */
  .how-it-works {
    padding: 80px 40px;
    background: rgba(0, 0, 0, 0.05);
  }

  .how-it-works h2 {
    font-size: 48px;
    text-align: center;
    margin-bottom: 60px;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  }

  .steps-container {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
    flex-wrap: wrap;
    max-width: 1000px;
    margin: 0 auto;
  }

  .step {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 12px;
    padding: 30px;
    text-align: center;
    min-width: 200px;
    transition: all 0.3s;
  }

  .step:hover {
    background: rgba(255, 255, 255, 0.15);
    transform: translateY(-5px);
  }

  .step-number {
    background: white;
    color: #1e3c72;
    width: 50px;
    height: 50px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    font-size: 24px;
    margin: 0 auto 15px;
  }

  .step h3 {
    font-size: 18px;
    margin: 0 0 10px 0;
  }

  .step p {
    margin: 0;
    opacity: 0.9;
    font-size: 14px;
  }

  .step-arrow {
    font-size: 24px;
    opacity: 0.7;
  }

  /* Security Section */
  .security {
    padding: 80px 40px;
    background: rgba(0, 0, 0, 0.1);
  }

  .security-content {
    max-width: 800px;
    margin: 0 auto;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 12px;
    padding: 50px;
  }

  .security-content h2 {
    font-size: 40px;
    margin: 0 0 30px 0;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  }

  .security-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .security-list li {
    font-size: 18px;
    margin-bottom: 15px;
    padding-left: 30px;
    position: relative;
    line-height: 1.6;
  }

  /* CTA Section */
  .cta-section {
    padding: 80px 40px;
    text-align: center;
    background: rgba(0, 0, 0, 0.05);
  }

  .cta-section h2 {
    font-size: 48px;
    margin-bottom: 20px;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  }

  .cta-section p {
    font-size: 20px;
    margin-bottom: 40px;
    opacity: 0.9;
  }

  .cta-btn-large {
    background: white;
    color: #1e3c72;
    border: none;
    padding: 18px 50px;
    border-radius: 8px;
    font-size: 18px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.2);
  }

  .cta-btn-large:hover {
    transform: translateY(-3px);
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.3);
  }

  /* Footer */
  .footer {
    padding: 40px;
    text-align: center;
    background: rgba(0, 0, 0, 0.2);
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    opacity: 0.8;
  }

  .footer p {
    margin: 0;
    font-size: 14px;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
      padding: 40px 20px;
    }

    .hero h1 {
      font-size: 40px;
    }

    .hero-visual {
      margin-top: 40px;
      height: 300px;
    }

    .navbar {
      flex-direction: column;
      gap: 20px;
      padding: 20px;
    }

    .steps-container {
      flex-direction: column;
    }

    .step-arrow {
      transform: rotate(90deg);
    }

    .features h2,
    .how-it-works h2,
    .cta-section h2 {
      font-size: 32px;
    }
  }
</style>
