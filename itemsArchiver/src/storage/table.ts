// Table storage
import { TableClient, TableEntityResult, TableServiceClientOptions } from "@azure/data-tables";

export default class TableStorage {
	private readonly _tableClient: TableClient;

	constructor(connectionString: string, tableName: string, options: TableServiceClientOptions = {}) {
		this._tableClient = TableClient.fromConnectionString(connectionString, tableName, { allowInsecureConnection: true });
	}

	public async getAllItems(): Promise<TableEntityResult<Record<string, unknown>>[]> {
		const items = [] as TableEntityResult<Record<string, unknown>>[];
		for await (const item of this._tableClient.listEntities()) {
			items.push(item);
		}
		return items;
	}

	public async deleteItem(rowKey: string, partitionKey: string): Promise<void> {
		await this._tableClient.deleteEntity(partitionKey, rowKey);
	}
}
