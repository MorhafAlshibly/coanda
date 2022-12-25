import { NextFunction, Request, Response } from "express";

const parse = (err: Error, req: Request, res: Response, next: NextFunction) => {
  if (err instanceof SyntaxError && "body" in err) {
    return res.status(400).send({ success: false, message: err.message });
  }
  next();
};

export default parse;
