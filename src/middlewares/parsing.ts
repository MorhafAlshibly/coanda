import { NextFunction, Request, Response } from "express";
import { failModel } from "../utils/responses";

const parse = (err: Error, req: Request, res: Response, next: NextFunction) => {
  if (err instanceof SyntaxError && "body" in err) {
    return res.status(400).jsend.fail(failModel("syntax_error", ["body"], err.message));
  }
  next();
};

export default parse;
