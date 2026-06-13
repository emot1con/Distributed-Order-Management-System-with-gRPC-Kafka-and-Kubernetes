'use client';

import React, { useState, useMemo, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { CreditCard, CheckCircle2, ArrowLeft, Clock, XCircle, Loader2, Copy } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Skeleton, ErrorCard } from '@/components/ui';
import { useAuthStore } from '@/store';
import { orderService } from '@/services/order-service';
import { paymentService } from '@/services/payment-service';
import { formatPrice, formatDate } from '@/lib/utils';
import type { OrderItem } from '@/types';

function OrderDetailSkeleton() {
  return (
    <div className="max-w-lg mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Back link */}
      <div className="h-5 w-32 bg-white/[0.06] rounded animate-pulse mb-6" />
      
      <div className="glass-card overflow-hidden border border-white/[0.06] shadow-xl p-6">
        <div className="space-y-3 mb-6">
          <Skeleton className="h-6 w-1/2 rounded" />
          <Skeleton className="h-4 w-1/4 rounded" />
        </div>
        <hr className="border-white/[0.06] mb-6" />
        <div className="space-y-4 mb-6">
          <Skeleton className="h-5 w-1/3 rounded" />
          <Skeleton className="h-4 w-full rounded" />
          <Skeleton className="h-4 w-5/6 rounded" />
        </div>
        <hr className="border-white/[0.06] mb-6" />
        <div className="flex justify-between items-center mb-6">
          <Skeleton className="h-5 w-16 rounded" />
          <Skeleton className="h-6 w-24 rounded" />
        </div>
        <hr className="border-white/[0.06] mb-6" />
        <Skeleton className="h-12 w-full rounded-lg" />
      </div>
    </div>
  );
}

function Sparkles() {
  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none z-0">
      <style dangerouslySetInnerHTML={{__html: `
        @keyframes sparkle-float {
          0% { transform: translateY(100px) scale(0) rotate(0deg); opacity: 0; }
          50% { opacity: 0.8; }
          100% { transform: translateY(-150px) scale(1) rotate(360deg); opacity: 0; }
        }
        .sparkle-particle {
          position: absolute;
          border-radius: 50%;
          animation: sparkle-float 4s ease-in-out infinite;
        }
      `}} />
      {Array.from({ length: 24 }).map((_, i) => {
        const left = `${Math.random() * 100}%`;
        const top = `${60 + Math.random() * 40}%`;
        const size = `${Math.random() * 8 + 4}px`;
        const delay = `${Math.random() * 4}s`;
        const duration = `${Math.random() * 3 + 3}s`;
        const color = ['bg-accent/40', 'bg-emerald-400/40', 'bg-cyan-400/40', 'bg-indigo-400/40'][i % 4];
        return (
          <div
            key={i}
            className={`sparkle-particle ${color}`}
            style={{
              left,
              top,
              width: size,
              height: size,
              animationDelay: delay,
              animationDuration: duration,
            }}
          />
        );
      })}
    </div>
  );
}

export default function OrderPaymentPage() {
  const params = useParams();
  const router = useRouter();
  const orderId = Number(params.id);
  const [isPaymentLoading, setIsPaymentLoading] = useState(false);
  const [isSnapOpen, setIsSnapOpen] = useState(false);
  
  // Get user info from auth store
  const user = useAuthStore((state) => state.user);
  
  const { data: orderData, isLoading, error, refetch } = useQuery({
    queryKey: ['order', orderId],
    queryFn: () => orderService.getOrderById(orderId),
    enabled: !isNaN(orderId),
    refetchOnMount: 'always',
    staleTime: 0,
  });

  const { data: paymentData, isLoading: isPaymentDataLoading, refetch: refetchPayment } = useQuery({
    queryKey: ['payment', orderId],
    queryFn: () => paymentService.getPaymentByOrderId(orderId),
    enabled: !isNaN(orderId),
    refetchOnMount: 'always',
    staleTime: 0,
    retry: 3, 
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 5000), 
  });
  
  const order = useMemo(() => {
    return orderData?.order || null;
  }, [orderData]);

  // Auto-poll status when it's pending to automatically detect backend webhook updates
  useQuery({
    queryKey: ['order-poll', orderId],
    queryFn: async () => {
      const result = await orderService.getOrderById(orderId);
      const status = result.order.status.toLowerCase();
      if (status === 'paid' || status === 'success') {
        refetch();
        refetchPayment();
        toast.success('Payment updated successfully!');
      }
      return result;
    },
    enabled: !!order && order.status.toLowerCase() === 'pending',
    refetchInterval: 5000, // Poll every 5s
  });

  // Countdown timer for payment expiration
  const [timeLeft, setTimeLeft] = useState<string>('');
  useEffect(() => {
    if (!paymentData?.expired_at) return;

    const interval = setInterval(() => {
      const difference = +new Date(paymentData.expired_at) - +new Date();
      if (difference <= 0) {
        setTimeLeft('Expired');
        clearInterval(interval);
      } else {
        const hours = Math.floor((difference / (1000 * 60 * 60)) % 24);
        const minutes = Math.floor((difference / 1000 / 60) % 60);
        const seconds = Math.floor((difference / 1000) % 60);
        
        const formatted = [
          hours.toString().padStart(2, '0'),
          minutes.toString().padStart(2, '0'),
          seconds.toString().padStart(2, '0')
        ].join(':');
        
        setTimeLeft(formatted);
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [paymentData?.expired_at]);

  const handleCopyVa = () => {
    if (paymentData?.va_number) {
      navigator.clipboard.writeText(paymentData.va_number);
      toast.success('VA Number copied to clipboard!');
    }
  };
  
  const handlePayWithMidtrans = async () => {
    // Prevent double-click and opening popup when already open
    if (isPaymentLoading || isSnapOpen) {
      return;
    }

    // If payment not loaded yet, try to refetch
    if (!paymentData?.id) {
      toast('Fetching payment info...', { icon: '⏳' });
      const result = await refetchPayment();
      if (!result.data?.id) {
        toast.error('Payment not found. Please wait a moment and try again.');
        return;
      }
    }

    // Get user info - use stored user or fallback
    const customerName = user?.full_name || user?.email?.split('@')[0] || 'Customer';
    const customerEmail = user?.email || '';
    
    if (!customerEmail) {
      toast.error('User information not available. Please login again.');
      return;
    }
    
    setIsPaymentLoading(true);
    
    // Get current payment data (might have been refetched)
    const currentPayment = paymentData || (await refetchPayment()).data;
    
    try {
      let token = currentPayment?.gateway_token;
      
      if (!token) {
        const response = await paymentService.initiatePayment({ 
          order_id: orderId,
          customer_name: customerName,
          customer_email: customerEmail,
        });
        token = response.token;
        
        if (!token) {
          throw new Error('Failed to get payment token');
        }
      }
      
      // Mark snap as open before calling
      setIsSnapOpen(true);
      setIsPaymentLoading(false);
      
      paymentService.openSnapPayment(token, {
        onSuccess: () => {
          setIsSnapOpen(false);
          toast.success('Payment successful! 🎉');
          refetch();
          refetchPayment();
        },
        onPending: () => {
          setIsSnapOpen(false);
          toast('Payment pending. Complete your payment.', { icon: '⏳' });
          refetchPayment();
        },
        onError: () => {
          setIsSnapOpen(false);
          toast.error('Payment failed. Please try again.');
        },
        onClose: () => {
          setIsSnapOpen(false);
          toast('Payment cancelled', { icon: '❌' });
          refetchPayment();
        },
      });
    } catch (err: unknown) {
      setIsSnapOpen(false);
      const errMsg = err instanceof Error ? err.message : 'Failed to initiate payment';
      toast.error(errMsg);
    } finally {
      setIsPaymentLoading(false);
    }
  };
  
  if (isLoading || isPaymentDataLoading) {
    return <OrderDetailSkeleton />;
  }
  
  if (error || !order) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center bg-surface px-4">
        <ErrorCard 
          message="Order not found or an error occurred while loading. Please try again."
          onRetry={() => router.push('/orders')}
        />
      </div>
    );
  }

  const renderStatusContent = () => {
    const status = order.status.toLowerCase();
    
    if (status === 'paid' || status === 'success') {
      return (
        <div className="glass-card p-8 text-center border border-white/[0.06] shadow-xl relative overflow-hidden">
          {/* Confetti particles */}
          <Sparkles />
          
          <div className="relative z-10">
            <CheckCircle2 className="w-16 h-16 text-emerald-400 mx-auto mb-4 animate-pulse-glow" />
            <h1 className="text-2xl font-bold text-gradient mb-2 animate-scale-up">Payment Successful!</h1>
            <p className="text-slate-400 mb-2">Order #{order.id} has been paid</p>
            <p className="text-2xl font-bold text-emerald-400 mb-2">{formatPrice(order.total_price)}</p>
            <p className="text-xs text-slate-500 mb-6">{formatDate(order.created_at)}</p>
            <Link href="/orders" className="block w-full">
              <Button className="w-full">
                <ArrowLeft className="w-4 h-4 mr-2" />
                Back to My Orders
              </Button>
            </Link>
          </div>
        </div>
      );
    }
    
    if (status === 'expired' || status === 'failed' || status === 'cancelled') {
      return (
        <div className="glass-card p-8 text-center border border-white/[0.06] shadow-xl">
          <XCircle className="w-16 h-16 text-red-400 mx-auto mb-4 animate-pulse-glow" />
          <h1 className="text-2xl font-bold text-gradient mb-2 animate-scale-up">
            Payment {status.charAt(0).toUpperCase() + status.slice(1)}
          </h1>
          <p className="text-slate-400 mb-6">Order #{order.id} - {formatPrice(order.total_price)}</p>
          <Link href="/products" className="block w-full">
            <Button className="w-full">Continue Shopping</Button>
          </Link>
        </div>
      );
    }
    
    return (
      <div className="glass-card overflow-hidden border border-white/[0.06] shadow-xl animate-slide-up">
        <div className="bg-gradient-to-br from-accent/20 to-surface-card p-6 border-b border-white/[0.06]">
          <div className="flex items-center gap-3 mb-2">
            <CreditCard className="w-6 h-6 text-accent" />
            <h1 className="text-xl font-bold text-slate-100">Complete Payment</h1>
          </div>
          <p className="text-slate-400 text-sm">Order #{order.id}</p>
        </div>

        <div className="p-6 border-b border-white/[0.06]">
          <h2 className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-4">Order Summary</h2>
          
          {order.order_items && order.order_items.length > 0 ? (
            <div className="space-y-3 mb-4">
              {order.order_items.map((item: OrderItem) => (
                <div key={item.id} className="flex justify-between text-sm">
                  <span className="text-slate-400">Product #{item.product_id} × {item.quantity}</span>
                  <span className="font-semibold text-slate-200">{formatPrice(item.price)}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-slate-500 mb-4">Order items not available</p>
          )}

          <div className="flex justify-between pt-4 border-t border-white/[0.06]">
            <span className="font-semibold text-slate-100">Total</span>
            <span className="font-bold text-xl text-accent-hover">{formatPrice(order.total_price)}</span>
          </div>
        </div>

        {/* VA Copyable section */}
        {paymentData?.va_number && (
          <div className="p-6 border-b border-white/[0.06] bg-white/[0.01]">
            <div className="flex justify-between items-center">
              <div>
                <p className="text-xs text-slate-500 uppercase font-semibold">Virtual Account Number</p>
                <p className="text-lg font-mono font-bold text-slate-200 tracking-wider mt-1">{paymentData.va_number}</p>
              </div>
              <Button variant="outline" size="sm" onClick={handleCopyVa} className="h-8 w-8 p-0">
                <Copy className="w-4 h-4" />
              </Button>
            </div>
          </div>
        )}

        <div className="p-6 border-b border-white/[0.06] bg-white/[0.02]">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-amber-400">
              <Clock className="w-5 h-5 animate-pulse" />
              <span className="font-medium text-sm">Awaiting Payment</span>
            </div>
            {timeLeft && (
              <span className="font-mono text-sm bg-amber-400/10 text-amber-400 px-2 py-0.5 rounded border border-amber-400/20">
                {timeLeft}
              </span>
            )}
          </div>
          {paymentData ? (
            <div className="mt-3 text-xs text-slate-400 space-y-1">
              <p>Payment ID: <span className="text-slate-300 font-mono">{paymentData.id}</span></p>
              {paymentData.expired_at && <p>Expires: <span className="text-slate-300">{formatDate(paymentData.expired_at)}</span></p>}
            </div>
          ) : (
            <div className="mt-3 text-xs text-amber-400/80">
              <p>Payment is being prepared. Click button below to continue.</p>
            </div>
          )}
        </div>

        <div className="p-6">
          <Button className="w-full shadow-lg shadow-accent/25" size="lg" onClick={handlePayWithMidtrans} disabled={isPaymentLoading || isPaymentDataLoading || isSnapOpen}>
            {isPaymentLoading ? (
              <>
                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                Processing...
              </>
            ) : isPaymentDataLoading ? (
              <>
                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                Loading...
              </>
            ) : isSnapOpen ? (
              <>
                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                Payment in progress...
              </>
            ) : (
              <>
                <CreditCard className="w-5 h-5 mr-2" />
                Pay Now - {formatPrice(order.total_price)}
              </>
            )}
          </Button>
          <p className="text-center text-xs text-slate-500 mt-4">Secure payment powered by Midtrans</p>
        </div>
      </div>
    );
  };
  
  return (
    <div className="max-w-lg mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <Link href="/orders" className="inline-flex items-center text-sm text-slate-400 hover:text-slate-200 mb-6 transition-colors">
        <ArrowLeft className="w-4 h-4 mr-2" />
        Back to My Orders
      </Link>
      {renderStatusContent()}
    </div>
  );
}
