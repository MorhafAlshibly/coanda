package flags

import (
	"github.com/peterbourgon/ff/v4"
)

type GlobalFlags struct {
	FlagSet       *ff.FlagSet
	Environment   *string
	VaultConn     *string
	TableSecret   *string
	TableConn     *string
	MongoSecret   *string
	MongoConn     *string
	MongoDatabase *string
}

func GetGlobalFlags() (*GlobalFlags, error) {
	fs := ff.NewFlagSet("global")
	gf := &GlobalFlags{
		FlagSet:       fs,
		Environment:   fs.String('e', "environment", "dev", "the environment the service is running in"),
		VaultConn:     fs.StringLong("vaultConn", "", "the secret connection string"),
		TableSecret:   fs.StringLong("tableSecret", "", "the name of the secret containing the table connection string"),
		TableConn:     fs.StringLong("tableConn", "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", "the connection string to the table storage"),
		MongoSecret:   fs.StringLong("mongoSecret", "", "the name of the secret containing the mongo connection string"),
		MongoConn:     fs.StringLong("mongoConn", "mongodb://localhost:27017", "the connection string to the mongo database"),
		MongoDatabase: fs.StringLong("mongoDatabase", "coanda", "the name of the mongo database"),
	}
	err := ff.Parse(fs, nil, ff.WithEnvVarPrefix("GLOBAL"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		return nil, err
	}
	return gf, nil
}
