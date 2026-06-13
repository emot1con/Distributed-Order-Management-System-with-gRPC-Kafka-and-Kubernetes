'use client';

import React, { Suspense } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { ProductCard } from '@/components/ui/product-card';
import { Button } from '@/components/ui/button';
import { SkeletonCard, ErrorCard, EmptyState } from '@/components/ui';
import { productService } from '@/services/product-service';
import { useCartStore, useAuthStore } from '@/store';
import { ChevronLeft, ChevronRight, Package } from 'lucide-react';
import type { Product } from '@/types';

function LoadingProductsSkeleton() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="mb-10">
        <div className="h-9 w-48 bg-white/[0.06] rounded animate-pulse mb-2" />
        <div className="h-5 w-64 bg-white/[0.06] rounded animate-pulse" />
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {Array.from({ length: 8 }).map((_, i) => (
          <SkeletonCard key={i} />
        ))}
      </div>
    </div>
  );
}

function ProductsContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  
  const page = Number(searchParams.get('page') || '1');
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const addItem = useCartStore((state) => state.addItem);
  const getItem = useCartStore((state) => state.getItem);
  
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['products', page],
    queryFn: () => productService.getProducts(page),
  });
  
  const handleAddToCart = (product: Product) => {
    if (!isAuthenticated) {
      toast.error('Please login to add items to cart');
      return;
    }
    addItem(product);
    toast.success(`${product.name} added to cart`);
  };

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('page', String(newPage));
    router.push(`/products?${params.toString()}`);
  };
  
  if (isLoading) {
    return <LoadingProductsSkeleton />;
  }
  
  if (error) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center bg-surface px-4">
        <ErrorCard 
          message="Failed to load products. Please check your network connection and try again." 
          onRetry={refetch}
        />
      </div>
    );
  }
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="mb-10">
        <h1 className="text-3xl font-bold text-gradient">Products</h1>
        <p className="text-slate-400 mt-2">
          Browse our collection of {data?.total || 0} products
        </p>
      </div>
      
      {data?.products && data.products.length > 0 ? (
        <>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 animate-fade-in">
            {data.products.map((product: Product) => (
              <ProductCard
                key={product.id}
                product={product}
                onAddToCart={handleAddToCart}
                isInCart={!!getItem(product.id)}
              />
            ))}
          </div>
          
          {/* Pagination */}
          {data.total_page > 1 && (
            <div className="mt-12 flex items-center justify-center gap-4">
              <Button
                variant="outline"
                disabled={page <= 1}
                onClick={() => handlePageChange(Math.max(1, page - 1))}
              >
                <ChevronLeft className="w-4 h-4 mr-1" />
                Previous
              </Button>
              
              <span className="text-slate-400 text-sm font-medium">
                Page {data.page} of {data.total_page}
              </span>
              
              <Button
                variant="outline"
                disabled={page >= data.total_page}
                onClick={() => handlePageChange(page + 1)}
              >
                Next
                <ChevronRight className="w-4 h-4 ml-1" />
              </Button>
            </div>
          )}
        </>
      ) : (
        <EmptyState 
          icon={Package}
          title="No products available"
          description="We couldn't find any products in our store right now. Please come back later!"
        />
      )}
    </div>
  );
}

export default function ProductsPage() {
  return (
    <Suspense fallback={<LoadingProductsSkeleton />}>
      <ProductsContent />
    </Suspense>
  );
}
