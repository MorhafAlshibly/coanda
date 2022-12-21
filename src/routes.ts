import { Express, Request, Response } from "express";
import { createReplayHandler } from "./controller/replay.controller";
import validator from "./middleware/validator";
import createReplaySchema from "./schema/replay.schema";

const routes = (app: Express) => {
  app.get("/ping", (req: Request, res: Response) => {
    res.sendStatus(200);
  });
  app.post("/replay/create", validator(createReplaySchema), createReplayHandler);
};

export default routes;
