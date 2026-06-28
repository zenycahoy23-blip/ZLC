import express from 'express';
import pkg from 'pg';
import jwt from 'jsonwebtoken';
import { register, Counter, Histogram } from 'prom-client';

const { Pool } = pkg;

const app = express();

// CORS middleware
app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization');
  if (req.method === 'OPTIONS') {
    res.sendStatus(200);
  } else {
    next();
  }
});

app.use(express.json());

const pool = new Pool({
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  host: process.env.DB_HOST,
  port: parseInt(process.env.DB_PORT),
  database: process.env.DB_NAME,
});

pool.on('error', (err) => {
  console.error('Pool error:', err.message);
});

const httpRequestsTotal = new Counter({
  name: 'http_requests_total',
  help: 'Total HTTP requests',
  labelNames: ['method', 'route', 'status_code'],
});

const loginAttemptsTotal = new Counter({
  name: 'login_attempts_total',
  help: 'Total login attempts',
  labelNames: ['result'],
});

const httpRequestDuration = new Histogram({
  name: 'http_request_duration_seconds',
  help: 'HTTP request latency',
  labelNames: ['method', 'route'],
  buckets: [0.1, 0.5, 1, 2, 5],
});

const activeUsersTotal = new Counter({
  name: 'active_users_total',
  help: 'Total active users',
});

app.use((req, res, next) => {
  const start = Date.now();
  res.on('finish', () => {
    const duration = (Date.now() - start) / 1000;
    httpRequestsTotal.inc({ method: req.method, route: req.route?.path || req.path, status_code: res.statusCode });
    httpRequestDuration.observe({ method: req.method, route: req.route?.path || req.path }, duration);
  });
  next();
});

app.post('/register', async (req, res) => {
  try {
    const { email, password } = req.body;
    if (!email || !password) {
      return res.status(400).json({ error: 'Email and password required' });
    }
    const hashedPassword = Buffer.from(password).toString('base64');
    
    const result = await pool.query(
      'INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id',
      [email, hashedPassword]
    );
    
    activeUsersTotal.inc();
    res.json({ message: 'User registered', id: result.rows[0].id });
  } catch (error) {
    console.error('Register error:', error.message);
    res.status(400).json({ error: error.message });
  }
});

app.post('/login', async (req, res) => {
  try {
    const { email, password } = req.body;
    if (!email || !password) {
      return res.status(400).json({ error: 'Email and password required' });
    }
    const hashedPassword = Buffer.from(password).toString('base64');
    
    const result = await pool.query(
      'SELECT id FROM users WHERE email = $1 AND password = $2',
      [email, hashedPassword]
    );
    
    if (result.rows.length > 0) {
      const token = jwt.sign({ userId: result.rows[0].id }, process.env.JWT_SECRET, { expiresIn: '24h' });
      loginAttemptsTotal.inc({ result: 'success' });
      res.json({ token });
    } else {
      loginAttemptsTotal.inc({ result: 'failed' });
      res.status(401).json({ error: 'Invalid credentials' });
    }
  } catch (error) {
    console.error('Login error:', error.message);
    loginAttemptsTotal.inc({ result: 'failed' });
    res.status(500).json({ error: error.message });
  }
});

app.get('/metrics', (req, res) => {
  res.set('Content-Type', register.contentType);
  res.end(register.metrics());
});

app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

const PORT = process.env.PORT || 8080;
app.listen(PORT, () => console.log(`Backend running on port ${PORT}`));
