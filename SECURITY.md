# Security & Secrets Management

**Important**: This document ensures we never accidentally commit sensitive data to the repository.

---

## 🔐 What We Protect

### Never Commit
- ❌ Private keys (.pem, .key files)
- ❌ API keys and tokens
- ❌ Database credentials
- ❌ AWS/cloud credentials
- ❌ Passwords and secrets
- ❌ Wrapped encryption keys (in production)
- ❌ Configuration with secrets

### Safe to Commit
- ✅ Source code
- ✅ Tests and documentation
- ✅ Build configuration
- ✅ Public keys (.pub files)
- ✅ Example configuration (.example files)
- ✅ Docker setup files (without secrets)

---

## 📋 Security Checklist

Before committing code:

- [ ] No `.env` files committed
- [ ] No hardcoded API keys
- [ ] No passwords in code
- [ ] No private keys in repo
- [ ] No credentials.json files
- [ ] No database backups with data
- [ ] No token files
- [ ] No AWS credentials
- [ ] No secrets in config files

---

## 🛡️ .gitignore Protection

Our `.gitignore` automatically prevents committing:

### Environment Files
```
.env
.env.local
.env.*.local
.env.*.yml
*.env
config.local.yml
```

### Private Keys
```
*.pem
*.key
*.pub
*.gpg
id_rsa*
id_ed25519*
```

### Credentials
```
credentials.json
secrets.json
.aws/credentials
.aws/config
.docker/config.json
```

### API Keys
```
.api-key
.api-keys
api-keys.txt
tokens.txt
*_tokens.json
*_credentials.json
```

### Database Files
```
*.db
*.sqlite
*.sqlite3
*.sql.bak
*.backup
```

---

## ⚙️ Environment Setup

### 1. Create .env File (Never Commit)

Create `.env` in project root:

```bash
# .env (NOT COMMITTED - use this for local development)

# PostgreSQL (development)
DATABASE_URL="postgres://auta:auta_dev_password@localhost/auta?sslmode=disable"

# Service
PORT=8000
LOG_LEVEL=debug

# API Keys (if needed in future)
API_KEY_SECRET="your-secret-here"
JWT_SECRET="your-jwt-secret"

# AWS/Cloud (if needed in future)
AWS_ACCESS_KEY_ID="xxx"
AWS_SECRET_ACCESS_KEY="xxx"
AWS_REGION="us-east-1"
```

**This file is in `.gitignore` - it won't be committed.**

### 2. Create .env.example (Commit This)

Create `.env.example` with no actual secrets:

```bash
# .env.example (SAFE TO COMMIT - shows structure only)

# PostgreSQL
DATABASE_URL="postgres://user:password@host/database?sslmode=disable"

# Service
PORT=8000
LOG_LEVEL=info

# API Keys (placeholder)
API_KEY_SECRET="your-api-key-here"
JWT_SECRET="your-jwt-secret"

# AWS/Cloud (placeholder)
AWS_ACCESS_KEY_ID="your-aws-key-id"
AWS_SECRET_ACCESS_KEY="your-aws-secret"
AWS_REGION="us-east-1"
```

**This file IS committed - shows the structure developers need.**

---

## 🚀 Development Workflow

### For Local Development

1. **Copy example file**:
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your values**:
   ```bash
   # Edit local secrets only
   vim .env
   ```

3. **`.env` is ignored by git**:
   ```bash
   # Verify it's in .gitignore
   git status .env  # Should say: ignored
   ```

4. **Run service with environment**:
   ```bash
   export $(cat .env | xargs)
   make run-metadata
   ```

---

## 🔍 Verify Nothing Sensitive is Tracked

### Check Current Repo

```bash
# List all files in git
git ls-files

# Should NOT show:
# .env
# .env.local
# *.key
# credentials.json
# Any secrets or private keys
```

### Check for Leaked Secrets

```bash
# Search for common patterns
git grep -i "password\|secret\|key\|token" -- ':!*.md' ':!*.txt' ':!tests'

# If this returns results (excluding docs), review them!
```

### Scan Before Committing

```bash
# Before committing, verify no secrets:
git diff --cached | grep -i "password\|api_key\|secret"

# If this returns anything, FIX IT before committing!
```

---

## ⚠️ If You Accidentally Commit a Secret

### Immediate Action

```bash
# 1. Remove from working directory
rm the-secret-file

# 2. Unstage if not yet committed
git reset HEAD the-secret-file

# 3. If already committed, use BFG or git-filter-repo:
# Install: brew install bfg
bfg --delete-files the-secret-file

# 4. Or remove from history:
git filter-branch --tree-filter 'rm -f the-secret-file' HEAD
```

### Then

1. ⚠️ **REVOKE the secret immediately** (change passwords, rotate keys, etc.)
2. Commit the fix
3. Force push (if needed): `git push --force` (use carefully!)
4. Alert the team

---

## 🔑 Handling Different Types of Secrets

### API Keys

```bash
# ❌ WRONG - Don't do this
const API_KEY = "sk-1234567890abcdef"

# ✅ RIGHT - Use environment variables
const API_KEY = process.env.API_KEY_SECRET
```

### Database Credentials

```bash
# ❌ WRONG - Hardcoded
DATABASE_URL="postgres://admin:password123@db.example.com/prod"

# ✅ RIGHT - Environment variable
DATABASE_URL=os.Getenv("DATABASE_URL")
```

### Private Keys

```bash
# ❌ WRONG - Committed to repo
id_rsa (in repo)

# ✅ RIGHT - Never commit, only reference
KEY_FILE=/home/user/.ssh/id_rsa (in .gitignore)
```

### Configuration with Secrets

```bash
# ❌ WRONG - Config with real secrets
{
  "database_password": "prod_password_123",
  "api_key": "sk-12345"
}

# ✅ RIGHT - Config references env vars
{
  "database_password": "${DB_PASSWORD}",
  "api_key": "${API_KEY_SECRET}"
}
```

---

## 📚 For Different Deployment Stages

### Development (.env - local only)
```bash
.env                           # Local secrets (ignored)
.env.example                   # Template (committed)
```

### Testing (.env.test)
```bash
.env.test                      # Test secrets (ignored)
.env.test.example              # Test template (committed)
```

### Production
```bash
# Never store secrets in repo!
# Use deployment platform's secret management:
# - Kubernetes Secrets
# - AWS Secrets Manager
# - HashiCorp Vault
# - GitHub Secrets (for CI/CD)
```

---

## 🔐 CI/CD Pipeline

### GitHub Actions Example

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - name: Setup test database
        env:
          DATABASE_URL: ${{ secrets.TEST_DATABASE_URL }}
        run: |
          make db-test-up
          make test
```

**Key**: Secrets are passed via GitHub Secrets, never committed.

---

## 📖 Best Practices

1. **Never commit secrets** - Always use environment variables
2. **Use `.example` files** - Show structure without secrets
3. **Keep `.gitignore` updated** - As you add new secret types
4. **Review before committing** - Check diffs for sensitive data
5. **Rotate secrets regularly** - Even if not leaked
6. **Use secure vaults** - For production secrets
7. **Audit logs** - Track who accessed what
8. **Document locations** - Where secrets are stored
9. **Team training** - Everyone knows the rules
10. **Automated checks** - Use pre-commit hooks

---

## 🛠️ Pre-commit Hook (Optional)

Create `.git/hooks/pre-commit` to prevent accidental commits:

```bash
#!/bin/bash

# Prevent commits with secrets

# Check for common patterns
git diff --cached | grep -E 'password|api_key|secret|token|credentials' && {
    echo "ERROR: Potential secrets found in staged changes!"
    echo "Please review the following lines:"
    git diff --cached | grep -E 'password|api_key|secret|token|credentials'
    exit 1
}

# Check for files that shouldn't be committed
if git diff --cached --name-only | grep -E '\.env($|\.)|\.key$|credentials\.json|secrets\.'; then
    echo "ERROR: Sensitive files detected in staging area!"
    exit 1
fi

exit 0
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

## ✅ Current Status

Our project currently has:

✅ Comprehensive `.gitignore` with all secret types  
✅ Environment variable support in Makefile  
✅ Example configuration pattern  
✅ Documentation of best practices  
✅ No secrets currently in repository  

### What to Do Now

1. ✅ Review `.gitignore` (already updated)
2. ✅ Review this security guide
3. → Create `.env.example` when needed
4. → Add secrets to `.env` (for local development)
5. → Team follows best practices

---

## 🔍 Verify Your Setup

```bash
# Check nothing sensitive is tracked
git ls-files | grep -E '\.env|\.key|credentials|secret' 
# Should return nothing

# Check .gitignore is working
echo "test-secret" > .env.test
git status
# Should show ".env.test" as "ignored"

# Clean up
rm .env.test
```

---

## 📞 Questions?

- **How do I handle database passwords?** → Use `.env` file + `.env.example`
- **What about production secrets?** → Use deployment platform's secret manager
- **Do I commit `.env.example`?** → Yes, with placeholder values only
- **What if I committed a secret?** → See "If You Accidentally Commit a Secret" section
- **How do I rotate secrets?** → Change in secret manager, deploy
- **How do I share secrets with team?** → Use secure channel (1Password, Vault, etc.), never email

---

## 🎯 Summary

✅ Never commit sensitive data  
✅ Always use environment variables  
✅ Use `.env.example` to show structure  
✅ Keep `.gitignore` comprehensive  
✅ Review diffs before committing  
✅ Use platform-specific secret management  

**Rule of thumb**: If you wouldn't write it on a whiteboard in a public cafe, don't commit it to the repo.
