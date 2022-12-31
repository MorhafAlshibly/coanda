const { name, version, license } = require("../package.json");

module.exports = {
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
};
