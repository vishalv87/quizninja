#!/bin/bash

# Setup Google Cloud Secrets for QuizNinja API
# This script reads from your .env file and creates secrets in GCP Secret Manager
#
# Prerequisites:
# 1. gcloud CLI installed and authenticated
# 2. Google Cloud project selected (gcloud config set project YOUR_PROJECT_ID)
# 3. Secret Manager API enabled (gcloud services enable secretmanager.googleapis.com)
# 4. .env file exists in the parent directory

set -e  # Exit on error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../.env"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== QuizNinja GCP Secret Manager Setup ===${NC}"
echo ""

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${RED}Error: .env file not found at $ENV_FILE${NC}"
    exit 1
fi

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}Error: gcloud CLI is not installed${NC}"
    echo "Install it from: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Check if user is authenticated
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" &> /dev/null; then
    echo -e "${RED}Error: Not authenticated with gcloud${NC}"
    echo "Run: gcloud auth login"
    exit 1
fi

# Get current project
PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
if [ -z "$PROJECT_ID" ]; then
    echo -e "${RED}Error: No GCP project selected${NC}"
    echo "Run: gcloud config set project YOUR_PROJECT_ID"
    exit 1
fi

echo -e "${YELLOW}Current GCP Project:${NC} $PROJECT_ID"
echo ""

# Function to extract value from .env file
get_env_value() {
    local key=$1
    # Extract value, removing comments and trimming whitespace
    grep "^${key}=" "$ENV_FILE" | cut -d '=' -f2- | sed 's/#.*//' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//'
}

# Function to create or update a secret
create_secret() {
    local secret_name=$1
    local secret_value=$2

    if [ -z "$secret_value" ]; then
        echo -e "${YELLOW}Warning: Skipping $secret_name (empty value)${NC}"
        return
    fi

    # Check if secret already exists
    if gcloud secrets describe "$secret_name" --project="$PROJECT_ID" &> /dev/null; then
        echo -e "${YELLOW}Secret $secret_name already exists. Creating new version...${NC}"
        echo -n "$secret_value" | gcloud secrets versions add "$secret_name" \
            --project="$PROJECT_ID" \
            --data-file=- &> /dev/null
        echo -e "${GREEN}✓ Updated $secret_name${NC}"
    else
        echo -e "${YELLOW}Creating new secret: $secret_name${NC}"
        echo -n "$secret_value" | gcloud secrets create "$secret_name" \
            --project="$PROJECT_ID" \
            --replication-policy="automatic" \
            --data-file=- &> /dev/null
        echo -e "${GREEN}✓ Created $secret_name${NC}"
    fi
}

echo "Reading secrets from .env file..."
echo ""

# Extract secrets from .env
SUPABASE_URL=$(get_env_value "SUPABASE_URL")
SUPABASE_ANON_KEY=$(get_env_value "SUPABASE_ANON_KEY")
SUPABASE_SERVICE_KEY=$(get_env_value "SUPABASE_SERVICE_KEY")
SUPABASE_DB_HOST=$(get_env_value "SUPABASE_DB_HOST")
SUPABASE_DB_USER=$(get_env_value "SUPABASE_DB_USER")
SUPABASE_DB_PASSWORD=$(get_env_value "SUPABASE_DB_PASSWORD")
SUPABASE_DB_NAME=$(get_env_value "SUPABASE_DB_NAME")

echo "Creating secrets in Secret Manager..."
echo ""

# Create all secrets
create_secret "SUPABASE_URL" "$SUPABASE_URL"
create_secret "SUPABASE_ANON_KEY" "$SUPABASE_ANON_KEY"
create_secret "SUPABASE_SERVICE_KEY" "$SUPABASE_SERVICE_KEY"
create_secret "SUPABASE_DB_HOST" "$SUPABASE_DB_HOST"
create_secret "SUPABASE_DB_USER" "$SUPABASE_DB_USER"
create_secret "SUPABASE_DB_PASSWORD" "$SUPABASE_DB_PASSWORD"
create_secret "SUPABASE_DB_NAME" "$SUPABASE_DB_NAME"

echo ""
echo -e "${GREEN}=== Setup Complete ===${NC}"
echo ""
echo "Next steps:"
echo "1. Grant Secret Manager access to Cloud Run service account:"
echo "   PROJECT_NUMBER=\$(gcloud projects describe $PROJECT_ID --format='value(projectNumber)')"
echo "   gcloud projects add-iam-policy-binding $PROJECT_ID \\"
echo "     --member=\"serviceAccount:\${PROJECT_NUMBER}-compute@developer.gserviceaccount.com\" \\"
echo "     --role=\"roles/secretmanager.secretAccessor\""
echo ""
echo "2. Deploy to Cloud Run using the secrets:"
echo "   cd $SCRIPT_DIR/.."
echo "   gcloud run deploy quizninja-api \\"
echo "     --source . \\"
echo "     --region us-central1 \\"
echo "     --platform managed \\"
echo "     --allow-unauthenticated \\"
echo "     --memory 512Mi \\"
echo "     --set-env-vars=\"GIN_MODE=release,USE_SUPABASE=true\" \\"
echo "     --set-secrets=\"SUPABASE_URL=SUPABASE_URL:latest,SUPABASE_ANON_KEY=SUPABASE_ANON_KEY:latest,SUPABASE_SERVICE_KEY=SUPABASE_SERVICE_KEY:latest,SUPABASE_DB_HOST=SUPABASE_DB_HOST:latest,SUPABASE_DB_USER=SUPABASE_DB_USER:latest,SUPABASE_DB_PASSWORD=SUPABASE_DB_PASSWORD:latest,SUPABASE_DB_NAME=SUPABASE_DB_NAME:latest\""
echo ""
echo "To verify secrets were created:"
echo "   gcloud secrets list --project=$PROJECT_ID"
