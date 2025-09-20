export interface Product {
  id: string;
  name: string;
  vendor: string;
  category: string;
  description: string;
  icon_url: string;
  website_url: string;
  created_at: string;
  updated_at: string;
}

export interface ProductVersion {
  id: string;
  product_id: string;
  version: string;
  platform: string;
  architecture: string;
  download_url: string;
  checksum: string;
  checksum_type: string;
  file_size: number;
  filename: string;
  is_latest: boolean;
  etag: string;
  last_fetched: string;
  created_at: string;
  updated_at: string;
}

export interface ProductWithVersions {
  product: Product;
  versions: ProductVersion[];
}

export interface FetchJob {
  id: string;
  product_id: string;
  status: string;
  started_at?: string;
  completed_at?: string;
  error?: string;
  created_at: string;
  updated_at: string;
}

export interface APIResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface ProductsResponse {
  products: Product[];
}

export type Category = 'os' | 'app' | 'tool';
export type Platform = 'windows' | 'linux' | 'macos' | 'web';
export type Architecture = 'amd64' | 'arm64' | '386' | 'arm';