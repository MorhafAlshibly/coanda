import zodToJsonSchema from "zod-to-json-schema";
import * as TJS from "typescript-json-schema";
import { fdir } from "fdir";
import config from "config";

import { createReplaySchema, getReplaySchema } from "../src/microservices/replays/schemas";

const settings: TJS.PartialArgs = {
	ref: false,
};

const paths = new fdir().withFullPaths().glob(config.get<string>("swagger.paths.responses")).crawl("./src").sync() as string[];
const program = TJS.getProgramFromFiles(paths);
export const generator = TJS.buildGenerator(program, settings);

export const responses = {
	Success: {
		description: config.get<string>("swagger.successMessage"),
		content: {
			"application/json": {
				schema: generator?.getSchemaForSymbol("SuccessRes"),
			},
		},
	},
	Invalid: {
		description: config.get<string>("swagger.invalidMessage"),
		content: {
			"application/json": {
				schema: generator?.getSchemaForSymbol("InvalidRes"),
			},
		},
	},
	Error: {
		description: config.get<string>("swagger.errorMessage"),
		content: {
			"application/json": {
				schema: generator?.getSchemaForSymbol("ErrorRes"),
			},
		},
	},
	Unauthorized: {
		description: config.get<string>("swagger.unauthorizedMessage"),
	},
};

export const requests = {
	...zodToJsonSchema(createReplaySchema.shape.body, "CreateReplayInput").definitions,
	...zodToJsonSchema(getReplaySchema.shape.body, "GetReplayInput").definitions,
};
