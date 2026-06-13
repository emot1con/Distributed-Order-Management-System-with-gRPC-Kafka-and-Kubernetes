'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ShoppingCart, User, LogOut, Menu, X, Store, Settings } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAuthStore, useCartStore } from '@/store';
import { Button } from '@/components/ui/button';
import { config } from '@/lib/config';

export function Navbar() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const pathname = usePathname();
  const { isAuthenticated, logout, user } = useAuthStore();
  const cartItems = useCartStore((state) => state.items);
  const totalItems = cartItems.reduce((sum, item) => sum + item.quantity, 0);
  
  const navLinks = [
    { href: '/', label: 'Home' },
    { href: '/products', label: 'Products' },
  ];
  
  const authLinks = isAuthenticated
    ? [
        { href: '/orders', label: 'My Orders' },
      ]
    : [];
  
  const allLinks = [...navLinks, ...authLinks];
  
  return (
    <nav className="bg-surface/80 backdrop-blur-xl border-b border-white/[0.06] sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center space-x-2">
              <Store className="w-6 h-6 text-accent" />
              <span className="font-bold text-xl text-gradient">
                {config.appName}
              </span>
            </Link>
          </div>
          
          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-8">
            {allLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={cn(
                  'text-sm font-medium transition-colors',
                  pathname === link.href
                    ? 'text-accent'
                    : 'text-slate-400 hover:text-slate-100'
                )}
              >
                {link.label}
              </Link>
            ))}
          </div>
          
          {/* Right side actions */}
          <div className="hidden md:flex items-center space-x-4">
            {/* Cart always visible if logged in */}
            {isAuthenticated && (
              <Link href="/cart" className="relative p-2 text-slate-400 hover:text-slate-100 transition-colors">
                <ShoppingCart className="w-6 h-6" />
                {totalItems > 0 && (
                  <span className="absolute -top-1 -right-1 bg-accent text-white text-xs rounded-full w-5 h-5 flex items-center justify-center font-bold">
                    {totalItems > 99 ? '99+' : totalItems}
                  </span>
                )}
              </Link>
            )}

            {isAuthenticated ? (
              /* User Dropdown */
              <div className="relative">
                <button
                  onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                  className="flex items-center space-x-2 p-1 rounded-full hover:bg-white/[0.06] border border-transparent hover:border-white/10 transition-all focus:outline-none"
                >
                  <div className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center text-accent font-semibold text-sm uppercase">
                    {user?.full_name?.charAt(0) || user?.email?.charAt(0) || 'U'}
                  </div>
                </button>

                {isDropdownOpen && (
                  <>
                    {/* Backdrop to close */}
                    <div className="fixed inset-0 z-10" onClick={() => setIsDropdownOpen(false)} />
                    <div className="absolute right-0 mt-2 w-48 glass-card border border-white/[0.06] py-1 shadow-2xl rounded-xl z-20 animate-scale-up">
                      <div className="px-4 py-2 border-b border-white/[0.06]">
                        <p className="text-[10px] text-slate-500 font-bold uppercase tracking-wider">Logged in as</p>
                        <p className="text-sm font-semibold text-slate-200 truncate mt-0.5">{user?.full_name || 'User'}</p>
                      </div>
                      <Link
                        href="/profile"
                        className="flex items-center px-4 py-2.5 text-sm text-slate-400 hover:text-slate-200 hover:bg-white/[0.04] transition-colors"
                        onClick={() => setIsDropdownOpen(false)}
                      >
                        <User className="w-4 h-4 mr-2" />
                        My Profile
                      </Link>
                      <Link
                        href="/orders"
                        className="flex items-center px-4 py-2.5 text-sm text-slate-400 hover:text-slate-200 hover:bg-white/[0.04] transition-colors"
                        onClick={() => setIsDropdownOpen(false)}
                      >
                        <ShoppingCart className="w-4 h-4 mr-2" />
                        My Orders
                      </Link>
                      <hr className="border-white/[0.06] my-1" />
                      <button
                        onClick={() => {
                          logout();
                          setIsDropdownOpen(false);
                        }}
                        className="flex items-center w-full px-4 py-2.5 text-sm text-red-400 hover:text-red-300 hover:bg-red-500/10 transition-colors text-left"
                      >
                        <LogOut className="w-4 h-4 mr-2" />
                        Logout
                      </button>
                    </div>
                  </>
                )}
              </div>
            ) : (
              <>
                <Link href="/login">
                  <Button variant="ghost" size="sm">
                    Login
                  </Button>
                </Link>
                <Link href="/register">
                  <Button variant="primary" size="sm">
                    Register
                  </Button>
                </Link>
              </>
            )}
          </div>
          
          {/* Mobile menu button */}
          <div className="md:hidden flex items-center">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="p-2 text-slate-400 hover:text-slate-100"
            >
              {isMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>
          </div>
        </div>
      </div>
      
      {/* Mobile menu */}
      {isMenuOpen && (
        <div className="md:hidden bg-surface-card border-t border-white/[0.06]">
          <div className="px-4 py-4 space-y-3">
            {allLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={cn(
                  'block px-3 py-2 rounded-md text-base font-medium transition-colors',
                  pathname === link.href
                    ? 'bg-accent/10 text-accent'
                    : 'text-slate-400 hover:bg-white/[0.04]'
                )}
                onClick={() => setIsMenuOpen(false)}
              >
                {link.label}
              </Link>
            ))}
            
            <hr className="my-2 border-white/[0.06]" />
            
            {isAuthenticated ? (
              <>
                <Link
                  href="/profile"
                  className="flex items-center px-3 py-2 text-slate-400 hover:text-slate-100"
                  onClick={() => setIsMenuOpen(false)}
                >
                  <User className="w-5 h-5 mr-2" />
                  My Profile
                </Link>
                <Link
                  href="/cart"
                  className="flex items-center px-3 py-2 text-slate-400 hover:text-slate-100"
                  onClick={() => setIsMenuOpen(false)}
                >
                  <ShoppingCart className="w-5 h-5 mr-2" />
                  Cart ({totalItems})
                </Link>
                <button
                  onClick={() => {
                    logout();
                    setIsMenuOpen(false);
                  }}
                  className="flex items-center w-full px-3 py-2 text-slate-400 hover:text-slate-100 text-left"
                >
                  <LogOut className="w-5 h-5 mr-2" />
                  Logout
                </button>
              </>
            ) : (
              <div className="space-y-2">
                <Link href="/login" onClick={() => setIsMenuOpen(false)}>
                  <Button variant="outline" className="w-full">Login</Button>
                </Link>
                <Link href="/register" onClick={() => setIsMenuOpen(false)}>
                  <Button variant="primary" className="w-full">Register</Button>
                </Link>
              </div>
            )}
          </div>
        </div>
      )}
    </nav>
  );
}
