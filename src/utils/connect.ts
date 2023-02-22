import config from "config";
import mongoose from "mongoose";
import logger from "./logger";
import { cosmosSecret } from "./secrets";

// Connecting to MongoDB database
const connect = async () => {
	try {
		mongoose.set("strictQuery", false);
		// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
		const cosmosUri = (await cosmosSecret())!.replace("/?", "/" + config.get<string>("terraform.cosmosdb_main_database_name") + "?");
		await mongoose.connect(cosmosUri);
		logger.info(config.get<string>("cosmosdb.message"));
	} catch (error) {
		logger.error(config.get<string>("cosmosdb.errorMessage"));
		process.exit(1);
	}
};

export default connect;
