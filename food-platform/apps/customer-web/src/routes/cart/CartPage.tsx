// Cart page — review cart items before checkout

import { useState } from 'react'
import { useNavigate, Link, useLocation } from 'react-router-dom'
import { AppLayout } from '../../components/AppLayout'
import { Button, Badge } from '@food-platform/ui'
import { formatEGP } from '@food-platform/utils'
import { useAuthStore } from '@food-platform/auth'

interface CartItem {
  id: string
  name: string
  image_url?: string
  quantity: number
  unit_price: number
  modifiers: { modifier_name: string; option_name: string }[]
  notes?: string
  line_total: number
}

// Mock cart data (will be replaced with real cart store)
const MOCK_CART: CartItem[] = [
  {
    id: '1',
    name: 'Margherita Pizza (وسط)',
    quantity: 1,
    unit_price: 145,
    modifiers: [
      { modifier_name: 'إضافات', option_name: 'جبن إضافي' },
      { modifier_name: 'إضافات', option_name: 'فطر' },
    ],
    notes: 'بدون بصل',
    line_total: 170,
  },
  {
    id: '2',
    name: 'Pepperoni Pizza (كبير)',
    quantity: 2,
    unit_price: 165,
    modifiers: [],
    line_total: 330,
  },
]

const SERVICE_FEE_RATE = 0.05
const VAT_RATE = 0.14
const DELIVERY_FEE = 25
const DISCOUNT = 0

export function CartPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const [cart, setCart] = useState(MOCK_CART)
  const [couponCode, setCouponCode] = useState('')
  const [couponApplied, setCouponApplied] = useState(false)

  const subtotal = cart.reduce((sum, item) => sum + item.line_total, 0)
  const serviceFee = subtotal * SERVICE_FEE_RATE
  const vat = (subtotal + DELIVERY_FEE + serviceFee - DISCOUNT) * VAT_RATE
  const total = subtotal + DELIVERY_FEE + serviceFee + vat - DISCOUNT

  const updateQuantity = (id: string, delta: number) => {
    setCart((prev) =>
      prev
        .map((item) =>
          item.id === id
            ? {
                ...item,
                quantity: Math.max(1, item.quantity + delta),
                line_total: Math.max(1, item.quantity + delta) * item.unit_price,
              }
            : item,
        )
        .filter((item) => item.quantity > 0),
    )
  }

  const removeItem = (id: string) => {
    setCart((prev) => prev.filter((item) => item.id !== id))
  }

  const applyCoupon = () => {
    if (couponCode.toUpperCase() === 'WELCOME50') {
      setCouponApplied(true)
    }
  }

  if (cart.length === 0) {
    return (
      <AppLayout>
        <div className="min-h-[60vh] flex items-center justify-center p-6">
          <div className="text-center">
            <span className="material-symbols-rounded text-text-tertiary text-6xl mb-4">
              shopping_cart
            </span>
            <h2 className="text-h3 font-bold text-text-primary mb-2">السلة فاضية</h2>
            <p className="text-body text-text-secondary mb-6">
              ابدأ تطلب من مطاعمك المفضلة
            </p>
            <Link to="/">
              <Button size="lg">تصفح المطاعم</Button>
            </Link>
          </div>
        </div>
      </AppLayout>
    )
  }

  return (
    <AppLayout>
      <div className="max-w-2xl mx-auto px-4 py-4 pb-28">
        {/* Header */}
        <h1 className="text-h2 font-bold text-text-primary mb-4">السلة</h1>

        {/* Address */}
        <div className="bg-surface rounded-lg border border-border p-4 mb-4">
          <div className="flex items-center gap-2 text-body-sm">
            <span className="material-symbols-rounded text-primary">location_on</span>
            <div>
              <p className="text-text-tertiary text-overline">التوصيل إلى</p>
              <p className="text-text-primary font-medium">الزمالك، 26 يوليو، شقة 5</p>
            </div>
            <button className="text-primary text-caption font-semibold mr-auto hover:underline">
              تغيير
            </button>
          </div>
        </div>

        {/* Restaurant + items */}
        <div className="bg-surface rounded-lg border border-border p-4 mb-4">
          <h2 className="text-body font-bold text-text-primary mb-3">Pizza Hut</h2>

          <div className="space-y-3">
            {cart.map((item) => (
              <div key={item.id} className="flex gap-3 pb-3 border-b border-border last:border-0 last:pb-0">
                {/* Image */}
                <div className="w-16 h-16 rounded-lg overflow-hidden bg-bg-tertiary flex-shrink-0">
                  {item.image_url ? (
                    <img src={item.image_url} alt={item.name} className="w-full h-full object-cover" />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center">
                      <span className="material-symbols-rounded text-text-tertiary">lunch_dining</span>
                    </div>
                  )}
                </div>

                {/* Details */}
                <div className="flex-1 min-w-0">
                  <h3 className="text-body-sm font-semibold text-text-primary">{item.name}</h3>
                  {item.modifiers.map((mod, i) => (
                    <p key={i} className="text-caption text-text-secondary">
                      + {mod.option_name}
                    </p>
                  ))}
                  {item.notes && (
                    <p className="text-caption text-text-tertiary italic">{item.notes}</p>
                  )}

                  <div className="flex items-center justify-between mt-2">
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => updateQuantity(item.id, -1)}
                        className="w-7 h-7 rounded-full border border-border flex items-center justify-center hover:border-primary hover:text-primary"
                      >
                        <span className="material-symbols-rounded text-sm">remove</span>
                      </button>
                      <span className="text-body font-semibold w-6 text-center">{item.quantity}</span>
                      <button
                        onClick={() => updateQuantity(item.id, 1)}
                        className="w-7 h-7 rounded-full border border-border flex items-center justify-center hover:border-primary hover:text-primary"
                      >
                        <span className="material-symbols-rounded text-sm">add</span>
                      </button>
                    </div>

                    <div className="flex items-center gap-3">
                      <span className="text-body font-bold text-text-primary">
                        {formatEGP(item.line_total)}
                      </span>
                      <button
                        onClick={() => removeItem(item.id)}
                        className="text-text-tertiary hover:text-error"
                        aria-label="حذف"
                      >
                        <span className="material-symbols-rounded text-sm">delete</span>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Coupon */}
        <div className="bg-surface rounded-lg border border-border p-4 mb-4">
          <div className="flex items-center gap-2 text-body-sm mb-2">
            <span className="material-symbols-rounded text-primary">local_offer</span>
            <span className="font-semibold text-text-primary">كوبون خصم</span>
          </div>
          <div className="flex gap-2">
            <input
              type="text"
              value={couponCode}
              onChange={(e) => setCouponCode(e.target.value)}
              placeholder="WELCOME50"
              className="flex-1 h-10 px-3 rounded-md border border-border bg-surface text-body placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
              disabled={couponApplied}
            />
            <Button
              size="md"
              variant={couponApplied ? 'secondary' : 'primary'}
              onClick={applyCoupon}
              disabled={couponApplied || !couponCode}
            >
              {couponApplied ? '✅ مطبق' : 'تطبيق'}
            </Button>
          </div>
          {couponApplied && (
            <p className="text-caption text-success mt-2">
              وفّرت EGP 50 على طلبك!
            </p>
          )}
        </div>

        {/* Price breakdown */}
        <div className="bg-surface rounded-lg border border-border p-4 mb-4">
          <h3 className="text-body font-semibold text-text-primary mb-3">تفاصيل الفاتورة</h3>
          <div className="space-y-2 text-body-sm">
            <PriceRow label="Subtotal (الأصناف)" value={formatEGP(subtotal)} />
            <PriceRow label="رسوم التوصيل" value={formatEGP(DELIVERY_FEE)} />
            <PriceRow label={`رسوم الخدمة (${SERVICE_FEE_RATE * 100}%)`} value={formatEGP(serviceFee)} />
            <PriceRow label={`ض.ق.م (${VAT_RATE * 100}%)`} value={formatEGP(vat)} />
            {DISCOUNT > 0 && (
              <PriceRow label="الخصم" value={`-${formatEGP(DISCOUNT)}`} variant="success" />
            )}
            <div className="border-t border-border pt-2 mt-2">
              <div className="flex items-center justify-between">
                <span className="text-body font-bold text-text-primary">الإجمالي</span>
                <span className="text-h3 font-bold text-primary">{formatEGP(total)}</span>
              </div>
            </div>
          </div>

          {/* Cashback */}
          <div className="mt-3 bg-success/5 border border-success/20 rounded-lg p-3 flex items-center gap-2">
            <span className="material-symbols-rounded text-success">savings</span>
            <p className="text-caption text-success">
              كاش باك: {formatEGP(total * 0.05)} (5%) — هيتضاف لمحفظتك بعد التوصيل
            </p>
          </div>
        </div>

        {/* Spend more prompt */}
        {total < 700 && (
          <div className="bg-warning/5 border border-warning/20 rounded-lg p-3 mb-4 text-center">
            <p className="text-body-sm text-warning font-medium">
              🎁 تستحق خصم EGP 30 لو طلبت بـ EGP 700
            </p>
          </div>
        )}

        {/* Checkout button (sticky) */}
        <div className="fixed bottom-16 md:bottom-0 left-0 right-0 bg-surface border-t border-border p-4 z-sticky">
          <div className="max-w-2xl mx-auto">
            <Button
              fullWidth
              size="lg"
              onClick={() => {
                // Guest users must login before checkout
                if (!isAuthenticated) {
                  navigate('/login', { state: { from: '/checkout', message: 'سجّل دخولك عشان تكمل الطلب' } })
                } else {
                  navigate('/checkout')
                }
              }}
            >
              🛒 اطلب الآن — {formatEGP(total)}
            </Button>
          </div>
        </div>
      </div>
    </AppLayout>
  )
}

function PriceRow({
  label,
  value,
  variant = 'default',
}: {
  label: string
  value: string
  variant?: 'default' | 'success'
}) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-text-secondary">{label}</span>
      <span className={variant === 'success' ? 'text-success font-medium' : 'text-text-primary font-medium'}>
        {value}
      </span>
    </div>
  )
}
