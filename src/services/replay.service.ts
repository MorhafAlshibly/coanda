import { DocumentDefinition } from "mongoose";
import ReplayModel, { ReplayDocument } from "../models/replay.model";

export async function createReplay(input: DocumentDefinition<Omit<ReplayDocument, "createdAt" | "updatedAt">>) {
  try {
    return await ReplayModel.create(input);
  } catch (e: any) {
    throw new Error(e);
  }
}
