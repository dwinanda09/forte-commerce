'use client'

import { useState } from 'react'
import type { Condition, Action, ConditionType, ActionType } from '@/lib/types'

interface CampaignFormProps {
  initialValues?: {
    name: string
    description: string
    is_active: boolean
    priority: number
    conditions: Condition[]
    actions: Action[]
  }
  onSubmit: (data: {
    name: string
    description: string
    is_active: boolean
    priority: number
    conditions: Condition[]
    actions: Action[]
  }) => void
  loading: boolean
  submitLabel: string
}

const CONDITION_TYPES: { value: ConditionType; label: string }[] = [
  { value: 'cart_has_sku', label: 'Cart has SKU' },
  { value: 'item_qty_gte', label: 'Item qty ≥' },
  { value: 'cart_total_gte', label: 'Cart total ≥' },
  { value: 'cart_item_count_gte', label: 'Cart item count ≥' },
]

const ACTION_TYPES: { value: ActionType; label: string }[] = [
  { value: 'free_item', label: 'Free item' },
  { value: 'buy_n_get_m', label: 'Buy N get M' },
  { value: 'pct_discount_on_sku', label: '% discount on SKU' },
  { value: 'pct_discount_on_cart', label: '% discount on cart' },
  { value: 'fixed_discount', label: 'Fixed discount' },
]

const defaultCondition = (): Condition => ({ type: 'cart_has_sku' })
const defaultAction = (): Action => ({ type: 'free_item' })

function ConditionRow({
  cond,
  index,
  onChange,
  onRemove,
}: {
  cond: Condition
  index: number
  onChange: (i: number, c: Condition) => void
  onRemove: (i: number) => void
}) {
  const set = (patch: Partial<Condition>) => onChange(index, { ...cond, ...patch })

  return (
    <div className="flex flex-wrap gap-3 items-start p-3 bg-surface rounded-md border border-mist">
      <select
        value={cond.type}
        onChange={e => set({ type: e.target.value as ConditionType, sku: undefined, min_qty: undefined, qty: undefined, amount: undefined, count: undefined })}
        className="border border-mist rounded px-2 py-1.5 text-sm text-graphite bg-white focus:outline-none focus:border-teal"
      >
        {CONDITION_TYPES.map(ct => (
          <option key={ct.value} value={ct.value}>{ct.label}</option>
        ))}
      </select>

      {(cond.type === 'cart_has_sku' || cond.type === 'item_qty_gte') && (
        <input
          placeholder="SKU"
          value={cond.sku ?? ''}
          onChange={e => set({ sku: e.target.value })}
          className="border border-mist rounded px-2 py-1.5 text-sm w-32 focus:outline-none focus:border-teal"
        />
      )}

      {cond.type === 'item_qty_gte' && (
        <input
          type="number"
          placeholder="Min qty"
          value={cond.min_qty ?? ''}
          onChange={e => set({ min_qty: Number(e.target.value) })}
          className="border border-mist rounded px-2 py-1.5 text-sm w-24 focus:outline-none focus:border-teal"
        />
      )}

      {cond.type === 'cart_total_gte' && (
        <input
          type="number"
          step="0.01"
          placeholder="Amount"
          value={cond.amount ?? ''}
          onChange={e => set({ amount: Number(e.target.value) })}
          className="border border-mist rounded px-2 py-1.5 text-sm w-28 focus:outline-none focus:border-teal"
        />
      )}

      {cond.type === 'cart_item_count_gte' && (
        <input
          type="number"
          placeholder="Count"
          value={cond.count ?? ''}
          onChange={e => set({ count: Number(e.target.value) })}
          className="border border-mist rounded px-2 py-1.5 text-sm w-24 focus:outline-none focus:border-teal"
        />
      )}

      <button
        type="button"
        onClick={() => onRemove(index)}
        className="ml-auto text-red-500 hover:text-red-700 text-sm font-medium"
      >
        Remove
      </button>
    </div>
  )
}

function ActionRow({
  action,
  index,
  onChange,
  onRemove,
}: {
  action: Action
  index: number
  onChange: (i: number, a: Action) => void
  onRemove: (i: number) => void
}) {
  const set = (patch: Partial<Action>) => onChange(index, { ...action, ...patch })

  return (
    <div className="flex flex-wrap gap-3 items-start p-3 bg-surface rounded-md border border-mist">
      <select
        value={action.type}
        onChange={e => set({ type: e.target.value as ActionType, sku: undefined, trigger_sku: undefined, buy_n: undefined, pay_m: undefined, pct: undefined, amount: undefined })}
        className="border border-mist rounded px-2 py-1.5 text-sm text-graphite bg-white focus:outline-none focus:border-teal"
      >
        {ACTION_TYPES.map(at => (
          <option key={at.value} value={at.value}>{at.label}</option>
        ))}
      </select>

      {action.type === 'free_item' && (
        <>
          <input
            placeholder="Trigger SKU"
            value={action.trigger_sku ?? ''}
            onChange={e => set({ trigger_sku: e.target.value })}
            className="border border-mist rounded px-2 py-1.5 text-sm w-32 focus:outline-none focus:border-teal"
          />
          <input
            placeholder="Free SKU"
            value={action.sku ?? ''}
            onChange={e => set({ sku: e.target.value })}
            className="border border-mist rounded px-2 py-1.5 text-sm w-32 focus:outline-none focus:border-teal"
          />
        </>
      )}

      {action.type === 'buy_n_get_m' && (
        <>
          <input
            placeholder="SKU"
            value={action.sku ?? ''}
            onChange={e => set({ sku: e.target.value })}
            className="border border-mist rounded px-2 py-1.5 text-sm w-32 focus:outline-none focus:border-teal"
          />
          <input
            type="number"
            placeholder="Buy N"
            value={action.buy_n ?? ''}
            onChange={e => set({ buy_n: Number(e.target.value) })}
            className="border border-mist rounded px-2 py-1.5 text-sm w-20 focus:outline-none focus:border-teal"
          />
          <input
            type="number"
            placeholder="Pay M"
            value={action.pay_m ?? ''}
            onChange={e => set({ pay_m: Number(e.target.value) })}
            className="border border-mist rounded px-2 py-1.5 text-sm w-20 focus:outline-none focus:border-teal"
          />
        </>
      )}

      {(action.type === 'pct_discount_on_sku') && (
        <>
          <input
            placeholder="SKU"
            value={action.sku ?? ''}
            onChange={e => set({ sku: e.target.value })}
            className="border border-mist rounded px-2 py-1.5 text-sm w-32 focus:outline-none focus:border-teal"
          />
          <input
            type="number"
            step="0.01"
            placeholder="Pct %"
            value={action.pct ?? ''}
            onChange={e => set({ pct: Number(e.target.value) })}
            className="border border-mist rounded px-2 py-1.5 text-sm w-24 focus:outline-none focus:border-teal"
          />
        </>
      )}

      {action.type === 'pct_discount_on_cart' && (
        <input
          type="number"
          step="0.01"
          placeholder="Pct %"
          value={action.pct ?? ''}
          onChange={e => set({ pct: Number(e.target.value) })}
          className="border border-mist rounded px-2 py-1.5 text-sm w-24 focus:outline-none focus:border-teal"
        />
      )}

      {action.type === 'fixed_discount' && (
        <input
          type="number"
          step="0.01"
          placeholder="Amount"
          value={action.amount ?? ''}
          onChange={e => set({ amount: Number(e.target.value) })}
          className="border border-mist rounded px-2 py-1.5 text-sm w-28 focus:outline-none focus:border-teal"
        />
      )}

      <button
        type="button"
        onClick={() => onRemove(index)}
        className="ml-auto text-red-500 hover:text-red-700 text-sm font-medium"
      >
        Remove
      </button>
    </div>
  )
}

export function CampaignForm({ initialValues, onSubmit, loading, submitLabel }: CampaignFormProps) {
  const [name, setName] = useState(initialValues?.name ?? '')
  const [description, setDescription] = useState(initialValues?.description ?? '')
  const [isActive, setIsActive] = useState(initialValues?.is_active ?? true)
  const [priority, setPriority] = useState(initialValues?.priority ?? 0)
  const [conditions, setConditions] = useState<Condition[]>(initialValues?.conditions ?? [])
  const [actions, setActions] = useState<Action[]>(initialValues?.actions ?? [])

  const updateCondition = (i: number, c: Condition) =>
    setConditions(prev => prev.map((x, idx) => (idx === i ? c : x)))

  const removeCondition = (i: number) =>
    setConditions(prev => prev.filter((_, idx) => idx !== i))

  const updateAction = (i: number, a: Action) =>
    setActions(prev => prev.map((x, idx) => (idx === i ? a : x)))

  const removeAction = (i: number) =>
    setActions(prev => prev.filter((_, idx) => idx !== i))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({ name, description, is_active: isActive, priority, conditions, actions })
  }

  const inputClass = 'w-full border border-mist rounded-md px-3 py-2 text-sm text-graphite bg-white focus:outline-none focus:ring-2 focus:ring-teal/30 focus:border-teal transition-colors'
  const labelClass = 'block text-sm font-medium text-graphite mb-1'

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="bg-white rounded-lg shadow-card p-6 space-y-5">
        <div>
          <label className={labelClass}>Name <span className="text-red-500">*</span></label>
          <input
            required
            value={name}
            onChange={e => setName(e.target.value)}
            placeholder="e.g. Summer Sale"
            className={inputClass}
          />
        </div>

        <div>
          <label className={labelClass}>Description</label>
          <textarea
            value={description}
            onChange={e => setDescription(e.target.value)}
            placeholder="Optional description"
            rows={2}
            className={`${inputClass} resize-none`}
          />
        </div>

        <div className="flex gap-6">
          <div className="flex-1">
            <label className={labelClass}>Priority</label>
            <input
              type="number"
              value={priority}
              onChange={e => setPriority(Number(e.target.value))}
              className={inputClass}
            />
          </div>

          <div className="flex items-center gap-3 pt-5">
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={isActive}
                onChange={e => setIsActive(e.target.checked)}
                className="sr-only peer"
              />
              <div className="w-10 h-6 bg-mist peer-checked:bg-teal rounded-full transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:after:translate-x-4" />
            </label>
            <span className="text-sm font-medium text-graphite">Active</span>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-card p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-sm font-semibold text-graphite uppercase tracking-wide">Conditions</h2>
          <button
            type="button"
            onClick={() => setConditions(prev => [...prev, defaultCondition()])}
            className="text-sm text-teal hover:text-teal-hover font-medium transition-colors"
          >
            + Add condition
          </button>
        </div>
        {conditions.length === 0 && (
          <p className="text-sm text-steel">No conditions — campaign applies to every cart.</p>
        )}
        <div className="space-y-2">
          {conditions.map((c, i) => (
            <ConditionRow key={i} cond={c} index={i} onChange={updateCondition} onRemove={removeCondition} />
          ))}
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-card p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-sm font-semibold text-graphite uppercase tracking-wide">Actions</h2>
          <button
            type="button"
            onClick={() => setActions(prev => [...prev, defaultAction()])}
            className="text-sm text-teal hover:text-teal-hover font-medium transition-colors"
          >
            + Add action
          </button>
        </div>
        {actions.length === 0 && (
          <p className="text-sm text-steel">No actions — add at least one to apply discounts.</p>
        )}
        <div className="space-y-2">
          {actions.map((a, i) => (
            <ActionRow key={i} action={a} index={i} onChange={updateAction} onRemove={removeAction} />
          ))}
        </div>
      </div>

      <div className="flex justify-end">
        <button
          type="submit"
          disabled={loading}
          className="bg-teal text-white px-6 py-2.5 rounded-md font-medium hover:bg-teal-hover transition-colors text-sm disabled:opacity-50"
        >
          {loading ? 'Saving…' : submitLabel}
        </button>
      </div>
    </form>
  )
}
