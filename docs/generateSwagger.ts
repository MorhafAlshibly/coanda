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

export const generateSwagger = async () => {
  const jsdocOptions = {
    definition,
    apis: ["src/routes/*.ts"],
  };
  const oas: any = await swaggerJsdoc(jsdocOptions as swaggerJsdoc.Options);

  const paths = new fdir().withFullPaths().crawl("src/responses").sync() as string[];
  const program = TJS.getProgramFromFiles(paths);
  const generator = TJS.buildGenerator(program);

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
      oas.components.schemas = { ...oas.components.schemas, ...failResponse?.definitions };
      // @ts-ignore
      const failCodes = failResponse?.definitions[Object.keys(failResponse?.definitions)[0]].enum;
      // @ts-ignore
      failResponse?.properties?.statusCode["$ref"] = failResponse?.properties?.statusCode["$ref"].replace("definitions", "components/schemas");

      for (let k = 0; k < failCodes.length; k++) {
        const specificFailure = failResponse;
        // @ts-ignore
        specificFailure.properties.data.default = specificFailure?.properties.data.enum[k];
        // @ts-ignore
        delete specificFailure?.properties.data.enum;
        oas.paths[endpoint][requestType].responses[failCodes[k]] = {
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
  fs.writeFileSync("./docs/src/swagger.json", JSON.stringify(oas));
});
