package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo(uri, caPath string) (*mongo.Collection, error) {
	roots := x509.NewCertPool()
	pem, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	roots.AppendCertsFromPEM(pem)

	clientOpts := options.Client().ApplyURI(uri).SetTLSConfig(&tls.Config{RootCAs: roots})
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		return nil, err
	}
	return client.Database("blog").Collection("commments"), nil

}
