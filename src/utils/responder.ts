import { Response } from "express";
import mongoose from "mongoose";
import { ZodError, ZodIssue } from "zod";

export class Responses {
  status: string;
  statusCode: number;
  data: any;
  constructor(status: string, statusCode: number, data: any) {
    this.status = status;
    this.statusCode = statusCode;
    this.data = data;
  }
  send(res: Response) {
    res.status(this.statusCode).json({
      status: this.status,
      data: this.data,
    });
  }
}

/**
 * @openapi
 * components:
 *  responses:
 *    Invalid:
 *      description: Invalid input data.
 *      content:
 *        application/json:
 *          schema:
 *            type: object
 *            properties:
 *              status:
 *                const: invalid
 *              data:
 *                type: array
 *                items:
 *                  type: object
 *                  description: ZodIssue
 */
export class Invalid extends Responses {
  constructor(e: ZodError) {
    super("invalid", 400, e.errors);
  }
}

/**
 * @openapi
 * components:
 *  responses:
 *    Error:
 *      description: Temporary server error.
 *      content:
 *        application/json:
 *          schema:
 *            type: object
 *            properties:
 *              status:
 *                const: error
 *              data:
 *                type: string
 */
export class Error extends Responses {
  constructor(message: string) {
    super("error", 400, message);
  }
}
