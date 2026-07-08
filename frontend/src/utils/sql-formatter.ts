import { format } from 'sql-formatter'

export function formatSQL(sql: string, dialect: string = 'postgresql'): string {
  try {
    return format(sql, {
      language: 'postgresql',
      tabWidth: 2,
      keywordCase: 'upper',
      identifierCase: 'lower',
      linesBetweenQueries: 1,
    })
  } catch {
    return sql
  }
}

export function truncateSQL(sql: string, maxLen: number = 100): string {
  const oneLine = sql.replace(/\s+/g, ' ').trim()
  if (oneLine.length <= maxLen) return oneLine
  return oneLine.substring(0, maxLen) + '...'
}

export function detectQueryType(sql: string): string {
  const trimmed = sql.trim().toUpperCase()
  if (trimmed.startsWith('SELECT') || trimmed.startsWith('WITH')) return 'select'
  if (trimmed.startsWith('INSERT')) return 'insert'
  if (trimmed.startsWith('UPDATE')) return 'update'
  if (trimmed.startsWith('DELETE')) return 'delete'
  if (trimmed.startsWith('CREATE') || trimmed.startsWith('ALTER') || trimmed.startsWith('DROP')) return 'ddl'
  return 'other'
}
