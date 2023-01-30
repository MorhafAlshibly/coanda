import { CosmosClient } from "@azure/cosmos";
import replay from "./replay/index.replay";

const main = async (client: CosmosClient) => {
	const { database } = await client.databases.createIfNotExists({ id: "Main", maxThroughput: 1000 });
	console.log(database.id);
	await replay(database);
};

export default main;
