'use client'

import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Progress } from '@/components/ui/progress'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { ChevronLeft, MapPin, Search, Star, Clock, ShoppingBag, Plus, Minus, X, Phone, MessageCircle, CheckCircle2, Bike, Home, User, Receipt } from 'lucide-react'

// ============ Types ============
type Screen = 'login' | 'otp' | 'home' | 'restaurant' | 'cart' | 'checkout' | 'tracking' | 'orders' | 'profile'
type PaymentMethod = 'vodafone_cash' | 'instapay' | 'card' | 'cod'

interface Restaurant {
  id: string
  name: string
  cuisine: string
  rating: number
  ratingCount: number
  etaMin: number
  etaMax: number
  deliveryFee: number
  priceRange: number
  isOpen: boolean
  promo?: string
  image: string
  distance: number
}

interface MenuItem {
  id: string
  name: string
  description: string
  price: number
  image: string
  isAvailable: boolean
  prepTime: number
  rating?: number
  isPopular?: boolean
  category: string
  modifiers?: Modifier[]
}

interface Modifier {
  id: string
  name: string
  required: boolean
  multiple: boolean
  options: { id: string; name: string; priceDelta: number }[]
}

interface CartItem {
  id: string
  menuItemId: string
  name: string
  price: number
  quantity: number
  modifiers: string[]
  notes?: string
}

// ============ Mock Data ============
const MOCK_RESTAURANTS: Restaurant[] = [
  { id: '1', name: 'Pizza Hut', cuisine: 'إيطالي • بيتزا', rating: 4.6, ratingCount: 1243, etaMin: 30, etaMax: 40, deliveryFee: 25, priceRange: 2, isOpen: true, promo: 'خصم 20%', image: '🍕', distance: 1.2 },
  { id: '2', name: 'McDonald\'s', cuisine: 'وجبات سريعة • برجر', rating: 4.5, ratingCount: 2150, etaMin: 25, etaMax: 35, deliveryFee: 20, priceRange: 2, isOpen: true, image: '🍔', distance: 0.8 },
  { id: '3', name: 'KFC', cuisine: 'وجبات سريعة • دجاج', rating: 4.4, ratingCount: 1890, etaMin: 28, etaMax: 38, deliveryFee: 20, priceRange: 2, isOpen: true, promo: 'وجبة مجانية', image: '🍗', distance: 1.5 },
  { id: '4', name: 'Sushi House', cuisine: 'آسيوي • سوشي', rating: 4.8, ratingCount: 567, etaMin: 35, etaMax: 45, deliveryFee: 30, priceRange: 4, isOpen: true, image: '🍣', distance: 2.3 },
  { id: '5', name: 'كشري أبوطارق', cuisine: 'مصري • كشري', rating: 4.9, ratingCount: 3200, etaMin: 20, etaMax: 30, deliveryFee: 15, priceRange: 1, isOpen: true, promo: 'خصم 15%', image: '🍜', distance: 0.5 },
  { id: '6', name: 'فول و طعمية الحاج محمود', cuisine: 'مصري • فطار', rating: 4.7, ratingCount: 890, etaMin: 15, etaMax: 25, deliveryFee: 10, priceRange: 1, isOpen: true, image: '🫘', distance: 0.3 },
  { id: '7', name: 'Domino\'s Pizza', cuisine: 'إيطالي • بيتزا', rating: 4.3, ratingCount: 1100, etaMin: 32, etaMax: 42, deliveryFee: 25, priceRange: 2, isOpen: false, image: '🍕', distance: 1.8 },
  { id: '8', name: 'Starbucks', cuisine: 'قهوة • مشروبات', rating: 4.5, ratingCount: 2300, etaMin: 18, etaMax: 28, deliveryFee: 18, priceRange: 3, isOpen: true, image: '☕', distance: 0.7 },
]

const MOCK_MENU: Record<string, MenuItem[]> = {
  '1': [
    { id: 'm1', name: 'Margherita Pizza', description: 'صلصة طماطم، موتزاريلا طازجة، ريحان', price: 145, image: '🍕', isAvailable: true, prepTime: 12, rating: 4.7, isPopular: true, category: 'بيتزا', modifiers: [
      { id: 'mod1', name: 'الحجم', required: true, multiple: false, options: [
        { id: 's', name: 'صغير', priceDelta: -25 }, { id: 'm', name: 'وسط', priceDelta: 0 }, { id: 'l', name: 'كبير', priceDelta: 30 }
      ]},
      { id: 'mod2', name: 'إضافات', required: false, multiple: true, options: [
        { id: 'cheese', name: 'جبن إضافي', priceDelta: 15 }, { id: 'bacon', name: 'بيكون', priceDelta: 20 }, { id: 'mushroom', name: 'فطر', priceDelta: 10 }
      ]}
    ]},
    { id: 'm2', name: 'Pepperoni Pizza', description: 'بسطرمة لحم، جبن إضافي', price: 165, image: '🍕', isAvailable: true, prepTime: 12, rating: 4.8, isPopular: true, category: 'بيتزا' },
    { id: 'm3', name: 'Veggie Supreme', description: 'خضار مشكل، زيتون، فطر', price: 155, image: '🥗', isAvailable: true, prepTime: 10, rating: 4.5, category: 'بيتزا' },
    { id: 'm4', name: 'Coca Cola', description: '330ml', price: 25, image: '🥤', isAvailable: true, prepTime: 1, category: 'مشروبات' },
    { id: 'm5', name: 'Apple Pie', description: 'فطيرة تفاح دافئة', price: 35, image: '🍰', isAvailable: false, prepTime: 5, category: 'حلويات' },
  ],
  '2': [
    { id: 'm6', name: 'Big Mac Meal', description: 'برجر لحم بقري، بطاطس، كوكا', price: 85, image: '🍔', isAvailable: true, prepTime: 8, rating: 4.6, isPopular: true, category: 'وجبات' },
    { id: 'm7', name: 'McChicken', description: 'ساندويتش دجاج مقرمش', price: 90, image: '🍔', isAvailable: true, prepTime: 8, rating: 4.5, category: 'وجبات' },
    { id: 'm8', name: 'Apple Pie', description: 'فطيرة تفاح', price: 25, image: '🍰', isAvailable: true, prepTime: 3, category: 'حلويات' },
  ],
}

const CUISINES = [
  { name: 'مصري', icon: '🍽️', color: 'bg-orange-100' },
  { name: 'إيطالي', icon: '🍕', color: 'bg-red-100' },
  { name: 'آسيوي', icon: '🍜', color: 'bg-yellow-100' },
  { name: 'وجبات سريعة', icon: '🍔', color: 'bg-amber-100' },
  { name: 'صحي', icon: '🥗', color: 'bg-green-100' },
  { name: 'حلويات', icon: '🍰', color: 'bg-pink-100' },
  { name: 'قهوة', icon: '☕', color: 'bg-brown-100' },
  { name: 'فطار', icon: '🍳', color: 'bg-yellow-100' },
]

const ORDER_STEPS = [
  { status: 'confirmed', label: 'تم تأكيد الطلب', icon: CheckCircle2 },
  { status: 'preparing', label: 'المطعم بيحضّر', icon: Clock },
  { status: 'ready', label: 'الطلب جاهز', icon: ShoppingBag },
  { status: 'picked_up', label: 'المندوب في الطريق', icon: Bike },
  { status: 'delivered', label: 'تم التوصيل', icon: Home },
]

// ============ Main Page ============
export default function CustomerWebApp() {
  const [screen, setScreen] = useState<Screen>('login')
  const [phone, setPhone] = useState('')
  const [otp, setOtp] = useState(['', '', '', '', '', ''])
  const [selectedRestaurant, setSelectedRestaurant] = useState<Restaurant | null>(null)
  const [cart, setCart] = useState<CartItem[]>([])
  const [selectedItem, setSelectedItem] = useState<MenuItem | null>(null)
  const [orderStep, setOrderStep] = useState(0)
  const [showItemModal, setShowItemModal] = useState(false)

  // Simulate order progression
  useEffect(() => {
    if (screen === 'tracking' && orderStep < 4) {
      const timer = setTimeout(() => setOrderStep(orderStep + 1), 3000)
      return () => clearTimeout(timer)
    }
  }, [screen, orderStep])

  const cartTotal = cart.reduce((sum, item) => sum + item.price * item.quantity, 0)
  const deliveryFee = 25
  const serviceFee = cartTotal * 0.05
  const vat = (cartTotal + deliveryFee + serviceFee) * 0.14
  const grandTotal = cartTotal + deliveryFee + serviceFee + vat

  // ============ Login Screen ============
  if (screen === 'login') {
    return (
      <div className="min-h-screen bg-gradient-to-b from-orange-50 to-white flex flex-col items-center justify-center p-6" dir="rtl">
        <div className="w-full max-w-md space-y-6">
          <div className="text-center">
            <div className="inline-flex items-center justify-center w-20 h-20 bg-orange-500 rounded-3xl mb-4 text-4xl">
              🍔
            </div>
            <h1 className="text-3xl font-bold text-gray-900">أهلاً! 👋</h1>
            <p className="text-gray-500 mt-2">سجّل رقم موبايلك عشان نبعتلك كود تفعيل</p>
          </div>

          <div className="space-y-4">
            <div>
              <Label className="text-sm text-gray-600 mb-1.5 block">رقم الموبايل</Label>
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1 px-3 h-12 rounded-lg border border-gray-200 bg-gray-50 text-sm font-medium">
                  🇪🇬 +20
                </div>
                <Input
                  type="tel"
                  inputMode="numeric"
                  maxLength={11}
                  placeholder="01 2345 6789"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value.replace(/\D/g, ''))}
                  className="h-12 text-lg"
                />
              </div>
            </div>

            <Button
              className="w-full h-12 text-lg bg-orange-500 hover:bg-orange-600"
              disabled={phone.length < 11}
              onClick={() => setScreen('otp')}
            >
              أرسل الكود
            </Button>
          </div>

          <p className="text-xs text-gray-400 text-center">
            بتسجيلك، أنت موافق على <span className="text-orange-500">الشروط والأحكام</span> و<span className="text-orange-500">سياسة الخصوصية</span>
          </p>
        </div>
      </div>
    )
  }

  // ============ OTP Screen ============
  if (screen === 'otp') {
    return (
      <div className="min-h-screen bg-white flex flex-col items-center justify-center p-6" dir="rtl">
        <div className="w-full max-w-md space-y-6">
          <button onClick={() => setScreen('login')} className="flex items-center gap-1 text-gray-500 hover:text-gray-900">
            <ChevronLeft className="w-4 h-4 rotate-180" />
            <span>رجوع</span>
          </button>

          <div className="text-center">
            <h1 className="text-2xl font-bold text-gray-900">أدخل الكود</h1>
            <p className="text-gray-500 mt-2">بعتناك كود على</p>
            <p className="font-semibold text-gray-900" dir="ltr">+20 {phone}</p>
          </div>

          <div className="flex justify-center gap-2" dir="ltr">
            {otp.map((digit, i) => (
              <input
                key={i}
                type="text"
                inputMode="numeric"
                maxLength={1}
                value={digit}
                onChange={(e) => {
                  const newOtp = [...otp]
                  newOtp[i] = e.target.value.replace(/\D/g, '')
                  setOtp(newOtp)
                  if (newOtp[i] && i < 5) {
                    const next = document.getElementById(`otp-${i + 1}`)
                    next?.focus()
                  }
                  if (newOtp.every(d => d) && newOtp.join('').length === 6) {
                    setTimeout(() => setScreen('home'), 500)
                  }
                }}
                id={`otp-${i}`}
                className={`w-12 h-14 text-center text-2xl font-bold rounded-lg border-2 transition-colors
                  ${digit ? 'border-orange-500 bg-orange-50' : 'border-gray-200'}`}
                autoFocus={i === 0}
              />
            ))}
          </div>

          <p className="text-center text-sm text-gray-400">
            <Clock className="inline w-4 h-4 ml-1" />
            02:00 متبقي
          </p>

          <Button
            className="w-full h-12 bg-orange-500 hover:bg-orange-600"
            disabled={otp.some(d => !d)}
            onClick={() => setScreen('home')}
          >
            تأكيد
          </Button>
        </div>
      </div>
    )
  }

  // ============ Home Screen ============
  if (screen === 'home') {
    return (
      <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
        {/* Header */}
        <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
          <div className="flex items-center gap-2 flex-1">
            <MapPin className="w-5 h-5 text-orange-500" />
            <div>
              <p className="text-xs text-gray-400">التوصيل إلى</p>
              <p className="text-sm font-semibold text-gray-900">الزمالك، القاهرة</p>
            </div>
          </div>
          <button onClick={() => setScreen('orders')} className="w-10 h-10 rounded-full hover:bg-gray-100 flex items-center justify-center">
            <Search className="w-5 h-5 text-gray-500" />
          </button>
          <button onClick={() => setScreen('profile')} className="w-10 h-10 rounded-full bg-orange-100 hover:bg-orange-200 flex items-center justify-center">
            <User className="w-5 h-5 text-orange-600" />
          </button>
        </header>

        <div className="max-w-4xl mx-auto px-4 py-4 space-y-6">
          {/* Welcome Banner */}
          <div className="bg-gradient-to-l from-orange-500 to-red-500 rounded-2xl p-5 text-white flex items-center justify-between">
            <div>
              <h2 className="text-lg font-bold">🎁 خصم 50% على أول طلب!</h2>
              <p className="text-sm opacity-90 mt-1">استخدم الكود: WELCOME50</p>
            </div>
            <span className="text-4xl">🎉</span>
          </div>

          {/* Trending */}
          <section>
            <h2 className="text-lg font-bold text-gray-900 mb-3 flex items-center gap-2">
              🔥 رائج قريب منك
            </h2>
            <ScrollArea className="w-full">
              <div className="flex gap-4 pb-2">
                {MOCK_RESTAURANTS.filter(r => r.isOpen).map((r) => (
                  <RestaurantCard key={r.id} restaurant={r} onClick={() => { setSelectedRestaurant(r); setScreen('restaurant') }} />
                ))}
              </div>
            </ScrollArea>
          </section>

          {/* Cuisines */}
          <section>
            <h2 className="text-lg font-bold text-gray-900 mb-3">🍽️ أنواع المطابخ</h2>
            <div className="grid grid-cols-4 gap-3">
              {CUISINES.map((c) => (
                <button key={c.name} className="flex flex-col items-center gap-2 p-3 bg-white rounded-xl border hover:shadow-md transition-shadow">
                  <span className="text-2xl">{c.icon}</span>
                  <span className="text-xs font-medium text-gray-700">{c.name}</span>
                </button>
              ))}
            </div>
          </section>

          {/* Top Rated */}
          <section>
            <h2 className="text-lg font-bold text-gray-900 mb-3 flex items-center gap-2">
              🏆 الأعلى تقييماً
            </h2>
            <div className="space-y-3">
              {[...MOCK_RESTAURANTS].sort((a, b) => b.rating - a.rating).slice(0, 3).map((r) => (
                <RestaurantCardHorizontal key={r.id} restaurant={r} onClick={() => { setSelectedRestaurant(r); setScreen('restaurant') }} />
              ))}
            </div>
          </section>
        </div>

        {/* Bottom Nav */}
        <BottomNav active="home" onNavigate={setScreen} cartCount={cart.length} />
      </div>
    )
  }

  // ============ Restaurant Detail Screen ============
  if (screen === 'restaurant' && selectedRestaurant) {
    const menuItems = MOCK_MENU[selectedRestaurant.id] ?? MOCK_MENU['1']
    const categories = [...new Set(menuItems.map(m => m.category))]

    return (
      <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
        {/* Cover */}
        <div className="relative w-full h-48 bg-gradient-to-br from-orange-200 to-red-200 flex items-center justify-center">
          <span className="text-7xl">{selectedRestaurant.image}</span>
          <button onClick={() => setScreen('home')} className="absolute top-4 right-4 w-10 h-10 bg-white/90 rounded-full flex items-center justify-center shadow-md hover:bg-white">
            <ChevronLeft className="w-5 h-5 rotate-180" />
          </button>
        </div>

        {/* Info */}
        <div className="max-w-4xl mx-auto px-4 -mt-6 relative">
          <Card className="shadow-lg">
            <CardContent className="p-5">
              <h1 className="text-2xl font-bold text-gray-900">{selectedRestaurant.name}</h1>
              <p className="text-sm text-gray-500 mt-1">{selectedRestaurant.cuisine}</p>
              <div className="flex items-center gap-4 mt-3 text-sm">
                <span className="flex items-center gap-1">
                  <Star className="w-4 h-4 text-amber-400 fill-amber-400" />
                  <span className="font-semibold">{selectedRestaurant.rating}</span>
                  <span className="text-gray-400">({selectedRestaurant.ratingCount})</span>
                </span>
                <span className="flex items-center gap-1 text-gray-500">
                  <Clock className="w-4 h-4" />
                  {selectedRestaurant.etaMin}-{selectedRestaurant.etaMax} دقيقة
                </span>
                <span className="text-gray-500">EGP {selectedRestaurant.deliveryFee} توصيل</span>
              </div>
              {selectedRestaurant.promo && (
                <Badge className="mt-3 bg-orange-100 text-orange-700 hover:bg-orange-100">🎁 {selectedRestaurant.promo}</Badge>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Menu */}
        <div className="max-w-4xl mx-auto px-4 mt-6 space-y-6">
          {categories.map((cat) => (
            <div key={cat}>
              <h2 className="text-lg font-bold text-gray-900 mb-3 sticky top-0 bg-gray-50 py-2">{cat}</h2>
              <div className="space-y-3">
                {menuItems.filter(m => m.category === cat).map((item) => (
                  <Card key={item.id} className={`overflow-hidden ${!item.isAvailable ? 'opacity-50' : 'cursor-pointer hover:shadow-md'}`} onClick={() => { if (item.isAvailable) { setSelectedItem(item); setShowItemModal(true) } }}>
                    <CardContent className="p-3 flex gap-3">
                      <div className="flex-1">
                        <div className="flex items-center gap-2">
                          <h3 className="font-semibold text-gray-900">{item.name}</h3>
                          {item.isPopular && <Badge className="bg-amber-100 text-amber-700 hover:bg-amber-100 text-xs">🔥 الأكثر طلباً</Badge>}
                        </div>
                        <p className="text-xs text-gray-500 mt-1 line-clamp-2">{item.description}</p>
                        <div className="flex items-center gap-3 mt-2">
                          <span className="font-bold text-gray-900">EGP {item.price}</span>
                          {item.rating && <span className="text-xs text-gray-400 flex items-center gap-0.5"><Star className="w-3 h-3 text-amber-400 fill-amber-400" />{item.rating}</span>}
                          <span className="text-xs text-gray-400">{item.prepTime} د</span>
                        </div>
                        {!item.isAvailable && <Badge variant="destructive" className="mt-2 text-xs">نفد</Badge>}
                      </div>
                      <div className="relative w-24 h-24 bg-gray-100 rounded-lg flex items-center justify-center text-3xl flex-shrink-0">
                        {item.image}
                        {item.isAvailable && (
                          <button
                            className="absolute -bottom-2 -left-2 w-8 h-8 bg-orange-500 text-white rounded-full flex items-center justify-center shadow-md hover:bg-orange-600"
                            onClick={(e) => { e.stopPropagation(); setSelectedItem(item); setShowItemModal(true) }}
                          >
                            <Plus className="w-4 h-4" />
                          </button>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          ))}
        </div>

        {/* Cart bar */}
        {cart.length > 0 && (
          <div className="fixed bottom-4 left-4 right-4 z-40">
            <Button className="w-full h-14 bg-orange-500 hover:bg-orange-600 text-white shadow-xl rounded-xl flex items-center justify-between px-6" onClick={() => setScreen('cart')}>
              <span className="flex items-center gap-2">
                <ShoppingBag className="w-5 h-5" />
                {cart.length} صنف
              </span>
              <span>EGP {cartTotal.toFixed(2)} • عرض السلة</span>
            </Button>
          </div>
        )}

        {/* Item Modal */}
        {showItemModal && selectedItem && (
          <ItemModal item={selectedItem} onClose={() => setShowItemModal(false)} onAdd={(qty, mods, notes) => {
            setCart([...cart, { id: Date.now().toString(), menuItemId: selectedItem.id, name: selectedItem.name, price: selectedItem.price, quantity: qty, modifiers: mods, notes }])
            setShowItemModal(false)
          }} />
        )}
      </div>
    )
  }

  // ============ Cart Screen ============
  if (screen === 'cart') {
    return (
      <div className="min-h-screen bg-gray-50 pb-24" dir="rtl">
        <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
          <button onClick={() => setScreen('restaurant')} className="flex items-center gap-1 text-gray-500 hover:text-gray-900">
            <ChevronLeft className="w-5 h-5 rotate-180" />
            <span>رجوع</span>
          </button>
          <h1 className="text-lg font-bold text-gray-900">السلة</h1>
        </header>

        <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
          {cart.length === 0 ? (
            <div className="text-center py-20">
              <ShoppingBag className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-500">السلة فاضية</p>
              <Button className="mt-4 bg-orange-500 hover:bg-orange-600" onClick={() => setScreen('home')}>تصفح المطاعم</Button>
            </div>
          ) : (
            <>
              {/* Items */}
              <Card>
                <CardHeader><CardTitle>{selectedRestaurant?.name ?? 'Pizza Hut'}</CardTitle></CardHeader>
                <CardContent className="space-y-3">
                  {cart.map((item) => (
                    <div key={item.id} className="flex gap-3 pb-3 border-b last:border-0 last:pb-0">
                      <div className="w-14 h-14 bg-gray-100 rounded-lg flex items-center justify-center text-2xl">🍕</div>
                      <div className="flex-1">
                        <h3 className="text-sm font-semibold">{item.name}</h3>
                        {item.modifiers.map((m, i) => <p key={i} className="text-xs text-gray-400">+ {m}</p>)}
                        {item.notes && <p className="text-xs text-gray-400 italic">{item.notes}</p>}
                        <div className="flex items-center justify-between mt-2">
                          <div className="flex items-center gap-2">
                            <button className="w-7 h-7 rounded-full border flex items-center justify-center" onClick={() => setCart(cart.map(c => c.id === item.id ? { ...c, quantity: Math.max(1, c.quantity - 1) } : c))}><Minus className="w-3 h-3" /></button>
                            <span className="font-semibold w-6 text-center">{item.quantity}</span>
                            <button className="w-7 h-7 rounded-full border flex items-center justify-center" onClick={() => setCart(cart.map(c => c.id === item.id ? { ...c, quantity: c.quantity + 1 } : c))}><Plus className="w-3 h-3" /></button>
                          </div>
                          <span className="font-bold">EGP {(item.price * item.quantity).toFixed(2)}</span>
                          <button onClick={() => setCart(cart.filter(c => c.id !== item.id))} className="text-red-400 hover:text-red-600"><X className="w-4 h-4" /></button>
                        </div>
                      </div>
                    </div>
                  ))}
                </CardContent>
              </Card>

              {/* Coupon */}
              <Card>
                <CardContent className="p-4 flex items-center gap-2">
                  <span className="text-lg">🎁</span>
                  <Input placeholder="WELCOME50" className="flex-1" />
                  <Button variant="outline">تطبيق</Button>
                </CardContent>
              </Card>

              {/* Price breakdown */}
              <Card>
                <CardContent className="p-4 space-y-2 text-sm">
                  <div className="flex justify-between"><span className="text-gray-500">Subtotal</span><span>EGP {cartTotal.toFixed(2)}</span></div>
                  <div className="flex justify-between"><span className="text-gray-500">رسوم التوصيل</span><span>EGP {deliveryFee.toFixed(2)}</span></div>
                  <div className="flex justify-between"><span className="text-gray-500">رسوم الخدمة (5%)</span><span>EGP {serviceFee.toFixed(2)}</span></div>
                  <div className="flex justify-between"><span className="text-gray-500">ض.ق.م (14%)</span><span>EGP {vat.toFixed(2)}</span></div>
                  <Separator className="my-2" />
                  <div className="flex justify-between font-bold text-lg"><span>الإجمالي</span><span className="text-orange-500">EGP {grandTotal.toFixed(2)}</span></div>
                  <div className="bg-green-50 border border-green-200 rounded-lg p-3 mt-2 flex items-center gap-2">
                    <span>💰</span>
                    <p className="text-xs text-green-700">كاش باك: EGP {(grandTotal * 0.05).toFixed(2)} (5%)</p>
                  </div>
                </CardContent>
              </Card>
            </>
          )}
        </div>

        {cart.length > 0 && (
          <div className="fixed bottom-0 left-0 right-0 bg-white border-t p-4">
            <div className="max-w-2xl mx-auto">
              <Button className="w-full h-12 bg-orange-500 hover:bg-orange-600" onClick={() => setScreen('checkout')}>🛒 اطلب الآن — EGP {grandTotal.toFixed(2)}</Button>
            </div>
          </div>
        )}
      </div>
    )
  }

  // ============ Checkout Screen ============
  if (screen === 'checkout') {
    const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('vodafone_cash')
    return (
      <div className="min-h-screen bg-gray-50 pb-24" dir="rtl">
        <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
          <button onClick={() => setScreen('cart')} className="flex items-center gap-1 text-gray-500 hover:text-gray-900">
            <ChevronLeft className="w-5 h-5 rotate-180" /><span>رجوع</span>
          </button>
          <h1 className="text-lg font-bold">إتمام الطلب</h1>
        </header>

        <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
          <Section title="📍 عنوان التوصيل">
            <Card><CardContent className="p-4">
              <p className="font-medium">الزمالك، 26 يوليو، شقة 5</p>
              <p className="text-xs text-gray-400 mt-1">علامة: باب أزرق</p>
            </CardContent></Card>
          </Section>

          <Section title="⏰ وقت التوصيل">
            <Card><CardContent className="p-4">
              <RadioGroup defaultValue="asap">
                <div className="flex items-center gap-3 mb-2"><RadioGroupItem value="asap" id="asap" /><Label htmlFor="asap">في أسرع وقت (35-45 دقيقة)</Label></div>
                <div className="flex items-center gap-3"><RadioGroupItem value="scheduled" id="scheduled" /><Label htmlFor="scheduled">مجدول لوقت لاحق</Label></div>
              </RadioGroup>
            </CardContent></Card>
          </Section>

          <Section title="💳 طريقة الدفع">
            <Card><CardContent className="p-4 space-y-3">
              {[
                { id: 'vodafone_cash', label: 'Vodafone Cash', icon: '💚', desc: 'رصيدك: EGP 1,250' },
                { id: 'instapay', label: 'InstaPay', icon: '🟣', desc: 'تحويل بنكي فوري' },
                { id: 'card', label: 'بطاقة بنكية', icon: '💳', desc: '**** 4521' },
                { id: 'cod', label: 'الدفع عند الاستلام', icon: '💵', desc: '+EGP 5 رسوم' },
              ].map((opt) => (
                <label key={opt.id} className={`flex items-center gap-3 p-3 rounded-lg border-2 cursor-pointer transition-colors ${paymentMethod === opt.id ? 'border-orange-500 bg-orange-50' : 'border-gray-200'}`}>
                  <RadioGroupItem value={opt.id} checked={paymentMethod === opt.id} onClick={() => setPaymentMethod(opt.id as PaymentMethod)} />
                  <span className="text-xl">{opt.icon}</span>
                  <div className="flex-1"><p className="font-semibold text-sm">{opt.label}</p><p className="text-xs text-gray-400">{opt.desc}</p></div>
                </label>
              ))}
            </CardContent></Card>
          </Section>

          <Section title="📝 ملاحظات للمطعم">
            <Input placeholder="مثلاً: صلصة إضافية لو سمحت" />
          </Section>
        </div>

        <div className="fixed bottom-0 left-0 right-0 bg-white border-t p-4">
          <div className="max-w-2xl mx-auto">
            <Button className="w-full h-12 bg-orange-500 hover:bg-orange-600" onClick={() => { setOrderStep(0); setScreen('tracking') }}>🔒 تأكيد الطلب — EGP {grandTotal.toFixed(2)}</Button>
          </div>
        </div>
      </div>
    )
  }

  // ============ Order Tracking Screen ============
  if (screen === 'tracking') {
    const progress = ((orderStep + 1) / 5) * 100
    return (
      <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
        <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center">
          <button onClick={() => setScreen('home')} className="flex items-center gap-1 text-gray-500 hover:text-gray-900">
            <ChevronLeft className="w-5 h-5 rotate-180" /><span>رجوع</span>
          </button>
          <h1 className="text-lg font-bold absolute left-1/2 -translate-x-1/2">الطلب #A7X92F</h1>
        </header>

        <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
          {/* Celebration */}
          {orderStep === 0 && (
            <div className="bg-green-50 border-b border-green-200 p-6 text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-green-500 rounded-full mb-2">
                <CheckCircle2 className="w-8 h-8 text-white" />
              </div>
              <h2 className="text-xl font-bold text-gray-900">تم الطلب! 🎉</h2>
              <p className="text-sm text-gray-500 mt-1">طلبك اتأكد والمطعم هيبدأ تحضيره</p>
            </div>
          )}

          {/* ETA */}
          <div className="bg-orange-50 border-b border-orange-200 p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Clock className="w-5 h-5 text-orange-500" />
                <div><p className="text-xs text-gray-400">الوقت المتوقع</p><p className="text-xl font-bold text-orange-500">8:35 PM</p></div>
              </div>
              <div className="text-left"><p className="text-xs text-gray-400">الإجمالي</p><p className="text-xl font-bold">EGP {grandTotal.toFixed(2)}</p></div>
            </div>
          </div>

          {/* Progress bar */}
          <Card>
            <CardContent className="p-5">
              <Progress value={progress} className="h-2 mb-4" />
              <div className="space-y-4">
                {ORDER_STEPS.map((step, i) => {
                  const Icon = step.icon
                  const isDone = i < orderStep
                  const isCurrent = i === orderStep
                  return (
                    <div key={step.status} className="flex items-center gap-3">
                      <div className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 transition-all ${isDone ? 'bg-green-500 text-white' : isCurrent ? 'bg-orange-500 text-white animate-pulse' : 'bg-gray-100 text-gray-300'}`}>
                        <Icon className="w-5 h-5" />
                      </div>
                      <p className={`text-sm font-medium ${isDone || isCurrent ? 'text-gray-900' : 'text-gray-300'}`}>{step.label}</p>
                      {isCurrent && <span className="text-xs text-orange-500">جاري...</span>}
                    </div>
                  )
                })}
              </div>
            </CardContent>
          </Card>

          {/* Driver */}
          {orderStep >= 3 && (
            <Card>
              <CardContent className="p-4">
                <div className="flex items-center gap-3">
                  <Avatar className="w-12 h-12 bg-orange-100"><AvatarFallback className="text-orange-600">MS</AvatarFallback></Avatar>
                  <div className="flex-1">
                    <p className="font-semibold">Mahmoud S.</p>
                    <div className="flex items-center gap-1 text-xs text-gray-400"><Star className="w-3 h-3 text-amber-400 fill-amber-400" />4.8 • موتوسيكل</div>
                  </div>
                  <button className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center"><Phone className="w-5 h-5 text-green-600" /></button>
                  <button className="w-10 h-10 rounded-full bg-orange-100 flex items-center justify-center"><MessageCircle className="w-5 h-5 text-orange-600" /></button>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Order details */}
          <Card>
            <CardContent className="p-4">
              <h3 className="font-bold mb-3">تفاصيل الطلب</h3>
              {cart.map((item, i) => (
                <div key={i} className="flex justify-between text-sm py-1">
                  <span>{item.quantity}× {item.name}</span>
                  <span className="text-gray-400">EGP {(item.price * item.quantity).toFixed(2)}</span>
                </div>
              ))}
              <Separator className="my-2" />
              <div className="flex justify-between font-bold"><span>الإجمالي</span><span className="text-orange-500">EGP {grandTotal.toFixed(2)}</span></div>
            </CardContent>
          </Card>

          {orderStep === 4 && (
            <Card className="text-center">
              <CardContent className="p-4">
                <h3 className="font-bold mb-2">قيّم طلبك 👇</h3>
                <div className="flex justify-center gap-2 my-3">
                  {[1,2,3,4,5].map(s => <button key={s} className="text-3xl text-amber-400 hover:scale-110 transition-transform">★</button>)}
                </div>
                <Button className="w-full bg-orange-500 hover:bg-orange-600" onClick={() => setScreen('home')}>العودة للرئيسية</Button>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    )
  }

  // ============ Orders History ============
  if (screen === 'orders') {
    return (
      <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
        <header className="sticky top-0 z-40 bg-white border-b px-4 py-3"><h1 className="text-lg font-bold">طلباتي</h1></header>
        <div className="max-w-2xl mx-auto px-4 py-4">
          <Card className="mb-3"><CardContent className="p-4">
            <div className="flex justify-between items-start mb-2">
              <div><p className="font-semibold">Pizza Hut</p><p className="text-xs text-gray-400">أمس • 9:20 PM</p></div>
              <Badge className="bg-green-100 text-green-700 hover:bg-green-100">تم التوصيل</Badge>
            </div>
            <p className="text-sm text-gray-500">2× Margherita + 1× Pepperoni</p>
            <div className="flex justify-between items-center mt-2">
              <span className="font-bold">EGP 290.00</span>
              <Button variant="outline" size="sm" onClick={() => { setOrderStep(0); setScreen('tracking') }}>تتبع الطلب</Button>
            </div>
          </CardContent></Card>
          <Card><CardContent className="p-4">
            <div className="flex justify-between items-start mb-2">
              <div><p className="font-semibold">McDonald's</p><p className="text-xs text-gray-400">قبل 3 أيام</p></div>
              <Badge className="bg-green-100 text-green-700 hover:bg-green-100">تم التوصيل</Badge>
            </div>
            <p className="text-sm text-gray-500">1× Big Mac Meal</p>
            <div className="flex justify-between items-center mt-2">
              <span className="font-bold">EGP 85.00</span>
              <Button variant="outline" size="sm">إعادة الطلب</Button>
            </div>
          </CardContent></Card>
        </div>
        <BottomNav active="orders" onNavigate={setScreen} cartCount={cart.length} />
      </div>
    )
  }

  // ============ Profile Screen ============
  if (screen === 'profile') {
    return (
      <div className="min-h-screen bg-gray-50 pb-20" dir="rtl">
        <header className="sticky top-0 z-40 bg-white border-b px-4 py-3"><h1 className="text-lg font-bold">حسابي</h1></header>
        <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
          <Card><CardContent className="p-5 text-center">
            <Avatar className="w-20 h-20 mx-auto bg-orange-100 mb-3"><AvatarFallback className="text-2xl text-orange-600">أ</AvatarFallback></Avatar>
            <h2 className="text-xl font-bold">أحمد محمد</h2>
            <Badge className="mt-2 bg-amber-100 text-amber-700 hover:bg-amber-100">🥇 Platinum Member</Badge>
            <p className="text-sm text-gray-400 mt-2">010 1234 5678</p>
          </CardContent></Card>

          <div className="grid grid-cols-3 gap-3">
            <Card><CardContent className="p-4 text-center"><p className="text-2xl font-bold text-orange-500">47</p><p className="text-xs text-gray-400">طلبات</p></CardContent></Card>
            <Card><CardContent className="p-4 text-center"><p className="text-2xl font-bold text-orange-500">8,520</p><p className="text-xs text-gray-400">إنفاق</p></CardContent></Card>
            <Card><CardContent className="p-4 text-center"><p className="text-2xl font-bold text-orange-500">425</p><p className="text-xs text-gray-400">محفظة</p></CardContent></Card>
          </div>

          <Card><CardContent className="p-4 space-y-3">
            <div className="flex items-center gap-3 text-sm"><MapPin className="w-4 h-4 text-gray-400" /><span>العناوين المحفوظة</span></div>
            <Separator />
            <div className="flex items-center gap-3 text-sm"><CreditCard className="w-4 h-4 text-gray-400" /><span>طرق الدفع</span></div>
            <Separator />
            <div className="flex items-center gap-3 text-sm"><Receipt className="w-4 h-4 text-gray-400" /><span>الفواتير</span></div>
            <Separator />
            <div className="flex items-center gap-3 text-sm text-red-500 cursor-pointer"><X className="w-4 h-4" /><span>تسجيل الخروج</span></div>
          </CardContent></Card>
        </div>
        <BottomNav active="profile" onNavigate={setScreen} cartCount={cart.length} />
      </div>
    )
  }

  return null
}

// ============ Components ============

function RestaurantCard({ restaurant, onClick }: { restaurant: Restaurant; onClick: () => void }) {
  return (
    <div onClick={onClick} className={`w-72 flex-shrink-0 bg-white rounded-xl border overflow-hidden cursor-pointer hover:shadow-lg transition-all ${!restaurant.isOpen ? 'opacity-60' : ''}`}>
      <div className="w-full h-32 bg-gradient-to-br from-orange-100 to-red-100 flex items-center justify-center relative">
        <span className="text-5xl">{restaurant.image}</span>
        {restaurant.promo && <Badge className="absolute top-2 right-2 bg-orange-500 hover:bg-orange-500 text-xs">🎁 {restaurant.promo}</Badge>}
        {!restaurant.isOpen && <div className="absolute inset-0 bg-black/50 flex items-center justify-center"><span className="text-white font-bold">مقفل</span></div>}
      </div>
      <div className="p-3">
        <h3 className="font-bold text-gray-900 truncate">{restaurant.name}</h3>
        <p className="text-xs text-gray-400 mt-0.5 truncate">{restaurant.cuisine}</p>
        <div className="flex items-center gap-3 mt-2 text-xs">
          <span className="flex items-center gap-0.5"><Star className="w-3 h-3 text-amber-400 fill-amber-400" /><span className="font-semibold">{restaurant.rating}</span><span className="text-gray-300">({restaurant.ratingCount})</span></span>
          <span className="text-gray-400">{restaurant.etaMin}-{restaurant.etaMax} د</span>
        </div>
        <div className="flex justify-between mt-2 text-xs text-gray-400">
          <span>EGP {restaurant.deliveryFee} توصيل</span>
          <span>{restaurant.distance} km</span>
        </div>
      </div>
    </div>
  )
}

function RestaurantCardHorizontal({ restaurant, onClick }: { restaurant: Restaurant; onClick: () => void }) {
  return (
    <div onClick={onClick} className="flex gap-3 p-3 bg-white rounded-xl border cursor-pointer hover:shadow-md transition-all">
      <div className="w-20 h-20 rounded-lg bg-orange-100 flex items-center justify-center text-3xl flex-shrink-0">{restaurant.image}</div>
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2">
          <h3 className="font-bold text-gray-900 truncate">{restaurant.name}</h3>
          {restaurant.promo && <Badge variant="outline" className="text-xs flex-shrink-0">{restaurant.promo}</Badge>}
        </div>
        <p className="text-xs text-gray-400 mt-0.5 truncate">{restaurant.cuisine}</p>
        <div className="flex items-center gap-3 mt-2 text-xs text-gray-400">
          <span className="flex items-center gap-0.5"><Star className="w-3 h-3 text-amber-400 fill-amber-400" />{restaurant.rating}</span>
          <span>{restaurant.distance} km</span>
          <span>{restaurant.etaMin}-{restaurant.etaMax} د</span>
        </div>
        <p className="text-xs text-gray-300 mt-1">EGP {restaurant.deliveryFee} توصيل</p>
      </div>
    </div>
  )
}

function ItemModal({ item, onClose, onAdd }: { item: MenuItem; onClose: () => void; onAdd: (qty: number, mods: string[], notes: string) => void }) {
  const [qty, setQty] = useState(1)
  const [notes, setNotes] = useState('')
  const [selectedMods, setSelectedMods] = useState<Record<string, string[]>>({})

  const toggleMod = (modId: string, optId: string, multiple: boolean) => {
    if (multiple) {
      const current = selectedMods[modId] ?? []
      setSelectedMods({ ...selectedMods, [modId]: current.includes(optId) ? current.filter(o => o !== optId) : [...current, optId] })
    } else {
      setSelectedMods({ ...selectedMods, [modId]: [optId] })
    }
  }

  const allRequired = (item.modifiers ?? []).every(m => !m.required || (selectedMods[m.id]?.length ?? 0) > 0)
  const modExtra = (item.modifiers ?? []).reduce((sum, m) => sum + (m.options.filter(o => selectedMods[m.id]?.includes(o.id)).reduce((s, o) => s + o.priceDelta, 0)), 0)
  const total = (item.price + modExtra) * qty

  const selectedModNames = Object.entries(selectedMods).flatMap(([modId, optIds]) => {
    const mod = item.modifiers?.find(m => m.id === modId)
    return optIds.map(optId => mod?.options.find(o => o.id === optId)?.name ?? '').filter(Boolean)
  })

  return (
    <Sheet open onOpenChange={onClose}>
      <SheetContent side="bottom" className="h-[85vh] rounded-t-2xl p-0" dir="rtl">
        <SheetHeader className="px-5 pt-3">
          <SheetTitle className="text-center">{item.name}</SheetTitle>
        </SheetHeader>
        <ScrollArea className="h-[calc(85vh-140px)] px-5">
          {item.image && <div className="w-full h-40 bg-gray-100 rounded-xl flex items-center justify-center text-6xl mb-4">{item.image}</div>}
          <p className="text-sm text-gray-500 mb-2">{item.description}</p>
          <p className="text-lg font-bold text-orange-500 mb-4">EGP {item.price}</p>

          {item.modifiers?.map((mod) => (
            <div key={mod.id} className="mb-5">
              <div className="flex items-center gap-2 mb-2">
                <h3 className="font-semibold text-sm">{mod.name}</h3>
                {mod.required && <Badge variant="destructive" className="text-xs">إجباري</Badge>}
              </div>
              <div className="space-y-2">
                {mod.options.map((opt) => {
                  const isSelected = selectedMods[mod.id]?.includes(opt.id)
                  return (
                    <label key={opt.id} className={`flex items-center justify-between p-3 rounded-lg border-2 cursor-pointer transition-colors ${isSelected ? 'border-orange-500 bg-orange-50' : 'border-gray-200'}`}>
                      <div className="flex items-center gap-3">
                        <input type={mod.multiple ? 'checkbox' : 'radio'} name={mod.id} checked={isSelected} onChange={() => toggleMod(mod.id, opt.id, mod.multiple)} className="w-5 h-5 accent-orange-500" />
                        <span className="text-sm">{opt.name}</span>
                      </div>
                      {opt.priceDelta > 0 && <span className="text-sm text-gray-400">+EGP {opt.priceDelta}</span>}
                    </label>
                  )
                })}
              </div>
            </div>
          ))}

          <div className="mb-5">
            <Label className="font-semibold text-sm mb-2 block">ملاحظات خاصة</Label>
            <Textarea placeholder="مثلاً: بدون بصل" maxLength={200} value={notes} onChange={(e) => setNotes(e.target.value)} />
          </div>
        </ScrollArea>

        <div className="absolute bottom-0 left-0 right-0 bg-white border-t p-4 flex items-center gap-4">
          <div className="flex items-center gap-3">
            <button className="w-10 h-10 rounded-full border-2 flex items-center justify-center" onClick={() => setQty(Math.max(1, qty - 1))} disabled={qty <= 1}><Minus className="w-4 h-4" /></button>
            <span className="text-xl font-bold w-8 text-center">{qty}</span>
            <button className="w-10 h-10 rounded-full border-2 flex items-center justify-center" onClick={() => setQty(Math.min(20, qty + 1))}><Plus className="w-4 h-4" /></button>
          </div>
          <Button className="flex-1 h-12 bg-orange-500 hover:bg-orange-600" disabled={!allRequired} onClick={() => onAdd(qty, selectedModNames, notes)}>
            {allRequired ? `🛒 أضف للسلة — EGP ${total.toFixed(2)}` : 'اختر الخيارات المطلوبة'}
          </Button>
        </div>
      </SheetContent>
    </Sheet>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return <div className="mb-4"><h2 className="font-semibold text-sm mb-2">{title}</h2>{children}</div>
}

function BottomNav({ active, onNavigate, cartCount }: { active: string; onNavigate: (s: Screen) => void; cartCount: number }) {
  const items = [
    { id: 'home', label: 'الرئيسية', icon: Home },
    { id: 'orders', label: 'طلباتي', icon: Receipt },
    { id: 'cart', label: 'السلة', icon: ShoppingBag, badge: cartCount },
    { id: 'profile', label: 'حسابي', icon: User },
  ]
  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 bg-white border-t h-16 flex items-center justify-around px-2">
      {items.map((item) => {
        const Icon = item.icon
        const isActive = active === item.id
        return (
          <button key={item.id} onClick={() => onNavigate(item.id as Screen)} className={`flex flex-col items-center gap-1 flex-1 h-full justify-center transition-colors ${isActive ? 'text-orange-500' : 'text-gray-400'}`}>
            <div className="relative">
              <Icon className="w-5 h-5" />
              {item.badge ? <span className="absolute -top-1 -right-2 w-4 h-4 bg-orange-500 text-white text-xs rounded-full flex items-center justify-center">{item.badge}</span> : null}
            </div>
            <span className="text-xs font-medium">{item.label}</span>
          </button>
        )
      })}
    </nav>
  )
}
