package main

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestCreateListRemove(t *testing.T) {
	const (
		secretName  = "TEST_SECRET_123"
		secretValue = "test-value-123"
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client := NewClient(os.Getenv("OWNER"), os.Getenv("REPO"), os.Getenv("TOKEN"), false)

	names, err := client.ApplySecrets(ctx, map[string]string{secretName: secretValue}, true)
	if err != nil {
		t.Fatalf("Error applying secrets: %v", err)
	}

	if !reflect.DeepEqual(names, []string{secretName}) {
		t.Fatalf("Expected 1 secret to be applied, got %d", len(names))
	}

	secrets, err := client.ListSecrets(ctx, []string{secretName})
	if err != nil {
		t.Fatalf("Error listing secrets: %v", err)
	}

	if len(secrets) != 1 {
		t.Fatalf("Expected 1 secret to be listed, got %d", len(secrets))
	}

	if secrets[0].Name != secretName {
		t.Fatalf("Expected secret name: %s, got: %s", secretName, secrets[0].Name)
	}

	deleted, err := client.DeleteSecrets(ctx, []string{secretName})
	if err != nil {
		t.Fatalf("Error deleting secrets: %v", err)
	}

	if !reflect.DeepEqual(deleted, []string{secretName}) {
		t.Fatalf("Expected 1 secret to be deleted, got %d", len(deleted))
	}
}
