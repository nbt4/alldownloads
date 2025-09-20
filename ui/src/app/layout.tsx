import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'AllDownloads - Official Software Downloads',
  description: 'Get the latest official downloads for operating systems and popular applications',
  keywords: 'software downloads, operating systems, applications, Ubuntu, Windows, Chrome, Firefox',
  authors: [{ name: 'AllDownloads Team' }],
  openGraph: {
    title: 'AllDownloads - Official Software Downloads',
    description: 'Get the latest official downloads for operating systems and popular applications',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'AllDownloads - Official Software Downloads',
    description: 'Get the latest official downloads for operating systems and popular applications',
  },
  viewport: 'width=device-width, initial-scale=1',
  themeColor: '#e11d48',
  manifest: '/manifest.json',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className={inter.className}>
        <div className="min-h-screen bg-gradient-to-br from-background via-background to-primary/5">
          {children}
        </div>
      </body>
    </html>
  );
}