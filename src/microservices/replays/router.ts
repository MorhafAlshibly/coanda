import express from "express";
import validator from "../../middlewares/validator";
import { createReplayHandler, getReplayHandler } from "./controller";
import { createReplaySchema, getReplaySchema } from "./schemas";
const router = express.Router();

/**
 * @openapi
 * '/replay':
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
 *      success: CreateReplaySuccess
 *    security:
 *      - ApiKeyAuth: []
 */
router.post("/", validator(createReplaySchema), createReplayHandler);

/**
 * @openapi
 * '/replay':
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
 *      success: GetReplaySuccess
 *      fail: GetReplayFail
 *    security:
 *      - ApiKeyAuth: []
 */
router.get("/", validator(getReplaySchema), getReplayHandler);

export default router;
