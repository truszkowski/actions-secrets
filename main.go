package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"
	"os"

	"github.com/google/go-github/v59/github"
	"github.com/jamesruan/sodium"
)

func main() {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(os.Getenv("TOKEN"))

	var owner, repo string
	flag.StringVar(&owner, "owner", "", "owner of the repository")
	flag.StringVar(&repo, "repo", "", "name of the repository")
	flag.Parse()

	pubKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		log.Fatalln("Error fetching public key:", err)
	}

	log.Println("public key:", pubKey.GetKeyID(), pubKey.GetKey())

	secrets, _, err := client.Actions.ListRepoSecrets(ctx, owner, repo, nil)
	if err != nil {
		log.Fatalln("Error fetching secrets:", err)
	}

	for _, secret := range secrets.Secrets {
		log.Println(secret.Name, secret.CreatedAt)
	}

	rawKey, err := base64.StdEncoding.DecodeString(pubKey.GetKey())
	if err != nil {
		log.Fatalln("Error decoding public key:", err)
	}

	kp := sodium.MakeBoxKP()
	_ = rawKey

	//kp.PublicKey
	//kp.SecretKey

	/*encryptedValue := base64.StdEncoding.EncodeToString(encryptedRaw)

	resp, err := client.Actions.CreateOrUpdateRepoSecret(ctx, "owner", "repo", &github.EncryptedSecret{
		Name:           "TEST_SECRET_TEST",
		KeyID:          pubKey.GetKeyID(),
		EncryptedValue: encryptedValue,
	})*/

	log.Println(resp, err)
}
