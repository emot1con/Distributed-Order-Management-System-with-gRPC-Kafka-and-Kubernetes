import React from 'react';
import Link from 'next/link';
import { Home, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';

export default function NotFound() {
  return (
    <div className="min-h-[70vh] flex items-center justify-center px-4 relative overflow-hidden">
      {/* Background glow */}
      <div className="absolute inset-0 bg-hero-glow opacity-40 pointer-events-none" />
      
      <div className="glass-card max-w-md w-full p-8 text-center border border-white/[0.06] shadow-2xl relative z-10">
        <div className="text-6xl font-extrabold text-gradient mb-4">404</div>
        <h2 className="text-2xl font-bold text-slate-100 mb-3">Page Not Found</h2>
        <p className="text-slate-400 text-sm mb-8 leading-relaxed">
          The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.
        </p>
        <div className="flex flex-col sm:flex-row gap-3 justify-center">
          <Link href="/" className="flex-1">
            <Button variant="outline" className="w-full">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back Home
            </Button>
          </Link>
          <Link href="/products" className="flex-1">
            <Button className="w-full">
              <Home className="w-4 h-4 mr-2" />
              Browse Products
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
