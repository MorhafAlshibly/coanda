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
 *    tags:
 *    - Replay
 *    summary: Create a replay
 *    operationId: createReplay
 *    requestBody:
 *      required: true
 *      content:
 *        application/json:
 *          schema:
 *            $ref: '#/components/schemas/CreateReplayInput'
 *    responses:
 *      200:
 *        $ref: '#/components/responses/CreateReplaySuccess'
 *      400:
 *        $ref: '#/components/responses/Invalid'
 *      500:
 *        $ref: '#/components/responses/Error'
 *    security:
 *      - ApiKeyAuth: []
 */
router.post("/create", validator(createReplaySchema), createReplayHandler);

/**
 * @openapi
 * '/replay/get':
 *  get:
 *    tags:
 *    - Replay
 *    summary: Get a replay
 *    operationId: getReplay
 *    requestBody:
 *      required: true
 *      content:
 *        application/json:
 *          schema:
 *            $ref: '#/components/schemas/GetReplayInput'
 *    responses:
 *      200:
 *        $ref: '#/components/responses/GetReplaySuccess'
 *      404:
 *        $ref: '#/components/responses/GetReplayFail'
 *      400:
 *        $ref: '#/components/responses/Invalid'
 *      500:
 *        $ref: '#/components/responses/Error'
 *    security:
 *      - ApiKeyAuth: []
 */
router.get("/get", validator(getReplaySchema), cacheMiddleware("_id"), getReplayHandler);

export default router;
