const S = { width: 16, height: 16, fill: 'none' as const, 'aria-hidden': true }

export function IconRuns() {
  return (
    <svg {...S} viewBox="0 0 16 16">
      <circle cx="3" cy="8" r="1.5" stroke="currentColor" strokeWidth="1.2" />
      <circle cx="8" cy="8" r="1.5" stroke="currentColor" strokeWidth="1.2" />
      <circle cx="13" cy="8" r="1.5" stroke="currentColor" strokeWidth="1.2" />
      <path d="M4.5 8h2M9.5 8h2" stroke="currentColor" strokeWidth="1.2" />
    </svg>
  )
}

export function IconAnalytics() {
  return (
    <svg {...S} viewBox="0 0 16 16">
      <rect x="1" y="9" width="3" height="6" rx="1" fill="currentColor" opacity="0.45" />
      <rect x="6" y="5" width="3" height="10" rx="1" fill="currentColor" opacity="0.7" />
      <rect x="11" y="1" width="3" height="14" rx="1" fill="currentColor" />
    </svg>
  )
}

export function IconSettings() {
  return (
    <svg {...S} viewBox="0 0 16 16">
      <circle cx="8" cy="8" r="2.2" stroke="currentColor" strokeWidth="1.2" />
      <path d="M8 1.5v1.3M8 13.2v1.3M1.5 8h1.3M13.2 8h1.3M3.5 3.5l.9.9M11.6 11.6l.9.9M3.5 12.5l.9-.9M11.6 4.4l.9-.9" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
    </svg>
  )
}

export function IconBranch() {
  return (
    <svg width={12} height={12} viewBox="0 0 16 16" fill="none" aria-hidden>
      <circle cx="4" cy="4" r="2" stroke="currentColor" strokeWidth="1.3" />
      <circle cx="4" cy="12" r="2" stroke="currentColor" strokeWidth="1.3" />
      <circle cx="12" cy="4" r="2" stroke="currentColor" strokeWidth="1.3" />
      <path d="M4 6v4M4 6c0 2 8 3 8-2" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
}

export function IconExternal() {
  return (
    <svg width={12} height={12} viewBox="0 0 16 16" fill="none" aria-hidden>
      <path d="M6 3h7v7M13 3L3 13" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function IconChevron() {
  return (
    <svg width={14} height={14} viewBox="0 0 14 14" fill="none" aria-hidden>
      <path d="M3 7h8M7.5 3.5L11 7l-3.5 3.5" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function IconSetup() {
  return (
    <svg {...S} viewBox="0 0 16 16">
      <path d="M3 4h10v8H3z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" />
      <path d="M6 8h4M6 10.5h2.5" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
    </svg>
  )
}

export function IconPipeline() {
  return (
    <svg {...S} viewBox="0 0 16 16">
      <path d="M4 2.5h5l2.5 2.5V13.5H4V2.5z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" />
      <path d="M9 2.5V5h2.5" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" />
    </svg>
  )
}

export function IconGeneral() {
  return <IconSettings />
}

export function IconEnv() {
  return (
    <svg {...S} viewBox="0 0 16 16">
      <path d="M3 4l3 4-3 4M8 12h5" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function IconSecrets() {
  return (
    <svg {...S} viewBox="0 0 16 16">
      <rect x="4" y="7" width="8" height="7" rx="1.5" stroke="currentColor" strokeWidth="1.2" />
      <path d="M6 7V5a2 2 0 114 0v2" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
    </svg>
  )
}

export function IconTeam() {
  return (
    <svg {...S} viewBox="0 0 16 16">
      <circle cx="6" cy="5" r="2.5" stroke="currentColor" strokeWidth="1.2" />
      <path d="M1 13c0-2.5 2.2-4 5-4s5 1.5 5 4" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
      <circle cx="12" cy="5" r="2" stroke="currentColor" strokeWidth="1.2" />
    </svg>
  )
}
