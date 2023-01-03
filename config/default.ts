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
	// Authentication
	auth: {
		message: "Invalid API key",
	},
	// Express options
	express: {
		port: 5050,
		timeout: 5000,
		sizeLimit: "100mb",
		message: "Coanda API has started",
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
	// MongoDB options
	mongodb: {
		message: "Connected to Coanda DB",
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
			routes: "src/routes/*.ts",
			responses: "src/responses",
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
	// Replay endpoints options
	replay: {
		createReplay: {
			minDate: "Must be a date in the future",
			properties: {
				data: "The replay data",
				expireAt: "Expiry timestamp",
			},
		},
		getReplay: {
			properties: {
				_id: "The replay _id",
			},
		},
	},
};
