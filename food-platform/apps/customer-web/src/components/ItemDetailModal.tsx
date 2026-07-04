// Item detail modal — customize item before adding to cart

import { useState } from 'react'
import type { MenuItem } from '@food-platform/types'
import { Button, Badge } from '@food-platform/ui'
import { formatEGP } from '@food-platform/utils'

interface ItemDetailModalProps {
  item: MenuItem
  onClose: () => void
  onAddToCart?: (item: MenuItem, quantity: number, selectedModifiers: Record<string, string[]>, notes: string) => void
}

export function ItemDetailModal({ item, onClose, onAddToCart }: ItemDetailModalProps) {
  const [quantity, setQuantity] = useState(1)
  const [notes, setNotes] = useState('')
  const [selectedModifiers, setSelectedModifiers] = useState<Record<string, string[]>>({})

  // Calculate total price
  const basePrice = item.price * quantity
  const modifiersTotal = calculateModifiersTotal(item, selectedModifiers) * quantity
  const total = basePrice + modifiersTotal

  const handleModifierToggle = (modifierId: string, optionId: string, required: boolean, multipleChoice: boolean) => {
    setSelectedModifiers((prev) => {
      const current = prev[modifierId] ?? []

      if (!multipleChoice) {
        // Single choice — replace
        return { ...prev, [modifierId]: [optionId] }
      }

      // Multiple choice — toggle
      if (current.includes(optionId)) {
        return { ...prev, [modifierId]: current.filter((id) => id !== optionId) }
      }
      return { ...prev, [modifierId]: [...current, optionId] }
    })
  }

  const handleAddToCart = () => {
    // Validate required modifiers
    for (const modifier of item.modifiers ?? []) {
      if (modifier.required && !(selectedModifiers[modifier.id]?.length)) {
        return // Don't add — required modifier not selected
      }
    }

    onAddToCart?.(item, quantity, selectedModifiers, notes)
    onClose()
  }

  const allRequiredSelected = (item.modifiers ?? []).every(
    (m) => !m.required || (selectedModifiers[m.id]?.length ?? 0) > 0,
  )

  return (
    <div className="fixed inset-0 z-modal flex items-end md:items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50 animate-fade-in"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative bg-surface rounded-t-2xl md:rounded-2xl w-full md:max-w-lg max-h-[90vh] overflow-y-auto animate-slide-up">
        {/* Drag handle */}
        <div className="sticky top-0 bg-surface pt-3 pb-2 z-10">
          <div className="w-10 h-1 bg-border rounded-full mx-auto" />
        </div>

        {/* Item image */}
        {item.image_url && (
          <div className="w-full h-48 bg-bg-tertiary">
            <img
              src={item.image_url}
              alt={item.name}
              className="w-full h-full object-cover"
            />
          </div>
        )}

        {/* Content */}
        <div className="p-5">
          {/* Name + rating */}
          <div className="flex items-start justify-between gap-2">
            <h2 className="text-h2 font-bold text-text-primary">{item.name}</h2>
            {item.is_most_ordered && (
              <Badge variant="warning">🔥 الأكثر طلباً</Badge>
            )}
          </div>

          {/* Rating */}
          {item.rating && (
            <div className="flex items-center gap-1 mt-1 text-caption text-text-secondary">
              <span className="text-warning">★</span>
              <span className="font-semibold">{item.rating.toFixed(1)}</span>
              <span className="text-text-tertiary">({item.rating_count ?? 0} تقييم)</span>
            </div>
          )}

          {/* Description */}
          {item.description && (
            <p className="text-body-sm text-text-secondary mt-2">{item.description}</p>
          )}

          {/* Price */}
          <p className="text-h3 font-bold text-primary mt-3">{formatEGP(item.price)}</p>

          {/* Modifiers */}
          {item.modifiers?.map((modifier) => (
            <div key={modifier.id} className="mt-5">
              <div className="flex items-center gap-2 mb-2">
                <h3 className="text-body font-semibold text-text-primary">{modifier.name}</h3>
                {modifier.required && (
                  <Badge variant="error">إجباري</Badge>
                )}
              </div>

              <div className="space-y-2">
                {modifier.options.map((option) => {
                  const isSelected = selectedModifiers[modifier.id]?.includes(option.id)
                  return (
                    <label
                      key={option.id}
                      className={`flex items-center justify-between p-3 rounded-lg border-2 cursor-pointer transition-colors
                        ${isSelected
                          ? 'border-primary bg-primary/5'
                          : 'border-border hover:border-border-strong'
                        }`}
                    >
                      <div className="flex items-center gap-3">
                        <input
                          type={modifier.multiple_choice ? 'checkbox' : 'radio'}
                          name={`modifier-${modifier.id}`}
                          checked={isSelected}
                          onChange={() =>
                            handleModifierToggle(
                              modifier.id,
                              option.id,
                              modifier.required,
                              modifier.multiple_choice,
                            )
                          }
                          className="w-5 h-5 accent-primary"
                        />
                        <span className="text-body-sm text-text-primary">{option.name}</span>
                      </div>
                      {option.price_delta > 0 && (
                        <span className="text-body-sm text-text-secondary">
                          +{formatEGP(option.price_delta)}
                        </span>
                      )}
                    </label>
                  )
                })}
              </div>
            </div>
          ))}

          {/* Notes */}
          <div className="mt-5">
            <label className="block text-body font-semibold text-text-primary mb-2">
              ملاحظات خاصة
            </label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="مثلاً: بدون بصل، صلصة إضافية..."
              maxLength={200}
              rows={2}
              className="w-full p-3 rounded-lg border border-border bg-surface text-body placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary resize-none"
            />
            <p className="text-overline text-text-tertiary mt-1">{notes.length}/200</p>
          </div>

          {/* Quantity + Add to cart */}
          <div className="flex items-center gap-4 mt-6 pt-4 border-t border-border">
            {/* Quantity */}
            <div className="flex items-center gap-3">
              <button
                onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                className="w-10 h-10 rounded-full border-2 border-border flex items-center justify-center hover:border-primary hover:text-primary transition-colors"
                disabled={quantity <= 1}
              >
                <span className="material-symbols-rounded">remove</span>
              </button>
              <span className="text-h3 font-bold w-8 text-center">{quantity}</span>
              <button
                onClick={() => setQuantity((q) => Math.min(20, q + 1))}
                className="w-10 h-10 rounded-full border-2 border-border flex items-center justify-center hover:border-primary hover:text-primary transition-colors"
                disabled={quantity >= 20}
              >
                <span className="material-symbols-rounded">add</span>
              </button>
            </div>

            {/* Add to cart */}
            <Button
              fullWidth
              size="lg"
              onClick={handleAddToCart}
              disabled={!allRequiredSelected}
            >
              {allRequiredSelected ? (
                <>🛒 أضف للسلة — {formatEGP(total)}</>
              ) : (
                'اختر الخيارات المطلوبة'
              )}
            </Button>
          </div>
        </div>

        {/* Close button */}
        <button
          onClick={onClose}
          className="absolute top-4 left-4 w-8 h-8 bg-surface/80 backdrop-blur rounded-full flex items-center justify-center hover:bg-bg-tertiary transition-colors"
          aria-label="إغلاق"
        >
          <span className="material-symbols-rounded text-text-secondary">close</span>
        </button>
      </div>
    </div>
  )
}

// ============ Helpers ============

function calculateModifiersTotal(
  item: MenuItem,
  selected: Record<string, string[]>,
): number {
  let total = 0
  for (const modifier of item.modifiers ?? []) {
    const selectedOptions = selected[modifier.id] ?? []
    for (const option of modifier.options) {
      if (selectedOptions.includes(option.id)) {
        total += option.price_delta
      }
    }
  }
  return total
}
