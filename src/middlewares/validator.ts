import logger from "../utils/logger";
import { Request, Response, NextFunction } from "express";
import { AnyZodObject, object, string, ZodIssueCode } from "zod";

const validate = (schema: AnyZodObject) => (req: Request, res: Response, next: NextFunction) => {
  try {
    const apiSchema = schema.extend({
      apikey: string().superRefine((val, ctx) => {
        if (val != process.env.APIKEY) {
          ctx.addIssue({
            code: ZodIssueCode.custom,
            message: "Invalid API key",
          });
        }
      }),
    });
    req.body = apiSchema.parse({
      body: req.body,
      apikey: req.headers.apikey,
    }).body;
    next();
  } catch (e: any) {
    return res.status(400).jsend.fail(e.errors);
  }
};

export default validate;
