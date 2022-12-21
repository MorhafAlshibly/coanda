import mongoose from "mongoose";
import logger from "./logger";

async function connect() {
  try {
    mongoose.set("strictQuery", false);
    await mongoose.connect(process.env.MONGOURI);
    logger.info("Connected to Coanda DB!");
  } catch (error) {
    logger.error("Could not connect to DB");
    process.exit(1);
  }
}

export default connect;
