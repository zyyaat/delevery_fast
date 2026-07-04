// Employee Portal — Full implementation
// 12 screens: Login, WebAuthn, Sensitive Action, Dashboard, Activity,
// Refund, Restaurant Approve, Driver Approve, Payout Approve,
// Audit Log, Anomaly Alerts, Access Review

import { useState, useEffect } from 'react'
import { cn } from '@food-platform/ui'

type Screen = 'login' | 'webauthn' | 'dashboard' | 'audit' | 'anomaly' | 'access' | 'approvals'

export default function App() {
  const [screen, setScreen] = useState<Screen>('login')
  const [authenticated, setAuthenticated] = useState(false)
  const [webauthnEnrolled, setWebauthnEnrolled] = useState(false)
  const [biometricModal, setBiometricModal] = useState<BiometricAction | null>(null)
  const [actionToken, setActionToken] = useState<string | null>(null)

  if (!authenticated && screen !== 'login') {
    setScreen('login')
  }

  return (
    <div className="min-h-screen bg-gray-50" dir="rtl">
      {screen === 'login' && (
        <LoginScreen onLogin={() => { setAuthenticated(true); setScreen('dashboard') }} />
      )}

      {authenticated && !webauthnEnrolled && screen !== 'login' && (
        <WebAuthnScreen onEnroll={() => { setWebauthnEnrolled(true); setScreen('dashboard') }} />
      )}

      {authenticated && webauthnEnrolled && (
        <>
          {screen === 'dashboard' && <DashboardScreen onNavigate={setScreen} onAction={(a) => setBiometricModal(a)} />}
          {screen === 'audit' && <AuditScreen />}
          {screen === 'anomaly' && <AnomalyScreen />}
          {screen === 'access' && <AccessReviewScreen />}
          {screen === 'approvals' && <ApprovalsScreen onAction={(a) => setBiometricModal(a)} />}
          <Sidebar active={screen} onNavigate={setScreen} />
        </>
      )}

      {biometricModal && (
        <BiometricModal
          action={biometricModal}
          onClose={() => setBiometricModal(null)}
          onVerified={() => { setActionToken(Date.now().toString()); setBiometricModal(null) }}
        />
      )}
    </div>
  )
}

type BiometricAction = {
  type: string
  title: string
  details: { label: string; value: string }[]
}

// ============ Mock Data ============
const mockPendingApprovals = [
  { id: 'APR-001', type: 'refund', title: 'Refund Request — EGP 750', from: 'Sarah M. (Support L2)', entity: 'Order #A7X92F', customer: 'Ahmed M.', amount: 750, reason: 'Order never delivered', urgent: true },
  { id: 'APR-002', type: 'restaurant', title: 'Restaurant Activation', from: 'Omar T. (Ops Manager)', entity: 'KFC - Tagamoa', details: 'Commission: 18%', urgent: false },
  { id: 'APR-003', type: 'driver', title: 'Driver Activation', from: 'Fatma A. (Field Supervisor)', entity: 'Mostafa A. — Motorcycle', details: 'KYC verified, license valid', urgent: false },
  { id: 'APR-004', type: 'payout', title: 'Driver Payout — EGP 5,000', from: 'Finance System', entity: 'Driver: Mahmoud S.', amount: 5000, reason: 'Weekly payout', urgent: false },
]

const mockAuditLogs = [
  { time: '14:32:15', actor: 'ahmed.k', action: 'refund.issued', entity: '#A7X92F', details: 'EGP 300', bio: true, ip: '192.168.1.100' },
  { time: '14:25:03', actor: 'ahmed.k', action: 'customer.viewed', entity: 'uuid-cust-123', details: '', bio: false, ip: '192.168.1.100' },
  { time: '14:18:42', actor: 'ahmed.k', action: 'order.cancelled', entity: '#B3K45L', details: 'reason: delayed', bio: false, ip: '192.168.1.100' },
  { time: '14:05:00', actor: 'ahmed.k', action: 'auth.login', entity: 'session', details: 'biometric', bio: true, ip: '192.168.1.100' },
  { time: '13:45:22', actor: 'sarah.m', action: 'restaurant.approved', entity: 'KFC-Tagamoa', details: 'commission: 18%', bio: true, ip: '10.0.0.55' },
  { time: '13:30:10', actor: 'omar.t', action: 'driver.suspended', entity: 'uuid-driver-456', details: 'fraud: GPS spoofing', bio: true, ip: '10.0.0.42' },
  { time: '13:15:30', actor: 'sarah.m', action: 'refund.issued', entity: '#C8M72N', details: 'EGP 145', bio: true, ip: '10.0.0.55' },
  { time: '12:50:00', actor: 'ahmed.k', action: 'payout.approved', entity: 'driver: Mahmoud S.', details: 'EGP 5,000', bio: true, ip: '192.168.1.100' },
]

const mockAnomalies = [
  { id: 'ANM-001', employee: 'ahmed.k', role: 'Support L2', score: -0.92, severity: 'critical', reasons: ['25 refunds today (avg: 8)', 'Avg amount EGP 200 (avg: 80)', 'Activity at 2 AM (unusual)', 'Login from new IP'] },
  { id: 'ANM-002', employee: 'sarah.m', role: 'Support L2', score: -0.72, severity: 'medium', reasons: ['5 refunds to same customer in 7 days'] },
  { id: 'ANM-003', employee: 'omar.t', role: 'Ops Manager', score: -0.55, severity: 'low', reasons: ['3 restaurant approvals in 1 hour (avg: 1)'] },
]

// ============ Login Screen ============
function LoginScreen({ onLogin }: { onLogin: () => void }) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [totp, setTotp] = useState('')

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 to-slate-800 flex items-center justify-center p-6">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-slate-700 rounded-2xl mb-4 border border-slate-600">
            <span className="text-white text-3xl">🔐</span>
          </div>
          <h1 className="text-2xl font-bold text-white">بوابة الموظفين</h1>
          <p className="text-slate-400 mt-2 text-sm">Internal Employee Portal</p>
        </div>

        <div className="bg-slate-800 rounded-2xl border border-slate-700 p-6 space-y-4">
          <div>
            <label className="text-sm text-slate-400 block mb-1">اسم المستخدم</label>
            <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} placeholder="ahmed.k@food-platform.com"
              className="w-full h-12 px-3 rounded-lg bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-2 focus:ring-blue-500" />
          </div>
          <div>
            <label className="text-sm text-slate-400 block mb-1">كلمة المرور</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••••••"
              className="w-full h-12 px-3 rounded-lg bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-2 focus:ring-blue-500" />
          </div>
          <div>
            <label className="text-sm text-slate-400 block mb-1">رمز TOTP (6 أرقام)</label>
            <input type="text" inputMode="numeric" maxLength={6} value={totp} onChange={(e) => setTotp(e.target.value.replace(/\D/g, ''))}
              placeholder="123456" className="w-full h-12 px-3 rounded-lg bg-slate-900 border border-slate-700 text-white text-center text-2xl tracking-widest focus:outline-none focus:ring-2 focus:ring-blue-500" />
          </div>
          <button onClick={onLogin} disabled={!username || !password || totp.length < 6}
            className="w-full h-12 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold disabled:opacity-50 transition-colors">
            تسجيل الدخول
          </button>
        </div>

        <div className="bg-amber-900/20 border border-amber-800/50 rounded-lg p-3 text-center">
          <p className="text-xs text-amber-400">⚠️ جميع الإجراءات مسجلة ومراقبة — Law 175/2018</p>
        </div>
      </div>
    </div>
  )
}

// ============ WebAuthn Enrollment ============
function WebAuthnScreen({ onEnroll }: { onEnroll: () => void }) {
  const [step, setStep] = useState<'intro' | 'enrolling' | 'success'>('intro')

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 to-slate-800 flex items-center justify-center p-6">
      <div className="w-full max-w-md space-y-6">
        <div className="bg-slate-800 rounded-2xl border border-slate-700 p-8 text-center">
          {step === 'intro' && (
            <>
              <div className="inline-flex items-center justify-center w-20 h-20 bg-blue-600 rounded-full mb-4">
                <span className="text-white text-4xl">🔑</span>
              </div>
              <h2 className="text-xl font-bold text-white mb-2">تفعيل البصمة الحيوية</h2>
              <p className="text-slate-400 text-sm mb-6">عشان نحسن الأمان، لازم تفعل البصمة (Touch ID / Face ID / Windows Hello) للإجراءات الحساسة</p>
              <div className="bg-slate-900 rounded-lg p-4 text-right mb-6 space-y-2">
                <p className="text-xs text-slate-300">✅ استرجاع الأموال (Refund)</p>
                <p className="text-xs text-slate-300">✅ تفعيل المطاعم</p>
                <p className="text-xs text-slate-300">✅ تفعيل المناديب</p>
                <p className="text-xs text-slate-300">✅ اعتماد المدفوعات</p>
                <p className="text-xs text-slate-300">✅ تعديل المنيو</p>
              </div>
              <button onClick={() => { setStep('enrolling'); setTimeout(() => setStep('success'), 2000) }}
                className="w-full h-12 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold">🔑 تفعيل الآن</button>
            </>
          )}

          {step === 'enrolling' && (
            <div className="py-8">
              <div className="inline-flex items-center justify-center w-20 h-20 bg-blue-600 rounded-full mb-4 animate-pulse">
                <span className="text-white text-4xl">🔑</span>
              </div>
              <p className="text-white font-semibold mb-2">جاري التفعيل...</p>
              <p className="text-slate-400 text-sm">ضع إصبعك على مستشعر البصمة</p>
            </div>
          )}

          {step === 'success' && (
            <>
              <div className="inline-flex items-center justify-center w-20 h-20 bg-green-500 rounded-full mb-4">
                <span className="text-white text-4xl">✅</span>
              </div>
              <h2 className="text-xl font-bold text-white mb-2">تم التفعيل!</h2>
              <p className="text-slate-400 text-sm mb-6">البصمة الحيوية مفعّلة على هذا الجهاز</p>
              <button onClick={onEnroll} className="w-full h-12 bg-green-600 hover:bg-green-700 text-white rounded-lg font-semibold">متابعة</button>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

// ============ Sidebar ============
function Sidebar({ active, onNavigate }: { active: string; onNavigate: (s: Screen) => void }) {
  const items = [
    { id: 'dashboard', label: 'لوحة التحكم', icon: '📊' },
    { id: 'approvals', label: 'اعتمادات معلقة', icon: '⏳', badge: 4 },
    { id: 'audit', label: 'سجل التدقيق', icon: '📜' },
    { id: 'anomaly', label: 'تنبيهات الأنومالي', icon: '⚠️', badge: 2 },
    { id: 'access', label: 'مراجعة الصلاحيات', icon: '🔑' },
  ]
  return (
    <aside className="fixed right-0 top-0 bottom-0 w-60 bg-white border-l border-gray-200 flex flex-col z-30">
      <div className="p-4 border-b">
        <div className="flex items-center gap-2">
          <div className="w-10 h-10 bg-slate-800 rounded-lg flex items-center justify-center text-xl">🔐</div>
          <div>
            <p className="font-bold text-sm">بوابة الموظفين</p>
            <p className="text-xs text-gray-400">Ahmed K. • Support L2</p>
          </div>
        </div>
      </div>
      <nav className="flex-1 p-3 space-y-1">
        {items.map(item => (
          <button key={item.id} onClick={() => onNavigate(item.id as Screen)}
            className={cn('w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors relative',
              active === item.id ? 'bg-slate-800 text-white' : 'text-gray-600 hover:bg-gray-100')}>
            <span className="text-lg">{item.icon}</span>{item.label}
            {item.badge && <span className="absolute left-3 w-5 h-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">{item.badge}</span>}
          </button>
        ))}
      </nav>
      <div className="p-3 border-t space-y-2">
        <div className="bg-green-50 rounded-lg p-2 text-center">
          <p className="text-xs text-green-600">🔑 البصمة مفعّلة</p>
          <p className="text-xs text-gray-400 mt-1">جلسة: 2س 15د</p>
        </div>
        <button className="w-full text-sm text-red-500 py-2 hover:bg-red-50 rounded-lg">🚪 تسجيل الخروج</button>
      </div>
    </aside>
  )
}

// ============ Dashboard ============
function DashboardScreen({ onNavigate, onAction }: { onNavigate: (s: Screen) => void; onAction: (a: BiometricAction) => void }) {
  return (
    <div className="mr-60 p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">لوحة التحكم</h1>
        <span className="text-sm text-gray-400">📅 4 يوليو 2026 | 🟢 متصل</span>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-4 gap-4">
        <StatCard icon="📋" label="طلبات اليوم" value="28" color="bg-blue-50 text-blue-600" />
        <StatCard icon="💰" label="استرجاعات" value="8" color="bg-green-50 text-green-600" />
        <StatCard icon="⭐" label="أداء CSAT" value="90%" color="bg-amber-50 text-amber-600" />
        <StatCard icon="⚠️" label="تنبيهات" value="0" color="bg-gray-50 text-gray-500" />
      </div>

      {/* Pending approvals */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold">⏳ اعتمادات معلقة (4)</h2>
          <button onClick={() => onNavigate('approvals')} className="text-sm text-blue-600 font-semibold hover:underline">عرض الكل</button>
        </div>
        <div className="space-y-3">
          {mockPendingApprovals.map(approval => (
            <div key={approval.id} className="bg-white rounded-xl border p-4 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className={cn('text-xs font-bold px-2 py-1 rounded',
                  approval.type === 'refund' ? 'bg-green-50 text-green-600' : approval.type === 'restaurant' ? 'bg-purple-50 text-purple-600' : approval.type === 'driver' ? 'bg-blue-50 text-blue-600' : 'bg-amber-50 text-amber-600')}>
                  {approval.type === 'refund' ? '💰' : approval.type === 'restaurant' ? '🍔' : approval.type === 'driver' ? '🛵' : '💸'} {approval.type}
                </span>
                <div>
                  <p className="font-semibold text-sm text-gray-900">{approval.title}</p>
                  <p className="text-xs text-gray-400">من: {approval.from} • {approval.entity}</p>
                </div>
                {approval.urgent && <span className="text-xs bg-red-50 text-red-500 px-2 py-0.5 rounded">عاجل</span>}
              </div>
              <div className="flex gap-2">
                <button onClick={() => onAction({
                  type: approval.type,
                  title: approval.title,
                  details: [
                    { label: 'الطلب', value: approval.entity },
                    { label: 'القيمة', value: approval.amount ? `EGP ${approval.amount}` : approval.details || '-' },
                    { label: 'السبب', value: approval.reason || approval.details || '-' },
                  ]
                })} className="px-4 py-1.5 bg-green-500 text-white rounded-lg text-sm font-semibold hover:bg-green-600">✅ اعتماد</button>
                <button className="px-4 py-1.5 border border-gray-200 text-gray-500 rounded-lg text-sm">❌ رفض</button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Recent activity */}
      <div className="bg-white rounded-xl border p-5">
        <h3 className="font-bold mb-3">📜 النشاط الأخير</h3>
        <div className="space-y-2 text-sm">
          {mockAuditLogs.slice(0, 5).map((log, i) => (
            <div key={i} className="flex items-center gap-3 py-1 border-b last:border-0">
              <span className="text-xs text-gray-400 w-16">{log.time}</span>
              <span className="font-mono text-xs text-gray-600">{log.action}</span>
              <span className="text-xs text-gray-400">{log.entity}</span>
              {log.details && <span className="text-xs text-gray-400">{log.details}</span>}
              {log.bio && <span className="text-xs text-green-500">🔑</span>}
            </div>
          ))}
        </div>
      </div>

      {/* Quick actions */}
      <div className="grid grid-cols-4 gap-3">
        <QuickAction icon="💰" label="استرجاع" onClick={() => onAction({ type: 'refund', title: 'معالجة استرجاع', details: [{ label: 'الطلب', value: '#A7X92F' }, { label: 'العميل', value: 'Ahmed M.' }, { label: 'المبلغ', value: 'EGP 300' }, { label: 'السبب', value: 'order_delayed_45min' }] })} />
        <QuickAction icon="🍔" label="تفعيل مطعم" onClick={() => onAction({ type: 'restaurant', title: 'تفعيل مطعم', details: [{ label: 'المطعم', value: 'KFC - Tagamoa' }, { label: 'العمولة', value: '18%' }] })} />
        <QuickAction icon="🛵" label="تفعيل مندوب" onClick={() => onAction({ type: 'driver', title: 'تفعيل مندوب', details: [{ label: 'المندوب', value: 'Mostafa A.' }, { label: 'المركبة', value: 'موتوسيكل' }] })} />
        <QuickAction icon="💸" label="اعتماد دفعة" onClick={() => onAction({ type: 'payout', title: 'اعتماد دفعة', details: [{ label: 'المستلم', value: 'Mahmoud S.' }, { label: 'المبلغ', value: 'EGP 5,000' }] })} />
      </div>
    </div>
  )
}

// ============ Approvals Screen ============
function ApprovalsScreen({ onAction }: { onAction: (a: BiometricAction) => void }) {
  return (
    <div className="mr-60 p-6 space-y-4">
      <h1 className="text-2xl font-bold">⏳ الاعتمادات المعلقة</h1>
      <div className="space-y-3">
        {mockPendingApprovals.map(approval => (
          <div key={approval.id} className="bg-white rounded-xl border p-5">
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-3">
                <span className={cn('text-xs font-bold px-2 py-1 rounded',
                  approval.type === 'refund' ? 'bg-green-50 text-green-600' : approval.type === 'restaurant' ? 'bg-purple-50 text-purple-600' : approval.type === 'driver' ? 'bg-blue-50 text-blue-600' : 'bg-amber-50 text-amber-600')}>
                  {approval.type}
                </span>
                {approval.urgent && <span className="text-xs bg-red-50 text-red-500 px-2 py-0.5 rounded">عاجل</span>}
              </div>
              <span className="text-xs text-gray-400">{approval.id}</span>
            </div>
            <h3 className="font-bold text-gray-900">{approval.title}</h3>
            <p className="text-sm text-gray-500 mt-1">من: {approval.from}</p>
            <div className="bg-gray-50 rounded-lg p-3 mt-3 text-sm space-y-1">
              <div className="flex justify-between"><span className="text-gray-400">الجهة</span><span className="font-medium">{approval.entity}</span></div>
              {approval.amount && <div className="flex justify-between"><span className="text-gray-400">المبلغ</span><span className="font-bold">EGP {approval.amount}</span></div>}
              {approval.reason && <div className="flex justify-between"><span className="text-gray-400">السبب</span><span className="font-medium">{approval.reason}</span></div>}
              {approval.details && <div className="flex justify-between"><span className="text-gray-400">التفاصيل</span><span className="font-medium">{approval.details}</span></div>}
            </div>
            <div className="flex gap-2 mt-4">
              <button onClick={() => onAction({
                type: approval.type, title: approval.title,
                details: [
                  { label: 'الجهة', value: approval.entity },
                  { label: 'المبلغ', value: approval.amount ? `EGP ${approval.amount}` : approval.details || '-' },
                  { label: 'السبب', value: approval.reason || approval.details || '-' },
                ]
              })} className="flex-1 h-10 bg-green-500 hover:bg-green-600 text-white rounded-lg font-semibold">✅ اعتماد</button>
              <button className="flex-1 h-10 border border-gray-200 text-gray-500 rounded-lg font-semibold">❌ رفض</button>
              <button className="px-4 h-10 border border-gray-200 text-gray-500 rounded-lg">👁️ مراجعة</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Biometric Modal ============
function BiometricModal({ action, onClose, onVerified }: { action: BiometricAction; onClose: () => void; onVerified: () => void }) {
  const [step, setStep] = useState<'confirm' | 'scanning' | 'verified'>('confirm')

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" dir="rtl">
      <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full overflow-hidden">
        {/* Header */}
        <div className="bg-amber-50 border-b border-amber-200 p-4">
          <div className="flex items-center gap-2">
            <span className="text-xl">⚠️</span>
            <h2 className="font-bold text-amber-800">إجراء حساس — تأكيد الهوية</h2>
          </div>
        </div>

        <div className="p-6 space-y-4">
          {step === 'confirm' && (
            <>
              {/* Action details */}
              <div className="bg-gray-50 rounded-lg p-4 space-y-2">
                <p className="text-xs text-gray-400">الإجراء</p>
                <p className="font-bold text-gray-900">{action.title}</p>
                {action.details.map((d, i) => (
                  <div key={i} className="flex justify-between text-sm border-t pt-2">
                    <span className="text-gray-400">{d.label}</span>
                    <span className="font-medium text-gray-900">{d.value}</span>
                  </div>
                ))}
              </div>

              <div className="bg-blue-50 rounded-lg p-3 text-xs text-blue-700">
                📝 سيتم تسجيل هذا الإجراء في سجل التدقيق مع بصمتك الحيوية وتفاصيل الجلسة
              </div>

              <button onClick={() => { setStep('scanning'); setTimeout(() => setStep('verified'), 2000) }}
                className="w-full h-14 bg-slate-800 hover:bg-slate-900 text-white rounded-xl font-bold text-lg flex items-center justify-center gap-2">
                <span className="text-2xl">🔑</span> أكد ببصمتك
              </button>
              <button onClick={onClose} className="w-full h-10 text-gray-400 hover:text-gray-600 text-sm">إلغاء</button>
            </>
          )}

          {step === 'scanning' && (
            <div className="text-center py-8">
              <div className="inline-flex items-center justify-center w-20 h-20 bg-blue-100 rounded-full mb-4 animate-pulse">
                <span className="text-4xl">🔑</span>
              </div>
              <p className="font-semibold text-gray-900 mb-1">جاري المسح...</p>
              <p className="text-sm text-gray-400">ضع إصبعك على مستشعر البصمة</p>
            </div>
          )}

          {step === 'verified' && (
            <div className="text-center py-8">
              <div className="inline-flex items-center justify-center w-20 h-20 bg-green-500 rounded-full mb-4">
                <span className="text-4xl">✅</span>
              </div>
              <p className="font-bold text-gray-900 mb-1">تم التأكيد!</p>
              <p className="text-sm text-gray-400">تم تنفيذ الإجراء وتسجيله في سجل التدقيق</p>
              <button onClick={onVerified} className="w-full h-12 mt-4 bg-green-600 hover:bg-green-700 text-white rounded-lg font-semibold">تم</button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

// ============ Audit Log Screen ============
function AuditScreen() {
  const [search, setSearch] = useState('')
  const [filterAction, setFilterAction] = useState('all')

  const filtered = mockAuditLogs.filter(l =>
    (filterAction === 'all' || l.action.includes(filterAction)) &&
    (!search || l.actor.includes(search) || l.entity.includes(search))
  )

  return (
    <div className="mr-60 p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">📜 سجل التدقيق</h1>
        <div className="flex items-center gap-2">
          <span className="text-xs text-green-600 bg-green-50 px-3 py-1 rounded-full">✅ Chain Verified</span>
          <span className="text-xs text-gray-400">1,247,392 سجل</span>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border p-4 flex gap-3">
        <input type="text" value={search} onChange={(e) => setSearch(e.target.value)} placeholder="🔍 بحث بالموظف أو الجهة..."
          className="flex-1 h-10 px-3 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400" />
        <select value={filterAction} onChange={(e) => setFilterAction(e.target.value)} className="h-10 px-3 rounded-lg border border-gray-200 text-sm">
          <option value="all">كل الإجراءات</option>
          <option value="refund">استرجاع</option>
          <option value="restaurant">مطعم</option>
          <option value="driver">مندوب</option>
          <option value="payout">دفعة</option>
          <option value="auth">مصادقة</option>
          <option value="order">طلب</option>
        </select>
        <button className="px-4 h-10 bg-gray-100 text-gray-600 rounded-lg text-sm font-semibold">📄 تصدير</button>
      </div>

      {/* Log table */}
      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="text-right p-3 font-semibold text-gray-600">الوقت</th>
              <th className="text-right p-3 font-semibold text-gray-600">الموظف</th>
              <th className="text-right p-3 font-semibold text-gray-600">الإجراء</th>
              <th className="text-right p-3 font-semibold text-gray-600">الجهة</th>
              <th className="text-right p-3 font-semibold text-gray-600">التفاصيل</th>
              <th className="text-right p-3 font-semibold text-gray-600">بصمة</th>
              <th className="text-right p-3 font-semibold text-gray-600">IP</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((log, i) => (
              <tr key={i} className="border-b last:border-0 hover:bg-gray-50">
                <td className="p-3 font-mono text-xs text-gray-400">{log.time}</td>
                <td className="p-3 font-medium text-gray-900">{log.actor}</td>
                <td className="p-3 font-mono text-xs text-blue-600">{log.action}</td>
                <td className="p-3 text-gray-600">{log.entity}</td>
                <td className="p-3 text-gray-400">{log.details}</td>
                <td className="p-3">{log.bio ? <span className="text-green-500">🔑</span> : <span className="text-gray-300">—</span>}</td>
                <td className="p-3 font-mono text-xs text-gray-400">{log.ip}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

// ============ Anomaly Screen ============
function AnomalyScreen() {
  return (
    <div className="mr-60 p-6 space-y-4">
      <h1 className="text-2xl font-bold">⚠️ تنبيهات الأنومالي (UEBA)</h1>
      <p className="text-sm text-gray-400">User & Entity Behavior Analytics — كشف السلوك غير الطبيعي للموظفين</p>

      <div className="space-y-3">
        {mockAnomalies.map(anom => (
          <div key={anom.id} className={cn('bg-white rounded-xl border-2 p-5',
            anom.severity === 'critical' ? 'border-red-300' : anom.severity === 'medium' ? 'border-amber-300' : 'border-blue-200')}>
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-3">
                <span className={cn('text-xs font-bold px-3 py-1 rounded-full',
                  anom.severity === 'critical' ? 'bg-red-100 text-red-600' : anom.severity === 'medium' ? 'bg-amber-100 text-amber-600' : 'bg-blue-100 text-blue-600')}>
                  {anom.severity === 'critical' ? '🔴 Critical' : anom.severity === 'medium' ? '🟡 Medium' : '🔵 Low'}
                </span>
                <div>
                  <p className="font-bold text-gray-900">{anom.employee} <span className="text-sm text-gray-400 font-normal">({anom.role})</span></p>
                  <p className="text-xs text-gray-400">Score: {anom.score}</p>
                </div>
              </div>
              <span className="text-xs text-gray-400">{anom.id}</span>
            </div>

            <div className="bg-gray-50 rounded-lg p-3 space-y-1 mb-3">
              <p className="text-xs font-semibold text-gray-600 mb-1">الأسباب:</p>
              {anom.reasons.map((r, i) => (
                <p key={i} className="text-sm text-gray-500">• {r}</p>
              ))}
            </div>

            <div className="flex gap-2">
              {anom.severity === 'critical' && (
                <button className="px-4 py-1.5 bg-red-500 text-white rounded-lg text-sm font-semibold hover:bg-red-600">🚫 تجميد الحساب</button>
              )}
              <button className="px-4 py-1.5 bg-blue-50 text-blue-600 rounded-lg text-sm font-semibold hover:bg-blue-100">📞 اتصال</button>
              <button className="px-4 py-1.5 bg-gray-50 text-gray-600 rounded-lg text-sm font-semibold hover:bg-gray-100">👁️ تحقيق</button>
              {anom.severity !== 'critical' && (
                <button className="px-4 py-1.5 bg-green-50 text-green-600 rounded-lg text-sm font-semibold hover:bg-green-100">✅ تجاهل</button>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Access Review Screen ============
function AccessReviewScreen() {
  const team = [
    { name: 'Sarah M.', role: 'Support L2', permissions: 12, lastActive: 'الآن', status: 'active' },
    { name: 'Ahmed K.', role: 'Support L2', permissions: 12, lastActive: 'الآن', status: 'active' },
    { name: 'Omar T.', role: 'Ops Manager', permissions: 25, lastActive: '5 دقايق', status: 'active' },
    { name: 'Fatma A.', role: 'Field Supervisor', permissions: 8, lastActive: '2 ساعة', status: 'active' },
    { name: 'Karim H.', role: 'Support L1', permissions: 5, lastActive: '30 يوم', status: 'suspended' },
    { name: 'Mariam F.', role: 'Finance', permissions: 15, lastActive: '1 ساعة', status: 'active' },
  ]

  return (
    <div className="mr-60 p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">🔑 مراجعة الصلاحيات</h1>
        <span className="text-sm text-gray-400">📅 مراجعة ربعية — Q3 2026</span>
      </div>

      <div className="bg-blue-50 rounded-lg p-4 text-sm text-blue-700">
        ℹ️ يجب على كل manager مراجعة صلاحيات فريقه وتأكيدها أو إلغائها كل 3 أشهر. الحسابات غير النشطة (+30 يوم) تتوقف تلقائياً.
      </div>

      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="text-right p-3 font-semibold text-gray-600">الموظف</th>
              <th className="text-right p-3 font-semibold text-gray-600">الدور</th>
              <th className="text-right p-3 font-semibold text-gray-600">الصلاحيات</th>
              <th className="text-right p-3 font-semibold text-gray-600">آخر نشاط</th>
              <th className="text-right p-3 font-semibold text-gray-600">الحالة</th>
              <th className="text-right p-3 font-semibold text-gray-600">إجراء</th>
            </tr>
          </thead>
          <tbody>
            {team.map((member, i) => (
              <tr key={i} className="border-b last:border-0 hover:bg-gray-50">
                <td className="p-3 font-medium text-gray-900">{member.name}</td>
                <td className="p-3 text-gray-600">{member.role}</td>
                <td className="p-3"><span className="font-bold text-gray-900">{member.permissions}</span> <span className="text-gray-400">صلاحية</span></td>
                <td className="p-3 text-gray-400">{member.lastActive}</td>
                <td className="p-3">
                  <span className={cn('text-xs px-2 py-1 rounded font-semibold',
                    member.status === 'active' ? 'bg-green-50 text-green-600' : 'bg-red-50 text-red-500')}>
                    {member.status === 'active' ? '🟢 نشط' : '🔴 موقوف'}
                  </span>
                </td>
                <td className="p-3">
                  <div className="flex gap-1">
                    <button className="px-2 py-1 bg-green-50 text-green-600 rounded text-xs">✅ إبقاء</button>
                    <button className="px-2 py-1 bg-red-50 text-red-500 rounded text-xs">❌ إلغاء</button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
        <p className="text-sm text-amber-700">⚠️ الموظف Karim H. غير نشط منذ 30 يوم — سيتم إيقافه تلقائياً</p>
      </div>
    </div>
  )
}

// ============ Helpers ============
function StatCard({ icon, label, value, color }: { icon: string; label: string; value: string; color: string }) {
  return (
    <div className="bg-white rounded-xl border p-4">
      <div className={cn('w-10 h-10 rounded-lg flex items-center justify-center mb-2 text-xl', color)}>{icon}</div>
      <p className="text-xs text-gray-400">{label}</p>
      <p className="text-2xl font-bold text-gray-900 mt-1">{value}</p>
    </div>
  )
}

function QuickAction({ icon, label, onClick }: { icon: string; label: string; onClick: () => void }) {
  return (
    <button onClick={onClick} className="bg-white rounded-xl border p-4 hover:shadow-md transition-shadow flex flex-col items-center gap-2">
      <span className="text-2xl">{icon}</span>
      <span className="text-sm font-medium text-gray-700">{label}</span>
    </button>
  )
}
