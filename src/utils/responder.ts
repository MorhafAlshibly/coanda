import { Response } from "express";
import mongoose from "mongoose";
import { ZodError, ZodIssue } from "zod";
import { IssueCode, IssueStatus } from "../schemas/issues";

export enum Status {
  success = "success",
  invalid = "invalid",
  fail = "fail",
  error = "error",
}

export class Success {
  data: object | Array<object> | mongoose.Types.ObjectId;
  constructor(data: object | Array<object> | mongoose.Types.ObjectId) {
    this.data = data;
  }

  send(res: Response) {
    res.status(200).json({
      status: Status.success,
      data: this.data,
    });
  }
}

export class Invalid {
  errors: ZodIssue[];
  constructor(e: ZodError) {
    this.errors = e.errors;
  }

  send(res: Response) {
    res.status(400).json({
      status: Status.invalid,
      errors: this.errors,
    });
  }
}

export class Fail {
  code: IssueCode;
  message: string;
  constructor(code: IssueCode, message: string) {
    this.code = code;
    this.message = message;
  }

  send(res: Response) {
    res.status(IssueStatus[this.code]).json({
      status: Status.fail,
      code: this.code,
      message: this.message,
    });
  }
}

export class Error {
  message: string;
  constructor(message: string) {
    this.message = message;
  }

  send(res: Response) {
    res.status(500).json({
      status: Status.error,
      message: this.message,
    });
  }
}
