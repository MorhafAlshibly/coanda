import config from "config";
import mongoose from "mongoose";
import { object, date, string, preprocess, ZodIssueCode, TypeOf } from "zod";
import { Responses } from "../utils/responder";

export const createReplaySchema = object({
  body: object({
    data: object({})
      .passthrough()
      .superRefine((val, ctx) => {
        if (Object.keys(val).length == 0)
          ctx.addIssue({
            code: ZodIssueCode.invalid_type,
            expected: "object",
            received: "undefined",
          });
      }),
    expireAt: preprocess(
      (arg) => {
        if (typeof arg == "number" || typeof arg == "string" || arg instanceof Date) return new Date(arg);
      },
      date().min(new Date(), {
        message: config.get<string>("replay.createReplay.minDate"),
      })
    ).default(new Date(9999999999999)),
  }),
});

/**
 * @openapi
 * components:
 *  schemas:
 *    CreateReplayInput:
 *      type: object
 *      required:
 *        - data
 *      properties:
 *        data:
 *          type: object
 *        expireAt:
 *          oneOf:
 *            - type: number
 *            - type: string
 *          description: Timestamp in ms
 *          default: 9999999999999
 */
export type CreateReplayInput = TypeOf<typeof createReplaySchema>;

/**
 * @openapi
 * components:
 *  responses:
 *    CreateReplaySuccess:
 *      description: Success.
 *      content:
 *        application/json:
 *          schema:
 *            type: object
 *            properties:
 *              status:
 *                type: string
 *                const: success
 *              data:
 *                type: string
 *                description: "_id of the replay"
 */
export class CreateReplaySuccess extends Responses {
  constructor(_id: mongoose.Types.ObjectId) {
    super("success", 200, _id);
  }
}

export const getReplaySchema = object({
  body: object({
    _id: string().transform((val, ctx) => {
      if (mongoose.Types.ObjectId.isValid(val)) return new mongoose.Types.ObjectId(val);
      else
        ctx.addIssue({
          code: ZodIssueCode.invalid_type,
          expected: "string",
          received: "unknown",
        });
    }),
  }),
});

/**
 * @openapi
 * components:
 *  schemas:
 *    GetReplayInput:
 *      type: object
 *      required:
 *        - _id
 *      properties:
 *        _id:
 *          type: string
 *          description: _id of the replay
 */

export type GetReplayInput = TypeOf<typeof getReplaySchema>;

/**
 * @openapi
 * components:
 *  responses:
 *    GetReplaySuccess:
 *      description: Success.
 *      content:
 *        application/json:
 *          schema:
 *            type: object
 *            properties:
 *              status:
 *                type: string
 *                const: success
 *              data:
 *                type: object
 *                description: "The data of the replay"
 */
export class GetReplaySuccess extends Responses {
  constructor(data: object) {
    super("success", 200, data);
  }
}

/**
 * @openapi
 * components:
 *  schemas:
 *    GetReplayIssue:
 *      oneOf:
 *        - const: replay_not_found
 *          type: string
 *          description: "Replay data not found"
 */
export enum GetReplayIssue {
  replay_not_found = 404,
}
export type GetReplayIssueCode = keyof typeof GetReplayIssue;

/**
 * @openapi
 * components:
 *  responses:
 *    GetReplayFail:
 *      description: There is an issue with the request.
 *      content:
 *        application/json:
 *          schema:
 *            type: object
 *            properties:
 *              status:
 *                type: string
 *                const: fail
 *              data:
 *                $ref: "#/components/schemas/GetReplayIssue"
 */
export class GetReplayFail extends Responses {
  constructor(issue: GetReplayIssueCode) {
    super("fail", GetReplayIssue[issue], issue);
  }
}
