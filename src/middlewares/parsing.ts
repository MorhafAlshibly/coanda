import config from "config";
import { NextFunction, Request, Response } from "express";
import { ZodError, ZodIssue } from "zod";
import { Invalid } from "../responses/index.response";

const parse = (err: Error, req: Request, res: Response, next: NextFunction) => {
  if (err instanceof SyntaxError && "body" in err) {
    return new Invalid(new ZodError(config.get<ZodIssue[]>("express.syntaxError")).errors).send(res);
  }
  next();
};

export default parse;
