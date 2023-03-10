/* eslint-disable @typescript-eslint/ban-ts-comment */
import * as TJS from "typescript-json-schema";
import { fdir } from "fdir";
import config from "config";
import fs from "fs";
import swaggerJsdoc from "swagger-jsdoc";
import requestSchemas from "./requestSchemas";
import { basicResponses } from "./responseSchemas";

const definition = config.get<object>("swagger.definition");
const successMessage = config.get<string>("swagger.successMessage");
const failMessage = config.get<string>("swagger.failMessage");

const settings: TJS.PartialArgs = {
	ref: false,
};

export const generateSwagger = async () => {
	const jsdocOptions = {
		definition,
		apis: [config.get<string>("swagger.paths.routes")],
	};
	const oas: any = await swaggerJsdoc(jsdocOptions as swaggerJsdoc.Options);

	const paths = new fdir().withFullPaths().crawl(config.get<string>("swagger.paths.responses")).sync() as string[];
	const program = TJS.getProgramFromFiles(paths);
	const generator = TJS.buildGenerator(program, settings);

	oas.components.schemas = requestSchemas;
	oas.components.responses = basicResponses;

	for (let i = 0; i < Object.keys(oas.paths).length; i++) {
		const endpoint = Object.keys(oas.paths)[i];
		for (let j = 0; j < Object.keys(oas.paths[endpoint]).length; j++) {
			const requestType = Object.keys(oas.paths[endpoint])[j];

			const successResponseName = oas.paths[endpoint][requestType].responses.success;
			const failResponseName = oas.paths[endpoint][requestType].responses.fail;

			if (!successResponseName && !failResponseName) continue;

			oas.paths[endpoint][requestType].responses = {
				"400": { $ref: "#/components/responses/Invalid" },
				"401": { $ref: "#/components/responses/Unauthorized" },
				"500": { $ref: "#/components/responses/Error" },
			};

			if (successResponseName) {
				oas.paths[endpoint][requestType].responses["200"] = {
					description: successMessage,
					content: {
						"application/json": {
							schema: generator?.getSchemaForSymbol(successResponseName),
						},
					},
				};
			}

			if (!failResponseName) continue;

			const failResponse = generator?.getSchemaForSymbol(failResponseName);

			// @ts-ignore
			const failStatusCodes = failResponse?.properties.statusCode.enum;
			// @ts-ignore
			const failCodes = failResponse?.properties.data.enum;

			for (let k = 400; k < 500; k++) {
				if (!failStatusCodes.includes(k)) continue;
				const specificFailure = JSON.parse(JSON.stringify(failResponse));
				// @ts-ignore
				specificFailure?.properties.data.enum = [];
				// @ts-ignore
				specificFailure?.properties.statusCode.enum = [k];
				for (let l = 0; l < failStatusCodes.length; l++) {
					if (failStatusCodes[l] == k) {
						// @ts-ignore
						specificFailure?.properties.data.enum.push(failCodes[l]);
					}
				}

				oas.paths[endpoint][requestType].responses[k] = {
					description: failMessage,
					content: {
						"application/json": {
							schema: specificFailure,
						},
					},
				};
			}
		}
	}
	return oas;
};

generateSwagger().then((oas) => {
	fs.writeFileSync(config.get<string>("swagger.paths.output"), JSON.stringify(oas));
});
