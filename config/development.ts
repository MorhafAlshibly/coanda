import { name, version } from "../package.json";
import { ZodIssueCode } from "zod";

export default {
	// Info about the project
	info: {
		license: {
			name: "Apache 2.0",
			url: "http://www.apache.org/licenses/LICENSE-2.0.html",
		},
	},
	// Express options
	express: {
		timeout: 5000,
		sizeLimit: "10kb",
		message: "Microservice has started",
		// Syntax error body
		syntaxError: [
			{
				code: ZodIssueCode.invalid_type,
				path: ["body"],
				expected: "object",
				received: "unknown",
				message: "Expected object, received unknown",
			},
		],
	},
	// CosmosDB options
	cosmosdb: {
		message: "Connected to Coanda DB",
		errorMessage: "Unable to connect to Coanda DB",
	},
	// Terraform variables
	terraform: {
		environment: "dev",
		location: "eastus",
		resource_group_name: "coanda-resources",
		key_vault_name: "coandakv",
		cosmosdb_account_name: "coanda-cosmosdb",
		cosmosdb_secret_name: "cosmosdb-connection-string",
		cosmosdb_main_database_name: "coanda-cosmosdb-main",
		cosmosdb_replays_collection_name: "Replays",
		cosmosdb_replays_collection_default_ttl_seconds: 86400,
	},
	// Utility modules
	utils: {
		secrets: {
			errorMessage: "Unable to connect to Azure Key Vault",
		},
	},
	// Swagger options
	swagger: {
		successMessage: "Success",
		failMessage: "There is an issue with the request",
		invalidMessage: "Invalid input data",
		errorMessage: "Temporary server error",
		unauthorizedMessage: "Unauthorized access",
		// Paths to schemas
		paths: {
			routes: "./src/microservices/**/router.ts",
			responses: "./**/responses.ts",
			output: "./docs/src/swagger.json",
		},
		// Base swagger definition
		definition: {
			openapi: "3.1.0",
			info: {
				title: name.charAt(0).toUpperCase() + name.slice(1),
				version: version,
				license: {
					name: "Apache 2.0",
					url: "http://www.apache.org/licenses/LICENSE-2.0.html",
				},
			},
			components: {
				securitySchemes: {
					ApiKeyAuth: {
						type: "apiKey",
						in: "header",
						name: "apikey",
					},
				},
			},
			servers: [
				{
					url: "https://{tenant}",
					variables: {
						tenant: {
							default: "api.example.com",
							description: "The server host",
						},
					},
				},
			],
		},
	},
	// Microservices config
	microservices: {
		// General endpoints
		general: {
			port: 5050,
		},
		// Replay endpoints options
		replays: {
			port: 5055,
			createReplay: {
				minDate: "Must be a date in the future",
				properties: {
					data: "The replay data",
					userId: "The user's unique identifier",
				},
			},
			getReplay: {
				properties: {
					_id: "The replay _id",
				},
			},
		},
	},
};
