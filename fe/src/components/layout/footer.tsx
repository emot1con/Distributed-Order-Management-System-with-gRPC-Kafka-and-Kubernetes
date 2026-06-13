import React from 'react';
import Link from 'next/link';
import { Store } from 'lucide-react';
import { config } from '@/lib/config';

export function Footer() {
  return (
    <footer className="bg-surface border-t border-white/[0.06]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Brand */}
          <div className="col-span-1 md:col-span-2">
            <div className="flex items-center space-x-2 mb-4">
              <Store className="w-6 h-6 text-accent" />
              <span className="font-bold text-xl text-slate-100">
                {config.appName}
              </span>
            </div>
            <p className="text-slate-400 text-sm max-w-md">
              {config.appDescription}. Curated premium products with instant payment processing and zero-delay delivery routing.
            </p>
          </div>
          
          {/* Quick Links */}
          <div>
            <h3 className="font-semibold text-slate-200 mb-4">Quick Links</h3>
            <ul className="space-y-2">
              <li>
                <Link href="/" className="text-slate-400 hover:text-slate-200 text-sm transition-colors">
                  Home
                </Link>
              </li>
              <li>
                <Link href="/products" className="text-slate-400 hover:text-slate-200 text-sm transition-colors">
                  Products
                </Link>
              </li>
              <li>
                <Link href="/cart" className="text-slate-400 hover:text-slate-200 text-sm transition-colors">
                  Cart
                </Link>
              </li>
            </ul>
          </div>
          
          {/* Account */}
          <div>
            <h3 className="font-semibold text-slate-200 mb-4">Account</h3>
            <ul className="space-y-2">
              <li>
                <Link href="/login" className="text-slate-400 hover:text-slate-200 text-sm transition-colors">
                  Login
                </Link>
              </li>
              <li>
                <Link href="/register" className="text-slate-400 hover:text-slate-200 text-sm transition-colors">
                  Register
                </Link>
              </li>
              <li>
                <Link href="/orders" className="text-slate-400 hover:text-slate-200 text-sm transition-colors">
                  My Orders
                </Link>
              </li>
            </ul>
          </div>
        </div>
        
        <div className="border-t border-white/[0.06] mt-8 pt-8 text-center text-sm text-slate-500">
          <p>&copy; {new Date().getFullYear()} {config.appName}. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
}
