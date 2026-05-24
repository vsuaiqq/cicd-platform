const S = { width: 20, height: 20, fill: 'none' as const, 'aria-hidden': true }

export function IconSearch() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <circle cx="9" cy="9" r="5.5" stroke="currentColor" strokeWidth="1.25" />
      <path d="M13 13l4 4" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconDashboard() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <rect x="2.5" y="2.5" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="11.5" y="2.5" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="2.5" y="11.5" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="11.5" y="11.5" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
    </svg>
  )
}

export function IconPlus() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <path d="M10 4.5v11M4.5 10h11" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconProject() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <path d="M5 3.5h7l3.5 3.5V16.5H5V3.5z" stroke="currentColor" strokeWidth="1.25" strokeLinejoin="round" />
      <path d="M12 3.5V7h3.5" stroke="currentColor" strokeWidth="1.25" strokeLinejoin="round" />
    </svg>
  )
}

export function IconClock() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <circle cx="10" cy="10" r="6.5" stroke="currentColor" strokeWidth="1.25" />
      <path d="M10 6.5V10l2.5 2" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconUser() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <circle cx="10" cy="7" r="3" stroke="currentColor" strokeWidth="1.25" />
      <path d="M4.5 16.5c0-3 2.5-5 5.5-5s5.5 2 5.5 5" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconSettings() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <circle cx="10" cy="10" r="2.5" stroke="currentColor" strokeWidth="1.25" />
      <path d="M10 4v1.2M10 14.8V16M4 10h1.2M15.8 10H17M6 6l.9.9M13.1 13.1l.9.9M6 14l.9-.9M13.1 6.9l.9-.9" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconEnter() {
  return (
    <svg width={14} height={14} viewBox="0 0 16 16" fill="none" aria-hidden>
      <path d="M3 8h7M8 5l3 3-3 3" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}
