import mongoose from "mongoose";
import { Responder, Success } from "./index.response";

export class CreateReplaySuccess extends Success {
  /**
   * The ObjectId of the replay.
   * @TJS-type string
   */
  data: mongoose.Types.ObjectId;
  constructor(data: mongoose.Types.ObjectId) {
    super();
    this.data = data;
  }
}

export class GetReplaySuccess extends Success {
  data: object;
  constructor(data: object) {
    super();
    this.data = data;
  }
}

export class GetReplayFail extends Responder {
  statusCode: GetReplayIssue;
  status = "fail";
  data: GetReplayIssueCode;
  constructor(issue: GetReplayIssueCode) {
    super();
    this.statusCode = GetReplayIssue[issue];
    this.data = issue;
  }
}

export enum GetReplayIssue {
  replay_not_found = 404,
}
export type GetReplayIssueCode = keyof typeof GetReplayIssue;
