import { Request, Response } from "express";
import server from "../../utils/server";

const app = server();

/**
 * @openapi
 * /ping:
 *  get:
 *    tags:
 *    - General
 *    operationId: ping
 *    summary: Ping the API
 *    description: Responds if the app is up and running
 *    responses:
 *      200:
 *        description: App is up and running
 */
app.get("/ping", (req: Request, res: Response) => {
	res.sendStatus(200);
});
