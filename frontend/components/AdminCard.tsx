import Link from 'next/link'

interface StatItem {
  key: string
  value: React.ReactNode
}

interface AdminCardProps {
  label?: string
  title: string
  badge?: React.ReactNode
  description?: string
  stats?: StatItem[]
  editHref: string
  onDelete: () => void
  isDeleting?: boolean
}

export function AdminCard({
  label,
  title,
  badge,
  description,
  stats,
  editHref,
  onDelete,
  isDeleting,
}: AdminCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-card border border-mist/70 p-5 flex flex-col gap-3 hover:shadow-modal hover:border-mist transition-all">
      {/* Header */}
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0">
          {label && (
            <p className="text-xs font-mono text-steel mb-0.5 uppercase tracking-wide truncate">
              {label}
            </p>
          )}
          <h3 className="text-sm font-semibold text-graphite leading-snug line-clamp-2">
            {title}
          </h3>
        </div>
        {badge && <div className="shrink-0">{badge}</div>}
      </div>

      {/* Description */}
      {description !== undefined && (
        <p className="text-xs text-steel line-clamp-2 min-h-[2rem]">
          {description || <span className="italic opacity-60">No description</span>}
        </p>
      )}

      {/* Stats */}
      {stats && stats.length > 0 && (
        <div className="flex flex-wrap gap-x-4 gap-y-1">
          {stats.map((s) => (
            <span key={s.key} className="text-xs text-steel">
              {s.key}{' '}
              <span className="font-semibold text-graphite tabular-nums">{s.value}</span>
            </span>
          ))}
        </div>
      )}

      {/* Footer */}
      <div className="flex items-center justify-end gap-3 pt-1 border-t border-mist/50 mt-auto">
        <Link
          href={editHref}
          className="text-xs text-teal hover:text-teal-hover font-medium transition-colors"
        >
          Edit
        </Link>
        <button
          onClick={onDelete}
          disabled={isDeleting}
          className="text-xs text-steel hover:text-red-600 transition-colors disabled:opacity-40"
        >
          {isDeleting ? '…' : 'Delete'}
        </button>
      </div>
    </div>
  )
}
