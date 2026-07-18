export namespace dbexport {
	
	export class BackupOptions {
	    connection_id: string;
	    database: string;
	    output_path: string;
	    format: string;
	    schema_only: boolean;
	    data_only: boolean;
	    clean: boolean;
	    create: boolean;
	    inserts: boolean;
	    column_inserts: boolean;
	    disable_triggers: boolean;
	    verbose: boolean;
	    pg_bin_path: string;
	
	    static createFrom(source: any = {}) {
	        return new BackupOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.database = source["database"];
	        this.output_path = source["output_path"];
	        this.format = source["format"];
	        this.schema_only = source["schema_only"];
	        this.data_only = source["data_only"];
	        this.clean = source["clean"];
	        this.create = source["create"];
	        this.inserts = source["inserts"];
	        this.column_inserts = source["column_inserts"];
	        this.disable_triggers = source["disable_triggers"];
	        this.verbose = source["verbose"];
	        this.pg_bin_path = source["pg_bin_path"];
	    }
	}
	export class MaintenanceOptions {
	    connection_id: string;
	    database: string;
	    task: string;
	    scope: string;
	    schema_name: string;
	    table_name: string;
	    verbose: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MaintenanceOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.database = source["database"];
	        this.task = source["task"];
	        this.scope = source["scope"];
	        this.schema_name = source["schema_name"];
	        this.table_name = source["table_name"];
	        this.verbose = source["verbose"];
	    }
	}

}

export namespace main {
	
	export class CatalogItem {
	    connection_id: string;
	    connection_name: string;
	    database_name: string;
	    schema: string;
	    name: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new CatalogItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.connection_name = source["connection_name"];
	        this.database_name = source["database_name"];
	        this.schema = source["schema"];
	        this.name = source["name"];
	        this.type = source["type"];
	    }
	}
	export class DefaultConnectionConfig {
	    host: string;
	    port: number;
	    username: string;
	    password: string;
	    database: string;
	
	    static createFrom(source: any = {}) {
	        return new DefaultConnectionConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.database = source["database"];
	    }
	}

}

export namespace models {
	
	export class FilterCondition {
	    column: string;
	    operator: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new FilterCondition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.column = source["column"];
	        this.operator = source["operator"];
	        this.value = source["value"];
	    }
	}
	export class BrowseOptions {
	    page: number;
	    page_size: number;
	    sort_by?: string;
	    sort_order?: string;
	    filters?: FilterCondition[];
	    columns?: string[];
	
	    static createFrom(source: any = {}) {
	        return new BrowseOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	        this.sort_by = source["sort_by"];
	        this.sort_order = source["sort_order"];
	        this.filters = this.convertValues(source["filters"], FilterCondition);
	        this.columns = source["columns"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FKReference {
	    target_schema: string;
	    target_table: string;
	    target_column: string;
	    on_update: string;
	    on_delete: string;
	
	    static createFrom(source: any = {}) {
	        return new FKReference(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.target_schema = source["target_schema"];
	        this.target_table = source["target_table"];
	        this.target_column = source["target_column"];
	        this.on_update = source["on_update"];
	        this.on_delete = source["on_delete"];
	    }
	}
	export class ColumnDefinition {
	    name: string;
	    data_type: string;
	    length?: number;
	    precision?: number;
	    scale?: number;
	    is_nullable: boolean;
	    default_value?: string;
	    is_primary_key: boolean;
	    is_unique: boolean;
	    is_auto_increment: boolean;
	    references?: FKReference;
	    comment?: string;
	
	    static createFrom(source: any = {}) {
	        return new ColumnDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.data_type = source["data_type"];
	        this.length = source["length"];
	        this.precision = source["precision"];
	        this.scale = source["scale"];
	        this.is_nullable = source["is_nullable"];
	        this.default_value = source["default_value"];
	        this.is_primary_key = source["is_primary_key"];
	        this.is_unique = source["is_unique"];
	        this.is_auto_increment = source["is_auto_increment"];
	        this.references = this.convertValues(source["references"], FKReference);
	        this.comment = source["comment"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ColumnInfo {
	    table_schema: string;
	    table_name: string;
	    column_name: string;
	    data_type: string;
	    is_nullable: boolean;
	    default_value?: any;
	    is_primary_key: boolean;
	    ordinal_position: number;
	    character_maximum_length?: number;
	    comment?: string;
	
	    static createFrom(source: any = {}) {
	        return new ColumnInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.table_schema = source["table_schema"];
	        this.table_name = source["table_name"];
	        this.column_name = source["column_name"];
	        this.data_type = source["data_type"];
	        this.is_nullable = source["is_nullable"];
	        this.default_value = source["default_value"];
	        this.is_primary_key = source["is_primary_key"];
	        this.ordinal_position = source["ordinal_position"];
	        this.character_maximum_length = source["character_maximum_length"];
	        this.comment = source["comment"];
	    }
	}
	export class ColumnMeta {
	    name: string;
	    data_type: string;
	    table_oid?: number;
	
	    static createFrom(source: any = {}) {
	        return new ColumnMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.data_type = source["data_type"];
	        this.table_oid = source["table_oid"];
	    }
	}
	export class Connection {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    database: string;
	    username: string;
	    password: string;
	    ssl_mode: string;
	    color: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Connection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.database = source["database"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.ssl_mode = source["ssl_mode"];
	        this.color = source["color"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConnectionTestResult {
	    ok: boolean;
	    latency_ms: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionTestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.latency_ms = source["latency_ms"];
	        this.message = source["message"];
	    }
	}
	export class ExplainNode {
	    node_type: string;
	    relation_name?: string;
	    alias?: string;
	    startup_cost: number;
	    total_cost: number;
	    plan_rows: number;
	    plan_width: number;
	    actual_time?: number;
	    actual_rows?: number;
	    actual_loops?: number;
	    filter?: string;
	    rows_removed_by_filter?: number;
	    shared_hit_blocks?: number;
	    shared_read_blocks?: number;
	    plans?: ExplainNode[];
	
	    static createFrom(source: any = {}) {
	        return new ExplainNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.node_type = source["node_type"];
	        this.relation_name = source["relation_name"];
	        this.alias = source["alias"];
	        this.startup_cost = source["startup_cost"];
	        this.total_cost = source["total_cost"];
	        this.plan_rows = source["plan_rows"];
	        this.plan_width = source["plan_width"];
	        this.actual_time = source["actual_time"];
	        this.actual_rows = source["actual_rows"];
	        this.actual_loops = source["actual_loops"];
	        this.filter = source["filter"];
	        this.rows_removed_by_filter = source["rows_removed_by_filter"];
	        this.shared_hit_blocks = source["shared_hit_blocks"];
	        this.shared_read_blocks = source["shared_read_blocks"];
	        this.plans = this.convertValues(source["plans"], ExplainNode);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExplainResult {
	    raw_text: string;
	    plan?: ExplainNode;
	
	    static createFrom(source: any = {}) {
	        return new ExplainResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.raw_text = source["raw_text"];
	        this.plan = this.convertValues(source["plan"], ExplainNode);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FKDefinition {
	    name: string;
	    source_columns: string[];
	    target_schema: string;
	    target_table: string;
	    target_columns: string[];
	    on_update: string;
	    on_delete: string;
	
	    static createFrom(source: any = {}) {
	        return new FKDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.source_columns = source["source_columns"];
	        this.target_schema = source["target_schema"];
	        this.target_table = source["target_table"];
	        this.target_columns = source["target_columns"];
	        this.on_update = source["on_update"];
	        this.on_delete = source["on_delete"];
	    }
	}
	export class FKInfo {
	    constraint_name: string;
	    source_schema: string;
	    source_table: string;
	    source_column: string;
	    target_schema: string;
	    target_table: string;
	    target_column: string;
	    on_update: string;
	    on_delete: string;
	
	    static createFrom(source: any = {}) {
	        return new FKInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.constraint_name = source["constraint_name"];
	        this.source_schema = source["source_schema"];
	        this.source_table = source["source_table"];
	        this.source_column = source["source_column"];
	        this.target_schema = source["target_schema"];
	        this.target_table = source["target_table"];
	        this.target_column = source["target_column"];
	        this.on_update = source["on_update"];
	        this.on_delete = source["on_delete"];
	    }
	}
	
	
	export class FunctionInfo {
	    schema: string;
	    name: string;
	    arguments: string;
	    return_type: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new FunctionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.name = source["name"];
	        this.arguments = source["arguments"];
	        this.return_type = source["return_type"];
	        this.language = source["language"];
	    }
	}
	export class IndexDefinition {
	    name: string;
	    table_name: string;
	    columns: string[];
	    is_unique: boolean;
	    index_type: string;
	    where?: string;
	
	    static createFrom(source: any = {}) {
	        return new IndexDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.table_name = source["table_name"];
	        this.columns = source["columns"];
	        this.is_unique = source["is_unique"];
	        this.index_type = source["index_type"];
	        this.where = source["where"];
	    }
	}
	export class IndexInfo {
	    index_name: string;
	    table_name: string;
	    columns: string[];
	    is_unique: boolean;
	    index_type: string;
	    definition: string;
	
	    static createFrom(source: any = {}) {
	        return new IndexInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index_name = source["index_name"];
	        this.table_name = source["table_name"];
	        this.columns = source["columns"];
	        this.is_unique = source["is_unique"];
	        this.index_type = source["index_type"];
	        this.definition = source["definition"];
	    }
	}
	export class QueryHistoryEntry {
	    id: number;
	    connection_id: string;
	    database: string;
	    sql: string;
	    duration_ms: number;
	    rows_returned: number;
	    // Go type: time
	    executed_at: any;
	    error?: string;
	    bookmarked: boolean;
	
	    static createFrom(source: any = {}) {
	        return new QueryHistoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.connection_id = source["connection_id"];
	        this.database = source["database"];
	        this.sql = source["sql"];
	        this.duration_ms = source["duration_ms"];
	        this.rows_returned = source["rows_returned"];
	        this.executed_at = this.convertValues(source["executed_at"], null);
	        this.error = source["error"];
	        this.bookmarked = source["bookmarked"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QueryResult {
	    columns: ColumnMeta[];
	    rows: any[][];
	    row_count: number;
	    total_count: number;
	    duration_ms: number;
	    error?: string;
	    query_type: string;
	    rows_affected?: number;
	    raw_sql?: string;
	
	    static createFrom(source: any = {}) {
	        return new QueryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.columns = this.convertValues(source["columns"], ColumnMeta);
	        this.rows = source["rows"];
	        this.row_count = source["row_count"];
	        this.total_count = source["total_count"];
	        this.duration_ms = source["duration_ms"];
	        this.error = source["error"];
	        this.query_type = source["query_type"];
	        this.rows_affected = source["rows_affected"];
	        this.raw_sql = source["raw_sql"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SavedQuery {
	    id: number;
	    name: string;
	    sql: string;
	    folder?: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SavedQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.sql = source["sql"];
	        this.folder = source["folder"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SchemaInfo {
	    name: string;
	    owner: string;
	
	    static createFrom(source: any = {}) {
	        return new SchemaInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.owner = source["owner"];
	    }
	}
	export class SequenceInfo {
	    schema: string;
	    name: string;
	    data_type: string;
	    min_value: number;
	    max_value: number;
	    increment: number;
	    cache_size: number;
	
	    static createFrom(source: any = {}) {
	        return new SequenceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.name = source["name"];
	        this.data_type = source["data_type"];
	        this.min_value = source["min_value"];
	        this.max_value = source["max_value"];
	        this.increment = source["increment"];
	        this.cache_size = source["cache_size"];
	    }
	}
	export class TableDefinition {
	    schema: string;
	    table_name: string;
	    columns: ColumnDefinition[];
	    indexes?: IndexDefinition[];
	    foreign_keys?: FKDefinition[];
	
	    static createFrom(source: any = {}) {
	        return new TableDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.table_name = source["table_name"];
	        this.columns = this.convertValues(source["columns"], ColumnDefinition);
	        this.indexes = this.convertValues(source["indexes"], IndexDefinition);
	        this.foreign_keys = this.convertValues(source["foreign_keys"], FKDefinition);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TableInfo {
	    schema: string;
	    name: string;
	    type: string;
	    row_count: number;
	    comment?: string;
	
	    static createFrom(source: any = {}) {
	        return new TableInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.row_count = source["row_count"];
	        this.comment = source["comment"];
	    }
	}
	export class TriggerInfo {
	    schema: string;
	    table_name: string;
	    trigger_name: string;
	    event: string;
	    timing: string;
	    definition: string;
	
	    static createFrom(source: any = {}) {
	        return new TriggerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.table_name = source["table_name"];
	        this.trigger_name = source["trigger_name"];
	        this.event = source["event"];
	        this.timing = source["timing"];
	        this.definition = source["definition"];
	    }
	}
	export class TypeInfo {
	    schema: string;
	    name: string;
	    type: string;
	    values?: string[];
	
	    static createFrom(source: any = {}) {
	        return new TypeInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.values = source["values"];
	    }
	}
	export class ViewInfo {
	    schema: string;
	    name: string;
	    read_only: boolean;
	    definition: string;
	
	    static createFrom(source: any = {}) {
	        return new ViewInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.name = source["name"];
	        this.read_only = source["read_only"];
	        this.definition = source["definition"];
	    }
	}

}

