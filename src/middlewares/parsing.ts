import { NextFunction, Request, Response } from "express";
import { fail } from "../utils/responder";

const parse = (err: Error, req: Request, res: Response, next: NextFunction) => {
  if (err instanceof SyntaxError && "body" in err) {
    return fail(res, { code: "syntax_error", path: ["body"], message: err.message });
  }
  next();
};

export default parse;
