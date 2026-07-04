// Branded types for type-safe IDs
export type Brand<T, B> = T & { __brand: B }

export type UUID = Brand<string, 'UUID'>
export type OrderID = Brand<string, 'OrderID'>
export type CustomerID = Brand<string, 'CustomerID'>
export type RestaurantID = Brand<string, 'RestaurantID'>
export type DriverID = Brand<string, 'DriverID'>
export type MenuItemID = Brand<string, 'MenuItemID'>
export type AddressID = Brand<string, 'AddressID'>
export type PaymentMethodID = Brand<string, 'PaymentMethodID'>
export type PromoID = Brand<string, 'PromoID'>
export type TicketID = Brand<string, 'TicketID'>
export type EmployeeID = Brand<string, 'EmployeeID'>
