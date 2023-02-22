import { SecretClient } from "@azure/keyvault-secrets";
import { DefaultAzureCredential } from "@azure/identity";
import config from "config";
import logger from "./logger";

// Create secret client for secrets
const secrets = () => {
	try {
		// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
		const credential = new DefaultAzureCredential();
		const url = "https://" + config.get<string>("terraform.key_vault_name") + ".vault.azure.net";
		return new SecretClient(url, credential);
	} catch (error: any) {
		logger.error(config.get<string>("utils.secrets.errorMessage"));
		throw new Error(error);
	}
};

// Export Cosmos secret
export const cosmosSecret = async () => {
	try {
		const client = secrets();
		const secret = await client.getSecret(config.get<string>("terraform.cosmosdb_secret_name"));
		return secret.value;
	} catch (error: any) {
		throw new Error(error);
	}
};
