import { NextFunction, Request, Response } from "express";
import { ZodError, ZodIssueCode } from "zod";
import { Invalid } from "../utils/responder";

const parse = (err: Error, req: Request, res: Response, next: NextFunction) => {
  if (err instanceof SyntaxError && "body" in err) {
    return new Invalid(
      new ZodError([
        {
          code: ZodIssueCode.invalid_type,
          path: ["body"],
          expected: "object",
          received: "unknown",
          message: "Expected object, received unknown",
        },
      ])
    ).send(res);
  }
  next();
};

export default parse;
