import { Routes, Route } from 'react-router-dom'
import { RequireAuth, RedirectIfAuthed } from '@food-platform/auth'
import { HomePage } from './routes/home/HomePage'
import { LoginPage } from './routes/auth/LoginPage'
import { OtpPage } from './routes/auth/OtpPage'
import { RestaurantDetailPage } from './routes/restaurant-detail/RestaurantDetailPage'
import { CartPage } from './routes/cart/CartPage'
import { CheckoutPage } from './routes/checkout/CheckoutPage'
import { OrderTrackingPage } from './routes/order-tracking/OrderTrackingPage'
import { OrdersHistoryPage } from './routes/orders-history/OrdersHistoryPage'
import { ProfilePage } from './routes/profile/ProfilePage'
import { NotFoundPage } from './routes/NotFoundPage'

export default function App() {
  return (
    <Routes>
      {/* Auth routes (redirect if already authenticated) */}
      <Route
        path="/login"
        element={
          <RedirectIfAuthed>
            <LoginPage />
          </RedirectIfAuthed>
        }
      />
      <Route
        path="/otp"
        element={
          <RedirectIfAuthed>
            <OtpPage />
          </RedirectIfAuthed>
        }
      />

      {/* ====== GUEST ROUTES (no auth required) ====== */}
      {/* Guests can browse restaurants, view menus, and add to cart */}
      <Route path="/" element={<HomePage />} />
      <Route path="/restaurants/:id" element={<RestaurantDetailPage />} />
      <Route path="/cart" element={<CartPage />} />

      {/* ====== PROTECTED ROUTES (auth required) ====== */}
      {/* Checkout, orders, and profile require login */}
      <Route
        path="/checkout"
        element={
          <RequireAuth>
            <CheckoutPage />
          </RequireAuth>
        }
      />
      <Route
        path="/orders/:id"
        element={
          <RequireAuth>
            <OrderTrackingPage />
          </RequireAuth>
        }
      />
      <Route
        path="/orders"
        element={
          <RequireAuth>
            <OrdersHistoryPage />
          </RequireAuth>
        }
      />
      <Route
        path="/profile"
        element={
          <RequireAuth>
            <ProfilePage />
          </RequireAuth>
        }
      />

      {/* 404 */}
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}
