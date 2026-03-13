# QuizNinja API Scripts

This directory contains helper scripts for deployment and operations.

## Scripts

### `setup-gcp-secrets.sh`

Automatically creates secrets in Google Cloud Secret Manager by reading from your `.env` file.

**Purpose:** Simplifies Step 3 of the containerization guide by automating secret creation.

**Prerequisites:**
```bash
# 1. Install Google Cloud SDK
brew install google-cloud-sdk  # macOS
# or visit: https://cloud.google.com/sdk/docs/install

# 2. Authenticate
gcloud auth login

# 3. Set your project
gcloud config set project YOUR_PROJECT_ID

# 4. Enable Secret Manager API
gcloud services enable secretmanager.googleapis.com
```

**Usage:**
```bash
# From the quizninja-api directory
cd /Users/vishalvaibhav/Code/quizninja/quizninja-api

# Run the script
./scripts/setup-gcp-secrets.sh
```

**What it does:**
1. Validates prerequisites (gcloud installed, authenticated, project selected)
2. Reads your `.env` file
3. Creates the following secrets in GCP Secret Manager:
   - `SUPABASE_URL`
   - `SUPABASE_ANON_KEY`
   - `SUPABASE_SERVICE_KEY`
   - `SUPABASE_DB_HOST`
   - `SUPABASE_DB_USER`
   - `SUPABASE_DB_PASSWORD`
   - `SUPABASE_DB_NAME`
4. If secrets already exist, it creates new versions instead
5. Provides next steps for IAM permissions and deployment

**Output:**
```
=== QuizNinja GCP Secret Manager Setup ===

Current GCP Project: quizninja-475703

Reading secrets from .env file...

Creating secrets in Secret Manager...

✓ Created SUPABASE_URL
✓ Created SUPABASE_ANON_KEY
✓ Created SUPABASE_SERVICE_KEY
✓ Created SUPABASE_DB_HOST
✓ Created SUPABASE_DB_USER
✓ Created SUPABASE_DB_PASSWORD
✓ Created SUPABASE_DB_NAME

=== Setup Complete ===

Next steps:
1. Grant Secret Manager access to Cloud Run service account
2. Deploy to Cloud Run using the secrets
```

**Verify Secrets:**
```bash
# List all secrets
gcloud secrets list

# View a specific secret (not the value, just metadata)
gcloud secrets describe SUPABASE_URL

# Access a secret value (for testing)
gcloud secrets versions access latest --secret="SUPABASE_URL"
```

**Security Notes:**
- This script never prints secret values to the console
- Secrets are piped directly to gcloud using `--data-file=-`
- Your `.env` file remains local and is never uploaded anywhere
- Only create secrets from a trusted machine with proper access controls

**Troubleshooting:**

**Error: "Permission denied"**
```bash
# Grant yourself Secret Manager Admin role
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="user:YOUR_EMAIL" \
  --role="roles/secretmanager.admin"
```

**Error: "Secret already exists"**
- The script automatically handles this by creating a new version
- To completely recreate a secret: `gcloud secrets delete SECRET_NAME`

**Error: "API not enabled"**
```bash
gcloud services enable secretmanager.googleapis.com
```

## Next Steps After Running Script

1. **Grant Cloud Run Access to Secrets:**
   ```bash
   PROJECT_NUMBER=$(gcloud projects describe YOUR_PROJECT_ID --format='value(projectNumber)')

   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
     --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
     --role="roles/secretmanager.secretAccessor"
   ```

2. **Deploy to Cloud Run:**
   ```bash
   cd /Users/vishalvaibhav/Code/quizninja/quizninja-api

   gcloud run deploy quizninja-api \
     --source . \
     --region us-central1 \
     --platform managed \
     --allow-unauthenticated \
     --memory 512Mi \
     --set-env-vars="GIN_MODE=release,USE_SUPABASE=true" \
     --set-secrets="SUPABASE_URL=SUPABASE_URL:latest,SUPABASE_ANON_KEY=SUPABASE_ANON_KEY:latest,SUPABASE_SERVICE_KEY=SUPABASE_SERVICE_KEY:latest,SUPABASE_DB_HOST=SUPABASE_DB_HOST:latest,SUPABASE_DB_USER=SUPABASE_DB_USER:latest,SUPABASE_DB_PASSWORD=SUPABASE_DB_PASSWORD:latest,SUPABASE_DB_NAME=SUPABASE_DB_NAME:latest"
   ```

## How Secrets Work in Cloud Run

**Important:** Your application does NOT directly fetch secrets from Secret Manager.

Instead:
1. You create secrets in Secret Manager (using this script)
2. During deployment, you use `--set-secrets=` flag
3. Cloud Run automatically:
   - Fetches the secrets from Secret Manager
   - Mounts them as environment variables in your container
   - Makes them available via `os.Getenv()` in your Go code

Your application (`config/config.go`) simply reads from environment variables - it has no awareness that the values come from Secret Manager. This is a Google Cloud platform feature, not something your code needs to implement.

## Related Documentation

- [Google Cloud Secret Manager](https://cloud.google.com/secret-manager/docs)
- [Cloud Run Secrets](https://cloud.google.com/run/docs/configuring/secrets)
- [Containerization Guide](../prod-readiness/CONTAINERIZATION_GUIDE.md)
