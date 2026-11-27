export interface ParsedCommand {
  raw: string
  cmd: string
  args: string[]
  rest: string
  flags: {
    self: boolean
    recursive: boolean
    all: boolean
    from?: string
  }
}

/** Expand `!!` (previous command) or `!n` (1-based history index). */
export function expandHistoryBang(raw: string, history: string[]): { expanded: string; error?: string } {
  const trimmed = raw.trim()
  if (trimmed === '!!') {
    if (!history.length) return { expanded: '', error: '暂无可重放的历史命令' }
    return { expanded: history[history.length - 1] }
  }
  const m = trimmed.match(/^!(\d+)$/)
  if (m) {
    const idx = Number(m[1]) - 1
    if (idx < 0 || idx >= history.length) {
      return { expanded: '', error: `未找到历史命令 #${m[1]}（共 ${history.length} 条）` }
    }
    return { expanded: history[idx] }
  }
  return { expanded: trimmed }
}

/** Tokenize a CLI line supporting single/double quotes and common flags (-s, -r, -a, -f <src>). */
export function parseCommandLine(raw: string): ParsedCommand {
  const trimmed = raw.trim()
  const tokens: string[] = []
  let cur = ''
  let quote: '"' | "'" | null = null

  for (let i = 0; i < trimmed.length; i++) {
    const ch = trimmed[i]
    if (quote) {
      if (ch === quote) {
        quote = null
      } else {
        cur += ch
      }
      continue
    }
    if (ch === '"' || ch === "'") {
      quote = ch
      continue
    }
    if (/\s/.test(ch)) {
      if (cur) {
        tokens.push(cur)
        cur = ''
      }
      continue
    }
    cur += ch
  }
  if (cur) tokens.push(cur)

  const cmd = (tokens[0] || '').toLowerCase()
  const args: string[] = []
  const flags: ParsedCommand['flags'] = {
    self: false,
    recursive: false,
    all: false,
  }

  for (let i = 1; i < tokens.length; i++) {
    const t = tokens[i]
    if (t === '-s' || t === '--self') {
      flags.self = true
    } else if (t === '-r' || t === '--recursive' || t === '-rf') {
      flags.recursive = true
    } else if (t === '-a' || t === '--all') {
      flags.all = true
    } else if (t === '-f' || t === '--from') {
      if (tokens[i + 1]) {
        flags.from = tokens[++i].toLowerCase()
      }
    } else if (t.startsWith('--from=')) {
      flags.from = t.slice(7).toLowerCase()
    } else {
      args.push(t)
    }
  }

  return {
    raw: trimmed,
    cmd,
    args,
    rest: args.join(' ').trim(),
    flags,
  }
}

/** Parse duration strings like `90`, `90s`, `25m`, `1h` into seconds. */
export function parseDurationSeconds(input: string): number | null {
  const s = (input || '').trim().toLowerCase()
  if (!s) return null
  const m = s.match(/^(\d+(?:\.\d+)?)\s*(s|sec|m|min|h)?$/)
  if (!m) return null
  const val = parseFloat(m[1])
  const unit = m[2] || 's'
  const mult = unit.startsWith('h') ? 3600 : unit.startsWith('m') ? 60 : 1
  const total = Math.round(val * mult)
  if (!Number.isFinite(total) || total <= 0 || total > 86400) return null
  return total
}
