package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
)

const projectID = "mlab-sandbox"
const namespace = "reboot-api"
const filename = "import"

type Credentials struct {
	Hostname string `datastore:"hostname"`
	Username string `datastore:"username"`
	Password string `datastore:"password"`
	Model    string `datastore:"model"`
	Address  string `datastore:"address"`
}

func AddCredentials(ctx context.Context, client *datastore.Client,
	host string, user string, pass string, model string,
	ip string) (*datastore.Key, error) {
	credentials := &Credentials{
		Hostname: host,
		Username: user,
		Password: pass,
		Model:    model,
		Address:  ip,
	}

	key := datastore.IncompleteKey("Credentials", nil)
	key.Namespace = namespace
	return client.Put(ctx, key, credentials)
}

func deleteAll(ctx context.Context, client *datastore.Client) {
	q := datastore.NewQuery("Credentials").Namespace(namespace)

	var creds []*Credentials
	keys, err := client.GetAll(ctx, q, &creds)

	if err != nil {
		log.Fatal(err)
	}

	client.DeleteMulti(ctx, keys)
}

func main() {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	deleteAll(ctx, client)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(), " ")
		_, err := AddCredentials(ctx, client, splitted[0], splitted[1],
			splitted[2], splitted[3], splitted[4])

		if err != nil {
			log.Fatal(err)
		}
	}

	if scanner.Err() != nil {
		log.Printf("Failed to read input file: %v\n", scanner.Err())
	}
}
