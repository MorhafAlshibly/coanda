import mongoose from "mongoose";
import { Responder, SuccessRes } from "../../utils/responses";

// A successful createReplay
export class CreateReplaySuccess extends SuccessRes {
	/**
	 * The replay _id
	 * @TJS-type string
	 */
	data: mongoose.Types.ObjectId;
	constructor(data: mongoose.Types.ObjectId) {
		super();
		this.data = data;
	}
}

// A successful getReplay
export class GetReplaySuccess extends SuccessRes {
	/**
	 * The replay data
	 */
	data: object;
	constructor(data: object) {
		super();
		this.data = data;
	}
}

// A failed getReplay
export class GetReplayFail extends Responder {
	statusCode: GetReplayIssueCode;
	status = "fail";
	data: GetReplayIssue;
	constructor(issue: GetReplayIssue) {
		super();
		this.statusCode = GetReplayIssueCode[issue];
		this.data = issue;
	}
}

// All possible getReplay issues
export enum GetReplayIssueCode {
	replay_not_found = 404,
}
export type GetReplayIssue = keyof typeof GetReplayIssueCode;
