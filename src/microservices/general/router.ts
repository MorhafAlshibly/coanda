import express from "express";
import { Response, Request } from "express";
const router = express.Router();

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
router.get("/ping", (req: Request, res: Response) => {
	res.sendStatus(200);
});

export default router;
