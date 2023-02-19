import config from "config";
import mongoose from "mongoose";
import logger from "./logger";
import { cosmosSecret } from "./secrets";

// Connecting to MongoDB database
const connect = async () => {
	try {
		mongoose.set("strictQuery", false);
		// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
		const cosmosUri = (await cosmosSecret())!.replace("/?", "/coanda-cosmosdb-main?");
		await mongoose.connect(cosmosUri);
		logger.info(config.get<string>("mongodb.message"));
	} catch (error) {
		logger.error("Unable to connect to Cosmos DB");
		process.exit(1);
	}
};

export default connect;
