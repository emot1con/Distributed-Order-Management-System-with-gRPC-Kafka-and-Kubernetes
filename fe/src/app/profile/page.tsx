'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { User, Mail, LogOut, Package, Calendar, Award } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui';
import { useAuthStore } from '@/store';
import { orderService } from '@/services/order-service';
import { formatDate, formatPrice } from '@/lib/utils';
import Link from 'next/link';

export default function ProfilePage() {
  const router = useRouter();
  const { user, isAuthenticated, logout, role } = useAuthStore();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Fetch orders to calculate stats and show recent orders
  const { data, isLoading } = useQuery({
    queryKey: ['profile-orders'],
    queryFn: () => orderService.getOrders(0),
    enabled: isAuthenticated,
  });

  if (!isAuthenticated || !user) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center px-4">
        <User className="w-16 h-16 text-slate-600 mb-4 animate-pulse-glow" />
        <h2 className="text-xl font-semibold text-slate-100 mb-2">Please login first</h2>
        <p className="text-slate-400 mb-6 text-center max-w-sm">You need to be logged in to view your profile</p>
        <Link href="/login">
          <Button size="lg">Login</Button>
        </Link>
      </div>
    );
  }

  const handleLogout = async () => {
    setIsLoggingOut(true);
    try {
      logout();
      router.push('/');
    } finally {
      setIsLoggingOut(false);
    }
  };

  const orders = data?.orders || [];
  const recentOrders = orders.slice(0, 3);
  const paidOrders = orders.filter(
    (o) => o.status.toLowerCase() === 'paid' || o.status.toLowerCase() === 'success'
  );
  const totalSpend = paidOrders.reduce((sum, o) => sum + o.total_price, 0);

  const getBadgeVariant = (status: string) => {
    const s = status.toLowerCase();
    if (s === 'paid' || s === 'success') return 'success';
    if (s === 'pending') return 'warning';
    return 'error';
  };

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-3xl font-bold text-gradient mb-8 animate-fade-in">My Profile</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {/* Profile Card */}
        <div className="md:col-span-1 animate-slide-up">
          <div className="glass-card p-6 text-center border border-white/[0.06] shadow-xl">
            <div className="w-24 h-24 bg-accent/15 rounded-full mx-auto mb-4 flex items-center justify-center border border-accent/20">
              <User className="w-12 h-12 text-accent" />
            </div>
            <h2 className="text-xl font-semibold text-slate-100 mb-1 truncate">
              {user.full_name}
            </h2>
            <p className="text-slate-400 text-sm mb-3 truncate">{user.email}</p>
            {role && (
              <span className="inline-block px-3 py-1 text-xs font-semibold bg-accent/15 text-accent-hover rounded-full mb-6">
                {role}
              </span>
            )}
            <div className="space-y-3">
              <Link href="/orders" className="block w-full">
                <Button variant="outline" className="w-full">
                  <Package className="w-4 h-4 mr-2" />
                  My Orders
                </Button>
              </Link>
              <Button
                variant="outline"
                className="w-full text-red-400 border-red-500/20 hover:bg-red-500/10 hover:border-red-500/40"
                onClick={handleLogout}
                isLoading={isLoggingOut}
              >
                <LogOut className="w-4 h-4 mr-2" />
                Logout
              </Button>
            </div>
          </div>
        </div>

        {/* Profile Details & Stats */}
        <div className="md:col-span-2 animate-slide-up" style={{ animationDelay: '0.1s' }}>
          <div className="glass-card p-6 border border-white/[0.06] shadow-xl">
            
            {/* Stats section */}
            <div className="grid grid-cols-2 gap-4 mb-8">
              <div className="bg-white/[0.02] border border-white/[0.06] p-4 rounded-xl text-center">
                <p className="text-xs text-slate-500 font-semibold uppercase tracking-wider">Total Orders</p>
                <p className="text-2xl font-bold text-slate-100 mt-1">{orders.length}</p>
              </div>
              <div className="bg-white/[0.02] border border-white/[0.06] p-4 rounded-xl text-center">
                <p className="text-xs text-slate-500 font-semibold uppercase tracking-wider">Total Spend</p>
                <p className="text-2xl font-bold text-accent-hover mt-1">{formatPrice(totalSpend)}</p>
              </div>
            </div>

            <h3 className="text-lg font-semibold text-slate-100 mb-6">Account Information</h3>
            
            <div className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-slate-400 mb-2">
                  Full Name
                </label>
                <div className="flex items-center space-x-3 p-3 bg-white/[0.04] border border-white/10 rounded-lg text-slate-200">
                  <User className="w-4 h-4 text-slate-500" />
                  <span className="text-sm font-medium">{user.full_name}</span>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-400 mb-2">
                  Email Address
                </label>
                <div className="flex items-center space-x-3 p-3 bg-white/[0.04] border border-white/10 rounded-lg text-slate-200">
                  <Mail className="w-4 h-4 text-slate-500" />
                  <span className="text-sm font-medium">{user.email}</span>
                </div>
              </div>

              {user.provider && (
                <div>
                  <label className="block text-sm font-medium text-slate-400 mb-2">
                    Login Provider
                  </label>
                  <div className="flex items-center space-x-3 p-3 bg-white/[0.04] border border-white/10 rounded-lg text-slate-200">
                    <span className="text-sm font-medium capitalize">{user.provider}</span>
                  </div>
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-slate-400 mb-2">
                  Member Since
                </label>
                <div className="flex items-center space-x-3 p-3 bg-white/[0.04] border border-white/10 rounded-lg text-slate-200">
                  <Calendar className="w-4 h-4 text-slate-500" />
                  <span className="text-sm font-medium">{formatDate(user.created_at)}</span>
                </div>
              </div>
            </div>

            {/* Recent Orders Section */}
            <div className="mt-8 border-t border-white/[0.06] pt-8">
              <h3 className="text-lg font-semibold text-slate-100 mb-4 flex items-center gap-2">
                <Award className="w-5 h-5 text-accent" />
                Recent Orders
              </h3>
              {isLoading ? (
                <div className="space-y-3">
                  {Array.from({ length: 3 }).map((_, i) => (
                    <div key={i} className="h-16 w-full bg-white/[0.06] rounded-xl animate-pulse" />
                  ))}
                </div>
              ) : recentOrders.length > 0 ? (
                <div className="space-y-3">
                  {recentOrders.map((order) => (
                    <div 
                      key={order.id} 
                      className="p-4 rounded-xl bg-white/[0.02] border border-white/[0.06] flex justify-between items-center hover:bg-white/[0.04] hover:border-white/10 transition-all cursor-pointer"
                      onClick={() => router.push(`/orders/${order.id}`)}
                    >
                      <div>
                        <p className="font-semibold text-slate-200">Order #{order.id}</p>
                        <p className="text-xs text-slate-400 mt-0.5">{formatDate(order.created_at)}</p>
                      </div>
                      <div className="text-right">
                        <p className="font-bold text-slate-200">{formatPrice(order.total_price)}</p>
                        <span className="inline-block mt-1">
                          <Badge variant={getBadgeVariant(order.status)}>
                            {order.status}
                          </Badge>
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-sm text-slate-500">No orders found.</p>
              )}
            </div>

            <div className="mt-8 pt-6 border-t border-white/[0.06]">
              <p className="text-xs text-slate-500 leading-relaxed">
                Note: Profile editing functionality is not available yet. 
                Contact support to update your information.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
