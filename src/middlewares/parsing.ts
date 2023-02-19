import config from "config";
import { NextFunction, Request, Response } from "express";
import { ZodError, ZodIssue } from "zod";
import { InvalidRes } from "../utils/response";

// Middleware to check response body is a valid JSON
const parse = (err: Error, req: Request, res: Response, next: NextFunction) => {
	if (err instanceof SyntaxError && "body" in err) {
		return new InvalidRes(new ZodError(config.get<ZodIssue[]>("express.syntaxError")).errors).send(res);
	}
	next();
};

export default parse;
