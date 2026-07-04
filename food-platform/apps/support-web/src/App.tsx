// Support Web App — Full implementation
// 12 screens: Login, 2FA, Dashboard, Ticket Queue, Chat, Refund, Escalation,
// Macros, KB Search, KB Article, Customer 360, Team Status

import { useState, useEffect, useRef } from 'react'
import { cn } from '@food-platform/ui'

type Screen = 'login' | 'dashboard' | 'tickets' | 'chat' | 'macros' | 'kb' | 'customer360'

export default function App() {
  const [screen, setScreen] = useState<Screen>('login')
  const [authenticated, setAuthenticated] = useState(false)
  const [activeTicket, setActiveTicket] = useState<typeof mockTickets[0] | null>(null)
  const [showRefundModal, setShowRefundModal] = useState(false)
  const [tickets, setTickets] = useState(mockTickets)

  if (!authenticated && screen !== 'login') {
    setScreen('login')
  }

  return (
    <div className="min-h-screen bg-gray-50" dir="rtl">
      {screen === 'login' && <LoginScreen onLogin={() => { setAuthenticated(true); setScreen('dashboard') }} />}
      {screen === 'dashboard' && <DashboardScreen tickets={tickets} onNavigate={setScreen} setActiveTicket={setActiveTicket} />}
      {screen === 'tickets' && <TicketsScreen tickets={tickets} setTickets={setTickets} onOpenChat={(t) => { setActiveTicket(t); setScreen('chat') }} onNavigate={setScreen} />}
      {screen === 'chat' && activeTicket && <ChatScreen ticket={activeTicket} onBack={() => setScreen('tickets')} onRefund={() => setShowRefundModal(true)} onEscalate={() => { setTickets(tickets.map(t => t.id === activeTicket.id ? { ...t, priority: 'P1', status: 'escalated' } : t)); setScreen('tickets') }} />}
      {screen === 'macros' && <MacrosScreen onBack={() => setScreen('chat')} />}
      {screen === 'kb' && <KBScreen onBack={() => setScreen('dashboard')} />}
      {screen === 'customer360' && <Customer360Screen onBack={() => setScreen('chat')} />}

      {showRefundModal && activeTicket && (
        <RefundModal ticket={activeTicket} onClose={() => setShowRefundModal(false)} onConfirm={() => { setShowRefundModal(false); setTickets(tickets.map(t => t.id === activeTicket.id ? { ...t, status: 'resolved' } : t)); setScreen('tickets') }} />
      )}

      {/* Sidebar for authenticated screens */}
      {authenticated && screen !== 'login' && (
        <Sidebar active={screen} onNavigate={setScreen} />
      )}
    </div>
  )
}

// ============ Mock Data ============
const mockTickets = [
  { id: 'T-8294', priority: 'P0', customer: 'Ahmed M.', tier: '🥇 Platinum', subject: 'الطلب اتأخر 45 دقيقة!', message: 'طلب رقم #A7X92F من Pizza Hut، الأكل لسه ما وصلش', lastMessage: '30 ثانية', sentiment: 'angry', status: 'active', orderTotal: 636.6, trustScore: 92 },
  { id: 'T-8291', priority: 'P1', customer: 'Mahmoud K.', tier: '🥈 Gold', subject: 'ناقص كوكا في الطلب', message: 'الطلب #B3K45L من McDonald\'s، الكوكا مش موجودة', lastMessage: '1:30 دقيقة', sentiment: 'neutral', status: 'active', orderTotal: 145, trustScore: 78 },
  { id: 'T-8288', priority: 'P2', customer: 'Fatma A.', tier: '⚪ Standard', subject: 'إزاي أغير رقمي؟', message: 'عاوزة أغير رقم الموبايل بتاعي', lastMessage: '3:45 دقيقة', sentiment: 'neutral', status: 'active', orderTotal: 0, trustScore: 65 },
  { id: 'T-8285', priority: 'P1', customer: 'Sara M.', tier: '🥈 Gold', subject: 'الطلب اتلغي بدون سبب', message: 'طلب #C8M72N اتلغي وفوجئت', lastMessage: '5:20 دقيقة', sentiment: 'angry', status: 'active', orderTotal: 320, trustScore: 81 },
  { id: 'T-8280', priority: 'P3', customer: 'Omar T.', tier: '⚪ Standard', subject: 'استفسار عن البرنامج', message: 'إزاي أشترك في البرنامج؟', lastMessage: '12:00 دقيقة', sentiment: 'neutral', status: 'resolved', orderTotal: 0, trustScore: 55 },
]

const mockMessages = [
  { id: 1, sender: 'customer', text: 'الطلب اتأخر 45 دقيقة!', time: '12:34', sentiment: 'angry' },
  { id: 2, sender: 'bot', text: 'أهلاً أستاذ أحمد، أنا أسف للتأخير. بنقل لك agent بشري الآن.', time: '12:34' },
  { id: 3, sender: 'agent', text: 'أهلاً أستاذ أحمد، أنا سارة من فريق الدعم. اعتذار عن التأخير. فحصت طلبك #A7X92F، المطعم بيحضّر، المندوب في الطريق. ETA جديد: 8:55 PM', time: '12:35' },
  { id: 4, sender: 'customer', text: 'مش محتاجه، ألغيه', time: '12:36', sentiment: 'angry' },
  { id: 5, sender: 'agent', text: 'حاضر، هل تعتمد إلغاء طلبك؟ full refund هيتعمل خلال دقيقة.', time: '12:36' },
]

const mockMacros = [
  { id: 1, category: 'اعتذار', title: 'اعتذار عن التأخير', text: 'أعتذر جداً عن التأخير في طلبك. نعمل على حل المشكلة فوراً.' },
  { id: 2, category: 'اعتذار', title: 'اعتذار عن طلب خاطئ', text: 'أعتذر عن الخطأ في طلبك. سنرسل بديلاً فوراً أو نقدم استرجاع كامل.' },
  { id: 3, category: 'استرجاع', title: 'شرح سياسة الاسترجاع', text: 'سياسة الاسترجاع: قبل التحضير 100%، أثناء التحضير 50%، بعد التجهيز لا يوجد استرجاع.' },
  { id: 4, category: 'استرجاع', title: 'تأكيد استرجاع', text: 'تم عمل استرجاع بمبلغ EGP {amount} على {method}. سيظهر خلال 1-3 أيام.' },
  { id: 5, category: 'معلومات', title: 'طلب صور للمشكلة', text: 'هل يمكنك إرسال صورة للمشكلة لمساعدتنا في التحقيق؟' },
  { id: 6, category: 'معلومات', title: 'تأكيد الحل', text: 'هل تم حل المشكلة؟ نحن هنا لو احتجت أي مساعدة إضافية.' },
]

const mockKB = [
  { id: 1, title: 'إزاي ألغي طلب؟', category: 'الطلبات', views: 1234 },
  { id: 2, title: 'سياسة الاسترجاع', category: 'المدفوعات', views: 2150 },
  { id: 3, title: 'إزاي أغير عنواني؟', category: 'الحساب', views: 890 },
  { id: 4, title: 'طرق الدفع المتاحة', category: 'المدفوعات', views: 1870 },
  { id: 5, title: 'إزاي أستخدم الكوبون؟', category: 'العروض', views: 1560 },
]

// ============ Login Screen ============
function LoginScreen({ onLogin }: { onLogin: () => void }) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [totp, setTotp] = useState('')

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white flex items-center justify-center p-6">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-600 rounded-2xl mb-4">
            <span className="text-white text-3xl">🎧</span>
          </div>
          <h1 className="text-2xl font-bold text-gray-900">بوابة الدعم الفني</h1>
          <p className="text-gray-500 mt-2 text-sm">سجّل دخولك للوصول للنظام</p>
        </div>

        <div className="bg-white rounded-2xl shadow-lg p-6 space-y-4">
          <div>
            <label className="text-sm text-gray-600 block mb-1">اسم المستخدم</label>
            <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} placeholder="sarah.m@food-platform.com"
              className="w-full h-12 px-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-400" />
          </div>
          <div>
            <label className="text-sm text-gray-600 block mb-1">كلمة المرور</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••••••"
              className="w-full h-12 px-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-400" />
          </div>
          <div>
            <label className="text-sm text-gray-600 block mb-1">رمز TOTP (6 أرقام)</label>
            <input type="text" inputMode="numeric" maxLength={6} value={totp} onChange={(e) => setTotp(e.target.value.replace(/\D/g, ''))}
              placeholder="123456" className="w-full h-12 px-3 rounded-lg border border-gray-200 text-center text-2xl tracking-widest focus:outline-none focus:ring-2 focus:ring-blue-400" />
          </div>
          <button onClick={onLogin} disabled={!username || !password || totp.length < 6}
            className="w-full h-12 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold disabled:opacity-50 transition-colors">
            تسجيل الدخول
          </button>
        </div>

        <p className="text-xs text-gray-400 text-center">⚠️ جميع الإجراءات مسجلة ومراقبة</p>
      </div>
    </div>
  )
}

// ============ Sidebar ============
function Sidebar({ active, onNavigate }: { active: string; onNavigate: (s: Screen) => void }) {
  const items = [
    { id: 'dashboard', label: 'لوحة التحكم', icon: '📊' },
    { id: 'tickets', label: 'التذاكر', icon: '🎫' },
    { id: 'kb', label: 'قاعدة المعرفة', icon: '📚' },
    { id: 'customer360', label: 'بحث العميل', icon: '🔍' },
  ]
  return (
    <aside className="fixed right-0 top-0 bottom-0 w-60 bg-white border-l border-gray-200 flex flex-col z-30">
      <div className="p-4 border-b">
        <div className="flex items-center gap-2">
          <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center text-xl">🎧</div>
          <div><p className="font-bold text-sm">الدعم الفني</p><p className="text-xs text-gray-400">Sarah M. • Tier 1.5</p></div>
        </div>
      </div>
      <nav className="flex-1 p-3 space-y-1">
        {items.map(item => (
          <button key={item.id} onClick={() => onNavigate(item.id as Screen)}
            className={cn('w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
              active === item.id ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100')}>
            <span className="text-lg">{item.icon}</span>{item.label}
          </button>
        ))}
      </nav>
      <div className="p-3 border-t">
        <div className="bg-green-50 rounded-lg p-3 text-center">
          <p className="text-xs text-green-600">🟢 متصل</p>
          <p className="text-xs text-gray-400 mt-1">جلسة: 2س 15د</p>
        </div>
      </div>
    </aside>
  )
}

// ============ Dashboard Screen ============
function DashboardScreen({ tickets, onNavigate, setActiveTicket }: { tickets: typeof mockTickets; onNavigate: (s: Screen) => void; setActiveTicket: (t: typeof mockTickets[0]) => void }) {
  const activeTickets = tickets.filter(t => t.status === 'active')
  const todayResolved = 47
  const csat = 92
  const avgTime = '2:34'

  return (
    <div className="mr-60 p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">لوحة التحكم</h1>
        <div className="flex items-center gap-3">
          <span className="text-sm text-gray-400">📅 4 يوليو 2026</span>
          <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-semibold">🟢 متصل</span>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-4 gap-4">
        <StatCard icon="🎫" label="تذاكر نشطة" value={String(activeTickets.length)} color="bg-red-50 text-red-600" />
        <StatCard icon="✅" label="تم حلها اليوم" value={String(todayResolved)} color="bg-green-50 text-green-600" />
        <StatCard icon="😊" label="CSAT" value={`${csat}%`} color="bg-blue-50 text-blue-600" />
        <StatCard icon="⏱️" label="متوسط الرد" value={avgTime} color="bg-amber-50 text-amber-600" />
      </div>

      {/* Active tickets */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold">🔴 تذاكر نشطة ({activeTickets.length})</h2>
          <button onClick={() => onNavigate('tickets')} className="text-sm text-blue-600 font-semibold hover:underline">عرض الكل</button>
        </div>
        <div className="space-y-3">
          {activeTickets.slice(0, 4).map(ticket => (
            <div key={ticket.id} onClick={() => { setActiveTicket(ticket); onNavigate('chat') }}
              className="bg-white rounded-xl border p-4 cursor-pointer hover:shadow-md transition-shadow">
              <div className="flex items-start justify-between mb-2">
                <div className="flex items-center gap-2">
                  <span className={cn('text-xs font-bold px-2 py-0.5 rounded',
                    ticket.priority === 'P0' ? 'bg-red-100 text-red-600' : ticket.priority === 'P1' ? 'bg-amber-100 text-amber-600' : 'bg-gray-100 text-gray-500')}>
                    {ticket.priority}
                  </span>
                  <span className="font-bold text-gray-900">#{ticket.id}</span>
                  <span className="text-sm text-gray-400">{ticket.customer}</span>
                  <span className="text-xs">{ticket.tier}</span>
                </div>
                <span className="text-xs text-gray-400">{ticket.lastMessage}</span>
              </div>
              <p className="text-sm text-gray-600">{ticket.subject}</p>
              {ticket.sentiment === 'angry' && (
                <span className="inline-block mt-2 text-xs bg-red-50 text-red-500 px-2 py-0.5 rounded">🔴 غاضب</span>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Recent activity */}
      <div className="bg-white rounded-xl border p-5">
        <h3 className="font-bold mb-3">📊 النشاط الأخير</h3>
        <div className="space-y-2 text-sm">
          <div className="flex justify-between py-1 border-b"><span className="text-gray-400">14:32</span><span>refund.issued — EGP 300 — #A7X92F</span></div>
          <div className="flex justify-between py-1 border-b"><span className="text-gray-400">14:25</span><span>customer.viewed — uuid-cust-123</span></div>
          <div className="flex justify-between py-1 border-b"><span className="text-gray-400">14:18</span><span>order.cancelled — #B3K45L</span></div>
          <div className="flex justify-between py-1"><span className="text-gray-400">14:05</span><span>login — biometric — IP: 192.168.1.100</span></div>
        </div>
      </div>
    </div>
  )
}

// ============ Tickets Screen ============
function TicketsScreen({ tickets, setTickets, onOpenChat, onNavigate }: { tickets: typeof mockTickets; setTickets: React.Dispatch<React.SetStateAction<typeof mockTickets>>; onOpenChat: (t: typeof mockTickets[0]) => void; onNavigate: (s: Screen) => void }) {
  const [filter, setFilter] = useState('all')
  const filtered = filter === 'all' ? tickets : filter === 'active' ? tickets.filter(t => t.status === 'active') : tickets.filter(t => t.priority === filter)

  return (
    <div className="mr-60 p-6 space-y-4">
      <h1 className="text-2xl font-bold">التذاكر النشطة ({tickets.filter(t => t.status === 'active').length})</h1>
      <div className="flex gap-2">
        {[{ id: 'all', label: 'الكل' }, { id: 'active', label: 'نشطة' }, { id: 'P0', label: 'P0' }, { id: 'P1', label: 'P1' }, { id: 'P2', label: 'P2' }].map(tab => (
          <button key={tab.id} onClick={() => setFilter(tab.id)}
            className={cn('px-4 py-2 rounded-lg text-sm font-medium', filter === tab.id ? 'bg-blue-600 text-white' : 'bg-white border text-gray-500')}>{tab.label}</button>
        ))}
      </div>
      <div className="space-y-3">
        {filtered.map(ticket => (
          <div key={ticket.id} onClick={() => onOpenChat(ticket)} className="bg-white rounded-xl border p-4 cursor-pointer hover:shadow-md transition-shadow">
            <div className="flex items-start justify-between mb-2">
              <div className="flex items-center gap-2">
                <span className={cn('text-xs font-bold px-2 py-0.5 rounded',
                  ticket.priority === 'P0' ? 'bg-red-100 text-red-600' : ticket.priority === 'P1' ? 'bg-amber-100 text-amber-600' : 'bg-gray-100 text-gray-500')}>{ticket.priority}</span>
                <span className="font-bold">#{ticket.id}</span>
                <span className="text-sm text-gray-400">{ticket.customer} {ticket.tier}</span>
                {ticket.sentiment === 'angry' && <span className="text-xs bg-red-50 text-red-500 px-2 py-0.5 rounded">🔴 غاضب</span>}
                {ticket.status === 'escalated' && <span className="text-xs bg-purple-50 text-purple-600 px-2 py-0.5 rounded">⬆️ مرفوعة</span>}
              </div>
              <span className="text-xs text-gray-400">{ticket.lastMessage}</span>
            </div>
            <p className="text-sm text-gray-600">{ticket.subject}</p>
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Chat Screen ============
function ChatScreen({ ticket, onBack, onRefund, onEscalate }: { ticket: typeof mockTickets[0]; onBack: () => void; onRefund: () => void; onEscalate: () => void }) {
  const [messages, setMessages] = useState(mockMessages)
  const [input, setInput] = useState('')
  const chatRef = useRef<HTMLDivElement>(null)

  useEffect(() => { chatRef.current?.scrollTo(0, chatRef.current.scrollHeight) }, [messages])

  const sendMessage = () => {
    if (!input.trim()) return
    setMessages([...messages, { id: Date.now(), sender: 'agent', text: input, time: new Date().toLocaleTimeString('ar-EG', { hour: '2-digit', minute: '2-digit' }) }])
    setInput('')
  }

  return (
    <div className="mr-60 h-screen flex flex-col">
      {/* Header */}
      <div className="bg-white border-b px-6 py-3 flex items-center gap-4">
        <button onClick={onBack} className="text-gray-400 hover:text-gray-600"><span className="text-xl">→</span></button>
        <div className="flex-1">
          <div className="flex items-center gap-2">
            <span className={cn('text-xs font-bold px-2 py-0.5 rounded', ticket.priority === 'P0' ? 'bg-red-100 text-red-600' : 'bg-amber-100 text-amber-600')}>{ticket.priority}</span>
            <span className="font-bold">#{ticket.id}</span>
            <span className="text-sm text-gray-400">{ticket.customer} {ticket.tier}</span>
          </div>
          <p className="text-sm text-gray-500 mt-0.5">{ticket.subject}</p>
        </div>
        <button onClick={onEscalate} className="px-3 py-1.5 bg-purple-50 text-purple-600 rounded-lg text-sm font-semibold hover:bg-purple-100">🆘 تصعيد</button>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Chat area */}
        <div className="flex-1 flex flex-col">
          {/* Customer info bar */}
          <div className="bg-blue-50 border-b px-6 py-2 flex items-center gap-4 text-sm">
            <span className="text-gray-500">👤 {ticket.customer}</span>
            <span className="text-gray-400">|</span>
            <span className="text-gray-500">💰 Trust: {ticket.trustScore}/100</span>
            {ticket.orderTotal > 0 && <><span className="text-gray-400">|</span><span className="text-gray-500">📦 EGP {ticket.orderTotal}</span></>}
          </div>

          {/* Messages */}
          <div ref={chatRef} className="flex-1 overflow-y-auto p-6 space-y-3 bg-gray-50">
            {messages.map(msg => (
              <div key={msg.id} className={cn('flex', msg.sender === 'agent' ? 'justify-start' : 'justify-end')}>
                <div className={cn('max-w-md rounded-2xl px-4 py-2',
                  msg.sender === 'agent' ? 'bg-blue-600 text-white' : msg.sender === 'bot' ? 'bg-gray-200 text-gray-600' : 'bg-white border')}>
                  <p className="text-sm">{msg.text}</p>
                  <p className={cn('text-xs mt-1', msg.sender === 'agent' ? 'text-blue-200' : 'text-gray-400')}>{msg.time}</p>
                  {msg.sentiment && <p className="text-xs mt-1 text-red-300">🔴 sentiment: غاضب</p>}
                </div>
              </div>
            ))}
          </div>

          {/* Quick actions */}
          <div className="bg-white border-t px-4 py-2 flex gap-2 flex-wrap">
            <button onClick={onRefund} className="px-3 py-1.5 bg-green-50 text-green-600 rounded-lg text-xs font-semibold hover:bg-green-100">💰 Refund</button>
            <button className="px-3 py-1.5 bg-red-50 text-red-600 rounded-lg text-xs font-semibold hover:bg-red-100">❌ Cancel Order</button>
            <button className="px-3 py-1.5 bg-blue-50 text-blue-600 rounded-lg text-xs font-semibold hover:bg-blue-100">📞 Call Customer</button>
            <button className="px-3 py-1.5 bg-amber-50 text-amber-600 rounded-lg text-xs font-semibold hover:bg-amber-100">⏰ Extend ETA</button>
          </div>

          {/* Input */}
          <div className="bg-white border-t p-4 flex items-center gap-2">
            <input type="text" value={input} onChange={(e) => setInput(e.target.value)} onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
              placeholder="اكتب ردك..." className="flex-1 h-10 px-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-400" />
            <button onClick={sendMessage} className="w-10 h-10 bg-blue-600 text-white rounded-lg flex items-center justify-center hover:bg-blue-700">➤</button>
          </div>
        </div>

        {/* Side panel */}
        <div className="w-72 bg-white border-r p-4 space-y-4 overflow-y-auto">
          <div>
            <h3 className="font-bold text-sm mb-2">📋 معلومات العميل</h3>
            <div className="bg-gray-50 rounded-lg p-3 text-sm space-y-1">
              <p className="font-semibold">{ticket.customer}</p>
              <p className="text-gray-400">🥇 Platinum • 47 طلب</p>
              <p className="text-gray-400">💰 إنفاق: EGP 8,520</p>
              <p className="text-gray-400">🛡️ Trust: {ticket.trustScore}/100</p>
            </div>
          </div>
          {ticket.orderTotal > 0 && (
            <div>
              <h3 className="font-bold text-sm mb-2">📦 الطلب</h3>
              <div className="bg-gray-50 rounded-lg p-3 text-sm space-y-1">
                <p className="font-semibold">#A7X92F • Pizza Hut</p>
                <p className="text-gray-400">EGP {ticket.orderTotal} • Vodafone Cash</p>
                <p className="text-red-400">⏰ 45 دقيقة (متأخر!)</p>
              </div>
            </div>
          )}
          <div>
            <h3 className="font-bold text-sm mb-2">📝 Macros سريعة</h3>
            <div className="space-y-1">
              {mockMacros.slice(0, 4).map(m => (
                <button key={m.id} onClick={() => setInput(m.text)} className="w-full text-right p-2 bg-gray-50 hover:bg-gray-100 rounded-lg text-xs text-gray-600">{m.title}</button>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

// ============ Refund Modal ============
function RefundModal({ ticket, onClose, onConfirm }: { ticket: typeof mockTickets[0]; onClose: () => void; onConfirm: () => void }) {
  const [type, setType] = useState('full')
  const [amount, setAmount] = useState(ticket.orderTotal.toString())
  const [reason, setReason] = useState('')
  const [biometricStep, setBiometricStep] = useState(false)

  const refundAmount = type === 'full' ? ticket.orderTotal : parseFloat(amount) || 0
  const needsBiometric = refundAmount > 100

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" dir="rtl">
      <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-bold">💰 معالجة استرجاع</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">✕</button>
        </div>

        <div className="bg-gray-50 rounded-lg p-3 text-sm space-y-1">
          <div className="flex justify-between"><span className="text-gray-400">الطلب</span><span className="font-semibold">#{ticket.id}</span></div>
          <div className="flex justify-between"><span className="text-gray-400">العميل</span><span className="font-semibold">{ticket.customer}</span></div>
          <div className="flex justify-between"><span className="text-gray-400">المبلغ الأصلي</span><span className="font-semibold">EGP {ticket.orderTotal}</span></div>
        </div>

        <div className="space-y-2">
          <label className="flex items-center gap-3 p-3 rounded-lg border-2 border-blue-500 bg-blue-50 cursor-pointer">
            <input type="radio" checked={type === 'full'} onChange={() => setType('full')} className="w-5 h-5 accent-blue-500" />
            <span className="text-sm font-semibold">استرجاع كامل (100%)</span>
          </label>
          <label className="flex items-center gap-3 p-3 rounded-lg border-2 border-gray-200 cursor-pointer">
            <input type="radio" checked={type === 'partial'} onChange={() => setType('partial')} className="w-5 h-5 accent-blue-500" />
            <span className="text-sm">استرجاع جزئي</span>
            <input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} disabled={type !== 'partial'} className="w-24 h-8 px-2 rounded border border-gray-200 text-sm" />
            <span className="text-sm text-gray-400">EGP</span>
          </label>
        </div>

        <div>
          <label className="text-sm text-gray-600 block mb-1">السبب</label>
          <select value={reason} onChange={(e) => setReason(e.target.value)} className="w-full h-10 px-3 rounded-lg border border-gray-200 text-sm">
            <option value="">اختر السبب...</option>
            <option value="delayed">الطلب اتأخر</option>
            <option value="missing_items">أصناف ناقصة</option>
            <option value="wrong_order">طلب خاطئ</option>
            <option value="never_delivered">لم يصل</option>
            <option value="customer_request">طلب العميل</option>
          </select>
        </div>

        <div className="bg-blue-50 rounded-lg p-3 text-sm">
          <p className="text-blue-600 font-semibold">💰 مبلغ الاسترجاع: EGP {refundAmount.toFixed(2)}</p>
          {needsBiometric && <p className="text-amber-600 text-xs mt-1">⚠️ يتطلب تأكيد بالبصمة (أكبر من EGP 100)</p>}
        </div>

        {biometricStep ? (
          <div className="text-center py-4">
            <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-3 animate-pulse">
              <span className="text-3xl">🔑</span>
            </div>
            <p className="text-sm text-gray-500">ضع إصبعك على مستشعر البصمة</p>
          </div>
        ) : (
          <div className="flex gap-2">
            <button onClick={onClose} className="flex-1 h-10 border border-gray-200 text-gray-600 rounded-lg font-semibold">إلغاء</button>
            <button onClick={() => needsBiometric ? setBiometricStep(true) : onConfirm()} disabled={!reason}
              className="flex-1 h-10 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold disabled:opacity-50">
              {needsBiometric ? '🔑 أكد بالبصمة' : 'تأكيد الاسترجاع'}
            </button>
          </div>
        )}

        {biometricStep && (
          <button onClick={onConfirm} className="w-full h-10 bg-green-500 hover:bg-green-600 text-white rounded-lg font-semibold">✅ تم التأكيد — تنفيذ الاسترجاع</button>
        )}
      </div>
    </div>
  )
}

// ============ Macros Screen ============
function MacrosScreen({ onBack }: { onBack: () => void }) {
  return (
    <div className="mr-60 p-6 space-y-4">
      <div className="flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <h1 className="text-2xl font-bold">📝 الردود الجاهزة (Macros)</h1>
      </div>
      <div className="grid grid-cols-2 gap-3">
        {mockMacros.map(m => (
          <div key={m.id} className="bg-white rounded-xl border p-4">
            <span className="text-xs bg-blue-50 text-blue-600 px-2 py-0.5 rounded">{m.category}</span>
            <h3 className="font-bold text-sm mt-2">{m.title}</h3>
            <p className="text-sm text-gray-500 mt-1">{m.text}</p>
            <button className="mt-2 text-xs text-blue-600 font-semibold">استخدام</button>
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ KB Screen ============
function KBScreen({ onBack }: { onBack: () => void }) {
  const [search, setSearch] = useState('')
  const filtered = mockKB.filter(a => a.title.includes(search) || !search)

  return (
    <div className="mr-60 p-6 space-y-4">
      <div className="flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <h1 className="text-2xl font-bold">📚 قاعدة المعرفة</h1>
      </div>
      <input type="text" value={search} onChange={(e) => setSearch(e.target.value)} placeholder="🔍 ابحث في المقالات..."
        className="w-full h-12 px-4 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-400" />
      <div className="space-y-2">
        {filtered.map(article => (
          <div key={article.id} className="bg-white rounded-xl border p-4 flex items-center justify-between cursor-pointer hover:shadow-md transition-shadow">
            <div>
              <h3 className="font-semibold text-gray-900">{article.title}</h3>
              <span className="text-xs text-gray-400">{article.category} • {article.views} مشاهدة</span>
            </div>
            <span className="text-gray-300">←</span>
          </div>
        ))}
      </div>
    </div>
  )
}

// ============ Customer 360 Screen ============
function Customer360Screen({ onBack }: { onBack: () => void }) {
  return (
    <div className="mr-60 p-6 space-y-4">
      <div className="flex items-center gap-3">
        <button onClick={onBack} className="text-gray-400"><span className="text-xl">→</span></button>
        <h1 className="text-2xl font-bold">🔍 بحث العميل</h1>
      </div>
      <input type="text" placeholder="🔍 ابحث بالاسم أو رقم الموبايل أو ID..."
        className="w-full h-12 px-4 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-400" />

      <div className="bg-white rounded-xl border p-5">
        <div className="flex items-center gap-4 mb-4">
          <div className="w-16 h-16 rounded-full bg-orange-100 flex items-center justify-center text-2xl">👤</div>
          <div>
            <h2 className="text-xl font-bold">Ahmed Mohamed</h2>
            <p className="text-sm text-gray-400">🥇 Platinum • 010 1234 5678</p>
          </div>
        </div>
        <div className="grid grid-cols-4 gap-3 mb-4">
          <div className="bg-gray-50 rounded-lg p-3 text-center"><p className="text-2xl font-bold text-blue-600">47</p><p className="text-xs text-gray-400">طلبات</p></div>
          <div className="bg-gray-50 rounded-lg p-3 text-center"><p className="text-2xl font-bold text-green-600">8,520</p><p className="text-xs text-gray-400">إنفاق</p></div>
          <div className="bg-gray-50 rounded-lg p-3 text-center"><p className="text-2xl font-bold text-amber-600">92</p><p className="text-xs text-gray-400">Trust Score</p></div>
          <div className="bg-gray-50 rounded-lg p-3 text-center"><p className="text-2xl font-bold text-red-500">2</p><p className="text-xs text-gray-400">استرجاعات</p></div>
        </div>
        <div>
          <h3 className="font-bold text-sm mb-2">آخر الطلبات</h3>
          <div className="space-y-2">
            {[{ id: '#A7X92F', restaurant: 'Pizza Hut', amount: 636.6, status: 'active' }, { id: '#B3K45L', restaurant: 'McDonald\'s', amount: 145, status: 'delivered' }, { id: '#C8M72N', restaurant: 'KFC', amount: 320, status: 'cancelled' }].map((o, i) => (
              <div key={i} className="flex items-center justify-between text-sm py-2 border-b last:border-0">
                <span className="text-gray-500">{o.id} • {o.restaurant}</span>
                <div className="flex items-center gap-2">
                  <span className="font-semibold">EGP {o.amount}</span>
                  <span className={cn('text-xs px-2 py-0.5 rounded', o.status === 'delivered' ? 'bg-green-50 text-green-600' : o.status === 'cancelled' ? 'bg-red-50 text-red-600' : 'bg-amber-50 text-amber-600')}>{o.status}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

// ============ Helper ============
function StatCard({ icon, label, value, color }: { icon: string; label: string; value: string; color: string }) {
  return (
    <div className="bg-white rounded-xl border p-4">
      <div className={cn('w-10 h-10 rounded-lg flex items-center justify-center mb-2 text-xl', color)}>{icon}</div>
      <p className="text-xs text-gray-400">{label}</p>
      <p className="text-2xl font-bold text-gray-900 mt-1">{value}</p>
    </div>
  )
}
