import { Express, Request, Response } from "express";

import replay from "./replay.route";

const routes = (app: Express) => {
  /**
   * @openapi
   * /ping:
   *  get:
   *     tags:
   *     - General
   *     description: Responds if the app is up and running
   *     responses:
   *       200:
   *         description: App is up and running
   */
  app.get("/ping", (req: Request, res: Response) => {
    res.sendStatus(200);
  });

  // Routes
  app.use("/replay", replay);
};

export default routes;
