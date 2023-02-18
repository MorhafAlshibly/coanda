import { SecretClient } from "@azure/keyvault-secrets";
import { DefaultAzureCredential } from "@azure/identity";

const secrets = async () => {
	const credential = new DefaultAzureCredential();

	const keyVaultName = "coandakv";
	const url = "https://" + keyVaultName + ".vault.azure.net";

	return new SecretClient(url, credential);
};

export const cosmosUri = async () => {
	const client = await secrets();
	const secret = await client.getSecret("cosmosdb-connection-string");
	return secret.value;
};
