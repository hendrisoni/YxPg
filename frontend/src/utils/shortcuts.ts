import type { Ref } from 'vue'

export interface Shortcut {
  key: string
  ctrl?: boolean
  shift?: boolean
  alt?: boolean
  action: string
  description: string
}

export const SHORTCUTS: Shortcut[] = [
  { key: 'q', ctrl: true, action: 'newQueryTab', description: 'New Query Tab' },
  { key: 't', ctrl: true, action: 'newQueryTab', description: 'New Query Tab' },
  { key: 'w', ctrl: true, action: 'closeTab', description: 'Close Tab' },
  { key: 'Enter', ctrl: true, action: 'runQuery', description: 'Run Query' },
  { key: 'F5', action: 'runQuery', description: 'Run Query' },
  { key: 'F', ctrl: true, shift: true, action: 'formatSql', description: 'Format SQL' },
  { key: 'H', ctrl: true, action: 'toggleHistory', description: 'Toggle History' },
  { key: 'P', ctrl: true, action: 'commandPalette', description: 'Command Palette' },

  { key: 'R', ctrl: true, action: 'refreshSchema', description: 'Refresh Schema' },
  { key: 'N', ctrl: true, action: 'newConnection', description: 'New Connection' },
  { key: 'k', ctrl: true, action: 'toggleTableSearch', description: 'Search Tables/Views/Functions' },
  { key: '/', ctrl: true, action: 'toggleComment', description: 'Toggle Comment' },
  { key: 'F1', action: 'openBuilder', description: 'Query Builder' },
  { key: 'F2', action: 'openDesigner', description: 'Table Designer' },
  { key: 'F3', action: 'findNext', description: 'Find Next' },
]

export function formatShortcut(shortcut: Shortcut): string {
  const parts: string[] = []
  if (shortcut.ctrl) parts.push('Ctrl')
  if (shortcut.shift) parts.push('Shift')
  if (shortcut.alt) parts.push('Alt')
  parts.push(shortcut.key)
  return parts.join('+')
}

export function setupKeyboardShortcuts(handlers: Record<string, () => void>) {
  const handler = (e: KeyboardEvent) => {
    for (const shortcut of SHORTCUTS) {
      const ctrlMatch = shortcut.ctrl ? (e.ctrlKey || e.metaKey) : !(e.ctrlKey || e.metaKey)
      const shiftMatch = shortcut.shift ? e.shiftKey : !e.shiftKey
      const altMatch = shortcut.alt ? e.altKey : !e.altKey
      const keyMatch = e.key.toLowerCase() === shortcut.key.toLowerCase() ||
                       e.code === shortcut.key

      if (ctrlMatch && shiftMatch && altMatch && keyMatch) {
        e.preventDefault()
        const actionHandler = handlers[shortcut.action]
        if (actionHandler) {
          actionHandler()
        }
        return
      }
    }
  }

  window.addEventListener('keydown', handler)
  return () => window.removeEventListener('keydown', handler)
}
