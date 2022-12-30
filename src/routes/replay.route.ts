import express from "express";
import { cacheMiddleware } from "../middlewares/cache";
import validator from "../middlewares/validator";
import { createReplayHandler, getReplayHandler } from "../controllers/replay.controller";
import { createReplaySchema, getReplaySchema } from "../schemas/replay.schema";
let router = express.Router();

/**
 * @openapi
 * '/replay/create':
 *  post:
 *     tags:
 *     - Replay
 *     summary: Create a replay
 *     requestBody:
 *      required: true
 *      content:
 *        application/json:
 *           schema:
 *              $ref: '#/components/schemas/CreateReplayInput'
 *     responses:
 *      200:
 *        description: Success
 *        content:
 *          application/json:
 *            schema:
 *              $ref: '#/components/schemas/CreateUserResponse'
 *      409:
 *        description: Conflict
 *      400:
 *        description: Bad request
 */
router.post("/create", validator(createReplaySchema), createReplayHandler);
router.get("/get", validator(getReplaySchema), cacheMiddleware("_id"), getReplayHandler);

export default router;
