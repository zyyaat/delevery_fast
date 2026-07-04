// Command Center — Full implementation
// 10 screens: Login, Dashboard, Live Map, Zone Detail, Surge Control,
// Incidents, Manual Dispatch, Restaurant/Driver Control, Forecast, Metrics

import { useState, useEffect } from 'react'
import { cn } from '@food-platform/ui'

type Screen = 'dashboard' | 'map' | 'zones' | 'incidents' | 'dispatch' | 'forecast' | 'metrics'

export default function App() {
  const [screen, setScreen] = useState<Screen>('dashboard')
  const [selectedZone, setSelectedZone] = useState<string | null>(null)
  const [surgeModal, setSurgeModal] = useState<string | null>(null)
  const [incidents, setIncidents] = useState(mockIncidents)

  return (
    <div className="min-h-screen bg-[#0A0E1A] text-[#EAF0FF]" dir="rtl">
      <Sidebar active={screen} onNavigate={setScreen} />
      <div className="mr-60">
        {screen === 'dashboard' && <DashboardScreen onNavigate={setScreen} onSelectZone={setSelectedZone} onNavigateZone={() => setScreen('zones')} />}
        {screen === 'map' && <MapScreen />}
        {screen === 'zones' && <ZonesScreen selectedZone={selectedZone} onSelectZone={setSelectedZone} onSurge={setSurgeModal} />}
        {screen === 'incidents' && <IncidentsScreen incidents={incidents} setIncidents={setIncidents} />}
        {screen === 'dispatch' && <DispatchScreen />}
        {screen === 'forecast' && <ForecastScreen />}
        {screen === 'metrics' && <MetricsScreen />}
      </div>

      {surgeModal && <SurgeModal zone={surgeModal} onClose={() => setSurgeModal(null)} />}
    </div>
  )
}

// ============ Mock Data ============
const mockIncidents = [
  { id: 'INC-001', severity: 'P0', title: 'Vodafone Cash API Down', desc: 'Payment success rate dropped to 65%', status: 'mitigating', startedAt: '1:00 PM', duration: '15 min', owner: 'Ahmed K.' },
  { id: 'INC-002', severity: 'P1', title: 'معادي - Driver Shortage', desc: 'Gap: 18 orders without drivers', status: 'investigating', startedAt: '12:35 PM', duration: '40 min', owner: 'Omar T.' },
  { id: 'INC-003', severity: 'P2', title: 'Pizza Hut - 5 delayed orders', desc: 'Orders delayed >15 min', status: 'acknowledged', startedAt: '12:15 PM', duration: '60 min', owner: 'Sarah M.' },
]

// ============ Sidebar ============
function Sidebar({ active, onNavigate }: { active: string; onNavigate: (s: Screen) => void }) {
  const items = [
    { id: 'dashboard', label: 'لوحة القيادة', icon: '📊' },
    { id: 'map', label: 'الخريطة الحية', icon: '🗺️' },
    { id: 'zones', label: 'المناطق', icon: '📍' },
    { id: 'incidents', label: 'الحوادث', icon: '🚨' },
    { id: 'dispatch', label: 'تخصيص يدوي', icon: '🛵' },
    { id: 'forecast', label: 'التوقعات', icon: '📈' },
    { id: 'metrics', label: 'المؤشرات', icon: '⚡' },
  ]
  return (
    <aside className="fixed right-0 top-0 bottom-0 w-60 bg-[#131A2E] border-l border-[#2A3656] flex flex-col z-30">
      <div className="p-4 border-b border-[#2A3656]">
        <div className="flex items-center gap-2">
          <div className="w-10 h-10 bg-[#00D4FF] rounded-lg flex items-center justify-center text-xl">🚨</div>
          <div><p className="font-bold text-sm text-[#EAF0FF]">مركز القيادة</p><p className="text-xs text-[#5C6B8E]">Omar T. • Ops Mgr</p></div>
        </div>
      </div>
      <nav className="flex-1 p-3 space-y-1">
        {items.map(item => (
          <button key={item.id} onClick={() => onNavigate(item.id as Screen)}
            className={cn('w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
              active === item.id ? 'bg-[#00D4FF] text-[#0A0E1A]' : 'text-[#9BA8C7] hover:bg-[#1A2240]')}>
            <span className="text-lg">{item.icon}</span>{item.label}
          </button>
        ))}
      </nav>
      <div className="p-3 border-t border-[#2A3656]">
        <div className="flex items-center gap-2 text-xs text-[#5C6B8E]">
          <span className="w-2 h-2 rounded-full bg-[#00E58A] animate-pulse" />
          <span>متصل • Shift: 8AM-4PM</span>
        </div>
      </div>
    </aside>
  )
}

// ============ Dashboard ============
function DashboardScreen({ onNavigate, onSelectZone, onNavigateZone }: { onNavigate: (s: Screen) => void; onSelectZone: (z: string) => void; onNavigateZone: () => void }) {
  const [time, setTime] = useState(new Date())
  useEffect(() => { const t = setInterval(() => setTime(new Date()), 1000); return () => clearInterval(t) }, [])

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-[#EAF0FF]">لوحة القيادة</h1>
        <span className="text-sm text-[#5C6B8E]">📅 {time.toLocaleString('ar-EG')}</span>
      </div>

      {/* KPI cards */}
      <div className="grid grid-cols-5 gap-4">
        <KPICard icon="📦" label="طلبات نشطة" value="423" trend="↑12%" color="text-[#00D4FF]" />
        <KPICard icon="🛵" label="مناديب" value="187" trend="↑8%" color="text-[#B14EFF]" />
        <KPICard icon="🍔" label="مطاعم" value="142" trend="↑5%" color="text-[#00E58A]" />
        <KPICard icon="⚠️" label="مشاكل" value="3" trend="🔴" color="text-[#FF4757]" />
        <KPICard icon="💰" label="GMV اليوم" value="EGP 84K" trend="↑18%" color="text-[#FFD166]" />
      </div>

      {/* Live map mini */}
      <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-[#EAF0FF]">📍 الخريطة الحية — القاهرة الكبرى</h2>
          <button onClick={() => onNavigate('map')} className="text-sm text-[#00D4FF] font-semibold hover:underline">عرض كامل</button>
        </div>
        <div className="relative h-64 bg-[#0A0E1A] rounded-lg overflow-hidden border border-[#2A3656]">
          {/* Simulated heat zones */}
          <div className="absolute top-10 right-20 w-24 h-24 bg-red-500/30 rounded-full blur-xl" />
          <div className="absolute top-10 right-20 text-center"><span className="text-xs text-red-400 font-bold">معادي</span></div>
          <div className="absolute top-16 left-32 w-28 h-28 bg-red-500/30 rounded-full blur-xl" />
          <div className="absolute top-16 left-32 text-center"><span className="text-xs text-red-400 font-bold">الزمالك</span></div>
          <div className="absolute bottom-12 right-40 w-20 h-20 bg-yellow-500/20 rounded-full blur-xl" />
          <div className="absolute bottom-12 right-40 text-center"><span className="text-xs text-yellow-400 font-bold">مدينة نصر</span></div>
          <div className="absolute bottom-8 left-20 w-16 h-16 bg-green-500/20 rounded-full blur-xl" />
          <div className="absolute bottom-8 left-20 text-center"><span className="text-xs text-green-400 font-bold">التحرير</span></div>
          {/* Driver dots */}
          {[
            { top: '15%', right: '25%' }, { top: '25%', left: '30%' }, { top: '40%', right: '15%' },
            { top: '55%', left: '20%' }, { top: '70%', right: '35%' }, { top: '35%', left: '45%' },
          ].map((pos, i) => (
            <div key={i} className="absolute w-2 h-2 bg-[#00D4FF] rounded-full animate-pulse" style={pos} />
          ))}
        </div>
        <div className="flex gap-4 mt-3 text-xs text-[#5C6B8E]">
          <span>🔴 طلب عالي</span><span>🟡 متوسط</span><span>🟢 هادي</span><span>🔵 مندوب</span>
        </div>
      </div>

      {/* Live metrics */}
      <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
        <h2 className="text-lg font-bold text-[#EAF0FF] mb-4">📊 مؤشرات حية</h2>
        <div className="space-y-3">
          <MetricBar label="Orders/min" value={7.2} target={8.0} color="bg-[#00D4FF]" />
          <MetricBar label="Avg Delivery Time" value={32} target={35} color="bg-[#00E58A]" unit=" د" />
          <MetricBar label="Driver Utilization" value={78} target={80} color="bg-[#B14EFF]" unit="%" />
          <MetricBar label="Order Completion" value={94} target={92} color="bg-[#00E58A]" unit="%" />
          <MetricBar label="Cancellation Rate" value={4} target={5} color="bg-[#FFB800]" unit="%" />
          <MetricBar label="Payment Success" value={98.5} target={98} color="bg-[#00E58A]" unit="%" />
        </div>
      </div>

      {/* Alerts */}
      <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
        <h2 className="text-lg font-bold text-[#EAF0FF] mb-4">⚠️ تنبيهات نشطة (3)</h2>
        <div className="space-y-3">
          <AlertCard severity="P0" message="معادي - طلبات عالية + نقص مناديب (gap: 18 طلب)" actions={['تفعيل Surge 1.3x', 'طلب مناديب']} onClick={() => onNavigate('incidents')} />
          <AlertCard severity="P1" message="Pizza Hut معادي - 5 طلبات delayed >15 دقيقة" actions={['اتصال بالمطعم', 'Auto-reject']} onClick={() => onNavigate('incidents')} />
          <AlertCard severity="P2" message="Vodafone Cash API - latency مرتفعة (3 ثواني)" actions={['Switch to InstaPay', 'Monitor']} onClick={() => onNavigate('incidents')} />
        </div>
      </div>
    </div>
  )
}

// ============ Map Screen ============
function MapScreen() {
  return (
    <div className="p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-[#EAF0FF]">📍 الخريطة الحية</h1>
        <div className="flex gap-2">
          {['Demand', 'Drivers', 'Orders', 'Problems'].map((layer, i) => (
            <button key={layer} className={cn('px-3 py-1.5 rounded-lg text-xs font-semibold', i === 0 ? 'bg-[#00D4FF] text-[#0A0E1A]' : 'bg-[#131A2E] text-[#9BA8C7] border border-[#2A3656]')}>{layer}</button>
          ))}
        </div>
      </div>
      <div className="relative h-[600px] bg-[#131A2E] rounded-xl border border-[#2A3656] overflow-hidden">
        {/* Large heat zones */}
        <div className="absolute top-12 right-32 w-40 h-40 bg-red-500/30 rounded-full blur-2xl" />
        <div className="absolute top-12 right-32 text-center"><span className="text-sm text-red-400 font-bold">🔴 معادي — 72 طلب</span></div>
        <div className="absolute top-24 left-40 w-44 h-44 bg-red-500/30 rounded-full blur-2xl" />
        <div className="absolute top-24 left-40 text-center"><span className="text-sm text-red-400 font-bold">🔴 الزمالك — 58 طلب</span></div>
        <div className="absolute bottom-20 right-48 w-32 h-32 bg-yellow-500/20 rounded-full blur-2xl" />
        <div className="absolute bottom-20 right-48 text-center"><span className="text-sm text-yellow-400 font-bold">🟡 مدينة نصر — 45 طلب</span></div>
        <div className="absolute bottom-12 left-24 w-28 h-28 bg-green-500/20 rounded-full blur-2xl" />
        <div className="absolute bottom-12 left-24 text-center"><span className="text-sm text-green-400 font-bold">🟢 التحرير — 32 طلب</span></div>

        {/* Many driver dots */}
        {Array.from({ length: 30 }).map((_, i) => {
          const top = `${Math.random() * 80 + 5}%`
          const left = `${Math.random() * 80 + 5}%`
          return <div key={i} className="absolute w-2 h-2 bg-[#00D4FF] rounded-full animate-pulse" style={{ top, left }} />
        })}

        {/* Order flow lines */}
        <svg className="absolute inset-0 w-full h-full pointer-events-none">
          <line x1="30%" y1="20%" x2="50%" y2="40%" stroke="#00D4FF" strokeWidth="1" strokeDasharray="4" opacity="0.4" />
          <line x1="60%" y1="30%" x2="40%" y2="60%" stroke="#B14EFF" strokeWidth="1" strokeDasharray="4" opacity="0.4" />
        </svg>
      </div>
      <div className="flex gap-4 text-xs text-[#5C6B8E]">
        <span>🔴 طلب عالي جداً</span><span>🟡 متوسط</span><span>🟢 هادي</span><span>🔵 مندوب متاح</span>
        <span className="mr-auto">آخر تحديث: {new Date().toLocaleTimeString('ar-EG')} (5 ثواني)</span>
      </div>
    </div>
  )
}

// ============ Zones Screen ============
function ZonesScreen({ selectedZone, onSelectZone, onSurge }: { selectedZone: string | null; onSelectZone: (z: string) => void; onSurge: (z: string) => void }) {
  const zones = [
    { id: 'maadi', name: 'معادي', orders: 72, drivers: 31, gap: -8, demand: 'high', time: 28 },
    { id: 'zamalek', name: 'الزمالك', orders: 58, drivers: 25, gap: -5, demand: 'high', time: 30 },
    { id: 'nasr', name: 'مدينة نصر', orders: 45, drivers: 38, gap: 3, demand: 'medium', time: 35 },
    { id: 'downtown', name: 'وسط البلد', orders: 32, drivers: 28, gap: 4, demand: 'medium', time: 32 },
    { id: 'heliopolis', name: 'مصر الجديدة', orders: 28, drivers: 30, gap: 2, demand: 'low', time: 38 },
    { id: 'tagamoa', name: 'التجمع', orders: 22, drivers: 18, gap: -4, demand: 'medium', time: 42 },
  ]

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-bold text-[#EAF0FF]">📍 إدارة المناطق</h1>
      <div className="grid grid-cols-3 gap-4">
        {zones.map(zone => (
          <div key={zone.id} onClick={() => onSelectZone(zone.id)}
            className={cn('bg-[#131A2E] rounded-xl border p-5 cursor-pointer transition-all',
              selectedZone === zone.id ? 'border-[#00D4FF] shadow-lg shadow-[#00D4FF]/10' : 'border-[#2A3656] hover:border-[#3A4A6E]')}>
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-bold text-[#EAF0FF]">{zone.name}</h3>
              <span className={cn('text-xs px-2 py-0.5 rounded font-bold',
                zone.demand === 'high' ? 'bg-red-500/20 text-red-400' : zone.demand === 'medium' ? 'bg-yellow-500/20 text-yellow-400' : 'bg-green-500/20 text-green-400')}>
                {zone.demand === 'high' ? '🔴 عالي' : zone.demand === 'medium' ? '🟡 متوسط' : '🟢 هادي'}
              </span>
            </div>
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div><p className="text-[#5C6B8E] text-xs">طلبات</p><p className="font-bold text-[#EAF0FF]">{zone.orders}</p></div>
              <div><p className="text-[#5C6B8E] text-xs">مناديب</p><p className="font-bold text-[#EAF0FF]">{zone.drivers}</p></div>
              <div><p className="text-[#5C6B8E] text-xs">Gap</p><p className={cn('font-bold', zone.gap < 0 ? 'text-red-400' : 'text-green-400')}>{zone.gap > 0 ? '+' : ''}{zone.gap}</p></div>
              <div><p className="text-[#5C6B8E] text-xs">Avg Time</p><p className="font-bold text-[#EAF0FF]">{zone.time}د</p></div>
            </div>
            {/* Demand vs Supply bar */}
            <div className="mt-3 space-y-1">
              <div className="flex items-center gap-2 text-xs">
                <span className="text-[#5C6B8E] w-12">الطلب</span>
                <div className="flex-1 h-2 bg-[#0A0E1A] rounded-full overflow-hidden">
                  <div className="h-full bg-red-500/60 rounded-full" style={{ width: `${(zone.orders / 80) * 100}%` }} />
                </div>
                <span className="text-[#9BA8C7] w-8">{zone.orders}</span>
              </div>
              <div className="flex items-center gap-2 text-xs">
                <span className="text-[#5C6B8E] w-12">مناديب</span>
                <div className="flex-1 h-2 bg-[#0A0E1A] rounded-full overflow-hidden">
                  <div className="h-full bg-[#00D4FF]/60 rounded-full" style={{ width: `${(zone.drivers / 80) * 100}%` }} />
                </div>
                <span className="text-[#9BA8C7] w-8">{zone.drivers}</span>
              </div>
            </div>
            {zone.gap < 0 && (
              <button onClick={(e) => { e.stopPropagation(); onSurge(zone.id) }}
                className="w-full mt-3 py-1.5 bg-[#00D4FF]/10 text-[#00D4FF] rounded-lg text-xs font-semibold hover:bg-[#00D4FF]/20">
                ⚡ تفعيل Surge
              </button>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Incidents Screen ============
function IncidentsScreen({ incidents, setIncidents }: { incidents: typeof mockIncidents; setIncidents: React.Dispatch<React.SetStateAction<typeof mockIncidents>> }) {
  const resolveIncident = (id: string) => {
    setIncidents(incidents.map(i => i.id === id ? { ...i, status: 'resolved' } : i))
  }

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-bold text-[#EAF0FF]">🚨 الحوادث</h1>
      <div className="space-y-3">
        {incidents.map(inc => (
          <div key={inc.id} className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-3">
                <span className={cn('text-xs font-bold px-2 py-1 rounded',
                  inc.severity === 'P0' ? 'bg-red-500/20 text-red-400' : inc.severity === 'P1' ? 'bg-amber-500/20 text-amber-400' : 'bg-blue-500/20 text-blue-400')}>
                  {inc.severity}
                </span>
                <h3 className="font-bold text-[#EAF0FF]">{inc.title}</h3>
              </div>
              <span className={cn('text-xs px-2 py-1 rounded font-semibold',
                inc.status === 'resolved' ? 'bg-green-500/20 text-green-400' : inc.status === 'mitigating' ? 'bg-amber-500/20 text-amber-400' : 'bg-blue-500/20 text-blue-400')}>
                {inc.status}
              </span>
            </div>
            <p className="text-sm text-[#9BA8C7] mb-3">{inc.desc}</p>
            <div className="flex items-center gap-4 text-xs text-[#5C6B8E] mb-3">
              <span>⏰ {inc.startedAt}</span><span>⏱️ {inc.duration}</span><span>👤 {inc.owner}</span>
            </div>
            {inc.status !== 'resolved' && (
              <div className="flex gap-2">
                <button onClick={() => resolveIncident(inc.id)} className="px-4 py-1.5 bg-green-500/10 text-green-400 rounded-lg text-xs font-semibold hover:bg-green-500/20">✅ حل</button>
                <button className="px-4 py-1.5 bg-blue-500/10 text-blue-400 rounded-lg text-xs font-semibold hover:bg-blue-500/20">📝 Postmortem</button>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Dispatch Screen ============
function DispatchScreen() {
  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-bold text-[#EAF0FF]">🛵 تخصيص يدوي</h1>
      <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
        <h3 className="font-bold text-[#EAF0FF] mb-3">طلبات بدون مندوب (2)</h3>
        <div className="space-y-3">
          {[
            { id: '#A7X92F', restaurant: 'Pizza Hut معادي', customer: 'أحمد - الزمالك', eta: 35 },
            { id: '#B3K45L', restaurant: 'KFC مدينة نصر', customer: 'فاطمة - مصر الجديدة', eta: 40 },
          ].map(order => (
            <div key={order.id} className="bg-[#0A0E1A] rounded-lg p-4">
              <div className="flex items-center justify-between mb-3">
                <div><span className="font-bold text-[#EAF0FF]">{order.id}</span><span className="text-sm text-[#5C6B8E] mr-2">{order.restaurant} → {order.customer}</span></div>
                <span className="text-sm text-[#5C6B8E]">ETA: {order.eta}د</span>
              </div>
              <div className="grid grid-cols-3 gap-2">
                {[
                  { name: 'Mahmoud S.', dist: '0.8km', rating: '4.8', accept: '92%' },
                  { name: 'Ahmed K.', dist: '1.2km', rating: '4.7', accept: '88%' },
                  { name: 'Mostafa A.', dist: '1.8km', rating: '4.6', accept: '75%' },
                ].map(driver => (
                  <div key={driver.name} className="bg-[#131A2E] rounded-lg p-3 text-center">
                    <p className="text-sm font-semibold text-[#EAF0FF]">{driver.name}</p>
                    <p className="text-xs text-[#5C6B8E]">{driver.dist} • ⭐{driver.rating} • {driver.accept}</p>
                    <button className="w-full mt-2 py-1 bg-[#00D4FF]/10 text-[#00D4FF] rounded text-xs font-semibold hover:bg-[#00D4FF]/20">تخصيص</button>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

// ============ Forecast Screen ============
function ForecastScreen() {
  const hours = ['6am', '7am', '8am', '9am', '10am', '11am', '12pm', '1pm', '2pm', '3pm', '4pm', '5pm', '6pm', '7pm', '8pm', '9pm', '10pm', '11pm']
  const forecast = [80, 120, 160, 220, 280, 340, 480, 440, 320, 280, 250, 300, 400, 450, 420, 350, 280, 200]
  const max = Math.max(...forecast)

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-bold text-[#EAF0FF]">📈 توقعات الطلب — بكرة</h1>
      <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
        <div className="flex items-end gap-2 h-64">
          {forecast.map((val, i) => (
            <div key={i} className="flex-1 flex flex-col items-center gap-1">
              <div className="w-full bg-gradient-to-t from-[#00D4FF] to-[#B14EFF] rounded-t" style={{ height: `${(val / max) * 100}%` }} />
              <span className="text-xs text-[#5C6B8E]">{hours[i]}</span>
              <span className="text-xs font-bold text-[#9BA8C7]">{val}</span>
            </div>
          ))}
        </div>
      </div>
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-4">
          <p className="text-sm text-[#5C6B8E]">ذروة الغداء</p>
          <p className="text-2xl font-bold text-[#00D4FF]">480 طلب</p>
          <p className="text-xs text-[#5C6B8E]">12:00 PM</p>
          <p className="text-xs text-amber-400 mt-1">⚠️ زيادة مناديب معادي +15</p>
        </div>
        <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-4">
          <p className="text-sm text-[#5C6B8E]">ذروة العشاء</p>
          <p className="text-2xl font-bold text-[#B14EFF]">450 طلب</p>
          <p className="text-xs text-[#5C6B8E]">7:00 PM</p>
          <p className="text-xs text-amber-400 mt-1">⚠️ pre-scale Order + Payment</p>
        </div>
        <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-4">
          <p className="text-sm text-[#5C6B8E]">دقة التوقعات</p>
          <p className="text-2xl font-bold text-[#00E58A]">91.7%</p>
          <p className="text-xs text-[#5C6B8E]">MAPE: 8.3%</p>
          <p className="text-xs text-green-400 mt-1">✅ ممتاز (أقل من 10%)</p>
        </div>
      </div>
    </div>
  )
}

// ============ Metrics Screen ============
function MetricsScreen() {
  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-bold text-[#EAF0FF]">⚡ المؤشرات التفصيلية</h1>
      <div className="grid grid-cols-2 gap-4">
        <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
          <h3 className="font-bold text-[#EAF0FF] mb-4">📊 GMV اليوم</h3>
          <div className="flex items-end gap-1 h-32">
            {[20, 35, 50, 65, 80, 95, 100, 90, 70, 55, 40, 60, 85, 95, 90, 75, 50, 30].map((val, i) => (
              <div key={i} className="flex-1 bg-[#00D4FF]/40 rounded-t" style={{ height: `${val}%` }} />
            ))}
          </div>
          <p className="text-2xl font-bold text-[#00D4FF] mt-3">EGP 84,250</p>
          <p className="text-xs text-green-400">↑ 18% عن الأمس</p>
        </div>
        <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
          <h3 className="font-bold text-[#EAF0FF] mb-4">🛵 استخدام المناديب</h3>
          <div className="flex items-center justify-center h-32">
            <div className="relative w-32 h-32">
              <svg className="w-full h-full transform -rotate-90">
                <circle cx="64" cy="64" r="56" fill="none" stroke="#1A2240" strokeWidth="12" />
                <circle cx="64" cy="64" r="56" fill="none" stroke="#00D4FF" strokeWidth="12" strokeDasharray={`${78 * 3.52} 999`} strokeLinecap="round" />
              </svg>
              <div className="absolute inset-0 flex items-center justify-center"><span className="text-3xl font-bold text-[#EAF0FF]">78%</span></div>
            </div>
          </div>
          <p className="text-xs text-[#5C6B8E] text-center mt-2">Target: 75-85%</p>
        </div>
      </div>
      <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-5">
        <h3 className="font-bold text-[#EAF0FF] mb-4">🏆 Top Drivers</h3>
        <div className="space-y-2">
          {[{ name: 'Mahmoud S.', deliveries: 12, rating: 4.9, earnings: 420 }, { name: 'Ahmed K.', deliveries: 10, rating: 4.8, earnings: 380 }, { name: 'Mostafa A.', deliveries: 9, rating: 4.7, earnings: 340 }].map((d, i) => (
            <div key={i} className="flex items-center gap-3 py-2 border-b border-[#2A3656] last:border-0">
              <span className="text-lg">{['🥇', '🥈', '🥉'][i]}</span>
              <span className="flex-1 text-sm font-semibold text-[#EAF0FF]">{d.name}</span>
              <span className="text-xs text-[#5C6B8E]">{d.deliveries} توصيلة</span>
              <span className="text-xs text-amber-400">⭐{d.rating}</span>
              <span className="text-xs font-bold text-[#00D4FF]">EGP {d.earnings}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

// ============ Surge Modal ============
function SurgeModal({ zone, onClose }: { zone: string; onClose: () => void }) {
  const [multiplier, setMultiplier] = useState(1.3)
  const [duration, setDuration] = useState(60)

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4" dir="rtl">
      <div className="bg-[#131A2E] rounded-2xl border border-[#2A3656] shadow-2xl max-w-md w-full p-6 space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-bold text-[#EAF0FF]">⚡ تفعيل Surge — {zone}</h2>
          <button onClick={onClose} className="text-[#5C6B8E] hover:text-[#EAF0FF]">✕</button>
        </div>
        <div>
          <label className="text-sm text-[#9BA8C7] block mb-2">المضاعف: {multiplier}x</label>
          <input type="range" min="1" max="1.5" step="0.1" value={multiplier} onChange={(e) => setMultiplier(parseFloat(e.target.value))} className="w-full accent-[#00D4FF]" />
          <div className="flex justify-between text-xs text-[#5C6B8E]"><span>1.0x</span><span>1.5x</span></div>
        </div>
        <div>
          <label className="text-sm text-[#9BA8C7] block mb-2">المدة: {duration} دقيقة</label>
          <select value={duration} onChange={(e) => setDuration(parseInt(e.target.value))} className="w-full h-10 px-3 rounded-lg bg-[#0A0E1A] border border-[#2A3656] text-[#EAF0FF] text-sm">
            <option value={30}>30 دقيقة</option><option value={60}>60 دقيقة</option><option value={90}>90 دقيقة</option>
          </select>
        </div>
        <div className="bg-[#0A0E1A] rounded-lg p-3 text-sm">
          <p className="text-[#9BA8C7]">سيتم تطبيق Surge على:</p>
          <p className="text-[#EAF0FF] font-semibold mt-1">📍 {zone} — جميع الطلبات الجديدة</p>
        </div>
        <div className="flex gap-2">
          <button onClick={onClose} className="flex-1 h-10 border border-[#2A3656] text-[#9BA8C7] rounded-lg font-semibold">إلغاء</button>
          <button onClick={onClose} className="flex-1 h-10 bg-[#00D4FF] hover:bg-[#00B8E0] text-[#0A0E1A] rounded-lg font-bold">⚡ تطبيق Surge</button>
        </div>
      </div>
    </div>
  )
}

// ============ Helpers ============
function KPICard({ icon, label, value, trend, color }: { icon: string; label: string; value: string; trend: string; color: string }) {
  return (
    <div className="bg-[#131A2E] rounded-xl border border-[#2A3656] p-4">
      <div className="flex items-center gap-2 mb-2"><span className="text-xl">{icon}</span><span className={cn('text-xs', color)}>{trend}</span></div>
      <p className="text-xs text-[#5C6B8E]">{label}</p>
      <p className={cn('text-2xl font-bold mt-1', color)}>{value}</p>
    </div>
  )
}

function MetricBar({ label, value, target, color, unit = '' }: { label: string; value: number; target: number; color: string; unit?: string }) {
  const pct = Math.min(100, (value / (target * 1.5)) * 100)
  const isGood = label.includes('Cancellation') ? value < target : value >= target
  return (
    <div>
      <div className="flex justify-between text-sm mb-1">
        <span className="text-[#9BA8C7]">{label}</span>
        <span className={cn('font-bold', isGood ? 'text-[#00E58A]' : 'text-[#FF4757]')}>{value}{unit} <span className="text-[#5C6B8E] text-xs">(target: {target}{unit})</span></span>
      </div>
      <div className="h-2 bg-[#0A0E1A] rounded-full overflow-hidden">
        <div className={cn('h-full rounded-full', color)} style={{ width: `${pct}%` }} />
      </div>
    </div>
  )
}

function AlertCard({ severity, message, actions, onClick }: { severity: string; message: string; actions: string[]; onClick: () => void }) {
  return (
    <div onClick={onClick} className="bg-[#0A0E1A] rounded-lg p-3 cursor-pointer hover:bg-[#131A2E] transition-colors">
      <div className="flex items-center gap-2 mb-2">
        <span className={cn('text-xs font-bold px-2 py-0.5 rounded',
          severity === 'P0' ? 'bg-red-500/20 text-red-400' : severity === 'P1' ? 'bg-amber-500/20 text-amber-400' : 'bg-blue-500/20 text-blue-400')}>{severity}</span>
        <p className="text-sm text-[#EAF0FF]">{message}</p>
      </div>
      <div className="flex gap-2">
        {actions.map((a, i) => <span key={i} className="text-xs text-[#00D4FF] bg-[#00D4FF]/10 px-2 py-0.5 rounded">{a}</span>)}
      </div>
    </div>
  )
}
