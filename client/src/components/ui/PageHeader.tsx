import { pageHeader, pageTitle, pageDesc } from '../../styles/theme'

export default function PageHeader({
  title,
  description,
}: {
  title: string
  description?: string
}) {
  return (
    <header css={pageHeader}>
      <h1 css={pageTitle}>{title}</h1>
      {description && <p css={pageDesc}>{description}</p>}
    </header>
  )
}
