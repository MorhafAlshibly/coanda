import { DocumentDefinition, FilterQuery, QueryOptions } from "mongoose";
import ReplayModel, { ReplayDocument } from "../models/replay.model";

// Creating the replay
export const createReplay = async (input: DocumentDefinition<Omit<ReplayDocument, "createdAt" | "updatedAt">>) => {
	try {
		return await ReplayModel.create(input);
	} catch (e: any) {
		throw new Error(e);
	}
};

// Getting the replay
export const getReplay = async (query: FilterQuery<ReplayDocument>, options: QueryOptions = { lean: true }) => {
	try {
		return await ReplayModel.findOne(query, {}, options);
	} catch (e: any) {
		throw new Error(e);
	}
};
