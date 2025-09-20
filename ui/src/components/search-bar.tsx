'use client';

import { motion, AnimatePresence } from 'framer-motion';
import { Search, X, Filter } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { useState, useEffect } from 'react';

interface SearchBarProps {
  searchTerm: string;
  onSearchChange: (term: string) => void;
  selectedCategory: string;
  onCategoryChange: (category: string) => void;
  selectedPlatform: string;
  onPlatformChange: (platform: string) => void;
}

const categories = [
  { value: '', label: 'All Categories' },
  { value: 'os', label: 'Operating Systems', icon: '💿' },
  { value: 'app', label: 'Applications', icon: '📱' },
  { value: 'tool', label: 'Tools', icon: '🛠️' },
];

const platforms = [
  { value: '', label: 'All Platforms' },
  { value: 'windows', label: 'Windows', icon: '🪟' },
  { value: 'linux', label: 'Linux', icon: '🐧' },
  { value: 'macos', label: 'macOS', icon: '🍎' },
  { value: 'web', label: 'Web', icon: '🌐' },
];

export function SearchBar({
  searchTerm,
  onSearchChange,
  selectedCategory,
  onCategoryChange,
  selectedPlatform,
  onPlatformChange,
}: SearchBarProps) {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [isFocused, setIsFocused] = useState(false);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === '/' && !e.ctrlKey && !e.metaKey && !e.altKey) {
        e.preventDefault();
        const searchInput = document.getElementById('search-input') as HTMLInputElement;
        searchInput?.focus();
      }
      if (e.key === 'Escape') {
        const searchInput = document.getElementById('search-input') as HTMLInputElement;
        searchInput?.blur();
        setIsFilterOpen(false);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, []);

  const hasActiveFilters = selectedCategory || selectedPlatform;

  return (
    <div className="w-full max-w-4xl mx-auto space-y-4">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="relative"
      >
        <div className={`relative transition-all duration-300 ${isFocused ? 'scale-105' : ''}`}>
          <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-muted-foreground w-5 h-5" />
          <Input
            id="search-input"
            type="text"
            placeholder="Search for software, OS, or tools... (Press '/' to focus)"
            value={searchTerm}
            onChange={(e) => onSearchChange(e.target.value)}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            className="pl-12 pr-24 h-14 text-lg glass border-white/20 focus:border-primary/50 focus:ring-primary/25"
          />
          <div className="absolute right-2 top-1/2 transform -translate-y-1/2 flex items-center space-x-2">
            {searchTerm && (
              <Button
                size="icon"
                variant="ghost"
                onClick={() => onSearchChange('')}
                className="h-8 w-8 rounded-full hover:bg-white/10"
              >
                <X className="w-4 h-4" />
              </Button>
            )}
            <Button
              size="sm"
              variant={hasActiveFilters ? 'default' : 'glass'}
              onClick={() => setIsFilterOpen(!isFilterOpen)}
              className={`relative ${hasActiveFilters ? 'bg-primary' : ''}`}
            >
              <Filter className="w-4 h-4 mr-2" />
              Filter
              {hasActiveFilters && (
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  className="absolute -top-1 -right-1 w-3 h-3 bg-pink-500 rounded-full"
                />
              )}
            </Button>
          </div>
        </div>

        <AnimatePresence>
          {isFilterOpen && (
            <motion.div
              initial={{ opacity: 0, y: -10, scale: 0.95 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -10, scale: 0.95 }}
              transition={{ duration: 0.2 }}
              className="absolute top-full mt-2 w-full z-50"
            >
              <Card className="glass border-white/20 p-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium mb-2 block">Category</label>
                    <div className="space-y-1">
                      {categories.map((category) => (
                        <motion.button
                          key={category.value}
                          whileHover={{ scale: 1.02 }}
                          whileTap={{ scale: 0.98 }}
                          onClick={() => onCategoryChange(category.value)}
                          className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors flex items-center space-x-2 ${
                            selectedCategory === category.value
                              ? 'bg-primary text-primary-foreground'
                              : 'hover:bg-white/10'
                          }`}
                        >
                          {category.icon && <span>{category.icon}</span>}
                          <span>{category.label}</span>
                        </motion.button>
                      ))}
                    </div>
                  </div>

                  <div>
                    <label className="text-sm font-medium mb-2 block">Platform</label>
                    <div className="space-y-1">
                      {platforms.map((platform) => (
                        <motion.button
                          key={platform.value}
                          whileHover={{ scale: 1.02 }}
                          whileTap={{ scale: 0.98 }}
                          onClick={() => onPlatformChange(platform.value)}
                          className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors flex items-center space-x-2 ${
                            selectedPlatform === platform.value
                              ? 'bg-primary text-primary-foreground'
                              : 'hover:bg-white/10'
                          }`}
                        >
                          {platform.icon && <span>{platform.icon}</span>}
                          <span>{platform.label}</span>
                        </motion.button>
                      ))}
                    </div>
                  </div>
                </div>

                {hasActiveFilters && (
                  <div className="mt-4 pt-4 border-t border-white/10">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        onCategoryChange('');
                        onPlatformChange('');
                      }}
                      className="w-full"
                    >
                      Clear all filters
                    </Button>
                  </div>
                )}
              </Card>
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>

      {hasActiveFilters && (
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          className="flex flex-wrap gap-2"
        >
          {selectedCategory && (
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              exit={{ scale: 0 }}
              className="flex items-center space-x-2 bg-primary/20 text-primary px-3 py-1 rounded-full text-sm"
            >
              <span>{categories.find(c => c.value === selectedCategory)?.label}</span>
              <button
                onClick={() => onCategoryChange('')}
                className="hover:bg-primary/30 rounded-full p-1"
              >
                <X className="w-3 h-3" />
              </button>
            </motion.div>
          )}
          {selectedPlatform && (
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              exit={{ scale: 0 }}
              className="flex items-center space-x-2 bg-primary/20 text-primary px-3 py-1 rounded-full text-sm"
            >
              <span>{platforms.find(p => p.value === selectedPlatform)?.label}</span>
              <button
                onClick={() => onPlatformChange('')}
                className="hover:bg-primary/30 rounded-full p-1"
              >
                <X className="w-3 h-3" />
              </button>
            </motion.div>
          )}
        </motion.div>
      )}
    </div>
  );
}