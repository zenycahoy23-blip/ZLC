import React, { useState } from 'react';
import ReactDOM from 'react-dom/client';

function App() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [token, setToken] = useState(localStorage.getItem('token') || '');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const register = async () => {
    setLoading(true);
    setMessage('');
    try {
      const res = await fetch('http://localhost:8082/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const data = await res.json();
      if (res.ok) {
        setMessage('Registered successfully! Now log in.');
        setEmail('');
        setPassword('');
      } else {
        setMessage(`Error: ${data.error}`);
      }
    } catch (err) {
      setMessage(`Registration failed: ${err.message}`);
    }
    setLoading(false);
  };

  const login = async () => {
    setLoading(true);
    setMessage('');
    try {
      const res = await fetch('http://localhost:8082/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const data = await res.json();
      if (res.ok && data.token) {
        localStorage.setItem('token', data.token);
        setToken(data.token);
        setEmail('');
        setPassword('');
        setMessage('Login successful!');
      } else {
        setMessage(`Login failed: ${data.error || 'Unknown error'}`);
      }
    } catch (err) {
      setMessage(`Login failed: ${err.message}`);
    }
    setLoading(false);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken('');
    setMessage('Logged out');
  };

  if (token) {
    return (
      <div style={{ maxWidth: '600px', margin: '50px auto', fontFamily: 'Arial', textAlign: 'center' }}>
        <h1>Dashboard</h1>
        <p>You are logged in!</p>
        <div style={{ backgroundColor: '#f0f0f0', padding: '20px', borderRadius: '5px', marginBottom: '20px' }}>
          <p style={{ wordBreak: 'break-all', fontSize: '12px' }}>
            <strong>Token:</strong> {token}
          </p>
        </div>
        <button 
          onClick={logout} 
          style={{ padding: '10px 20px', backgroundColor: '#dc3545', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer' }}
        >
          Logout
        </button>
        <p style={{ marginTop: '20px', fontSize: '14px', color: '#666' }}>
          Your authentication token is securely stored. Use it to make authenticated API requests.
        </p>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '400px', margin: '50px auto', fontFamily: 'Arial' }}>
      <h1 style={{ textAlign: 'center' }}>Auth System</h1>
      {message && (
        <div style={{ 
          padding: '10px', 
          marginBottom: '15px', 
          backgroundColor: message.includes('Error') || message.includes('failed') ? '#f8d7da' : '#d4edda',
          color: message.includes('Error') || message.includes('failed') ? '#721c24' : '#155724',
          borderRadius: '5px',
          border: '1px solid',
          borderColor: message.includes('Error') || message.includes('failed') ? '#f5c6cb' : '#c3e6cb'
        }}>
          {message}
        </div>
      )}
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        style={{ width: '100%', padding: '10px', marginBottom: '10px', boxSizing: 'border-box', borderRadius: '5px', border: '1px solid #ccc' }}
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        style={{ width: '100%', padding: '10px', marginBottom: '15px', boxSizing: 'border-box', borderRadius: '5px', border: '1px solid #ccc' }}
      />
      <div style={{ display: 'flex', gap: '10px' }}>
        <button 
          onClick={register} 
          disabled={loading} 
          style={{ flex: 1, padding: '10px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer', opacity: loading ? 0.5 : 1 }}
        >
          Register
        </button>
        <button 
          onClick={login} 
          disabled={loading} 
          style={{ flex: 1, padding: '10px', backgroundColor: '#28a745', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer', opacity: loading ? 0.5 : 1 }}
        >
          Login
        </button>
      </div>
    </div>
  );
}

ReactDOM.createRoot(document.getElementById('root')).render(<App />);
