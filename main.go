package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("TOKEN")

	var owner, repo, applyPath, deletePath, listPath string
	flag.StringVar(&owner, "owner", "", "owner of the repository")
	flag.StringVar(&repo, "repo", "", "name of the repository")
	flag.StringVar(&applyPath, "apply", "", "path to the file containing secrets to apply")
	flag.StringVar(&deletePath, "delete", "", "path to the file containing secrets to delete")
	flag.StringVar(&listPath, "list", "", "path to the file containing secrets to list")

	var verbose, override, listAll bool
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.BoolVar(&override, "override", false, "override existing secrets")
	flag.BoolVar(&listAll, "list-all", false, "list all secrets")
	flag.Parse()

	client := NewClient(owner, repo, token, verbose)

	if deletePath != "" {
		names, err := loadEnvVarsNames(deletePath)
		if err != nil {
			log.Fatalf("Error loading secret names from file: %v", err)
		}

		deleted, err := client.DeleteSecrets(ctx, names)
		if err != nil {
			log.Fatalln("Error deleting secrets:", err)
		}
		_ = deleted
	}

	if applyPath != "" {
		secrets, err := loadEnvVars(applyPath)
		if err != nil {
			log.Fatalf("Error loading secrets from file: %v", err)
		}

		applied, err := client.ApplySecrets(ctx, secrets, override)
		if err != nil {
			log.Fatalln("Error applying secrets:", err)
		}
		_ = applied
	}

	if listPath != "" {
		names, err := loadEnvVarsNames(listPath)
		if err != nil {
			log.Fatalf("Error loading secret names from file: %v", err)
		}

		secrets, err := client.ListSecrets(ctx, names)
		if err != nil {
			log.Fatalln("Error listing secrets:", err)
		}

		for _, secret := range secrets {
			if verbose {
				fmt.Println(secret.Name, "created:", secret.CreatedAt.Format("2006-01-02T15:04:05"), "updated:", secret.UpdatedAt.Format("2006-01-02T15:04:05"))
			} else {
				fmt.Println(secret.Name)
			}
		}
	}

	if listAll {
		secrets, err := client.ListAllSecrets(ctx)
		if err != nil {
			log.Fatalln("Error listing secrets:", err)
		}

		for _, secret := range secrets {
			if verbose {
				fmt.Println(secret.Name, "created:", secret.CreatedAt.Format("2006-01-02T15:04:05"), "updated:", secret.UpdatedAt.Format("2006-01-02T15:04:05"))
			} else {
				fmt.Println(secret.Name)
			}
		}
	}
}
