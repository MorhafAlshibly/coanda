import { Request, Response, NextFunction } from "express";
import { AnyZodObject, ZodError } from "zod";
import { InvalidRes } from "../utils/responses";

// Middleware to validate request body based on Zod schema
const validate = (schema: AnyZodObject) => (req: Request, res: Response, next: NextFunction) => {
	try {
		// Parsing request body into validation schema
		req.body = schema.parse({
			body: req.body,
		}).body;
		next();
	} catch (e: unknown) {
		if (e instanceof ZodError) return new InvalidRes(e.errors).send(res);
	}
};

export default validate;
