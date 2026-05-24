
export function matchScore(haystack: string, query: string): number {
  const text = haystack.toLowerCase()
  const q = query.trim().toLowerCase()
  if (!q) return 1
  if (text === q) return 100
  if (text.startsWith(q)) return 85
  if (text.includes(q)) return 55

  let from = 0
  for (const ch of q) {
    const at = text.indexOf(ch, from)
    if (at === -1) return 0
    from = at + 1
  }
  return 25
}

export function bestScore(fields: string[], query: string): number {
  return Math.max(0, ...fields.map((f) => matchScore(f, query)))
}
