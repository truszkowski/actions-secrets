package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	sodium "github.com/GoKillers/libsodium-go/cryptobox"
	"github.com/google/go-github/v59/github"
)

// Assumes a function that encrypts the secret value correctly
func encryptSecretValue(publicKey, secretValue string) (string, error) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %w", err)
	}

	encryptedBytes, exit := sodium.CryptoBoxSeal([]byte(secretValue), pubKeyBytes)
	if exit != 0 {
		log.Fatalf("Failed to encrypt secret with libsodium, exit code: %d", exit)
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

type Client struct {
	client  *github.Client
	owner   string
	repo    string
	token   string
	verbose bool
}

func NewClient(owner, repo, token string, verbose bool) *Client {
	return &Client{
		client:  github.NewClient(nil).WithAuthToken(token),
		owner:   owner,
		repo:    repo,
		token:   token,
		verbose: verbose,
	}
}

func (cli *Client) ApplySecrets(ctx context.Context, secrets map[string]string, override bool) ([]string, error) {
	publicKey, _, err := cli.client.Actions.GetRepoPublicKey(ctx, cli.owner, cli.repo)
	if err != nil {
		return nil, fmt.Errorf("error fetching public key: %w", err)
	}

	var applied []string

	for name, value := range secrets {
		encryptedValue, err := encryptSecretValue(publicKey.GetKey(), value)
		if err != nil {
			return applied, fmt.Errorf("error encrypting secret '%s': %w", name, err)
		}

		if !override {
			secret, err := cli.GetSecret(ctx, name)
			if err != nil {
				return applied, fmt.Errorf("error fetching secret '%s': %w", name, err)
			}

			if secret != nil {
				applied = append(applied, name)
				if cli.verbose {
					fmt.Printf("Secret '%s' already exists, skipping\n", name)
				}
				continue
			}
		}

		secret := &github.EncryptedSecret{
			Name:           name,
			KeyID:          publicKey.GetKeyID(),
			EncryptedValue: encryptedValue,
		}

		res, err := cli.client.Actions.CreateOrUpdateRepoSecret(ctx, cli.owner, cli.repo, secret)
		if err != nil {
			if res != nil && res.Body != nil {
				msg, _ := io.ReadAll(res.Body)
				return applied, fmt.Errorf("error applying '%s': %w\nmessage: %s", name, err, string(msg))
			}
			return applied, fmt.Errorf("error applying '%s': %w", name, err)
		}

		applied = append(applied, name)
		if cli.verbose {
			fmt.Printf("Secret '%s' applied successfully\n", name)
		}
	}

	return applied, nil
}

func (cli *Client) GetSecret(ctx context.Context, name string) (*github.Secret, error) {
	secret, res, err := cli.client.Actions.GetRepoSecret(ctx, cli.owner, cli.repo, name)
	if err != nil {
		if res.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching secret '%s': %w", name, err)
	}

	return secret, nil
}

func (cli *Client) ListSecrets(ctx context.Context, names []string) ([]*github.Secret, error) {
	opts := &github.ListOptions{PerPage: 20, Page: 1}
	var secrets []*github.Secret

	for {
		page, _, err := cli.client.Actions.ListRepoSecrets(ctx, cli.owner, cli.repo, opts)
		if err != nil {
			return nil, fmt.Errorf("error fetching secrets: %w", err)
		}

		for _, secret := range page.Secrets {
			for _, name := range names {
				if secret.Name == name {
					secrets = append(secrets, secret)
				}
			}
		}

		if opts.PerPage > len(page.Secrets) {
			break
		}

		opts.Page++
	}

	return secrets, nil
}

func (cli *Client) ListAllSecrets(ctx context.Context) ([]*github.Secret, error) {
	opts := &github.ListOptions{PerPage: 20, Page: 1}
	var secrets []*github.Secret

	for {
		page, _, err := cli.client.Actions.ListRepoSecrets(ctx, cli.owner, cli.repo, opts)
		if err != nil {
			return nil, fmt.Errorf("error fetching secrets: %w", err)
		}

		secrets = append(secrets, page.Secrets...)

		if opts.PerPage > len(page.Secrets) {
			break
		}

		opts.Page++
	}

	return secrets, nil
}

func (cli *Client) DeleteSecrets(ctx context.Context, names []string) ([]string, error) {
	var deleted []string
	for _, name := range names {
		res, err := cli.client.Actions.DeleteRepoSecret(ctx, cli.owner, cli.repo, name)
		if err != nil {
			if res != nil && res.StatusCode == 404 {
				if cli.verbose {
					fmt.Printf("Secret '%s' does not exist, skipping\n", name)
				}
				continue
			}
			return deleted, fmt.Errorf("error deleting secret '%s': %w", name, err)
		}
		deleted = append(deleted, name)
		if cli.verbose {
			fmt.Printf("Secret '%s' deleted successfully\n", name)
		}
	}
	return deleted, nil
}
