// Order tracking page — real-time order status with live map

import { useState, useEffect } from 'react'
import { useParams, useLocation, useNavigate, Link } from 'react-router-dom'
import { AppLayout } from '../../components/AppLayout'
import { Button, Badge } from '@food-platform/ui'
import { formatEGP } from '@food-platform/utils'
import { useWebSocket } from '@food-platform/hooks'

// Order status steps
const ORDER_STEPS = [
  { status: 'pending', label: 'تم استلام الطلب', icon: 'receipt' },
  { status: 'confirmed', label: 'المطعم بدأ التحضير', icon: 'restaurant' },
  { status: 'preparing', label: 'الأكل جاهز', icon: 'check_circle' },
  { status: 'ready', label: 'المندوب في الطريق', icon: 'delivery_dining' },
  { status: 'picked_up', label: 'في الطريق إليك', icon: 'location_on' },
  { status: 'delivered', label: 'تم التوصيل', icon: 'celebration' },
] as const

// Mock order data (will come from API)
const MOCK_ORDER = {
  id: 'mock-order-id',
  orderNumber: 'A7X92F',
  restaurantName: 'Pizza Hut',
  total: 636.6,
  paymentMethod: 'vodafone_cash',
  eta: '8:35 PM',
  items: [
    { name: 'Margherita Pizza (وسط)', quantity: 1, price: 170 },
    { name: 'Pepperoni Pizza (كبير)', quantity: 2, price: 330 },
    { name: 'Coca Cola', quantity: 2, price: 50 },
  ],
}

export function OrderTrackingPage() {
  const { id } = useParams<{ id: string }>()
  const location = useLocation()
  const navigate = useNavigate()

  const state = (location.state ?? {}) as {
    orderConfirmed?: boolean
    total?: number
    eta?: string
  }

  const [currentStep, setCurrentStep] = useState(0)
  const [showCelebration, setShowCelebration] = useState(state.orderConfirmed ?? false)

  // Simulate order progression (in production: WebSocket events from backend)
  useEffect(() => {
    if (!showCelebration) return

    // Auto-advance through steps for demo
    const timers: ReturnType<typeof setTimeout>[] = []
    const delays = [0, 3000, 8000, 12000, 16000, 20000]

    for (let i = 0; i < ORDER_STEPS.length; i++) {
      const timer = setTimeout(() => {
        setCurrentStep(i)
        if (i === ORDER_STEPS.length - 1) {
          // Order delivered
          setShowCelebration(false)
        }
      }, delays[i])
      timers.push(timer)
    }

    return () => timers.forEach(clearTimeout)
  }, [showCelebration])

  // WebSocket connection (will be used in production for real-time updates)
  // const { isConnected } = useWebSocket(`order.${id}`, {
  //   onMessage: (msg) => {
  //     if (msg.event === 'order.status_changed') {
  //       const newStatus = msg.payload.status
  //       const stepIndex = ORDER_STEPS.findIndex(s => s.status === newStatus)
  //       if (stepIndex >= 0) setCurrentStep(stepIndex)
  //     }
  //   }
  // })

  return (
    <AppLayout showHeader={false} showBottomNav={false}>
      <div className="min-h-screen pb-24">
        {/* Header */}
        <div className="sticky top-0 z-sticky bg-surface border-b border-border h-14 flex items-center px-4">
          <button
            onClick={() => navigate('/')}
            className="flex items-center gap-1 text-text-secondary hover:text-text-primary"
          >
            <span className="material-symbols-rounded">arrow_forward</span>
            <span className="text-body">رجوع</span>
          </button>
          <h1 className="text-body font-bold text-text-primary absolute left-1/2 -translate-x-1/2">
            الطلب #{MOCK_ORDER.orderNumber}
          </h1>
        </div>

        {/* Order confirmation celebration */}
        {showCelebration && (
          <div className="bg-success/5 border-b border-success/20 p-4 text-center">
            <div className="inline-flex items-center justify-center w-16 h-16 bg-success rounded-full mb-2">
              <span className="material-symbols-rounded text-white text-3xl">check</span>
            </div>
            <h2 className="text-h3 font-bold text-text-primary">تم الطلب! 🎉</h2>
            <p className="text-body-sm text-text-secondary mt-1">
              طلبك اتأكد والمطعم هيبدأ تحضيره
            </p>
          </div>
        )}

        {/* ETA banner */}
        <div className="bg-primary/5 border-b border-primary/20 p-4">
          <div className="flex items-center justify-between max-w-2xl mx-auto">
            <div className="flex items-center gap-2">
              <span className="material-symbols-rounded text-primary">schedule</span>
              <div>
                <p className="text-overline text-text-tertiary">الوقت المتوقع للتوصيل</p>
                <p className="text-h3 font-bold text-primary">{MOCK_ORDER.eta}</p>
              </div>
            </div>
            <div className="text-left">
              <p className="text-overline text-text-tertiary">الإجمالي</p>
              <p className="text-h3 font-bold text-text-primary">{formatEGP(MOCK_ORDER.total)}</p>
            </div>
          </div>
        </div>

        <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
          {/* Status timeline */}
          <div className="bg-surface rounded-lg border border-border p-5">
            <h2 className="text-body font-bold text-text-primary mb-4">حالة الطلب</h2>

            <div className="space-y-4">
              {ORDER_STEPS.map((step, index) => {
                const isCompleted = index < currentStep
                const isCurrent = index === currentStep
                const isFuture = index > currentStep

                return (
                  <div key={step.status} className="flex items-center gap-3">
                    {/* Icon */}
                    <div
                      className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 transition-all
                        ${isCompleted
                          ? 'bg-success text-white'
                          : isCurrent
                            ? 'bg-primary text-white animate-pulse'
                            : 'bg-bg-tertiary text-text-tertiary'
                        }`}
                    >
                      <span className="material-symbols-rounded text-xl">
                        {isCompleted ? 'check' : step.icon}
                      </span>
                    </div>

                    {/* Line (connecting) */}
                    {index < ORDER_STEPS.length - 1 && (
                      <div
                        className={`absolute mr-5 mt-10 w-0.5 h-6 transition-colors
                          ${isCompleted ? 'bg-success' : 'bg-border'}`}
                      />
                    )}

                    {/* Label */}
                    <div>
                      <p
                        className={`text-body-sm font-medium transition-colors
                          ${isCompleted || isCurrent ? 'text-text-primary' : 'text-text-tertiary'}`}
                      >
                        {step.label}
                      </p>
                      {isCurrent && (
                        <p className="text-caption text-primary mt-0.5">جاري...</p>
                      )}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>

          {/* Live map placeholder */}
          {currentStep >= 3 && (
            <div className="bg-surface rounded-lg border border-border overflow-hidden">
              <div className="h-48 bg-bg-tertiary flex items-center justify-center relative">
                <div className="text-center">
                  <span className="material-symbols-rounded text-primary text-5xl mb-2">
                    delivery_dining
                  </span>
                  <p className="text-body-sm text-text-secondary">المندوب في الطريق إليك</p>
                  <p className="text-caption text-text-tertiary mt-1">ETA: 10 دقايق</p>
                </div>

                {/* Animated marker */}
                <div className="absolute bottom-4 right-4 animate-bounce">
                  <span className="material-symbols-rounded text-primary text-3xl">location_on</span>
                </div>
              </div>
            </div>
          )}

          {/* Driver info (when assigned) */}
          {currentStep >= 3 && currentStep < 5 && (
            <div className="bg-surface rounded-lg border border-border p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                  <span className="material-symbols-rounded text-primary">person</span>
                </div>
                <div className="flex-1">
                  <p className="text-body font-semibold text-text-primary">Mahmoud S.</p>
                  <div className="flex items-center gap-1 text-caption text-text-secondary">
                    <span className="text-warning">★</span>
                    <span>4.8</span>
                    <span className="text-text-tertiary">• موتوسيكل</span>
                  </div>
                </div>
                <a
                  href="tel:+201000000000"
                  className="w-10 h-10 rounded-full bg-success/10 flex items-center justify-center hover:bg-success/20 transition-colors"
                >
                  <span className="material-symbols-rounded text-success">call</span>
                </a>
                <button className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center hover:bg-primary/20 transition-colors">
                  <span className="material-symbols-rounded text-primary">chat</span>
                </button>
              </div>
            </div>
          )}

          {/* Order details */}
          <div className="bg-surface rounded-lg border border-border p-4">
            <h2 className="text-body font-bold text-text-primary mb-3">تفاصيل الطلب</h2>
            <p className="text-body-sm text-text-secondary mb-3">{MOCK_ORDER.restaurantName}</p>

            <div className="space-y-2">
              {MOCK_ORDER.items.map((item, i) => (
                <div key={i} className="flex items-center justify-between text-body-sm">
                  <span className="text-text-primary">
                    {item.quantity}× {item.name}
                  </span>
                  <span className="text-text-secondary">{formatEGP(item.price)}</span>
                </div>
              ))}
            </div>

            <div className="border-t border-border mt-3 pt-3 flex items-center justify-between">
              <span className="text-body font-bold text-text-primary">الإجمالي</span>
              <span className="text-body font-bold text-primary">{formatEGP(MOCK_ORDER.total)}</span>
            </div>

            <div className="mt-2 flex items-center gap-2">
              <Badge variant="success">{MOCK_ORDER.paymentMethod === 'vodafone_cash' ? '💚 Vodafone Cash' : 'مدفوع'}</Badge>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-3">
            {currentStep < 2 && (
              <Button
                variant="outline"
                fullWidth
                onClick={() => navigate('/')}
              >
                🚫 إلغاء الطلب
              </Button>
            )}
            <Button
              variant="ghost"
              fullWidth
              onClick={() => navigate('/')}
            >
              ❓ مساعدة
            </Button>
          </div>

          {/* Rate order (when delivered) */}
          {currentStep === 5 && (
            <div className="bg-surface rounded-lg border border-border p-4 text-center">
              <h3 className="text-body font-bold text-text-primary mb-2">قيّم طلبك 👇</h3>
              <div className="flex justify-center gap-2 my-3">
                {[1, 2, 3, 4, 5].map((star) => (
                  <button
                    key={star}
                    className="text-3xl text-warning hover:scale-110 transition-transform"
                  >
                    ★
                  </button>
                ))}
              </div>
              <Link to="/">
                <Button fullWidth>العودة للرئيسية</Button>
              </Link>
            </div>
          )}
        </div>
      </div>
    </AppLayout>
  )
}
