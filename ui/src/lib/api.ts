import { Product, ProductWithVersions, ProductsResponse } from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || process.env.API_BASE_URL || 'http://localhost:8080';

class APIError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'APIError';
  }
}

async function fetchAPI<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    throw new APIError(response.status, `API Error: ${response.status} ${response.statusText}`);
  }

  return response.json();
}

export async function getProducts(): Promise<Product[]> {
  const response = await fetchAPI<ProductsResponse>('/api/products');
  return response.products;
}

export async function getProduct(id: string): Promise<ProductWithVersions> {
  return fetchAPI<ProductWithVersions>(`/api/products/${id}`);
}

export async function refreshProducts(authToken: string): Promise<{ message: string; jobs_queued: number }> {
  return fetchAPI('/api/refresh', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${authToken}`,
    },
  });
}

export async function getHealthCheck(): Promise<{ status: string; timestamp: string; version: string }> {
  return fetchAPI('/api/health');
}