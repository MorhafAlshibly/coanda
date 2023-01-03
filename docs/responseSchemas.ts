import * as TJS from "typescript-json-schema";
import { fdir } from "fdir";
import config from "config";

const paths = new fdir().withFullPaths().crawl(config.get<string>("swagger.paths.responses")).sync() as string[];
const program = TJS.getProgramFromFiles(paths);
const generator = TJS.buildGenerator(program);

export const basicResponses = {
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
