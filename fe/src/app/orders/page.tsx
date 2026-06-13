'use client';

import React, { Suspense, useState } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Package, ChevronLeft, ChevronRight, ClipboardList } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { SkeletonTableRow, ErrorCard, EmptyState, Badge } from '@/components/ui';
import { useAuthStore } from '@/store';
import { orderService } from '@/services/order-service';
import { formatPrice, formatDate, cn } from '@/lib/utils';
import type { Order } from '@/types';

function LoadingOrdersSkeleton() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="h-9 w-40 bg-white/[0.06] rounded animate-pulse mb-8" />
      <div className="glass-card overflow-hidden border border-white/[0.06] shadow-xl">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-white/[0.06]">
            <thead className="bg-white/[0.02]">
              <tr>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase">Order ID</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase">Total Price</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase">Status</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase">Created At</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase">Action</th>
              </tr>
            </thead>
            <tbody>
              {Array.from({ length: 5 }).map((_, i) => (
                <SkeletonTableRow key={i} />
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function OrdersContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const [activeTab, setActiveTab] = useState<'all' | 'pending' | 'paid' | 'failed'>('all');
  
  const offset = Number(searchParams.get('offset') || '0');
  const pageSize = 15; // Match backend page size
  
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['orders', offset],
    queryFn: () => orderService.getOrders(offset),
    enabled: isAuthenticated,
    refetchOnMount: 'always',
    staleTime: 0,
  });

  const handlePageChange = (newOffset: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('offset', String(newOffset));
    router.push(`/orders?${params.toString()}`);
  };
  
  if (!isAuthenticated) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center px-4">
        <Package className="w-16 h-16 text-slate-600 mb-4 animate-pulse-glow" />
        <h2 className="text-xl font-semibold text-slate-100 mb-2">Please login first</h2>
        <p className="text-slate-400 mb-6 text-center max-w-sm">You need to be logged in to view your orders</p>
        <Link href="/login">
          <Button size="lg">Login</Button>
        </Link>
      </div>
    );
  }

  if (isLoading) {
    return <LoadingOrdersSkeleton />;
  }

  if (error) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center bg-surface px-4">
        <ErrorCard 
          message="Failed to load your orders. Please try again."
          onRetry={refetch}
        />
      </div>
    );
  }

  const orders = data?.orders || [];

  const filteredOrders = orders.filter((order) => {
    if (activeTab === 'all') return true;
    const s = order.status.toLowerCase();
    if (activeTab === 'pending') return s === 'pending';
    if (activeTab === 'paid') return s === 'paid' || s === 'success';
    if (activeTab === 'failed') return s === 'failed' || s === 'expired' || s === 'cancelled';
    return true;
  });
  
  const getBadgeVariant = (status: string) => {
    const s = status.toLowerCase();
    if (s === 'paid' || s === 'success') return 'success';
    if (s === 'pending') return 'warning';
    return 'error';
  };
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-3xl font-bold text-gradient mb-6">My Orders</h1>

      {/* Tabs */}
      <div className="flex gap-2 border-b border-white/[0.06] mb-8 pb-px overflow-x-auto scrollbar-none">
        {(['all', 'pending', 'paid', 'failed'] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={cn(
              'px-4 py-2.5 text-sm font-medium border-b-2 transition-all capitalize whitespace-nowrap',
              activeTab === tab
                ? 'border-accent text-accent'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            )}
          >
            {tab}
          </button>
        ))}
      </div>
      
      {filteredOrders.length === 0 ? (
        <EmptyState 
          icon={ClipboardList}
          title={activeTab === 'all' ? 'No orders yet' : `No ${activeTab} orders`}
          description={
            activeTab === 'all' 
              ? 'Start shopping to create your first order!' 
              : `We couldn't find any orders matching the status "${activeTab}".`
          }
          action={activeTab === 'all' ? { label: 'Browse Products', href: '/products' } : undefined}
        />
      ) : (
        <>
          {/* Mobile Layout: Stacked Cards */}
          <div className="md:hidden space-y-4 animate-fade-in">
            {filteredOrders.map((order: Order) => (
              <div key={order.id} className="glass-card p-5 border border-white/[0.06] space-y-3">
                <div className="flex justify-between items-center">
                  <span className="font-semibold text-slate-200">#{order.id}</span>
                  <Badge variant={getBadgeVariant(order.status)}>
                    {order.status}
                  </Badge>
                </div>
                <div className="flex justify-between items-center text-xs text-slate-400">
                  <span>Created: {formatDate(order.created_at)}</span>
                  <span className="text-accent-hover font-bold text-sm">
                    {formatPrice(order.total_price)}
                  </span>
                </div>
                <div className="pt-2 border-t border-white/[0.06] flex justify-end">
                  {order.status.toLowerCase() === 'pending' ? (
                    <Button 
                      variant="primary" 
                      size="sm"
                      onClick={() => router.push(`/orders/${order.id}`)}
                    >
                      Pay Now
                    </Button>
                  ) : (
                    <Link href={`/orders/${order.id}`} className="block">
                      <Button variant="outline" size="sm">View Details</Button>
                    </Link>
                  )}
                </div>
              </div>
            ))}
          </div>

          {/* Desktop Layout: Table */}
          <div className="hidden md:block glass-card overflow-hidden border border-white/[0.06] shadow-xl animate-fade-in">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-white/[0.06]">
                <thead className="bg-white/[0.02]">
                  <tr>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">
                      Order ID
                    </th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">
                      Total Price
                    </th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">
                      Created At
                    </th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">
                      Action
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/[0.06]">
                  {filteredOrders.map((order: Order) => (
                    <tr key={order.id} className="hover:bg-white/[0.02] transition-colors">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="font-semibold text-slate-200">#{order.id}</span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="text-accent-hover font-bold">
                          {formatPrice(order.total_price)}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Badge variant={getBadgeVariant(order.status)}>
                          {order.status}
                        </Badge>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-400">
                        {formatDate(order.created_at)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {order.status.toLowerCase() === 'pending' ? (
                          <Button 
                            variant="primary" 
                            size="sm"
                            onClick={() => router.push(`/orders/${order.id}`)}
                          >
                            Pay Now
                          </Button>
                        ) : (
                          <Link href={`/orders/${order.id}`}>
                            <Button variant="outline" size="sm">View</Button>
                          </Link>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
          
          {/* Pagination */}
          <div className="flex items-center justify-between mt-8">
            <Button
              variant="outline"
              onClick={() => handlePageChange(Math.max(0, offset - pageSize))}
              disabled={offset === 0}
            >
              <ChevronLeft className="w-4 h-4 mr-1" />
              Previous
            </Button>
            <span className="text-sm text-slate-400 font-medium">
              Showing {offset + 1} - {offset + orders.length}
            </span>
            <Button
              variant="outline"
              onClick={() => handlePageChange(offset + pageSize)}
              disabled={orders.length < pageSize}
            >
              Next
              <ChevronRight className="w-4 h-4 ml-1" />
            </Button>
          </div>
        </>
      )}
    </div>
  );
}

export default function OrdersPage() {
  return (
    <Suspense fallback={<LoadingOrdersSkeleton />}>
      <OrdersContent />
    </Suspense>
  );
}
