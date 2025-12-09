package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Secret represents a single secret entry from the CSV.
type Secret struct {
	Name  string
	Value string
}

func main() {
	projectID := flag.String("project-id", "", "your Google cloud project")
	secretsFile := flag.String("secrets-file", "", "path to the CSV file containing secrets (name,value)")
	location := flag.String("secrets-location", "", "location of the secrets created")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if *projectID == "" || *secretsFile == "" || *location == "" {
		logger.Error("project-id and secrets-file are required and so the locaton")
		flag.Usage()
		os.Exit(1)
	}

	secrets, err := readSecretsCSV(*secretsFile)
	if err != nil {
		logger.Error("error reading secrets file", "error", err)
		os.Exit(1)
	}

	logger.Info("found secrets to upload", "count", len(secrets))

	err = uploadSecrets(logger, *projectID, *location, secrets)
	if err != nil {
		logger.Error("error uploading secrets", "error", err)
		os.Exit(1)
	}

	logger.Info("successfully uploaded all secrets")
}

func readSecretsCSV(filePath string) ([]Secret, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %q: %w", filePath, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 2
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("unable to read CSV header: %w", err)
	}
	if strings.ToLower(header[0]) != "name" || strings.ToLower(header[1]) != "value" {
		return nil, fmt.Errorf("CSV header must be 'name,value', got '%s,%s'", header[0], header[1])
	}

	var secrets []Secret
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("unable to read CSV record: %w", err)
		}

		secrets = append(secrets, Secret{
			Name:  record[0],
			Value: record[1],
		})
	}

	return secrets, nil
}

func uploadSecrets(logger *slog.Logger, projectID string, location string, secrets []Secret) error {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	for _, secret := range secrets {
		parent := fmt.Sprintf("projects/%s", projectID)
		secretName := secret.Name
		secretValue := []byte(secret.Value)

		_, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
			Name: fmt.Sprintf("%s/secrets/%s", parent, secretName),
		})

		if err != nil {
			if strings.Contains(err.Error(), "NotFound") {

				logger.Info("creating secret", "secret_name", secretName)
				_, err = client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
					Parent:   parent,
					SecretId: secretName,
					Secret: &secretmanagerpb.Secret{
						Replication: &secretmanagerpb.Replication{
							Replication: &secretmanagerpb.Replication_UserManaged_{
								UserManaged: &secretmanagerpb.Replication_UserManaged{
									Replicas: []*secretmanagerpb.Replication_UserManaged_Replica{
										{
											Location: location,
										},
									},
								},
							},
						},
					},
				})
				if err != nil {
					return fmt.Errorf("failed to create secret %s: %w", secretName, err)
				}
			} else {
				return fmt.Errorf("failed to check existence of secret %s: %w", secretName, err)
			}
		}

		logger.Info("adding version to secret", "secret_name", secretName)
		_, err = client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
			Parent: fmt.Sprintf("%s/secrets/%s", parent, secretName),
			Payload: &secretmanagerpb.SecretPayload{
				Data: secretValue,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to add secret version to %s: %w", secretName, err)
		}
	}

	return nil
}
