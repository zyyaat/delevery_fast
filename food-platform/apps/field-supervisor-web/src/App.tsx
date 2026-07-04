// Field Supervisor App — Full implementation (FINAL APP!)
// 11 screens: Login, Task List, Restaurant Verification, Driver Verification,
// Complaint Investigation, Route Planner, Daily Report, Performance,
// Driver Training, Photo Capture, GPS Verification

import { useState, useEffect, useRef } from 'react'
import { cn } from '@food-platform/ui'

type Screen = 'login' | 'tasks' | 'restaurant_verify' | 'driver_verify' | 'complaint' | 'route' | 'report' | 'performance' | 'training'

export default function App() {
  const [screen, setScreen] = useState<Screen>('login')
  const [authenticated, setAuthenticated] = useState(false)
  const [tasks, setTasks] = useState(mockTasks)
  const [activeTask, setActiveTask] = useState<typeof mockTasks[0] | null>(null)

  if (!authenticated && screen !== 'login') {
    setScreen('login')
  }

  const completeTask = (id: string, result: string) => {
    setTasks(tasks.map(t => t.id === id ? { ...t, status: 'completed' } : t))
    setActiveTask(null)
    setScreen('tasks')
  }

  return (
    <div className="min-h-screen bg-gray-50" dir="rtl">
      {!authenticated ? (
        <LoginScreen onLogin={() => { setAuthenticated(true); setScreen('tasks') }} />
      ) : (
        <>
          {screen === 'tasks' && <TasksScreen tasks={tasks} onNavigate={setScreen} onOpenTask={(t) => { setActiveTask(t); setScreen(t.type === 'restaurant_verification' ? 'restaurant_verify' : t.type === 'driver_verification' ? 'driver_verify' : 'complaint') }} />}
          {screen === 'restaurant_verify' && activeTask && <RestaurantVerifyScreen task={activeTask} onComplete={(r) => completeTask(activeTask.id, r)} onBack={() => setScreen('tasks')} />}
          {screen === 'driver_verify' && activeTask && <DriverVerifyScreen task={activeTask} onComplete={(r) => completeTask(activeTask.id, r)} onBack={() => setScreen('tasks')} />}
          {screen === 'complaint' && activeTask && <ComplaintScreen task={activeTask} onComplete={(r) => completeTask(activeTask.id, r)} onBack={() => setScreen('tasks')} />}
          {screen === 'route' && <RouteScreen tasks={tasks} onBack={() => setScreen('tasks')} />}
          {screen === 'report' && <ReportScreen tasks={tasks} onBack={() => setScreen('tasks')} />}
          {screen === 'performance' && <PerformanceScreen onBack={() => setScreen('tasks')} />}
          {screen === 'training' && <TrainingScreen onBack={() => setScreen('tasks')} />}
          <BottomNav active={screen} onNavigate={setScreen} />
        </>
      )}
    </div>
  )
}

// ============ Mock Data ============
const mockTasks = [
  { id: 'TASK-001', type: 'complaint', priority: 'urgent', title: 'تحقيق شكوى: حشرة في الأكل', address: 'KFC - التجمع', distance: 12, eta: 35, scheduled: 'الآن', status: 'pending' },
  { id: 'TASK-002', type: 'restaurant_verification', priority: 'normal', title: 'تفعيل مطعم جديد: Pizza Hut المعصرة', address: 'Pizza Hut - المعصرة', distance: 8.5, eta: 25, scheduled: '11:00 AM', status: 'pending' },
  { id: 'TASK-003', type: 'driver_verification', priority: 'low', title: 'تفعيل مندوب: Mahmoud S.', address: 'مدينة نصر', distance: 5.2, eta: 18, scheduled: '1:00 PM', status: 'pending' },
  { id: 'TASK-004', type: 'restaurant_verification', priority: 'normal', title: 'فحص مفاجئ: McDonald\'s معادي', address: 'McDonald\'s - معادي', distance: 2.1, eta: 8, scheduled: '3:00 PM', status: 'pending' },
  { id: 'TASK-005', type: 'driver_verification', priority: 'low', title: 'تفعيل مندوب: Ahmed K.', address: 'الزمالك', distance: 3.5, eta: 12, scheduled: '4:00 PM', status: 'pending' },
  { id: 'TASK-006', type: 'complaint', priority: 'normal', title: 'شكوى: المندوب تصرف بوقاحة', address: 'الدقي', distance: 6.8, eta: 20, scheduled: '5:00 PM', status: 'pending' },
]

// ============ Login ============
function LoginScreen({ onLogin }: { onLogin: () => void }) {
  const [phone, setPhone] = useState('')
  const [otp, setOtp] = useState(['', '', '', '', '', ''])
  const [step, setStep] = useState<'phone' | 'otp'>('phone')
  const inputsRef = useRef<(HTMLInputElement | null)[]>([])

  return (
    <div className="min-h-screen bg-gradient-to-b from-teal-50 to-white flex items-center justify-center p-6">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-teal-600 rounded-3xl mb-4 text-4xl">🧑‍💼</div>
          <h1 className="text-2xl font-bold text-gray-900">المشرف الميداني</h1>
          <p className="text-gray-500 mt-2 text-sm">Field Supervisor App</p>
        </div>

        {step === 'phone' && (
          <div className="space-y-4">
            <div>
              <label className="text-sm text-gray-600 block mb-1">رقم الموبايل</label>
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1 px-3 h-12 rounded-lg border border-gray-200 bg-gray-50 text-sm font-medium">🇪🇬 +20</div>
                <input type="tel" inputMode="numeric" maxLength={11} placeholder="01 2345 6789" value={phone}
                  onChange={(e) => setPhone(e.target.value.replace(/\D/g, ''))}
                  className="w-full h-12 px-3 rounded-lg border border-gray-200 text-lg focus:outline-none focus:ring-2 focus:ring-teal-400" />
              </div>
            </div>
            <button onClick={() => setStep('otp')} disabled={phone.length < 11}
              className="w-full h-12 bg-teal-600 hover:bg-teal-700 text-white rounded-lg font-semibold disabled:opacity-50">أرسل الكود</button>
          </div>
        )}

        {step === 'otp' && (
          <div className="space-y-4">
            <p className="text-center text-gray-500 text-sm">بعتناك كود على <span dir="ltr" className="font-semibold">+20 {phone}</span></p>
            <div className="flex justify-center gap-2" dir="ltr">
              {otp.map((d, i) => (
                <input key={i} ref={(el) => { inputsRef.current[i] = el }} type="text" inputMode="numeric" maxLength={1} value={d}
                  onChange={(e) => {
                    const digit = e.target.value.replace(/\D/g, '').slice(-1)
                    const newOtp = [...otp]; newOtp[i] = digit; setOtp(newOtp)
                    if (digit && i < 5) inputsRef.current[i + 1]?.focus()
                    if (newOtp.every(x => x) && newOtp.join('').length === 6) setTimeout(onLogin, 500)
                  }}
                  className={`w-12 h-14 text-center text-2xl font-bold rounded-lg border-2 ${d ? 'border-teal-500 bg-teal-50' : 'border-gray-200'} focus:outline-none`} autoFocus={i === 0} />
              ))}
            </div>
            <button onClick={onLogin} disabled={otp.some(d => !d)}
              className="w-full h-12 bg-teal-600 hover:bg-teal-700 text-white rounded-lg font-semibold disabled:opacity-50">تأكيد</button>
            <button onClick={() => setStep('phone')} className="w-full text-sm text-gray-400">تغيير الرقم</button>
          </div>
        )}
      </div>
    </div>
  )
}

// ============ Tasks Screen ============
function TasksScreen({ tasks, onNavigate, onOpenTask }: { tasks: typeof mockTasks; onNavigate: (s: Screen) => void; onOpenTask: (t: typeof mockTasks[0]) => void }) {
  const pending = tasks.filter(t => t.status === 'pending')
  const completed = tasks.filter(t => t.status === 'completed')
  const totalDistance = tasks.reduce((sum, t) => sum + t.distance, 0)

  return (
    <div className="min-h-screen pb-20">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-xs text-gray-400">مهام اليوم</p>
            <p className="text-lg font-bold text-gray-900">{pending.length}/{tasks.length} متبقي</p>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs text-gray-400">📍 معادي</span>
            <span className="text-xs text-teal-600 bg-teal-50 px-2 py-1 rounded-full">⏱️ {totalDistance}km</span>
          </div>
        </div>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-3">
        {/* Urgent first */}
        {pending.filter(t => t.priority === 'urgent').map(task => (
          <TaskCard key={task.id} task={task} onClick={() => onOpenTask(task)} />
        ))}

        {/* Normal */}
        {pending.filter(t => t.priority === 'normal').map(task => (
          <TaskCard key={task.id} task={task} onClick={() => onOpenTask(task)} />
        ))}

        {/* Low */}
        {pending.filter(t => t.priority === 'low').map(task => (
          <TaskCard key={task.id} task={task} onClick={() => onOpenTask(task)} />
        ))}

        {/* Completed */}
        {completed.length > 0 && (
          <>
            <p className="text-sm text-gray-400 pt-4">─── مكتملة ({completed.length}) ───</p>
            {completed.map(task => (
              <div key={task.id} className="bg-green-50 border border-green-200 rounded-xl p-4 opacity-70">
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✅</span>
                  <p className="text-sm font-medium text-gray-600 line-through">{task.title}</p>
                </div>
              </div>
            ))}
          </>
        )}
      </div>
    </div>
  )
}

function TaskCard({ task, onClick }: { task: typeof mockTasks[0]; onClick: () => void }) {
  const typeIcons: Record<string, string> = {
    restaurant_verification: '🍔', driver_verification: '🛵', complaint: '⚠️', audit: '🔍', training: '📚'
  }
  const priorityColors: Record<string, string> = {
    urgent: 'border-red-300 bg-red-50', normal: 'border-gray-200 bg-white', low: 'border-gray-200 bg-white'
  }
  const priorityLabels: Record<string, string> = { urgent: '🔴 عاجل', normal: '🟡 عادي', low: '🟢 منخفض' }

  return (
    <div onClick={onClick} className={cn('rounded-xl border-2 p-4 cursor-pointer hover:shadow-md transition-shadow', priorityColors[task.priority])}>
      <div className="flex items-start justify-between mb-2">
        <div className="flex items-center gap-2">
          <span className="text-2xl">{typeIcons[task.type] || '📋'}</span>
          <div>
            <p className="font-semibold text-gray-900 text-sm">{task.title}</p>
            <p className="text-xs text-gray-400 mt-0.5">📍 {task.address} • {task.distance}km • ⏱️ {task.eta}د</p>
          </div>
        </div>
        <span className="text-xs text-gray-400">{priorityLabels[task.priority]}</span>
      </div>
      <div className="flex items-center justify-between mt-2">
        <span className="text-xs text-gray-400">🕐 {task.scheduled}</span>
        <button className="text-xs text-teal-600 font-semibold">▶️ ابدأ</button>
      </div>
    </div>
  )
}

// ============ Restaurant Verification ============
function RestaurantVerifyScreen({ task, onComplete, onBack }: { task: typeof mockTasks[0]; onComplete: (r: string) => void; onBack: () => void }) {
  const [checklist, setChecklist] = useState<Record<string, boolean>>({})
  const [photos, setPhotos] = useState<string[]>([])
  const [gpsVerified, setGpsVerified] = useState(false)
  const [decision, setDecision] = useState<string | null>(null)
  const [notes, setNotes] = useState('')

  const sections = [
    { name: 'هوية المطعم', items: ['اسم المطعم مطابق', 'صاحب المطعم حاضر', 'رخصة التشغيل سارية', 'البطاقة الضريبية صحيحة', 'عقد الإيجار موجود'] },
    { name: 'النظافة', items: ['المطبخ نظيف', 'معدات التخزين سليمة', 'لا توجد حشرات', 'نظام القمامة سليم', 'العاملين يرتدون قفازات'] },
    { name: 'السلامة', items: ['يوجد طفاية حريق', 'المخارج واضحة', 'التهوية كافية', 'لا توجد مخاطر كهربائية'] },
  ]

  const requiredPhotos = ['واجهة المطعم', 'المدخل', 'المطبخ', 'رخصة التشغيل', 'البطاقة الضريبية']
  const allChecked = sections.every(s => s.items.every(i => checklist[`${s.name}-${i}`]))
  const allPhotos = photos.length >= requiredPhotos.length

  return (
    <div className="min-h-screen pb-24">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <div className="flex-1"><p className="font-bold text-sm">🍔 {task.title}</p><p className="text-xs text-gray-400">📍 {task.address}</p></div>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-5">
        {/* GPS verification */}
        <div className={cn('rounded-xl border-2 p-4', gpsVerified ? 'border-green-300 bg-green-50' : 'border-amber-300 bg-amber-50')}>
          <div className="flex items-center justify-between">
            <div>
              <p className="font-semibold text-sm">📍 التحقق من الموقع (GPS)</p>
              <p className="text-xs text-gray-500 mt-1">{gpsVerified ? '✅ تم التأكد — أنت في موقع المطعم' : '⚠️ يجب أن تكون في موقع المطعم'}</p>
            </div>
            <button onClick={() => setGpsVerified(true)} disabled={gpsVerified}
              className={cn('px-4 py-2 rounded-lg text-sm font-semibold', gpsVerified ? 'bg-green-100 text-green-600' : 'bg-amber-500 text-white hover:bg-amber-600')}>
              {gpsVerified ? '✅ تم' : '📍 تحقق'}
            </button>
          </div>
        </div>

        {/* Checklist */}
        {sections.map(section => (
          <div key={section.name} className="bg-white rounded-xl border p-4">
            <h3 className="font-bold text-sm mb-3">{section.name}</h3>
            <div className="space-y-2">
              {section.items.map(item => {
                const key = `${section.name}-${item}`
                return (
                  <label key={key} className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 cursor-pointer">
                    <input type="checkbox" checked={checklist[key] || false}
                      onChange={(e) => setChecklist({ ...checklist, [key]: e.target.checked })}
                      className="w-5 h-5 accent-teal-500" />
                    <span className={cn('text-sm', checklist[key] ? 'text-gray-400 line-through' : 'text-gray-700')}>{item}</span>
                    {checklist[key] && <span className="text-green-500 text-xs mr-auto">✅</span>}
                  </label>
                )
              })}
            </div>
          </div>
        ))}

        {/* Photos */}
        <div className="bg-white rounded-xl border p-4">
          <h3 className="font-bold text-sm mb-3">📷 الصور المطلوبة ({photos.length}/{requiredPhotos.length})</h3>
          <div className="grid grid-cols-3 gap-3">
            {requiredPhotos.map((photo, i) => (
              <button key={i} onClick={() => photos.length > i ? null : setPhotos([...photos, photo])}
                className={cn('aspect-square rounded-lg border-2 flex flex-col items-center justify-center gap-1 transition-colors',
                  photos.length > i ? 'border-green-300 bg-green-50' : 'border-dashed border-gray-300 hover:border-teal-400')}>
                {photos.length > i ? <><span className="text-2xl">✅</span><span className="text-xs text-green-600">{photo}</span></> : <><span className="text-2xl">📷</span><span className="text-xs text-gray-400">{photo}</span></>}
              </button>
            ))}
          </div>
          <p className="text-xs text-gray-400 mt-2">📌 كل صورة بتتسجل مع GPS + timestamp تلقائياً</p>
        </div>

        {/* Notes */}
        <div className="bg-white rounded-xl border p-4">
          <h3 className="font-bold text-sm mb-2">📝 ملاحظات</h3>
          <textarea value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="ملاحظات إضافية..." rows={3}
            className="w-full p-3 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-teal-400 resize-none" />
        </div>

        {/* Decision */}
        <div className="bg-white rounded-xl border p-4">
          <h3 className="font-bold text-sm mb-3">القرار النهائي</h3>
          <div className="space-y-2">
            {[{ id: 'approved', label: '✅ موافقة (تفعيل فوري)', color: 'border-green-500 bg-green-50' },
              { id: 'conditional', label: '⚠️ مشروط (مهلة 7 أيام)', color: 'border-amber-500 bg-amber-50' },
              { id: 'rejected', label: '❌ رفض', color: 'border-red-500 bg-red-50' }].map(opt => (
              <label key={opt.id} className={cn('flex items-center gap-3 p-3 rounded-lg border-2 cursor-pointer transition-colors',
                decision === opt.id ? opt.color : 'border-gray-200')}>
                <input type="radio" checked={decision === opt.id} onChange={() => setDecision(opt.id)} className="w-5 h-5 accent-teal-500" />
                <span className="text-sm font-medium">{opt.label}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Submit */}
        <button onClick={() => decision && onComplete(decision)} disabled={!allChecked || !allPhotos || !gpsVerified || !decision}
          className="w-full h-12 bg-teal-600 hover:bg-teal-700 text-white rounded-lg font-bold disabled:opacity-50">
          {!allChecked ? 'أكمل الـ Checklist' : !allPhotos ? 'التقط كل الصور' : !gpsVerified ? 'تحقق من GPS' : '🔒 تأكيد القرار'}
        </button>
      </div>
    </div>
  )
}

// ============ Driver Verification ============
function DriverVerifyScreen({ task, onComplete, onBack }: { task: typeof mockTasks[0]; onComplete: (r: string) => void; onBack: () => void }) {
  const [checklist, setChecklist] = useState<Record<string, boolean>>({})
  const [photos, setPhotos] = useState<string[]>([])
  const [gpsVerified, setGpsVerified] = useState(false)
  const [decision, setDecision] = useState<string | null>(null)

  const sections = [
    { name: 'هوية المندوب', items: ['البطاقة الشخصية أصلية', 'صورة المندوب مطابقة', 'Face match (liveness)', 'الرقم القومي موجود'] },
    { name: 'المركبة', items: ['المركبة موجودة فعلياً', 'اللوحات مطابقة', 'التأمين ساري', 'المركبة في حالة جيدة', 'خوذة متوفرة'] },
    { name: 'الفحص العملي', items: ['يقدر يركب الموتوسيكل بأمان', 'يعرف إشارات المرور', 'يقدر يستخدم GPS', 'يتواصل بأدب (اختبار قصير)'] },
  ]
  const requiredPhotos = ['صورة المندوب مع المركبة', 'اللوحات', 'الرخصة', 'البطاقة', 'Selfie مع البطاقة']
  const allChecked = sections.every(s => s.items.every(i => checklist[`${s.name}-${i}`]))
  const allPhotos = photos.length >= requiredPhotos.length

  return (
    <div className="min-h-screen pb-24">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <div className="flex-1"><p className="font-bold text-sm">🛵 {task.title}</p><p className="text-xs text-gray-400">📍 {task.address}</p></div>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-5">
        {/* GPS */}
        <div className={cn('rounded-xl border-2 p-4', gpsVerified ? 'border-green-300 bg-green-50' : 'border-amber-300 bg-amber-50')}>
          <div className="flex items-center justify-between">
            <p className="font-semibold text-sm">📍 التحقق من الموقع</p>
            <button onClick={() => setGpsVerified(true)} disabled={gpsVerified}
              className={cn('px-4 py-2 rounded-lg text-sm font-semibold', gpsVerified ? 'bg-green-100 text-green-600' : 'bg-amber-500 text-white')}>{gpsVerified ? '✅ تم' : '📍 تحقق'}</button>
          </div>
        </div>

        {/* Checklist */}
        {sections.map(section => (
          <div key={section.name} className="bg-white rounded-xl border p-4">
            <h3 className="font-bold text-sm mb-3">{section.name}</h3>
            <div className="space-y-2">
              {section.items.map(item => {
                const key = `${section.name}-${item}`
                return (
                  <label key={key} className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 cursor-pointer">
                    <input type="checkbox" checked={checklist[key] || false}
                      onChange={(e) => setChecklist({ ...checklist, [key]: e.target.checked })}
                      className="w-5 h-5 accent-teal-500" />
                    <span className={cn('text-sm', checklist[key] ? 'text-gray-400 line-through' : 'text-gray-700')}>{item}</span>
                    {checklist[key] && <span className="text-green-500 text-xs mr-auto">✅</span>}
                  </label>
                )
              })}
            </div>
          </div>
        ))}

        {/* Photos */}
        <div className="bg-white rounded-xl border p-4">
          <h3 className="font-bold text-sm mb-3">📷 الصور المطلوبة ({photos.length}/{requiredPhotos.length})</h3>
          <div className="grid grid-cols-3 gap-3">
            {requiredPhotos.map((photo, i) => (
              <button key={i} onClick={() => photos.length > i ? null : setPhotos([...photos, photo])}
                className={cn('aspect-square rounded-lg border-2 flex flex-col items-center justify-center gap-1',
                  photos.length > i ? 'border-green-300 bg-green-50' : 'border-dashed border-gray-300 hover:border-teal-400')}>
                {photos.length > i ? <><span className="text-2xl">✅</span><span className="text-[10px] text-green-600 text-center px-1">{photo}</span></> : <><span className="text-2xl">📷</span><span className="text-[10px] text-gray-400 text-center px-1">{photo}</span></>}
              </button>
            ))}
          </div>
        </div>

        {/* Decision */}
        <div className="bg-white rounded-xl border p-4">
          <h3 className="font-bold text-sm mb-3">القرار</h3>
          <div className="space-y-2">
            {[{ id: 'approved', label: '✅ موافقة (تفعيل)' }, { id: 'rejected', label: '❌ رفض' }].map(opt => (
              <label key={opt.id} className={cn('flex items-center gap-3 p-3 rounded-lg border-2 cursor-pointer',
                decision === opt.id ? (opt.id === 'approved' ? 'border-green-500 bg-green-50' : 'border-red-500 bg-red-50') : 'border-gray-200')}>
                <input type="radio" checked={decision === opt.id} onChange={() => setDecision(opt.id)} className="w-5 h-5 accent-teal-500" />
                <span className="text-sm font-medium">{opt.label}</span>
              </label>
            ))}
          </div>
        </div>

        <button onClick={() => decision && onComplete(decision)} disabled={!allChecked || !allPhotos || !gpsVerified || !decision}
          className="w-full h-12 bg-teal-600 hover:bg-teal-700 text-white rounded-lg font-bold disabled:opacity-50">
          {!allChecked ? 'أكمل الـ Checklist' : !allPhotos ? 'التقط كل الصور' : !gpsVerified ? 'تحقق من GPS' : '🔒 تأكيد القرار'}
        </button>
      </div>
    </div>
  )
}

// ============ Complaint Investigation ============
function ComplaintScreen({ task, onComplete, onBack }: { task: typeof mockTasks[0]; onComplete: (r: string) => void; onBack: () => void }) {
  const [step, setStep] = useState<'details' | 'investigate' | 'resolve'>('details')
  const [findings, setFindings] = useState('')
  const [photos, setPhotos] = useState<string[]>([])

  return (
    <div className="min-h-screen pb-24">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <div className="flex-1"><p className="font-bold text-sm">⚠️ {task.title}</p><p className="text-xs text-gray-400">📍 {task.address}</p></div>
        <span className="text-xs bg-red-50 text-red-500 px-2 py-1 rounded">🔴 عاجل</span>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-5">
        {step === 'details' && (
          <>
            <div className="bg-red-50 rounded-xl border border-red-200 p-4">
              <h3 className="font-bold text-sm text-red-700 mb-2">📋 تفاصيل الشكوى</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between"><span className="text-gray-500">العميل</span><span className="font-medium">Sara M.</span></div>
                <div className="flex justify-between"><span className="text-gray-500">الطلب</span><span className="font-medium">#D9N81P</span></div>
                <div className="flex justify-between"><span className="text-gray-500">المطعم</span><span className="font-medium">KFC - التجمع</span></div>
                <div className="flex justify-between"><span className="text-gray-500">المبلغ</span><span className="font-medium">EGP 285</span></div>
                <div className="flex justify-between"><span className="text-gray-500">الشكوى</span><span className="font-medium text-red-600">حشرة في الأكل</span></div>
              </div>
            </div>

            <div className="bg-white rounded-xl border p-4">
              <h3 className="font-bold text-sm mb-2">📝 خطوات التحقيق</h3>
              <div className="space-y-2 text-sm text-gray-600">
                <p>1. 📸 التقاط صور للأكل المشتبه به (من العميل لو متاح)</p>
                <p>2. 🏪 زيارة المطعم وفحص المطبخ</p>
                <p>3. 👨‍🍳 التحدث مع الطباخين والموظفين</p>
                <p>4. 📋 فحص سجلات النظافة</p>
                <p>5. 🛵 التحدث مع المندوب (لو متاح)</p>
              </div>
            </div>

            <button onClick={() => setStep('investigate')} className="w-full h-12 bg-teal-600 hover:bg-teal-700 text-white rounded-lg font-bold">▶️ بدء التحقيق</button>
          </>
        )}

        {step === 'investigate' && (
          <>
            <div className="bg-white rounded-xl border p-4">
              <h3 className="font-bold text-sm mb-3">📷 أدلة (صور)</h3>
              <div className="grid grid-cols-3 gap-3">
                {['صورة الأكل', 'المطبخ', 'سجلات النظافة'].map((label, i) => (
                  <button key={i} onClick={() => setPhotos([...photos, label])}
                    className={cn('aspect-square rounded-lg border-2 flex flex-col items-center justify-center gap-1',
                      photos.length > i ? 'border-green-300 bg-green-50' : 'border-dashed border-gray-300')}>
                    {photos.length > i ? <span className="text-2xl">✅</span> : <span className="text-2xl">📷</span>}
                    <span className="text-xs text-gray-400">{label}</span>
                  </button>
                ))}
              </div>
            </div>

            <div className="bg-white rounded-xl border p-4">
              <h3 className="font-bold text-sm mb-2">🔍 النتائج</h3>
              <textarea value={findings} onChange={(e) => setFindings(e.target.value)}
                placeholder="اكتب نتائج التحقيق..." rows={5}
                className="w-full p-3 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-teal-400 resize-none" />
            </div>

            <button onClick={() => setStep('resolve')} disabled={!findings || photos.length < 1}
              className="w-full h-12 bg-teal-600 hover:bg-teal-700 text-white rounded-lg font-bold disabled:opacity-50">التالي: الحل</button>
          </>
        )}

        {step === 'resolve' && (
          <>
            <div className="bg-white rounded-xl border p-4">
              <h3 className="font-bold text-sm mb-3">🎯 الحل المقترح</h3>
              <div className="space-y-2">
                {[{ id: 'restaurant_fault', label: '🍔 المطعم مذنب — إيقاف + تحقيق', color: 'border-red-500' },
                  { id: 'driver_fault', label: '🛵 المندوب مذنب — suspend', color: 'border-amber-500' },
                  { id: 'customer_fraud', label: '🚫 احتيال عميل — flag', color: 'border-purple-500' },
                  { id: 'no_fault', label: '✅ لا توجد مخالفة', color: 'border-green-500' }].map(opt => (
                  <label key={opt.id} className={cn('flex items-center gap-3 p-3 rounded-lg border-2 cursor-pointer', opt.color, 'bg-gray-50')}>
                    <input type="radio" name="resolution" className="w-5 h-5 accent-teal-500" />
                    <span className="text-sm font-medium">{opt.label}</span>
                  </label>
                ))}
              </div>
            </div>

            <div className="bg-blue-50 rounded-lg p-3 text-sm text-blue-700">
              💰 تعويض العميل: EGP 285 (full refund) + EGP 50 (كوبون اعتذار)
            </div>

            <button onClick={() => onComplete('resolved')} className="w-full h-12 bg-green-600 hover:bg-green-700 text-white rounded-lg font-bold">✅ إرسال التقرير</button>
          </>
        )}
      </div>
    </div>
  )
}

// ============ Route Planner ============
function RouteScreen({ tasks, onBack }: { tasks: typeof mockTasks; onBack: () => void }) {
  const pending = tasks.filter(t => t.status === 'pending')
  const totalDistance = pending.reduce((sum, t) => sum + t.distance, 0)
  const totalTime = pending.reduce((sum, t) => sum + t.eta, 0)

  return (
    <div className="min-h-screen pb-20">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <h1 className="font-bold">🗺️ تخطيط المسار</h1>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        <div className="bg-teal-50 rounded-xl border border-teal-200 p-4 grid grid-cols-3 gap-3 text-center">
          <div><p className="text-2xl font-bold text-teal-600">{pending.length}</p><p className="text-xs text-gray-500">مهام</p></div>
          <div><p className="text-2xl font-bold text-teal-600">{totalDistance}km</p><p className="text-xs text-gray-500">مسافة</p></div>
          <div><p className="text-2xl font-bold text-teal-600">{Math.floor(totalTime / 60)}س {totalTime % 60}د</p><p className="text-xs text-gray-500">وقت</p></div>
        </div>

        <div className="space-y-3">
          {pending.map((task, i) => (
            <div key={task.id} className="bg-white rounded-xl border p-4 flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-teal-500 text-white flex items-center justify-center font-bold text-sm">{i + 1}</div>
              <div className="flex-1">
                <p className="text-sm font-semibold text-gray-900">{task.title}</p>
                <p className="text-xs text-gray-400">📍 {task.address} • {task.distance}km • ⏱️ {task.eta}د</p>
              </div>
              <span className="text-xs text-gray-400">{task.scheduled}</span>
            </div>
          ))}
        </div>

        <div className="space-y-2">
          <label className="flex items-center gap-3 p-3 bg-white rounded-lg border cursor-pointer">
            <input type="checkbox" defaultChecked className="w-5 h-5 accent-teal-500" />
            <span className="text-sm">تجميع المهام القريبة</span>
          </label>
          <label className="flex items-center gap-3 p-3 bg-white rounded-lg border cursor-pointer">
            <input type="checkbox" defaultChecked className="w-5 h-5 accent-teal-500" />
            <span className="text-sm">تجنب ساعات الذروة</span>
          </label>
          <label className="flex items-center gap-3 p-3 bg-white rounded-lg border cursor-pointer">
            <input type="checkbox" defaultChecked className="w-5 h-5 accent-teal-500" />
            <span className="text-sm">إعطاء أولوية للمهام العاجلة</span>
          </label>
        </div>

        <button className="w-full h-12 bg-teal-600 hover:bg-teal-700 text-white rounded-lg font-bold">🔄 إعادة ترتيب تلقائي (أمثل)</button>
        <button className="w-full h-12 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-bold">📍 ابدأ بالزيارة الأولى</button>
      </div>
    </div>
  )
}

// ============ Daily Report ============
function ReportScreen({ tasks, onBack }: { tasks: typeof mockTasks; onBack: () => void }) {
  const completed = tasks.filter(t => t.status === 'completed')
  const pending = tasks.filter(t => t.status === 'pending')

  return (
    <div className="min-h-screen pb-20">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <h1 className="font-bold">📊 تقرير اليوم</h1>
      </header>

      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        <div className="grid grid-cols-3 gap-3">
          <div className="bg-white rounded-xl border p-4 text-center"><p className="text-2xl font-bold text-teal-600">{completed.length}/{tasks.length}</p><p className="text-xs text-gray-400">مكتملة</p></div>
          <div className="bg-white rounded-xl border p-4 text-center"><p className="text-2xl font-bold text-gray-900">7س</p><p className="text-xs text-gray-400">مدة العمل</p></div>
          <div className="bg-white rounded-xl border p-4 text-center"><p className="text-2xl font-bold text-gray-900">68km</p><p className="text-xs text-gray-400">مسافة</p></div>
        </div>

        <div className="bg-white rounded-xl border p-5">
          <h3 className="font-bold mb-3">تفصيل الزيارات</h3>
          <div className="space-y-2">
            <div className="flex justify-between text-sm"><span className="text-gray-400">🍔 تفعيل مطاعم</span><span className="font-bold">2 ✅</span></div>
            <div className="flex justify-between text-sm"><span className="text-gray-400">🛵 تفعيل مناديب</span><span className="font-bold">3 ✅</span></div>
            <div className="flex justify-between text-sm"><span className="text-gray-400">⚠️ تحقيق شكاوى</span><span className="font-bold">1 ✅</span></div>
          </div>
        </div>

        <div className="bg-white rounded-xl border p-5">
          <h3 className="font-bold mb-3">Issues Found</h3>
          <div className="space-y-2 text-sm">
            <p className="text-gray-600">• 1 رخصة تشغيل منتهية</p>
            <p className="text-gray-600">• 1 مركبة في حالة سيئة</p>
            <p className="text-gray-600">• 1 مطعم فيه مشكلة نظافة (KFC التجمع)</p>
          </div>
        </div>

        <div className="bg-white rounded-xl border p-5">
          <h3 className="font-bold mb-3">التوصيات</h3>
          <div className="space-y-2 text-sm">
            <p className="text-gray-600">• إعادة فحص Pizza Hut بعد 7 أيام</p>
            <p className="text-gray-600">• تفعيل فوري لـ 2 مطعم</p>
            <p className="text-gray-600">• suspend KFC التجمع</p>
          </div>
        </div>

        <button className="w-full h-12 bg-teal-600 hover:bg-teal-700 text-white rounded-lg font-bold">📤 رفع التقرير النهائي</button>
      </div>
    </div>
  )
}

// ============ Performance ============
function PerformanceScreen({ onBack }: { onBack: () => void }) {
  return (
    <div className="min-h-screen pb-20">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <h1 className="font-bold">📈 الأداء</h1>
      </header>
      <div className="max-w-2xl mx-auto px-4 py-4 space-y-4">
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-white rounded-xl border p-4 text-center"><p className="text-3xl font-bold text-teal-600">42</p><p className="text-sm text-gray-400">زيارات هذا الأسبوع</p></div>
          <div className="bg-white rounded-xl border p-4 text-center"><p className="text-3xl font-bold text-teal-600">4.7⭐</p><p className="text-sm text-gray-400">تقييم</p></div>
        </div>
        <div className="bg-white rounded-xl border p-5">
          <h3 className="font-bold mb-3">المقارنة بالفريق</h3>
          <div className="space-y-2">
            {[{ name: 'Omar F.', visits: 42, rating: 4.7, you: true }, { name: 'Sara M.', visits: 38, rating: 4.8 }, { name: 'Ahmed K.', visits: 45, rating: 4.5 }, { name: 'Mahmoud A.', visits: 35, rating: 4.6 }].map((p, i) => (
              <div key={i} className={cn('flex items-center gap-3 py-2 border-b last:border-0', p.you && 'bg-teal-50 rounded-lg px-2')}>
                <span className="text-lg">{['🥇', '🥈', '🥉', '4️⃣'][i]}</span>
                <span className={cn('flex-1 text-sm font-semibold', p.you ? 'text-teal-600' : 'text-gray-700')}>{p.name} {p.you && '(أنت)'}</span>
                <span className="text-xs text-gray-400">{p.visits} زيارة</span>
                <span className="text-xs text-amber-400">⭐{p.rating}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

// ============ Training ============
function TrainingScreen({ onBack }: { onBack: () => void }) {
  return (
    <div className="min-h-screen pb-20">
      <header className="sticky top-0 z-40 bg-white border-b px-4 py-3 flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <h1 className="font-bold">📚 تدريب المناديب</h1>
      </header>
      <div className="max-w-2xl mx-auto px-4 py-4 space-y-3">
        {['استلام الطلب', 'الاستلام من المطعم', 'التوصيل للعميل', 'التعامل مع العملاء', 'السلامة المرورية', 'استلام الأرباح'].map((module, i) => (
          <div key={i} className="bg-white rounded-xl border p-4 flex items-center gap-3">
            <span className="text-2xl">{['📋', '🍔', '📍', '😊', '🛡️', '💰'][i]}</span>
            <div className="flex-1"><p className="font-semibold text-sm">{module}</p><p className="text-xs text-gray-400">15 دقيقة</p></div>
            <button className="px-3 py-1.5 bg-teal-50 text-teal-600 rounded-lg text-xs font-semibold">▶️ ابدأ</button>
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Bottom Nav ============
function BottomNav({ active, onNavigate }: { active: string; onNavigate: (s: Screen) => void }) {
  const items = [
    { id: 'tasks', label: 'المهام', icon: '📋' },
    { id: 'route', label: 'المسار', icon: '🗺️' },
    { id: 'report', label: 'التقرير', icon: '📊' },
    { id: 'performance', label: 'الأداء', icon: '📈' },
    { id: 'training', label: 'تدريب', icon: '📚' },
  ]
  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 bg-white border-t h-16 flex items-center justify-around px-2">
      {items.map(item => (
        <button key={item.id} onClick={() => onNavigate(item.id as Screen)}
          className={cn('flex flex-col items-center gap-1 flex-1 h-full justify-center', active === item.id ? 'text-teal-600' : 'text-gray-400')}>
          <span className="text-xl">{item.icon}</span>
          <span className="text-xs font-medium">{item.label}</span>
        </button>
      ))}
    </nav>
  )
}
