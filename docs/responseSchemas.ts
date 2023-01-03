import * as TJS from "typescript-json-schema";
import { fdir } from "fdir";
import config from "config";

const successMessage = config.get<string>("swagger.successMessage");
const invalidMessage = config.get<string>("swagger.invalidMessage");
const errorMessage = config.get<string>("swagger.errorMessage");

const paths = new fdir().withFullPaths().crawl(config.get<string>("swagger.paths.responses")).sync() as string[];
const program = TJS.getProgramFromFiles(paths);
const generator = TJS.buildGenerator(program);

export const basicResponses = {
	Success: {
		description: successMessage,
		content: {
			"application/json": {
				schema: generator?.getSchemaForSymbol("SuccessRes"),
			},
		},
	},
	Invalid: {
		description: invalidMessage,
		content: {
			"application/json": {
				schema: generator?.getSchemaForSymbol("InvalidRes"),
			},
		},
	},
	Error: {
		description: errorMessage,
		content: {
			"application/json": {
				schema: generator?.getSchemaForSymbol("ErrorRes"),
			},
		},
	},
};
