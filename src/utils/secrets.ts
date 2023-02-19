import { SecretClient } from "@azure/keyvault-secrets";
import { DefaultAzureCredential } from "@azure/identity";
import logger from "./logger";

// Create secret client for secrets
const secrets = async () => {
	try {
		const credential = new DefaultAzureCredential();
		const keyVaultName = "coandakv";
		const url = "https://" + keyVaultName + ".vault.azure.net";
		return new SecretClient(url, credential);
	} catch (error) {
		logger.error("Unable to connect to Azure Key Vault");
		process.exit(1);
	}
};

// Export Cosmos secret
export const cosmosSecret = async () => {
	try {
		const client = await secrets();
		const secret = await client.getSecret("cosmosdb-connection-string");
		return secret.value;
	} catch (error) {
		logger.error("Unable to get Cosmos Secret");
		process.exit(1);
	}
};
