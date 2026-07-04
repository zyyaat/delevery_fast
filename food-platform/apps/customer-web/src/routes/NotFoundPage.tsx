import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <div className="min-h-screen bg-bg-primary flex items-center justify-center p-4">
      <div className="text-center">
        <h1 className="text-display-md font-extrabold text-primary mb-4">404</h1>
        <p className="text-body-lg text-text-secondary mb-6">الصفحة مش موجودة</p>
        <Link
          to="/"
          className="inline-flex items-center justify-center px-6 py-3 bg-primary text-white rounded-md font-semibold hover:bg-primary-dark transition-colors"
        >
          العودة للرئيسية
        </Link>
      </div>
    </div>
  )
}
