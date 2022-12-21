import dotenv from "dotenv";
dotenv.config();

import express from "express";
import cors from "cors";
import config from "config";
import connect from "./utils/connect";
import logger from "./utils/logger";
import routes from "./routes";

const port = config.get<number>("port");

const app = express();
app.use(cors());
app.use(
  express.json({
    limit: config.get<string>("sizeLimit"),
  })
);

const listener = app.listen(port, async () => {
  logger.info(`Coanda API is running on port ${port}`);
  await connect();
  routes(app);
});
listener.setTimeout(config.get<number>("timeout"));
