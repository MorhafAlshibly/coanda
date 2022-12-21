import { DocumentDefinition } from "mongoose";
import ReplayModel, { ReplayDocument } from "../models/replay.model";

const createReplay = async (input: DocumentDefinition<Omit<ReplayDocument, "createdAt" | "updatedAt" | "expireAt">>) => {
  try {
    return await ReplayModel.create(input);
  } catch (e: any) {
    throw new Error(e);
  }
};

export default createReplay;
