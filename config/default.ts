import { name, version, license } from "../package.json";

export default {
  info: {
    license: {
      name: "Apache 2.0",
      url: "http://www.apache.org/licenses/LICENSE-2.0.html",
    },
  },
  express: {
    port: 5050,
    timeout: 5000,
    sizeLimit: "100mb",
    message: "Coanda API has started",
  },
  mongodb: {
    message: "Connected to Coanda DB",
  },
  swagger: {
    successMessage: "Success",
    failMessage: "There is an issue with the request",
    invalidMessage: "Invalid input data",
    errorMessage: "Temporary server error",
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
  replay: {
    createReplay: {
      minDate: "Must be a date in the future",
    },
  },
};
