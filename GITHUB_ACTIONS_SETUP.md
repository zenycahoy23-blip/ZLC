# GitHub Actions Setup Guide

## 🚀 Quick Start

To enable GitHub Actions CI/CD for UnoReverse, follow these steps:

---

## Step 1: Add Slack Webhook Secret

**Purpose:** Enable Slack notifications for all pipeline events

1. Go to your GitHub repository
2. Navigate to: **Settings → Secrets and variables → Actions**
3. Click **New repository secret**
4. Fill in:
   - **Name:** `SLACK_WEBHOOK_URL`
   - **Value:** Your Slack webhook URL
5. Click **Add secret**

### Getting your Slack Webhook URL:

1. Go to [Slack API Apps](https://api.slack.com/apps)
2. Select your app (or create a new one)
3. Go to **Incoming Webhooks**
4. Toggle **Activate Incoming Webhooks** to ON
5. Click **Add New Webhook to Workspace**
6. Select a channel (e.g., #deployments)
7. Copy the generated Webhook URL
8. Use it in the GitHub secret

---

## Step 2: Enable GitHub Container Registry (GHCR)

**Purpose:** Push Docker images to GitHub's container registry

### Option A: Using GITHUB_TOKEN (Recommended)

1. Go to **Settings → Actions → General**
2. Under **Workflow permissions**, select:
   - ✓ **Read and write permissions**
   - ✓ **Allow GitHub Actions to create and approve pull requests**
3. Click **Save**

The pipeline will automatically use `GITHUB_TOKEN` to push images.

### Option B: Using Personal Access Token (PAT)

1. Go to **Settings → Developer settings → Personal access tokens → Tokens (classic)**
2. Click **Generate new token (classic)**
3. Select scopes:
   - ✓ `repo` (full control)
   - ✓ `write:packages`
   - ✓ `read:packages`
4. Copy the token
5. Go to repository **Settings → Secrets and variables → Actions**
6. Add secret:
   - **Name:** `GHCR_PAT`
   - **Value:** `<your-token>`

---

## Step 3: Configure Repository Actions Permissions

1. Go to **Settings → Actions → General**
2. Under **Actions permissions**, select:
   - ✓ **Allow all actions and reusable workflows**
3. Click **Save**

---

## Step 4: Test the Pipeline

### Method 1: Push to Develop Branch

```bash
git add .
git commit -m "Test GitHub Actions"
git push origin develop
```

The pipeline will:
- Run full CI/CD on develop
- Deploy to staging
- Send Slack notification

### Method 2: Push to Main Branch

```bash
git add .
git commit -m "Release version"
git push origin main
```

The pipeline will:
- Run full CI/CD on main
- Build and push images to GHCR
- Deploy to production
- Send Slack notification

---

## Step 5: Monitor Pipeline

1. Go to **Actions** tab in your repository
2. Click the latest workflow run
3. View real-time logs for each stage
4. Check Slack for notifications

---

## 📊 Pipeline Overview

### Stages:

| Stage | Trigger | Purpose |
|-------|---------|---------|
| **build-and-test** | All | Build backend (Go) & frontend (Node.js), run unit tests |
| **lint-backend** | All | Run golangci-lint on Go code |
| **lint-frontend** | All | Run ESLint on JavaScript code |
| **security-scan** | All | Trivy vulnerability scan + npm audit |
| **database-validation** | All | Validate PostgreSQL schema |
| **performance-test** | All | Baseline performance metrics |
| **build-and-push-images** | main/develop | Build & push to GHCR |
| **notify-slack** | All | Send detailed Slack notification |
| **deploy-staging** | develop | Deploy to staging environment |
| **deploy-production** | main | Deploy to production environment |

### Branch Behavior:

- **develop**: Full pipeline + staging deployment
- **main**: Full pipeline + production deployment
- **pull_request**: CI checks only (no deployment)

---

## 🔑 Secrets Required

| Secret | Description | Required |
|--------|-------------|----------|
| `SLACK_WEBHOOK_URL` | Slack incoming webhook | ✓ Yes |
| `GITHUB_TOKEN` | Auto-provided by GitHub Actions | ✓ Auto |

---

## 🚨 Troubleshooting

### Workflow fails with "Resource not accessible by integration"

**Solution:**
1. Go to **Settings → Actions → General**
2. Under **Workflow permissions**, enable **Read and write permissions**
3. Save and re-run workflow

### "SLACK_WEBHOOK_URL secret not found"

**Solution:**
1. Verify the secret is added in **Settings → Secrets and variables → Actions**
2. Check the secret name exactly matches: `SLACK_WEBHOOK_URL`
3. Re-run workflow

### Slack notifications not sending

**Solution:**
1. Verify webhook URL is correct
2. Test webhook manually:
   ```bash
   curl -X POST -H 'Content-type: application/json' \
     --data '{"text":"test"}' \
     YOUR_WEBHOOK_URL
   ```
3. Check Slack workspace app permissions
4. Verify channel exists and bot has access

### Images not pushing to GHCR

**Solution:**
1. Ensure `GITHUB_TOKEN` has write permissions
2. Check **Settings → Actions → General → Workflow permissions**
3. Verify repository is public (or PAT has write:packages scope)
4. Check workflow logs for specific error

---

## 📝 Customization

### Change notification channel:

In `.github/workflows/ci-cd.yml`, update Slack payload channel:

```yaml
"channel": "#your-channel-name"
```

### Add deployment environment:

Edit `deploy-production` job to add `environment`:

```yaml
environment:
  name: production
  url: https://your-domain.com
```

### Disable specific stages:

Comment out jobs in `.github/workflows/ci-cd.yml` YAML (not recommended)

---

## ✅ Verification Checklist

- [ ] `SLACK_WEBHOOK_URL` secret added
- [ ] Repository Actions permissions enabled
- [ ] `.github/workflows/ci-cd.yml` is committed
- [ ] Slack workspace configured with incoming webhook
- [ ] Test push to develop branch
- [ ] Verify Slack notification received
- [ ] Logs visible in Actions tab

---

## 📚 Useful Resources

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Slack API Apps](https://api.slack.com/apps)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [GitHub Secrets Guide](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions)
- [Workflow Syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

---

## 🎯 Next Steps

1. **Complete setup checklist** ✓
2. **Add `SLACK_WEBHOOK_URL` secret** ✓
3. **Push code to develop** ✓
4. **Monitor first workflow run** ✓
5. **Verify Slack notification** ✓
6. **Set up production deployment** (optional)

