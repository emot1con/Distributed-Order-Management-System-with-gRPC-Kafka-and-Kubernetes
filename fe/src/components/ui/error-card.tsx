import React from 'react';
import { AlertCircle } from 'lucide-react';
import { Button } from './button';

interface ErrorCardProps {
  message: string;
  onRetry?: () => void;
}

export function ErrorCard({ message, onRetry }: ErrorCardProps) {
  return (
    <div className="glass-card p-8 text-center border border-red-500/10 shadow-xl max-w-md mx-auto my-8 bg-red-500/[0.02]">
      <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-red-500/10 mb-4 text-red-400">
        <AlertCircle className="w-6 h-6" />
      </div>
      <h3 className="text-lg font-semibold text-slate-100 mb-2">Error Occurred</h3>
      <p className="text-slate-400 text-sm mb-6">{message}</p>
      {onRetry && (
        <Button variant="outline" size="sm" onClick={onRetry}>
          Try Again
        </Button>
      )}
    </div>
  );
}
