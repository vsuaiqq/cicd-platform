import { HighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { oneDark } from '@codemirror/theme-one-dark'
import { EditorView } from '@codemirror/view'
import { tags as t } from '@lezer/highlight'
import type { Extension } from '@codemirror/state'

const flowEditorLight = EditorView.theme(
  {
    '&': {
      color: 'var(--text-primary)',
      backgroundColor: 'var(--editor-bg)',
    },
    '.cm-content': { caretColor: 'var(--accent)' },
    '.cm-cursor, .cm-dropCursor': { borderLeftColor: 'var(--accent)' },
    '&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection': {
      backgroundColor: 'var(--accent-muted) !important',
    },
    '.cm-gutters': {
      backgroundColor: 'var(--editor-bg)',
      color: 'var(--text-disabled)',
      border: 'none',
    },
    '.cm-activeLineGutter': { backgroundColor: 'var(--editor-line-active)' },
    '.cm-activeLine': { backgroundColor: 'var(--editor-line-active)' },
    '.cm-lintRange-error': {
      backgroundImage: 'none',
      borderBottom: '2px wavy var(--danger)',
    },
    '.cm-lintRange-warning': {
      backgroundImage: 'none',
      borderBottom: '2px wavy var(--warning)',
    },
  },
  { dark: false },
)

const flowYamlLight = syntaxHighlighting(
  HighlightStyle.define([
    { tag: t.keyword, color: 'var(--syntax-key)' },
    { tag: [t.name, t.propertyName, t.definition(t.name)], color: 'var(--syntax-key)' },
    { tag: t.comment, color: 'var(--syntax-comment)', fontStyle: 'italic' },
    { tag: [t.string, t.special(t.string)], color: 'var(--syntax-string)' },
    { tag: [t.number, t.bool, t.null], color: 'var(--syntax-number)' },
    { tag: t.operator, color: 'var(--syntax-punctuation)' },
    { tag: t.punctuation, color: 'var(--syntax-punctuation)' },
    { tag: t.meta, color: 'var(--syntax-value)' },
  ]),
)


export function getCodeMirrorTheme(isDark: boolean): Extension {
  return isDark ? oneDark : [flowEditorLight, flowYamlLight]
}
