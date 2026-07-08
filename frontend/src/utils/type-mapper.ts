export interface PGType {
  label: string
  value: string
  category: string
}

export const PG_TYPES: PGType[] = [
  // Numeric
  { label: 'smallint', value: 'smallint', category: 'Numeric' },
  { label: 'integer', value: 'integer', category: 'Numeric' },
  { label: 'bigint', value: 'bigint', category: 'Numeric' },
  { label: 'decimal', value: 'decimal', category: 'Numeric' },
  { label: 'numeric', value: 'numeric', category: 'Numeric' },
  { label: 'real', value: 'real', category: 'Numeric' },
  { label: 'double precision', value: 'double precision', category: 'Numeric' },
  { label: 'smallserial', value: 'smallserial', category: 'Numeric' },
  { label: 'serial', value: 'serial', category: 'Numeric' },
  { label: 'bigserial', value: 'bigserial', category: 'Numeric' },

  // Text
  { label: 'varchar(n)', value: 'varchar', category: 'Text' },
  { label: 'char(n)', value: 'char', category: 'Text' },
  { label: 'text', value: 'text', category: 'Text' },
  { label: 'citext', value: 'citext', category: 'Text' },

  // Binary
  { label: 'bytea', value: 'bytea', category: 'Binary' },

  // Boolean
  { label: 'boolean', value: 'boolean', category: 'Boolean' },

  // Date/Time
  { label: 'date', value: 'date', category: 'Date/Time' },
  { label: 'time', value: 'time', category: 'Date/Time' },
  { label: 'timetz', value: 'timetz', category: 'Date/Time' },
  { label: 'timestamp', value: 'timestamp', category: 'Date/Time' },
  { label: 'timestamptz', value: 'timestamptz', category: 'Date/Time' },
  { label: 'interval', value: 'interval', category: 'Date/Time' },

  // JSON
  { label: 'json', value: 'json', category: 'JSON' },
  { label: 'jsonb', value: 'jsonb', category: 'JSON' },

  // UUID
  { label: 'uuid', value: 'uuid', category: 'UUID' },

  // Network
  { label: 'inet', value: 'inet', category: 'Network' },
  { label: 'cidr', value: 'cidr', category: 'Network' },
  { label: 'macaddr', value: 'macaddr', category: 'Network' },

  // Geometric
  { label: 'point', value: 'point', category: 'Geometric' },
  { label: 'line', value: 'line', category: 'Geometric' },
  { label: 'box', value: 'box', category: 'Geometric' },
  { label: 'polygon', value: 'polygon', category: 'Geometric' },
  { label: 'circle', value: 'circle', category: 'Geometric' },

  // Special
  { label: 'tsvector', value: 'tsvector', category: 'Special' },
  { label: 'tsquery', value: 'tsquery', category: 'Special' },
  { label: 'xml', value: 'xml', category: 'Special' },
]

export const PG_TYPE_CATEGORIES = [...new Set(PG_TYPES.map(t => t.category))]

export function getTypesByCategory(category: string): PGType[] {
  return PG_TYPES.filter(t => t.category === category)
}

export function formatPGType(type: string, length?: number): string {
  if (length && (type === 'varchar' || type === 'char')) {
    return `${type}(${length})`
  }
  return type
}

export function inferPGType(value: any): string {
  if (value === null || value === undefined) return 'text'
  if (typeof value === 'boolean') return 'boolean'
  if (typeof value === 'number') {
    if (Number.isInteger(value)) {
      if (Math.abs(value) < 32768) return 'smallint'
      if (Math.abs(value) < 2147483647) return 'integer'
      return 'bigint'
    }
    return 'double precision'
  }
  if (typeof value === 'string') {
    // Try to detect UUID
    if (/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(value)) {
      return 'uuid'
    }
    // Try to detect date
    if (/^\d{4}-\d{2}-\d{2}/.test(value)) {
      if (value.includes('T') || value.includes(' ')) return 'timestamptz'
      return 'date'
    }
    // Try to detect JSON
    if ((value.startsWith('{') && value.endsWith('}')) || (value.startsWith('[') && value.endsWith(']'))) {
      try {
        JSON.parse(value)
        return 'jsonb'
      } catch {
        // not JSON
      }
    }
    return 'text'
  }
  return 'text'
}
