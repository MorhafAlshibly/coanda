import mongoose from "mongoose";
import { object, date, string, preprocess, ZodIssueCode, TypeOf } from "zod";
import { Responses } from "../utils/responder";
import { Issue } from "./issues";

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
        message: "Must be a date in the future",
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
 *                const: success
 *              data:
 *                type: string
 *                description: "_id of the replay"
 */
export class CreateReplaySuccess extends Responses {
  constructor(data: object) {
    super("success", 200, data);
  }
}

/**
 * @openapi
 * components:
 *  schemas:
 *    CreateReplayIssue:
 *      type: string
 *      oneOf:
 *        - const: replay_not_found
 *          description: "Replay data not found"
 */
export enum CreateReplayIssue {
  replay_not_found = 404,
}
export type CreateReplayIssueCode = keyof typeof CreateReplayIssue;

/**
 * @openapi
 * components:
 *  responses:
 *    CreateReplayFail:
 *      description: There is an issue with the request.
 *      content:
 *        application/json:
 *          schema:
 *            type: object
 *            properties:
 *              status:
 *                const: fail
 *              data:
 *                $ref: "#/components/schemas/CreateReplayIssue"
 */
export class CreateReplayFail extends Responses {
  constructor(issue: CreateReplayIssueCode) {
    super("fail", CreateReplayIssue[issue], issue);
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

export type GetReplayInput = TypeOf<typeof getReplaySchema>;
