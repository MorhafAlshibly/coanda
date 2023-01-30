import { Database } from "@azure/cosmos";

const replay = async (database: Database) => {
	const { container, resource: definition } = await database.containers.createIfNotExists({ id: "Replay", partitionKey: "/id" });
	if (definition != undefined) {
		definition.defaultTtl = 86400;
		await container.replace(definition);
	}
	console.log(container.id);
};

export default replay;
