import type { Metadata, Viewport } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';

const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-inter',
});

export const metadata: Metadata = {
  title: {
    default: 'App | Built with Next.js',
    template: '%s | App',
  },
  description:
    'A modern, animated web application built with Next.js 14, Tailwind CSS, and Framer Motion.',
  keywords: ['Next.js', 'React', 'TypeScript', 'Tailwind CSS', 'Framer Motion'],
  authors: [{ name: 'App Team' }],
  creator: 'App Team',
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_APP_URL ?? 'http://localhost:3000'
  ),
  openGraph: {
    type: 'website',
    locale: 'en_US',
    url: process.env.NEXT_PUBLIC_APP_URL ?? 'http://localhost:3000',
    siteName: 'App',
    title: 'App | Built with Next.js',
    description:
      'A modern, animated web application built with Next.js 14, Tailwind CSS, and Framer Motion.',
    images: [
      {
        url: '/og-image.png',
        width: 1200,
        height: 630,
        alt: 'App Preview',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'App | Built with Next.js',
    description:
      'A modern, animated web application built with Next.js 14, Tailwind CSS, and Framer Motion.',
    images: ['/og-image.png'],
    creator: '@appteam',
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  icons: {
    icon: '/favicon.svg',
    shortcut: '/favicon.svg',
  },
};

export const viewport: Viewport = {
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#ffffff' },
    { media: '(prefers-color-scheme: dark)', color: '#0f172a' },
  ],
  width: 'device-width',
  initialScale: 1,
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${inter.variable} font-sans antialiased bg-white dark:bg-slate-950 text-slate-900 dark:text-slate-50 min-h-screen`}
      >
        {children}
      </body>
    </html>
  );
}
