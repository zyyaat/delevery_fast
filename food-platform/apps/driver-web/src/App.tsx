// Driver Web App — Full implementation
// 14 screens: Login, OTP, KYC, Home, Order Offer, Pickup, Dropoff, Completed,
// Earnings, Payout, History, Profile, Heat Map

import { useState, useEffect, useRef } from 'react'
import { cn } from '@food-platform/ui'

type Screen = 'login' | 'otp' | 'kyc' | 'home' | 'pickup' | 'dropoff' | 'completed' | 'earnings' | 'payout' | 'history' | 'profile' | 'heatmap'

export default function App() {
  const [screen, setScreen] = useState<Screen>('login')
  const [phone, setPhone] = useState('')
  const [otp, setOtp] = useState(['', '', '', '', '', ''])
  const [isOnline, setIsOnline] = useState(false)
  const [orderOffer, setOrderOffer] = useState<typeof mockOrderOffer | null>(null)
  const [activeOrder, setActiveOrder] = useState<typeof mockActiveOrder | null>(null)
  const [deliveryStep, setDeliveryStep] = useState(0)
  const [todayEarnings, setTodayEarnings] = useState(245.5)
  const [todayDeliveries, setTodayDeliveries] = useState(7)

  // Simulate incoming order when online
  useEffect(() => {
    if (isOnline && !orderOffer && !activeOrder && screen === 'home') {
      const timer = setTimeout(() => {
        setOrderOffer({ ...mockOrderOffer })
      }, 5000)
      return () => clearTimeout(timer)
    }
  }, [isOnline, orderOffer, activeOrder, screen])

  // Simulate delivery step progression
  useEffect(() => {
    if (screen === 'pickup' && deliveryStep < 2) {
      const timer = setTimeout(() => setDeliveryStep(deliveryStep + 1), 3000)
      return () => clearTimeout(timer)
    }
  }, [screen, deliveryStep])

  const acceptOrder = () => {
    if (orderOffer) {
      setActiveOrder({ ...mockActiveOrder, ...orderOffer })
      setOrderOffer(null)
      setDeliveryStep(0)
      setScreen('pickup')
    }
  }

  const completeDelivery = () => {
    setTodayEarnings(todayEarnings + 28.5)
    setTodayDeliveries(todayDeliveries + 1)
    setActiveOrder(null)
    setScreen('completed')
  }

  return (
    <div className="min-h-screen bg-bg-primary" dir="rtl">
      {screen === 'login' && <LoginScreen phone={phone} setPhone={setPhone} onNext={() => setScreen('otp')} />}
      {screen === 'otp' && <OTPScreen phone={phone} otp={otp} setOtp={setOtp} onNext={() => setScreen('kyc')} onBack={() => setScreen('login')} />}
      {screen === 'kyc' && <KYCScreen onNext={() => setScreen('home')} />}
      {screen === 'home' && <HomeScreen isOnline={isOnline} setIsOnline={setIsOnline} todayEarnings={todayEarnings} todayDeliveries={todayDeliveries} onNavigate={setScreen} />}
      {screen === 'pickup' && activeOrder && <PickupScreen order={activeOrder} deliveryStep={deliveryStep} onArrived={() => setDeliveryStep(2)} onPickedUp={() => setScreen('dropoff')} />}
      {screen === 'dropoff' && activeOrder && <DropoffScreen order={activeOrder} onDelivered={completeDelivery} />}
      {screen === 'completed' && <CompletedScreen earnings={28.5} onContinue={() => setScreen('home')} />}
      {screen === 'earnings' && <EarningsScreen todayEarnings={todayEarnings} todayDeliveries={todayDeliveries} onNavigate={setScreen} />}
      {screen === 'payout' && <PayoutScreen pendingPayout={todayEarnings} onBack={() => setScreen('earnings')} />}
      {screen === 'history' && <HistoryScreen onNavigate={setScreen} />}
      {screen === 'profile' && <ProfileScreen onNavigate={setScreen} />}
      {screen === 'heatmap' && <HeatMapScreen onBack={() => setScreen('home')} />}

      {/* Order offer modal (15s timer) */}
      {orderOffer && screen === 'home' && (
        <OrderOfferModal order={orderOffer} onAccept={acceptOrder} onReject={() => setOrderOffer(null)} />
      )}

      {/* Bottom nav (only on certain screens) */}
      {['home', 'earnings', 'history', 'profile'].includes(screen) && (
        <BottomNav active={screen} onNavigate={setScreen} />
      )}
    </div>
  )
}

// ============ Mock Data ============
const mockOrderOffer = {
  id: 'order-1',
  number: '#A7X92F',
  restaurant: { name: 'Pizza Hut', address: 'معادي - 1.2km', distance: 1.2 },
  customer: { area: 'الدقي - 4.5km', distance: 4.5 },
  eta: 22,
  earnings: { total: 28.5, base: 20, distance: 5.5, peak: 3 },
  pickupCode: '4892',
}

const mockActiveOrder = {
  id: 'order-1',
  number: '#A7X92F',
  restaurant: { name: 'Pizza Hut', address: 'معادي - 1.2km', distance: 1.2 },
  customer: { name: 'Ahmed M.', phone: '01012345678', area: 'الدقي', address: '12 شارع التحرير، شقة 5', notes: 'باب أزرق، اتصل قبل الوصول' },
  eta: 22,
  earnings: { total: 28.5, base: 20, distance: 5.5, peak: 3 },
  pickupCode: '4892',
  dropoffCode: '7294',
  items: [
    { name: 'Margherita Pizza (وسط)', qty: 1, notes: 'بدون بصل' },
    { name: 'Pepperoni Pizza (كبير)', qty: 2 },
    { name: 'Coca Cola', qty: 2 },
  ],
}

// ============ Login Screen ============
function LoginScreen({ phone, setPhone, onNext }: { phone: string; setPhone: (v: string) => void; onNext: () => void }) {
  return (
    <div className="min-h-screen bg-gradient-to-b from-orange-50 to-white flex flex-col items-center justify-center p-6">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-orange-500 rounded-3xl mb-4 text-4xl">🛵</div>
          <h1 className="text-3xl font-bold text-gray-900">مندوب توصيل</h1>
          <p className="text-gray-500 mt-2">سجّل رقم موبايلك عشان تبعتلك كود</p>
        </div>
        <div className="space-y-4">
          <div>
            <label className="text-sm text-gray-600 mb-1.5 block">رقم الموبايل</label>
            <div className="flex items-center gap-2">
              <div className="flex items-center gap-1 px-3 h-12 rounded-lg border border-gray-200 bg-gray-50 text-sm font-medium">🇪🇬 +20</div>
              <input
                type="tel" inputMode="numeric" maxLength={11} placeholder="01 2345 6789"
                value={phone} onChange={(e) => setPhone(e.target.value.replace(/\D/g, ''))}
                className="w-full h-12 px-3 rounded-lg border border-gray-200 text-lg focus:outline-none focus:ring-2 focus:ring-orange-400"
              />
            </div>
          </div>
          <button
            className="w-full h-12 text-lg bg-orange-500 hover:bg-orange-600 text-white rounded-lg font-semibold transition-colors disabled:opacity-50"
            disabled={phone.length < 11} onClick={onNext}
          >أرسل الكود</button>
        </div>
        <p className="text-xs text-gray-400 text-center">بتسجيلك، أنت موافق على الشروط والأحكام</p>
      </div>
    </div>
  )
}

// ============ OTP Screen ============
function OTPScreen({ phone, otp, setOtp, onNext, onBack }: { phone: string; otp: string[]; setOtp: (v: string[]) => void; onNext: () => void; onBack: () => void }) {
  const inputsRef = useRef<(HTMLInputElement | null)[]>([])

  useEffect(() => { inputsRef.current[0]?.focus() }, [])

  const handleChange = (i: number, val: string) => {
    const digit = val.replace(/\D/g, '').slice(-1)
    const newOtp = [...otp]; newOtp[i] = digit; setOtp(newOtp)
    if (digit && i < 5) inputsRef.current[i + 1]?.focus()
    if (newOtp.every(d => d) && newOtp.join('').length === 6) setTimeout(onNext, 500)
  }

  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-center p-6">
      <div className="w-full max-w-md space-y-6">
        <button onClick={onBack} className="flex items-center gap-1 text-gray-500 hover:text-gray-900">
          <span className="text-xl">→</span><span>رجوع</span>
        </button>
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900">أدخل الكود</h1>
          <p className="text-gray-500 mt-2">بعتناك كود على</p>
          <p className="font-semibold text-gray-900" dir="ltr">+20 {phone}</p>
        </div>
        <div className="flex justify-center gap-2" dir="ltr">
          {otp.map((d, i) => (
            <input key={i} ref={(el) => { inputsRef.current[i] = el }} type="text" inputMode="numeric" maxLength={1} value={d}
              onChange={(e) => handleChange(i, e.target.value)}
              className={`w-12 h-14 text-center text-2xl font-bold rounded-lg border-2 ${d ? 'border-orange-500 bg-orange-50' : 'border-gray-200'} focus:outline-none`}
            />
          ))}
        </div>
        <button className="w-full h-12 bg-orange-500 hover:bg-orange-600 text-white rounded-lg font-semibold disabled:opacity-50"
          disabled={otp.some(d => !d)} onClick={onNext}>تأكيد</button>
      </div>
    </div>
  )
}

// ============ KYC Screen ============
function KYCScreen({ onNext }: { onNext: () => void }) {
  const [step, setStep] = useState(1)
  const [name, setName] = useState('')
  const [vehicleType, setVehicleType] = useState('motorcycle')

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-md mx-auto space-y-5">
        <div className="flex items-center gap-2 mb-4">
          <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
            <div className="h-full bg-orange-500 rounded-full transition-all" style={{ width: `${(step / 3) * 100}%` }} />
          </div>
          <span className="text-sm text-gray-400">خطوة {step} من 3</span>
        </div>

        {step === 1 && (
          <div className="space-y-4">
            <h1 className="text-2xl font-bold text-gray-900">المعلومات الأساسية</h1>
            <div>
              <label className="text-sm text-gray-600 block mb-1">الاسم الكامل</label>
              <input type="text" value={name} onChange={(e) => setName(e.target.value)} placeholder="محمود سعيد"
                className="w-full h-12 px-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-orange-400" />
            </div>
            <div>
              <label className="text-sm text-gray-600 block mb-2">نوع المركبة</label>
              <div className="space-y-2">
                {[{ id: 'motorcycle', label: '🛵 موتوسيكل' }, { id: 'car', label: '🚗 عربية' }, { id: 'bicycle', label: '🚲 دراجة' }].map(v => (
                  <label key={v.id} className={`flex items-center gap-3 p-3 rounded-lg border-2 cursor-pointer ${vehicleType === v.id ? 'border-orange-500 bg-orange-50' : 'border-gray-200'}`}>
                    <input type="radio" checked={vehicleType === v.id} onChange={() => setVehicleType(v.id)} className="w-5 h-5 accent-orange-500" />
                    <span className="text-lg">{v.label}</span>
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="text-sm text-gray-600 block mb-2">📷 صورة البطاقة (أمامي)</label>
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center cursor-pointer hover:border-orange-400">
                <span className="text-4xl">📷</span>
                <p className="text-sm text-gray-400 mt-2">ارفع الصورة</p>
              </div>
            </div>
            <button className="w-full h-12 bg-orange-500 hover:bg-orange-600 text-white rounded-lg font-semibold" onClick={() => setStep(2)}>التالي</button>
          </div>
        )}

        {step === 2 && (
          <div className="space-y-4">
            <h1 className="text-2xl font-bold text-gray-900">معلومات المركبة</h1>
            {['📷 رخصة القيادة', '📷 لوحات المركبة', '📷 تأمين المركبة'].map(label => (
              <div key={label}>
                <label className="text-sm text-gray-600 block mb-1">{label}</label>
                <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center cursor-pointer hover:border-orange-400">
                  <span className="text-3xl">📷</span><p className="text-sm text-gray-400 mt-1">ارفع الصورة</p>
                </div>
              </div>
            ))}
            <div className="flex gap-2">
              <button className="flex-1 h-12 border-2 border-gray-200 text-gray-600 rounded-lg font-semibold" onClick={() => setStep(1)}>رجوع</button>
              <button className="flex-1 h-12 bg-orange-500 hover:bg-orange-600 text-white rounded-lg font-semibold" onClick={() => setStep(3)}>التالي</button>
            </div>
          </div>
        )}

        {step === 3 && (
          <div className="space-y-4">
            <h1 className="text-2xl font-bold text-gray-900">طريقة استلام الأرباح</h1>
            <label className="flex items-center gap-3 p-4 rounded-lg border-2 border-orange-500 bg-orange-50 cursor-pointer">
              <input type="radio" defaultChecked className="w-5 h-5 accent-orange-500" />
              <span className="text-2xl">💚</span>
              <div><p className="font-semibold">Vodafone Cash</p><p className="text-sm text-gray-400">رقم المحفظة: {phone || '01012345678'}</p></div>
            </label>
            <label className="flex items-center gap-3 p-4 rounded-lg border-2 border-gray-200 cursor-pointer">
              <input type="radio" className="w-5 h-5 accent-orange-500" />
              <span className="text-2xl">🟣</span>
              <span className="font-semibold">InstaPay</span>
            </label>
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <p className="text-sm text-blue-700">⏰ جاري المراجعة — هنتواصل معاك في 48 ساعة لتفعيل حسابك</p>
            </div>
            <button className="w-full h-12 bg-orange-500 hover:bg-orange-600 text-white rounded-lg font-semibold" onClick={onNext}>خلينا نبدأ 🚀</button>
          </div>
        )}
      </div>
    </div>
  )
}

// ============ Home Screen ============
function HomeScreen({ isOnline, setIsOnline, todayEarnings, todayDeliveries, onNavigate }: { isOnline: boolean; setIsOnline: (v: boolean) => void; todayEarnings: number; todayDeliveries: number; onNavigate: (s: Screen) => void }) {
  return (
    <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={() => setIsOnline(!isOnline)} className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-semibold ${isOnline ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'}`}>
          <span className={`w-2 h-2 rounded-full ${isOnline ? 'bg-green-500 animate-pulse' : 'bg-gray-400'}`} />
          {isOnline ? 'متاح' : 'غير متاح'}
        </button>
        <div className="flex-1 text-center">
          <p className="text-xs text-gray-400">أرباح اليوم</p>
          <p className="text-lg font-bold text-orange-500">EGP {todayEarnings.toFixed(2)}</p>
        </div>
        <button onClick={() => onNavigate('profile')} className="w-10 h-10 rounded-full bg-orange-100 flex items-center justify-center">
          <span className="text-xl">👤</span>
        </button>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        {/* Earnings cards */}
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-white rounded-xl border p-4">
            <p className="text-sm text-gray-400">💰 اليوم</p>
            <p className="text-2xl font-bold text-gray-900 mt-1">EGP {todayEarnings.toFixed(2)}</p>
            <p className="text-xs text-gray-400 mt-1">{todayDeliveries} توصيلات • 4.2 ساعة</p>
            <p className="text-xs text-green-500 mt-1">المعدل: EGP {(todayEarnings / 4.2).toFixed(1)}/ساعة</p>
          </div>
          <div className="bg-white rounded-xl border p-4">
            <p className="text-sm text-gray-400">📅 هذا الأسبوع</p>
            <p className="text-2xl font-bold text-gray-900 mt-1">EGP 1,820</p>
            <p className="text-xs text-gray-400 mt-1">52 توصيلة</p>
          </div>
        </div>

        {/* Heat map */}
        <button onClick={() => onNavigate('heatmap')} className="w-full bg-white rounded-xl border p-4 text-right">
          <div className="flex items-center justify-between mb-2">
            <h3 className="font-bold text-gray-900">🔥 مناطق الطلب العالي</h3>
            <span className="text-xs text-orange-500 font-semibold">عرض الكل</span>
          </div>
          <div className="flex gap-2">
            <div className="flex-1 bg-red-50 border border-red-200 rounded-lg p-2 text-center">
              <p className="text-xs text-red-600 font-semibold">معادي</p><p className="text-xs text-red-400">3 طلبات نشطة</p>
            </div>
            <div className="flex-1 bg-red-50 border border-red-200 rounded-lg p-2 text-center">
              <p className="text-xs text-red-600 font-semibold">الزمالك</p><p className="text-xs text-red-400">5 طلبات نشطة</p>
            </div>
            <div className="flex-1 bg-yellow-50 border border-yellow-200 rounded-lg p-2 text-center">
              <p className="text-xs text-yellow-600 font-semibold">مدينة نصر</p><p className="text-xs text-yellow-400">2 طلبات</p>
            </div>
          </div>
        </button>

        {/* Tier */}
        <div className="bg-gradient-to-l from-amber-400 to-yellow-500 rounded-xl p-4 text-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm opacity-90">المستوى</p>
              <p className="text-2xl font-bold">🥇 Platinum</p>
              <p className="text-xs opacity-80 mt-1">⭐ 4.92 (148) | قبول: 91% | إكمال: 96%</p>
            </div>
            <span className="text-4xl">🏆</span>
          </div>
        </div>

        {/* Toggle online */}
        <button
          onClick={() => setIsOnline(!isOnline)}
          className={`w-full h-14 rounded-xl font-bold text-lg transition-colors ${isOnline ? 'bg-red-500 hover:bg-red-600 text-white' : 'bg-orange-500 hover:bg-orange-600 text-white'}`}
        >
          {isOnline ? '🛑 إنهاء العمل' : '▶️ ابدأ العمل'}
        </button>

        {isOnline && !orderOffer && (
          <div className="text-center py-8">
            <div className="inline-block animate-spin h-8 w-8 border-4 border-orange-200 border-t-orange-500 rounded-full mb-3" />
            <p className="text-gray-500">بنبحث لك على طلبات...</p>
          </div>
        )}
      </div>
    </div>
  )
}

// ============ Order Offer Modal (15s timer) ============
function OrderOfferModal({ order, onAccept, onReject }: { order: typeof mockOrderOffer; onAccept: () => void; onReject: () => void }) {
  const [seconds, setSeconds] = useState(15)

  useEffect(() => {
    if (seconds <= 0) { onReject(); return }
    const timer = setTimeout(() => setSeconds(seconds - 1), 1000)
    return () => clearTimeout(timer)
  }, [seconds, onReject])

  const progress = (seconds / 15) * 100
  const isUrgent = seconds <= 5

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" dir="rtl">
      <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full overflow-hidden">
        {/* Timer */}
        <div className={`p-4 text-center ${isUrgent ? 'bg-red-500' : 'bg-orange-500'}`}>
          <p className="text-white text-4xl font-bold">{seconds}</p>
          <p className="text-white/80 text-sm">ثانية متبقية</p>
          <div className="w-full h-1 bg-white/30 rounded-full mt-2">
            <div className="h-full bg-white rounded-full transition-all" style={{ width: `${progress}%` }} />
          </div>
        </div>

        <div className="p-5">
          <h2 className="text-xl font-bold text-gray-900 mb-4">🔔 طلب جديد {order.number}</h2>

          {/* Route */}
          <div className="space-y-3 mb-4">
            <div className="flex items-center gap-3">
              <span className="text-2xl">🍔</span>
              <div className="flex-1">
                <p className="font-semibold text-gray-900">{order.restaurant.name}</p>
                <p className="text-sm text-gray-400">{order.restaurant.address}</p>
              </div>
            </div>
            <div className="border-r-2 border-dashed border-gray-200 mr-5 h-4" />
            <div className="flex items-center gap-3">
              <span className="text-2xl">📍</span>
              <div className="flex-1">
                <p className="font-semibold text-gray-900">العميل</p>
                <p className="text-sm text-gray-400">{order.customer.area}</p>
              </div>
            </div>
          </div>

          <div className="bg-gray-50 rounded-lg p-3 mb-4 flex items-center justify-between text-sm">
            <span className="text-gray-400">⏱️ ETA التوصيل</span>
            <span className="font-bold text-gray-900">{order.eta} دقيقة</span>
          </div>

          {/* Earnings */}
          <div className="bg-orange-50 rounded-lg p-4 mb-4">
            <p className="text-sm text-gray-400 mb-2">💰 أرباحك من الطلب</p>
            <p className="text-3xl font-bold text-orange-500 mb-2">EGP {order.earnings.total.toFixed(2)}</p>
            <div className="space-y-1 text-xs text-gray-500">
              <div className="flex justify-between"><span>رسوم التوصيل</span><span>EGP {order.earnings.base.toFixed(2)}</span></div>
              <div className="flex justify-between"><span>مكافأة المسافة</span><span>EGP {order.earnings.distance.toFixed(2)}</span></div>
              <div className="flex justify-between"><span>Peak bonus</span><span>EGP {order.earnings.peak.toFixed(2)}</span></div>
            </div>
            <p className="text-xs text-green-600 mt-2">⚡ Instant payout متاح</p>
          </div>

          <div className="flex gap-3">
            <button onClick={onReject} className="flex-1 h-12 border-2 border-gray-200 text-gray-600 rounded-lg font-semibold hover:bg-gray-50">رفض</button>
            <button onClick={onAccept} className="flex-1 h-12 bg-green-500 hover:bg-green-600 text-white rounded-lg font-bold">✅ قبول</button>
          </div>
        </div>
      </div>
    </div>
  )
}

// ============ Pickup Screen ============
function PickupScreen({ order, deliveryStep, onArrived, onPickedUp }: { order: typeof mockActiveOrder; deliveryStep: number; onArrived: () => void; onPickedUp: () => void }) {
  const steps = ['الذهاب للمطعم', 'استلام الطلب', 'التوصيل للعميل']

  return (
    <div className="min-h-screen bg-gray-50 pb-24" dir="rtl">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center justify-between">
        <span className="font-bold text-gray-900">{order.number}</span>
        <span className="text-sm text-gray-400">خطوة {deliveryStep + 1} من 3</span>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        {/* Progress */}
        <div className="flex items-center gap-2">
          {steps.map((step, i) => (
            <div key={i} className="flex items-center gap-2 flex-1">
              <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm ${i < deliveryStep ? 'bg-green-500 text-white' : i === deliveryStep ? 'bg-orange-500 text-white animate-pulse' : 'bg-gray-200 text-gray-400'}`}>
                {i < deliveryStep ? '✓' : i + 1}
              </div>
              <span className={`text-xs ${i <= deliveryStep ? 'text-gray-900 font-medium' : 'text-gray-300'}`}>{step}</span>
              {i < 2 && <div className={`flex-1 h-0.5 ${i < deliveryStep ? 'bg-green-500' : 'bg-gray-200'}`} />}
            </div>
          ))}
        </div>

        {/* Map placeholder */}
        <div className="bg-gradient-to-br from-blue-100 to-green-100 rounded-xl h-48 flex items-center justify-center">
          <div className="text-center">
            <span className="text-5xl">🗺️</span>
            <p className="text-sm text-gray-500 mt-2">{order.restaurant.name} — {order.restaurant.distance}km</p>
            <button className="mt-2 px-4 py-2 bg-blue-500 text-white rounded-lg text-sm font-semibold">📍 ابدأ الملاحة</button>
          </div>
        </div>

        {/* Restaurant info */}
        <div className="bg-white rounded-xl border p-4">
          <h3 className="font-bold text-gray-900">🍔 {order.restaurant.name}</h3>
          <p className="text-sm text-gray-400 mt-1">{order.restaurant.address}</p>
        </div>

        {/* Items */}
        <div className="bg-white rounded-xl border p-4">
          <h3 className="font-bold text-gray-900 mb-2">📋 تفاصيل الطلب</h3>
          {order.items.map((item, i) => (
            <div key={i} className="flex justify-between text-sm py-1">
              <span>{item.qty}× {item.name}</span>
              {item.notes && <span className="text-gray-400 text-xs">→ {item.notes}</span>}
            </div>
          ))}
          <div className="mt-3 bg-orange-50 rounded-lg p-3 text-center">
            <p className="text-xs text-gray-400">رمز الاستلام</p>
            <p className="text-2xl font-bold tracking-widest text-orange-500">{order.pickupCode}</p>
            <p className="text-xs text-gray-400 mt-1">اعرضه للمطعم</p>
          </div>
        </div>

        {/* Actions */}
        <div className="space-y-2">
          {deliveryStep === 0 && (
            <button onClick={onArrived} className="w-full h-12 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-semibold">🟢 أنا في المطعم</button>
          )}
          {deliveryStep === 1 && (
            <button onClick={onPickedUp} className="w-full h-12 bg-green-500 hover:bg-green-600 text-white rounded-lg font-bold">✅ استلمت الطلب</button>
          )}
          {deliveryStep < 2 && (
            <button className="w-full h-10 border border-gray-200 text-gray-500 rounded-lg text-sm">⚠️ المطعم مش جاهز</button>
          )}
        </div>
      </div>
    </div>
  )
}

// ============ Dropoff Screen ============
function DropoffScreen({ order, onDelivered }: { order: typeof mockActiveOrder; onDelivered: () => void }) {
  const [otpDigits, setOtpDigits] = useState(['', '', '', ''])
  const inputsRef = useRef<(HTMLInputElement | null)[]>([])
  const otpComplete = otpDigits.every(d => d)

  const handleOtp = (i: number, val: string) => {
    const digit = val.replace(/\D/g, '').slice(-1)
    const newOtp = [...otpDigits]; newOtp[i] = digit; setOtpDigits(newOtp)
    if (digit && i < 3) inputsRef.current[i + 1]?.focus()
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-24" dir="rtl">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center justify-between">
        <span className="font-bold text-gray-900">{order.number}</span>
        <span className="text-sm text-green-500">خطوة 3 من 3 — التوصيل</span>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        {/* Map */}
        <div className="bg-gradient-to-br from-green-100 to-yellow-100 rounded-xl h-48 flex items-center justify-center">
          <div className="text-center">
            <span className="text-5xl">🗺️</span>
            <p className="text-sm text-gray-500 mt-2">{order.customer.area} — {order.customer.address}</p>
            <button className="mt-2 px-4 py-2 bg-blue-500 text-white rounded-lg text-sm font-semibold">📍 ابدأ الملاحة</button>
          </div>
        </div>

        {/* Customer info */}
        <div className="bg-white rounded-xl border p-4">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-12 h-12 rounded-full bg-orange-100 flex items-center justify-center text-xl">👤</div>
            <div className="flex-1">
              <p className="font-bold text-gray-900">{order.customer.name}</p>
              <p className="text-sm text-gray-400">{order.customer.phone}</p>
            </div>
            <a href={`tel:${order.customer.phone}`} className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">📞</a>
            <button className="w-10 h-10 rounded-full bg-orange-100 flex items-center justify-center">💬</button>
          </div>
          <div className="bg-gray-50 rounded-lg p-3 text-sm">
            <p className="text-gray-500">📍 {order.customer.address}</p>
            {order.customer.notes && <p className="text-gray-400 mt-1">📝 {order.customer.notes}</p>}
          </div>
        </div>

        {/* OTP */}
        <div className="bg-white rounded-xl border p-4 text-center">
          <h3 className="font-bold text-gray-900 mb-2">🔢 رمز التسليم</h3>
          <p className="text-sm text-gray-400 mb-3">اطلبه من العميل</p>
          <div className="flex justify-center gap-3 mb-4" dir="ltr">
            {otpDigits.map((d, i) => (
              <input key={i} ref={(el) => { inputsRef.current[i] = el }} type="text" inputMode="numeric" maxLength={1} value={d}
                onChange={(e) => handleOtp(i, e.target.value)}
                className={`w-14 h-16 text-center text-2xl font-bold rounded-lg border-2 ${d ? 'border-orange-500 bg-orange-50' : 'border-gray-200'} focus:outline-none`}
              />
            ))}
          </div>
          <button
            onClick={onDelivered}
            disabled={!otpComplete}
            className="w-full h-12 bg-green-500 hover:bg-green-600 text-white rounded-lg font-bold disabled:opacity-50"
          >✅ تم التسليم</button>
        </div>

        {/* Alternative */}
        <div className="bg-white rounded-xl border p-4">
          <p className="text-sm text-gray-400 mb-2">أو التقط صورة للطلب كدليل على التسليم</p>
          <button className="w-full h-10 border-2 border-gray-200 rounded-lg text-sm text-gray-500">📷 التقاط صورة</button>
        </div>

        <button className="w-full h-10 border border-gray-200 text-gray-500 rounded-lg text-sm">⚠️ العميل مش موجود</button>
      </div>
    </div>
  )
}

// ============ Completed Screen ============
function CompletedScreen({ earnings, onContinue }: { earnings: number; onContinue: () => void }) {
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-6" dir="rtl">
      <div className="max-w-md w-full text-center">
        <div className="inline-flex items-center justify-center w-20 h-20 bg-green-500 rounded-full mb-4">
          <span className="text-5xl">✅</span>
        </div>
        <h1 className="text-2xl font-bold text-gray-900">تم التوصيل! 🎉</h1>

        <div className="bg-white rounded-xl border p-5 mt-6">
          <p className="text-sm text-gray-400">أرباحك من الطلب</p>
          <p className="text-4xl font-bold text-orange-500 mt-1">EGP {earnings.toFixed(2)}</p>
          <div className="space-y-1 text-sm mt-3 text-left">
            <div className="flex justify-between"><span className="text-gray-400">⏱️ المدة</span><span className="font-medium">28 دقيقة</span></div>
            <div className="flex justify-between"><span className="text-gray-400">📏 المسافة</span><span className="font-medium">5.7km</span></div>
          </div>
          <div className="border-t mt-3 pt-3">
            <p className="text-sm text-gray-400">💰 رصيدك القابل للسحب</p>
            <p className="text-xl font-bold text-gray-900">EGP {(245.5 + earnings).toFixed(2)}</p>
          </div>
        </div>

        <button onClick={onContinue} className="w-full h-12 mt-4 bg-orange-500 hover:bg-orange-600 text-white rounded-lg font-bold">طلب تاني 🛵</button>
        <button onClick={onContinue} className="w-full h-10 mt-2 text-gray-400 text-sm">العودة للرئيسية</button>
      </div>
    </div>
  )
}

// ============ Earnings Screen ============
function EarningsScreen({ todayEarnings, todayDeliveries, onNavigate }: { todayEarnings: number; todayDeliveries: number; onNavigate: (s: Screen) => void }) {
  const days = [
    { day: 'سبت', amount: 320 }, { day: 'أحد', amount: 280 }, { day: 'إثن', amount: 450 },
    { day: 'ثلا', amount: 380 }, { day: 'أرب', amount: 520 }, { day: 'خمي', amount: 610 }, { day: 'جمع', amount: todayEarnings },
  ]
  const maxAmount = Math.max(...days.map(d => d.amount))
  const weekTotal = days.reduce((sum, d) => sum + d.amount, 0)

  return (
    <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3"><h1 className="text-lg font-bold">💰 الأرباح</h1></header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        {/* Tabs */}
        <div className="flex gap-2">
          {['اليوم', 'الأسبوع', 'الشهر'].map((tab, i) => (
            <button key={i} className={`px-4 py-2 rounded-lg text-sm font-medium ${i === 1 ? 'bg-orange-500 text-white' : 'bg-white border text-gray-500'}`}>{tab}</button>
          ))}
        </div>

        {/* Summary */}
        <div className="bg-white rounded-xl border p-5">
          <p className="text-sm text-gray-400">هذا الأسبوع</p>
          <p className="text-3xl font-bold text-gray-900 mt-1">EGP {weekTotal.toFixed(2)}</p>
          <p className="text-sm text-gray-400 mt-1">{todayDeliveries * 7} توصيلة • 38 ساعة</p>
          <p className="text-sm text-green-500 mt-1">المعدل: EGP {(weekTotal / 38).toFixed(1)}/ساعة</p>
        </div>

        {/* Chart */}
        <div className="bg-white rounded-xl border p-5">
          <h3 className="font-bold mb-4">📊 الرسم البياني</h3>
          <div className="flex items-end gap-2 h-32">
            {days.map((d, i) => (
              <div key={i} className="flex-1 flex flex-col items-center gap-1">
                <div className="w-full bg-orange-200 rounded-t" style={{ height: `${(d.amount / maxAmount) * 100}%` }}>
                  <div className="w-full bg-orange-500 rounded-t h-full opacity-70" />
                </div>
                <span className="text-xs text-gray-400">{d.day}</span>
                <span className="text-xs font-bold text-gray-600">{d.amount}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Payout */}
        <div className="bg-white rounded-xl border p-5">
          <div className="flex items-center justify-between mb-3">
            <h3 className="font-bold">💵 رصيدك القابل للسحب</h3>
            <span className="text-2xl font-bold text-orange-500">EGP {weekTotal.toFixed(2)}</span>
          </div>
          <div className="flex gap-2">
            <button onClick={() => onNavigate('payout')} className="flex-1 h-10 bg-orange-500 hover:bg-orange-600 text-white rounded-lg text-sm font-semibold">⚡ سحب فوري</button>
            <button className="flex-1 h-10 border border-gray-200 text-gray-500 rounded-lg text-sm">📅 سحب يومي</button>
          </div>
        </div>

        {/* Recent deliveries */}
        <div className="bg-white rounded-xl border p-5">
          <h3 className="font-bold mb-3">📋 آخر التوصيلات</h3>
          {[
            { time: '14:30', order: '#A7X92F', amount: 28.5 },
            { time: '13:15', order: '#B3K45L', amount: 22.0 },
            { time: '12:00', order: '#C8M72N', amount: 35.0 },
          ].map((d, i) => (
            <div key={i} className="flex justify-between items-center py-2 border-b last:border-0 text-sm">
              <span className="text-gray-400">{d.time} — {d.order}</span>
              <span className="font-bold text-gray-900">EGP {d.amount.toFixed(2)} ✅</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

// ============ Payout Screen ============
function PayoutScreen({ pendingPayout, onBack }: { pendingPayout: number; onBack: () => void }) {
  const [amount, setAmount] = useState(pendingPayout.toString())
  const [method, setMethod] = useState('vodafone_cash')
  const [processing, setProcessing] = useState(false)
  const [success, setSuccess] = useState(false)

  const handlePayout = () => {
    setProcessing(true)
    setTimeout(() => { setProcessing(false); setSuccess(true) }, 2000)
  }

  if (success) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center p-6" dir="rtl">
        <div className="text-center">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-green-500 rounded-full mb-4"><span className="text-4xl">✅</span></div>
          <h1 className="text-2xl font-bold text-gray-900">تم! 💚</h1>
          <p className="text-gray-500 mt-2">EGP {(parseFloat(amount) - 2).toFixed(2)} في محفظتك</p>
          <button onClick={onBack} className="w-full h-12 mt-6 bg-orange-500 text-white rounded-lg font-semibold">تم</button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50" dir="rtl">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-500"><span className="text-xl">→</span></button>
        <h1 className="text-lg font-bold">⚡ سحب فوري</h1>
      </header>
      <div className="max-w-md mx-auto px-4 py-4 space-y-4">
        <div className="bg-white rounded-xl border p-5 text-center">
          <p className="text-sm text-gray-400">رصيدك القابل للسحب</p>
          <p className="text-3xl font-bold text-gray-900">EGP {pendingPayout.toFixed(2)}</p>
        </div>
        <div>
          <label className="text-sm text-gray-600 block mb-1">المبلغ</label>
          <input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} className="w-full h-12 px-3 rounded-lg border border-gray-200 text-lg font-bold focus:outline-none focus:ring-2 focus:ring-orange-400" />
        </div>
        <div className="space-y-2">
          <label className="flex items-center gap-3 p-3 rounded-lg border-2 border-orange-500 bg-orange-50 cursor-pointer">
            <input type="radio" checked={method === 'vodafone_cash'} onChange={() => setMethod('vodafone_cash')} className="w-5 h-5 accent-orange-500" />
            <span className="text-xl">💚</span><span className="font-semibold">Vodafone Cash</span>
          </label>
          <label className="flex items-center gap-3 p-3 rounded-lg border-2 border-gray-200 cursor-pointer">
            <input type="radio" checked={method === 'instapay'} onChange={() => setMethod('instapay')} className="w-5 h-5 accent-orange-500" />
            <span className="text-xl">🟣</span><span className="font-semibold">InstaPay</span>
          </label>
        </div>
        <div className="bg-gray-50 rounded-lg p-3 text-sm flex justify-between">
          <span className="text-gray-400">رسوم الخدمة</span><span>EGP 2.00</span>
        </div>
        <div className="bg-gray-50 rounded-lg p-3 text-sm flex justify-between font-bold">
          <span>هتستلم</span><span className="text-orange-500">EGP {(parseFloat(amount || '0') - 2).toFixed(2)}</span>
        </div>
        <button onClick={handlePayout} disabled={processing} className="w-full h-12 bg-orange-500 hover:bg-orange-600 text-white rounded-lg font-bold disabled:opacity-50">
          {processing ? '⏳ جاري التحويل...' : '⚡ تأكيد السحب'}
        </button>
      </div>
    </div>
  )
}

// ============ History Screen ============
function HistoryScreen({ onNavigate }: { onNavigate: (s: Screen) => void }) {
  const deliveries = [
    { time: '14:30', order: '#A7X92F', restaurant: 'Pizza Hut', customer: 'Ahmed M.', distance: 5.7, duration: 28, amount: 28.5, rating: 5 },
    { time: '13:15', order: '#B3K45L', restaurant: 'McDonald\'s', customer: 'Fatma A.', distance: 3.2, duration: 22, amount: 22.0, rating: 4 },
    { time: '12:00', order: '#C8M72N', restaurant: 'KFC', customer: 'Mahmoud K.', distance: 6.1, duration: 35, amount: 35.0, rating: 5 },
    { time: '10:45', order: '#D4N91P', restaurant: 'Sushi House', customer: 'Sara M.', distance: 2.8, duration: 18, amount: 25.0, rating: 5 },
  ]

  return (
    <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={() => onNavigate('home')} className="text-gray-500"><span className="text-xl">→</span></button>
        <h1 className="text-lg font-bold">📦 سجل التوصيلات</h1>
      </header>
      <div className="max-w-2xl mx-auto px-4 py-4 space-y-3">
        <div className="flex gap-2 mb-2">
          <select className="h-10 px-3 rounded-lg border border-gray-200 text-sm bg-white">
            <option>📅 اليوم</option><option>أمس</option><option>هذا الأسبوع</option><option>هذا الشهر</option>
          </select>
        </div>
        <p className="text-sm text-gray-400">─── اليوم (7) ───</p>
        {deliveries.map((d, i) => (
          <div key={i} className="bg-white rounded-xl border p-4">
            <div className="flex items-center justify-between mb-2">
              <div><span className="font-bold text-gray-900">{d.order}</span><span className="text-sm text-gray-400 mr-2">{d.time}</span></div>
              <span className="font-bold text-orange-500">EGP {d.amount.toFixed(2)}</span>
            </div>
            <p className="text-sm text-gray-500">{d.restaurant} → {d.customer}</p>
            <div className="flex items-center gap-3 mt-2 text-xs text-gray-400">
              <span>📏 {d.distance}km</span><span>⏱️ {d.duration} د</span>
              <span className="text-amber-400">{'★'.repeat(d.rating)}{'☆'.repeat(5 - d.rating)}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Profile Screen ============
function ProfileScreen({ onNavigate }: { onNavigate: (s: Screen) => void }) {
  return (
    <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3"><h1 className="text-lg font-bold">👤 حسابي</h1></header>
      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        <div className="bg-white rounded-xl border p-5 text-center">
          <div className="w-20 h-20 mx-auto rounded-full bg-orange-100 flex items-center justify-center text-3xl mb-3">🧑</div>
          <h2 className="text-xl font-bold text-gray-900">محمود سعيد</h2>
          <div className="inline-block mt-2 px-3 py-1 bg-amber-100 text-amber-700 rounded-full text-sm font-semibold">🥇 Platinum</div>
          <p className="text-sm text-gray-400 mt-2" dir="ltr">010 1234 5678</p>
          <p className="text-sm text-gray-400">🛵 موتوسيكل - لوحة 1234</p>
        </div>
        <div className="grid grid-cols-3 gap-3">
          <div className="bg-white rounded-xl border p-3 text-center"><p className="text-xl font-bold text-orange-500">247</p><p className="text-xs text-gray-400">توصيلات</p></div>
          <div className="bg-white rounded-xl border p-3 text-center"><p className="text-xl font-bold text-orange-500">8,420</p><p className="text-xs text-gray-400">أرباح</p></div>
          <div className="bg-white rounded-xl border p-3 text-center"><p className="text-xl font-bold text-orange-500">4.92</p><p className="text-xs text-gray-400">⭐ تقييم</p></div>
        </div>
        <div className="bg-white rounded-xl border p-4 space-y-3">
          {[{ icon: '⚙️', label: 'الملاحة (Google Maps)' }, { icon: '🔔', label: 'الإشعارات' }, { icon: '🌍', label: 'اللغة' }, { icon: '❓', label: 'مساعدة' }, { icon: '🚪', label: 'تسجيل الخروج', danger: true }].map((item, i) => (
            <button key={i} className={`w-full flex items-center gap-3 py-2 text-sm ${item.danger ? 'text-red-500' : 'text-gray-700'}`}>
              <span className="text-lg">{item.icon}</span><span>{item.label}</span>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}

// ============ Heat Map Screen ============
function HeatMapScreen({ onBack }: { onBack: () => void }) {
  return (
    <div className="min-h-screen bg-gray-50" dir="rtl">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-500"><span className="text-xl">→</span></button>
        <h1 className="text-lg font-bold">🔥 مناطق الطلب العالي</h1>
      </header>
      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        <div className="bg-gradient-to-br from-red-100 via-yellow-50 to-green-100 rounded-xl h-80 relative overflow-hidden">
          <div className="absolute top-4 right-4 text-center"><div className="w-16 h-16 bg-red-500 rounded-full opacity-30 flex items-center justify-center text-2xl">🔴</div><p className="text-xs font-bold text-red-600 mt-1">معادي</p><p className="text-xs text-red-400">3 طلبات</p></div>
          <div className="absolute top-20 left-12 text-center"><div className="w-20 h-20 bg-red-500 rounded-full opacity-30 flex items-center justify-center text-2xl">🔴</div><p className="text-xs font-bold text-red-600 mt-1">الزمالك</p><p className="text-xs text-red-400">5 طلبات</p></div>
          <div className="absolute bottom-20 right-12 text-center"><div className="w-14 h-14 bg-yellow-500 rounded-full opacity-30 flex items-center justify-center text-xl">🟡</div><p className="text-xs font-bold text-yellow-600 mt-1">مدينة نصر</p><p className="text-xs text-yellow-400">2 طلبات</p></div>
          <div className="absolute bottom-8 left-8 text-center"><div className="w-12 h-12 bg-green-500 rounded-full opacity-30 flex items-center justify-center text-lg">🟢</div><p className="text-xs font-bold text-green-600 mt-1">التحرير</p><p className="text-xs text-green-400">هادي</p></div>
        </div>
        <div className="flex gap-4 justify-center text-sm">
          <span className="flex items-center gap-1"><span className="w-3 h-3 bg-red-500 rounded-full" /> طلب عالي</span>
          <span className="flex items-center gap-1"><span className="w-3 h-3 bg-yellow-500 rounded-full" /> متوسط</span>
          <span className="flex items-center gap-1"><span className="w-3 h-3 bg-green-500 rounded-full" /> هادي</span>
        </div>
        <p className="text-xs text-gray-400 text-center">آخر تحديث: 14:32 — يتم التحديث كل 2 دقيقة</p>
      </div>
    </div>
  )
}

// ============ Bottom Nav ============
function BottomNav({ active, onNavigate }: { active: string; onNavigate: (s: Screen) => void }) {
  const items = [
    { id: 'home', label: 'الرئيسية', icon: '🏠' },
    { id: 'earnings', label: 'الأرباح', icon: '💰' },
    { id: 'history', label: 'السجل', icon: '📦' },
    { id: 'profile', label: 'حسابي', icon: '👤' },
  ]
  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 bg-white border-t h-16 flex items-center justify-around px-2">
      {items.map((item) => (
        <button key={item.id} onClick={() => onNavigate(item.id as Screen)} className={`flex flex-col items-center gap-1 flex-1 h-full justify-center ${active === item.id ? 'text-orange-500' : 'text-gray-400'}`}>
          <span className="text-xl">{item.icon}</span>
          <span className="text-xs font-medium">{item.label}</span>
        </button>
      ))}
    </nav>
  )
}
