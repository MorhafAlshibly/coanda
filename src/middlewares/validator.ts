import { Request, Response, NextFunction } from "express";
import { AnyZodObject, string, ZodIssueCode, ZodError } from "zod";
import { Invalid } from "../utils/responder";

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
  } catch (e: unknown) {
    if (e instanceof ZodError) return new Invalid(e).send(res);
  }
};

export default validate;
