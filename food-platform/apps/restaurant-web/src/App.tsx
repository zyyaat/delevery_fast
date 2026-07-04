// Restaurant Web App — Full implementation as Vite React app
// This replaces the placeholder App.tsx with a working restaurant dashboard

import { useState, useEffect } from 'react'
import { cn } from '@food-platform/ui'

type Screen = 'dashboard' | 'orders' | 'menu' | 'kds' | 'analytics' | 'schedule'

export default function App() {
  const [screen, setScreen] = useState<Screen>('dashboard')
  const [isOnline, setIsOnline] = useState(true)
  const [incomingOrder, setIncomingOrder] = useState<typeof mockIncomingOrder | null>(null)
  const [orders, setOrders] = useState(mockOrders)

  // Simulate incoming order every 30 seconds
  useEffect(() => {
    const timer = setInterval(() => {
      if (!incomingOrder && isOnline) {
        setIncomingOrder({ ...mockIncomingOrder })
      }
    }, 30000)
    return () => clearInterval(timer)
  }, [incomingOrder, isOnline])

  const acceptOrder = () => {
    if (incomingOrder) {
      setOrders([{
        id: incomingOrder.id,
        number: incomingOrder.number,
        status: 'preparing',
        timer: 0,
        items: incomingOrder.items,
        total: incomingOrder.total,
        payment: incomingOrder.payment,
        customer: incomingOrder.customer,
      }, ...orders])
      setIncomingOrder(null)
    }
  }

  const rejectOrder = () => {
    setIncomingOrder(null)
  }

  return (
    <div className="min-h-screen bg-bg-primary flex" dir="rtl">
      {/* Sidebar */}
      <aside className="w-64 bg-surface border-l border-border flex flex-col h-screen sticky top-0">
        {/* Logo */}
        <div className="p-5 border-b border-border">
          <div className="flex items-center gap-2">
            <span className="text-2xl">🍔</span>
            <div>
              <h1 className="text-body font-bold text-text-primary">Pizza Hut</h1>
              <p className="text-caption text-text-tertiary">معادي</p>
            </div>
          </div>
        </div>

        {/* Status toggle */}
        <div className="p-4 border-b border-border">
          <button
            onClick={() => setIsOnline(!isOnline)}
            className={cn(
              'w-full flex items-center justify-between p-3 rounded-lg border-2 transition-colors',
              isOnline
                ? 'border-success bg-success/5'
                : 'border-error bg-error/5'
            )}
          >
            <span className={cn('text-body-sm font-semibold', isOnline ? 'text-success' : 'text-error')}>
              {isOnline ? '🟢 مفتوح' : '🔴 مقفل'}
            </span>
            <span className={cn('material-symbols-rounded', isOnline ? 'text-success' : 'text-error')}>
              {isOnline ? 'toggle_on' : 'toggle_off'}
            </span>
          </button>
        </div>

        {/* Nav */}
        <nav className="flex-1 p-3 space-y-1">
          {navItems.map((item) => (
            <button
              key={item.id}
              onClick={() => setScreen(item.id as Screen)}
              className={cn(
                'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-body-sm font-medium transition-colors',
                screen === item.id
                  ? 'bg-primary text-white'
                  : 'text-text-secondary hover:bg-bg-tertiary'
              )}
            >
              <span className="material-symbols-rounded text-xl">{item.icon}</span>
              {item.label}
            </button>
          ))}
        </nav>

        {/* Settings */}
        <div className="p-3 border-t border-border">
          <button className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-body-sm text-text-secondary hover:bg-bg-tertiary">
            <span className="material-symbols-rounded text-xl">settings</span>
            الإعدادات
          </button>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 overflow-auto">
        {screen === 'dashboard' && <DashboardScreen orders={orders} isOnline={isOnline} onNavigate={setScreen} />}
        {screen === 'orders' && <OrdersScreen orders={orders} setOrders={setOrders} />}
        {screen === 'menu' && <MenuScreen />}
        {screen === 'kds' && <KDSScreen orders={orders} setOrders={setOrders} />}
        {screen === 'analytics' && <AnalyticsScreen />}
        {screen === 'schedule' && <ScheduleScreen />}
      </main>

      {/* Inbound order modal */}
      {incomingOrder && (
        <InboundOrderModal order={incomingOrder} onAccept={acceptOrder} onReject={rejectOrder} />
      )}
    </div>
  )
}

// ============ Navigation ============
const navItems = [
  { id: 'dashboard', label: 'لوحة التحكم', icon: 'dashboard' },
  { id: 'orders', label: 'الطلبات', icon: 'receipt_long' },
  { id: 'menu', label: 'المنيو', icon: 'restaurant_menu' },
  { id: 'kds', label: 'شاشة المطبخ', icon: 'display_settings' },
  { id: 'analytics', label: 'التقارير', icon: 'analytics' },
  { id: 'schedule', label: 'المواعيد', icon: 'schedule' },
]

// ============ Mock Data ============
const mockOrders = [
  { id: '1', number: '#A7X92F', status: 'preparing', timer: 482, items: [{ name: 'Margherita Pizza', qty: 1, notes: 'بدون بصل' }, { name: 'Coca Cola', qty: 2 }], total: 290, payment: 'vodafone_cash', customer: 'Ahmed M.' },
  { id: '2', number: '#B3K45L', status: 'ready', timer: 325, items: [{ name: 'Pepperoni Pizza (كبير)', qty: 2 }], total: 330, payment: 'card', customer: 'Fatma A.' },
  { id: '3', number: '#C8M72N', status: 'preparing', timer: 187, items: [{ name: 'Margherita Pizza', qty: 1 }, { name: 'Apple Pie', qty: 2 }], total: 215, payment: 'cod', customer: 'Mahmoud K.' },
]

const mockIncomingOrder = {
  id: 'new-1',
  number: '#D9N81P',
  items: [
    { name: 'Big Mac Meal', qty: 1, notes: 'بدون مخلل' },
    { name: 'McChicken', qty: 2 },
    { name: 'Apple Pie', qty: 1 },
  ],
  total: 285,
  payment: 'vodafone_cash',
  customer: 'Sara M.',
  address: 'الزمالك، 26 يوليو',
  eta: 22,
}

const mockMenu = [
  { id: '1', name: 'Margherita Pizza', price: 145, category: 'بيتزا', available: true, popular: true },
  { id: '2', name: 'Pepperoni Pizza', price: 165, category: 'بيتزا', available: true, popular: true },
  { id: '3', name: 'Veggie Supreme', price: 155, category: 'بيتزا', available: true, popular: false },
  { id: '4', name: 'Coca Cola', price: 25, category: 'مشروبات', available: true, popular: false },
  { id: '5', name: 'Apple Pie', price: 35, category: 'حلويات', available: false, popular: false },
]

// ============ Dashboard Screen ============
function DashboardScreen({ orders, isOnline, onNavigate }: { orders: typeof mockOrders; isOnline: boolean; onNavigate: (s: Screen) => void }) {
  const todayRevenue = 4250
  const todayOrders = 28
  const avgPrepTime = 12
  const rating = 4.6

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-h2 font-bold text-text-primary">لوحة التحكم</h1>

      {/* Stats */}
      <div className="grid grid-cols-4 gap-4">
        <StatCard icon="payments" label="مبيعات اليوم" value={`EGP ${todayRevenue}`} trend="+12%" color="primary" />
        <StatCard icon="receipt_long" label="طلبات اليوم" value={String(todayOrders)} trend="+8%" color="success" />
        <StatCard icon="timer" label="متوسط التحضير" value={`${avgPrepTime} د`} trend="-2%" color="warning" />
        <StatCard icon="star" label="التقييم" value={String(rating)} trend="+0.1" color="accent" />
      </div>

      {/* Active orders */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-h3 font-bold text-text-primary">🔴 طلبات نشطة ({orders.length})</h2>
          <button onClick={() => onNavigate('orders')} className="text-caption text-primary font-semibold hover:underline">
            عرض الكل
          </button>
        </div>
        <div className="grid grid-cols-3 gap-4">
          {orders.map((order) => (
            <OrderCard key={order.id} order={order} />
          ))}
        </div>
      </div>
    </div>
  )
}

// ============ Orders Screen ============
function OrdersScreen({ orders, setOrders }: { orders: typeof mockOrders; setOrders: React.Dispatch<React.SetStateAction<typeof mockOrders>> }) {
  const updateStatus = (id: string, status: string) => {
    setOrders(orders.map(o => o.id === id ? { ...o, status } : o))
  }

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-h2 font-bold text-text-primary">الطلبات النشطة</h1>
      <div className="flex gap-2">
        {['الكل', 'جديد', 'يتحضّر', 'جاهز'].map((tab, i) => (
          <button key={i} className={cn('px-4 py-2 rounded-lg text-body-sm font-medium', i === 0 ? 'bg-primary text-white' : 'bg-surface text-text-secondary border border-border')}>{tab}</button>
        ))}
      </div>
      <div className="space-y-3">
        {orders.map((order) => (
          <OrderDetailCard key={order.id} order={order} onUpdateStatus={updateStatus} />
        ))}
      </div>
    </div>
  )
}

// ============ Menu Screen ============
function MenuScreen() {
  const [menu, setMenu] = useState(mockMenu)

  const toggleAvailability = (id: string) => {
    setMenu(menu.map(m => m.id === id ? { ...m, available: !m.available } : m))
  }

  return (
    <div className="p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-h2 font-bold text-text-primary">🍽️ المنيو</h1>
        <button className="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg text-body-sm font-semibold hover:bg-primary-dark">
          <span className="material-symbols-rounded text-xl">add</span>
          إضافة صنف
        </button>
      </div>

      {['بيتزا', 'مشروبات', 'حلويات'].map((category) => (
        <div key={category}>
          <h2 className="text-body font-bold text-text-primary mb-2">{category} ({menu.filter(m => m.category === category).length})</h2>
          <div className="grid grid-cols-3 gap-3">
            {menu.filter(m => m.category === category).map((item) => (
              <div key={item.id} className={cn('bg-surface rounded-lg border p-4', !item.available && 'opacity-60')}>
                <div className="flex items-start justify-between">
                  <div>
                    <h3 className="text-body-sm font-semibold text-text-primary">{item.name}</h3>
                    <p className="text-caption text-text-secondary mt-1">EGP {item.price}</p>
                  </div>
                  {item.popular && <span className="text-xs bg-warning/10 text-warning px-2 py-0.5 rounded">🔥</span>}
                </div>
                <button
                  onClick={() => toggleAvailability(item.id)}
                  className={cn('w-full mt-3 py-1.5 rounded-lg text-caption font-semibold transition-colors',
                    item.available ? 'bg-success/10 text-success hover:bg-success/20' : 'bg-error/10 text-error hover:bg-error/20'
                  )}
                >
                  {item.available ? '🟢 متاح' : '🔴 نفد (86)'}
                </button>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}

// ============ KDS Screen ============
function KDSScreen({ orders, setOrders }: { orders: typeof mockOrders; setOrders: React.Dispatch<React.SetStateAction<typeof mockOrders>> }) {
  const advanceOrder = (id: string) => {
    const order = orders.find(o => o.id === id)
    if (!order) return
    const next = order.status === 'preparing' ? 'ready' : 'picked_up'
    if (next === 'picked_up') {
      setOrders(orders.filter(o => o.id !== id))
    } else {
      setOrders(orders.map(o => o.id === id ? { ...o, status: next } : o))
    }
  }

  const getColor = (timer: number) => {
    if (timer > 600) return 'border-error'
    if (timer > 300) return 'border-warning'
    return 'border-success'
  }

  return (
    <div className="p-4">
      <h1 className="text-h3 font-bold text-text-primary mb-4">🍳 شاشة المطبخ (KDS)</h1>
      <div className="grid grid-cols-4 gap-3">
        {orders.filter(o => o.status === 'preparing' || o.status === 'ready').map((order) => (
          <div key={order.id} className={cn('bg-surface rounded-lg border-2 p-3', getColor(order.timer), order.status === 'ready' && 'border-success bg-success/5')}>
            <div className="flex items-center justify-between mb-2">
              <span className="text-body-sm font-bold">{order.number}</span>
              <span className={cn('text-caption font-bold', order.timer > 600 ? 'text-error' : order.timer > 300 ? 'text-warning' : 'text-success')}>
                {Math.floor(order.timer / 60)}:{String(order.timer % 60).padStart(2, '0')}
              </span>
            </div>
            <div className="space-y-1 mb-3">
              {order.items.map((item, i) => (
                <div key={i} className="text-caption text-text-primary">
                  {item.qty}× {item.name}
                  {item.notes && <p className="text-text-tertiary text-xs">→ {item.notes}</p>}
                </div>
              ))}
            </div>
            <button
              onClick={() => advanceOrder(order.id)}
              className={cn('w-full py-2 rounded-lg text-caption font-bold', order.status === 'preparing' ? 'bg-warning/10 text-warning' : 'bg-success/10 text-success')}
            >
              {order.status === 'preparing' ? '✅ جاهز' : '🛵 استلمه المندوب'}
            </button>
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Analytics Screen ============
function AnalyticsScreen() {
  return (
    <div className="p-6 space-y-6">
      <h1 className="text-h2 font-bold text-text-primary">📈 التقارير</h1>
      <div className="grid grid-cols-2 gap-4">
        <div className="bg-surface rounded-lg border p-5">
          <h3 className="text-body font-semibold mb-3">📊 مبيعات الأسبوع</h3>
          <div className="flex items-end gap-2 h-40">
            {[320, 280, 450, 380, 520, 610, 425].map((val, i) => (
              <div key={i} className="flex-1 flex flex-col items-center gap-1">
                <div className="w-full bg-primary/20 rounded-t" style={{ height: `${(val / 700) * 100}%` }}>
                  <div className="w-full bg-primary rounded-t h-full opacity-60" />
                </div>
                <span className="text-overline text-text-tertiary">{['سبت', 'أحد', 'إثن', 'ثلا', 'أرب', 'خمي', 'جمع'][i]}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="bg-surface rounded-lg border p-5">
          <h3 className="text-body font-semibold mb-3">🍔 الأصناف الأكثر مبيعاً</h3>
          <div className="space-y-2">
            {[{ name: 'Margherita Pizza', count: 42 }, { name: 'Pepperoni Pizza', count: 38 }, { name: 'Apple Pie', count: 35 }, { name: 'Coca Cola', count: 28 }].map((item, i) => (
              <div key={i} className="flex items-center justify-between text-body-sm">
                <span className="text-text-primary">{i + 1}. {item.name}</span>
                <span className="font-bold text-primary">×{item.count}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

// ============ Schedule Screen ============
function ScheduleScreen() {
  const days = ['السبت', 'الأحد', 'الإثنين', 'الثلاثاء', 'الأربعاء', 'الخميس', 'الجمعة']
  return (
    <div className="p-6 space-y-4">
      <h1 className="text-h2 font-bold text-text-primary">📅 ساعات العمل</h1>
      <div className="space-y-2 max-w-lg">
        {days.map((day, i) => (
          <div key={i} className="flex items-center gap-3 bg-surface rounded-lg border p-3">
            <span className="text-body-sm font-medium w-24">{day}</span>
            <button className={cn('w-12 h-6 rounded-full', i < 6 ? 'bg-success' : 'bg-gray-300')}>
              <span className={cn('block w-5 h-5 bg-white rounded-full transition-transform', i < 6 ? 'translate-x-0' : 'translate-x-6')} />
            </button>
            {i < 6 ? (
              <span className="text-body-sm text-text-secondary">10:00 ص — 02:00 ص</span>
            ) : (
              <span className="text-body-sm text-text-tertiary">مقفل</span>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Components ============

function StatCard({ icon, label, value, trend, color }: { icon: string; label: string; value: string; trend: string; color: string }) {
  const colors: Record<string, string> = {
    primary: 'bg-primary/10 text-primary',
    success: 'bg-success/10 text-success',
    warning: 'bg-warning/10 text-warning',
    accent: 'bg-purple-100 text-purple-600',
  }
  return (
    <div className="bg-surface rounded-lg border p-4">
      <div className={cn('w-10 h-10 rounded-lg flex items-center justify-center mb-2', colors[color])}>
        <span className="material-symbols-rounded">{icon}</span>
      </div>
      <p className="text-caption text-text-tertiary">{label}</p>
      <p className="text-h3 font-bold text-text-primary mt-1">{value}</p>
      <p className="text-overline text-success mt-1">{trend}</p>
    </div>
  )
}

function OrderCard({ order }: { order: typeof mockOrders[0] }) {
  const statusColors: Record<string, string> = {
    preparing: 'bg-warning/10 text-warning',
    ready: 'bg-success/10 text-success',
  }
  return (
    <div className="bg-surface rounded-lg border p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-body font-bold">{order.number}</span>
        <span className={cn('text-caption px-2 py-0.5 rounded', statusColors[order.status])}>
          {order.status === 'preparing' ? '🟡 يتحضّر' : '🟢 جاهز'}
        </span>
      </div>
      <div className="space-y-0.5 mb-2">
        {order.items.map((item, i) => (
          <p key={i} className="text-caption text-text-secondary">{item.qty}× {item.name}</p>
        ))}
      </div>
      <div className="flex items-center justify-between text-caption">
        <span className="font-bold text-text-primary">EGP {order.total}</span>
        <span className="text-text-tertiary">{order.payment === 'vodafone_cash' ? '💚 VF Cash' : order.payment}</span>
      </div>
    </div>
  )
}

function OrderDetailCard({ order, onUpdateStatus }: { order: typeof mockOrders[0]; onUpdateStatus: (id: string, status: string) => void }) {
  return (
    <div className="bg-surface rounded-lg border p-4">
      <div className="flex items-center justify-between mb-3">
        <div>
          <span className="text-body font-bold">{order.number}</span>
          <span className="text-caption text-text-tertiary mr-2">{order.customer}</span>
        </div>
        <span className="text-caption text-text-tertiary">{Math.floor(order.timer / 60)}:{String(order.timer % 60).padStart(2, '0')}</span>
      </div>
      <div className="space-y-1 mb-3">
        {order.items.map((item, i) => (
          <div key={i} className="text-body-sm">
            <span className="font-medium">{item.qty}× {item.name}</span>
            {item.notes && <span className="text-text-tertiary mr-2">— {item.notes}</span>}
          </div>
        ))}
      </div>
      <div className="flex items-center justify-between">
        <span className="text-body-sm font-bold">EGP {order.total} • {order.payment}</span>
        <div className="flex gap-2">
          {order.status === 'preparing' && (
            <button onClick={() => onUpdateStatus(order.id, 'ready')} className="px-4 py-1.5 bg-success/10 text-success rounded-lg text-caption font-semibold hover:bg-success/20">
              ✅ جاهز للاستلام
            </button>
          )}
          <button className="px-4 py-1.5 bg-warning/10 text-warning rounded-lg text-caption font-semibold hover:bg-warning/20">
            ⚠️ تأخير
          </button>
        </div>
      </div>
    </div>
  )
}

// ============ Inbound Order Modal (90s timer) ============
function InboundOrderModal({ order, onAccept, onReject }: { order: typeof mockIncomingOrder; onAccept: () => void; onReject: () => void }) {
  const [seconds, setSeconds] = useState(90)

  useEffect(() => {
    if (seconds <= 0) {
      onReject()
      return
    }
    const timer = setTimeout(() => setSeconds(seconds - 1), 1000)
    return () => clearTimeout(timer)
  }, [seconds, onReject])

  const progress = (seconds / 90) * 100
  const isUrgent = seconds < 30

  return (
    <div className="fixed inset-0 z-modal flex items-center justify-center bg-black/50" dir="rtl">
      <div className="bg-surface rounded-2xl shadow-2xl max-w-md w-full mx-4 overflow-hidden">
        {/* Timer */}
        <div className={cn('p-4 text-center', isUrgent ? 'bg-error' : seconds < 60 ? 'bg-warning' : 'bg-primary')}>
          <p className="text-white text-h2 font-bold">{seconds}</p>
          <p className="text-white/80 text-caption">ثانية متبقية</p>
          <div className="w-full h-1 bg-white/30 rounded-full mt-2">
            <div className="h-full bg-white rounded-full transition-all" style={{ width: `${progress}%` }} />
          </div>
        </div>

        {/* Order details */}
        <div className="p-5">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-h3 font-bold">🔔 طلب جديد {order.number}</h2>
          </div>

          <div className="space-y-2 mb-4">
            {order.items.map((item, i) => (
              <div key={i} className="flex justify-between text-body-sm">
                <div>
                  <span className="font-medium">{item.qty}× {item.name}</span>
                  {item.notes && <p className="text-text-tertiary text-xs">→ {item.notes}</p>}
                </div>
              </div>
            ))}
          </div>

          <div className="border-t pt-3 space-y-1 text-body-sm">
            <div className="flex justify-between"><span className="text-text-secondary">العميل</span><span className="font-medium">{order.customer}</span></div>
            <div className="flex justify-between"><span className="text-text-secondary">العنوان</span><span className="font-medium">{order.address}</span></div>
            <div className="flex justify-between"><span className="text-text-secondary">ETA</span><span className="font-medium">{order.eta} دقيقة</span></div>
            <div className="flex justify-between"><span className="text-text-secondary">الدفع</span><span className="font-medium">💚 Vodafone Cash</span></div>
            <div className="flex justify-between font-bold"><span>الإجمالي</span><span className="text-primary">EGP {order.total}</span></div>
          </div>

          {/* Buttons */}
          <div className="flex gap-3 mt-5">
            <button onClick={onReject} className="flex-1 py-3 border-2 border-border text-text-secondary rounded-lg font-semibold hover:bg-bg-tertiary">
              ❌ رفض
            </button>
            <button onClick={onAccept} className="flex-1 py-3 bg-success text-white rounded-lg font-semibold hover:bg-success/90">
              ✅ قبول
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
