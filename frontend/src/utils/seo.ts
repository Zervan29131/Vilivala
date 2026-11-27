import { siteConfig } from '@/site.config'

export interface SeoMetaOptions {
  title?: string
  description?: string
  path?: string
  type?: 'website' | 'article'
  publishedAt?: string
  updatedAt?: string
  tags?: string[]
}

function upsertMeta(attr: 'name' | 'property', key: string, content: string) {
  let el = document.head.querySelector(`meta[${attr}="${key}"]`) as HTMLMetaElement | null
  if (!el) {
    el = document.createElement('meta')
    el.setAttribute(attr, key)
    document.head.appendChild(el)
  }
  el.setAttribute('content', content)
}

function upsertCanonical(href: string) {
  let el = document.head.querySelector('link[rel="canonical"]') as HTMLLinkElement | null
  if (!el) {
    el = document.createElement('link')
    el.setAttribute('rel', 'canonical')
    document.head.appendChild(el)
  }
  el.setAttribute('href', href)
}

function upsertJsonLd(data: Record<string, unknown> | null) {
  const id = 'vilivala-jsonld'
  let el = document.getElementById(id) as HTMLScriptElement | null
  if (!data) {
    el?.remove()
    return
  }
  if (!el) {
    el = document.createElement('script')
    el.id = id
    el.type = 'application/ld+json'
    document.head.appendChild(el)
  }
  el.textContent = JSON.stringify(data)
}

export function updatePageMeta(opts: SeoMetaOptions = {}) {
  const fullTitle = opts.title ? `${opts.title} - ${siteConfig.title}` : siteConfig.title
  const desc = (opts.description || `${siteConfig.hero.title} · ${siteConfig.description}`).slice(0, 160)
  const origin = window.location.origin
  const url = `${origin}${opts.path ?? window.location.pathname}`
  const type = opts.type || 'website'

  document.title = fullTitle
  upsertMeta('name', 'description', desc)
  upsertCanonical(url)

  upsertMeta('property', 'og:type', type)
  upsertMeta('property', 'og:site_name', siteConfig.title)
  upsertMeta('property', 'og:title', fullTitle)
  upsertMeta('property', 'og:description', desc)
  upsertMeta('property', 'og:url', url)

  upsertMeta('name', 'twitter:card', 'summary')
  upsertMeta('name', 'twitter:title', fullTitle)
  upsertMeta('name', 'twitter:description', desc)

  if (type === 'article') {
    upsertJsonLd({
      '@context': 'https://schema.org',
      '@type': 'BlogPosting',
      headline: opts.title || siteConfig.title,
      description: desc,
      url,
      datePublished: opts.publishedAt,
      dateModified: opts.updatedAt || opts.publishedAt,
      keywords: opts.tags?.join(', '),
      author: { '@type': 'Person', name: siteConfig.blogger.name },
    })
  } else {
    upsertJsonLd(null)
  }
}
