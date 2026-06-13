import React from 'react';
import { cn, formatPrice } from '@/lib/utils';
import { Button } from './button';
import { ShoppingCart, Package } from 'lucide-react';
import type { Product } from '@/types';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product) => void;
  isInCart?: boolean;
}

export function ProductCard({ product, onAddToCart, isInCart }: ProductCardProps) {
  const isOutOfStock = product.stock <= 0;
  
  return (
    <div className="glass-card-hover group overflow-hidden">
      {/* Product Image Area */}
      <div className="h-48 bg-gradient-to-br from-accent/10 to-surface-elevated flex items-center justify-center relative overflow-hidden group">
        <div className="absolute inset-0 bg-gradient-to-t from-surface-card to-transparent z-10" />
        {product.image_url ? (
          <img 
            src={product.image_url} 
            alt={product.name}
            className="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
            onError={(e) => {
              (e.target as HTMLImageElement).style.display = 'none';
              (e.target as HTMLImageElement).nextElementSibling?.classList.remove('hidden');
            }}
          />
        ) : null}
        <Package className={cn(
          "w-12 h-12 text-slate-600 transition-colors duration-300 relative z-0",
          product.image_url ? "hidden group-hover:text-accent/60" : "group-hover:text-accent/60"
        )} />
      </div>
      
      {/* Product Info */}
      <div className="p-5">
        <h3 className="text-base font-semibold text-slate-100 mb-1 line-clamp-1">
          {product.name}
        </h3>
        <p className="text-sm text-slate-400 mb-4 line-clamp-2 h-10">
          {product.description || 'No description available'}
        </p>
        
        <div className="flex items-center justify-between mb-4">
          <span className="text-lg font-bold text-accent-hover">
            {formatPrice(product.price)}
          </span>
          <span className={cn(
            'text-xs px-2.5 py-1 rounded-full font-medium',
            isOutOfStock
              ? 'bg-red-500/15 text-red-400'
              : product.stock < 10
                ? 'bg-amber-500/15 text-amber-400'
                : 'bg-emerald-500/15 text-emerald-400'
          )}>
            {isOutOfStock ? 'Sold Out' : `${product.stock} left`}
          </span>
        </div>
        
        <Button
          variant={isInCart ? 'secondary' : 'primary'}
          className="w-full"
          disabled={isOutOfStock}
          onClick={() => onAddToCart?.(product)}
        >
          <ShoppingCart className="w-4 h-4 mr-2" />
          {isInCart ? 'In Cart' : 'Add to Cart'}
        </Button>
      </div>
    </div>
  );
}
