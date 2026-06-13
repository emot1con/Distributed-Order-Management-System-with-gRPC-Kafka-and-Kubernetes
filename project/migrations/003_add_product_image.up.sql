-- Migration: Add image_url column to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS image_url TEXT DEFAULT '';

-- Seed existing products with name-specific high quality Unsplash images
UPDATE products 
SET image_url = CASE 
  WHEN name = 'bag' AND id % 3 = 0 THEN 'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?auto=format&fit=crop&w=400&h=300&q=80'
  WHEN name = 'bag' AND id % 3 = 1 THEN 'https://images.unsplash.com/photo-1544816155-12df9643f363?auto=format&fit=crop&w=400&h=300&q=80'
  WHEN name = 'bag' AND id % 3 = 2 THEN 'https://images.unsplash.com/photo-1622560480605-d83c853bc5c3?auto=format&fit=crop&w=400&h=300&q=80'
  WHEN name = 'pencil' AND id % 3 = 0 THEN 'https://images.unsplash.com/photo-1513542789411-b6a5d4f31634?auto=format&fit=crop&w=400&h=300&q=80'
  WHEN name = 'pencil' AND id % 3 = 1 THEN 'https://images.unsplash.com/photo-1569003339405-ea396a5a8a90?auto=format&fit=crop&w=400&h=300&q=80'
  WHEN name = 'pencil' AND id % 3 = 2 THEN 'https://images.unsplash.com/photo-1519750783826-e2420f4d687f?auto=format&fit=crop&w=400&h=300&q=80'
  ELSE 'https://images.unsplash.com/photo-1586075010923-2dd4570fb338?auto=format&fit=crop&w=400&h=300&q=80' -- generic placeholder
END;
