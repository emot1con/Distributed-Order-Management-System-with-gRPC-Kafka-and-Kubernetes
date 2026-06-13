'use client';

import React, { useEffect } from 'react';
import { ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function Error({ error, reset }: ErrorProps) {
  useEffect(() => {
    // Log the error to an error reporting service if needed
    console.error('Next.js Global Error:', error);
  }, [error]);

  return (
    <div className="min-h-[70vh] flex items-center justify-center px-4 relative overflow-hidden">
      {/* Background glow */}
      <div className="absolute inset-0 bg-hero-glow opacity-40 pointer-events-none" />
      
      <div className="glass-card max-w-md w-full p-8 text-center border border-red-500/10 shadow-2xl relative z-10 bg-red-500/[0.01]">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-red-500/10 mb-6 text-red-400 animate-pulse-glow">
          <span className="text-2xl font-bold">!</span>
        </div>
        <h2 className="text-2xl font-bold text-gradient mb-3">Something went wrong</h2>
        <p className="text-slate-400 text-sm mb-8 leading-relaxed">
          An unexpected error occurred while processing your request. Please try reloading the page or go back home.
        </p>
        <div className="flex flex-col sm:flex-row gap-3 justify-center">
          <Button variant="outline" onClick={() => window.location.href = '/'} className="flex-1">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Home
          </Button>
          <Button onClick={reset} className="flex-1">
            Try Again
          </Button>
        </div>
      </div>
    </div>
  );
}
