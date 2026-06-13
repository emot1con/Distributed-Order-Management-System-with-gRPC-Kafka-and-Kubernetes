import React from 'react';
import { cn } from '@/lib/utils';

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'success' | 'warning' | 'error' | 'neutral';
}

export function Badge({ children, variant = 'neutral', className, ...props }: BadgeProps) {
  const getVariantStyles = () => {
    switch (variant) {
      case 'success':
        return 'bg-emerald-500/15 text-emerald-400';
      case 'warning':
        return 'bg-amber-500/15 text-amber-400';
      case 'error':
        return 'bg-red-500/15 text-red-400';
      case 'neutral':
      default:
        return 'bg-white/[0.06] text-slate-400';
    }
  };

  return (
    <span
      className={cn(
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border border-transparent transition-colors duration-200',
        getVariantStyles(),
        className
      )}
      {...props}
    >
      {children}
    </span>
  );
}
