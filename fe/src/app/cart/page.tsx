'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Trash2, Plus, Minus, ShoppingBag, Package, AlertTriangle } from 'lucide-react';
import toast from 'react-hot-toast';
import { Button } from '@/components/ui/button';
import { Modal } from '@/components/ui';
import { useCartStore, useAuthStore } from '@/store';
import { formatPrice } from '@/lib/utils';

export default function CartPage() {
  const router = useRouter();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const items = useCartStore((state) => state.items);
  const removeItem = useCartStore((state) => state.removeItem);
  const updateQuantity = useCartStore((state) => state.updateQuantity);
  const totalPrice = useCartStore((state) => state.totalPrice);
  const clearCart = useCartStore((state) => state.clearCart);

  // Modal State
  const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null);
  const [confirmClear, setConfirmClear] = useState(false);
  
  if (!isAuthenticated) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center px-4">
        <ShoppingBag className="w-16 h-16 text-slate-600 mb-4 animate-pulse-glow" />
        <h2 className="text-xl font-semibold text-slate-100 mb-2">Please login first</h2>
        <p className="text-slate-400 mb-6 text-center max-w-sm">You need to be logged in to view your cart</p>
        <Link href="/login">
          <Button size="lg">Login</Button>
        </Link>
      </div>
    );
  }
  
  if (items.length === 0) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center px-4">
        <ShoppingBag className="w-16 h-16 text-slate-600 mb-4 animate-pulse-glow" />
        <h2 className="text-xl font-semibold text-slate-100 mb-2">Your cart is empty</h2>
        <p className="text-slate-400 mb-6 text-center max-w-sm">Add some products to get started</p>
        <Link href="/products">
          <Button size="lg">Browse Products</Button>
        </Link>
      </div>
    );
  }
  
  const handleCheckout = () => {
    router.push('/checkout');
  };

  const handleRemoveClick = (id: number) => {
    setConfirmDeleteId(id);
  };

  const confirmRemove = () => {
    if (confirmDeleteId !== null) {
      removeItem(confirmDeleteId);
      toast.success('Item removed from cart');
      setConfirmDeleteId(null);
    }
  };

  const confirmClearCart = () => {
    clearCart();
    toast.success('Cart cleared');
    setConfirmClear(false);
  };

  const itemToDelete = items.find(item => item.product.id === confirmDeleteId);
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
        <h1 className="text-3xl font-bold text-gradient">Shopping Cart</h1>
        <Button 
          variant="outline" 
          size="sm" 
          className="text-red-400 hover:text-red-300 border-red-500/20 hover:border-red-500/40 hover:bg-red-500/10"
          onClick={() => setConfirmClear(true)}
        >
          <Trash2 className="w-4 h-4 mr-2" />
          Clear Cart
        </Button>
      </div>
      
      <div className="lg:grid lg:grid-cols-12 lg:gap-8">
        {/* Cart Items */}
        <div className="lg:col-span-8">
          <div className="glass-card divide-y divide-white/[0.06] overflow-hidden animate-slide-up">
            {items.map((item) => (
              <div key={item.product.id} className="p-6 flex flex-col sm:flex-row items-start sm:items-center gap-4">
                {/* Product Image */}
                <div className="w-20 h-20 bg-gradient-to-br from-accent/10 to-surface-elevated rounded-lg flex items-center justify-center flex-shrink-0 border border-white/[0.06] overflow-hidden relative">
                  {item.product.image_url ? (
                    <img 
                      src={item.product.image_url} 
                      alt={item.product.name}
                      className="absolute inset-0 w-full h-full object-cover"
                      onError={(e) => {
                        (e.target as HTMLImageElement).style.display = 'none';
                        (e.target as HTMLImageElement).nextElementSibling?.classList.remove('hidden');
                      }}
                    />
                  ) : null}
                  <Package className={`w-8 h-8 text-slate-500 relative z-0 ${item.product.image_url ? "hidden" : ""}`} />
                </div>
                
                {/* Product Info */}
                <div className="flex-1 min-w-0">
                  <h3 className="text-lg font-medium text-slate-100 truncate">
                    {item.product.name}
                  </h3>
                  <p className="text-sm text-slate-400 mt-1">
                    {formatPrice(item.product.price)} each
                  </p>
                  <p className="text-xs text-slate-500 mt-1">
                    Available stock: <span className="font-semibold text-slate-400">{item.product.stock}</span>
                  </p>
                </div>
                
                {/* Controls and pricing */}
                <div className="flex items-center justify-between w-full sm:w-auto gap-6 sm:gap-4 mt-4 sm:mt-0">
                  {/* Quantity Controls */}
                  <div className="flex flex-col items-center gap-1">
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => {
                          if (item.quantity > 1) {
                            updateQuantity(item.product.id, item.quantity - 1);
                          } else {
                            setConfirmDeleteId(item.product.id);
                          }
                        }}
                        className="p-1 rounded-md hover:bg-white/[0.06] text-slate-400 hover:text-slate-100 transition-colors"
                      >
                        <Minus className="w-4 h-4" />
                      </button>
                      <span className="w-8 text-center font-medium text-slate-200">{item.quantity}</span>
                      <button
                        onClick={() => {
                          if (item.quantity < item.product.stock) {
                            updateQuantity(item.product.id, item.quantity + 1);
                          } else {
                            toast.error(`Only ${item.product.stock} units available in stock`);
                          }
                        }}
                        className="p-1 rounded-md hover:bg-white/[0.06] text-slate-400 hover:text-slate-100 transition-colors"
                      >
                        <Plus className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                  
                  {/* Item Total */}
                  <div className="text-right min-w-[80px]">
                    <p className="font-semibold text-slate-200">
                      {formatPrice(item.product.price * item.quantity)}
                    </p>
                  </div>
                  
                  {/* Remove Button */}
                  <button
                    onClick={() => handleRemoveClick(item.product.id)}
                    className="p-2 text-red-400/80 hover:bg-red-500/10 hover:text-red-400 rounded-md transition-colors"
                    title="Remove item"
                  >
                    <Trash2 className="w-5 h-5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
        
        {/* Order Summary */}
        <div className="lg:col-span-4 mt-8 lg:mt-0">
          <div className="glass-card p-6 sticky top-24 shadow-xl border border-white/[0.06] animate-slide-up" style={{ animationDelay: '0.1s' }}>
            <h2 className="text-lg font-semibold text-slate-100 mb-4">Order Summary</h2>
            
            <div className="space-y-4 mb-6">
              <div className="flex justify-between text-slate-400 text-sm">
                <span>Subtotal ({items.reduce((s, i) => s + i.quantity, 0)} items)</span>
                <span className="text-slate-200">{formatPrice(totalPrice())}</span>
              </div>
              <div className="flex justify-between text-slate-400 text-sm">
                <span>Shipping</span>
                <span className="text-emerald-400 font-medium">Free</span>
              </div>
              <hr className="border-white/[0.06]" />
              <div className="flex justify-between text-base font-semibold text-slate-100">
                <span>Total</span>
                <span className="text-accent-hover">{formatPrice(totalPrice())}</span>
              </div>
            </div>
            
            <Button className="w-full" onClick={handleCheckout}>
              Proceed to Checkout
            </Button>
            
            <Link href="/products" className="block mt-3">
              <Button variant="outline" className="w-full">
                Continue Shopping
              </Button>
            </Link>
          </div>
        </div>
      </div>

      {/* Confirmation Modals */}
      <Modal
        isOpen={confirmDeleteId !== null}
        onClose={() => setConfirmDeleteId(null)}
        title="Remove Item?"
        footer={
          <>
            <Button variant="outline" onClick={() => setConfirmDeleteId(null)}>
              Cancel
            </Button>
            <Button variant="danger" onClick={confirmRemove}>
              Remove
            </Button>
          </>
        }
      >
        <div className="flex items-center gap-3">
          <AlertTriangle className="w-8 h-8 text-red-400 flex-shrink-0" />
          <p>Are you sure you want to remove <span className="font-semibold text-slate-200">{itemToDelete?.product.name}</span> from your cart?</p>
        </div>
      </Modal>

      <Modal
        isOpen={confirmClear}
        onClose={() => setConfirmClear(false)}
        title="Clear Shopping Cart?"
        footer={
          <>
            <Button variant="outline" onClick={() => setConfirmClear(false)}>
              Cancel
            </Button>
            <Button variant="danger" onClick={confirmClearCart}>
              Clear Cart
            </Button>
          </>
        }
      >
        <div className="flex items-center gap-3">
          <AlertTriangle className="w-8 h-8 text-red-400 flex-shrink-0" />
          <p>Are you sure you want to remove all items from your shopping cart? This action cannot be undone.</p>
        </div>
      </Modal>
    </div>
  );
}
