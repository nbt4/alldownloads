'use client';

import { motion } from 'framer-motion';
import { ExternalLink, Download, Copy, Clock, Package } from 'lucide-react';
import { Product, ProductVersion } from '@/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  formatFileSize,
  getTimeAgo,
  copyToClipboard,
  getPlatformIcon,
  getArchitectureLabel,
  getFreshnessColor,
  getCategoryIcon
} from '@/lib/utils';
import { useState } from 'react';

interface ProductCardProps {
  product: Product;
  versions?: ProductVersion[];
  onViewDetails: (productId: string) => void;
}

export function ProductCard({ product, versions = [], onViewDetails }: ProductCardProps) {
  const [copied, setCopied] = useState<string | null>(null);

  const latestVersions = versions.filter(v => v.is_latest);
  const lastUpdated = versions.length > 0
    ? Math.max(...versions.map(v => new Date(v.last_fetched).getTime()))
    : new Date(product.updated_at).getTime();

  const handleCopy = async (text: string, type: string) => {
    try {
      await copyToClipboard(text);
      setCopied(type);
      setTimeout(() => setCopied(null), 2000);
    } catch (error) {
      console.error('Failed to copy:', error);
    }
  };

  const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.5,
        ease: 'easeOut'
      }
    },
    hover: {
      y: -5,
      transition: {
        duration: 0.2,
        ease: 'easeOut'
      }
    }
  };

  return (
    <motion.div
      variants={cardVariants}
      initial="hidden"
      animate="visible"
      whileHover="hover"
      className="h-full"
    >
      <Card className="glass glow-hover h-full flex flex-col group cursor-pointer"
            onClick={() => onViewDetails(product.id)}>
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-primary/20 to-pink-500/20 flex items-center justify-center text-2xl">
                {getCategoryIcon(product.category)}
              </div>
              <div>
                <CardTitle className="text-lg leading-none mb-1">
                  {product.name}
                </CardTitle>
                <CardDescription className="text-sm">
                  {product.vendor}
                </CardDescription>
              </div>
            </div>
            <motion.div
              className={`text-xs px-2 py-1 rounded-full ${getFreshnessColor(new Date(lastUpdated).toISOString())}`}
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.3 }}
            >
              <Clock className="w-3 h-3 inline mr-1" />
              {getTimeAgo(new Date(lastUpdated).toISOString())}
            </motion.div>
          </div>
        </CardHeader>

        <CardContent className="flex-1 flex flex-col">
          <p className="text-sm text-muted-foreground mb-4 flex-1">
            {product.description}
          </p>

          {latestVersions.length > 0 && (
            <div className="space-y-3">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Latest versions:</span>
                <span className="text-primary font-mono">
                  {latestVersions[0]?.version}
                </span>
              </div>

              <div className="grid grid-cols-2 gap-2">
                {latestVersions.slice(0, 4).map((version, index) => (
                  <motion.div
                    key={version.id}
                    initial={{ opacity: 0, x: -10 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 0.1 * index }}
                    className="flex items-center justify-between p-2 rounded-md bg-secondary/50 hover:bg-secondary/80 transition-colors group/version"
                  >
                    <div className="flex items-center space-x-2 text-xs">
                      <span>{getPlatformIcon(version.platform)}</span>
                      <span className="text-muted-foreground">
                        {getArchitectureLabel(version.architecture)}
                      </span>
                    </div>
                    <div className="flex items-center space-x-1 opacity-0 group-hover/version:opacity-100 transition-opacity">
                      <Button
                        size="icon"
                        variant="ghost"
                        className="h-6 w-6"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleCopy(version.download_url, `download-${version.id}`);
                        }}
                      >
                        {copied === `download-${version.id}` ? (
                          <motion.div
                            initial={{ scale: 0 }}
                            animate={{ scale: 1 }}
                            className="text-green-500"
                          >
                            ✓
                          </motion.div>
                        ) : (
                          <Copy className="w-3 h-3" />
                        )}
                      </Button>
                      <Button
                        size="icon"
                        variant="ghost"
                        className="h-6 w-6"
                        onClick={(e) => {
                          e.stopPropagation();
                          window.open(version.download_url, '_blank');
                        }}
                      >
                        <Download className="w-3 h-3" />
                      </Button>
                    </div>
                  </motion.div>
                ))}
              </div>

              {latestVersions.length > 4 && (
                <p className="text-xs text-muted-foreground text-center">
                  +{latestVersions.length - 4} more versions
                </p>
              )}
            </div>
          )}

          <div className="flex items-center justify-between mt-4 pt-3 border-t border-border">
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                window.open(product.website_url, '_blank');
              }}
              className="text-primary hover:text-primary/80"
            >
              <ExternalLink className="w-4 h-4 mr-2" />
              Visit Site
            </Button>
            <Button
              variant="glass"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                onViewDetails(product.id);
              }}
            >
              <Package className="w-4 h-4 mr-2" />
              View Details
            </Button>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}