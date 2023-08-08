// Table storage
import { TableClient, AzureNamedKeyCredential, TableEntity, TableEntityResult } from "@azure/data-tables";

export class TableStorage {
	private readonly _tableClient: TableClient;

	constructor(connectionString: string, tableName: string) {
		this._tableClient = TableClient.fromConnectionString(connectionString, tableName, { allowInsecureConnection: true });
	}

	public async getAllItems(): Promise<TableEntityResult<Record<string, unknown>>[]> {
		const items = [];
		for await (const item of this._tableClient.listEntities()) {
			items.push(item);
		}
		return items;
	}
}
