import { NextFunction, Request, Response } from "express";
import { Fail } from "../utils/responder";

const parse = (err: Error, req: Request, res: Response, next: NextFunction) => {
  if (err instanceof SyntaxError && "body" in err) {
    return new Fail("syntax_error", err.message).send(res);
  }
  next();
};

export default parse;
