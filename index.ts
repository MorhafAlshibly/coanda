import dotenv from "dotenv";
dotenv.config();

import { Request, Response } from "express";
import config from "config";
import connect from "./src/utils/connect";
import logger from "./src/utils/logger";
import { server } from "./src/utils/server";
import swaggerDocs from "./src/utils/swagger";

const port = config.get<number>("port");

const app = server();

app.get("/ping", (req: Request, res: Response) => {
  res.sendStatus(200);
});

const listener = app.listen(port, async () => {
  logger.info(`Coanda API is running on port ${port}`);
  await connect();
  swaggerDocs(app, port);
});
listener.setTimeout(config.get<number>("timeout"));
