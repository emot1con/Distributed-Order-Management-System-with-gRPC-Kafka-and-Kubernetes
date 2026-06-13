'use client';

import React from 'react';
import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { useAuthStore, useCartStore } from '@/store';
import { productService } from '@/services/product-service';
import { ProductCard, SkeletonCard } from '@/components/ui';
import { ArrowRight } from 'lucide-react';
import type { Product } from '@/types';

export default function FeaturedProducts() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const addItem = useCartStore((state) => state.addItem);
  const getItem = useCartStore((state) => state.getItem);

  const { data, isLoading, error } = useQuery({
    queryKey: ['featured-products'],
    queryFn: () => productService.getProducts(1),
  });

  const handleAddToCart = (product: Product) => {
    if (!isAuthenticated) {
      toast.error('Please login to add items to cart');
      return;
    }
    addItem(product);
    toast.success(`${product.name} added to cart`);
  };

  if (error) return null; // Hide the section if there's an API error

  const featured = data?.products?.slice(0, 4) || [];

  return (
    <section className="py-20 bg-surface/50 border-t border-b border-white/[0.06] relative">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div className="flex flex-col sm:flex-row items-start sm:items-end justify-between mb-12">
          <div>
            <h2 className="text-3xl font-bold text-gradient">Featured Products</h2>
            <p className="text-slate-400 mt-2">Discover some of our trending high-quality items</p>
          </div>
          <Link href="/products" className="mt-4 sm:mt-0 group inline-flex items-center text-accent hover:text-accent-hover font-semibold transition-colors">
            View All Products
            <ArrowRight className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
          </Link>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {Array.from({ length: 4 }).map((_, i) => (
              <SkeletonCard key={i} />
            ))}
          </div>
        ) : featured.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 animate-fade-in">
            {featured.map((product: Product) => (
              <ProductCard
                key={product.id}
                product={product}
                onAddToCart={handleAddToCart}
                isInCart={!!getItem(product.id)}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-12 text-slate-500">
            No products found
          </div>
        )}
      </div>
    </section>
  );
}
