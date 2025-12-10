# Google Secrets Uploader

A simple command-line utility written in Go to process a spreadsheet and upload bulk secrets into Google Cloud Secret Manager.

## Usage

You must specify a replication policy for the secrets. You can choose between automatic (global) replication or user-managed replication in one or more regions.

**Flags:**

*   `--project-id`: (Required) Your Google Cloud project ID.
*   `--secrets-file`: (Required) Path to the CSV file containing the secrets.
*   `--global`: Use automatic replication. The secret will be replicated globally.
*   `--secrets-location`: Use user-managed replication. Provide a single region or a comma-separated list of regions.

**Note:** You must use either `--global` or `--secrets-location`, but not both.

### Examples

**1. Global Secret (Automatic Replication)**

```bash
go run main.go \
  --project-id <your-gcp-project-id> \
  --secrets-file <path-to-your-secrets.csv> \
  --global
```

**2. Single-Region Secret (User-Managed Replication)**

```bash
go run main.go \
  --project-id <your-gcp-project-id> \
  --secrets-file <path-to-your-secrets.csv> \
  --secrets-location us-central1
```

**3. Multi-Region Secret (User-Managed Replication)**

```bash
go run main.go \
  --project-id <your-gcp-project-id> \
  --secrets-file <path-to-your-secrets.csv> \
  --secrets-location us-central1,us-east1,us-west1
```

## Spreadsheet Format

The secrets file should be a CSV with two columns: `name` and `value`.

- `name`: The name of the secret in Google Cloud Secret Manager.
- `value`: The secret content to be uploaded.

**Example `secrets.csv`:**

```csv
name,value
my-secret-1,my-super-secret-value
my-secret-2,another-secret
```

## Security Note

:warning: **Do not check in CSV files with real passwords, secrets, or any other confidential information into version control.** The sample CSV files in this repository are for testing purposes only.

## Setup

1.  Ensure Go is installed (version 1.16 or higher recommended).
2.  Authenticate with Google Cloud. The application will use Application Default Credentials. Make sure your environment is configured (e.g., `gcloud auth application-default login`).

## Running with Docker

You can also run this utility within a Docker container.

### 1. Build the Docker Image

First, build the Docker image from the root of the project:

```bash
docker build -t google-secrets-uploader .
```

### 2. Run the Container

When running the container, you need to mount your secrets file and your Google Cloud credentials into the container.

**Example (Single-Region):**

```bash
docker run --rm \
  -v $(pwd)/secrets.csv:/app/secrets.csv \
  -v ~/.config/gcloud/application_default_credentials.json:/root/.config/gcloud/application_default_credentials.json \
  google-secrets-uploader \
  --project-id <your-gcp-project-id> \
  --secrets-file /app/secrets.csv \
  --secrets-location us-central1
```

**Example (Global):**

```bash
docker run --rm \
  -v $(pwd)/secrets.csv:/app/secrets.csv \
  -v ~/.config/gcloud/application_default_credentials.json:/root/.config/gcloud/application_default_credentials.json \
  google-secrets-uploader \
  --project-id <your-gcp-project-id> \
  --secrets-file /app/secrets.csv \
  --global
```

**Volume Mounts Explained:**

*   `-v $(pwd)/secrets.csv:/app/secrets.csv`: This mounts the `secrets.csv` file from your current directory on your host machine to the `/app/secrets.csv` path inside the container. The application reads the secrets from this file.
*   `-v ~/.config/gcloud/application_default_credentials.json:/root/.config/gcloud/application_default_credentials.json`: This mounts your local Google Cloud Application Default Credentials into the container. This is how the application authenticates with Google Cloud. The path `~/.config/gcloud/application_default_credentials.json` is the default location for these credentials.
