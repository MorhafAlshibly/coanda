package flags

import (
	"flag"
	"os"

	"github.com/peterbourgon/ff"
)

type GlobalFlags struct {
	Environment   *string
	VaultConn     *string
	TableSecret   *string
	TableConn     *string
	MongoSecret   *string
	MongoConn     *string
	MongoDatabase *string
}

func GetGlobalFlags() (*GlobalFlags, error) {
	fs := flag.NewFlagSet("global", flag.ContinueOnError)
	gf := &GlobalFlags{
		Environment:   fs.String("environment", "dev", "the environment the service is running in"),
		VaultConn:     fs.String("vaultConn", "", "the secret connection string"),
		TableSecret:   fs.String("tableSecret", "", "the name of the secret containing the table connection string"),
		TableConn:     fs.String("tableConn", "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", "the connection string to the table storage"),
		MongoSecret:   fs.String("mongoSecret", "", "the name of the secret containing the mongo connection string"),
		MongoConn:     fs.String("mongoConn", "mongodb://localhost:27017", "the connection string to the mongo database"),
		MongoDatabase: fs.String("mongoDatabase", "coanda", "the name of the mongo database"),
	}
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("GLOBAL"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		return nil, err
	}
	return gf, nil
}
