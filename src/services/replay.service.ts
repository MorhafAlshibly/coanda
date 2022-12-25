import { DocumentDefinition, FilterQuery, QueryOptions } from "mongoose";
import ReplayModel, { ReplayDocument } from "../models/replay.model";

export const createReplay = async (input: DocumentDefinition<Omit<ReplayDocument, "createdAt" | "updatedAt">>) => {
  try {
    return await ReplayModel.create(input);
  } catch (e: any) {
    throw new Error(e);
  }
};

export const getReplay = async (query: FilterQuery<ReplayDocument>, options: QueryOptions = { lean: true }) => {
  try {
    const replay = await ReplayModel.findOne(query, {}, options);
    if (!replay)
      return [
        {
          code: "not_found",
          path: ["body", "_id"],
          message: "Data not found",
        },
      ];
    else return replay;
  } catch (e: any) {
    throw new Error(e);
  }
};
