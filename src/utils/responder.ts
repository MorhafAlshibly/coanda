import { Response } from "express";
import mongoose from "mongoose";
import { IssueCode, IssueStatus } from "./issues";

export type Success = mongoose.Types.ObjectId | object | Array<object>;
export interface Fail {
  code: IssueCode;
  path: Array<string>;
  message: string;
}

export const success = (res: Response, input: Success) => {
  res.status(200).json({
    status: "success",
    data: input,
  });
};

export const fail = (res: Response, input: Fail | Array<Fail>) => {
  const failResult = {
    status: "fail",
    data: input instanceof Array ? input : [input],
  };
  res.status(IssueStatus[failResult.data[0].code]).json(failResult);
};

export const error = (res: Response, message: string) => {
  res.status(500).json({
    status: "error",
    message,
  });
};
