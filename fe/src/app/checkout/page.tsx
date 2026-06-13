'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import toast from 'react-hot-toast';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useCartStore, useAuthStore } from '@/store';
import { orderService } from '@/services/order-service';
import { formatPrice } from '@/lib/utils';
import { MapPin, Phone, User, CreditCard } from 'lucide-react';
import type { CreateOrderRequest } from '@/types';

export default function CheckoutPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const items = useCartStore((state) => state.items);
  const totalPrice = useCartStore((state) => state.totalPrice);
  const clearCart = useCartStore((state) => state.clearCart);

  // Form states
  const [shippingName, setShippingName] = useState('');
  const [shippingPhone, setShippingPhone] = useState('');
  const [shippingAddress, setShippingAddress] = useState('');
  const [shippingCity, setShippingCity] = useState('');
  const [shippingPostal, setShippingPostal] = useState('');
  const [errors, setErrors] = useState<{ [key: string]: string }>({});

  if (!isAuthenticated) {
    router.push('/login');
    return null;
  }
  
  if (items.length === 0) {
    router.push('/cart');
    return null;
  }

  const validateForm = () => {
    const newErrors: { [key: string]: string } = {};
    if (!shippingName.trim()) newErrors.name = 'Recipient name is required';
    if (!shippingPhone.trim()) newErrors.phone = 'Phone number is required';
    if (!shippingAddress.trim()) newErrors.address = 'Shipping address is required';
    if (!shippingCity.trim()) newErrors.city = 'City is required';
    if (!shippingPostal.trim()) newErrors.postal = 'Postal code is required';
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };
  
  const handlePlaceOrder = async () => {
    if (!validateForm()) {
      toast.error('Please fill in all shipping details');
      return;
    }

    setIsLoading(true);
    try {
      const orderPayload: CreateOrderRequest = {
        items: items.map((item) => ({
          product_id: item.product.id,
          quantity: item.quantity,
        })),
      };
      
      const response = await orderService.createOrder(orderPayload);
      
      clearCart();
      toast.success('Order placed! Redirecting to payment...');
      
      // Redirect to order detail page where user can pay
      router.push(`/orders/${response.order.id}`);
    } catch (error: unknown) {
      const errMsg = error instanceof Error ? error.message : 'Failed to place order';
      toast.error(errMsg);
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-3xl font-bold text-gradient mb-8 animate-fade-in">Checkout</h1>
      
      <div className="lg:grid lg:grid-cols-12 lg:gap-8">
        {/* Shipping Form & Payment Info */}
        <div className="lg:col-span-7 space-y-6 animate-slide-up">
          {/* Shipping Form */}
          <div className="glass-card p-6 border border-white/[0.06] shadow-xl">
            <h2 className="text-lg font-semibold text-slate-100 mb-6 flex items-center gap-2">
              <MapPin className="w-5 h-5 text-accent" />
              Shipping Information
            </h2>

            <div className="space-y-4">
              <Input
                label="Recipient Name"
                placeholder="John Doe"
                value={shippingName}
                onChange={(e) => setShippingName(e.target.value)}
                error={errors.name}
              />
              <Input
                label="Phone Number"
                placeholder="+62 812 3456 7890"
                value={shippingPhone}
                onChange={(e) => setShippingPhone(e.target.value)}
                error={errors.phone}
              />
              <Input
                label="Full Address"
                placeholder="Street Name, Building, Unit Number"
                value={shippingAddress}
                onChange={(e) => setShippingAddress(e.target.value)}
                error={errors.address}
              />
              <div className="grid grid-cols-2 gap-4">
                <Input
                  label="City"
                  placeholder="Jakarta"
                  value={shippingCity}
                  onChange={(e) => setShippingCity(e.target.value)}
                  error={errors.city}
                />
                <Input
                  label="Postal Code"
                  placeholder="12345"
                  value={shippingPostal}
                  onChange={(e) => setShippingPostal(e.target.value)}
                  error={errors.postal}
                />
              </div>
            </div>
          </div>

          {/* Payment Info */}
          <div className="glass-card p-6 border border-white/[0.06] shadow-xl bg-white/[0.01]">
            <h2 className="text-lg font-semibold text-slate-100 mb-3 flex items-center gap-2">
              <CreditCard className="w-5 h-5 text-accent" />
              Payment Information
            </h2>
            <p className="text-slate-400 text-sm leading-relaxed">
              After placing your order, you will be redirected to complete payment 
              using our secure payment gateway (Midtrans). You can pay using various methods 
              including Credit Card, Bank Transfer, GoPay, OVO, and more.
            </p>
          </div>
        </div>

        {/* Order Summary & Action */}
        <div className="lg:col-span-5 mt-8 lg:mt-0 animate-slide-up" style={{ animationDelay: '0.1s' }}>
          <div className="glass-card p-6 sticky top-24 shadow-xl border border-white/[0.06]">
            <h2 className="text-lg font-semibold text-slate-100 mb-4">Order Summary</h2>
            
            <div className="divide-y divide-white/[0.06] max-h-60 overflow-y-auto mb-4 pr-2">
              {items.map((item) => (
                <div key={item.product.id} className="py-3 flex justify-between text-sm">
                  <div>
                    <p className="font-medium text-slate-200 truncate max-w-[200px]">{item.product.name}</p>
                    <p className="text-xs text-slate-400 mt-0.5">
                      {formatPrice(item.product.price)} × {item.quantity}
                    </p>
                  </div>
                  <p className="font-medium text-slate-200">
                    {formatPrice(item.product.price * item.quantity)}
                  </p>
                </div>
              ))}
            </div>
            
            <hr className="my-4 border-white/[0.06]" />
            
            <div className="flex justify-between text-lg font-semibold text-slate-100 mb-6">
              <span>Total</span>
              <span className="text-accent-hover">{formatPrice(totalPrice())}</span>
            </div>

            {/* Actions */}
            <div className="flex flex-col gap-3">
              <Button
                className="w-full shadow-lg shadow-accent/20"
                onClick={handlePlaceOrder}
                isLoading={isLoading}
              >
                Place Order
              </Button>
              <Button
                variant="outline"
                className="w-full"
                onClick={() => router.back()}
              >
                Back to Cart
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
