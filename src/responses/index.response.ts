import { Response } from "express";
import { ZodIssue } from "zod";

// Parent class with send response method
export class Responder {
	statusCode: any;
	status: any;
	data: any;
	send(res: Response) {
		res.status(this.statusCode).send({
			statusCode: this.statusCode,
			status: this.status,
			data: this.data,
		});
	}
}

// A successful response
export class SuccessRes extends Responder {
	statusCode = 200;
	status = "success";
	data: any;
	constructor() {
		super();
	}
}

// An invalid response
export class InvalidRes extends Responder {
	statusCode = 400;
	status = "invalid";
	/**
	 * The array of ZodIssue objects
	 * @TJS-type array
	 * @items.type object
	 */
	data: ZodIssue[];
	constructor(data: ZodIssue[]) {
		super();
		this.data = data;
	}
}

// An errored response
export class ErrorRes extends Responder {
	statusCode = 500;
	status = "error";
	/**
	 * Error message
	 */
	data: string;
	constructor(data: string) {
		super();
		this.data = data;
	}
}
